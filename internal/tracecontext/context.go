// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package tracecontext

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tene-ai/tene-codex/internal/codeintel"
	"github.com/tene-ai/tene-codex/internal/domain"
)

const ContextSchemaVersion = "1.0.0"

type ContextItem struct {
	Kind           string `json:"kind"`
	Category       string `json:"category"`
	RefID          string `json:"ref_id"`
	Priority       int    `json:"priority"`
	Mandatory      bool   `json:"mandatory"`
	Content        string `json:"content"`
	Locator        string `json:"locator,omitempty"`
	SourceRevision uint64 `json:"source_revision"`
	ContentHash    string `json:"content_hash"`
	Size           int    `json:"size"`
}
type SourceRef struct {
	Locator  string `json:"locator"`
	Revision uint64 `json:"revision"`
	SHA256   string `json:"sha256,omitempty"`
	Status   string `json:"status"`
}
type ExcludedSummary struct {
	Reason string `json:"reason"`
	Count  int    `json:"count"`
	Bytes  int    `json:"bytes"`
}
type ContextPack struct {
	SchemaVersion    string                 `json:"schema_version"`
	ID               string                 `json:"id"`
	Phase            domain.Phase           `json:"phase"`
	StateRevision    uint64                 `json:"state_revision"`
	Budget           int                    `json:"budget"`
	BudgetUnit       string                 `json:"budget_unit"`
	Used             int                    `json:"used"`
	Objective        string                 `json:"objective"`
	Items            []ContextItem          `json:"items"`
	ToolCapabilities []codeintel.Capability `json:"tool_capabilities"`
	Provenance       []SourceRef            `json:"provenance"`
	ExcludedSummary  []ExcludedSummary      `json:"excluded_summary"`
	ContentHash      string                 `json:"content_hash"`
	BudgetAllocation map[string]int         `json:"budget_allocation"`
	CategoryUsage    map[string]int         `json:"category_usage"`
}
type BuildOptions struct {
	Phase  domain.Phase
	Budget int
}
type FreshnessResult struct {
	Fresh                 bool        `json:"fresh"`
	StateRevisionExpected uint64      `json:"state_revision_expected"`
	StateRevisionActual   uint64      `json:"state_revision_actual"`
	Changed               []SourceRef `json:"changed"`
	Missing               []SourceRef `json:"missing"`
}

