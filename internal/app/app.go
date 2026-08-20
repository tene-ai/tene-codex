// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tene-ai/tene-codex/internal/codeintel"
	"github.com/tene-ai/tene-codex/internal/document"
	"github.com/tene-ai/tene-codex/internal/domain"
	"github.com/tene-ai/tene-codex/internal/projectconfig"
	"github.com/tene-ai/tene-codex/internal/qaadapter"
	"github.com/tene-ai/tene-codex/internal/router"
	"github.com/tene-ai/tene-codex/internal/secret"
	"github.com/tene-ai/tene-codex/internal/state"
	"github.com/tene-ai/tene-codex/internal/tracecontext"
	"github.com/tene-ai/tene-codex/internal/workflow"
)

type runtime struct {
	root      string
	json      bool
	expected  *uint64
	requestID string
	out, err  io.Writer
	version   string
}
type envelope struct {
	OK            bool       `json:"ok"`
	SchemaVersion string     `json:"schema_version"`
	RequestID     string     `json:"request_id,omitempty"`
	Revision      uint64     `json:"revision"`
	Result        any        `json:"result,omitempty"`
	Warnings      []string   `json:"warnings"`
	Errors        []apiError `json:"errors"`
}
type apiError struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Remediation string `json:"remediation,omitempty"`
}
type commandError struct {
	code, message, remediation string
	exit                       int
}

func (e *commandError) Error() string { return e.message }

func Run(args []string, stdout, stderr io.Writer, version string) int {
	rt := &runtime{out: stdout, err: stderr, version: version}
	args = rt.global(args)
	if rt.root == "" {
		rt.root = discoverRoot()
	}
	if len(args) == 0 {
		return rt.fail(&commandError{"CLI_USAGE", usage(), "Run tene-workflow help.", 2})
	}
	var result any
	var revision uint64
	var err error
	switch args[0] {
	case "help", "--help", "-h":
		fmt.Fprintln(stdout, usage())
		return 0
	case "version", "--version":
		result = map[string]string{"version": version}
	case "init":
		result, revision, err = rt.init(args[1:])
	case "status":
		result, revision, err = rt.status()
	case "route":
		result, revision, err = rt.route(args[1:])
	case "sprint":
		result, revision, err = rt.sprint(args[1:])
	case "master":
		result, revision, err = rt.master(args[1:])
	case "phase":
		result, revision, err = rt.phase(args[1:])
	case "approval":
		result, revision, err = rt.approval(args[1:])
	case "task":
		result, revision, err = rt.task(args[1:])
	case "intent":
		result, revision, err = rt.intent(args[1:])
	case "document":
		result, revision, err = rt.documents(args[1:])
	case "graph":
		result, revision, err = rt.graph(args[1:])
	case "context":
		result, revision, err = rt.context(args[1:])
	case "loop":
		result, revision, err = rt.loop(args[1:])
	case "waiver":
		result, revision, err = rt.waiver(args[1:])
	case "evidence":
		result, revision, err = rt.evidence(args[1:])
	case "qa":
		result, revision, err = rt.qa(args[1:])
	case "report":
		result, revision, err = rt.report(args[1:])
	case "secret":
		result, revision, err = rt.secrets(args[1:])
	case "migrate":
		result, revision, err = rt.migrate(args[1:])
	case "doctor":
		result, revision, err = rt.doctor(args[1:])
	case "compact":
		result, revision, err = rt.compact()
	case "clear":
		result, revision, err = rt.clear()
	default:
		err = &commandError{"CLI_UNKNOWN_COMMAND", "unknown command: " + args[0], "Run tene-workflow help.", 2}
	}
	if err != nil {
		return rt.fail(err)
	}
	return rt.success(result, revision)
}

func (rt *runtime) master(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("use master create|status|validate")
	}
	p, err := state.New(rt.root).Load()
	if err != nil {
		return nil, 0, err
	}
	view := map[string]any{"project_id": p.ProjectID, "active_sprint_id": p.ActiveSprintID, "revision": p.Revision, "sprints": sortedSprints(p)}
	switch args[0] {
	case "create", "status":
		return view, p.Revision, nil
	case "validate":
		findings := validateMaster(p)
		if workflow.Blocking(findings) {
			return map[string]any{"valid": false, "findings": findings}, p.Revision, &commandError{"MASTER_INVALID", "master plan dependencies are invalid", "Repair missing, self-referential, or cyclic dependencies.", 3}
		}
		return map[string]any{"valid": true, "findings": findings, "master": view}, p.Revision, nil
	default:
		return nil, 0, usageErr("use master create|status|validate")
	}
}

func validateMaster(p *domain.Project) []domain.Finding {
	var fs []domain.Finding
	for id, sp := range p.Sprints {
		for _, d := range sp.Predecessors {
			if d == id || p.Sprints[d] == nil {
				fs = append(fs, domain.Finding{Code: "MASTER_SPRINT_DEP_INVALID", Severity: "blocker", SubjectRefs: []string{id, d}, Message: "sprint predecessor is missing or self-referential", Remediation: "Reference an existing distinct sprint."})
			}
		}
	}
	for id, t := range p.Tasks {
		for _, d := range t.DependsOn {
			if d == id || p.Tasks[d] == nil {
				fs = append(fs, domain.Finding{Code: "MASTER_TASK_DEP_INVALID", Severity: "blocker", SubjectRefs: []string{id, d}, Message: "task dependency is missing or self-referential", Remediation: "Reference an existing distinct task."})
			}
		}
	}
	if cycleMap(p.Sprints, func(s *domain.Sprint) []string { return s.Predecessors }) {
		fs = append(fs, domain.Finding{Code: "MASTER_SPRINT_CYCLE", Severity: "blocker", Message: "sprint predecessor graph contains a cycle", Remediation: "Remove a cyclic predecessor edge."})
	}
	if cycleMap(p.Tasks, func(t *domain.Task) []string { return t.DependsOn }) {
		fs = append(fs, domain.Finding{Code: "MASTER_TASK_CYCLE", Severity: "blocker", Message: "task dependency graph contains a cycle", Remediation: "Remove a cyclic task dependency."})
	}
	return fs
}
func cycleMap[T any](items map[string]*T, deps func(*T) []string) bool {
	state := map[string]int{}
	var visit func(string) bool
	visit = func(id string) bool {
		if state[id] == 1 {
			return true
		}
		if state[id] == 2 {
			return false
		}
		x := items[id]
		if x == nil {
			return false
		}
		state[id] = 1
		for _, d := range deps(x) {
			if visit(d) {
				return true
			}
		}
		state[id] = 2
		return false
	}
	for id := range items {
		if visit(id) {
			return true
		}
	}
	return false
}

func (rt *runtime) route(args []string) (any, uint64, error) {
	fs := flag.NewFlagSet("route", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input := fs.String("text", "", "")
	phase := fs.String("phase", "", "")
	activeFlag := fs.String("active", "auto", "")
	if err := fs.Parse(args); err != nil {
		return nil, 0, err
	}
	if *input == "" {
		*input = strings.Join(fs.Args(), " ")
	}
	if strings.TrimSpace(*input) == "" {
		return nil, 0, usageErr("route --text is required")
	}
	active := false
	ph := domain.Phase(*phase)
	rev := uint64(0)
	if p, err := state.New(rt.root).Load(); err == nil {
		rev = p.Revision
		if p.ActiveSprintID != "" {
			active = true
			if sp := p.Sprints[p.ActiveSprintID]; sp != nil && ph == "" {
				ph = sp.Phase
			}
		}
	}
	if *activeFlag == "true" {
		active = true
	}
	if *activeFlag == "false" {
		active = false
	}
	return router.Route(*input, active, ph), rev, nil
}

func (rt *runtime) global(args []string) []string {
	var out []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			rt.json = true
		case "--root":
			if i+1 < len(args) {
				i++
				rt.root = args[i]
			}
		case "--expected-revision":
			if i+1 < len(args) {
				i++
				if n, err := strconv.ParseUint(args[i], 10, 64); err == nil {
					rt.expected = &n
				}
			}
		case "--request-id":
			if i+1 < len(args) {
				i++
				rt.requestID = args[i]
			}
		default:
			out = append(out, args[i])
		}
	}
	return out
}

func (rt *runtime) success(result any, revision uint64) int {
	if rt.json {
		_ = json.NewEncoder(rt.out).Encode(envelope{OK: true, SchemaVersion: domain.SchemaVersion, RequestID: rt.requestID, Revision: revision, Result: result, Warnings: []string{}, Errors: []apiError{}})
	} else {
		printHuman(rt.out, result)
	}
	return 0
}

func (rt *runtime) fail(err error) int {
	ce := classify(err)
	if rt.json {
		_ = json.NewEncoder(rt.out).Encode(envelope{OK: false, SchemaVersion: domain.SchemaVersion, RequestID: rt.requestID, Warnings: []string{}, Errors: []apiError{{Code: ce.code, Message: ce.message, Remediation: ce.remediation}}})
	} else {
		fmt.Fprintf(rt.err, "%s: %s\n", ce.code, ce.message)
		if ce.remediation != "" {
			fmt.Fprintf(rt.err, "Next: %s\n", ce.remediation)
		}
	}
	return ce.exit
}

func classify(err error) *commandError {
	var ce *commandError
	if errors.As(err, &ce) {
		return ce
	}
	switch {
	case errors.Is(err, state.ErrNotInitialized):
		return &commandError{"STATE_NOT_INITIALIZED", "project is not initialized", "Run tene-workflow init.", 5}
	case errors.Is(err, state.ErrConflict):
		return &commandError{"STATE_REVISION_CONFLICT", err.Error(), "Reload status and retry with the current revision.", 4}
	case errors.Is(err, state.ErrLocked):
		return &commandError{"STATE_LOCKED", err.Error(), "Wait for the other mutation or run doctor.", 4}
	case errors.Is(err, state.ErrCorrupt):
		return &commandError{"STATE_CORRUPT", err.Error(), "Run doctor and restore from a verified backup.", 7}
	}
	msg := err.Error()
	code, exit := "INTERNAL_ERROR", 10
	if i := strings.Index(msg, ":"); i > 0 && strings.Contains(msg[:i], "_") {
		code = msg[:i]
		if strings.HasPrefix(code, "SEC_") {
			exit = 6
		} else {
			exit = 2
		}
		msg = strings.TrimSpace(msg[i+1:])
	}
	return &commandError{code, msg, "", exit}
}

func (rt *runtime) init(args []string) (any, uint64, error) {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	name := fs.String("name", filepath.Base(rt.root), "")
	profile := fs.String("profile", "standard", "")
	if err := fs.Parse(args); err != nil {
		return nil, 0, err
	}
	if !workflow.ValidProfile(*profile) {
		return nil, 0, &commandError{"CFG_PROFILE_INVALID", "profile must be strict, standard, light, or off", "Choose a supported workflow profile.", 2}
	}
	p := domain.NewProject(domain.NewID("project"), *name, *profile, time.Now())
	s := state.New(rt.root)
	if err := s.Initialize(p); err != nil {
		return nil, 0, err
	}
	agents, err := projectconfig.ScaffoldAgents(rt.root)
	if err != nil {
		return nil, 0, err
	}
	created := append([]string{state.DirName}, agents...)
	return map[string]any{"project": p, "created": created}, p.Revision, nil
}

