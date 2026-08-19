---
schema_version: 1.0.0
document_type: analysis
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v2b2q0668t7tajkj3m0
phase: loop-check
status: draft
revision: 88
intent_ids: []
generated_at: 2026-08-19T17:30:10Z
generated_by: tene-workflow
---

# loop-check — Code Intelligence and Intent QA Adapters

<!-- tene:section:purpose -->
## Purpose

Compare the WP-08/WP-10 implementation against the confirmed intent, PRD and design before QA.

<!-- tene:section:scope -->
## Scope

Code intelligence providers/materialization, QA discovery/execution/observation, traceability validation, CLI integration and public contracts.

<!-- tene:section:layers -->
## Layers

Interface: new graph and QA actions are routed by app. Business Logic: codeintel/qaadapter keep provider logic outside CLI. Persistence: derived graph and hash-bound evidence reuse the state store. Infrastructure: git/Go/CodeGraph/npm/npx are probed through argv-only commands.

<!-- tene:section:six-questions -->
## Six questions

`Analyze` is defined in `internal/codeintel/codeintel.go`, imported/called by app graph commands, accepts context/root/path selection and returns Report/error. `Execute`, `Discover` and `ReadObservation` are defined in `internal/qaadapter/qaadapter.go`, called by app QA commands, accept bounded adapter/artifact inputs and return structured capabilities/results/observations. State mutations occur only in app after validation.

<!-- tene:section:traceability -->
## Traceability

AC1 → codeintel tests and graph understand output. AC2 → qaadapter tests/capabilities/observe/execute. AC3 → full verification evidence. All three implementation tasks now reference valid AC IDs; invalid references found during dogfood were removed with `task link --replace`.

<!-- tene:section:decisions -->
## Decisions

Keep CodeGraph explicit through `--codegraph-query`; never create indexes. Keep external browser execution outside the CLI but make its observation format first-class and validated.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:baseline -->
## Baseline

MVP graph held only authored intent/task/evidence nodes. QA cases held variant/layers/status but no actor, journey steps, observer checkpoints, risk, layer dispositions or adapter capability model.

<!-- tene:section:changed-artifacts -->
## Changed artifacts

Added `internal/codeintel`, `internal/qaadapter`, QA observation schema and tests. Extended domain QACase and app graph/task/QA commands. Updated QA skill, CLI/QA references and bilingual README.

<!-- tene:section:gap-matrix -->
## Gap matrix

- AC1: matched; deterministic Go fallback emits all Six Questions and explicit semantic-index unknown.
- AC2: matched; capability probe, allowlisted execution and validated external observation import are implemented.
- AC3: pending final QA evidence, no implementation gap.
- Newly found: task add accepted dangling AC/dependency references. Resolved with creation-time validation and replace-mode repair.

<!-- tene:section:iterations -->
## Iterations

Iteration 1 added provider/materializer and observation contracts. Iteration 2 added real allowlisted adapter execution, explicit CodeGraph query, path containment and task-reference validation.

<!-- tene:section:regression -->
## Regression

Package tests and vet pass. Existing archive/evidence tests remain green; final race/plugin/schema/skill checks run in QA.
