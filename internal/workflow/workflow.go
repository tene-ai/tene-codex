// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package workflow

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/tene-ai/tene-codex/internal/domain"
)

func ValidProfile(profile string) bool {
	return slices.Contains([]string{"strict", "standard", "light", "off"}, profile)
}

func ParsePhase(value string) (domain.Phase, error) {
	p := domain.Phase(strings.ToLower(value))
	if slices.Contains(domain.PhaseOrder, p) {
		return p, nil
	}
	return "", fmt.Errorf("unknown phase %q", value)
}

func CanTransition(p *domain.Project, sprint *domain.Sprint, target domain.Phase, docExists func(domain.Phase) bool) []domain.Finding {
	return CanTransitionWithApproval(p, sprint, target, "", time.Now().UTC(), docExists)
}

func CanTransitionWithApproval(p *domain.Project, sprint *domain.Sprint, target domain.Phase, approvalID string, now time.Time, docExists func(domain.Phase) bool) []domain.Finding {
	var findings []domain.Finding
	if sprint.Phase == target {
		return findings
	}
	allowed := map[domain.Phase][]domain.Phase{
		domain.PhaseDraft:     {domain.PhasePRD},
		domain.PhasePRD:       {domain.PhasePlan},
		domain.PhasePlan:      {domain.PhaseDesign},
		domain.PhaseDesign:    {domain.PhaseDo},
		domain.PhaseDo:        {domain.PhaseLoopCheck},
		domain.PhaseLoopCheck: {domain.PhaseDo, domain.PhaseQA},
		domain.PhaseQA:        {domain.PhaseDo, domain.PhaseLoopCheck, domain.PhaseReport},
		domain.PhaseReport:    {domain.PhaseLoopCheck, domain.PhaseArchived},
	}
	if !slices.Contains(allowed[sprint.Phase], target) {
		return []domain.Finding{finding("WF_INVALID_TRANSITION", "blocker", sprint.SprintID, fmt.Sprintf("cannot transition from %s to %s", sprint.Phase, target), "Follow the fixed sprint lifecycle.", false)}
	}
	currentDocPhase := sprint.Phase
	if currentDocPhase != domain.PhaseDraft && currentDocPhase != domain.PhaseDo && currentDocPhase != domain.PhaseArchived && !docExists(currentDocPhase) {
		findings = append(findings, finding("DOC_REQUIRED", "blocker", sprint.SprintID, fmt.Sprintf("required %s document is missing or incomplete", currentDocPhase), "Run document scaffold and complete the required sections.", false))
	}
	switch {
	case sprint.Phase == domain.PhasePRD && target == domain.PhasePlan:
		confirmed := 0
		for _, id := range sprint.IntentIDs {
			if in := p.Intents[id]; in != nil && in.Status == "confirmed" {
				confirmed++
			}
		}
		if confirmed == 0 {
			findings = append(findings, finding("INTENT_UNCONFIRMED", "blocker", sprint.SprintID, "no confirmed product intent", "Capture and confirm at least one intent.", false))
		}
		blocking := 0
		for _, ac := range p.Criteria {
			if slices.Contains(sprint.IntentIDs, ac.IntentID) && ac.Priority == "blocking" {
				blocking++
			}
		}
		if blocking == 0 {
			findings = append(findings, finding("AC_MISSING", "blocker", sprint.SprintID, "no blocking acceptance criterion", "Add an observable blocking acceptance criterion.", false))
		}
	case sprint.Phase == domain.PhasePlan && target == domain.PhaseDesign:
		if len(sprint.TaskIDs) == 0 {
			findings = append(findings, finding("TASK_MISSING", "blocker", sprint.SprintID, "implementation plan has no tasks", "Add at least one traceable task.", false))
		}
		for _, id := range sprint.TaskIDs {
			task := p.Tasks[id]
			if task == nil {
				findings = append(findings, finding("TASK_MISSING", "blocker", id, "sprint references a missing task", "Repair the task reference.", false))
				continue
			}
			if len(task.IntentIDs) == 0 && len(task.CriterionIDs) == 0 {
				findings = append(findings, finding("TASK_UNTRACED", "blocker", id, "task is not linked to intent or acceptance criteria", "Add --ac or an intent link.", false))
			}
			for _, dep := range task.DependsOn {
				if p.Tasks[dep] == nil {
					findings = append(findings, finding("TASK_DEPENDENCY_MISSING", "blocker", id, "task dependency does not exist: "+dep, "Repair the dependency ID.", false))
				}
			}
		}
	case sprint.Phase == domain.PhaseDo && target == domain.PhaseLoopCheck:
		for _, id := range sprint.TaskIDs {
			if t := p.Tasks[id]; t != nil && t.Status == "doing" {
				findings = append(findings, finding("TASK_ACTIVE", "blocker", id, "a task is still in progress", "Complete or block the active task.", true))
			}
		}
	case sprint.Phase == domain.PhaseLoopCheck && target == domain.PhaseQA:
		for _, id := range sprint.OpenGapIDs {
			if g := p.Gaps[id]; g != nil && g.Status == "open" && g.Severity == "blocker" && !ActiveWaiver(p, g, time.Now().UTC()) {
				findings = append(findings, finding("GAP_OPEN", "blocker", id, g.Description, "Resolve the blocking gap before QA.", false))
			}
		}
	case sprint.Phase == domain.PhaseQA && target == domain.PhaseReport:
		if sprint.LastQAStatus != "passed" {
			findings = append(findings, finding("QA_NOT_PASSED", "blocker", sprint.SprintID, "the latest QA gate has not passed", "Run and evaluate QA with valid evidence.", false))
		}
	case sprint.Phase == domain.PhaseReport && target == domain.PhaseArchived:
		if sprint.ReportPath == "" || !docExists(domain.PhaseReport) {
			findings = append(findings, finding("REPORT_REQUIRED", "blocker", sprint.SprintID, "validated sprint report is missing", "Generate and validate the report.", false))
		}
	}
	if RequiredApproval(p.Profile, sprint.Phase, target) {
		if code := ApprovalValidity(p, sprint, approvalID, sprint.Phase, target, now); code != "" {
			findings = append(findings, finding(code, "blocker", sprint.SprintID, "a valid human approval is required for this profile boundary", "Request and approve the exact transition, then retry with --approval ID.", false))
		}
	}
	return findings
}

