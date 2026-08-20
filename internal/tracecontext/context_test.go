package tracecontext

import (
	"github.com/tene-ai/tene-codex/internal/domain"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func testTime() time.Time { return time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC) }

func TestContextBudgetPhaseAndFreshness(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".tene-workflow"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".tene-workflow", "project.json"), []byte("state"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".tene-workflow", "policies.yaml"), []byte("policy"), 0644); err != nil {
		t.Fatal(err)
	}
	p := domain.NewProject("project_x", "x", "strict", testTime())
	p.Revision = 7
	in := &domain.Intent{IntentID: "intent_x", Status: "confirmed", Statement: "intent"}
	p.Intents[in.IntentID] = in
	ac := &domain.Criterion{CriterionID: "ac_x", IntentID: in.IntentID, Priority: "blocking", Statement: "works", Observable: "visible"}
	p.Criteria[ac.CriterionID] = ac
	p.Tasks["task_x"] = &domain.Task{TaskID: "task_x", Title: "optional task", CriterionIDs: []string{"ac_x"}}
	sp := &domain.Sprint{SprintID: "sprint_x", Title: "Objective", Phase: domain.PhaseDesign, IntentIDs: []string{in.IntentID}, TaskIDs: []string{"task_x"}}
	pack, err := BuildContextPack(root, p, sp, BuildOptions{Phase: domain.PhasePRD, Budget: 4096}, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range pack.Items {
		if item.Kind == "task" {
			t.Fatal("PRD pack included task")
		}
	}
	if pack.BudgetUnit != "utf8-bytes" || pack.ContentHash == "" {
		t.Fatalf("%#v", pack)
	}
	if fresh := ValidateContextPack(root, p, pack); !fresh.Fresh {
		t.Fatalf("%#v", fresh)
	}
	if err := os.WriteFile(filepath.Join(root, ".tene-workflow", "policies.yaml"), []byte("changed"), 0644); err != nil {
		t.Fatal(err)
	}
	if fresh := ValidateContextPack(root, p, pack); fresh.Fresh || len(fresh.Changed) == 0 {
		t.Fatalf("%#v", fresh)
	}
}

func TestContextMandatoryOverflow(t *testing.T) {
	p := domain.NewProject("p", "x", "strict", testTime())
	p.Revision = 1
	sp := &domain.Sprint{Title: "x", Phase: domain.PhasePRD}
	if _, err := BuildContextPack(t.TempDir(), p, sp, BuildOptions{Budget: 1}, nil); err == nil {
		t.Fatal("expected mandatory overflow")
	}
}
