---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380y1g1eyc1xeym60b0rtc
phase: plan
status: complete
revision: 1005
intent_ids: [intent_0000380y1xbkhac1273cdc3284]
generated_at: 2026-08-20T07:21:39Z
generated_by: tene-workflow
---

# plan — Portable Clean Workflow Matrix

<!-- tene:section:purpose -->
## Purpose

Build a reusable portable lifecycle runner, embed it in release smoke, map AC-PRODUCT-07 to it and dogfood the same Sprint gates.

<!-- tene:section:scope -->
## Scope

Python black-box orchestrator, release integration, semantic contract, docs, tests, QA/report/archive.

<!-- tene:section:layers -->
## Layers

Public CLI interface; workflow/gate business rules; isolated state/docs/evidence persistence; staged package and temp Git infrastructure.

<!-- tene:section:six-questions -->
## Six questions

Plan STACKS/run_stack definitions, release-smoke reference/caller, staged CLI input and stack summary/archive mutations.

<!-- tene:section:traceability -->
## Traceability

One implementation task covers runner; one release task covers stage integration and semantic evidence; both map ac_0000380y1xbkgdqhqk45j21axg.

<!-- tene:section:decisions -->
## Decisions

Run actual 32-revision lifecycle rather than mock commands. Seed a Git baseline so loop analysis evaluates Sprint changes, not fixture creation.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:work-packages -->
## Work packages

Implement deterministic runner; validate local wrapper; integrate staged release; run matrix/release/full gates; create evidence; report/archive/commit.

<!-- tene:section:dependencies -->
## Dependencies

Existing public CLI, document templates, loop analyzer, QA observation importer, package and release scripts.

<!-- tene:section:verification -->
## Verification

Direct local matrix, staged release smoke, semantic audit, make check, race/vet, doctor/evidence and post-archive final audit.

<!-- tene:section:risks -->
## Risks

Runtime cost grows about 30 seconds per release smoke. Keep fixtures minimal and deterministic; fail on any command or invariant.

<!-- tene:section:yagni -->
## Yagni

No containers, network package install, app servers or framework test suites for this portability contract.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `1005`
- Sprint: `sprint_0000380y1g1eyc1xeym60b0rtc`
- Intents: `intent_0000380y1xbkhac1273cdc3284`
- Tasks: `task_0000380y21dtymdms0ys2mcqyw`, `task_0000380y21f450p879vr1qgpv4`

<!-- tene:generated:traceability:end -->
