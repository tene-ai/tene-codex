// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package loopcheck

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/tene-ai/tene-codex/internal/domain"
)

type Candidate struct {
	Fingerprint string   `json:"fingerprint"`
	Category    string   `json:"category"`
	Severity    string   `json:"severity"`
	Description string   `json:"description"`
	SubjectRefs []string `json:"subject_refs"`
}

type Options struct {
	ChangedFiles    []string
	DiscoverChanges bool
}
type ReconcileStats struct {
	Created  int `json:"created"`
	Resolved int `json:"resolved"`
	Reopened int `json:"reopened"`
}

func Reconcile(p *domain.Project, sp *domain.Sprint, candidates []Candidate) ReconcileStats {
	stats := ReconcileStats{}
	byFingerprint := map[string]string{}
	for id, g := range p.Gaps {
		if g.SprintID == sp.SprintID && g.DetectedBy == "loop-analyzer" && g.Fingerprint != "" {
			byFingerprint[g.Fingerprint] = id
		}
	}
	active := map[string]bool{}
	for _, candidate := range candidates {
		id := byFingerprint[candidate.Fingerprint]
		if id == "" {
			id = domain.NewID("gap")
			p.Gaps[id] = &domain.Gap{GapID: id, SprintID: sp.SprintID, Category: candidate.Category, Severity: candidate.Severity, Status: "open", Description: candidate.Description, SubjectRefs: candidate.SubjectRefs, Fingerprint: candidate.Fingerprint, DetectedBy: "loop-analyzer", DetectedAtRevision: p.Revision}
			stats.Created++
		} else if p.Gaps[id].Status != "open" {
			p.Gaps[id].Status = "open"
			p.Gaps[id].Resolution = ""
			stats.Reopened++
		}
		active[id] = true
		sp.OpenGapIDs = unique(append(sp.OpenGapIDs, id))
	}
	for id, g := range p.Gaps {
		if g.SprintID == sp.SprintID && g.DetectedBy == "loop-analyzer" && g.Status == "open" && !active[id] {
			g.Status = "resolved"
			g.Resolution = "automatic loop recheck no longer detects the condition"
			sp.OpenGapIDs = remove(sp.OpenGapIDs, id)
			stats.Resolved++
		}
	}
	return stats
}

var contractPattern = regexp.MustCompile(`<!--\s*tene:contract\s+path="([^"]+)"\s+symbol="([^"]+)"\s*-->`)
var forbidPattern = regexp.MustCompile(`<!--\s*tene:forbid\s+path="([^"]+)"\s+contains="([^"]+)"\s*-->`)

