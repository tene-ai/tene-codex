---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380y3ew66xhrtn9e3sdy2r
phase: prd
status: complete
revision: 1081
intent_ids: [intent_0000380y3hff9tbrqh38mvehb8]
generated_at: 2026-08-20T07:38:48Z
generated_by: tene-workflow
---

# prd — Doctor Fail Closed Exit

<!-- tene:section:purpose -->
## Purpose

Make doctor a fail-closed CI/harness gate without losing diagnostic detail.

<!-- tene:section:scope -->
## Scope

Doctor healthy/unhealthy exit behavior, workflow/state classification and regression tests.

<!-- tene:section:layers -->
## Layers

Interface doctor JSON/exit; business blocker classification; persistence read-only diagnostics; infrastructure CI shell semantics.

<!-- tene:section:six-questions -->
## Six questions

Names runtime.doctor and commandError; app.go/app_test.go; called by CLI/release; inputs state/evidence; output success or structured error details.

<!-- tene:section:traceability -->
## Traceability

FR-07/08/09, CLI exit contract and ac_0000380y3hff8thtbtry95g2zc.

<!-- tene:section:decisions -->
## Decisions

Exit 3 for non-state blockers, 7 for STATE_* integrity blockers, 0 only when healthy.

<!-- tene:section:freeform -->
## Freeform

This section records the fail-closed doctor gate from the freeform perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:problem -->
## Problem

This section records the fail-closed doctor gate from the problem perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:actors -->
## Actors

This section records the fail-closed doctor gate from the actors perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:journeys -->
## Journeys

This section records the fail-closed doctor gate from the journeys perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

This section records the fail-closed doctor gate from the acceptance criteria perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:non-goals -->
## Non goals

This section records the fail-closed doctor gate from the non goals perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `1081`
- Sprint: `sprint_0000380y3ew66xhrtn9e3sdy2r`
- Intents: `intent_0000380y3hff9tbrqh38mvehb8`
- Tasks: `task_0000380y3mjbzq6kbm80a0y2mg`

<!-- tene:generated:traceability:end -->