func RequiredApproval(profile string, from, to domain.Phase) bool {
	switch profile {
	case "strict":
		return from == domain.PhaseDesign && to == domain.PhaseDo || from == domain.PhaseReport && to == domain.PhaseArchived
	case "standard":
		return from == domain.PhaseReport && to == domain.PhaseArchived
	default:
		return false
	}
}

func ApprovalValidity(p *domain.Project, sprint *domain.Sprint, id string, from, to domain.Phase, now time.Time) string {
	if id == "" {
		return "WF_APPROVAL_REQUIRED"
	}
	a := p.Approvals[id]
	if a == nil || a.SprintID != sprint.SprintID {
		return "WF_APPROVAL_INVALID"
	}
	if a.From != from || a.To != to {
		return "WF_APPROVAL_SCOPE_MISMATCH"
	}
	if !a.ExpiresAt.After(now) {
		return "WF_APPROVAL_EXPIRED"
	}
	if a.Status == "consumed" || a.ConsumedAt != nil {
		return "WF_APPROVAL_CONSUMED"
	}
	if a.Status != "approved" || a.ApprovedAt == nil || a.Approver == "" {
		return "WF_APPROVAL_NOT_APPROVED"
	}
	return ""
}

func ActiveWaiver(p *domain.Project, gap *domain.Gap, now time.Time) bool {
	if gap == nil || gap.Category == "security" || gap.Category == "evidence-integrity" {
		return false
	}
	for _, w := range p.Waivers {
		if w.GapID == gap.GapID && w.SprintID == gap.SprintID && w.Status == "active" && w.ExpiresAt.After(now) {
			return true
		}
	}
	return false
}

func Blocking(findings []domain.Finding) bool {
	for _, f := range findings {
		if f.Severity == "blocker" {
			return true
		}
	}
	return false
}

func finding(code, severity, ref, message, remediation string, waivable bool) domain.Finding {
	return domain.Finding{Code: code, Severity: severity, SubjectRefs: []string{ref}, Message: message, Remediation: remediation, Waivable: waivable}
}

func EvaluateQAGate(p *domain.Project, sprint *domain.Sprint, run *domain.QARun) []domain.Finding {
	return EvaluateQAGateAtRoot("", p, sprint, run)
}

