// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"errors"
	"testing"
	"time"

	"github.com/tene-ai/tene-codex/internal/domain"
)

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
