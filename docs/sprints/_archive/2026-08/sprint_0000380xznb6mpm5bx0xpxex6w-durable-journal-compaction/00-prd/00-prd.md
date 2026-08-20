---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380xznb6mpm5bx0xpxex6w
phase: prd
status: complete
revision: 961
intent_ids: [intent_0000380y046tkxxteemz74gg58]
generated_at: 2026-08-20T07:05:37Z
generated_by: tene-workflow
---

# prd — Durable Journal Compaction

<!-- tene:section:purpose -->
## Purpose

Close the FR-08/FR-09 retention gap: compact the active event journal without deleting canonical history or changing the replayed workflow meaning.

<!-- tene:section:scope -->
## Scope

Archive exact journal bytes, persist checksum/sequence/hash-chain metadata, replace the active journal with one anchored full-projection checkpoint, verify all archived segments in doctor, and fail closed on tampering.

<!-- tene:section:layers -->
## Layers

- Interface: `tene-workflow compact` result and `doctor` archive health.
- Business Logic: archive anchoring, repeated compaction, semantic replay equality.
- Persistence: `event-archive/*.ndjson`, manifest, checkpoint snapshot, atomic active journal replacement.
- Infrastructure: filesystem fsync/rename and corruption exit behavior.

<!-- tene:section:six-questions -->
## Six questions

- Names: `JournalArchive`, `VerifyArchivedSegments`, `CreateCheckpoint`, `verifyEventBytes`.
- Defined in: `internal/state/replay.go` and `internal/state/store.go`.
- Referenced by: compact and doctor application commands plus state tests.
- Called at: explicit compact, pre-compaction validation, and doctor.
- Input: current project projection and hash-chained NDJSON journal.
- Output/mutation: immutable archived segment+manifest, one-event active checkpoint, projections and snapshot; corruption returns `STATE_CORRUPT`.

<!-- tene:section:traceability -->
## Traceability

FR-08 Durable State; FR-09 Archive/Compact/Clear; WP-02 State Store; AC `ac_0000380y046tkpdm8ey9panj70`.

<!-- tene:section:decisions -->
## Decisions

Global event sequence numbers and previous hashes remain unchanged across segment boundaries. The first active checkpoint carries the immediate archive anchor; manifests independently authenticate archived bytes.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:problem -->
## Problem

Before this Sprint, compact only appended a full checkpoint. The repository's active journal remained 7.8MB and grew forever despite the PRD requiring completed events to move into an archive.

<!-- tene:section:actors -->
## Actors

Workflow users, Codex sessions resuming state, auditors, and recovery tooling.

<!-- tene:section:journeys -->
## Journeys

User compacts → old bytes are verified and archived → manifest is written → active journal becomes one checkpoint → replay/doctor verify the same project → later mutations continue the original global sequence. Tampered history blocks doctor and further compaction.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

- `ac_0000380y046tkpdm8ey9panj70` (blocking): active journal shrinks to one replayable checkpoint; prior bytes are checksum- and chain-verifiable.
- Repeated compaction retains every segment and sequence anchor.
- Archived-byte tampering is detected and no replacement occurs after the failure.

<!-- tene:section:non-goals -->
## Non goals

Deleting history, external databases, log compression codecs, or changing Sprint document archives.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `961`
- Sprint: `sprint_0000380xznb6mpm5bx0xpxex6w`
- Intents: `intent_0000380y046tkxxteemz74gg58`
- Tasks: `task_0000380y0cqhsx7fjvbj0zwk94`, `task_0000380y0crx5qtznz6hsj9pw8`

<!-- tene:generated:traceability:end -->
