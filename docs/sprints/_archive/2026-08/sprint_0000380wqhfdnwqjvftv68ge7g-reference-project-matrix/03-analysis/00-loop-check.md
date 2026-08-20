---
schema_version: 1.0.0
document_type: analysis
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wqhfdnwqjvftv68ge7g
phase: loop-check
status: draft
revision: 411
intent_ids: []
generated_at: 2026-08-20T01:15:02Z
generated_by: tene-workflow
---

# loop-check — Reference Project Portability Matrix

<!-- tene:section:purpose -->
## Purpose

Compare portability implementation to AC-PRODUCT-03/07 and Scenario D.

<!-- tene:section:scope -->
## Scope

Discovery, layers, generic components, fixtures and matrix assertions.

<!-- tene:section:layers -->
## Layers

All four layers are present in mature fixture; polyglot mappings and unknowns are explicit.

<!-- tene:section:six-questions -->
## Six questions

Analyze/sourceExtension/classify live in codeintel, are invoked by CLI/tests, accept roots and return complete or explicitly unknown component fields with provider provenance.

<!-- tene:section:traceability -->
## Traceability

Blocking AC has one deterministic matrix test plus focused fallback test.

<!-- tene:section:decisions -->
## Decisions

Unknown is an acceptable answer; omission or fabricated semantic certainty is a gap.

<!-- tene:section:freeform -->
## Freeform

Review found `.tene` was not excluded and corrected it before gate.

<!-- tene:section:baseline -->
## Baseline

Only Go files and greenfield runtime journey were visible.

<!-- tene:section:changed-artifacts -->
## Changed artifacts

`internal/codeintel/*`, `testdata/reference-mature/*`, `testdata/reference-polyglot/*`, Sprint/state files.

<!-- tene:section:gap-matrix -->
## Gap matrix

Three types pass; mature 4/4 layers; polyglot 4 files with explicit unknowns; vault files zero; open gaps zero.

<!-- tene:section:iterations -->
## Iterations

Initial fallback test exposed vault indexing; skip rule repaired and regression added.

<!-- tene:section:regression -->
## Regression

Focused matrix/fallback tests pass.
