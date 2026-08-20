---
schema_version: 1.0.0
document_type: report
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380xznb6mpm5bx0xpxex6w
phase: report
status: draft
revision: 961
intent_ids: [intent_0000380y046tkxxteemz74gg58]
generated_at: 2026-08-20T07:05:37Z
generated_by: tene-workflow
---

# report — Durable Journal Compaction

<!-- tene:section:purpose -->
## Purpose

Record implementation, proof, continuity and retrospective for durable compaction.

<!-- tene:section:scope -->
## Scope

State/app/tests/docs, actual compact, QA evidence, remaining policy and predecessor.

<!-- tene:section:layers -->
## Layers

Interface compact/doctor; business archive invariants; persistence segments/manifests/checkpoints; infrastructure atomic filesystem.

<!-- tene:section:six-questions -->
## Six questions

Final report enumerates symbols, definition files, references/callers, journal/project inputs and archive/checkpoint/error outputs.

<!-- tene:section:traceability -->
## Traceability

`ac_0000380y046tkpdm8ey9panj70`, FR-08/09 and WP-02.

<!-- tene:section:decisions -->
## Decisions

Preserve global sequence/hash continuity and independently verify immutable archived source.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:previous-sprints -->
## Previous sprints

Follows bounded-active-projection: it reduced resume projection size; this bounds the event working set while retaining the audit trail.

<!-- tene:section:changed-files -->
## Changed files

`internal/state/replay.go` defines archive metadata and compaction transaction; `internal/state/store.go` verifies active/archived chains; `internal/app/app.go` exposes compact/doctor results. State/app tests cover shrink, replay, repeated compact, write failure, invalid anchor and tamper. The QA runner, semantic manifest, storage design, READMEs and Sprint documents were added/updated. Dogfood state now includes the archived segment and bounded active checkpoint.

<!-- tene:section:intent-fulfillment -->
## Intent fulfillment

Compact archived the exact 7,542,532-byte, 952-event journal under SHA-256 `6c0642b93e12794ebd47a44c42a889c657ba7c7f0c174fae68130147488dbeec` and replaced it with a 2,134,981-byte full checkpoint. Subsequent seven evidence mutations plus evaluation continued at global sequence 953–960. Replay and all three projections match; doctor independently verifies the archive.

<!-- tene:section:qa-verdict -->
## Qa verdict

PASS: 7/7 variants, L1–L7, valid evidence, `make check`, race, vet, release smoke, semantic audit and healthy doctor.

<!-- tene:section:deferred-work -->
## Deferred work

No in-scope deferment or policy decision. Archive compression and external event stores remain explicit non-goals.

<!-- tene:section:next-sprint -->
## Next sprint

Independently prove the complete public Sprint lifecycle in three clean stack repositories; current reference tests prove stack analysis and product journeys but not yet the whole workflow lifecycle per stack.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `961`
- Sprint: `sprint_0000380xznb6mpm5bx0xpxex6w`
- Intents: `intent_0000380y046tkxxteemz74gg58`
- Tasks: `task_0000380y0cqhsx7fjvbj0zwk94`, `task_0000380y0crx5qtznz6hsj9pw8`

<!-- tene:generated:traceability:end -->

<!-- tene:generated:summary:start -->
### Generated Sprint Summary

- Sprint: `sprint_0000380xznb6mpm5bx0xpxex6w`
- Previous sprints: `sprint_0000380x5ytygknc3dgjzabepm`
- Intent IDs: `intent_0000380y046tkxxteemz74gg58`
- Tasks: 2
- QA verdict: `passed`
- Open gaps: 0
- State revision: 962

<!-- tene:generated:summary:end -->
