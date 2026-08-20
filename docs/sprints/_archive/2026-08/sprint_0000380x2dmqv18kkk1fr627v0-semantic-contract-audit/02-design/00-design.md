---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380x2dmqv18kkk1fr627v0
phase: design
status: complete
revision: 871
intent_ids: [intent_0000380x2dqcf9ctzpr01bb7jr]
generated_at: 2026-08-20T02:50:07Z
generated_by: tene-workflow
---

# design — Semantic Contract Completion Audit

<!-- tene:section:purpose -->
## Purpose

Specify the remaining canonical state and a fail-closed semantic acceptance engine.

<!-- tene:section:scope -->
## Scope

Master/intent data, graph/context projections, contract manifest schema, executable audit and final gate integration.

<!-- tene:section:layers -->
## Layers

CLI flags/commands→domain validation/projections→atomic event state/derived graph→test and release subprocesses.

<!-- tene:section:six-questions -->
## Six questions

New types/helpers are defined in domain/app/tracecontext; public commands and audit call them; inputs are validated flags, state and manifest entries; outputs are persisted metadata, graph/context items and per-contract pass/fail details.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380x2dqcf9ctzpr01bb7jr`; AC `ac_0000380x2dqce7nbcsa2c3b3g8`; tasks `task_0000380x2g10m4b5hef6hemvb4`, `task_0000380x2g280ctd6py5fq0bk4`, and `task_0000380x2g3eszpq5cxa1vtk8g`. Each design contract comment is task-owned before Loop Check.

<!-- tene:section:decisions -->
## Decisions

Manifest entries contain source requirement locator, implementation symbol regexes and command IDs. Auditor uses a fixed command registry, subprocess argument arrays and exact exit status.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:contract path="internal/domain/types.go" symbol="type MasterPlan" -->
<!-- tene:contract path="internal/tracecontext/context.go" symbol="func BuildContextPack" -->
<!-- tene:contract path="scripts/requirements-audit.py" symbol="def audit_contract" -->

<!-- tene:section:components -->
## Components

Master plan/intents, graph projection, context category allocator, semantic contract manifest, audit runner and mutation fixtures.

<!-- tene:section:interfaces -->
## Interfaces

`master create|status|validate`; enriched `intent capture|revise|supersede`; existing `graph build/trace/impact`; `context build/validate`; `requirements-audit.py [--final]`.

<!-- tene:section:data -->
## Data

Master goals/releases/milestones/risks/invariants; structured intent actors/outcomes/policies/rules/states/invariants/constraints/assumptions/questions/source; graph nodes/edges; contract requirement/source/symbols/commands.

<!-- tene:section:state-transitions -->
## State transitions

Intent candidate→confirmed→superseded/deprecated; master metadata revisioned by events; audit pending→each semantic proof→passed only if all proofs and workflow integrity pass.

<!-- tene:section:failures -->
## Failures

Unknown contract/command, missing source/symbol, command timeout/nonzero, open/deferred gap, unfinished task, unverified blocking AC, invalid archive or final active Sprint fail nonzero with structured detail.

<!-- tene:section:security -->
## Security

No manifest shell strings, no secret-bearing environment output, fixed command allowlist, bounded subprocess timeout and existing canary checks.

<!-- tene:section:tests -->
## Tests

Master/intent round trip and invalid inputs; graph node/edge completeness; context phase-doc/category/freshness/budget; audit missing symbol/nonzero command/state blocker mutations; full suite.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `871`
- Sprint: `sprint_0000380x2dmqv18kkk1fr627v0`
- Intents: `intent_0000380x2dqcf9ctzpr01bb7jr`
- Tasks: `task_0000380x2g10m4b5hef6hemvb4`, `task_0000380x2g280ctd6py5fq0bk4`, `task_0000380x2g3eszpq5cxa1vtk8g`

<!-- tene:generated:traceability:end -->
