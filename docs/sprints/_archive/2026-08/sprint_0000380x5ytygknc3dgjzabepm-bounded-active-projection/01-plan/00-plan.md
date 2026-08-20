---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380x5ytygknc3dgjzabepm
phase: plan
status: complete
revision: 910
intent_ids: [intent_0000380x5z1z88fessgksyy8ac]
generated_at: 2026-08-20T03:21:02Z
generated_by: tene-workflow
---

# plan — Bounded Active Projection Repair

<!-- tene:section:purpose -->
## Purpose

Replace full-state active writes with one tested projection function.

<!-- tene:section:scope -->
## Scope

All four write/repair paths, projection contract test, semantic manifest update and final repository measurement.

<!-- tene:section:layers -->
## Layers

Resume interface, selection logic, JSON projection and replay/doctor verification.

<!-- tene:section:six-questions -->
## Six questions

Definition/callers/inputs/output/effects are the `activeProjection` helper, Store initialization/mutation/migration and replay checkpoint/repair writers, Project input, bounded JSON output and atomic replacement.

<!-- tene:section:traceability -->
## Traceability

AC `ac_0000380x5z1z9x273dcf4hk5vw`; task `task_0000380x5znht0btzehgzrht2w`.

<!-- tene:section:decisions -->
## Decisions

Retain current QA run for active gating, but never copy graph, evidence collection or archived Sprint maps.

<!-- tene:section:freeform -->
## Freeform

No compatibility read path changes because commands load canonical `project.json`.

<!-- tene:section:work-packages -->
## Work packages

Implement helper; route all writes/expected replay/repair; add bounded and archive tests; update semantic audit; run full proof.

<!-- tene:section:dependencies -->
## Dependencies

State projection helper before writer replacement and tests.

<!-- tene:section:verification -->
## Verification

State/app tests, mutation+replay equivalence, doctor, active/project byte measurement, final audit/check.

<!-- tene:section:risks -->
## Risks

Omitting a live guard field could weaken resume; explicitly include active AC, gaps, approvals/waivers and last QA checkpoint.

<!-- tene:section:yagni -->
## Yagni

No external state database or journal pruning.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `910`
- Sprint: `sprint_0000380x5ytygknc3dgjzabepm`
- Intents: `intent_0000380x5z1z88fessgksyy8ac`
- Tasks: `task_0000380x5znht0btzehgzrht2w`

<!-- tene:generated:traceability:end -->
