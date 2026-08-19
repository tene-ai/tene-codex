---
schema_version: 1.0.0
document_type: report
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v060y01n452j1rk3xr4
phase: report
status: draft
revision: 1
intent_ids: []
generated_at: 2026-08-19T17:11:18Z
generated_by: tene-workflow
---

# report — MVP workflow engine and Codex plugin

<!-- tene:section:purpose -->
## Purpose

Review the first executable tene-codex vertical slice, its evidence, limitations and continuity into the next product Sprint.

<!-- tene:section:scope -->
## Scope

Reports the Go workflow engine, Codex plugin/skills/hooks/subagents, QA/security gates, schemas, CI/release packaging, docs and dogfood state created from the documentation-only baseline.

<!-- tene:section:layers -->
## Layers

- **Interface:** `cmd/tene-workflow/main.go`, `internal/app/app.go`, `.codex-plugin/plugin.json`, nine `skills/tene-*`, `hooks/hooks.json`, `hooks/tene_hook.py`, and generated `.codex/agents/*.toml` expose CLI, natural-language, lifecycle, and delegated-worker entry points.
- **Business Logic:** `internal/workflow/workflow.go` enforces transitions, task traceability and all-case QA; app use cases build graph/context, manage intent/gaps/evidence and generate reports.
- **Persistence:** `internal/state/store.go` writes hash-chained events, atomic projections, master plan, policy and snapshots; `internal/document` creates and validates authored Sprint documents; schemas describe public state/evidence contracts.
- **Infrastructure:** `internal/secret/runner.go` delegates to tene; scripts locate/package binaries; Makefile and GitHub Actions run tests and produce four-platform bundles.

<!-- tene:section:six-questions -->
## Six questions

| Name | Defined at | Imported/referenced | Called/used | Input | Output/mutation |
|---|---|---|---|---|---|
| `Run` | `internal/app/app.go` | `cmd/tene-workflow/main.go`, app tests | process entry | argv, writers, version | exit code, JSON/human output, use-case mutations |
| `Store.Mutate` | `internal/state/store.go` | command use cases | every canonical mutation | expected revision, actor, event, reducer | journal event, project/active/master projections |
| `CanTransition` | `internal/workflow/workflow.go` | phase command/tests | transition/dry-run | project, sprint, target, doc readiness | findings without mutation |
| `EvaluateQAGate` | `internal/workflow/workflow.go` | QA evaluate/tests | final QA gate | project, sprint, QA run | AC findings without mutation |
| `Scaffold`/`Validate` | `internal/document/document.go` | sprint/document/report commands | phase document lifecycle | project, sprint, phase/path | Markdown or validation findings |
| `ScaffoldAgents` | `internal/projectconfig/agents.go` | init command/tests | project initialization | repository root | five non-overwriting custom agent files |
| `secret.Run` | `internal/secret/runner.go` | secret command | secret-dependent child execution | environment alias, argv | sanitized child result; no stored secret value |
| Hook handlers | `hooks/tene_hook.py` | `hooks/hooks.json` | Codex lifecycle | JSON stdin/event action | context or narrow deny decision |

<!-- tene:section:traceability -->
## Traceability

One confirmed intent connects to five blocking ACs, five completed tasks, 15 QA cases and one hash-verified evidence artifact. Rebuilt graph has 12 nodes and 15 edges with no invariant finding. The Sprint itself proved phase, task, document, QA and archive gates.

<!-- tene:section:decisions -->
## Decisions

Keep the core dependency-free and local; keep secret manager and workflow binaries separate; install subagent profiles on init without overwriting; discover hooks through the default plugin path; treat hooks as optional defense in depth; require all generated QA variants; ship pre-alpha before adding remote orchestration.

<!-- tene:section:freeform -->
## Freeform

Dogfooding was materially useful: it exposed missing task trace repair and overly weak document/QA gates, all of which were corrected before final QA.

<!-- tene:section:previous-sprints -->
## Previous sprints

This is the first implementation Sprint. It continues the RND, product PRD, master implementation plan and detailed design committed in baseline `77af284`; no prior executable Sprint exists.

<!-- tene:section:changed-files -->
## Changed files

- Go: `go.mod`, `cmd/tene-workflow`, and `internal/{app,document,domain,projectconfig,secret,state,workflow}` with tests.
- Plugin: `.codex-plugin/plugin.json`, nine skill folders, shared references, hooks and wrapper/package scripts.
- Contracts: three JSON Schemas, template contract, `.gitignore`, Makefile and SECURITY policy.
- Automation: `.github/workflows/ci.yml` and `release.yml`.
- Dogfood: `.tene-workflow`, `.codex/agents`, and this Sprint's PRD/plan/design/analysis/QA/report tree.
- README files are updated separately in the same working tree to describe runnable pre-alpha usage.

<!-- tene:section:intent-fulfillment -->
## Intent fulfillment

Deterministic phase guards, durable repository state/context, valid plugin integration, strict evidence QA and the tene secret boundary each have implementation and passing evidence matching their blocking criteria.

<!-- tene:section:qa-verdict -->
## Qa verdict

Passed at run `run_0000380v0gf3be5e8yd4tv3s7m`: 15/15 cases passed, evidence hash/redaction passed, graph invariant passed, race/vet/unit/integration/hook/validator/package gates passed. No live credential or production system was used.

<!-- tene:section:deferred-work -->
## Deferred work

- Browser/Playwright/Chrome adapter interfaces and runtime observers: no UI product exists in this Sprint; target Sprint 2.
- Semantic CodeGraph/LSP/data-flow providers and automatic Six Questions materialization: authored graph MVP is sufficient for this vertical slice; target Sprint 2.
- Waiver, migration apply/rollback and event-reducer reconstruction: contracts exist but full commands are not complete; target stabilization Sprint.
- Marketplace portal submission, screenshots, privacy/terms URLs and signed release: defer until public beta quality.
- Windows binary/hook support: packaging currently targets macOS/Linux; defer based on demand.

<!-- tene:section:next-sprint -->
## Next sprint

Build code-intelligence providers and real intent-driven QA adapters against a reference web application, including Playwright/browser UX journeys and API/data observers. Add formal schema validation/migration and independent evaluator orchestration through custom subagents.


<!-- tene:generated:summary:start -->
### Generated Sprint Summary

- Sprint: `sprint_0000380v060y01n452j1rk3xr4`
- Previous sprints: ``
- Intent IDs: `intent_0000380v0879x1qrpeygtbas6g`
- Tasks: 5
- QA verdict: `passed`
- Open gaps: 0
- State revision: 56

<!-- tene:generated:summary:end -->
