---
schema_version: 1.0.0
document_type: analysis
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wmbmks792hdhe03j7bg
phase: loop-check
status: draft
revision: 290
intent_ids: []
generated_at: 2026-08-20T00:47:14Z
generated_by: tene-workflow
---

# loop-check — State Recovery and Deterministic Projection

<!-- tene:section:purpose -->
## Purpose

Verify implementation against recovery PRD/design before QA.

<!-- tene:section:scope -->
## Scope

Checkpoint, merge patch, replay, drift/repair, deterministic graph and concurrency.

<!-- tene:section:layers -->
## Layers

Doctor/compact interface; replay business rules; journal/projection persistence; atomic lock/fsync infrastructure.

<!-- tene:section:six-questions -->
## Six questions

Recovery functions are defined in state/replay.go, called by store/app doctor/compact, consume verified events/JSON projections and return replayed Project/drift/paths. Only checkpoint, normal mutation and repair write files.

<!-- tene:section:traceability -->
## Traceability

AC-ST-01 replay tests/live doctor; AC-ST-02 corruption repair test; AC-ST-03 hash verification; AC-ST-04 concurrent writer and deterministic graph tests; AC-ST-05 regressions.

<!-- tene:section:decisions -->
## Decisions

Legacy state was anchored at checkpoint revision 303; patches through revision 318 replay to exact matches.

<!-- tene:section:freeform -->
## Freeform

No `.codegraph/` was created. Graph IDs are semantic SHA-256 prefixes.

<!-- tene:section:baseline -->
## Baseline

Repair copied project.json and graph edges churned random IDs.

<!-- tene:section:changed-artifacts -->
## Changed artifacts

state/replay.go and tests; store/app/graph tests; CLI/skill/design docs; Sprint state.

<!-- tene:section:gap-matrix -->
## Gap matrix

No blocking gap. Live doctor reports all three projections match after multiple post-checkpoint mutations.

<!-- tene:section:iterations -->
## Iterations

Iteration 1 passed. Initial compile found doctor finding declaration order and legacy migration test checkpoint assumptions; both were corrected.

<!-- tene:section:regression -->
## Regression

Focused Go suite passes; full gates follow in QA.
