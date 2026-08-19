---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v060y01n452j1rk3xr4
phase: design
status: draft
revision: 1
intent_ids: []
generated_at: 2026-08-19T17:11:18Z
generated_by: tene-workflow
---

# design — MVP workflow engine and Codex plugin

<!-- tene:section:purpose -->
## Purpose

Define the implemented vertical slice and its dependency boundaries closely enough for independent loop-check and QA.

<!-- tene:section:scope -->
## Scope

The design instantiates `docs/02-design` with standard-library Go packages, Python hooks, Markdown skills, JSON Schemas, CI and release packaging.

<!-- tene:section:layers -->
## Layers

- Interface: `cmd/tene-workflow`, `internal/app`, nine `skills/tene-*`, `hooks/hooks.json` and `tene_hook.py`.
- Business Logic: `internal/workflow`, graph/context and QA use cases in `internal/app`.
- Persistence: `internal/state`, `internal/document`, `.tene-workflow`, `docs/sprints`, schemas.
- Infrastructure: `internal/secret`, Makefile, GitHub Actions, packaging wrapper and platform binaries.

<!-- tene:section:six-questions -->
## Six questions

`app.Run` is defined in `internal/app/app.go`, called by `cmd/tene-workflow/main.go`, receives argv/writers/version and returns an exit code while mutating state through `state.Store`. `Store.Mutate` is defined in `internal/state/store.go`, called by command use cases, receives expected revision/event/reducer and returns the new project projection. `CanTransition` and `EvaluateQAGate` are defined in `internal/workflow/workflow.go`, called by phase/QA commands, receive typed state and return findings without mutation. Hook handlers are defined in `hooks/tene_hook.py`, called by Codex hook events, receive JSON stdin and return supported hook JSON.

<!-- tene:section:traceability -->
## Traceability

Tasks are linked one-to-one to five blocking acceptance criteria; graph build materializes intent→AC→task and evidence→AC edges.

<!-- tene:section:decisions -->
## Decisions

Canonical mutation crosses the CLI only. Hooks are optional and advisory except narrow secret denials. Custom agents are installed without overwriting existing profiles. Builder and evaluator responsibilities remain isolated.

<!-- tene:section:freeform -->
## Freeform

Provider-specific browser and CodeGraph adapters are represented by skill capability routing in this pre-alpha; first-class Go adapter ports are follow-up work.

<!-- tene:section:components -->
## Components

Domain types/IDs; atomic hash-chained state store; workflow guard engine; document generator/validator; command application; secret runner; plugin skills/references; lifecycle hooks; project agent scaffolder; test and release harness.

<!-- tene:section:interfaces -->
## Interfaces

CLI commands return `{ok,schema_version,request_id,revision,result,warnings,errors}` under `--json`. Skills call the wrapper rather than state files. Hooks use Codex `SessionStart`, `PreToolUse`, `PreCompact`, `SubagentStart`, and `Stop` schemas.

<!-- tene:section:data -->
## Data

`Project` contains Sprint, Task, Intent, Criterion, Gap, Evidence, QARun and Graph maps. Events carry sequence, revision, previous hash and hash. Markdown sections use stable `tene:section:*` markers and require authored content or evidence-based N/A.

<!-- tene:section:state-transitions -->
## State transitions

Only draft→prd→plan→design→do→loop-check→qa→report→archived is forward-valid, with explicit repair edges from loop-check/qa. Archive relocates documents under `_archive/YYYY-MM` and clears the active pointer.

<!-- tene:section:failures -->
## Failures

Stable exit classes cover validation, guard, conflict, capability, security, corruption and child execution. Atomic temp-write/rename, repository lock, expected revision, hash verification and archive rollback prevent silent state loss.

<!-- tene:section:security -->
## Security

The model never receives a secret value. The runner invokes `tene run` with argv, rejects shell/environment dumps, and hooks block direct `tene get`, export and `.tene` access. Evidence is hashed and scanned before registration.

<!-- tene:section:tests -->
## Tests

Table and integration tests cover journal chain, stale revisions, document content, invalid transitions, full Sprint archive, all QA variants, agent scaffold preservation and hooks. Race, vet, validators and four-target package smoke tests are release gates.
