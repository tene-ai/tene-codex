---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380x5ytygknc3dgjzabepm
phase: prd
status: complete
revision: 910
intent_ids: [intent_0000380x5z1z88fessgksyy8ac]
generated_at: 2026-08-20T03:21:02Z
generated_by: tene-workflow
---

# prd — Bounded Active Projection Repair

<!-- tene:section:purpose -->
## Purpose

Keep cross-session resume state bounded without losing current workflow context.

<!-- tene:section:scope -->
## Scope

`active.json` projection, mutation/checkpoint/replay/repair writers, size/content/recovery tests and FR-08 semantic mapping.

<!-- tene:section:layers -->
## Layers

Status/resume interface, projection selection logic, active file persistence and journal/doctor infrastructure.

<!-- tene:section:six-questions -->
## Six questions

`activeProjection` is defined in state, called by every active projection writer, consumes canonical Project, and returns current master/Sprint/tasks/intents/AC/gaps/decisions/QA checkpoint without graph/history.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380x5z1z88fessgksyy8ac`; AC `ac_0000380x5z1z9x273dcf4hk5vw`; task `task_0000380x5znht0btzehgzrht2w`.

<!-- tene:section:decisions -->
## Decisions

`project.json` and the journal stay canonical; `active.json` is a replaceable derived projection.

<!-- tene:section:freeform -->
## Freeform

No source history or evidence is deleted.

<!-- tene:section:problem -->
## Problem

The prior active projection duplicated the complete multi-megabyte graph/evidence project state.

<!-- tene:section:actors -->
## Actors

Plugin user and a new Codex session resuming work.

<!-- tene:section:journeys -->
## Journeys

Mutation→bounded active write→restart/replay→doctor comparison/repair→archive→minimal inactive projection.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

`ac_0000380x5z1z9x273dcf4hk5vw`: no graph/history in active projection; current resume fields retained; smaller than project; replay and doctor consistent; minimal after archive.

<!-- tene:section:non-goals -->
## Non goals

No deletion/segmentation of canonical project state or journal history.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `910`
- Sprint: `sprint_0000380x5ytygknc3dgjzabepm`
- Intents: `intent_0000380x5z1z88fessgksyy8ac`
- Tasks: `task_0000380x5znht0btzehgzrht2w`

<!-- tene:generated:traceability:end -->
