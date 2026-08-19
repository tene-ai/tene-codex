---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v3xjbejnjk8w65p9q2w
phase: plan
status: draft
revision: 127
intent_ids: []
generated_at: 2026-08-19T17:43:57Z
generated_by: tene-workflow
---

# plan — Waiver Schema Migration and Recovery

<!-- tene:section:purpose -->
## Purpose

Deliver policy exceptions, migration and repair as one storage-safety vertical slice.

<!-- tene:section:scope -->
## Scope

Domain/schema, state APIs, gate logic, CLI, tests, references and recovery documentation.

<!-- tene:section:layers -->
## Layers

Interface commands call workflow/state ports; state owns atomic backup/write; workflow owns waiver gate semantics.

<!-- tene:section:six-questions -->
## Six questions

Tests exercise definition, import/caller, request fields, result fields and journal/projection mutations for every public operation.

<!-- tene:section:traceability -->
## Traceability

Waiver domain/CLI/gate → AC1; migration planner/apply → AC2; repair/strict validation → AC3.

<!-- tene:section:decisions -->
## Decisions

Legacy fixture is schema 0.9.0. Migration is idempotent at 1.0.0. Repair is derived-only and refuses a corrupt journal.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:work-packages -->
## Work packages

1. Add Waiver domain and workflow validity.
2. Add CLI lifecycle and graph/state projection.
3. Add migration inspect/dry-run/apply/backup/event.
4. Add doctor repair and referential/unknown-field checks.
5. Test, document, QA, archive, commit and push.

<!-- tene:section:dependencies -->
## Dependencies

Domain before state/workflow; state before CLI; fixtures before QA.

<!-- tene:section:verification -->
## Verification

Unit/CLI fixtures, expired/revoked/security negatives, migration backup equivalence, corrupt-journal refusal, race/vet/plugin/skill/evidence/doctor.

<!-- tene:section:risks -->
## Risks

Repair could imply journal replay that does not exist: explicitly limit it to derived projections. Migration failure could strand state: backup before replacement and atomic writes.

<!-- tene:section:yagni -->
## Yagni

No generic migration DSL, remote storage, arbitrary JSON patch or journal rewriting.
