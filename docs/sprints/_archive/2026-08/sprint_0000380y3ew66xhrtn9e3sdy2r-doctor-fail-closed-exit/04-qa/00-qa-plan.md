---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380y3ew66xhrtn9e3sdy2r
phase: qa
status: active
revision: 1081
intent_ids: [intent_0000380y3hff9tbrqh38mvehb8]
generated_at: 2026-08-20T07:38:48Z
generated_by: tene-workflow
---

# qa — Doctor Fail Closed Exit

<!-- tene:section:purpose -->
## Purpose

Prove doctor can be safely used as a shell gate.

<!-- tene:section:scope -->
## Scope

Focused app integration, full regression, release, audit, evidence and healthy dogfood.

<!-- tene:section:layers -->
## Layers

L1–L7.

<!-- tene:section:six-questions -->
## Six questions

Observe exit, error code/details/findings and healthy success.

<!-- tene:section:traceability -->
## Traceability

Seven variants bind ac_0000380y3hff8thtbtry95g2zc.

<!-- tene:section:decisions -->
## Decisions

No manual pass; state tamper must be exit 7.

<!-- tene:section:freeform -->
## Freeform

This section records the fail-closed doctor gate from the freeform perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:environment -->
## Environment

This section records the fail-closed doctor gate from the environment perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:capabilities -->
## Capabilities

This section records the fail-closed doctor gate from the capabilities perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:charters -->
## Charters

This section records the fail-closed doctor gate from the charters perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:ux-data-flow -->
## Ux data flow

This section records the fail-closed doctor gate from the ux data flow perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:evidence -->
## Evidence

Seven run/case-bound observations and the focused healthy/tampered integration transcript are under `04-qa/observations`.

<!-- tene:section:verdict -->
## Verdict

PASS. Healthy doctor returns exit 0; tampered archive returns `STATE_CORRUPT` exit 7 with `healthy:false` and findings in error details. `make check`, race, vet, portable staged release, semantic audit and dogfood doctor pass.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `1081`
- Sprint: `sprint_0000380y3ew66xhrtn9e3sdy2r`
- Intents: `intent_0000380y3hff9tbrqh38mvehb8`
- Tasks: `task_0000380y3mjbzq6kbm80a0y2mg`

<!-- tene:generated:traceability:end -->
