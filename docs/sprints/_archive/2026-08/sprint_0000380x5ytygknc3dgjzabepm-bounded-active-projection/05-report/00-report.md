---
schema_version: 1.0.0
document_type: report
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380x5ytygknc3dgjzabepm
phase: report
status: active
revision: 910
intent_ids: [intent_0000380x5z1z88fessgksyy8ac]
generated_at: 2026-08-20T03:21:02Z
generated_by: tene-workflow
---

# report — Bounded Active Projection Repair

<!-- tene:section:purpose -->
## Purpose

Record repair of the final bounded-state gap discovered after semantic audit.

<!-- tene:section:scope -->
## Scope

Active projection content/writers/replay/repair/tests and final size.

<!-- tene:section:layers -->
## Layers

Resume interface; activeProjection logic; atomic active JSON; journal/doctor recovery.

<!-- tene:section:six-questions -->
## Six questions

`activeProjection` in state/store is called by Initialize/Migrate/Mutate/CreateCheckpoint/ProjectionDrift/Repair; accepts Project and returns current-only map atomically written to active.json.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380x5z1z88fessgksyy8ac`; AC `ac_0000380x5z1z9x273dcf4hk5vw`; task `task_0000380x5znht0btzehgzrht2w`; run `run_0000380x68hrhap7fkvxh4epxc`.

<!-- tene:section:decisions -->
## Decisions

Do not prune canonical graph/evidence/history; project and journal remain complete sources.

<!-- tene:section:freeform -->
## Freeform

This Sprint exists because file-size inspection found a requirement gap after the prior audit archive.

<!-- tene:section:previous-sprints -->
## Previous sprints

Follows semantic audit Sprint `sprint_0000380x2dmqv18kkk1fr627v0` and repairs its discovered FR-08/09 active-state weakness.

<!-- tene:section:changed-files -->
## Changed files

State store/replay writers, app integration test, FR-08 semantic mapping, QA fixture and Sprint state/docs/evidence.

<!-- tene:section:intent-fulfillment -->
## Intent fulfillment

Active state excludes graph/evidence/history while retaining master, active Sprint/task/intent/AC/gap/decision/QA checkpoint and remains replay-repairable.

<!-- tene:section:qa-verdict -->
## Qa verdict

Passed seven L1–L7 variants including injected drift, missing file and doctor recovery; 14.5KB during QA versus 2.9MB canonical project.

<!-- tene:section:deferred-work -->
## Deferred work

None. Canonical project/journal compaction policies remain implemented without destructive history deletion.

<!-- tene:section:next-sprint -->
## Next sprint

No required MVP Sprint remains.


<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `910`
- Sprint: `sprint_0000380x5ytygknc3dgjzabepm`
- Intents: `intent_0000380x5z1z88fessgksyy8ac`
- Tasks: `task_0000380x5znht0btzehgzrht2w`

<!-- tene:generated:traceability:end -->

<!-- tene:generated:summary:start -->
### Generated Sprint Summary

- Sprint: `sprint_0000380x5ytygknc3dgjzabepm`
- Previous sprints: `sprint_0000380x2dmqv18kkk1fr627v0`
- Intent IDs: `intent_0000380x5z1z88fessgksyy8ac`
- Tasks: 1
- QA verdict: `passed`
- Open gaps: 0
- State revision: 908

<!-- tene:generated:summary:end -->
