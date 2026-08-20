---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wmbmks792hdhe03j7bg
phase: prd
status: draft
revision: 290
intent_ids: []
generated_at: 2026-08-20T00:47:14Z
generated_by: tene-workflow
---

# prd — State Recovery and Deterministic Projection

<!-- tene:section:purpose -->
## Purpose

Make the append-only journal and checksummed checkpoints the recoverable source of truth, and make derived graph rebuilds deterministic.

<!-- tene:section:scope -->
## Scope

Checkpoint+merge-patch events, semantic replay, projection comparison/repair with backup, crash-prefix and stale-writer tests, deterministic graph edge IDs, compact/clear compatibility and recovery diagnostics.

<!-- tene:section:layers -->
## Layers

Interface doctor/compact; business replay/diff; persistence journal/checkpoint/project/active/master; infrastructure fsync/rename/lock.

<!-- tene:section:six-questions -->
## Six questions

Replay, Checkpoint and RepairProjections are defined in state and called by doctor/compact; they consume verified events and return Project/diff/repaired paths. Mutations append hashed events and atomically replace projections.

<!-- tene:section:traceability -->
## Traceability

WP-02 state-store gates, design state source-of-truth/recovery/concurrency contracts, AC-PRODUCT-02 resume fidelity, and reliability budgets.

<!-- tene:section:decisions -->
## Decisions

Legacy history is anchored by an explicit full checkpoint. Subsequent events store RFC7396-style projection patches, not repeated full projects. Replay fails closed without a usable base. Graph edge IDs derive from semantic edge content.

<!-- tene:section:freeform -->
## Freeform

Graph is derived but remains in projection; deterministic IDs eliminate no-op rebuild churn.

<!-- tene:section:problem -->
## Problem

Current doctor verifies the journal but rebuilds active/master from possibly corrupted project.json. Historical payloads are insufficient for full replay and graph rebuilds generate random edge IDs.

<!-- tene:section:actors -->
## Actors

Developer recovering a repository, CI checking drift, Codex resuming after interruption.

<!-- tene:section:journeys -->
## Journeys

Create checkpoint, mutate state, corrupt/delete projections, run doctor diff then repair from replay; interrupt projection write and replay committed journal prefix; rebuild graph twice and observe semantic equality.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

- AC-ST-01: verified checkpoint+patch replay equals current semantic project state.
- AC-ST-02: doctor reports projection drift and repair backs up then rebuilds project/active/master exclusively from replay.
- AC-ST-03: journal corruption/truncation fails closed and committed-prefix crash recovery is deterministic.
- AC-ST-04: stale concurrent writers yield exactly one commit and graph rebuild IDs are stable.
- AC-ST-05: compact/clear preserve source and all regressions/validators pass.

<!-- tene:section:non-goals -->
## Non goals

Distributed locks, remote event store, or automatic repair of a corrupt journal.
