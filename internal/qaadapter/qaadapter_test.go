package qaadapter

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDiscoverGoAndPlaywright(t *testing.T) {
	r := t.TempDir()
	os.WriteFile(filepath.Join(r, "go.mod"), []byte("module x\n"), 0644)
	os.WriteFile(filepath.Join(r, "package.json"), []byte(`{"scripts":{"test":"vitest"},"devDependencies":{"@playwright/test":"1"}}`), 0644)
	c := Discover(r)
	if !c[0].Available || !c[1].Available || !c[2].Available || !c[3].Available {
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
