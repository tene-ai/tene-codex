---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wmbmks792hdhe03j7bg
phase: design
status: draft
revision: 290
intent_ids: []
generated_at: 2026-08-20T00:47:14Z
generated_by: tene-workflow
---

# design — State Recovery and Deterministic Projection

<!-- tene:section:purpose -->
## Purpose

Specify journal projection checkpoints, compact patches and safe repair.

<!-- tene:section:scope -->
## Scope

State store and app doctor/compact plus graph ID generation.

<!-- tene:section:layers -->
## Layers

Doctor/compact interface; JSON tree diff/apply/replay; journal and projection files; atomic backup/fsync/rename.

<!-- tene:section:six-questions -->
## Six questions

`CreateCheckpoint`, `Replay`, `ProjectionDrift`, `RepairFromJournal`, `mergePatch`, `applyMergePatch`, and deterministicEdgeID; state/app/buildGraph call them; inputs are verified events/current files; outputs are reconstructed state and repair records.

<!-- tene:section:traceability -->
## Traceability

Each component maps directly to AC-ST-01..05.

<!-- tene:section:decisions -->
## Decisions

Event payload envelope has `data` and optional `projection_patch`; checkpoint payload has `projection`. Hash chain covers both. Replay begins at latest checkpoint/ProjectInitialized and applies later patches in sequence.

<!-- tene:section:freeform -->
## Freeform

Unknown legacy events before checkpoint are intentionally irrelevant; unknown post-checkpoint event without a patch is corruption.

<!-- tene:section:components -->
## Components

state/replay.go, store mutation/checkpoint/repair, app doctor/compact, deterministic graph edge helper, tests.

<!-- tene:section:interfaces -->
## Interfaces

`compact` creates a checkpoint and snapshot; `doctor` reports replay/project/active/master drift; `doctor --repair` backs up and replaces projections.

<!-- tene:section:data -->
## Data

Payload envelope `{data,projection_patch}`; checkpoint `{projection}`; drift `{path,expected_hash,actual_hash,status}`.

<!-- tene:section:state-transitions -->
## State transitions

load current→mutate clone→calculate patch→append fsynced event→atomic project/active/master. Repair verify→replay→backup→atomic projections.

<!-- tene:section:failures -->
## Failures

Invalid chain/base/patch, missing checkpoint, future schema, backup/write failure and conflict return STATE_* errors without journal mutation.

<!-- tene:section:security -->
## Security

Never include `.tene`; project workflow state is already non-secret. Backups use repository workflow directory and restrictive scope.

<!-- tene:section:tests -->
## Tests

Round-trip nested maps/arrays/deletes, legacy checkpoint anchor, corruption/truncation, projection deletion/corruption, failed prefix, two-writer conflict and deterministic graph rebuild.
