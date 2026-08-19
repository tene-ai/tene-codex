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

func TestFixedLifecycleRejectsSkip(t *testing.T) {
	p, sp := fixture()
	findings := CanTransition(p, sp, domain.PhaseDesign, func(domain.Phase) bool { return true })
	if !Blocking(findings) || findings[0].Code != "WF_INVALID_TRANSITION" {
		t.Fatalf("expected invalid transition: %#v", findings)
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
	p.Evidence["ev"] = &domain.Evidence{EvidenceID: "ev", SHA256: "hash", RedactionStatus: "passed", CriterionIDs: []string{"ac_test"}}
	run := &domain.QARun{Cases: []domain.QACase{{CaseID: "happy", CriterionIDs: []string{"ac_test"}, Status: "passed", EvidenceIDs: []string{"ev"}}, {CaseID: "recovery", CriterionIDs: []string{"ac_test"}, Status: "pending"}}}
	if got := EvaluateQAGate(p, sp, run); !Blocking(got) {
		t.Fatal("pending recovery case must block")
	}
	run.Cases[1].Status = "passed"
	run.Cases[1].EvidenceIDs = []string{"ev"}
	if got := EvaluateQAGate(p, sp, run); Blocking(got) {
		t.Fatalf("unexpected findings: %#v", got)
	}
}
