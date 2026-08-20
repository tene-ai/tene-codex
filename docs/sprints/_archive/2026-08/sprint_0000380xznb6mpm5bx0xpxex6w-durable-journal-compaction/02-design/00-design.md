---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380xznb6mpm5bx0xpxex6w
phase: design
status: complete
revision: 961
intent_ids: [intent_0000380y046tkxxteemz74gg58]
generated_at: 2026-08-20T07:05:37Z
generated_by: tene-workflow
---

# design — Durable Journal Compaction

<!-- tene:section:purpose -->
## Purpose

Define executable segmented-journal compaction that genuinely bounds active history.

<!-- tene:section:scope -->
## Scope

`JournalArchive`, active anchors, archive validator, transaction order, doctor and failure behavior.

<!-- tene:section:layers -->
## Layers

Interface returns archive/count; business preserves sequence and meaning; persistence stores exact bytes/manifest/checkpoint; infrastructure uses atomic writers.

<!-- tene:section:six-questions -->
## Six questions

`CreateCheckpoint` reads project+journal, calls both verifiers, writes archive/manifest/checkpoint/snapshot/projections. `VerifyArchivedSegments` reads manifests+bytes, calls `verifyEventBytes`, returns metadata or corruption. Compact/doctor call them.

<!-- tene:section:traceability -->
## Traceability

`ac_0000380y046tkpdm8ey9panj70` maps to state shrink/equality/tamper tests and dogfood evidence.

<!-- tene:section:decisions -->
## Decisions

Store `event-archive/segment-<first>-<last>-<sha-prefix>.ndjson` plus manifest with paths, SHA-256, counts, sequence boundary, incoming anchor and terminal hash.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:components -->
## Components

State archive schema/verifiers; checkpoint compactor; app compact/doctor adapters; unit and integration tests.

<!-- tene:section:interfaces -->
## Interfaces

Compact returns snapshot, archived segment and `active_events:1`. Doctor returns segment count and `STATE_ARCHIVE_CORRUPT` finding.

<!-- tene:section:data -->
## Data

Exact current NDJSON is hashed/archived. First active checkpoint embeds the manifest reference and full project projection.

<!-- tene:section:state-transitions -->
## State transitions

`lock → verify active+archives → archive bytes+manifest → revision++ → checkpoint(sequence=last+1, previous=last hash) → replace active journal → projections+snapshot`. Later mutations continue global sequence.

<!-- tene:section:failures -->
## Failures

Bad anchor, missing/tampered archive, checksum/count/terminal mismatch fail with `STATE_CORRUPT` before active replacement.

<!-- tene:section:security -->
## Security

No secret content is introduced; manifests contain relative paths and hashes only. Clear cannot delete source archives.

<!-- tene:section:tests -->
## Tests

Shrink/equality, mutation continuation, repeated compaction, initialized-only boundary, tamper rejection and unchanged active journal on failure.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `961`
- Sprint: `sprint_0000380xznb6mpm5bx0xpxex6w`
- Intents: `intent_0000380y046tkxxteemz74gg58`
- Tasks: `task_0000380y0cqhsx7fjvbj0zwk94`, `task_0000380y0crx5qtznz6hsj9pw8`

<!-- tene:generated:traceability:end -->