func (rt *runtime) status() (any, uint64, error) {
	p, err := state.New(rt.root).Load()
	if err != nil {
		return nil, 0, err
	}
	var active *domain.Sprint
	if p.ActiveSprintID != "" {
		active = p.Sprints[p.ActiveSprintID]
	}
	counts := map[string]int{"sprints": len(p.Sprints), "tasks": len(p.Tasks), "intents": len(p.Intents), "open_gaps": openGaps(p), "deferred_gaps": deferredGaps(p)}
	result := map[string]any{"project_id": p.ProjectID, "name": p.Name, "profile": p.Profile, "active_sprint": active, "counts": counts}
	if active != nil {
		result["workflow"] = map[string]any{"effective_blockers": effectiveBlockers(p, active), "active_waivers": activeWaiverCount(p, active), "loop_iteration": active.LoopIteration, "max_loop_iterations": active.MaxLoopIterations, "loop_remaining": max(0, active.MaxLoopIterations-active.LoopIteration), "last_loop_outcome": active.LastLoopOutcome}
	}
	return result, p.Revision, nil
}

func (rt *runtime) sprint(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("sprint action required")
	}
	s := state.New(rt.root)
	p, err := s.Load()
	if err != nil {
		return nil, 0, err
	}
	switch args[0] {
	case "create":
		fs := flag.NewFlagSet("sprint create", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		title := fs.String("title", "", "")
		slug := fs.String("slug", "", "")
		pred := fs.String("predecessors", "", "")
		maxIterations := fs.Int("max-iterations", 5, "")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, 0, err
		}
		if *title == "" {
			return nil, 0, usageErr("--title is required")
		}
		if *maxIterations < 1 || *maxIterations > 20 {
			return nil, 0, usageErr("--max-iterations must be between 1 and 20")
		}
		if *slug == "" {
			*slug = slugify(*title)
		}
		id := domain.NewID("sprint")
		sp := &domain.Sprint{SprintID: id, Slug: *slug, Title: *title, Phase: domain.PhaseDraft, Predecessors: csv(*pred), DocumentRoot: filepath.ToSlash(filepath.Join("docs", "sprints", id+"-"+*slug)), MaxLoopIterations: *maxIterations}
		result, err := s.Mutate(rt.expected, actor(), "SprintCreated", id, sp, func(p *domain.Project) error { p.Sprints[id] = sp; p.ActiveSprintID = id; return nil })
		if err != nil {
			return nil, 0, err
		}
		paths, err := document.ScaffoldAll(rt.root, result, sp)
		if err != nil {
			return nil, 0, err
		}
		return map[string]any{"sprint": sp, "documents": paths}, result.Revision, nil
	case "list":
		list := sortedSprints(p)
		return list, p.Revision, nil
	case "master-plan":
		return map[string]any{"project_id": p.ProjectID, "active_sprint_id": p.ActiveSprintID, "revision": p.Revision, "sprints": sortedSprints(p)}, p.Revision, nil
	case "show":
		sp, err := selectSprint(p, arg(args, 1))
		return sp, p.Revision, err
	case "start":
		id := arg(args, 1)
		sp, err := selectSprint(p, id)
		if err != nil {
			return nil, 0, err
		}
		now := time.Now().UTC()
		result, err := s.Mutate(rt.expected, actor(), "SprintStarted", sp.SprintID, map[string]string{"sprint_id": sp.SprintID}, func(p *domain.Project) error {
			p.ActiveSprintID = sp.SprintID
			if p.Sprints[sp.SprintID].StartedAt == nil {
				p.Sprints[sp.SprintID].StartedAt = &now
			}
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		return result.Sprints[sp.SprintID], result.Revision, nil
	case "archive":
		fs := flag.NewFlagSet("sprint archive", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		approvalID := fs.String("approval", "", "")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, 0, err
		}
		return rt.transition(domain.PhaseArchived, false, *approvalID)
	default:
		return nil, 0, usageErr("unknown sprint action")
	}
}

func (rt *runtime) phase(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("phase action required")
	}
	p, err := state.New(rt.root).Load()
	if err != nil {
		return nil, 0, err
	}
	sp, err := activeSprint(p)
	if err != nil {
		return nil, 0, err
	}
	if args[0] == "show" {
		return map[string]any{"sprint_id": sp.SprintID, "phase": sp.Phase}, p.Revision, nil
	}
	if args[0] != "transition" || len(args) < 2 {
		return nil, 0, usageErr("use phase transition <phase> [--dry-run]")
	}
	target, err := workflow.ParsePhase(args[1])
	if err != nil {
		return nil, 0, err
	}
	fs := flag.NewFlagSet("phase transition", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	dry := fs.Bool("dry-run", false, "")
	approvalID := fs.String("approval", "", "")
	if err := fs.Parse(args[2:]); err != nil {
		return nil, 0, err
	}
	return rt.transition(target, *dry, *approvalID)
}

func (rt *runtime) transition(target domain.Phase, dry bool, approvalID string) (any, uint64, error) {
	s := state.New(rt.root)
	p, err := s.Load()
	if err != nil {
		return nil, 0, err
	}
	sp, err := activeSprint(p)
	if err != nil {
		return nil, 0, err
	}
	now := time.Now().UTC()
	findings := workflow.CanTransitionWithApproval(p, sp, target, approvalID, now, func(ph domain.Phase) bool {
		if !document.Exists(rt.root, sp, ph) {
			return false
		}
		return !workflow.Blocking(document.Validate(document.Path(rt.root, sp, ph), ph))
	})
	if workflow.Blocking(findings) {
		code := "WF_GUARD_FAILED"
		message := fmt.Sprintf("%d blocking guard(s) failed", len(findings))
		remediation := "Resolve the listed findings and retry."
		if len(findings) == 1 && strings.HasPrefix(findings[0].Code, "WF_APPROVAL_") {
			code = findings[0].Code
			message = findings[0].Message
			remediation = findings[0].Remediation
		}
		return map[string]any{"allowed": false, "findings": findings}, p.Revision, &commandError{code, message, remediation, 3}
	}
	if dry {
		return map[string]any{"allowed": true, "from": sp.Phase, "to": target, "findings": findings}, p.Revision, nil
	}
	from := sp.Phase
	oldDocumentRoot := sp.DocumentRoot
	newDocumentRoot := oldDocumentRoot
	var oldPath, newPath string
	if target == domain.PhaseArchived {
		newDocumentRoot = filepath.ToSlash(filepath.Join("docs", "sprints", "_archive", now.Format("2006-01"), filepath.Base(oldDocumentRoot)))
		oldPath, newPath = filepath.Join(rt.root, oldDocumentRoot), filepath.Join(rt.root, newDocumentRoot)
		if err := os.MkdirAll(filepath.Dir(newPath), 0o755); err != nil {
			return nil, 0, err
		}
		if err := os.Rename(oldPath, newPath); err != nil {
			return nil, 0, fmt.Errorf("STATE_ARCHIVE_MOVE_FAILED: %w", err)
		}
		manifest := map[string]any{"schema_version": domain.SchemaVersion, "sprint_id": sp.SprintID, "archived_at": now, "state_revision": p.Revision + 1, "qa_run_id": sp.LastQAID, "qa_status": sp.LastQAStatus, "report_path": newDocumentRoot + strings.TrimPrefix(sp.ReportPath, oldDocumentRoot)}
		manifestBytes, marshalErr := json.MarshalIndent(manifest, "", "  ")
		if marshalErr != nil {
			_ = os.Rename(newPath, oldPath)
			return nil, 0, marshalErr
		}
		manifestPath := filepath.Join(newPath, "99-archive", "archive-manifest.json")
		if err := os.MkdirAll(filepath.Dir(manifestPath), 0o755); err != nil {
			_ = os.Rename(newPath, oldPath)
			return nil, 0, err
		}
		if err := os.WriteFile(manifestPath, append(manifestBytes, '\n'), 0o644); err != nil {
			_ = os.Rename(newPath, oldPath)
			return nil, 0, err
		}
	}
	result, err := s.Mutate(rt.expected, actor(), "PhaseTransitioned", sp.SprintID, map[string]any{"from": from, "to": target}, func(p *domain.Project) error {
		x := p.Sprints[sp.SprintID]
		x.Phase = target
		if approvalID != "" && workflow.RequiredApproval(p.Profile, from, target) {
			a := p.Approvals[approvalID]
			a.Status = "consumed"
			a.ConsumedAt = &now
			x.ApprovalRefs = unique(append(x.ApprovalRefs, approvalID))
		}
		if target == domain.PhaseArchived {
			x.ArchivedAt = &now
			x.DocumentRoot = newDocumentRoot
			if strings.HasPrefix(x.ReportPath, oldDocumentRoot+"/") {
				x.ReportPath = newDocumentRoot + strings.TrimPrefix(x.ReportPath, oldDocumentRoot)
			}
			for _, evidence := range p.Evidence {
				if evidence.SprintID == x.SprintID && strings.HasPrefix(evidence.URI, oldDocumentRoot+"/") {
					evidence.URI = newDocumentRoot + strings.TrimPrefix(evidence.URI, oldDocumentRoot)
				}
			}
			for id, node := range p.Graph.Nodes {
				if strings.HasPrefix(node.Locator, oldDocumentRoot+"/") {
					node.Locator = newDocumentRoot + strings.TrimPrefix(node.Locator, oldDocumentRoot)
					p.Graph.Nodes[id] = node
				}
			}
			p.ActiveSprintID = ""
		}
		return nil
	})
	if err != nil {
		if target == domain.PhaseArchived {
			_ = os.Rename(newPath, oldPath)
		}
		return nil, 0, err
	}
	return map[string]any{"from": from, "to": target, "sprint": result.Sprints[sp.SprintID]}, result.Revision, nil
}

func (rt *runtime) approval(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("approval action required")
	}
	s := state.New(rt.root)
	p, err := s.Load()
	if err != nil {
		return nil, 0, err
	}
	sp, err := activeSprint(p)
	if err != nil {
		return nil, 0, err
	}
	switch args[0] {
	case "list":
		var out []*domain.Approval
		for _, a := range p.Approvals {
			if a.SprintID == sp.SprintID {
				out = append(out, a)
			}
		}
		sort.Slice(out, func(i, j int) bool { return out[i].ApprovalID < out[j].ApprovalID })
		return out, p.Revision, nil
	case "request":
		fs := flag.NewFlagSet("approval request", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		fromText := fs.String("from", string(sp.Phase), "")
		toText := fs.String("to", "", "")
		reason := fs.String("reason", "", "")
		requester := fs.String("requester", "", "")
		expires := fs.String("expires", "", "")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, 0, err
		}
		if *toText == "" || *reason == "" || *requester == "" || *expires == "" {
			return nil, 0, usageErr("--to --reason --requester --expires are required")
		}
		from, err := workflow.ParsePhase(*fromText)
		if err != nil {
			return nil, 0, err
		}
		to, err := workflow.ParsePhase(*toText)
		if err != nil {
			return nil, 0, err
		}
		if from != sp.Phase {
			return nil, p.Revision, &commandError{"WF_APPROVAL_SCOPE_MISMATCH", "approval from phase is not current sprint phase", "Request approval for the active transition.", 3}
		}
		expiry, err := time.Parse(time.RFC3339, *expires)
		if err != nil || !expiry.After(time.Now().UTC()) {
			return nil, p.Revision, &commandError{"WF_APPROVAL_EXPIRY_INVALID", "expiry must be a future RFC3339 timestamp", "Provide a bounded future expiry.", 2}
		}
		if !workflow.RequiredApproval(p.Profile, from, to) {
			return nil, p.Revision, &commandError{"WF_APPROVAL_NOT_REQUIRED", "this profile transition does not require approval", "Run phase transition without an approval.", 2}
		}
		id := domain.NewID("approval")
		a := &domain.Approval{ApprovalID: id, SprintID: sp.SprintID, From: from, To: to, Reason: *reason, Requester: *requester, Status: "requested", RequestedAt: time.Now().UTC(), ExpiresAt: expiry.UTC()}
		r, err := s.Mutate(rt.expected, actor(), "ApprovalRequested", id, a, func(p *domain.Project) error { p.Approvals[id] = a; return nil })
		if err != nil {
			return nil, 0, err
		}
		return r.Approvals[id], r.Revision, nil
	case "approve":
		if len(args) < 2 {
			return nil, 0, usageErr("approval id required")
		}
		id := args[1]
		a := p.Approvals[id]
		if a == nil || a.SprintID != sp.SprintID {
			return nil, 0, notFound("approval", id)
		}
		fs := flag.NewFlagSet("approval approve", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		approver := fs.String("approver", "", "")
		if err := fs.Parse(args[2:]); err != nil {
			return nil, 0, err
		}
		if *approver == "" {
			return nil, 0, usageErr("--approver is required")
		}
		if a.Status != "requested" {
			return nil, p.Revision, &commandError{"WF_APPROVAL_INVALID", "only requested approvals can be approved", "Create a new approval request.", 3}
		}
		now := time.Now().UTC()
		if !a.ExpiresAt.After(now) {
			return nil, p.Revision, &commandError{"WF_APPROVAL_EXPIRED", "approval request has expired", "Create a new bounded approval request.", 3}
		}
		r, err := s.Mutate(rt.expected, actor(), "ApprovalGranted", id, map[string]string{"approver": *approver}, func(p *domain.Project) error {
			x := p.Approvals[id]
			x.Status = "approved"
			x.Approver = *approver
			x.ApprovedAt = &now
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		return r.Approvals[id], r.Revision, nil
	default:
		return nil, 0, usageErr("unknown approval action")
	}
}

func (rt *runtime) task(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("task action required")
	}
	s := state.New(rt.root)
	p, err := s.Load()
	if err != nil {
		return nil, 0, err
	}
	sp, err := activeSprint(p)
	if err != nil {
		return nil, 0, err
	}
	switch args[0] {
	case "list":
		var out []*domain.Task
		for _, id := range sp.TaskIDs {
			if t := p.Tasks[id]; t != nil {
				out = append(out, t)
			}
		}
		return out, p.Revision, nil
	case "add":
		fs := flag.NewFlagSet("task add", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		title := fs.String("title", "", "")
		layer := fs.String("layer", "business", "")
		acs := fs.String("ac", "", "")
		deps := fs.String("depends-on", "", "")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, 0, err
		}
		if *title == "" {
			return nil, 0, usageErr("--title is required")
		}
		id := domain.NewID("task")
		t := &domain.Task{TaskID: id, SprintID: sp.SprintID, Title: *title, Status: "todo", Layer: *layer, CriterionIDs: csv(*acs), DependsOn: csv(*deps)}
		for _, ac := range t.CriterionIDs {
			if p.Criteria[ac] == nil {
				return nil, 0, notFound("acceptance criterion", ac)
			}
		}
		for _, dep := range t.DependsOn {
			if p.Tasks[dep] == nil {
				return nil, 0, notFound("task dependency", dep)
			}
		}
		r, err := s.Mutate(rt.expected, actor(), "TaskChanged", id, t, func(p *domain.Project) error {
			p.Tasks[id] = t
			p.Sprints[sp.SprintID].TaskIDs = append(p.Sprints[sp.SprintID].TaskIDs, id)
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		return r.Tasks[id], r.Revision, nil
	case "start", "complete", "block", "defer":
		if len(args) < 2 {
			return nil, 0, usageErr("task id is required")
		}
		id := args[1]
		if p.Tasks[id] == nil {
			return nil, 0, notFound("task", id)
		}
		status := map[string]string{"start": "doing", "complete": "done", "block": "blocked", "defer": "deferred"}[args[0]]
		r, err := s.Mutate(rt.expected, actor(), "TaskChanged", id, map[string]string{"status": status}, func(p *domain.Project) error { p.Tasks[id].Status = status; return nil })
		if err != nil {
			return nil, 0, err
		}
		return r.Tasks[id], r.Revision, nil
	case "link":
		if len(args) < 2 {
			return nil, 0, usageErr("task id is required")
		}
		id := args[1]
		if p.Tasks[id] == nil {
			return nil, 0, notFound("task", id)
		}
		fs := flag.NewFlagSet("task link", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		acs := fs.String("ac", "", "")
		intents := fs.String("intent", "", "")
		replace := fs.Bool("replace", false, "")
		if err := fs.Parse(args[2:]); err != nil {
			return nil, 0, err
		}
		if *acs == "" && *intents == "" {
			return nil, 0, usageErr("--ac or --intent is required")
		}
		for _, ac := range csv(*acs) {
			if p.Criteria[ac] == nil {
				return nil, 0, notFound("acceptance criterion", ac)
			}
		}
		for _, intent := range csv(*intents) {
			if p.Intents[intent] == nil {
				return nil, 0, notFound("intent", intent)
			}
		}
		r, err := s.Mutate(rt.expected, actor(), "TaskLinked", id, map[string]any{"ac_ids": csv(*acs), "intent_ids": csv(*intents)}, func(p *domain.Project) error {
			task := p.Tasks[id]
			if *replace {
				task.CriterionIDs = csv(*acs)
				task.IntentIDs = csv(*intents)
			} else {
				task.CriterionIDs = unique(append(task.CriterionIDs, csv(*acs)...))
				task.IntentIDs = unique(append(task.IntentIDs, csv(*intents)...))
			}
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		return r.Tasks[id], r.Revision, nil
	default:
		return nil, 0, usageErr("unknown task action")
	}
}

func (rt *runtime) intent(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("intent action required")
	}
	s := state.New(rt.root)
	p, err := s.Load()
	if err != nil {
		return nil, 0, err
	}
	sp, err := activeSprint(p)
	if err != nil {
		return nil, 0, err
	}
	switch args[0] {
	case "list":
		var out []*domain.Intent
		for _, id := range sp.IntentIDs {
			if x := p.Intents[id]; x != nil {
				out = append(out, x)
			}
		}
		return out, p.Revision, nil
	case "capture":
		fs := flag.NewFlagSet("intent capture", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		statement := fs.String("statement", "", "")
		rationale := fs.String("rationale", "", "")
		acText := fs.String("ac", "", "")
		observable := fs.String("observable", "", "")
		priority := fs.String("priority", "blocking", "")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, 0, err
		}
		if *statement == "" {
			return nil, 0, usageErr("--statement is required")
		}
		id := domain.NewID("intent")
		in := &domain.Intent{IntentID: id, SprintID: sp.SprintID, Revision: 1, Status: "candidate", Statement: *statement, Rationale: *rationale, Source: "user"}
		var ac *domain.Criterion
		if *acText != "" {
			ac = &domain.Criterion{CriterionID: domain.NewID("ac"), IntentID: id, Statement: *acText, Observable: *observable, Priority: *priority}
		}
		r, err := s.Mutate(rt.expected, actor(), "IntentCaptured", id, map[string]any{"intent": in, "criterion": ac}, func(p *domain.Project) error {
			p.Intents[id] = in
			p.Sprints[sp.SprintID].IntentIDs = append(p.Sprints[sp.SprintID].IntentIDs, id)
			if ac != nil {
				p.Criteria[ac.CriterionID] = ac
			}
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		return map[string]any{"intent": r.Intents[id], "criterion": ac}, r.Revision, nil
	case "confirm", "deprecate":
		if len(args) < 2 {
			return nil, 0, usageErr("intent id required")
		}
		id := args[1]
		if p.Intents[id] == nil {
			return nil, 0, notFound("intent", id)
		}
		status := map[string]string{"confirm": "confirmed", "deprecate": "deprecated"}[args[0]]
		now := time.Now().UTC()
		r, err := s.Mutate(rt.expected, actor(), "IntentChanged", id, map[string]string{"status": status}, func(p *domain.Project) error {
			in := p.Intents[id]
			in.Status = status
			in.Revision++
			if status == "confirmed" {
				in.ConfirmedAt = &now
				in.ConfirmedBy = "user"
			}
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		return r.Intents[id], r.Revision, nil
	case "revise":
		if len(args) < 2 {
			return nil, 0, usageErr("intent id required")
		}
		id := args[1]
		if p.Intents[id] == nil {
			return nil, 0, notFound("intent", id)
		}
		fs := flag.NewFlagSet("intent revise", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		statement := fs.String("statement", "", "")
		rationale := fs.String("rationale", "", "")
		if err := fs.Parse(args[2:]); err != nil {
			return nil, 0, err
		}
		if *statement == "" && *rationale == "" {
			return nil, 0, usageErr("--statement or --rationale is required")
		}
		r, err := s.Mutate(rt.expected, actor(), "IntentRevised", id, map[string]string{"statement": *statement, "rationale": *rationale}, func(p *domain.Project) error {
			in := p.Intents[id]
			if *statement != "" {
				in.Statement = *statement
			}
			if *rationale != "" {
				in.Rationale = *rationale
			}
			in.Revision++
			in.Status = "candidate"
			in.ConfirmedAt = nil
			in.ConfirmedBy = ""
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		return r.Intents[id], r.Revision, nil
	case "add-ac":
		if len(args) < 2 {
			return nil, 0, usageErr("intent id required")
		}
		id := args[1]
		if p.Intents[id] == nil {
			return nil, 0, notFound("intent", id)
		}
		fs := flag.NewFlagSet("intent add-ac", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		statement := fs.String("statement", "", "")
		observable := fs.String("observable", "", "")
		priority := fs.String("priority", "blocking", "")
		if err := fs.Parse(args[2:]); err != nil {
			return nil, 0, err
		}
		if *statement == "" || *observable == "" {
			return nil, 0, usageErr("--statement and --observable are required")
		}
		if *priority != "blocking" && *priority != "non-blocking" {
			return nil, 0, usageErr("priority must be blocking or non-blocking")
		}
		ac := &domain.Criterion{CriterionID: domain.NewID("ac"), IntentID: id, Statement: *statement, Observable: *observable, Priority: *priority}
		r, err := s.Mutate(rt.expected, actor(), "CriterionAdded", ac.CriterionID, ac, func(p *domain.Project) error { p.Criteria[ac.CriterionID] = ac; return nil })
		if err != nil {
			return nil, 0, err
		}
		return r.Criteria[ac.CriterionID], r.Revision, nil
	default:
		return nil, 0, usageErr("unknown intent action")
	}
}

func (rt *runtime) documents(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("document action required")
	}
	p, err := state.New(rt.root).Load()
	if err != nil {
		return nil, 0, err
	}
	sp, err := activeSprint(p)
	if err != nil {
		return nil, 0, err
	}
	switch args[0] {
	case "scaffold":
		paths, err := document.ScaffoldAll(rt.root, p, sp)
		return map[string]any{"created": paths}, p.Revision, err
	case "validate":
		var findings []domain.Finding
		phases := []domain.Phase{domain.PhasePRD, domain.PhasePlan, domain.PhaseDesign, domain.PhaseLoopCheck, domain.PhaseQA, domain.PhaseReport}
		if len(args) > 1 {
			phase, parseErr := workflow.ParsePhase(args[1])
			if parseErr != nil {
				return nil, 0, parseErr
			}
			phases = []domain.Phase{phase}
		}
		for _, ph := range phases {
			findings = append(findings, document.Validate(document.Path(rt.root, sp, ph), ph)...)
		}
		if workflow.Blocking(findings) {
			return map[string]any{"valid": false, "findings": findings}, p.Revision, &commandError{"DOC_VALIDATION_FAILED", "one or more documents are invalid", "Repair the listed sections.", 3}
		}
		return map[string]any{"valid": true, "findings": findings}, p.Revision, nil
	default:
		return nil, 0, usageErr("unknown document action")
	}
}

func (rt *runtime) graph(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("graph action required")
	}
	s := state.New(rt.root)
	p, err := s.Load()
	if err != nil {
		return nil, 0, err
	}
	switch args[0] {
	case "providers":
		return codeintel.Discover(rt.root), p.Revision, nil
	case "build":
		report, analyzeErr := codeintel.Analyze(context.Background(), rt.root, nil, false)
		if analyzeErr != nil {
			return nil, 0, analyzeErr
		}
		r, err := s.Mutate(rt.expected, actor(), "GraphRebuilt", p.ProjectID, map[string]string{"provider": "workflow-state+codeintel"}, func(p *domain.Project) error { p.Graph = buildGraph(p); mergeCodeGraph(&p.Graph, report); return nil })
		if err != nil {
			return nil, 0, err
		}
		return map[string]int{"nodes": len(r.Graph.Nodes), "edges": len(r.Graph.Edges)}, r.Revision, nil
	case "understand":
		fs := flag.NewFlagSet("graph understand", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		changed := fs.Bool("changed", false, "")
		paths := fs.String("path", "", "")
		query := fs.String("codegraph-query", "", "")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, 0, err
		}
		report, err := codeintel.Analyze(context.Background(), rt.root, csv(*paths), *changed)
		if err == nil && *query != "" {
			report.SemanticContext, err = codeintel.ExploreCodeGraph(context.Background(), rt.root, *query)
		}
		return report, p.Revision, err
	case "trace":
		if len(args) < 2 {
			return nil, 0, usageErr("node id required")
		}
		return traceGraph(p, args[1]), p.Revision, nil
	case "impact":
		if len(args) < 2 {
			return nil, 0, usageErr("use graph impact <node-id> [--depth N] [--call-depth N]")
		}
		fs := flag.NewFlagSet("graph impact", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		depth := fs.Int("depth", 8, "")
		callDepth := fs.Int("call-depth", 4, "")
		if err := fs.Parse(args[2:]); err != nil {
			return nil, 0, err
		}
		result, err := tracecontext.Impact(p.Graph, args[1], *depth, *callDepth)
		if err != nil {
			return nil, 0, &commandError{"GRAPH_IMPACT_INVALID", err.Error(), "Build the graph and provide valid traversal limits.", 2}
		}
		return result, p.Revision, nil
	case "validate":
		findings := validateGraph(p)
		if workflow.Blocking(findings) {
			return map[string]any{"valid": false, "findings": findings}, p.Revision, &commandError{"GRAPH_INVARIANT_FAILED", "graph invariants failed", "Build missing traceability edges.", 3}
		}
		return map[string]any{"valid": true, "findings": findings}, p.Revision, nil
	default:
		return nil, 0, usageErr("unknown graph action")
	}
}

func (rt *runtime) context(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("use context build|validate")
	}
	p, err := state.New(rt.root).Load()
	if err != nil {
		return nil, 0, err
	}
	switch args[0] {
	case "build":
		sp, err := activeSprint(p)
		if err != nil {
			return nil, 0, err
		}
		fs := flag.NewFlagSet("context build", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		phaseText := fs.String("phase", string(sp.Phase), "")
		budget := fs.Int("budget", 32768, "")
		output := fs.String("output", "", "")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, 0, err
		}
		phase, err := workflow.ParsePhase(*phaseText)
		if err != nil {
			return nil, 0, err
		}
		pack, err := tracecontext.BuildContextPack(rt.root, p, sp, tracecontext.BuildOptions{Phase: phase, Budget: *budget}, codeintel.Discover(rt.root))
		if err != nil {
			return nil, p.Revision, &commandError{"CONTEXT_BUILD_FAILED", err.Error(), "Increase the budget or repair the referenced workflow state.", 3}
		}
		if *output != "" {
			if err := writeDerivedJSON(rt.root, *output, pack); err != nil {
				return nil, p.Revision, err
			}
		}
		return pack, p.Revision, nil
	case "validate":
		fs := flag.NewFlagSet("context validate", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		input := fs.String("input", "", "")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, 0, err
		}
		if *input == "" {
			return nil, 0, usageErr("--input is required")
		}
		path, err := safeRootPath(rt.root, *input)
		if err != nil {
			return nil, p.Revision, err
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, p.Revision, err
		}
		var pack tracecontext.ContextPack
		if err := json.Unmarshal(b, &pack); err != nil {
			return nil, p.Revision, &commandError{"CONTEXT_PACK_INVALID", err.Error(), "Rebuild the context pack.", 2}
		}
		result := tracecontext.ValidateContextPack(rt.root, p, pack)
		if !result.Fresh {
			return result, p.Revision, &commandError{"CONTEXT_STALE", "context pack no longer matches state or provenance", "Run context build again before mutation.", 3}
		}
		return result, p.Revision, nil
	default:
		return nil, 0, usageErr("use context build|validate")
	}
}

func (rt *runtime) loop(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("loop action required")
	}
	s := state.New(rt.root)
	p, err := s.Load()
	if err != nil {
		return nil, 0, err
	}
	sp, err := activeSprint(p)
	if err != nil {
		return nil, 0, err
	}
	switch args[0] {
	case "check":
		var gaps []*domain.Gap
		for _, id := range sp.OpenGapIDs {
			if g := p.Gaps[id]; g != nil && g.Status == "open" {
				gaps = append(gaps, g)
			}
		}
		return map[string]any{"passed": len(gaps) == 0, "open_gaps": gaps, "iteration": sp.LoopIteration, "max_iterations": sp.MaxLoopIterations, "remaining": max(0, sp.MaxLoopIterations-sp.LoopIteration), "exhausted": sp.LoopIteration >= sp.MaxLoopIterations}, p.Revision, nil
	case "iterate":
		if sp.Phase != domain.PhaseLoopCheck {
			return nil, p.Revision, &commandError{"WF_LOOP_PHASE_REQUIRED", "loop iterations can only be recorded in loop-check", "Transition to loop-check first.", 3}
		}
		fs := flag.NewFlagSet("loop iterate", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		outcome := fs.String("outcome", "repair", "")
		summary := fs.String("summary", "", "")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, 0, err
		}
		if !slices.Contains([]string{"repair", "passed", "blocked"}, *outcome) || *summary == "" {
			return nil, 0, usageErr("--outcome repair|passed|blocked and --summary are required")
		}
		if sp.LoopIteration >= sp.MaxLoopIterations {
			return nil, p.Revision, &commandError{"WF_LOOP_EXHAUSTED", "loop iteration budget is exhausted", "Leave unresolved gaps visible and request a policy decision.", 3}
		}
		now := time.Now().UTC()
		r, err := s.Mutate(rt.expected, actor(), "LoopIterated", sp.SprintID, map[string]any{"outcome": *outcome, "summary": *summary, "iteration": sp.LoopIteration + 1}, func(p *domain.Project) error {
			x := p.Sprints[sp.SprintID]
			x.LoopIteration++
			x.LastLoopOutcome = *outcome
			x.LastLoopSummary = *summary
			x.LastLoopAt = &now
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		x := r.Sprints[sp.SprintID]
		return map[string]any{"iteration": x.LoopIteration, "max_iterations": x.MaxLoopIterations, "remaining": max(0, x.MaxLoopIterations-x.LoopIteration), "outcome": x.LastLoopOutcome, "exhausted": x.LoopIteration >= x.MaxLoopIterations}, r.Revision, nil
	case "record-gap":
		fs := flag.NewFlagSet("loop record-gap", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		desc := fs.String("description", "", "")
		category := fs.String("category", "mismatch", "")
		severity := fs.String("severity", "blocker", "")
		subjects := fs.String("subject", "", "")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, 0, err
		}
		if *desc == "" {
			return nil, 0, usageErr("--description required")
		}
		if !slices.Contains([]string{"missing", "mismatch", "unverified", "regression", "debt", "security", "evidence-integrity"}, *category) {
			return nil, 0, usageErr("unsupported gap category")
		}
		if !slices.Contains([]string{"blocker", "warning"}, *severity) {
			return nil, 0, usageErr("severity must be blocker or warning")
		}
		id := domain.NewID("gap")
		g := &domain.Gap{GapID: id, SprintID: sp.SprintID, Category: *category, Severity: *severity, Status: "open", Description: *desc, SubjectRefs: csv(*subjects)}
		r, err := s.Mutate(rt.expected, actor(), "GapRecorded", id, g, func(p *domain.Project) error {
			p.Gaps[id] = g
			p.Sprints[sp.SprintID].OpenGapIDs = append(p.Sprints[sp.SprintID].OpenGapIDs, id)
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		return r.Gaps[id], r.Revision, nil
	case "resolve-gap":
		if len(args) < 2 {
			return nil, 0, usageErr("gap id required")
		}
		id := args[1]
		fs := flag.NewFlagSet("loop resolve-gap", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		resolution := fs.String("resolution", "", "")
		evidenceIDs := fs.String("evidence", "", "")
		crossSprint := fs.Bool("cross-sprint", false, "")
		if err := fs.Parse(args[2:]); err != nil {
			return nil, 0, err
		}
		if *resolution == "" || *evidenceIDs == "" {
			return nil, 0, usageErr("--resolution and --evidence are required")
		}
		if p.Gaps[id] == nil || (p.Gaps[id].SprintID != sp.SprintID && !*crossSprint) {
			return nil, 0, notFound("gap", id)
		}
		if p.Gaps[id].SprintID != sp.SprintID && (p.Gaps[id].Category != "debt" || p.Sprints[p.Gaps[id].SprintID] == nil || p.Sprints[p.Gaps[id].SprintID].Phase != domain.PhaseArchived) {
			return nil, p.Revision, &commandError{"WF_CROSS_SPRINT_GAP_FORBIDDEN", "only inherited debt from an archived sprint can be resolved cross-sprint", "Resolve current work in its owning sprint or archive the predecessor first.", 3}
		}
		if p.Gaps[id].Status != "open" {
			return nil, p.Revision, &commandError{"WF_GAP_NOT_OPEN", "only open gaps can be resolved", "Inspect the current gap disposition.", 3}
		}
		for _, evidenceID := range csv(*evidenceIDs) {
			if p.Evidence[evidenceID] == nil {
				return nil, 0, notFound("evidence", evidenceID)
			}
		}
		r, err := s.Mutate(rt.expected, actor(), "GapResolved", id, map[string]any{"resolution": *resolution, "evidence_ids": csv(*evidenceIDs)}, func(p *domain.Project) error {
			p.Gaps[id].Status = "resolved"
			p.Gaps[id].Resolution = *resolution
			p.Gaps[id].ResolutionEvidenceIDs = csv(*evidenceIDs)
			p.Sprints[p.Gaps[id].SprintID].OpenGapIDs = remove(p.Sprints[p.Gaps[id].SprintID].OpenGapIDs, id)
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		return r.Gaps[id], r.Revision, nil
	case "defer-gap":
		if len(args) < 2 {
			return nil, 0, usageErr("gap id required")
		}
		id := args[1]
		g := p.Gaps[id]
		if g == nil || g.SprintID != sp.SprintID {
			return nil, 0, notFound("gap", id)
		}
		if g.Status != "open" {
			return nil, p.Revision, &commandError{"WF_GAP_NOT_OPEN", "only open gaps can be deferred", "Inspect the current gap disposition.", 3}
		}
		if g.Category == "security" || g.Category == "evidence-integrity" {
			return nil, p.Revision, &commandError{"WF_GAP_DEFER_FORBIDDEN", "security and evidence-integrity gaps cannot be deferred", "Resolve the gap with valid evidence.", 3}
		}
		fs := flag.NewFlagSet("loop defer-gap", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		reason := fs.String("reason", "", "")
		owner := fs.String("owner", "", "")
		target := fs.String("target-sprint", "", "")
		if err := fs.Parse(args[2:]); err != nil {
			return nil, 0, err
		}
		if *reason == "" || *owner == "" || *target == "" {
			return nil, 0, usageErr("--reason --owner --target-sprint are required")
		}
		now := time.Now().UTC()
		r, err := s.Mutate(rt.expected, actor(), "GapDeferred", id, map[string]string{"reason": *reason, "owner": *owner, "target_sprint": *target}, func(p *domain.Project) error {
			x := p.Gaps[id]
			x.Status = "deferred"
			x.DeferredReason = *reason
			x.DeferredOwner = *owner
			x.DeferredTargetSprint = *target
			x.DeferredAt = &now
			p.Sprints[sp.SprintID].OpenGapIDs = remove(p.Sprints[sp.SprintID].OpenGapIDs, id)
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		return r.Gaps[id], r.Revision, nil
	default:
		return nil, 0, usageErr("unknown loop action")
	}
}

func (rt *runtime) waiver(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("waiver action required")
	}
	s := state.New(rt.root)
	p, err := s.Load()
	if err != nil {
		return nil, 0, err
	}
	sp, err := activeSprint(p)
	if err != nil {
		return nil, 0, err
	}
	switch args[0] {
	case "list":
		var out []*domain.Waiver
		for _, w := range p.Waivers {
			if w.SprintID == sp.SprintID {
				out = append(out, w)
			}
		}
		sort.Slice(out, func(i, j int) bool { return out[i].WaiverID < out[j].WaiverID })
		return out, p.Revision, nil
	case "create":
		fs := flag.NewFlagSet("waiver create", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		gapID := fs.String("gap", "", "")
		reason := fs.String("reason", "", "")
		approver := fs.String("approver", "", "")
		expires := fs.String("expires", "", "")
		scope := fs.String("scope", "phase-transition", "")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, 0, err
		}
		if *gapID == "" || *reason == "" || *approver == "" || *expires == "" {
			return nil, 0, usageErr("--gap --reason --approver --expires are required")
		}
		gap := p.Gaps[*gapID]
		if gap == nil || gap.SprintID != sp.SprintID {
			return nil, 0, notFound("gap", *gapID)
		}
		if gap.Category == "security" || gap.Category == "evidence-integrity" {
			return nil, 0, &commandError{"WAIVER_FORBIDDEN", "security and evidence-integrity gaps cannot be waived", "Resolve the gap with valid evidence.", 3}
		}
		expiry, parseErr := time.Parse(time.RFC3339, *expires)
		if parseErr != nil || !expiry.After(time.Now().UTC()) {
			return nil, 0, &commandError{"WAIVER_EXPIRY_INVALID", "expiry must be a future RFC3339 timestamp", "Provide a bounded future expiry.", 2}
		}
		id := domain.NewID("waiver")
		w := &domain.Waiver{WaiverID: id, SprintID: sp.SprintID, GapID: *gapID, Reason: *reason, Scope: *scope, Approver: *approver, Status: "active", CreatedAt: time.Now().UTC(), ExpiresAt: expiry.UTC()}
		r, err := s.Mutate(rt.expected, actor(), "WaiverCreated", id, w, func(p *domain.Project) error { p.Waivers[id] = w; return nil })
		if err != nil {
			return nil, 0, err
		}
		return r.Waivers[id], r.Revision, nil
	case "revoke":
		if len(args) < 2 {
			return nil, 0, usageErr("waiver id required")
		}
		id := args[1]
		if p.Waivers[id] == nil {
			return nil, 0, notFound("waiver", id)
		}
		now := time.Now().UTC()
		r, err := s.Mutate(rt.expected, actor(), "WaiverRevoked", id, nil, func(p *domain.Project) error {
			p.Waivers[id].Status = "revoked"
			p.Waivers[id].RevokedAt = &now
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		return r.Waivers[id], r.Revision, nil
	default:
		return nil, 0, usageErr("unknown waiver action")
	}
}

func (rt *runtime) evidence(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("evidence action required")
	}
	s := state.New(rt.root)
	p, err := s.Load()
	if err != nil {
		return nil, 0, err
	}
	switch args[0] {
	case "list":
		sprintID := p.ActiveSprintID
		if len(args) > 1 {
			sprintID = args[1]
		}
		var out []*domain.Evidence
		for _, e := range p.Evidence {
			if sprintID == "" || e.SprintID == sprintID {
				out = append(out, e)
			}
		}
		return out, p.Revision, nil
	case "register":
		sp, err := activeSprint(p)
		if err != nil {
			return nil, 0, err
		}
		fs := flag.NewFlagSet("evidence register", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		path := fs.String("path", "", "")
		kind := fs.String("kind", "test-output", "")
		acs := fs.String("ac", "", "")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, 0, err
		}
		if *path == "" {
			return nil, 0, usageErr("--path required")
		}
		abs := *path
		if !filepath.IsAbs(abs) {
			abs = filepath.Join(rt.root, *path)
		}
		b, err := os.ReadFile(abs)
		if err != nil {
			return nil, 0, err
		}
		sum := sha256.Sum256(b)
		id := domain.NewID("evidence")
		ev := &domain.Evidence{EvidenceID: id, SprintID: sp.SprintID, Kind: *kind, URI: relative(rt.root, abs), SHA256: hex.EncodeToString(sum[:]), Size: int64(len(b)), CriterionIDs: csv(*acs), CreatedAt: time.Now().UTC(), RedactionStatus: "passed"}
		if looksSecret(b) {
			ev.RedactionStatus = "failed"
			return nil, 0, &commandError{"SEC_EVIDENCE_LEAK", "potential secret detected in evidence", "Remove or rotate the secret and create sanitized evidence.", 6}
		}
		r, err := s.Mutate(rt.expected, actor(), "EvidenceRegistered", id, ev, func(p *domain.Project) error { p.Evidence[id] = ev; return nil })
		if err != nil {
			return nil, 0, err
		}
		return r.Evidence[id], r.Revision, nil
	case "verify":
		bad := invalidEvidence(rt.root, p)
		if len(bad) > 0 {
			return map[string]any{"valid": false, "invalid": bad}, p.Revision, &commandError{"QA_EVIDENCE_INVALID", "evidence verification failed", "Regenerate the invalid artifacts.", 3}
		}
		return map[string]any{"valid": true, "count": len(p.Evidence)}, p.Revision, nil
	default:
		return nil, 0, usageErr("unknown evidence action")
	}
}

func (rt *runtime) qa(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("qa action required")
	}
	s := state.New(rt.root)
	p, err := s.Load()
	if err != nil {
		return nil, 0, err
	}
	sp, err := activeSprint(p)
	if err != nil {
		return nil, 0, err
	}
	switch args[0] {
	case "capabilities":
		return qaadapter.Discover(rt.root), p.Revision, nil
	case "plan":
		run := &domain.QARun{RunID: domain.NewID("run"), SprintID: sp.SprintID, Status: "planned", Environment: qaadapter.EnvironmentFingerprint(rt.root), StartedAt: time.Now().UTC(), StateRevision: p.Revision, SpecHash: workflow.QASpecHash(p, sp)}
		for _, ac := range p.Criteria {
			if slices.Contains(sp.IntentIDs, ac.IntentID) && ac.Priority == "blocking" {
				for _, variant := range qaVariants() {
					dispositions := map[string]string{"L1": "required", "L2": "required", "L3": "required", "L4": "required", "L5": "required", "L6": "required", "L7": "required"}
					run.Cases = append(run.Cases, domain.QACase{CaseID: domain.NewID("case"), CriterionIDs: []string{ac.CriterionID}, Title: ac.Statement + " — " + variant, Variant: variant, Layers: []string{"L1", "L2", "L3", "L4", "L5", "L6", "L7"}, Status: "pending", Actor: "project user", Preconditions: append([]string(nil), ac.Preconditions...), Steps: []domain.QAStep{{Action: qaVariantAction(variant), ExpectedUI: ac.Observable, ExpectedData: strings.Join(ac.Expected, "; "), ObserverIDs: []string{"interface", "boundary", "persistence"}}}, ForbiddenOutcomes: append([]string(nil), ac.Forbidden...), RequiredLayers: dispositions, Risk: "high"})
				}
			}
		}
		r, err := s.Mutate(rt.expected, actor(), "QAPlanned", run.RunID, run, func(p *domain.Project) error {
			p.QARuns[run.RunID] = run
			p.Sprints[sp.SprintID].LastQAID = run.RunID
			p.Sprints[sp.SprintID].LastQAStatus = "planned"
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		return r.QARuns[run.RunID], r.Revision, nil
	case "observe":
		if len(args) < 2 {
			return nil, 0, usageErr("use qa observe <case-id> --input FILE")
		}
		caseID := args[1]
		fs := flag.NewFlagSet("qa observe", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		input := fs.String("input", "", "")
		if err := fs.Parse(args[2:]); err != nil {
			return nil, 0, err
		}
		if *input == "" {
			return nil, 0, usageErr("--input required")
		}
		run := p.QARuns[sp.LastQAID]
		if run == nil {
			return nil, 0, notFound("qa run", sp.LastQAID)
		}
		var target *domain.QACase
		for i := range run.Cases {
			if run.Cases[i].CaseID == caseID {
				target = &run.Cases[i]
				break
			}
		}
		if target == nil {
			return nil, 0, notFound("qa case", caseID)
		}
		obs, b, err := qaadapter.ReadObservation(rt.root, *input)
		if err != nil {
			return nil, 0, err
		}
		if obs.RunID != run.RunID || obs.CaseID != caseID {
			return nil, 0, &commandError{"QA_OBSERVATION_MISMATCH", "observation run_id or case_id does not match active QA case", "Regenerate the observation for the active run and case.", 3}
		}
		if obs.SpecHash != run.SpecHash || obs.StateRevision != run.StateRevision {
			return nil, 0, &commandError{"QA_OBSERVATION_STALE", "observation specification hash or state revision is stale", "Regenerate the observation from the active QA plan.", 3}
		}
		if looksSecret(b) {
			return nil, 0, &commandError{"SEC_EVIDENCE_LEAK", "potential secret detected in observation", "Sanitize the observation and rotate exposed credentials.", 6}
		}
		abs := *input
		if !filepath.IsAbs(abs) {
			abs = filepath.Join(rt.root, *input)
		}
		sum := sha256.Sum256(b)
		id := domain.NewID("evidence")
		assertions := make([]domain.EvidenceAssertion, 0, len(obs.Assertions))
		for _, assertion := range obs.Assertions {
			assertions = append(assertions, domain.EvidenceAssertion{Statement: assertion.Statement, Passed: assertion.Passed, Layer: assertion.Layer, RequirementRefs: append([]string(nil), assertion.RequirementRefs...), Actual: assertion.Actual, Expected: assertion.Expected})
		}
		ev := &domain.Evidence{EvidenceID: id, SprintID: sp.SprintID, RunID: run.RunID, CaseID: caseID, SpecHash: run.SpecHash, StateRevision: run.StateRevision, Kind: "journey-observation", URI: relative(rt.root, abs), SHA256: hex.EncodeToString(sum[:]), Size: int64(len(b)), CriterionIDs: append([]string(nil), target.CriterionIDs...), CreatedAt: time.Now().UTC(), RedactionStatus: "passed", Layers: append([]string(nil), obs.Layers...), Assertions: assertions, Tool: obs.Adapter, ToolVersion: obs.ToolVersion, Environment: obs.Environment, StartedAt: &obs.StartedAt, FinishedAt: &obs.FinishedAt}
		r, err := s.Mutate(rt.expected, actor(), "QAObservationImported", id, obs, func(p *domain.Project) error {
			p.Evidence[id] = ev
			rr := p.QARuns[sp.LastQAID]
			for i := range rr.Cases {
				if rr.Cases[i].CaseID == caseID {
					rr.Cases[i].Status = "evidenced"
					rr.Cases[i].EvidenceIDs = unique(append(rr.Cases[i].EvidenceIDs, id))
				}
			}
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		return map[string]any{"observation": obs, "evidence": r.Evidence[id], "case_id": caseID}, r.Revision, nil
	case "execute":
		if len(args) < 2 {
			return nil, 0, usageErr("use qa execute <case-id> --adapter NAME")
		}
		caseID := args[1]
		fs := flag.NewFlagSet("qa execute", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		adapter := fs.String("adapter", "", "")
		if err := fs.Parse(args[2:]); err != nil {
			return nil, 0, err
		}
		if *adapter == "" {
			return nil, 0, usageErr("--adapter required")
		}
		run := p.QARuns[sp.LastQAID]
		if run == nil {
			return nil, 0, notFound("qa run", sp.LastQAID)
		}
		var target *domain.QACase
		for i := range run.Cases {
			if run.Cases[i].CaseID == caseID {
				target = &run.Cases[i]
				break
			}
		}
		if target == nil {
			return nil, 0, notFound("qa case", caseID)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		execution, executeErr := qaadapter.Execute(ctx, rt.root, *adapter)
		artifact, marshalErr := json.MarshalIndent(execution, "", "  ")
		if marshalErr != nil {
			return nil, 0, marshalErr
		}
		if looksSecret(artifact) {
			return nil, 0, &commandError{"SEC_EVIDENCE_LEAK", "potential secret detected in adapter output", "Sanitize test output and rotate exposed credentials.", 6}
		}
		dir := filepath.Join(rt.root, filepath.FromSlash(sp.DocumentRoot), "04-qa", "evidence")
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, 0, err
		}
		path := filepath.Join(dir, caseID+"-"+slugify(*adapter)+".json")
		if err := os.WriteFile(path, artifact, 0644); err != nil {
			return nil, 0, err
		}
		sum := sha256.Sum256(artifact)
		id := domain.NewID("evidence")
		capability, _ := qaadapter.CapabilityByName(rt.root, *adapter)
		assertions := []domain.EvidenceAssertion{}
		for _, layer := range capability.Layers {
			assertions = append(assertions, domain.EvidenceAssertion{Statement: "adapter completed successfully", Passed: executeErr == nil, Layer: layer, RequirementRefs: []string{"adapter-exit"}, Actual: fmt.Sprint(execution.Passed), Expected: "true"})
		}
		started, finished := execution.StartedAt, execution.FinishedAt
		ev := &domain.Evidence{EvidenceID: id, SprintID: sp.SprintID, RunID: run.RunID, CaseID: caseID, SpecHash: run.SpecHash, StateRevision: run.StateRevision, Kind: "adapter-execution", URI: relative(rt.root, path), SHA256: hex.EncodeToString(sum[:]), Size: int64(len(artifact)), CriterionIDs: append([]string(nil), target.CriterionIDs...), CreatedAt: time.Now().UTC(), RedactionStatus: "passed", Layers: append([]string(nil), capability.Layers...), Assertions: assertions, Tool: *adapter, ToolVersion: qaadapter.ToolVersion(capability.Command), Environment: run.Environment, StartedAt: &started, FinishedAt: &finished}
		statusValue := "passed"
		if executeErr != nil {
			statusValue = "failed"
		}
		r, err := s.Mutate(rt.expected, actor(), "QAAdapterExecuted", id, execution, func(p *domain.Project) error {
			p.Evidence[id] = ev
			rr := p.QARuns[sp.LastQAID]
			for i := range rr.Cases {
				if rr.Cases[i].CaseID == caseID {
					rr.Cases[i].Status = statusValue
					rr.Cases[i].EvidenceIDs = unique(append(rr.Cases[i].EvidenceIDs, id))
				}
			}
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		result := map[string]any{"execution": execution, "evidence": r.Evidence[id], "case_id": caseID}
		if executeErr != nil {
			return result, r.Revision, executeErr
		}
		return result, r.Revision, nil
	case "case":
		if len(args) < 3 {
			return nil, 0, usageErr("use qa case <case-id> <passed|failed> --evidence id,id")
		}
		caseID, statusValue := args[1], args[2]
		if statusValue == "passed" {
			return nil, p.Revision, &commandError{"QA_MANUAL_PASS_FORBIDDEN", "a QA case cannot be manually marked passed", "Attach structured evidence and run qa evaluate.", 3}
		}
		if statusValue != "failed" {
			return nil, 0, usageErr("manual case status must be failed")
		}
		evs := ""
		for i := 3; i < len(args)-1; i++ {
			if args[i] == "--evidence" {
				evs = args[i+1]
			}
		}
		run := p.QARuns[sp.LastQAID]
		if run == nil {
			return nil, 0, notFound("qa run", sp.LastQAID)
		}
		found := false
		for i := range run.Cases {
			if run.Cases[i].CaseID == caseID {
				found = true
			}
		}
		if !found {
			return nil, 0, notFound("qa case", caseID)
		}
		r, err := s.Mutate(rt.expected, actor(), "QACaseRecorded", caseID, map[string]any{"status": statusValue, "evidence": csv(evs)}, func(p *domain.Project) error {
			run := p.QARuns[sp.LastQAID]
			for i := range run.Cases {
				if run.Cases[i].CaseID == caseID {
					run.Cases[i].Status = statusValue
					run.Cases[i].EvidenceIDs = csv(evs)
				}
			}
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		return r.QARuns[sp.LastQAID], r.Revision, nil
	case "disposition":
		if len(args) < 3 {
			return nil, 0, usageErr("use qa disposition <case-id> <L1..L7> --status required|not-applicable|waived [--reason TEXT --approver ID]")
		}
		caseID, layer := args[1], strings.ToUpper(args[2])
		if !slices.Contains([]string{"L1", "L2", "L3", "L4", "L5", "L6", "L7"}, layer) {
			return nil, 0, usageErr("layer must be L1..L7")
		}
		fs := flag.NewFlagSet("qa disposition", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		statusValue := fs.String("status", "required", "")
		reason := fs.String("reason", "", "")
		approver := fs.String("approver", "", "")
		if err := fs.Parse(args[3:]); err != nil {
			return nil, 0, err
		}
		if !slices.Contains([]string{"required", "not-applicable", "waived"}, *statusValue) {
			return nil, 0, usageErr("status must be required, not-applicable, or waived")
		}
		if *statusValue != "required" && (*reason == "" || *approver == "") {
			return nil, 0, usageErr("--reason and --approver are required for non-required layers")
		}
		run := p.QARuns[sp.LastQAID]
		if run == nil {
			return nil, 0, notFound("qa run", sp.LastQAID)
		}
		found := false
		for i := range run.Cases {
			if run.Cases[i].CaseID == caseID {
				found = true
			}
		}
		if !found {
			return nil, 0, notFound("qa case", caseID)
		}
		disposition := *statusValue
		if *statusValue != "required" {
			disposition = *statusValue + ":" + *approver + ":" + *reason
		}
		r, err := s.Mutate(rt.expected, actor(), "QALayerDispositionChanged", caseID, map[string]string{"layer": layer, "disposition": disposition}, func(p *domain.Project) error {
			for i := range p.QARuns[sp.LastQAID].Cases {
				c := &p.QARuns[sp.LastQAID].Cases[i]
				if c.CaseID == caseID {
					if c.RequiredLayers == nil {
						c.RequiredLayers = map[string]string{}
					}
					c.RequiredLayers[layer] = disposition
				}
			}
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		return r.QARuns[sp.LastQAID], r.Revision, nil
	case "evaluate":
		run := p.QARuns[sp.LastQAID]
		findings := workflow.EvaluateQAGateAtRoot(rt.root, p, sp, run)
		statusValue := "passed"
		if workflow.Blocking(findings) {
			statusValue = "failed"
		}
		now := time.Now().UTC()
		r, err := s.Mutate(rt.expected, actor(), "QAEvaluated", sp.LastQAID, map[string]any{"status": statusValue, "findings": findings}, func(p *domain.Project) error {
			run := p.QARuns[sp.LastQAID]
			for i := range run.Cases {
				run.Cases[i].Status = "failed"
			}
			if statusValue == "passed" {
				for i := range run.Cases {
					run.Cases[i].Status = "passed"
				}
			}
			run.Status = statusValue
			run.FinishedAt = &now
			p.Sprints[sp.SprintID].LastQAStatus = statusValue
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
		if statusValue != "passed" {
			return map[string]any{"status": statusValue, "findings": findings, "run": r.QARuns[sp.LastQAID]}, r.Revision, &commandError{"QA_GATE_FAILED", "blocking QA criteria are not fully evidenced", "Resolve findings and rerun evaluation.", 3}
		}
		return map[string]any{"status": statusValue, "findings": findings, "run": r.QARuns[sp.LastQAID]}, r.Revision, nil
	case "status":
		return map[string]any{"status": sp.LastQAStatus, "run": p.QARuns[sp.LastQAID]}, p.Revision, nil
	default:
		return nil, 0, usageErr("unknown qa action")
	}
}

func qaVariants() []string {
	return []string{"happy", "alternate", "empty", "validation", "permission", "failure", "recovery"}
}
func qaVariantAction(v string) string {
	return map[string]string{"happy": "Complete the primary user and data journey", "alternate": "Complete a valid alternate or back-navigation journey", "empty": "Exercise empty initial and no-result states", "validation": "Submit invalid or boundary input and confirm no forbidden write", "permission": "Exercise least-privilege denial and authorized recovery", "failure": "Induce a downstream failure and observe user-visible and data-consistency behavior", "recovery": "Retry or roll back after failure and verify idempotent recovery"}[v]
}

func (rt *runtime) report(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("report action required")
	}
	s := state.New(rt.root)
	p, err := s.Load()
	if err != nil {
		return nil, 0, err
	}
	sp, err := activeSprint(p)
	if err != nil {
		return nil, 0, err
	}
	path := document.Path(rt.root, sp, domain.PhaseReport)
	switch args[0] {
	case "generate":
		if _, _, err := document.Scaffold(rt.root, p, sp, domain.PhaseReport); err != nil {
			return nil, 0, err
		}
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, 0, err
		}
		fmt.Fprintf(f, "\n<!-- tene:generated:summary:start -->\n### Generated Sprint Summary\n\n- Sprint: `%s`\n- Previous sprints: `%s`\n- Intent IDs: `%s`\n- Tasks: %d\n- QA verdict: `%s`\n- Open gaps: %d\n- State revision: %d\n\n<!-- tene:generated:summary:end -->\n", sp.SprintID, strings.Join(sp.Predecessors, ", "), strings.Join(sp.IntentIDs, ", "), len(sp.TaskIDs), sp.LastQAStatus, len(sp.OpenGapIDs), p.Revision)
		_ = f.Close()
		r, err := s.Mutate(rt.expected, actor(), "ReportGenerated", sp.SprintID, map[string]string{"path": relative(rt.root, path)}, func(p *domain.Project) error { p.Sprints[sp.SprintID].ReportPath = relative(rt.root, path); return nil })
		if err != nil {
			return nil, 0, err
		}
		return map[string]string{"path": path}, r.Revision, nil
	case "validate":
		findings := document.Validate(path, domain.PhaseReport)
		if workflow.Blocking(findings) {
			return map[string]any{"valid": false, "findings": findings}, p.Revision, &commandError{"DOC_VALIDATION_FAILED", "report validation failed", "Complete the report sections.", 3}
		}
		return map[string]any{"valid": true, "findings": findings}, p.Revision, nil
	default:
		return nil, 0, usageErr("unknown report action")
	}
}

func (rt *runtime) secrets(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("secret action required")
	}
	switch args[0] {
	case "check":
		path, err := secret.Check()
		return map[string]string{"path": path}, 0, err
	case "list":
		if len(args) < 2 {
			return nil, 0, usageErr("environment required")
		}
		v, err := secret.ListNames(context.Background(), args[1])
		return v, 0, err
	case "run":
		if len(args) < 4 || args[2] != "--" {
			return nil, 0, usageErr("use secret run <environment> -- <command> [args]")
		}
		v, err := secret.Run(context.Background(), args[1], args[3:])
		if err != nil {
			return v, 0, &commandError{"SEC_CHILD_FAILED", err.Error(), "Inspect the sanitized child error without exposing secrets.", 8}
		}
		return v, 0, nil
	default:
		return nil, 0, usageErr("unknown secret action")
	}
}

func (rt *runtime) migrate(args []string) (any, uint64, error) {
	if len(args) == 0 {
		return nil, 0, usageErr("migrate action required")
	}
	s := state.New(rt.root)
	plan, err := s.PlanMigration()
	if err != nil {
		return nil, 0, err
	}
	switch args[0] {
	case "status", "dry-run":
		return plan, 0, nil
	case "apply":
		updated, err := s.Migrate()
		if err != nil {
			return updated, 0, err
		}
		p, loadErr := s.Load()
		if loadErr != nil {
			return updated, 0, loadErr
		}
		return updated, p.Revision, nil
	default:
		return nil, 0, usageErr("unknown migrate action")
	}
}

func (rt *runtime) doctor(args []string) (any, uint64, error) {
	s := state.New(rt.root)
	var repaired []string
	if len(args) > 0 {
		if len(args) != 1 || args[0] != "--repair" {
			return nil, 0, usageErr("use doctor [--repair]")
		}
		var repairErr error
		repaired, repairErr = s.RepairDerived()
		if repairErr != nil {
			return nil, 0, repairErr
		}
	}
	p, err := s.Load()
	if err != nil {
		return nil, 0, err
	}
	events, err := s.VerifyJournal()
	if err != nil {
		return nil, p.Revision, err
	}
	findings := validateGraph(p)
	drift, _, replayErr := s.ProjectionDrift()
	if replayErr != nil {
		findings = append(findings, domain.Finding{Code: "STATE_REPLAY_FAILED", Severity: "blocker", Message: replayErr.Error(), Remediation: "Run compact on a known-good projection to create a checkpoint."})
	} else {
		for _, d := range drift {
			if d.Status != "match" {
				findings = append(findings, domain.Finding{Code: "STATE_PROJECTION_DRIFT", Severity: "blocker", SubjectRefs: []string{relative(rt.root, d.Path)}, Message: "projection does not match journal replay: " + d.Status, Remediation: "Run doctor --repair after reviewing the journal."})
			}
		}
	}
	for _, id := range invalidEvidence(rt.root, p) {
		findings = append(findings, domain.Finding{Code: "QA_EVIDENCE_INVALID", Severity: "blocker", SubjectRefs: []string{id}, Message: "evidence is missing, modified, or contains a secret pattern", Remediation: "Restore or regenerate sanitized evidence."})
	}
	for _, sp := range p.Sprints {
		for _, ph := range []domain.Phase{domain.PhasePRD, domain.PhasePlan, domain.PhaseDesign, domain.PhaseLoopCheck, domain.PhaseQA, domain.PhaseReport} {
			if documentDue(sp.Phase, ph) && document.Exists(rt.root, sp, ph) {
				findings = append(findings, document.Validate(document.Path(rt.root, sp, ph), ph)...)
			}
		}
	}
	return map[string]any{"healthy": !workflow.Blocking(findings), "events": len(events), "revision": p.Revision, "findings": findings, "projection_drift": drift, "repaired": repaired, "capabilities": map[string]any{"tene_cli": func() bool { _, e := secret.Check(); return e == nil }(), "codex": projectconfig.ProbeCodex(rt.root)}}, p.Revision, nil
}

func invalidEvidence(root string, p *domain.Project) []string {
	var bad []string
	for id, e := range p.Evidence {
		path := e.URI
		if !filepath.IsAbs(path) {
			path = filepath.Join(root, path)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			bad = append(bad, id)
			continue
		}
		sum := sha256.Sum256(b)
		if hex.EncodeToString(sum[:]) != e.SHA256 || looksSecret(b) {
			bad = append(bad, id)
		}
	}
	sort.Strings(bad)
	return bad
}

func documentDue(current, documentPhase domain.Phase) bool {
	rank := map[domain.Phase]int{domain.PhaseDraft: 0, domain.PhasePRD: 1, domain.PhasePlan: 2, domain.PhaseDesign: 3, domain.PhaseDo: 4, domain.PhaseLoopCheck: 5, domain.PhaseQA: 6, domain.PhaseReport: 7, domain.PhaseArchived: 8}
	documentRank := map[domain.Phase]int{domain.PhasePRD: 1, domain.PhasePlan: 2, domain.PhaseDesign: 3, domain.PhaseLoopCheck: 5, domain.PhaseQA: 6, domain.PhaseReport: 7}
	return documentRank[documentPhase] <= rank[current]
}
func (rt *runtime) compact() (any, uint64, error) {
	s := state.New(rt.root)
	path, p, err := s.CreateCheckpoint()
	if err != nil {
		return nil, 0, err
	}
	return map[string]string{"snapshot": relative(rt.root, path), "checkpoint": "created"}, p.Revision, nil
}
func (rt *runtime) clear() (any, uint64, error) {
	s := state.New(rt.root)
	p, err := s.Load()
	if err != nil {
		return nil, 0, err
	}
	err = s.ClearDerived()
	return map[string]bool{"derived_state_cleared": err == nil}, p.Revision, err
}

func buildGraph(p *domain.Project) domain.Graph {
	g := domain.Graph{Nodes: map[string]domain.Node{}, Edges: map[string]domain.Edge{}}
	for id, in := range p.Intents {
		g.Nodes[id] = domain.Node{ID: id, Kind: "Intent", Label: in.Statement, Source: "authored", Confidence: 1}
	}
	for id, ac := range p.Criteria {
		g.Nodes[id] = domain.Node{ID: id, Kind: "AcceptanceCriterion", Label: ac.Statement, Source: "authored", Confidence: 1}
		addEdge(&g, ac.IntentID, id, "realizes")
	}
	for id, t := range p.Tasks {
		g.Nodes[id] = domain.Node{ID: id, Kind: "Task", Label: t.Title, Source: "authored", Confidence: 1}
		for _, ac := range t.CriterionIDs {
			addEdge(&g, ac, id, "realizes")
		}
	}
	for id, e := range p.Evidence {
		g.Nodes[id] = domain.Node{ID: id, Kind: "Evidence", Label: e.Kind, Locator: e.URI, Source: "observed", Confidence: 1}
		for _, ac := range e.CriterionIDs {
			addEdge(&g, id, ac, "verifies")
		}
	}
	for id, gap := range p.Gaps {
		g.Nodes[id] = domain.Node{ID: id, Kind: "Gap", Label: gap.Description, Source: "authored", Confidence: 1, Attributes: map[string]any{"category": gap.Category, "severity": gap.Severity, "status": gap.Status}}
		for _, subject := range gap.SubjectRefs {
			addEdge(&g, subject, id, "depends_on")
		}
	}
	for id, sprint := range p.Sprints {
		g.Nodes[id] = domain.Node{ID: id, Kind: "Sprint", Label: sprint.Title, Locator: sprint.DocumentRoot, Source: "authored", Confidence: 1, Attributes: map[string]any{"phase": sprint.Phase}}
	}
	for id, approval := range p.Approvals {
		g.Nodes[id] = domain.Node{ID: id, Kind: "Approval", Label: string(approval.From) + " → " + string(approval.To), Source: "authored", Confidence: 1, Attributes: map[string]any{"status": approval.Status, "requester": approval.Requester, "approver": approval.Approver, "expires_at": approval.ExpiresAt}}
		addEdge(&g, approval.SprintID, id, "belongs_to")
	}
	for id, w := range p.Waivers {
		g.Nodes[id] = domain.Node{ID: id, Kind: "Waiver", Label: w.Reason, Source: "authored", Confidence: 1, Attributes: map[string]any{"status": w.Status, "expires_at": w.ExpiresAt, "approver": w.Approver, "scope": w.Scope}}
		addEdge(&g, w.GapID, id, "waived_by")
	}
	return g
}
func addEdge(g *domain.Graph, from, to, kind string) {
	if from == "" || to == "" {
		return
	}
	id := deterministicEdgeID(from, to, kind, "workflow-state", "workflow-state")
	g.Edges[id] = domain.Edge{ID: id, From: from, To: to, Kind: kind, SourceLocator: "workflow-state", Provider: "workflow-state", Confidence: 1}
}
func mergeCodeGraph(g *domain.Graph, report codeintel.Report) {
	for _, f := range report.Files {
		g.Nodes[f.ID] = domain.Node{ID: f.ID, Kind: "File", Label: f.Path, Locator: f.Path, Source: "derived", Confidence: f.Confidence, Attributes: map[string]any{"primary_layer": f.Layer, "layer_reason": f.LayerReason}}
	}
	for _, c := range report.Components {
		g.Nodes[c.ID] = domain.Node{ID: c.ID, Kind: "Symbol", Label: c.Name, Locator: c.Locator, Source: "derived", Confidence: c.Confidence, Attributes: map[string]any{"kind": c.Kind, "primary_layer": c.PrimaryLayer, "layer_reason": c.LayerReason, "imports": c.Imports, "references": c.References, "calls": c.Calls, "inputs": c.Inputs, "outputs": c.Outputs, "effects": c.Effects, "unknown": c.Unknown, "provider": c.Provider}}
	}
	for _, e := range report.Edges {
		if _, ok := g.Nodes[e.From]; !ok {
			g.Nodes[e.From] = domain.Node{ID: e.From, Kind: "Unknown", Label: e.From, Source: "derived", Confidence: e.Confidence}
		}
		if _, ok := g.Nodes[e.To]; !ok {
			g.Nodes[e.To] = domain.Node{ID: e.To, Kind: "ExternalSymbol", Label: strings.TrimPrefix(e.To, "call:"), Source: "derived", Confidence: e.Confidence, Attributes: map[string]any{"resolution": "unresolved by current provider"}}
		}
		id := deterministicEdgeID(e.From, e.To, e.Kind, e.Provider, e.Locator)
		g.Edges[id] = domain.Edge{ID: id, From: e.From, To: e.To, Kind: e.Kind, SourceLocator: e.Locator, Provider: e.Provider, Confidence: e.Confidence}
	}
}
func deterministicEdgeID(from, to, kind, provider, locator string) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{from, to, kind, provider, locator}, "\x00")))
	return "edge_" + hex.EncodeToString(sum[:16])
}
func traceGraph(p *domain.Project, id string) map[string]any {
	nodes := map[string]domain.Node{}
	edges := []domain.Edge{}
	if n, ok := p.Graph.Nodes[id]; ok {
		nodes[id] = n
	}
	changed := true
	for changed {
		changed = false
		for _, e := range p.Graph.Edges {
			_, from := nodes[e.From]
			_, to := nodes[e.To]
			if from || to {
				edges = appendUniqueEdge(edges, e)
				if n, ok := p.Graph.Nodes[e.From]; ok {
					if _, x := nodes[e.From]; !x {
						nodes[e.From] = n
						changed = true
					}
				}
				if n, ok := p.Graph.Nodes[e.To]; ok {
					if _, x := nodes[e.To]; !x {
						nodes[e.To] = n
						changed = true
					}
				}
			}
		}
	}
	return map[string]any{"nodes": nodes, "edges": edges}
}
func validateGraph(p *domain.Project) []domain.Finding {
	return tracecontext.ValidateGraph(p)
}

func safeRootPath(root, requested string) (string, error) {
	if filepath.IsAbs(requested) {
		return "", &commandError{"CONTEXT_PATH_OUTSIDE_ROOT", "absolute paths are not allowed", "Use a repository-relative path.", 2}
	}
	clean := filepath.Clean(requested)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", &commandError{"CONTEXT_PATH_OUTSIDE_ROOT", "path escapes repository root", "Use a repository-relative path.", 2}
	}
	return filepath.Join(root, clean), nil
}
func writeDerivedJSON(root, requested string, value any) error {
	path, err := safeRootPath(root, requested)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}
func appendUniqueEdge(v []domain.Edge, e domain.Edge) []domain.Edge {
	for _, x := range v {
		if x.ID == e.ID {
			return v
		}
	}
	return append(v, e)
}
func activeSprint(p *domain.Project) (*domain.Sprint, error) {
	if p.ActiveSprintID == "" {
		return nil, &commandError{"WF_NO_ACTIVE_SPRINT", "there is no active sprint", "Create or start a sprint.", 3}
	}
	sp := p.Sprints[p.ActiveSprintID]
	if sp == nil {
		return nil, &commandError{"STATE_CORRUPT", "active sprint does not exist", "Run doctor.", 7}
	}
	return sp, nil
}
func selectSprint(p *domain.Project, id string) (*domain.Sprint, error) {
	if id == "" {
		return activeSprint(p)
	}
	if sp := p.Sprints[id]; sp != nil {
		return sp, nil
	}
	for _, sp := range p.Sprints {
		if sp.Slug == id {
			return sp, nil
		}
	}
	return nil, notFound("sprint", id)
}
func sortedSprints(p *domain.Project) []*domain.Sprint {
	out := make([]*domain.Sprint, 0, len(p.Sprints))
	for _, x := range p.Sprints {
		out = append(out, x)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].SprintID < out[j].SprintID })
	return out
}
func openGaps(p *domain.Project) int {
	n := 0
	for _, g := range p.Gaps {
		if g.Status == "open" {
			n++
		}
	}
	return n
}
func deferredGaps(p *domain.Project) int {
	n := 0
	for _, g := range p.Gaps {
		if g.Status == "deferred" {
			n++
		}
	}
	return n
}
func effectiveBlockers(p *domain.Project, sp *domain.Sprint) int {
	n := 0
	now := time.Now().UTC()
	for _, id := range sp.OpenGapIDs {
		if g := p.Gaps[id]; g != nil && g.Status == "open" && g.Severity == "blocker" && !workflow.ActiveWaiver(p, g, now) {
			n++
		}
	}
	return n
}
func activeWaiverCount(p *domain.Project, sp *domain.Sprint) int {
	n := 0
	now := time.Now().UTC()
	for _, id := range sp.OpenGapIDs {
		if workflow.ActiveWaiver(p, p.Gaps[id], now) {
			n++
		}
	}
	return n
}
func actor() domain.Actor { return domain.Actor{Kind: "codex"} }
func arg(a []string, i int) string {
	if i < len(a) {
		return a[i]
	}
	return ""
}
func csv(v string) []string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if x := strings.TrimSpace(p); x != "" {
			out = append(out, x)
		}
	}
	return out
}
func remove(v []string, s string) []string {
	out := v[:0]
	for _, x := range v {
		if x != s {
			out = append(out, x)
		}
	}
	return out
}
func unique(values []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" && !seen[value] {
			seen[value] = true
			out = append(out, value)
		}
	}
	return out
}
func slugify(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	var b strings.Builder
	dash := false
	for _, r := range v {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
			dash = false
		} else if !dash && b.Len() > 0 {
			b.WriteByte('-')
			dash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
func relative(root, path string) string {
	r, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(r)
}
func looksSecret(b []byte) bool {
	s := string(b)
	patterns := []string{"-----BEGIN PRIVATE KEY-----", "-----BEGIN RSA PRIVATE KEY-----", "AKIA", "ghp_", "github_pat_"}
	for _, p := range patterns {
		if strings.Contains(s, p) {
			return true
		}
	}
	return secret.DetectLeak(b)
}
func discoverRoot() string {
	cwd, _ := os.Getwd()
	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, state.DirName)); err == nil {
			return dir
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		next := filepath.Dir(dir)
		if next == dir {
			return cwd
		}
		dir = next
	}
}
func usageErr(msg string) error { return &commandError{"CLI_USAGE", msg, "Run tene-workflow help.", 2} }
func notFound(kind, id string) error {
	return &commandError{"STATE_NOT_FOUND", kind + " not found: " + id, "List available objects and retry.", 2}
}
func printHuman(w io.Writer, v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Fprintln(w, string(b))
}
func usage() string {
	return `tene-workflow — spec-driven sprint workflow

Usage:
  tene-workflow [--root PATH] [--json] <command>

Commands:
  init, status, route, master, sprint, phase, approval, task, intent, document, graph, context,
  loop, waiver, evidence, qa, report, secret, migrate, doctor, compact, clear, version`
}
