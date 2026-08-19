// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package workflow

import (
	"fmt"
	"slices"
	"strings"

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
		domain.PhaseReport:    {domain.PhaseArchived},
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
			if g := p.Gaps[id]; g != nil && g.Status == "open" && g.Severity == "blocker" {
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
	return findings
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
	var findings []domain.Finding
	if run == nil {
		return []domain.Finding{finding("QA_RUN_MISSING", "blocker", sprint.SprintID, "QA run is missing", "Create a QA plan and attach evidence.", false)}
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
			if c.Status != "passed" || len(c.EvidenceIDs) == 0 {
				allPassed = false
				continue
			}
			valid := true
			for _, id := range c.EvidenceIDs {
				ev := p.Evidence[id]
				if ev == nil || ev.SHA256 == "" || ev.RedactionStatus != "passed" || !slices.Contains(ev.CriterionIDs, ac.CriterionID) {
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
