// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package document

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tene-ai/tene-codex/internal/domain"
)

type Spec struct {
	Dir, File, Type string
	Phase           domain.Phase
	Extra           []string
}

var specs = map[domain.Phase]Spec{
	domain.PhasePRD:       {"00-prd", "00-prd.md", "prd", domain.PhasePRD, []string{"problem", "actors", "journeys", "acceptance-criteria", "non-goals"}},
	domain.PhasePlan:      {"01-plan", "00-plan.md", "plan", domain.PhasePlan, []string{"work-packages", "dependencies", "verification", "risks", "yagni"}},
	domain.PhaseDesign:    {"02-design", "00-design.md", "design", domain.PhaseDesign, []string{"components", "interfaces", "data", "state-transitions", "failures", "security", "tests"}},
	domain.PhaseLoopCheck: {"03-analysis", "00-loop-check.md", "analysis", domain.PhaseLoopCheck, []string{"baseline", "changed-artifacts", "gap-matrix", "iterations", "regression"}},
	domain.PhaseQA:        {"04-qa", "00-qa-plan.md", "qa", domain.PhaseQA, []string{"environment", "capabilities", "charters", "ux-data-flow", "evidence", "verdict"}},
	domain.PhaseReport:    {"05-report", "00-report.md", "report", domain.PhaseReport, []string{"previous-sprints", "changed-files", "intent-fulfillment", "qa-verdict", "deferred-work", "next-sprint"}},
}

var common = []string{"purpose", "scope", "layers", "six-questions", "traceability", "decisions", "freeform"}

func SprintRoot(root string, sprint *domain.Sprint) string {
	return filepath.Join(root, sprint.DocumentRoot)
}

func Path(root string, sprint *domain.Sprint, phase domain.Phase) string {
	spec, ok := specs[phase]
	if !ok {
		return ""
	}
	return filepath.Join(SprintRoot(root, sprint), spec.Dir, spec.File)
}

func Exists(root string, sprint *domain.Sprint, phase domain.Phase) bool {
	path := Path(root, sprint, phase)
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func ScaffoldAll(root string, project *domain.Project, sprint *domain.Sprint) ([]string, error) {
	var paths []string
	for _, phase := range []domain.Phase{domain.PhasePRD, domain.PhasePlan, domain.PhaseDesign, domain.PhaseLoopCheck, domain.PhaseQA, domain.PhaseReport} {
		path, created, err := Scaffold(root, project, sprint, phase)
		if err != nil {
			return paths, err
		}
		if created {
			paths = append(paths, path)
		}
	}
	if err := os.MkdirAll(filepath.Join(SprintRoot(root, sprint), "99-archive"), 0o755); err != nil {
		return paths, err
	}
	return paths, nil
}

func Scaffold(root string, project *domain.Project, sprint *domain.Sprint, phase domain.Phase) (string, bool, error) {
	spec, ok := specs[phase]
	if !ok {
		return "", false, fmt.Errorf("no document template for phase %s", phase)
	}
	path := Path(root, sprint, phase)
	if _, err := os.Stat(path); err == nil {
		return path, false, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return path, false, err
	}
	var b strings.Builder
	fmt.Fprintf(&b, "---\nschema_version: %s\ndocument_type: %s\nproject_id: %s\nsprint_id: %s\nphase: %s\nstatus: draft\nrevision: %d\nintent_ids: []\ngenerated_at: %s\ngenerated_by: tene-workflow\n---\n\n# %s — %s\n\n", domain.SchemaVersion, spec.Type, project.ProjectID, sprint.SprintID, phase, project.Revision, time.Now().UTC().Format(time.RFC3339), phase, sprint.Title)
	for _, id := range append(common, spec.Extra...) {
		fmt.Fprintf(&b, "<!-- tene:section:%s -->\n## %s\n\n", id, title(id))
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		return path, false, err
	}
	return path, true, nil
}

func Validate(path string, phase domain.Phase) []domain.Finding {
	b, err := os.ReadFile(path)
	if err != nil {
		return []domain.Finding{{Code: "DOC_MISSING", Severity: "blocker", Message: err.Error(), Remediation: "Scaffold the document."}}
	}
	text := string(b)
	if !strings.HasPrefix(text, "---\n") {
		return []domain.Finding{{Code: "DOC_FRONTMATTER_MISSING", Severity: "blocker", Message: "YAML frontmatter is missing", Remediation: "Regenerate or repair the document frontmatter."}}
	}
	spec := specs[phase]
	var out []domain.Finding
	for _, id := range append(common, spec.Extra...) {
		marker := "<!-- tene:section:" + id + " -->"
		if !strings.Contains(text, marker) {
			out = append(out, domain.Finding{Code: "DOC_SECTION_MISSING", Severity: "blocker", SubjectRefs: []string{path}, Message: "missing section " + id, Remediation: "Add the stable section marker and content."})
			continue
		}
		if id != "freeform" && !sectionHasContent(text, marker) {
			out = append(out, domain.Finding{Code: "DOC_SECTION_EMPTY", Severity: "blocker", SubjectRefs: []string{path}, Message: "section has no authored content: " + id, Remediation: "Document the decision or write N/A with an evidence-based reason."})
		}
	}
	return out
}

func RequiredMarkers(phase domain.Phase) []string {
	return append(append([]string{}, common...), specs[phase].Extra...)
}

func title(id string) string {
	id = strings.ReplaceAll(id, "-", " ")
	return strings.ToUpper(id[:1]) + id[1:]
}

func sectionHasContent(text, marker string) bool {
	start := strings.Index(text, marker)
	if start < 0 {
		return false
	}
	segment := text[start+len(marker):]
	if end := strings.Index(segment, "<!-- tene:section:"); end >= 0 {
		segment = segment[:end]
	}
	for _, line := range strings.Split(segment, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "<!--") || strings.HasSuffix(line, "-->") {
			continue
		}
		return true
	}
	return false
}
