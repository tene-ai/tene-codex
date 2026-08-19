// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tene-ai/tene-codex/internal/state"
)

func execute(t *testing.T, root string, args ...string) (int, envelope) {
	t.Helper()
	var out, err bytes.Buffer
	all := append([]string{"--root", root, "--json"}, args...)
	code := Run(all, &out, &err, "test")
	var env envelope
	if e := json.Unmarshal(out.Bytes(), &env); e != nil {
		t.Fatalf("decode %q stderr=%q: %v", out.String(), err.String(), e)
	}
	return code, env
}

func completeDocument(t *testing.T, root, relativePath string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relativePath))
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	for _, line := range strings.Split(string(b), "\n") {
		lines = append(lines, line)
		if strings.HasPrefix(line, "## ") {
			lines = append(lines, "", "Completed with traceable evidence.")
		}
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCLISprintHappyPathToDesign(t *testing.T) {
	root := t.TempDir()
	if code, _ := execute(t, root, "init", "--name", "fixture"); code != 0 {
		t.Fatal(code)
	}
	code, created := execute(t, root, "sprint", "create", "--title", "Login")
	if code != 0 {
		t.Fatal(code)
	}
	result := created.Result.(map[string]any)
	sprint := result["sprint"].(map[string]any)
	if sprint["phase"] != "draft" {
		t.Fatal(sprint)
	}
	if code, _ := execute(t, root, "phase", "transition", "prd"); code != 0 {
		t.Fatal(code)
	}
	completeDocument(t, root, filepath.ToSlash(filepath.Join(sprint["document_root"].(string), "00-prd", "00-prd.md")))
	code, captured := execute(t, root, "intent", "capture", "--statement", "Users can log in", "--ac", "Valid users reach the dashboard", "--observable", "dashboard visible")
	if code != 0 {
		t.Fatal(code)
	}
	intent := captured.Result.(map[string]any)["intent"].(map[string]any)
	id := intent["intent_id"].(string)
	acID := captured.Result.(map[string]any)["criterion"].(map[string]any)["ac_id"].(string)
	if code, _ := execute(t, root, "intent", "confirm", id); code != 0 {
		t.Fatal(code)
	}
	if code, _ := execute(t, root, "phase", "transition", "plan"); code != 0 {
		t.Fatal(code)
	}
	completeDocument(t, root, filepath.ToSlash(filepath.Join(sprint["document_root"].(string), "01-plan", "00-plan.md")))
	if code, _ := execute(t, root, "task", "add", "--title", "Implement login", "--layer", "business", "--ac", acID); code != 0 {
		t.Fatal(code)
	}
	if code, _ := execute(t, root, "phase", "transition", "design"); code != 0 {
		t.Fatal(code)
	}
}

func TestCLIGuardRejectsPRDWithoutIntent(t *testing.T) {
	root := t.TempDir()
	execute(t, root, "init")
	_, created := execute(t, root, "sprint", "create", "--title", "No intent")
	sprint := created.Result.(map[string]any)["sprint"].(map[string]any)
	execute(t, root, "phase", "transition", "prd")
	completeDocument(t, root, filepath.ToSlash(filepath.Join(sprint["document_root"].(string), "00-prd", "00-prd.md")))
	code, env := execute(t, root, "phase", "transition", "plan")
	if code != 3 || env.OK {
		t.Fatalf("expected guard failure: code=%d env=%#v", code, env)
	}
}

func TestCLICompleteSprintArchivesDocuments(t *testing.T) {
	root := t.TempDir()
	execute(t, root, "init", "--name", "fixture")
	_, created := execute(t, root, "sprint", "create", "--title", "Complete flow")
	sprint := created.Result.(map[string]any)["sprint"].(map[string]any)
	originalRoot := filepath.Join(root, filepath.FromSlash(sprint["document_root"].(string)))
	execute(t, root, "phase", "transition", "prd")
	completeDocument(t, root, filepath.ToSlash(filepath.Join(sprint["document_root"].(string), "00-prd", "00-prd.md")))
	_, captured := execute(t, root, "intent", "capture", "--statement", "User completes the flow", "--ac", "The flow completes", "--observable", "completion is visible")
	intentResult := captured.Result.(map[string]any)
	intentID := intentResult["intent"].(map[string]any)["intent_id"].(string)
	acID := intentResult["criterion"].(map[string]any)["ac_id"].(string)
	execute(t, root, "intent", "confirm", intentID)
	execute(t, root, "phase", "transition", "plan")
	completeDocument(t, root, filepath.ToSlash(filepath.Join(sprint["document_root"].(string), "01-plan", "00-plan.md")))
	execute(t, root, "task", "add", "--title", "Implement flow", "--layer", "business", "--ac", acID)
	execute(t, root, "phase", "transition", "design")
	completeDocument(t, root, filepath.ToSlash(filepath.Join(sprint["document_root"].(string), "02-design", "00-design.md")))
	execute(t, root, "phase", "transition", "do")
	execute(t, root, "phase", "transition", "loop-check")
	completeDocument(t, root, filepath.ToSlash(filepath.Join(sprint["document_root"].(string), "03-analysis", "00-loop-check.md")))
	execute(t, root, "phase", "transition", "qa")
	completeDocument(t, root, filepath.ToSlash(filepath.Join(sprint["document_root"].(string), "04-qa", "00-qa-plan.md")))
	_, planned := execute(t, root, "qa", "plan")
	cases := planned.Result.(map[string]any)["cases"].([]any)
	evidencePath := filepath.Join(originalRoot, "04-qa", "evidence", "qa-result.txt")
	if err := os.MkdirAll(filepath.Dir(evidencePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(evidencePath, []byte("all assertions passed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, registered := execute(t, root, "evidence", "register", "--path", evidencePath, "--ac", acID)
	evidenceID := registered.Result.(map[string]any)["evidence_id"].(string)
	for _, raw := range cases {
		caseID := raw.(map[string]any)["case_id"].(string)
		if code, _ := execute(t, root, "qa", "case", caseID, "passed", "--evidence", evidenceID); code != 0 {
			t.Fatalf("case %s failed", caseID)
		}
	}
	if code, _ := execute(t, root, "qa", "evaluate"); code != 0 {
		t.Fatal("qa evaluate failed")
	}
	execute(t, root, "phase", "transition", "report")
	completeDocument(t, root, filepath.ToSlash(filepath.Join(sprint["document_root"].(string), "05-report", "00-report.md")))
	execute(t, root, "report", "generate")
	if code, _ := execute(t, root, "report", "validate"); code != 0 {
		t.Fatal("report validation failed")
	}
	if code, _ := execute(t, root, "sprint", "archive"); code != 0 {
		t.Fatal("archive failed")
	}
	if _, err := os.Stat(originalRoot); !os.IsNotExist(err) {
		t.Fatalf("original document root still exists: %v", err)
	}
	matches, err := filepath.Glob(filepath.Join(root, "docs", "sprints", "_archive", "*", filepath.Base(originalRoot)))
	if err != nil || len(matches) != 1 {
		t.Fatalf("archive not found: %v %v", matches, err)
	}
	if _, err := os.Stat(filepath.Join(matches[0], "99-archive", "archive-manifest.json")); err != nil {
		t.Fatalf("archive manifest missing: %v", err)
	}
	p, err := state.New(root).Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, evidence := range p.Evidence {
		if !strings.HasPrefix(evidence.URI, filepath.ToSlash(filepath.Join("docs", "sprints", "_archive"))) {
			t.Fatalf("stale evidence URI: %s", evidence.URI)
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(evidence.URI))); err != nil {
			t.Fatalf("archived evidence missing: %v", err)
		}
	}
	if code, env := execute(t, root, "doctor"); code != 0 || !env.Result.(map[string]any)["healthy"].(bool) {
		t.Fatalf("post-archive doctor failed: %#v", env)
	}
}

func TestCLICodeIntelAndTaskReferenceValidation(t *testing.T) {
	root := t.TempDir()
	execute(t, root, "init")
	if err := os.WriteFile(filepath.Join(root, "sample.go"), []byte("package sample\nfunc Run(input string) string { return input }\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if code, env := execute(t, root, "graph", "understand", "--path", "sample.go"); code != 0 || len(env.Result.(map[string]any)["components"].([]any)) != 1 {
		t.Fatalf("code=%d env=%#v", code, env)
	}
	_, created := execute(t, root, "sprint", "create", "--title", "References")
	_ = created
	if code, _ := execute(t, root, "task", "add", "--title", "invalid", "--ac", "ac_missing"); code == 0 {
		t.Fatal("missing AC accepted")
	}
}
