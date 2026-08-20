---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wmbmks792hdhe03j7bg
phase: plan
status: draft
revision: 290
intent_ids: []
generated_at: 2026-08-20T00:47:14Z
generated_by: tene-workflow
---

# plan — State Recovery and Deterministic Projection

<!-- tene:section:purpose -->
## Purpose

Deliver AC-ST-01..05 as a recoverable state vertical slice.

<!-- tene:section:scope -->
## Scope

Event envelope/patch codec, checkpoint, replay, compare/repair, deterministic graph IDs, tests, docs and lifecycle evidence.

<!-- tene:section:layers -->
## Layers

CLI, replay logic, durable files, atomic filesystem operations.

<!-- tene:section:six-questions -->
## Six questions

Report every exported recovery function, its state/app callers, verified inputs, projections and file mutations.

<!-- tene:section:traceability -->
## Traceability

Tasks cover replay/checkpoint, repair/doctor, atomicity/concurrency, deterministic graph, integration/docs.

<!-- tene:section:decisions -->
## Decisions

Patch generation compares canonical JSON trees. Checkpoint is created before relying on patch events and stores its own post-revision projection.

<!-- tene:section:freeform -->
## Freeform

Do not pretend pre-checkpoint legacy deltas are replayable.

<!-- tene:section:work-packages -->
## Work packages

1. Projection patch codec and event envelope.
2. Checkpoint and replay.
3. Projection diff, backup and repair.
4. Deterministic graph IDs and recovery tests.
5. Full QA, report, archive, commit/push.

<!-- tene:section:dependencies -->
## Dependencies

Hash-chain journal, atomic writer, workflow state, Graph Context and Workflow Approval predecessors.

<!-- tene:section:verification -->
## Verification

Property/snapshot replay, corruption/crash/concurrency tests, make/race/vet/Playwright, validators, evidence and doctor.

<!-- tene:section:risks -->
## Risks

Large patches or legacy gaps: checkpoint bounds history; patches are only post-checkpoint. Failed repair preserves backup and journal.

<!-- tene:section:yagni -->
## Yagni

No network replication or event schema generator.
