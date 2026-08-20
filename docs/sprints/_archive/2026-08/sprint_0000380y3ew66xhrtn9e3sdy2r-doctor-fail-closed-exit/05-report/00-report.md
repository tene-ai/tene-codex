---
schema_version: 1.0.0
document_type: report
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380y3ew66xhrtn9e3sdy2r
phase: report
status: draft
revision: 1081
intent_ids: [intent_0000380y3hff9tbrqh38mvehb8]
generated_at: 2026-08-20T07:38:48Z
generated_by: tene-workflow
---

# report — Doctor Fail Closed Exit

<!-- tene:section:purpose -->
## Purpose

Record fail-closed doctor gate implementation and proof.

<!-- tene:section:scope -->
## Scope

App/tests/docs/state/QA/archive.

<!-- tene:section:layers -->
## Layers

Four layers.

<!-- tene:section:six-questions -->
## Six questions

Final report covers runtime.doctor and test names/files/callers/inputs/outputs.

<!-- tene:section:traceability -->
## Traceability

ac_0000380y3hff8thtbtry95g2zc.

<!-- tene:section:decisions -->
## Decisions

Preserve diagnostic result under structured error details.

<!-- tene:section:freeform -->
## Freeform

This section records the fail-closed doctor gate from the freeform perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:previous-sprints -->
## Previous sprints

This section records the fail-closed doctor gate from the previous sprints perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:changed-files -->
## Changed files

`internal/app/app.go` builds one doctor result and classifies blocker exits; `internal/app/app_test.go` asserts state corruption exit/details; `scripts/qa-doctor-exit-observations.py` generates QA evidence; Sprint state/docs record the lifecycle.

<!-- tene:section:intent-fulfillment -->
## Intent fulfillment

Doctor is now safe as a shell/CI harness gate: success means healthy, workflow blockers use exit 3, state-integrity blockers use exit 7, and the full result remains available at `errors[0].details`.

<!-- tene:section:qa-verdict -->
## Qa verdict

PASS: 7/7 variants, focused success/tamper integration, `make check`, race, vet, portable release and healthy doctor.

<!-- tene:section:deferred-work -->
## Deferred work

None.

<!-- tene:section:next-sprint -->
## Next sprint

This section records the fail-closed doctor gate from the next sprint perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `1081`
- Sprint: `sprint_0000380y3ew66xhrtn9e3sdy2r`
- Intents: `intent_0000380y3hff9tbrqh38mvehb8`
- Tasks: `task_0000380y3mjbzq6kbm80a0y2mg`

<!-- tene:generated:traceability:end -->

<!-- tene:generated:summary:start -->
### Generated Sprint Summary

- Sprint: `sprint_0000380y3ew66xhrtn9e3sdy2r`
- Previous sprints: `sprint_0000380y2tr1qg8yzy7gc6dgm0`
- Intent IDs: `intent_0000380y3hff9tbrqh38mvehb8`
- Tasks: 1
- QA verdict: `passed`
- Open gaps: 0
- State revision: 1082

<!-- tene:generated:summary:end -->
