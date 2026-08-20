package qaadapter

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDiscoverGoAndPlaywright(t *testing.T) {
	r := t.TempDir()
	os.WriteFile(filepath.Join(r, "go.mod"), []byte("module x\n"), 0644)
	os.WriteFile(filepath.Join(r, "package.json"), []byte(`{"scripts":{"test":"vitest"},"devDependencies":{"@playwright/test":"1"}}`), 0644)
	os.Mkdir(filepath.Join(r, "tests"), 0755)
	c := Discover(r)
	if !c[0].Available || !c[1].Available || !c[2].Available || !c[3].Available || !c[4].Available {
		t.Fatalf("%#v", c)
	}
}

func TestExecuteRejectsUnknown(t *testing.T) {
	if _, err := Execute(context.Background(), t.TempDir(), "shell"); err == nil {
		t.Fatal("expected allowlist rejection")
	}
}
func TestReadObservation(t *testing.T) {
	r := t.TempDir()
	now := time.Now().UTC()
	o := Observation{SchemaVersion: "1.0.0", Adapter: "chrome", RunID: "run", CaseID: "case", Environment: "local", StartedAt: now, FinishedAt: now.Add(time.Second), SpecHash: "spec", StateRevision: 1, Layers: []string{"L5"}, ToolVersion: "chrome-test", Checkpoints: []Checkpoint{{Name: "journey", Kind: "ui-data", Before: map[string]any{"state": "start"}, After: map[string]any{"state": "done"}}}, Assertions: []Assertion{{Statement: "flow completes", Passed: true, Layer: "L5", RequirementRefs: []string{"observable"}, Actual: "done", Expected: "done"}}, RedactionStatus: "passed"}
	b, _ := json.Marshal(o)
	os.WriteFile(filepath.Join(r, "observation.json"), b, 0644)
	got, _, err := ReadObservation(r, "observation.json")
	if err != nil || got.CaseID != "case" {
		t.Fatalf("%#v %v", got, err)
	}
	o.Assertions[0].Passed = false
	b, _ = json.Marshal(o)
	os.WriteFile(filepath.Join(r, "bad.json"), b, 0644)
	if got, _, err := ReadObservation(r, "bad.json"); err != nil || got.Assertions[0].Passed {
		t.Fatal("failed assertions must be preserved as evidence")
	}
	o.Assertions[0].Layer = ""
	b, _ = json.Marshal(o)
	os.WriteFile(filepath.Join(r, "invalid.json"), b, 0644)
	if _, _, err := ReadObservation(r, "invalid.json"); err == nil {
		t.Fatal("expected structural rejection")
	}
}

func TestReadObservationRejectsAssertionCopiedAcrossLayers(t *testing.T) {
	now := time.Now().UTC()
	o := Observation{
		SchemaVersion: "1.0.0", Adapter: "codex-host", RunID: "run", CaseID: "case",
		Environment: "local", StartedAt: now, FinishedAt: now.Add(time.Second),
		SpecHash: "spec", StateRevision: 1, Layers: []string{"L1", "L2"},
		ToolVersion: "1", RedactionStatus: "passed",
		Checkpoints: []Checkpoint{{Name: "journey", Kind: "host", Before: map[string]any{"state": "before"}, After: map[string]any{"state": "after"}}},
		Assertions: []Assertion{
			{Statement: "L1 host contract", Passed: true, Layer: "L1", RequirementRefs: []string{"observable"}, Actual: "path passed", Expected: "path passed"},
			{Statement: "L2 host contract", Passed: true, Layer: "L2", RequirementRefs: []string{"observable"}, Actual: "path passed", Expected: "path passed"},
		},
	}
	b, err := json.Marshal(o)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "duplicated.json")
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := ReadObservation(dir, path); err == nil || !strings.Contains(err.Error(), "QA_OBSERVATION_DUPLICATE_LAYER_ASSERTION") {
		t.Fatalf("expected duplicated-layer assertion rejection, got %v", err)
	}
}
