---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380xznb6mpm5bx0xpxex6w
phase: plan
status: complete
revision: 961
intent_ids: [intent_0000380y046tkxxteemz74gg58]
generated_at: 2026-08-20T07:05:37Z
generated_by: tene-workflow
---

# plan — Durable Journal Compaction

<!-- tene:section:purpose -->
## Purpose

Implement durable segmented compaction and prove it against loss, tampering, repeated runs, replay, doctor, and repository regression.

<!-- tene:section:scope -->
## Scope

State archive contract, compact/doctor surfaces, tests, storage docs, dogfood lifecycle, full QA and archive.

<!-- tene:section:layers -->
## Layers

Interface compact/doctor; business retention; persistence segment/manifest/checkpoint; infrastructure atomic filesystem writes.

<!-- tene:section:six-questions -->
## Six questions

Plan `JournalArchive`, `VerifyArchivedSegments`, `CreateCheckpoint` and `verifyEventBytes`, their callers, journal/project inputs and archive/checkpoint outputs.

<!-- tene:section:traceability -->
## Traceability

Tasks implement archive+anchor persistence and verification+doctor, both realizing `ac_0000380y046tkpdm8ey9panj70`.

<!-- tene:section:decisions -->
## Decisions

Use immutable segments plus one active checkpoint without renumbering events. Validate all earlier archives before producing a new segment.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:work-packages -->
## Work packages

1. Add archive metadata and anchored chain verification. 2. Archive exact bytes and replace active journal atomically. 3. Add doctor verification. 4. Add shrink/equality/repeated/tamper tests and docs. 5. Full QA, dogfood compact, report/archive.

<!-- tene:section:dependencies -->
## Dependencies

Existing event hash, merge-patch replay, atomic writers, doctor and workflow state.

<!-- tene:section:verification -->
## Verification

Focused tests plus `make check`, race, vet, release smoke, requirements audit, evidence verify and doctor.

<!-- tene:section:risks -->
## Risks

A crash can leave an unreferenced duplicate archive but never lose source. Corrupted history must stop later compaction.

<!-- tene:section:yagni -->
## Yagni

No external event database, codec, background compactor or automatic canonical-history repair.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `961`
- Sprint: `sprint_0000380xznb6mpm5bx0xpxex6w`
- Intents: `intent_0000380y046tkxxteemz74gg58`
- Tasks: `task_0000380y0cqhsx7fjvbj0zwk94`, `task_0000380y0crx5qtznz6hsj9pw8`

<!-- tene:generated:traceability:end -->