func Analyze(ctx context.Context, root string, p *domain.Project, sp *domain.Sprint, options Options) ([]Candidate, error) {
	docs := map[string]string{}
	for key, rel := range map[string]string{"prd": "00-prd/00-prd.md", "plan": "01-plan/00-plan.md", "design": "02-design/00-design.md"} {
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(sp.DocumentRoot), filepath.FromSlash(rel)))
		if err != nil {
			docs[key] = ""
		} else {
			docs[key] = string(b)
		}
	}
	var out []Candidate
	add := func(category, severity, description string, refs ...string) {
		refs = clean(refs)
		sum := sha256.Sum256([]byte(strings.Join(append([]string{category, severity, description}, refs...), "\x00")))
		out = append(out, Candidate{Fingerprint: hex.EncodeToString(sum[:]), Category: category, Severity: severity, Description: description, SubjectRefs: refs})
	}
	for _, intentID := range sp.IntentIDs {
		in := p.Intents[intentID]
		if in == nil {
			add("missing", "blocker", "sprint intent is absent", intentID)
			continue
		}
		// Non-confirmed records remain durable audit history, but only confirmed
		// intent is allowed to drive implementation and verification gates.
		if in.Status != "confirmed" {
			continue
		}
		if !strings.Contains(docs["prd"], intentID) {
			add("mismatch", "blocker", "confirmed intent is not referenced by the PRD", intentID)
		}
		acs := criteriaForIntent(p, intentID)
		if len(acs) == 0 {
			add("missing", "blocker", "confirmed intent has no acceptance criterion", intentID)
		}
		for _, ac := range acs {
			if !strings.Contains(docs["prd"], ac.CriterionID) {
				add("mismatch", "blocker", "acceptance criterion is not referenced by the PRD", ac.CriterionID)
			}
			tasks := tasksForCriterion(p, sp, ac.CriterionID)
			if len(tasks) == 0 {
				add("missing", "blocker", "acceptance criterion has no implementation task", ac.CriterionID)
			}
			if !strings.Contains(docs["design"], ac.CriterionID) {
				add("mismatch", "blocker", "acceptance criterion is not referenced by the design", ac.CriterionID)
			}
		}
	}
	artifactOwners := map[string][]string{}
	for _, taskID := range sp.TaskIDs {
		t := p.Tasks[taskID]
		if t == nil {
			add("missing", "blocker", "sprint references a missing task", taskID)
			continue
		}
		if !strings.Contains(docs["plan"], taskID) {
			add("mismatch", "blocker", "task is not referenced by the plan", taskID)
		}
		if !strings.Contains(docs["design"], taskID) {
			add("mismatch", "blocker", "task is not referenced by the design", taskID)
		}
		if len(t.Artifacts) == 0 {
			add("unverified", "blocker", "task has no linked implementation artifact", taskID)
		}
		for _, path := range t.Artifacts {
			cleanPath := filepath.ToSlash(filepath.Clean(path))
			artifactOwners[cleanPath] = append(artifactOwners[cleanPath], taskID)
			if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(cleanPath))); err != nil {
				add("missing", "blocker", "linked task artifact does not exist", taskID, cleanPath)
			}
		}
	}
	changed := append([]string(nil), options.ChangedFiles...)
	if options.DiscoverChanges {
		var err error
		changed, err = gitChanges(ctx, root)
		if err != nil {
			return nil, err
		}
	}
	for _, path := range changed {
		path = filepath.ToSlash(filepath.Clean(path))
		if relevantChange(path, sp.DocumentRoot) && len(artifactOwners[path]) == 0 {
			add("mismatch", "blocker", "changed implementation artifact is not linked to a Sprint task", path)
		}
	}
	for _, match := range contractPattern.FindAllStringSubmatch(docs["design"], -1) {
		path, symbol := match[1], match[2]
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			add("missing", "blocker", "design contract path does not exist", path, symbol)
		} else if !strings.Contains(string(b), symbol) {
			add("mismatch", "blocker", "design contract symbol is absent from its declared file", path, symbol)
		}
	}
	for _, match := range forbidPattern.FindAllStringSubmatch(docs["design"], -1) {
		path, needle := match[1], match[2]
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			add("missing", "blocker", "forbidden-dependency contract path does not exist", path, needle)
		} else if strings.Contains(string(b), needle) {
			add("mismatch", "blocker", "design forbidden dependency is present", path, needle)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Fingerprint < out[j].Fingerprint })
	return out, nil
}

func gitChanges(ctx context.Context, root string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "status", "--porcelain=v1")
	cmd.Dir = root
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("LOOP_GIT_STATUS: %w", err)
	}
	var out []string
	for _, line := range strings.Split(string(b), "\n") {
		if len(line) < 4 {
			continue
		}
		status := line[:2]
		path := strings.TrimSpace(line[3:])
		if i := strings.Index(path, " -> "); i >= 0 {
			path = path[i+4:]
		}
		// Pure deletion tombstones have no inspectable path and cannot be linked
		// by the current task-artifact schema. Rename destinations and all
		// created/modified files remain subject to ownership checks.
		if strings.Contains(status, "D") && !strings.Contains(line[3:], " -> ") {
			continue
		}
		out = append(out, filepath.ToSlash(path))
	}
	return clean(out), nil
}
func relevantChange(path, documentRoot string) bool {
	if path == "" || strings.HasPrefix(path, ".tene-workflow/") || strings.HasPrefix(path, ".tene/") || strings.HasPrefix(path, documentRoot+"/") || strings.HasPrefix(path, "docs/sprints/_archive/") {
		return false
	}
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".go" || ext == ".py" || ext == ".js" || ext == ".jsx" || ext == ".ts" || ext == ".tsx" || ext == ".json" || ext == ".yaml" || ext == ".yml" || ext == ".sh" || ext == ".md"
}
func criteriaForIntent(p *domain.Project, id string) []*domain.Criterion {
	var out []*domain.Criterion
	for _, x := range p.Criteria {
		if x.IntentID == id {
			out = append(out, x)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CriterionID < out[j].CriterionID })
	return out
}
func tasksForCriterion(p *domain.Project, sp *domain.Sprint, id string) []*domain.Task {
	var out []*domain.Task
	for _, taskID := range sp.TaskIDs {
		if t := p.Tasks[taskID]; t != nil && contains(t.CriterionIDs, id) {
			out = append(out, t)
		}
	}
	return out
}
func contains(values []string, want string) bool {
	for _, x := range values {
		if x == want {
			return true
		}
	}
	return false
}
func clean(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, x := range values {
		x = strings.TrimSpace(x)
		if x != "" && !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	sort.Strings(out)
	return out
}
func unique(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, x := range values {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}
func remove(values []string, want string) []string {
	var out []string
	for _, x := range values {
		if x != want {
			out = append(out, x)
		}
	}
	return out
}
