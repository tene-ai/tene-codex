// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/tene-ai/tene-codex/internal/domain"
	"github.com/tene-ai/tene-codex/internal/state"
	"github.com/tene-ai/tene-codex/internal/workflow"
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

func TestEvidenceRejectsCredentialAndCanaryPatternsBeforeMutation(t *testing.T) {
	patterns := []string{"TENE_TEST_CANARY_0123456789", "ghp_12345678901234567890", "token=plaintext-value"}
	for _, value := range patterns {
		t.Run(value[:min(8, len(value))], func(t *testing.T) {
			root := t.TempDir()
			if code, _ := execute(t, root, "init", "--name", "security"); code != 0 {
				t.Fatal(code)
			}
			if code, _ := execute(t, root, "sprint", "create", "--title", "Security"); code != 0 {
				t.Fatal(code)
			}
			path := filepath.Join(root, "poisoned.txt")
			if err := os.WriteFile(path, []byte(value), 0o600); err != nil {
				t.Fatal(err)
			}
			code, got := execute(t, root, "evidence", "register", "--path", path)
			if code != 6 || len(got.Errors) == 0 || got.Errors[0].Code != "SEC_EVIDENCE_LEAK" {
				t.Fatalf("code=%d result=%+v", code, got)
			}
			p, err := state.New(root).Load()
			if err != nil {
				t.Fatal(err)
			}
			if len(p.Evidence) != 0 {
				t.Fatal("poisoned evidence mutated state")
			}
		})
	}
}

func TestEvidenceRegistrationDeduplicatesContentContract(t *testing.T) {
	root := t.TempDir()
	execute(t, root, "init", "--name", "evidence")
	execute(t, root, "sprint", "create", "--title", "Evidence")
	path := filepath.Join(root, "result.txt")
	if err := os.WriteFile(path, []byte("deterministic result"), 0644); err != nil {
		t.Fatal(err)
	}
	code, first := execute(t, root, "evidence", "register", "--path", path, "--kind", "test-output")
	if code != 0 {
		t.Fatal(first)
	}
	code, second := execute(t, root, "evidence", "register", "--path", path, "--kind", "test-output")
	if code != 0 {
		t.Fatal(second)
	}
	a := first.Result.(map[string]any)["evidence_id"]
	b := second.Result.(map[string]any)["evidence_id"]
	if a != b || first.Revision != second.Revision {
		t.Fatalf("not deduplicated: %#v %#v", first, second)
	}
}

func TestMasterPlanStatusAndDependencyValidation(t *testing.T) {
	root := t.TempDir()
	if c, _ := execute(t, root, "init", "--name", "master"); c != 0 {
		t.Fatal(c)
	}
	if c, _ := execute(t, root, "sprint", "create", "--title", "One"); c != 0 {
		t.Fatal(c)
	}
	if c, e := execute(t, root, "master", "validate"); c != 0 {
		t.Fatalf("%d %+v", c, e)
	}
	if c, e := execute(t, root, "master", "status"); c != 0 || e.Result.(map[string]any)["project_id"] == nil {
		t.Fatalf("%d %+v", c, e)
	}
}

