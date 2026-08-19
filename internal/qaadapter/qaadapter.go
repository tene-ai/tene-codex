// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package qaadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Capability struct {
	Adapter   string   `json:"adapter"`
	Available bool     `json:"available"`
	Kind      string   `json:"kind"`
	Layers    []string `json:"layers"`
	Command   []string `json:"command,omitempty"`
	Reason    string   `json:"reason,omitempty"`
}
type Checkpoint struct {
	Name   string `json:"name"`
	Kind   string `json:"kind"`
	Before any    `json:"before,omitempty"`
	After  any    `json:"after,omitempty"`
}
type Assertion struct {
	Statement string `json:"statement"`
	Passed    bool   `json:"passed"`
	Actual    string `json:"actual,omitempty"`
	Expected  string `json:"expected,omitempty"`
}
type Observation struct {
	SchemaVersion   string       `json:"schema_version"`
	Adapter         string       `json:"adapter"`
	RunID           string       `json:"run_id"`
	CaseID          string       `json:"case_id"`
	Environment     string       `json:"environment"`
	StartedAt       time.Time    `json:"started_at"`
	FinishedAt      time.Time    `json:"finished_at"`
	Checkpoints     []Checkpoint `json:"checkpoints"`
	Assertions      []Assertion  `json:"assertions"`
	RedactionStatus string       `json:"redaction_status"`
}
type ExecutionResult struct {
	Adapter    string    `json:"adapter"`
	Command    []string  `json:"command"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	Passed     bool      `json:"passed"`
	Output     string    `json:"output"`
}

func Execute(ctx context.Context, root, adapter string) (ExecutionResult, error) {
	var selected *Capability
	for _, c := range Discover(root) {
		if c.Adapter == adapter {
			x := c
			selected = &x
			break
		}
	}
	if selected == nil {
		return ExecutionResult{}, fmt.Errorf("QA_ADAPTER_UNKNOWN: %s", adapter)
	}
	if !selected.Available || len(selected.Command) == 0 {
		return ExecutionResult{}, fmt.Errorf("QA_ADAPTER_UNAVAILABLE: %s", selected.Reason)
	}
	started := time.Now().UTC()
	cmd := exec.CommandContext(ctx, selected.Command[0], selected.Command[1:]...)
	cmd.Dir = root
	b, err := cmd.CombinedOutput()
	if len(b) > 1<<20 {
		b = b[len(b)-(1<<20):]
	}
	result := ExecutionResult{Adapter: adapter, Command: selected.Command, StartedAt: started, FinishedAt: time.Now().UTC(), Passed: err == nil, Output: string(b)}
	if err != nil {
		return result, fmt.Errorf("QA_ADAPTER_FAILED: %s", adapter)
	}
	return result, nil
}

func Discover(root string) []Capability {
	_, goBin := exec.LookPath("go")
	_, goMod := os.Stat(filepath.Join(root, "go.mod"))
	goCap := Capability{Adapter: "go-test", Available: goBin == nil && goMod == nil, Kind: "native-test", Layers: []string{"L1", "L2", "L3"}, Command: []string{"go", "test", "./..."}}
	if goCap.Available == false {
		goCap.Command = nil
		goCap.Reason = "go executable or go.mod is unavailable"
	}
	npm := Capability{Adapter: "npm-test", Kind: "native-test", Layers: []string{"L1", "L2", "L3"}}
	if _, err := exec.LookPath("npm"); err == nil {
		if hasScript(filepath.Join(root, "package.json"), "test") {
			npm.Available = true
			npm.Command = []string{"npm", "test"}
		}
	}
	if !npm.Available {
		npm.Reason = "npm or package.json test script is unavailable"
	}
	pw := Capability{Adapter: "playwright", Kind: "browser-e2e", Layers: []string{"L4", "L5", "L6", "L7"}}
	if _, err := exec.LookPath("npx"); err == nil && (exists(filepath.Join(root, "playwright.config.ts")) || exists(filepath.Join(root, "playwright.config.js")) || packageContains(filepath.Join(root, "package.json"), "@playwright/test")) {
		pw.Available = true
		pw.Command = []string{"npx", "playwright", "test"}
	} else {
		pw.Reason = "Playwright configuration/dependency is unavailable"
	}
	browser := Capability{Adapter: "external-browser-observation", Available: true, Kind: "observation-import", Layers: []string{"L3", "L4", "L5", "L6", "L7"}, Reason: "execution is provided by Codex browser/Chrome/Playwright capability; CLI validates and hashes the observation"}
	return []Capability{goCap, npm, pw, browser}
}

func ReadObservation(root, path string) (Observation, []byte, error) {
	abs := path
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(root, path)
	}
	cleanRoot, err := filepath.Abs(root)
	if err != nil {
		return Observation{}, nil, err
	}
	cleanPath, err := filepath.Abs(abs)
	if err != nil {
		return Observation{}, nil, err
	}
	rel, err := filepath.Rel(cleanRoot, cleanPath)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return Observation{}, nil, fmt.Errorf("QA_OBSERVATION_OUTSIDE_ROOT: input must be inside project")
	}
	b, err := os.ReadFile(cleanPath)
	if err != nil {
		return Observation{}, nil, err
	}
	var o Observation
	if err := json.Unmarshal(b, &o); err != nil {
		return o, nil, fmt.Errorf("QA_OBSERVATION_INVALID: %w", err)
	}
	if o.SchemaVersion != "1.0.0" || o.Adapter == "" || o.RunID == "" || o.CaseID == "" || o.StartedAt.IsZero() || o.FinishedAt.IsZero() || o.FinishedAt.Before(o.StartedAt) || len(o.Assertions) == 0 {
		return o, nil, fmt.Errorf("QA_OBSERVATION_INVALID: required fields or timestamps are invalid")
	}
	if o.RedactionStatus != "passed" {
		return o, nil, fmt.Errorf("QA_OBSERVATION_REDACTION_FAILED: artifact was not sanitized")
	}
	for _, a := range o.Assertions {
		if !a.Passed {
			return o, nil, fmt.Errorf("QA_OBSERVATION_ASSERTION_FAILED: %s", a.Statement)
		}
	}
	return o, b, nil
}
func exists(path string) bool { _, err := os.Stat(path); return err == nil }
func hasScript(path, name string) bool {
	b, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	var p struct {
		Scripts map[string]string `json:"scripts"`
	}
	return json.Unmarshal(b, &p) == nil && p.Scripts[name] != ""
}
func packageContains(path, needle string) bool {
	b, err := os.ReadFile(path)
	return err == nil && strings.Contains(string(b), needle)
}