func BuildContextPack(root string, p *domain.Project, sp *domain.Sprint, options BuildOptions, caps []codeintel.Capability) (ContextPack, error) {
	if options.Budget <= 0 {
		return ContextPack{}, fmt.Errorf("CONTEXT_BUDGET_INVALID")
	}
	phase := options.Phase
	if phase == "" {
		phase = sp.Phase
	}
	if !validPhase(phase) {
		return ContextPack{}, fmt.Errorf("CONTEXT_PHASE_INVALID: %s", phase)
	}
	var candidates []ContextItem
	add := func(category, kind, id, content, locator string, priority int, mandatory bool) {
		content = strings.TrimSpace(content)
		if content == "" {
			return
		}
		sum := sha256.Sum256([]byte(content))
		candidates = append(candidates, ContextItem{Category: category, Kind: kind, RefID: id, Priority: priority, Mandatory: mandatory, Content: content, Locator: locator, SourceRevision: p.Revision, ContentHash: hex.EncodeToString(sum[:]), Size: len([]byte(content))})
	}
	masterPolicy, _ := json.Marshal(map[string]any{"profile": p.Profile, "objective": p.MasterPlan.Objective, "invariants": p.MasterPlan.CrossSprintInvariants})
	add("policy", "policy", "workflow-profile", string(masterPolicy), ".tene-workflow/project.json", 100, true)
	for _, id := range sortedStrings(sp.IntentIDs) {
		if in := p.Intents[id]; in != nil && in.Status == "confirmed" {
			b, _ := json.Marshal(in)
			add("spec", "intent", id, string(b), ".tene-workflow/project.json", 90, true)
		}
	}
	for _, id := range sortedCriteria(p, sp) {
		ac := p.Criteria[id]
		b, _ := json.Marshal(ac)
		add("spec", "acceptance-criterion", id, string(b), ".tene-workflow/project.json", 85, ac.Priority == "blocking")
	}
	for _, id := range sortedStrings(sp.OpenGapIDs) {
		if gap := p.Gaps[id]; gap != nil && gap.Status == "open" {
			b, _ := json.Marshal(gap)
			add("evidence", "gap", id, string(b), ".tene-workflow/project.json", 80, gap.Severity == "blocker")
		}
	}
	if phaseAllowsTasks(phase) {
		for _, id := range sortedStrings(sp.TaskIDs) {
			if task := p.Tasks[id]; task != nil {
				b, _ := json.Marshal(task)
				add("work", "task", id, string(b), ".tene-workflow/project.json", 60, false)
			}
		}
		if locator := phaseDocument(sp, phase); locator != "" {
			if b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(locator))); err == nil {
				add("work", "phase-document", string(phase), string(b), locator, 70, false)
			}
		}
	}
	if phaseAllowsGraph(phase) {
		for _, id := range sortedKeys(p.Graph.Nodes) {
			node := p.Graph.Nodes[id]
			if node.Kind == "File" || node.Kind == "Symbol" {
				b, _ := json.Marshal(node)
				add("graph", "graph-node", id, string(b), node.Locator, 40, false)
			}
		}
	}
	for id, approval := range p.Approvals {
		if approval.SprintID == sp.SprintID {
			b, _ := json.Marshal(approval)
			add("evidence", "decision", id, string(b), ".tene-workflow/project.json", 55, false)
		}
	}
	for id, waiver := range p.Waivers {
		if waiver.SprintID == sp.SprintID {
			b, _ := json.Marshal(waiver)
			add("evidence", "decision", id, string(b), ".tene-workflow/project.json", 55, false)
		}
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].Priority == candidates[j].Priority {
			if candidates[i].Kind == candidates[j].Kind {
				return candidates[i].RefID < candidates[j].RefID
			}
			return candidates[i].Kind < candidates[j].Kind
		}
		return candidates[i].Priority > candidates[j].Priority
	})
	allocations := map[string]int{"policy": options.Budget * 15 / 100, "spec": options.Budget * 25 / 100, "work": options.Budget * 20 / 100, "graph": options.Budget * 25 / 100, "evidence": options.Budget * 10 / 100, "reserve": options.Budget * 5 / 100}
	pack := ContextPack{SchemaVersion: ContextSchemaVersion, ID: contextID(p.Revision, phase), Phase: phase, StateRevision: p.Revision, Budget: options.Budget, BudgetUnit: "utf8-bytes", Objective: sp.Title, ToolCapabilities: caps, Items: []ContextItem{}, Provenance: []SourceRef{}, ExcludedSummary: []ExcludedSummary{}, BudgetAllocation: allocations, CategoryUsage: map[string]int{}}
	seen := map[string]bool{}
	excluded := ExcludedSummary{Reason: "budget", Count: 0, Bytes: 0}
	for _, item := range candidates {
		if seen[item.ContentHash] {
			continue
		}
		seen[item.ContentHash] = true
		if pack.Used+item.Size > pack.Budget || (!item.Mandatory && pack.CategoryUsage[item.Category]+item.Size > allocations[item.Category]) {
			if item.Mandatory {
				return ContextPack{}, fmt.Errorf("CONTEXT_MANDATORY_OVERFLOW: %s requires %d bytes with %d remaining", item.RefID, item.Size, pack.Budget-pack.Used)
			}
			excluded.Count++
			excluded.Bytes += item.Size
			continue
		}
		pack.Items = append(pack.Items, item)
		pack.Used += item.Size
		pack.CategoryUsage[item.Category] += item.Size
	}
	if excluded.Count > 0 {
		pack.ExcludedSummary = append(pack.ExcludedSummary, excluded)
	}
	locators := map[string]bool{".tene-workflow/project.json": true, ".tene-workflow/policies.yaml": true}
	if doc := phaseDocument(sp, phase); doc != "" {
		locators[doc] = true
	}
	for _, item := range pack.Items {
		if item.Locator != "" && safeProvenance(item.Locator) {
			locators[item.Locator] = true
		}
	}
	for _, locator := range sortedKeys(locators) {
		ref := SourceRef{Locator: locator, Revision: p.Revision, Status: "recorded"}
		path := filepath.Join(root, filepath.FromSlash(provenanceFile(locator)))
		if b, err := os.ReadFile(path); err == nil {
			sum := sha256.Sum256(b)
			ref.SHA256 = hex.EncodeToString(sum[:])
		} else {
			ref.Status = "missing"
		}
		pack.Provenance = append(pack.Provenance, ref)
	}
	pack.ContentHash = packHash(pack)
	return pack, nil
}

