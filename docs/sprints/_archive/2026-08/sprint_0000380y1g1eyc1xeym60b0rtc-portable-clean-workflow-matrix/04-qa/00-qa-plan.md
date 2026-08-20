---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380y1g1eyc1xeym60b0rtc
phase: qa
status: active
revision: 1005
intent_ids: [intent_0000380y1xbkhac1273cdc3284]
generated_at: 2026-08-20T07:21:39Z
generated_by: tene-workflow
---

# qa — Portable Clean Workflow Matrix

<!-- tene:section:purpose -->
## Purpose

Verify portable workflow proof at stack, phase, evidence, archive and release boundaries.

<!-- tene:section:scope -->
## Scope

Local and staged matrix, seven variants, release/package regression, full tests/race/vet/audit/doctor.

<!-- tene:section:layers -->
## Layers

L1 build through L7 post-archive audit.

<!-- tene:section:six-questions -->
## Six questions

Observe runner/release callers, clean roots, command inputs, state/doc/evidence/archive outputs and nonzero failure propagation.

<!-- tene:section:traceability -->
## Traceability

All seven cases bind ac_0000380y1xbkgdqhqk45j21axg; stack summary binds AC-PRODUCT-07.

<!-- tene:section:decisions -->
## Decisions

No manual pass. Every stack must independently reach archive and healthy doctor.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:environment -->
## Environment

macOS local packaged platform, temp isolated Git repositories, staged self-contained plugin bundle.

<!-- tene:section:capabilities -->
## Capabilities

Python 3, Git, Go-built staged binary and filesystem; no network or secret provider needed.

<!-- tene:section:charters -->
## Charters

Happy three-stack completion; alternate new-process resume; empty clean state; validation/phase guards; isolated permissions; command failure propagation; post-archive recovery/doctor.

<!-- tene:section:ux-data-flow -->
## Ux data flow

Adopter installs package → invokes natural/public workflow → durable docs/state across processes → evidence gate → archive/report and resumable inactive status.

<!-- tene:section:evidence -->
## Evidence

Seven run/case-bound observations, a three-stack JSON summary and local+staged release transcript live under `04-qa/observations`.

<!-- tene:section:verdict -->
## Verdict

PASS. Local and staged matrices each completed Go CLI, Next.js web and Python API-worker lifecycles. Every clean project reached revision 32 with 33 events, seven evidence items, one archive manifest, inactive status and healthy doctor. `make check`, full race, vet, release smoke and semantic audit passed.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `1005`
- Sprint: `sprint_0000380y1g1eyc1xeym60b0rtc`
- Intents: `intent_0000380y1xbkhac1273cdc3284`
- Tasks: `task_0000380y21dtymdms0ys2mcqyw`, `task_0000380y21f450p879vr1qgpv4`

<!-- tene:generated:traceability:end -->
