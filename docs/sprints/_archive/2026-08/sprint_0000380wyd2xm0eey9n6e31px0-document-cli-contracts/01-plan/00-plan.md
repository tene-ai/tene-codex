---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wyd2xm0eey9n6e31px0
phase: plan
status: complete
revision: 719
intent_ids: [intent_0000380wyd7d99vkcqmwbmz158]
generated_at: 2026-08-20T02:15:00Z
generated_by: tene-workflow
---

# plan — Document Sync and CLI Contract Completion

<!-- tene:section:purpose -->
## Purpose

Close each public CLI/document contract mismatch with executable tests.

<!-- tene:section:scope -->
## Scope

Document engine, runtime dispatch, domain lifecycle, schemas, references and contract matrix.

<!-- tene:section:layers -->
## Layers

Interface flags/aliases, business lifecycle/dedup rules, persisted request/waiver data and document filesystem writes.

<!-- tene:section:six-questions -->
## Six questions

Definitions reside in document/domain/app; callers are CLI commands and skills; inputs are args/state/files; outputs are previews, mutations, cached responses and stable errors.

<!-- tene:section:traceability -->
## Traceability

AC `ac_0000380wyd7d8zhc4mwqd86wkm`; tasks `task_0000380wz9ygjx9rs8dkgkrtf0` and `task_0000380wza05c2zmfcjch8s92m`.

<!-- tene:section:decisions -->
## Decisions

Request retry returns the cached first response; evidence dedup keys content hash, kind and AC set.

<!-- tene:section:freeform -->
## Freeform

Compatibility aliases do not weaken authorization or validation.

<!-- tene:section:work-packages -->
## Work packages

Document sync; CLI flags/aliases; request/evidence dedup; waiver request/approve/expire; schemas/docs/tests.

<!-- tene:section:dependencies -->
## Dependencies

Domain schema before runtime behavior, then integration tests and reference updates.

<!-- tene:section:verification -->
## Verification

Golden sync, retry conflict, evidence identity, waiver approval, quiet/alias matrix, race/vet/full suite.

<!-- tene:section:risks -->
## Risks

File/event cross-resource atomicity is bounded by explicit event and recovery checks.

<!-- tene:section:yagni -->
## Yagni

No editor UI or distributed idempotency service.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `719`
- Sprint: `sprint_0000380wyd2xm0eey9n6e31px0`
- Intents: `intent_0000380wyd7d99vkcqmwbmz158`
- Tasks: `task_0000380wz9ygjx9rs8dkgkrtf0`, `task_0000380wza05c2zmfcjch8s92m`

<!-- tene:generated:traceability:end -->