func ValidateContextPack(root string, p *domain.Project, pack ContextPack) FreshnessResult {
	result := FreshnessResult{Fresh: pack.StateRevision == p.Revision, StateRevisionExpected: pack.StateRevision, StateRevisionActual: p.Revision, Changed: []SourceRef{}, Missing: []SourceRef{}}
	for _, ref := range pack.Provenance {
		if ref.SHA256 == "" {
			continue
		}
		path := filepath.Join(root, filepath.FromSlash(provenanceFile(ref.Locator)))
		b, err := os.ReadFile(path)
		if err != nil {
			result.Missing = append(result.Missing, ref)
			result.Fresh = false
			continue
		}
		sum := sha256.Sum256(b)
		actual := hex.EncodeToString(sum[:])
		if actual != ref.SHA256 {
			ref.Status = "changed:" + actual
			result.Changed = append(result.Changed, ref)
			result.Fresh = false
		}
	}
	return result
}

func phaseAllowsTasks(p domain.Phase) bool {
	return p == domain.PhasePlan || p == domain.PhaseDesign || p == domain.PhaseDo || p == domain.PhaseLoopCheck || p == domain.PhaseQA || p == domain.PhaseReport
}
func phaseAllowsGraph(p domain.Phase) bool {
	return p == domain.PhaseDesign || p == domain.PhaseDo || p == domain.PhaseLoopCheck || p == domain.PhaseQA || p == domain.PhaseReport
}
func validPhase(p domain.Phase) bool {
	for _, x := range domain.PhaseOrder {
		if p == x {
			return true
		}
	}
	return false
}
func phaseDocument(sp *domain.Sprint, p domain.Phase) string {
	names := map[domain.Phase]string{domain.PhasePRD: "00-prd/00-prd.md", domain.PhasePlan: "01-plan/00-plan.md", domain.PhaseDesign: "02-design/00-design.md", domain.PhaseDo: "02-design/00-design.md", domain.PhaseLoopCheck: "03-analysis/00-loop-check.md", domain.PhaseQA: "04-qa/00-qa-plan.md", domain.PhaseReport: "05-report/00-report.md"}
	if n := names[p]; n != "" {
		return filepath.ToSlash(filepath.Join(sp.DocumentRoot, n))
	}
	return ""
}
func sortedStrings(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}
func sortedCriteria(p *domain.Project, sp *domain.Sprint) []string {
	var out []string
	for _, ac := range domain.ConfirmedCriteria(p, sp) {
		out = append(out, ac.CriterionID)
	}
	return out
}
func safeProvenance(locator string) bool {
	clean := filepath.ToSlash(filepath.Clean(locator))
	return clean != ".tene" && !strings.HasPrefix(clean, ".tene/") && clean != ".." && !strings.HasPrefix(clean, "../") && !filepath.IsAbs(clean)
}
func provenanceFile(locator string) string {
	for i := len(locator) - 1; i >= 0; i-- {
		if locator[i] != ':' {
			continue
		}
		line := locator[i+1:]
		if line == "" {
			break
		}
		digits := true
		for _, r := range line {
			if r < '0' || r > '9' {
				digits = false
				break
			}
		}
		if digits {
			return locator[:i]
		}
		break
	}
	return locator
}
func contextID(rev uint64, p domain.Phase) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%d:%s", rev, p)))
	return "context_" + hex.EncodeToString(sum[:8])
}
func packHash(pack ContextPack) string {
	clone := pack
	clone.ContentHash = ""
	b, _ := json.Marshal(clone)
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