func TestMasterMetadataAndStructuredIntentLifecycle(t *testing.T) {
	root := t.TempDir()
	execute(t, root, "init", "--name", "semantic")
	if code, env := execute(t, root, "master", "create", "--objective", "ship safely", "--milestones", "mvp,ga", "--releases", "v1", "--risks", "false-pass", "--invariants", "evidence required,secret safe"); code != 0 || env.Result.(map[string]any)["plan"].(map[string]any)["objective"] != "ship safely" {
		t.Fatalf("code=%d env=%#v", code, env)
	}
	_, created := execute(t, root, "sprint", "create", "--title", "Rich intent", "--milestone", "mvp", "--release", "v1")
	sprint := created.Result.(map[string]any)["sprint"].(map[string]any)
	if sprint["milestone"] != "mvp" || sprint["release"] != "v1" {
		t.Fatalf("%#v", sprint)
	}
	execute(t, root, "phase", "transition", "prd")
	_, captured := execute(t, root, "intent", "capture", "--statement", "remember semantics", "--actors", "buyer,operator", "--outcomes", "visible success", "--policies", "tenant isolation", "--business-rules", "one active order", "--ux-states", "empty,confirmed", "--data-invariants", "one write", "--constraints", "offline", "--assumptions", "local", "--open-questions", "retention", "--source-locator", "conversation:42", "--ac", "works", "--observable", "confirmed")
	intent := captured.Result.(map[string]any)["intent"].(map[string]any)
	if len(intent["actors"].([]any)) != 2 || intent["source_locator"] != "conversation:42" {
		t.Fatalf("%#v", intent)
	}
	id := intent["intent_id"].(string)
	if code, _ := execute(t, root, "intent", "confirm", id); code != 0 {
		t.Fatal(code)
	}
	code, superseded := execute(t, root, "intent", "supersede", id, "--statement", "remember revised semantics")
	if code != 0 || superseded.Result.(map[string]any)["superseded"].(map[string]any)["status"] != "superseded" || superseded.Result.(map[string]any)["replacement"].(map[string]any)["supersedes"] != id {
		t.Fatalf("%#v", superseded)
	}
	if code, _ := execute(t, root, "graph", "build"); code != 0 {
		t.Fatal(code)
	}
	p, err := state.New(root).Load()
	if err != nil {
		t.Fatal(err)
	}
	kinds := map[string]bool{}
	for _, node := range p.Graph.Nodes {
		kinds[node.Kind] = true
	}
	if !kinds["Policy"] || !kinds["DocumentSection"] || !kinds["Sprint"] || !kinds["Intent"] {
		t.Fatalf("graph kinds=%#v", kinds)
	}
}

func TestGraphBuildKeepsJournalProjectionEquivalent(t *testing.T) {
	root := t.TempDir()
	execute(t, root, "init", "--name", "graph replay")
	execute(t, root, "sprint", "create", "--title", "Graph replay")
	execute(t, root, "graph", "build")
	store := state.New(root)
	live, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	replayed, err := store.Replay()
	if err != nil {
		t.Fatal(err)
	}
	a, _ := json.MarshalIndent(live, "", "  ")
	b, _ := json.MarshalIndent(replayed, "", "  ")
	if !bytes.Equal(a, b) {
		t.Fatalf("live/replay mismatch\nlive=%s\nreplayed=%s", a, b)
	}
	projectBytes, _ := os.ReadFile(filepath.Join(root, ".tene-workflow", "project.json"))
	activeBytes, _ := os.ReadFile(filepath.Join(root, ".tene-workflow", "active.json"))
	if len(activeBytes) >= len(projectBytes) || bytes.Contains(activeBytes, []byte(`"graph"`)) || !bytes.Contains(activeBytes, []byte(`"active_sprint"`)) {
		t.Fatalf("active projection is not bounded or resumable: active=%d project=%d", len(activeBytes), len(projectBytes))
	}
}

func TestQAPlanVariantContract(t *testing.T) {
	want := []string{"happy", "alternate", "empty", "validation", "permission", "failure", "recovery"}
	if got := qaVariants(); !slices.Equal(got, want) {
		t.Fatalf("%v", got)
	}
	for _, v := range want {
		if qaVariantAction(v) == "" {
			t.Fatalf("missing action for %s", v)
		}
	}
}

func TestRequestIDDeduplicatesAndRejectsReuse(t *testing.T) {
	root := t.TempDir()
	execute(t, root, "init", "--name", "dedup")
	code, first := execute(t, root, "sprint", "create", "--title", "Once", "--request-id", "req-once")
	if code != 0 {
		t.Fatal(first)
	}
	code, second := execute(t, root, "sprint", "create", "--title", "Once", "--request-id", "req-once")
	if code != 0 || first.Revision != second.Revision {
		t.Fatalf("not deduplicated: %#v %#v", first, second)
	}
	p, err := state.New(root).Load()
	if err != nil || len(p.Sprints) != 1 {
		t.Fatalf("duplicate mutation: %v %#v", err, p)
	}
	code, conflict := execute(t, root, "sprint", "create", "--title", "Different", "--request-id", "req-once")
	if code != 4 || len(conflict.Errors) == 0 || conflict.Errors[0].Code != "REQUEST_ID_CONFLICT" {
		t.Fatalf("reuse accepted: %d %#v", code, conflict)
	}
}

