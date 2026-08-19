// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package document

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tene-ai/tene-codex/internal/domain"
)

func TestScaffoldAndValidateAllDocuments(t *testing.T) {
	root := t.TempDir()
	p := domain.NewProject("project_test", "test", "standard", time.Now())
	sp := &domain.Sprint{SprintID: "sprint_test", Title: "Feature", DocumentRoot: "docs/sprints/sprint_test-feature"}
	paths, err := ScaffoldAll(root, p, sp)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 6 {
		t.Fatalf("expected 6 documents, got %d", len(paths))
	}
	for _, path := range paths {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		text := strings.ReplaceAll(string(b), "\n\n", "\n\n")
		var lines []string
		for _, line := range strings.Split(text, "\n") {
			lines = append(lines, line)
			if strings.HasPrefix(line, "## ") {
				lines = append(lines, "", "Completed with evidence.")
			}
		}
		if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, ph := range []domain.Phase{domain.PhasePRD, domain.PhasePlan, domain.PhaseDesign, domain.PhaseLoopCheck, domain.PhaseQA, domain.PhaseReport} {
		if got := Validate(Path(root, sp, ph), ph); len(got) != 0 {
			t.Fatalf("%s invalid: %#v", ph, got)
		}
	}
}

func TestValidationFindsEmptySections(t *testing.T) {
	root := t.TempDir()
	p := domain.NewProject("project_test", "test", "standard", time.Now())
	sp := &domain.Sprint{SprintID: "sprint_test", Title: "Feature", DocumentRoot: "docs/sprints/sprint_test-feature"}
	path, _, err := Scaffold(root, p, sp, domain.PhasePRD)
	if err != nil {
		t.Fatal(err)
	}
	if got := Validate(path, domain.PhasePRD); len(got) == 0 {
		t.Fatal("expected empty-section findings")
	}
}
