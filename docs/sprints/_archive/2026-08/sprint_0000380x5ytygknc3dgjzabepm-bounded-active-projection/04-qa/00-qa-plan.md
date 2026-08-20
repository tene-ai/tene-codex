---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380x5ytygknc3dgjzabepm
phase: qa
status: complete
revision: 910
intent_ids: [intent_0000380x5z1z88fessgksyy8ac]
generated_at: 2026-08-20T03:21:02Z
generated_by: tene-workflow
---

# qa — Bounded Active Projection Repair

<!-- tene:section:purpose -->
## Purpose

Prove bounded active state, failure detection and deterministic repair.

<!-- tene:section:scope -->
## Scope

Current, alternate and empty states; corruption/missing/permission/recovery; L1–L7.

<!-- tene:section:layers -->
## Layers

Unit/contract/integration/system/resume journey/regression/operational recovery.

<!-- tene:section:six-questions -->
## Six questions

Evidence identifies activeProjection and writers, fixture commands/input states, JSON outputs and atomic repair effects.

<!-- tene:section:traceability -->
## Traceability

AC `ac_0000380x5z1z9x273dcf4hk5vw`; run `run_0000380x68hrhap7fkvxh4epxc`; seven observations.

<!-- tene:section:decisions -->
## Decisions

Size is measured with active QA included; archive shrink is verified after transition.

<!-- tene:section:freeform -->
## Freeform

No waiver/manual pass.

<!-- tene:section:environment -->
## Environment

Local temp repositories plus the real 2.9MB canonical project.

<!-- tene:section:capabilities -->
## Capabilities

Go state/app tests, CLI init/sprint/doctor/repair, filesystem size/content checks and structured observations.

<!-- tene:section:charters -->
## Charters

Bounded happy, alternate Sprint, no active Sprint, injected drift, forbidden key scan, missing file, repair recovery.

<!-- tene:section:ux-data-flow -->
## Ux data flow

Session resume reads current-only context; corruption is visible; repair derives from journal without altering canonical source.

<!-- tene:section:evidence -->
## Evidence

Case-bound L1–L7 observations plus current 14.5KB vs 2.9MB measurement and temp-repo failure/recovery outputs.

<!-- tene:section:verdict -->
## Verdict

Passed with zero findings; all seven required layers/variants evidenced.


<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `910`
- Sprint: `sprint_0000380x5ytygknc3dgjzabepm`
- Intents: `intent_0000380x5z1z88fessgksyy8ac`
- Tasks: `task_0000380x5znht0btzehgzrht2w`

<!-- tene:generated:traceability:end -->
