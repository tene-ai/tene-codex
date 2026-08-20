// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package workflow

import (
	"testing"
	"time"

	"github.com/tene-ai/tene-codex/internal/domain"
)

func fixture() (*domain.Project, *domain.Sprint) {
	p := domain.NewProject("project_test", "test", "standard", time.Now())
	sp := &domain.Sprint{SprintID: "sprint_test", Phase: domain.PhasePRD, IntentIDs: []string{"intent_test"}}
	p.Sprints[sp.SprintID] = sp
	p.Intents["intent_test"] = &domain.Intent{IntentID: "intent_test", Status: "confirmed"}
	p.Criteria["ac_test"] = &domain.Criterion{CriterionID: "ac_test", IntentID: "intent_test", Priority: "blocking"}
	return p, sp
}

func TestActiveWaiverFailsClosed(t *testing.T) {
	now := time.Now().UTC()
	p := domain.NewProject("project_test", "x", "strict", now)
	g := &domain.Gap{GapID: "gap_test", SprintID: "sprint_test", Category: "mismatch", Severity: "blocker", Status: "open"}
	p.Gaps[g.GapID] = g
	p.Waivers["waiver_test"] = &domain.Waiver{WaiverID: "waiver_test", SprintID: g.SprintID, GapID: g.GapID, Status: "active", ExpiresAt: now.Add(time.Hour)}
	if !ActiveWaiver(p, g, now) {
		t.Fatal("active waiver rejected")
	}
	p.Waivers["waiver_test"].ExpiresAt = now.Add(-time.Second)
	if ActiveWaiver(p, g, now) {
		t.Fatal("expired waiver accepted")
	}
	g.Category = "security"
	p.Waivers["waiver_test"].ExpiresAt = now.Add(time.Hour)
	if ActiveWaiver(p, g, now) {
		t.Fatal("security waiver accepted")
	}
}

func TestFixedLifecycleRejectsSkip(t *testing.T) {
	p, sp := fixture()
	findings := CanTransition(p, sp, domain.PhaseDesign, func(domain.Phase) bool { return true })
	if !Blocking(findings) || findings[0].Code != "WF_INVALID_TRANSITION" {
		t.Fatalf("expected invalid transition: %#v", findings)
	}
}

func TestProfileApprovalMatrixAndValidity(t *testing.T) {
	p, sp := fixture()
	sp.Phase = domain.PhaseDesign
	now := time.Now().UTC()
	for _, tc := range []struct {
		profile  string
		from, to domain.Phase
		required bool
	}{{"strict", domain.PhaseDesign, domain.PhaseDo, true}, {"strict", domain.PhaseReport, domain.PhaseArchived, true}, {"standard", domain.PhaseDesign, domain.PhaseDo, false}, {"standard", domain.PhaseReport, domain.PhaseArchived, true}, {"light", domain.PhaseReport, domain.PhaseArchived, false}, {"off", domain.PhaseDesign, domain.PhaseDo, false}} {
		if got := RequiredApproval(tc.profile, tc.from, tc.to); got != tc.required {
			t.Fatalf("%#v got %v", tc, got)
		}
	}
	approved := now
	a := &domain.Approval{ApprovalID: "approval_x", SprintID: sp.SprintID, From: domain.PhaseDesign, To: domain.PhaseDo, Status: "approved", Approver: "human", ApprovedAt: &approved, ExpiresAt: now.Add(time.Hour)}
	p.Approvals[a.ApprovalID] = a
	if code := ApprovalValidity(p, sp, a.ApprovalID, a.From, a.To, now); code != "" {
		t.Fatal(code)
	}
	if code := ApprovalValidity(p, sp, a.ApprovalID, a.From, domain.PhaseQA, now); code != "WF_APPROVAL_SCOPE_MISMATCH" {
		t.Fatal(code)
	}
	a.ExpiresAt = now.Add(-time.Second)
	if code := ApprovalValidity(p, sp, a.ApprovalID, a.From, a.To, now); code != "WF_APPROVAL_EXPIRED" {
		t.Fatal(code)
	}
}

