---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380y1g1eyc1xeym60b0rtc
phase: design
status: complete
revision: 1005
intent_ids: [intent_0000380y1xbkhac1273cdc3284]
generated_at: 2026-08-20T07:21:39Z
generated_by: tene-workflow
---

# design — Portable Clean Workflow Matrix

<!-- tene:section:purpose -->
## Purpose

Specify the clean-stack black-box workflow matrix at implementable command and invariant level.

<!-- tene:section:scope -->
## Scope

STACKS fixture schema, call boundary, lifecycle sequence, observation generation, archive assertions and release wiring.

<!-- tene:section:layers -->
## Layers

Interface staged scripts/tene-workflow; business phase/loop/QA/report gates; persistence isolated .tene-workflow/docs; infrastructure package/temp/Git/Python.

<!-- tene:section:six-questions -->
## Six questions

STACKS and run_stack are defined in portable-workflow-smoke.py; release-smoke invokes them. Inputs are --cli and optional workspace; outputs JSON stack summaries while mutating only isolated roots.

<!-- tene:section:traceability -->
## Traceability

ac_0000380y1xbkgdqhqk45j21axg → three stack entries → full command chain → seven evidence cases → archive/doctor assertions.

<!-- tene:section:decisions -->
## Decisions

Seed Go cmd, Next app/route, Python API/worker. Public wrapper is absolute. JSON envelopes are parsed and any nonzero/not-ok fails immediately.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:components -->
## Components

Argument parser, call, complete_documents, run_stack, STACKS, release-smoke integration and semantic audit locator.

<!-- tene:section:interfaces -->
## Interfaces

portable-workflow-smoke.py --cli PATH [--workspace PATH]; stdout JSON with passed/stacks/workspace.

<!-- tene:section:data -->
## Data

Each result includes stack, revision, event count, evidence count and archive manifest. Observations bind run/case/spec/state plus L1–L7 and before/after.

<!-- tene:section:state-transitions -->
## State transitions

init → draft/prd/plan/design/do/loop-check/qa/report/archived; task and intent states, QA evaluate and post-archive verification are mandatory.

<!-- tene:section:failures -->
## Failures

Invalid JSON, nonzero envelope, loop failure, QA failure, doctor/evidence failure or missing archive raises nonzero and blocks release.

<!-- tene:section:security -->
## Security

Fixtures contain no secrets; observations contain only stack/path identifiers. Temporary workspace is removed unless explicit for debugging.

<!-- tene:section:tests -->
## Tests

Local wrapper and staged package must each return three revision-32 archived results; release smoke retains tamper/update/uninstall checks.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `1005`
- Sprint: `sprint_0000380y1g1eyc1xeym60b0rtc`
- Intents: `intent_0000380y1xbkhac1273cdc3284`
- Tasks: `task_0000380y21dtymdms0ys2mcqyw`, `task_0000380y21f450p879vr1qgpv4`

<!-- tene:generated:traceability:end -->
