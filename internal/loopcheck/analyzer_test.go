// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package loopcheck

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tene-ai/tene-codex/internal/domain"
)

func TestAnalyzeMutationDetectionAndCleanConvergence(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(string, *domain.Project, *domain.Sprint)
	}{
		{"prd-intent-missing", func(root string, p *domain.Project, sp *domain.Sprint) {
			write(t, filepath.Join(root, sp.DocumentRoot, "00-prd/00-prd.md"), "ac_x")
		}},
		{"ac-task-missing", func(root string, p *domain.Project, sp *domain.Sprint) { p.Tasks["task_x"].CriterionIDs = nil }},
		{"changed-file-untraced", func(root string, p *domain.Project, sp *domain.Sprint) {}},
		{"contract-symbol-missing", func(root string, p *domain.Project, sp *domain.Sprint) {
			write(t, filepath.Join(root, "internal/service.go"), "package internal")
		}},
		{"forbidden-dependency", func(root string, p *domain.Project, sp *domain.Sprint) {
			write(t, filepath.Join(root, "internal/service.go"), "package internal\n// Service\n// direct-db")
		}},
		{"linked-artifact-missing", func(root string, p *domain.Project, sp *domain.Sprint) {
			os.Remove(filepath.Join(root, "internal/service.go"))
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root, p, sp := fixture(t)
			tc.mutate(root, p, sp)
			changed := []string{}
			if tc.name == "changed-file-untraced" {
				write(t, filepath.Join(root, "internal/extra.go"), "package internal")
				changed = []string{"internal/extra.go"}
			}
			got, err := Analyze(context.Background(), root, p, sp, Options{ChangedFiles: changed})
			if err != nil || len(got) == 0 {
				t.Fatalf("mutation escaped: %#v %v", got, err)
			}
		})
	}
	root, p, sp := fixture(t)
	got, err := Analyze(context.Background(), root, p, sp, Options{})
	if err != nil || len(got) != 0 {
		t.Fatalf("clean fixture did not converge: %#v %v", got, err)
	}
}

func TestReconcileIsStableResolvesAndReopens(t *testing.T) {
	_, p, sp := fixture(t)
	candidate := Candidate{Fingerprint: "stable", Category: "mismatch", Severity: "blocker", Description: "drift", SubjectRefs: []string{"x"}}
	first := Reconcile(p, sp, []Candidate{candidate})
	if first.Created != 1 || len(sp.OpenGapIDs) != 1 {
		t.Fatalf("first %#v %#v", first, sp.OpenGapIDs)
	}
	id := sp.OpenGapIDs[0]
	second := Reconcile(p, sp, []Candidate{candidate})
	if second.Created != 0 || len(sp.OpenGapIDs) != 1 || sp.OpenGapIDs[0] != id {
		t.Fatalf("duplicate %#v %#v", second, sp.OpenGapIDs)
	}
	third := Reconcile(p, sp, nil)
	if third.Resolved != 1 || p.Gaps[id].Status != "resolved" || len(sp.OpenGapIDs) != 0 {
		t.Fatalf("resolve %#v %#v", third, p.Gaps[id])
	}
	fourth := Reconcile(p, sp, []Candidate{candidate})
	if fourth.Reopened != 1 || p.Gaps[id].Status != "open" || sp.OpenGapIDs[0] != id {
		t.Fatalf("reopen %#v %#v", fourth, p.Gaps[id])
	}
}

func fixture(t *testing.T) (string, *domain.Project, *domain.Sprint) {
	t.Helper()
	root := t.TempDir()
	p := domain.NewProject("project_x", "x", "strict", time.Now())
	sp := &domain.Sprint{SprintID: "sprint_x", IntentIDs: []string{"intent_x"}, TaskIDs: []string{"task_x"}, DocumentRoot: "docs/sprints/sprint_x"}
	p.Sprints[sp.SprintID] = sp
	p.Intents["intent_x"] = &domain.Intent{IntentID: "intent_x", Status: "confirmed"}
	p.Criteria["ac_x"] = &domain.Criterion{CriterionID: "ac_x", IntentID: "intent_x", Priority: "blocking"}
	p.Tasks["task_x"] = &domain.Task{TaskID: "task_x", SprintID: sp.SprintID, CriterionIDs: []string{"ac_x"}, Artifacts: []string{"internal/service.go"}}
	write(t, filepath.Join(root, sp.DocumentRoot, "00-prd/00-prd.md"), "intent_x ac_x")
	write(t, filepath.Join(root, sp.DocumentRoot, "01-plan/00-plan.md"), "task_x ac_x")
	write(t, filepath.Join(root, sp.DocumentRoot, "02-design/00-design.md"), "ac_x task_x\n<!-- tene:contract path=\"internal/service.go\" symbol=\"Service\" -->\n<!-- tene:forbid path=\"internal/service.go\" contains=\"direct-db\" -->")
	write(t, filepath.Join(root, "internal/service.go"), "package internal\n// Service")
	return root, p, sp
}
func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