func TestRequestCrashWindowFailsClosedWithoutReexecution(t *testing.T) {
	root := t.TempDir()
	execute(t, root, "init", "--name", "crash")
	args := []string{"sprint", "create", "--title", "Crash"}
	hash := hashStrings(args)
	s := state.New(root)
	if _, err := s.Mutate(nil, domain.Actor{Kind: "codex", ID: hash, SessionID: "req-crash"}, "SyntheticCommittedMutation", "project", nil, func(*domain.Project) error { return nil }); err != nil {
		t.Fatal(err)
	}
	code, env := execute(t, root, "sprint", "create", "--title", "Crash", "--request-id", "req-crash")
	if code != 4 || len(env.Errors) == 0 || env.Errors[0].Code != "REQUEST_RECOVERY_PENDING" {
		t.Fatalf("crash retry did not fail closed: %d %#v", code, env)
	}
	p, _ := s.Load()
	if len(p.Sprints) != 0 {
		t.Fatal("handler re-executed after committed placeholder")
	}
}

func TestDocumentSyncAliasesFlagsAndPreservesAuthoredText(t *testing.T) {
	root := t.TempDir()
	execute(t, root, "init", "--name", "docs")
	_, created := execute(t, root, "sprint", "create", "--title", "Sync")
	sp := created.Result.(map[string]any)["sprint"].(map[string]any)
	path := filepath.Join(root, filepath.FromSlash(sp["document_root"].(string)), "00-prd", "00-prd.md")
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	_, _ = f.WriteString("\nAUTHORED-SENTINEL\n")
	_ = f.Close()
	code, preview := execute(t, root, "docs", "sync", "--phase", "prd", "--no-color", "--verbose")
	if code != 0 || preview.Result.(map[string]any)["applied"].(bool) {
		t.Fatalf("preview failed %#v", preview)
	}
	before, _ := os.ReadFile(path)
	if !strings.Contains(string(before), "AUTHORED-SENTINEL") || strings.Contains(string(before), "tene:generated:traceability:start") {
		t.Fatal("preview mutated")
	}
	if code, _ := execute(t, root, "document", "sync", "--phase", "prd", "--apply"); code != 0 {
		t.Fatal(code)
	}
	after, _ := os.ReadFile(path)
	if !strings.Contains(string(after), "AUTHORED-SENTINEL") || !strings.Contains(string(after), "tene:generated:traceability:start") {
		t.Fatal("sync contract failed")
	}
	var out, errOut bytes.Buffer
	if code := Run([]string{"--root", root, "--quiet", "status"}, &out, &errOut, "test"); code != 0 || out.Len() != 0 {
		t.Fatalf("quiet output %d %q", code, out.String())
	}
}

