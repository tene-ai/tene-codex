---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wmbmks792hdhe03j7bg
phase: qa
status: draft
revision: 290
intent_ids: []
generated_at: 2026-08-20T00:47:14Z
generated_by: tene-workflow
---

# qa — State Recovery and Deterministic Projection

<!-- tene:section:purpose -->
## Purpose

Prove recovery from journal rather than the damaged projection.

<!-- tene:section:scope -->
## Scope

Patch round-trip, checkpoint replay, corrupted projection repair, hash corruption rejection, concurrency, deterministic graph and full regression.

<!-- tene:section:layers -->
## Layers

CLI diagnostics, replay engine, files, atomic locking.

<!-- tene:section:six-questions -->
## Six questions

Assert recovery public functions, app callers, JSON inputs, Project outputs and exact backup/projection mutations.

<!-- tene:section:traceability -->
## Traceability

Four blocking ACs receive happy/error/recovery cases and hash-valid evidence.

<!-- tene:section:decisions -->
## Decisions

Unit temp repositories safely corrupt projections/journals; live repository remains verified and matched.

<!-- tene:section:freeform -->
## Freeform

Playwright is regression only.

<!-- tene:section:environment -->
## Environment

macOS local Go/Node/Codex validators.

<!-- tene:section:capabilities -->
## Capabilities

Unit/race/vet/filesystem crash fixtures, Playwright, plugin/skill, evidence and doctor.

<!-- tene:section:charters -->
## Charters

Replay equality; missing/drift repair; corrupt chain fail; concurrent stale writers; repeated graph IDs; clear/compact preservation.

<!-- tene:section:ux-data-flow -->
## Ux data flow

compact→checkpoint→patched mutations→doctor replay/diff→optional backup+repair→resume.

<!-- tene:section:evidence -->
## Evidence

`evidence_0000380wn6ah0a492y3bff470g` records replay/repair/patch/concurrency/determinism and all regression gates.

<!-- tene:section:verdict -->
## Verdict

PASS: run `run_0000380wn16nndrx67dynzxp6w`, 12/12 cases, four blocking ACs.
