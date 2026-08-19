---
schema_version: 1.0.0
document_type: analysis
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v060y01n452j1rk3xr4
phase: loop-check
status: draft
revision: 1
intent_ids: []
generated_at: 2026-08-19T17:11:18Z
generated_by: tene-workflow
---

# loop-check — MVP workflow engine and Codex plugin

<!-- tene:section:purpose -->
## Purpose

Compare the implemented repository against the canonical product, plan and design documents before intent-driven QA.

<!-- tene:section:scope -->
## Scope

Inspect core command behavior, state invariants, document gates, plugin surfaces, subagent and hook behavior, secret boundary, tests, CI and packaging.

<!-- tene:section:layers -->
## Layers

Interface, Business Logic, Persistence and Infrastructure all have changed artifacts and direct verification. No layer is N/A.

<!-- tene:section:six-questions -->
## Six questions

Definitions and call paths for `Run`, `Store.Mutate`, `CanTransition`, `EvaluateQAGate`, `ScaffoldAgents`, `secret.Run`, and hook handlers were checked against their inputs and state/output mutations. Tests exercise these public seams rather than private wording.

<!-- tene:section:traceability -->
## Traceability

The graph contains one confirmed intent, five blocking criteria, five linked tasks and their realizes edges. Evidence edges are intentionally pending until QA.

<!-- tene:section:decisions -->
## Decisions

Treat empty required document sections and untraced tasks as blockers. Add a post-creation `task link` command because dogfood exposed that traceability needed a repair operation.

<!-- tene:section:freeform -->
## Freeform

The current graph provider covers authored workflow relationships; semantic source-code/data-flow provider depth is a documented post-MVP extension.

<!-- tene:section:baseline -->
## Baseline

Baseline is commit `77af284` containing research, PRD, plan and design only. The working tree adds the first executable pre-alpha implementation.

<!-- tene:section:changed-artifacts -->
## Changed artifacts

Added Go module and packages under `cmd/` and `internal/`; plugin manifest, skills, references, hooks and scripts; schemas/templates; CI/release workflows; security policy; tests; and dogfood Sprint state/documents.

<!-- tene:section:gap-matrix -->
## Gap matrix

- Resolved: tasks created without AC links after a shell extraction failure. The plan→design gate blocked the transition and `task link` was implemented and tested.
- Resolved: empty generated documents initially passed structural validation. Authored-content checks now block the relevant phase.
- Resolved: QA initially allowed one successful variant to satisfy an AC. The evaluator now requires every generated case and matching criterion-linked evidence.
- Open non-blocking: first-class browser, CodeGraph and external observer adapter ports remain future work.

<!-- tene:section:iterations -->
## Iterations

Iteration 1 built the vertical core. Iteration 2 fixed event hash normalization. Iteration 3 added plugin surfaces and full Sprint test. Iteration 4 strengthened document/task/QA gates and added subagent profiles. Iteration 5 dogfooded the workflow and added task trace repair.

<!-- tene:section:regression -->
## Regression

`go test -race ./...`, `go vet ./...`, hook tests, plugin/skill validators and four-platform cross-build passed before QA. Existing documentation and license files remain intact.
