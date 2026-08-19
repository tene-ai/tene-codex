// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package projectconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScaffoldAgentsPreservesExistingFile(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".codex", "agents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "tene-builder.toml")
	if err := os.WriteFile(path, []byte("custom=true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	created, err := ScaffoldAgents(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(created) != 4 {
		t.Fatalf("expected 4 new agents, got %d", len(created))
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "custom=true\n" {
		t.Fatalf("existing agent overwritten: %q", b)
	}
}
