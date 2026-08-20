---
schema_version: 1.0.0
document_type: analysis
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wjfxtbq8sgcngq0tmz4
phase: loop-check
status: draft
revision: 242
intent_ids: []
generated_at: 2026-08-20T00:30:55Z
generated_by: tene-workflow
---

# loop-check — Workflow Approval Loop Completion

<!-- tene:section:purpose -->
## Purpose

Compare profile, approval, iteration, and gap state machines against AC-WF-01..04 before QA.

<!-- tene:section:scope -->
## Scope

Domain defaults, schema, workflow guards, CLI/state events, status projection, documentation and tests.

<!-- tene:section:layers -->
## Layers

Interface approval/phase/loop/status; business profile/validity/budget/disposition; persistence project+journal; no new infrastructure.

<!-- tene:section:six-questions -->
## Six questions

Approval and expanded Sprint/Gap are defined in domain types and defaulted by state. RequiredApproval, ApprovalValidity and CanTransitionWithApproval are defined in workflow, imported/called by app. Commands accept bounded metadata and return records/findings; only Store.Mutate changes state.

<!-- tene:section:traceability -->
## Traceability

AC-WF-01/02 use workflow matrix and CLI strict-boundary tests; AC-WF-03 uses max=2 exhaustion integration; AC-WF-04 uses debt defer and security rejection integration; AC-WF-05 uses full regressions and validators.

<!-- tene:section:decisions -->
## Decisions

Dry-run and mutation share one guard. Approval is consumed only inside successful transition mutation. Existing 1.0 state defaults max iterations to five and approvals to an empty map.

<!-- tene:section:freeform -->
## Freeform

The user instruction “계속해” is recorded as the human authorization for this Sprint's strict design→do boundary; the exact approval is now consumed.

<!-- tene:section:baseline -->
## Baseline

Profiles were labels only, ApprovalRefs had no record, loop iterations were unbounded, and gap resolve/defer metadata was incomplete.

<!-- tene:section:changed-artifacts -->
## Changed artifacts

Domain/state/workflow/app/schema/tests, README/CLI reference, tene-sprint and tene-loop-check skills, Sprint documents and state events.

<!-- tene:section:gap-matrix -->
## Gap matrix

No blocking spec-code gap remains. Graph validation is valid and only reports the four expected pre-QA evidence warnings.

<!-- tene:section:iterations -->
## Iterations

Iteration 1 recorded `passed`. Focused tests caught the expected standard-profile archive regression and updated the integration journey to request, approve, and consume a report→archive approval.

<!-- tene:section:regression -->
## Regression

All Go tests and vet pass. Full race, Playwright, plugin/skill, evidence and doctor checks run in QA.
