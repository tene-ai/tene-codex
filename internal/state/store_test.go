// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
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
	if _, _, err := s.CreateCheckpoint(); err != nil {
		t.Fatal(err)
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

func TestMigrationSupportWindowFailsClosedWithoutMutation(t *testing.T) {
	for _, version := range []string{"0.8.0", "2.0.0"} {
		t.Run(version, func(t *testing.T) {
			root := t.TempDir()
			s := New(root)
			p := domain.NewProject("project_test", "fixture", "strict", time.Now())
			if err := s.Initialize(p); err != nil {
				t.Fatal(err)
			}
			b, _ := os.ReadFile(s.ProjectPath())
			var raw map[string]any
			_ = json.Unmarshal(b, &raw)
			raw["schema_version"] = version
			b, _ = json.Marshal(raw)
			_ = os.WriteFile(s.ProjectPath(), b, 0644)
			before, _ := os.ReadFile(s.ProjectPath())
			plan, err := s.PlanMigration()
			if err != nil || plan.Supported || !plan.Required {
				t.Fatalf("%+v %v", plan, err)
			}
			if _, err = s.Migrate(); err == nil {
				t.Fatal("unsupported migration accepted")
			}
			after, _ := os.ReadFile(s.ProjectPath())
			if string(before) != string(after) {
				t.Fatal("unsupported migration changed projection")
			}
		})
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
	for i := 0; i < 4; i++ {
		if _, err := s.Mutate(nil, domain.Actor{Kind: "test"}, "TestChanged", p.ProjectID, map[string]int{"iteration": i}, func(p *domain.Project) error {
			p.Name += "x"
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	before, _ := os.Stat(s.EventsPath())
	if _, err := s.Compact(); err != nil {
		t.Fatal(err)
	}
	after, _ := os.Stat(s.EventsPath())
	if after.Size() >= before.Size() {
		t.Fatalf("active journal did not shrink: before=%d after=%d", before.Size(), after.Size())
	}
	events, err := s.VerifyJournal()
	if err != nil || len(events) != 1 || events[0].EventType != "ProjectionCheckpoint" || events[0].Sequence != 6 {
		t.Fatalf("unexpected compacted journal: %#v %v", events, err)
	}
	archives, err := s.VerifyArchivedSegments()
	if err != nil || len(archives) != 1 || archives[0].EventCount != 5 || archives[0].LastSequence != 5 {
		t.Fatalf("unexpected archive: %#v %v", archives, err)
	}
	replayed, err := s.Replay()
	if err != nil || replayed.Name != "testxxxx" {
		t.Fatalf("semantic replay changed: %#v %v", replayed, err)
	}
	if _, err := s.Mutate(nil, domain.Actor{Kind: "test"}, "AfterCompact", p.ProjectID, nil, func(p *domain.Project) error { p.Name = "after"; return nil }); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Compact(); err != nil {
		t.Fatal(err)
	}
	archives, err = s.VerifyArchivedSegments()
	if err != nil || len(archives) != 2 {
		t.Fatalf("repeated compaction did not retain both segments: %#v %v", archives, err)
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

func TestArchivedSegmentTamperFailsClosed(t *testing.T) {
	s := New(t.TempDir())
	p := domain.NewProject("project_test", "test", "standard", time.Now())
	if err := s.Initialize(p); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.CreateCheckpoint(); err != nil {
		t.Fatal(err)
	}
	archives, err := s.VerifyArchivedSegments()
	if err != nil || len(archives) != 1 {
		t.Fatalf("archive missing: %#v %v", archives, err)
	}
	path := filepath.Join(s.Root, filepath.FromSlash(archives[0].Path))
	if err := os.WriteFile(path, []byte("tampered\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := s.VerifyArchivedSegments(); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("tampered archive accepted: %v", err)
	}
	activeBefore, _ := os.ReadFile(s.EventsPath())
	if _, _, err := s.CreateCheckpoint(); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("compaction continued over corrupt history: %v", err)
	}
	activeAfter, _ := os.ReadFile(s.EventsPath())
	if string(activeBefore) != string(activeAfter) {
		t.Fatal("failed compaction changed active journal")
	}
}

func TestCompactWriteFailurePreservesActiveJournal(t *testing.T) {
	s := New(t.TempDir())
	p := domain.NewProject("project_test", "test", "standard", time.Now())
	if err := s.Initialize(p); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(s.EventsPath())
	if err := os.WriteFile(s.EventArchiveDir(), []byte("not-a-directory"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.CreateCheckpoint(); err == nil {
		t.Fatal("archive write failure was accepted")
	}
	after, _ := os.ReadFile(s.EventsPath())
	if string(before) != string(after) {
		t.Fatal("failed archive write changed active journal")
	}
}

func TestCompactedJournalRequiresValidArchiveAnchor(t *testing.T) {
	s := New(t.TempDir())
	p := domain.NewProject("project_test", "test", "standard", time.Now())
	if err := s.Initialize(p); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.CreateCheckpoint(); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(s.EventsPath())
	var event domain.Event
	if err := json.Unmarshal(bytes.TrimSpace(b), &event); err != nil {
		t.Fatal(err)
	}
	event.Payload = map[string]any{"projection": p}
	event.Hash = eventHash(event)
	line, _ := json.Marshal(event)
	if err := os.WriteFile(s.EventsPath(), append(line, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := s.VerifyJournal(); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("unanchored compacted journal accepted: %v", err)
	}
}
