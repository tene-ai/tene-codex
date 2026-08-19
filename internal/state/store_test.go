// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/tene-ai/tene-codex/internal/domain"
)

func TestMigrateLegacyAndRepairDerived(t *testing.T) {
	root := t.TempDir()
	s := New(root)
	p := domain.NewProject("project_test", "fixture", "strict", time.Now())
	if err := s.Initialize(p); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(s.ProjectPath())
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatal(err)
	}
	raw["schema_version"] = "0.9.0"
	delete(raw, "waivers")
	b, _ = json.Marshal(raw)
	if err := os.WriteFile(s.ProjectPath(), b, 0644); err != nil {
		t.Fatal(err)
	}
	plan, err := s.PlanMigration()
	if err != nil || !plan.Required || !plan.Supported {
		t.Fatalf("%#v %v", plan, err)
	}
	applied, err := s.Migrate()
	if err != nil || applied.Backup == "" {
		t.Fatalf("%#v %v", applied, err)
	}
	got, err := s.Load()
	if err != nil || got.SchemaVersion != domain.SchemaVersion || got.Waivers == nil {
		t.Fatalf("%#v %v", got, err)
	}
	os.Remove(s.ActivePath())
	os.Remove(s.MasterPlanPath())
	paths, err := s.RepairDerived()
	if err != nil || len(paths) != 2 {
		t.Fatalf("%v %v", paths, err)
	}
}

func TestLoadRejectsUnknownPersistedField(t *testing.T) {
	root := t.TempDir()
	s := New(root)
	p := domain.NewProject("project_test", "fixture", "strict", time.Now())
	if err := s.Initialize(p); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(s.ProjectPath())
	var raw map[string]any
	json.Unmarshal(b, &raw)
	raw["mystery"] = true
	b, _ = json.Marshal(raw)
	os.WriteFile(s.ProjectPath(), b, 0644)
	if _, err := s.Load(); err == nil {
		t.Fatal("unknown field accepted")
	}
}

func TestStoreInitializeMutateAndVerify(t *testing.T) {
	root := t.TempDir()
	s := New(root)
	p := domain.NewProject("project_test", "test", "standard", time.Unix(1, 0))
	if err := s.Initialize(p); err != nil {
		t.Fatal(err)
	}
	expected := uint64(0)
	updated, err := s.Mutate(&expected, domain.Actor{Kind: "test"}, "TestChanged", p.ProjectID, map[string]string{"value": "ok"}, func(p *domain.Project) error { p.Name = "changed"; return nil })
	if err != nil {
		t.Fatal(err)
	}
	if updated.Revision != 1 || updated.Name != "changed" {
		t.Fatalf("unexpected state: %#v", updated)
	}
	events, err := s.VerifyJournal()
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 || events[1].PreviousHash != events[0].Hash {
		t.Fatalf("invalid event chain: %#v", events)
	}
}

func TestStoreRejectsStaleRevision(t *testing.T) {
	s := New(t.TempDir())
	p := domain.NewProject("project_test", "test", "standard", time.Now())
	if err := s.Initialize(p); err != nil {
		t.Fatal(err)
	}
	stale := uint64(2)
	_, err := s.Mutate(&stale, domain.Actor{Kind: "test"}, "Noop", p.ProjectID, nil, func(*domain.Project) error { return nil })
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestCompactAndClearPreserveSource(t *testing.T) {
	s := New(t.TempDir())
	p := domain.NewProject("project_test", "test", "standard", time.Now())
	if err := s.Initialize(p); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Compact(); err != nil {
		t.Fatal(err)
	}
	if err := s.ClearDerived(); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Load(); err != nil {
		t.Fatal(err)
	}
	if _, err := s.VerifyJournal(); err != nil {
		t.Fatal(err)
	}
}
