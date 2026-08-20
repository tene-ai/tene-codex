---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380x5ytygknc3dgjzabepm
phase: design
status: complete
revision: 910
intent_ids: [intent_0000380x5z1z88fessgksyy8ac]
generated_at: 2026-08-20T03:21:02Z
generated_by: tene-workflow
---

# design — Bounded Active Projection Repair

<!-- tene:section:purpose -->
## Purpose

Define one canonical bounded active projection contract.

<!-- tene:section:scope -->
## Scope

Projection fields, writers, replay comparison, repair and tests.

<!-- tene:section:layers -->
## Layers

Resume JSON→selection helper→atomic derived file→journal replay/doctor.

<!-- tene:section:six-questions -->
## Six questions

`activeProjection(Project)` in `internal/state/store.go` is used by Initialize/Migrate/Mutate/CreateCheckpoint/Repair; returns a map containing project/profile/master, current Sprint objects and checkpoint.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380x5z1z88fessgksyy8ac`; AC `ac_0000380x5z1z9x273dcf4hk5vw`; task `task_0000380x5znht0btzehgzrht2w`.

<!-- tene:section:decisions -->
## Decisions

Archived/no-active state returns only identity, revision, profile, master and empty active ID.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:contract path="internal/state/store.go" symbol="func activeProjection" -->

<!-- tene:section:components -->
## Components

Active selector, atomic writers, replay comparator/repair and integration test.

<!-- tene:section:interfaces -->
## Interfaces

Internal `activeProjection(*domain.Project) map[string]any`; unchanged public status/context commands.

<!-- tene:section:data -->
## Data

Identity/master/current Sprint/tasks/intents/criteria/open gaps/waivers/approvals/last QA/checkpoint; graph/evidence/history excluded.

<!-- tene:section:state-transitions -->
## State transitions

Active Sprint updates replace bounded content; archive removes the active detail set.

<!-- tene:section:failures -->
## Failures

Replay difference reports paths and doctor repairs derived active file from journal.

<!-- tene:section:security -->
## Security

No secret vault/source/environment data is included.

<!-- tene:section:tests -->
## Tests

Active smaller than project, lacks graph/evidence/history, retains required resume fields, matches replay, repairs and shrinks after archive.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `910`
- Sprint: `sprint_0000380x5ytygknc3dgjzabepm`
- Intents: `intent_0000380x5z1z88fessgksyy8ac`
- Tasks: `task_0000380x5znht0btzehgzrht2w`

<!-- tene:generated:traceability:end -->