func EvaluateQAGateAtRoot(root string, p *domain.Project, sprint *domain.Sprint, run *domain.QARun) []domain.Finding {
	var findings []domain.Finding
	if run == nil {
		return []domain.Finding{finding("QA_RUN_MISSING", "blocker", sprint.SprintID, "QA run is missing", "Create a QA plan and attach evidence.", false)}
	}
	currentSpecHash := QASpecHash(p, sprint)
	if run.SpecHash == "" || run.SpecHash != currentSpecHash {
		findings = append(findings, finding("QA_SPEC_STALE", "blocker", run.RunID, "QA plan no longer matches the confirmed intent and acceptance criteria", "Create a new QA plan after the specification change.", false))
	}
	caseByAC := map[string][]domain.QACase{}
	for _, c := range run.Cases {
		for _, ac := range c.CriterionIDs {
			caseByAC[ac] = append(caseByAC[ac], c)
		}
	}
	for _, ac := range p.Criteria {
		if !slices.Contains(sprint.IntentIDs, ac.IntentID) || ac.Priority != "blocking" {
			continue
		}
		cases := caseByAC[ac.CriterionID]
		if len(cases) == 0 {
			findings = append(findings, finding("QA_CHARTER_MISSING", "blocker", ac.CriterionID, "blocking criterion has no QA case", "Regenerate the QA plan.", false))
			continue
		}
		allPassed := true
		for _, c := range cases {
			if len(c.EvidenceIDs) == 0 {
				allPassed = false
				continue
			}
			valid := true
			coveredLayers := map[string]bool{}
			coveredRequirements := map[string]bool{}
			for _, id := range c.EvidenceIDs {
				ev := p.Evidence[id]
				if ev == nil || ev.SHA256 == "" || ev.RedactionStatus != "passed" || !slices.Contains(ev.CriterionIDs, ac.CriterionID) || ev.RunID != run.RunID || ev.CaseID != c.CaseID || ev.SpecHash != run.SpecHash || ev.StateRevision != run.StateRevision || ev.Tool == "" || ev.ToolVersion == "" || ev.Environment == "" || ev.StartedAt == nil || ev.FinishedAt == nil || ev.FinishedAt.Before(*ev.StartedAt) {
					valid = false
					continue
				}
				if root != "" && !evidenceContentValid(root, ev) {
					valid = false
					continue
				}
				for _, assertion := range ev.Assertions {
					if !assertion.Passed || assertion.Layer == "" || !slices.Contains(ev.Layers, assertion.Layer) {
						valid = false
						continue
					}
					coveredLayers[assertion.Layer] = true
					for _, ref := range assertion.RequirementRefs {
						coveredRequirements[ref] = true
					}
				}
			}
			for layer, disposition := range c.RequiredLayers {
				switch {
				case disposition == "required":
					if !coveredLayers[layer] {
						valid = false
					}
				case strings.HasPrefix(disposition, "not-applicable:") || strings.HasPrefix(disposition, "waived:"):
					parts := strings.SplitN(disposition, ":", 3)
					if len(parts) != 3 || strings.TrimSpace(parts[1]) == "" || strings.TrimSpace(parts[2]) == "" {
						valid = false
					}
				default:
					valid = false
				}
			}
			for _, ref := range RequiredQARequirementRefs(c, ac) {
				if !coveredRequirements[ref] {
					valid = false
				}
			}
			allPassed = allPassed && valid
		}
		if !allPassed {
			findings = append(findings, finding("AC_UNVERIFIED", "blocker", ac.CriterionID, "blocking criterion lacks passed, redaction-safe evidence", "Attach valid evidence and mark the QA case passed.", false))
		}
	}
	return findings
}

func QASpecHash(p *domain.Project, sprint *domain.Sprint) string {
	type spec struct {
		Intents  []*domain.Intent    `json:"intents"`
		Criteria []*domain.Criterion `json:"criteria"`
	}
	v := spec{}
	for _, id := range append([]string(nil), sprint.IntentIDs...) {
		if in := p.Intents[id]; in != nil {
			x := *in
			v.Intents = append(v.Intents, &x)
		}
	}
	for _, ac := range p.Criteria {
		if slices.Contains(sprint.IntentIDs, ac.IntentID) {
			x := *ac
			v.Criteria = append(v.Criteria, &x)
		}
	}
	sort.Slice(v.Intents, func(i, j int) bool { return v.Intents[i].IntentID < v.Intents[j].IntentID })
	sort.Slice(v.Criteria, func(i, j int) bool { return v.Criteria[i].CriterionID < v.Criteria[j].CriterionID })
	b, _ := json.Marshal(v)
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func RequiredQARequirementRefs(c domain.QACase, ac *domain.Criterion) []string {
	refs := []string{"observable", "variant:" + c.Variant}
	for i := range ac.Expected {
		refs = append(refs, fmt.Sprintf("expected:%d", i))
	}
	for i := range ac.Forbidden {
		refs = append(refs, fmt.Sprintf("forbidden:%d", i))
	}
	return refs
}

func evidenceContentValid(root string, ev *domain.Evidence) bool {
	path := ev.URI
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, filepath.FromSlash(path))
	}
	b, err := os.ReadFile(path)
	if err != nil || int64(len(b)) != ev.Size {
		return false
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]) == ev.SHA256
}
