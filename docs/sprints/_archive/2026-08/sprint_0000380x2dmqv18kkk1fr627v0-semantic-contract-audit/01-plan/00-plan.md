---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380x2dmqv18kkk1fr627v0
phase: plan
status: complete
revision: 871
intent_ids: [intent_0000380x2dqcf9ctzpr01bb7jr]
generated_at: 2026-08-20T02:50:07Z
generated_by: tene-workflow
---

# plan — Semantic Contract Completion Audit

<!-- tene:section:purpose -->
## Purpose

Close audited semantic gaps, then prove every contract through deterministic commands.

<!-- tene:section:scope -->
## Scope

Domain/schema/CLI extensions, graph/context materialization, semantic manifest/auditor, mutation tests, reference/security/release verification and final traceability.

<!-- tene:section:layers -->
## Layers

CLI contracts; master/intent/graph/context business rules; journal/projection/schema persistence; audit/test/package runtime.

<!-- tene:section:six-questions -->
## Six questions

Task 1 owns domain and runtime symbols; task 2 owns manifest/auditor definitions and command execution; task 3 owns test/reference evidence. Inputs are docs/state/source; outputs are canonical state, graph/context packs and fail-closed audit results.

<!-- tene:section:traceability -->
## Traceability

All three tasks realize `ac_0000380x2dqce7nbcsa2c3b3g8` and the 33 source contract IDs.

<!-- tene:section:decisions -->
## Decisions

Implement only material MVP mismatches found by direct design-to-source review; preserve documented post-MVP boundaries.

<!-- tene:section:freeform -->
## Freeform

Audit commands are deduplicated and run without shell interpolation.

<!-- tene:section:work-packages -->
## Work packages

1. Add master metadata and complete intent/spec memory fields and lifecycle.
2. Materialize missing graph relationships, phase document/context categories and bounded budget behavior.
3. Define contract manifest with symbol patterns and executable commands; make the auditor execute and report each ID.
4. Add negative mutations and run race/vet/check/release/reference/security proof.

<!-- tene:section:dependencies -->
## Dependencies

Domain/schema first, CLI and graph/context second, manifest/audit third, final proof last.

<!-- tene:section:verification -->
## Verification

Focused package tests, schema JSON validation, audit self-mutations, `go test -race ./...`, `go vet ./...`, `make check`, release smoke and `requirements-audit.py --final` after archive.

<!-- tene:section:risks -->
## Risks

Self-authored mappings can overclaim; require symbol regex resolution plus commands that exercise public behaviors and explicit negative tests.

<!-- tene:section:yagni -->
## Yagni

No generalized external PM database, graph UI, LSP host or hosted audit service.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `871`
- Sprint: `sprint_0000380x2dmqv18kkk1fr627v0`
- Intents: `intent_0000380x2dqcf9ctzpr01bb7jr`
- Tasks: `task_0000380x2g10m4b5hef6hemvb4`, `task_0000380x2g280ctd6py5fq0bk4`, `task_0000380x2g3eszpq5cxa1vtk8g`

<!-- tene:generated:traceability:end -->
