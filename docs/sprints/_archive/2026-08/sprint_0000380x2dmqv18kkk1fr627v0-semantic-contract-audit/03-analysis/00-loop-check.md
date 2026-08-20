---
schema_version: 1.0.0
document_type: analysis
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380x2dmqv18kkk1fr627v0
phase: loop-check
status: complete
revision: 871
intent_ids: [intent_0000380x2dqcf9ctzpr01bb7jr]
generated_at: 2026-08-20T02:50:07Z
generated_by: tene-workflow
---

# loop-check — Semantic Contract Completion Audit

<!-- tene:section:purpose -->
## Purpose

Compare every MVP requirement, plan package and design contract against source symbols, public behavior and current workflow integrity.

<!-- tene:section:scope -->
## Scope

33 named contracts, all changed artifacts, master/intent/graph/context/error/replay semantics and audit negative mutations.

<!-- tene:section:layers -->
## Layers

CLI entry contracts, business invariants, persisted/replayed state and executable tool/release infrastructure were checked as one system.

<!-- tene:section:six-questions -->
## Six questions

The semantic manifest records definition name/path and commands; package tests cover callers, accepted inputs, returned envelopes/state and persistence/runtime effects.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380x2dqcf9ctzpr01bb7jr`; AC `ac_0000380x2dqce7nbcsa2c3b3g8`; three completed tasks and resolved gap `gap_0000380x4qdsp7zcbx8z5z6bg8`.

<!-- tene:section:decisions -->
## Decisions

Archived legacy evidence is accepted only when its hash is valid, redaction passed and its Sprint has a passed authoritative QA run; active evidence requires run/case binding.

<!-- tene:section:freeform -->
## Freeform

No waiver, deferred item or post-MVP scope promotion.

<!-- tene:section:baseline -->
## Baseline

The old audit verified path existence and state summaries but could not detect a missing symbol, unknown command or behavior-test failure.

<!-- tene:section:changed-artifacts -->
## Changed artifacts

Domain/app/state/tracecontext packages and tests, project schema, CLI reference, semantic contract manifest, audit runner/tests and canonical Sprint state/docs.

<!-- tene:section:gap-matrix -->
## Gap matrix

Initial automatic gap: design omitted the literal AC ID. It was repaired; the next analysis detected 0 and resolved 1. Semantic audit reports 0 contract failures and 0 command failures; only this Sprint's pre-QA AC remains intentionally unverified.

<!-- tene:section:iterations -->
## Iterations

Iteration 1 passed after enriched contracts, projection null normalization, semantic replay comparison and full executable audit.

<!-- tene:section:regression -->
## Regression

All 19 fixed command groups pass; missing-symbol and unknown-command mutations fail closed; projection doctor is healthy apart from the expected pre-QA evidence warning.


<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `871`
- Sprint: `sprint_0000380x2dmqv18kkk1fr627v0`
- Intents: `intent_0000380x2dqcf9ctzpr01bb7jr`
- Tasks: `task_0000380x2g10m4b5hef6hemvb4`, `task_0000380x2g280ctd6py5fq0bk4`, `task_0000380x2g3eszpq5cxa1vtk8g`

<!-- tene:generated:traceability:end -->
