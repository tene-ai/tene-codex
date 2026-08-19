---
schema_version: 1.0.0
document_type: analysis
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v1ack4pe5nddv67gxaw
phase: loop-check
status: draft
revision: 59
intent_ids: []
generated_at: 2026-08-19T17:21:14Z
generated_by: tene-workflow
---

# loop-check — Archive evidence integrity stabilization

<!-- tene:section:purpose -->
## Purpose

Confirm archive integrity fixes match the stabilization PRD and design.

<!-- tene:section:scope -->
## Scope

Review the application diff, integration assertions and real predecessor evidence paths.

<!-- tene:section:layers -->
## Layers

Interface remains compatible; business helper is shared; persistence updates all locators; infrastructure doctor re-verifies history.

<!-- tene:section:six-questions -->
## Six questions

Definitions, callers, input state and locator mutations for `transition` and `invalidEvidence` match Design. No new public symbol is untraced.

<!-- tene:section:traceability -->
## Traceability

Two completed tasks map to two blocking ACs and focused integration assertions.

<!-- tene:section:decisions -->
## Decisions

No waiver is required; archive manifest and post-archive evidence are blocking.

<!-- tene:section:freeform -->
## Freeform

The fix preserves the predecessor evidence hash and content.

<!-- tene:section:baseline -->
## Baseline

The preceding MVP Sprint at state revision 58 exposed stale-path risk after its archive transition.

<!-- tene:section:changed-artifacts -->
## Changed artifacts

`internal/app/app.go`, its integration test, CLI reference, archive manifest behavior and this stabilization Sprint.

<!-- tene:section:gap-matrix -->
## Gap matrix

No blocking gap remains. The integration test now uses evidence inside the Sprint directory and asserts both manifest and relocated URI. Dogfood QA also found that doctor validated blank future-phase documents; doctor is now phase-aware and validates the complete set only for archived Sprints.

<!-- tene:section:iterations -->
## Iterations

One implementation iteration, one archive locator verification, and one phase-aware doctor repair iteration.

<!-- tene:section:regression -->
## Regression

Go tests and vet pass; the previous archived evidence remains verifiable and doctor remains healthy.