func TestWaiverRequestRequiresApproval(t *testing.T) {
	root := t.TempDir()
	execute(t, root, "init", "--name", "waiver")
	execute(t, root, "sprint", "create", "--title", "Waiver")
	_, gapEnv := execute(t, root, "loop", "record-gap", "--description", "accepted risk", "--category", "debt", "--severity", "blocker")
	gapID := gapEnv.Result.(map[string]any)["gap_id"].(string)
	expiry := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	code, requested := execute(t, root, "waiver", "request", "--gap", gapID, "--reason", "bounded", "--requester", "owner", "--expires", expiry)
	if code != 0 {
		t.Fatal(requested)
	}
	id := requested.Result.(map[string]any)["waiver_id"].(string)
	p, _ := state.New(root).Load()
	if p.Waivers[id].Status != "requested" || workflow.ActiveWaiver(p, p.Gaps[gapID], time.Now()) {
		t.Fatal("request became active")
	}
	if code, _ := execute(t, root, "waiver", "approve", id, "--approver", "reviewer"); code != 0 {
		t.Fatal(code)
	}
	p, _ = state.New(root).Load()
	if !workflow.ActiveWaiver(p, p.Gaps[gapID], time.Now()) {
		t.Fatal("approved waiver inactive")
	}
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
	plannedRun := planned.Result.(map[string]any)
	cases := plannedRun["cases"].([]any)
	firstCaseID := cases[0].(map[string]any)["case_id"].(string)
	if code, env := execute(t, root, "qa", "case", firstCaseID, "passed", "--evidence", "anything"); code != 3 || len(env.Errors) == 0 || env.Errors[0].Code != "QA_MANUAL_PASS_FORBIDDEN" {
		t.Fatalf("manual pass was not rejected: code=%d env=%#v", code, env)
	}
	evidenceDir := filepath.Join(originalRoot, "04-qa", "evidence")
	if err := os.MkdirAll(evidenceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, raw := range cases {
		qaCase := raw.(map[string]any)
		caseID, variant := qaCase["case_id"].(string), qaCase["variant"].(string)
		evidencePath := filepath.Join(evidenceDir, caseID+".json")
		now := time.Now().UTC()
		assertions := []map[string]any{}
		for _, layer := range []string{"L1", "L2", "L3", "L4", "L5", "L6", "L7"} {
			refs := []string{"layer:" + layer}
			if layer == "L5" {
				refs = append(refs, "observable", "variant:"+variant)
			}
			assertions = append(assertions, map[string]any{"statement": "verified " + layer, "passed": true, "layer": layer, "requirement_refs": refs, "actual": "verified", "expected": "verified"})
		}
		observation := map[string]any{"schema_version": "1.0.0", "adapter": "test-observer", "run_id": plannedRun["run_id"], "case_id": caseID, "environment": plannedRun["environment"], "started_at": now, "finished_at": now.Add(time.Second), "redaction_status": "passed", "spec_hash": plannedRun["spec_hash"], "state_revision": plannedRun["state_revision"], "layers": []string{"L1", "L2", "L3", "L4", "L5", "L6", "L7"}, "tool_version": "test-1", "checkpoints": []map[string]any{{"name": "journey", "kind": "ui-api-data", "before": map[string]any{"state": "start"}, "after": map[string]any{"state": "done"}}}, "assertions": assertions}
		b, _ := json.Marshal(observation)
		if err := os.WriteFile(evidencePath, b, 0o644); err != nil {
			t.Fatal(err)
		}
		if code, env := execute(t, root, "qa", "observe", caseID, "--input", evidencePath); code != 0 {
			t.Fatalf("case %s failed: %#v", caseID, env)
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
	expiry := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	code, requested := execute(t, root, "approval", "request", "--from", "report", "--to", "archived", "--reason", "final report reviewed", "--requester", "test", "--expires", expiry)
	if code != 0 {
		t.Fatalf("approval request failed: %#v", requested)
	}
	approvalID := requested.Result.(map[string]any)["approval_id"].(string)
	if code, _ := execute(t, root, "approval", "approve", approvalID, "--approver", "reviewer"); code != 0 {
		t.Fatal("approval failed")
	}
	if code, _ := execute(t, root, "sprint", "archive", "--approval", approvalID); code != 0 {
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

func TestCLICompactBoundsJournalAndDoctorDetectsArchiveTamper(t *testing.T) {
	root := t.TempDir()
	if code, _ := execute(t, root, "init", "--name", "compact-fixture"); code != 0 {
		t.Fatal("init failed")
	}
	if code, _ := execute(t, root, "sprint", "create", "--title", "Compaction"); code != 0 {
		t.Fatal("mutation fixture failed")
	}
	before, err := os.Stat(filepath.Join(root, state.DirName, "events.ndjson"))
	if err != nil {
		t.Fatal(err)
	}
	code, compacted := execute(t, root, "compact")
	if code != 0 {
		t.Fatalf("compact failed: %#v", compacted)
	}
	result := compacted.Result.(map[string]any)
	if result["active_events"].(float64) != 1 {
		t.Fatalf("unexpected compact result: %#v", result)
	}
	after, _ := os.Stat(filepath.Join(root, state.DirName, "events.ndjson"))
	if after.Size() >= before.Size() {
		t.Fatalf("active journal did not shrink: %d >= %d", after.Size(), before.Size())
	}
	if code, checked := execute(t, root, "doctor"); code != 0 || !checked.Result.(map[string]any)["healthy"].(bool) || checked.Result.(map[string]any)["archived_event_segments"].(float64) != 1 {
		t.Fatalf("doctor did not verify archive: %#v", checked)
	}
	archive := result["archived_segment"].(map[string]any)
	segmentPath := filepath.Join(root, filepath.FromSlash(archive["path"].(string)))
	if err := os.WriteFile(segmentPath, []byte("tampered\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code, checked := execute(t, root, "doctor"); code != 0 || checked.Result.(map[string]any)["healthy"].(bool) {
		t.Fatalf("doctor accepted tampered archive: code=%d result=%#v", code, checked)
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

func TestCLIGraphImpactAndContextFreshness(t *testing.T) {
	root := t.TempDir()
	execute(t, root, "init", "--name", "context fixture")
	_, created := execute(t, root, "sprint", "create", "--title", "Fresh context")
	sprint := created.Result.(map[string]any)["sprint"].(map[string]any)
	execute(t, root, "phase", "transition", "prd")
	completeDocument(t, root, filepath.ToSlash(filepath.Join(sprint["document_root"].(string), "00-prd", "00-prd.md")))
	_, captured := execute(t, root, "intent", "capture", "--statement", "fresh context", "--ac", "impact is traced", "--observable", "an AC is returned")
	intentID := captured.Result.(map[string]any)["intent"].(map[string]any)["intent_id"].(string)
	acID := captured.Result.(map[string]any)["criterion"].(map[string]any)["ac_id"].(string)
	execute(t, root, "intent", "confirm", intentID)
	if code, _ := execute(t, root, "graph", "build"); code != 0 {
		t.Fatal(code)
	}
	code, impacted := execute(t, root, "graph", "impact", intentID, "--depth", "4")
	if code != 0 {
		t.Fatalf("impact failed: %#v", impacted)
	}
	ids := impacted.Result.(map[string]any)["impacted_ac_ids"].([]any)
	if len(ids) != 1 || ids[0] != acID {
		t.Fatalf("unexpected impact: %#v", impacted.Result)
	}
	packPath := filepath.Join(".tene-workflow", "cache", "prd-context.json")
	code, built := execute(t, root, "context", "build", "--phase", "prd", "--budget", "8192", "--output", packPath)
	if code != 0 || built.Result.(map[string]any)["budget_unit"] != "utf8-bytes" {
		t.Fatalf("build failed: %#v", built)
	}
	if code, valid := execute(t, root, "context", "validate", "--input", packPath); code != 0 || !valid.Result.(map[string]any)["fresh"].(bool) {
		t.Fatalf("validate failed: %#v", valid)
	}
	policy := filepath.Join(root, ".tene-workflow", "policies.yaml")
	if err := os.WriteFile(policy, []byte("changed: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code, stale := execute(t, root, "context", "validate", "--input", packPath); code != 3 || stale.OK {
		t.Fatalf("expected stale context: code=%d %#v", code, stale)
	}
}

func TestCLIApprovalLoopAndGapLifecycle(t *testing.T) {
	root := t.TempDir()
	execute(t, root, "init", "--profile", "strict")
	_, created := execute(t, root, "sprint", "create", "--title", "Guarded", "--max-iterations", "2")
	sprint := created.Result.(map[string]any)["sprint"].(map[string]any)
	execute(t, root, "phase", "transition", "prd")
	completeDocument(t, root, filepath.ToSlash(filepath.Join(sprint["document_root"].(string), "00-prd", "00-prd.md")))
	_, captured := execute(t, root, "intent", "capture", "--statement", "guarded flow", "--ac", "flow is controlled", "--observable", "state is visible")
	intentID := captured.Result.(map[string]any)["intent"].(map[string]any)["intent_id"].(string)
	acID := captured.Result.(map[string]any)["criterion"].(map[string]any)["ac_id"].(string)
	execute(t, root, "intent", "confirm", intentID)
	execute(t, root, "phase", "transition", "plan")
	completeDocument(t, root, filepath.ToSlash(filepath.Join(sprint["document_root"].(string), "01-plan", "00-plan.md")))
	_, taskEnv := execute(t, root, "task", "add", "--title", "work", "--ac", acID)
	taskID := taskEnv.Result.(map[string]any)["task_id"].(string)
	execute(t, root, "phase", "transition", "design")
	completeDocument(t, root, filepath.ToSlash(filepath.Join(sprint["document_root"].(string), "02-design", "00-design.md")))
	if code, env := execute(t, root, "phase", "transition", "do", "--dry-run"); code != 3 || env.OK || env.Errors[0].Details == nil || env.Revision == 0 || env.Errors[0].Retryable {
		t.Fatalf("approval boundary passed: %#v", env)
	}
	expires := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	_, request := execute(t, root, "approval", "request", "--to", "do", "--reason", "design reviewed", "--requester", "agent", "--expires", expires)
	approvalID := request.Result.(map[string]any)["approval_id"].(string)
	execute(t, root, "approval", "approve", approvalID, "--approver", "human")
	if code, _ := execute(t, root, "phase", "transition", "do", "--approval", approvalID, "--dry-run"); code != 0 {
		t.Fatal("approved dry-run failed")
	}
	if code, _ := execute(t, root, "phase", "transition", "do", "--approval", approvalID); code != 0 {
		t.Fatal("approved transition failed")
	}
	if code, _ := execute(t, root, "approval", "approve", approvalID, "--approver", "human"); code == 0 {
		t.Fatal("consumed approval reused")
	}
	execute(t, root, "task", "complete", taskID)
	execute(t, root, "phase", "transition", "loop-check")
	if code, _ := execute(t, root, "loop", "iterate", "--outcome", "repair", "--summary", "first repair"); code != 0 {
		t.Fatal(code)
	}
	if code, _ := execute(t, root, "loop", "iterate", "--outcome", "blocked", "--summary", "policy needed"); code != 0 {
		t.Fatal(code)
	}
	if code, _ := execute(t, root, "loop", "iterate", "--outcome", "repair", "--summary", "too many"); code != 3 {
		t.Fatalf("expected exhausted, got %d", code)
	}
	_, gapEnv := execute(t, root, "loop", "record-gap", "--description", "carry debt", "--category", "debt", "--severity", "warning")
	gapID := gapEnv.Result.(map[string]any)["gap_id"].(string)
	if code, _ := execute(t, root, "loop", "defer-gap", gapID, "--reason", "bounded scope", "--owner", "team", "--target-sprint", "next"); code != 0 {
		t.Fatal(code)
	}
	_, securityEnv := execute(t, root, "loop", "record-gap", "--description", "security defect", "--category", "security")
	securityID := securityEnv.Result.(map[string]any)["gap_id"].(string)
	if code, _ := execute(t, root, "loop", "defer-gap", securityID, "--reason", "no", "--owner", "team", "--target-sprint", "next"); code != 3 {
		t.Fatalf("security defer code=%d", code)
	}
}

func TestGraphBuildUsesDeterministicEdgeIDs(t *testing.T) {
	root := t.TempDir()
	execute(t, root, "init")
	execute(t, root, "sprint", "create", "--title", "Stable graph")
	_, captured := execute(t, root, "intent", "capture", "--statement", "stable", "--ac", "stable edge", "--observable", "same id")
	intentID := captured.Result.(map[string]any)["intent"].(map[string]any)["intent_id"].(string)
	execute(t, root, "intent", "confirm", intentID)
	execute(t, root, "graph", "build")
	first, err := state.New(root).Load()
	if err != nil {
		t.Fatal(err)
	}
	execute(t, root, "graph", "build")
	second, err := state.New(root).Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Graph.Edges) != len(second.Graph.Edges) {
		t.Fatal("edge count drift")
	}
	for id := range first.Graph.Edges {
		if _, ok := second.Graph.Edges[id]; !ok {
			t.Fatalf("edge id drift: %s", id)
		}
	}
}