func TestPRDGateRequiresConfirmedIntentAndCriterion(t *testing.T) {
	p, sp := fixture()
	if got := CanTransition(p, sp, domain.PhasePlan, func(domain.Phase) bool { return true }); Blocking(got) {
		t.Fatalf("unexpected findings: %#v", got)
	}
	p.Intents["intent_test"].Status = "candidate"
	if got := CanTransition(p, sp, domain.PhasePlan, func(domain.Phase) bool { return true }); !Blocking(got) {
		t.Fatal("expected blocker")
	}
}

func TestQAGateRequiresEveryVariantAndMatchingEvidence(t *testing.T) {
	p, sp := fixture()
	sp.Phase = domain.PhaseQA
	now := time.Now().UTC()
	run := &domain.QARun{RunID: "run", StateRevision: p.Revision, SpecHash: QASpecHash(p, sp), Cases: []domain.QACase{{CaseID: "happy", Variant: "happy", CriterionIDs: []string{"ac_test"}, RequiredLayers: map[string]string{"L5": "required"}, EvidenceIDs: []string{"ev-happy"}}, {CaseID: "recovery", Variant: "recovery", CriterionIDs: []string{"ac_test"}, RequiredLayers: map[string]string{"L5": "required"}}}}
	p.Evidence["ev-happy"] = validEvidence("ev-happy", run, run.Cases[0], now)
	if got := EvaluateQAGate(p, sp, run); !Blocking(got) {
		t.Fatal("pending recovery case must block")
	}
	run.Cases[1].RequiredLayers = map[string]string{"L5": "not-applicable:independent-evaluator:no recovery boundary"}
	if got := EvaluateQAGate(p, sp, run); Blocking(got) {
		t.Fatalf("fully not-applicable case should not require synthetic evidence: %#v", got)
	}
	run.Cases[1].RequiredLayers = map[string]string{"L5": "required"}
	run.Cases[1].EvidenceIDs = []string{"ev-recovery"}
	p.Evidence["ev-recovery"] = validEvidence("ev-recovery", run, run.Cases[1], now)
	if got := EvaluateQAGate(p, sp, run); Blocking(got) {
		t.Fatalf("unexpected findings: %#v", got)
	}
	mutations := []struct {
		name   string
		mutate func(*domain.Evidence)
	}{
		{"wrong-run", func(e *domain.Evidence) { e.RunID = "other" }},
		{"wrong-case", func(e *domain.Evidence) { e.CaseID = "other" }},
		{"stale-spec", func(e *domain.Evidence) { e.SpecHash = "old" }},
		{"no-layer", func(e *domain.Evidence) { e.Layers = nil }},
		{"failed-assertion", func(e *domain.Evidence) { e.Assertions[0].Passed = false }},
		{"no-tool-version", func(e *domain.Evidence) { e.ToolVersion = "" }},
		{"redaction-failed", func(e *domain.Evidence) { e.RedactionStatus = "failed" }},
	}
	for _, tc := range mutations {
		t.Run(tc.name, func(t *testing.T) {
			original := *p.Evidence["ev-happy"]
			original.Layers = append([]string(nil), original.Layers...)
			original.Assertions = append([]domain.EvidenceAssertion(nil), original.Assertions...)
			tc.mutate(&original)
			p.Evidence["ev-happy"] = &original
			if got := EvaluateQAGate(p, sp, run); !Blocking(got) {
				t.Fatalf("mutation passed: %#v", original)
			}
			p.Evidence["ev-happy"] = validEvidence("ev-happy", run, run.Cases[0], now)
		})
	}
}

func validEvidence(id string, run *domain.QARun, c domain.QACase, now time.Time) *domain.Evidence {
	return &domain.Evidence{EvidenceID: id, RunID: run.RunID, CaseID: c.CaseID, SpecHash: run.SpecHash, StateRevision: run.StateRevision, SHA256: "hash", Size: 1, RedactionStatus: "passed", CriterionIDs: append([]string(nil), c.CriterionIDs...), Layers: []string{"L5"}, Assertions: []domain.EvidenceAssertion{{Statement: "observed", Passed: true, Layer: "L5", RequirementRefs: []string{"observable", "variant:" + c.Variant}}}, Tool: "observer", ToolVersion: "1", Environment: "local", StartedAt: &now, FinishedAt: &now}
}
