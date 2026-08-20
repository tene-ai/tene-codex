---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wjfxtbq8sgcngq0tmz4
phase: design
status: draft
revision: 242
intent_ids: []
generated_at: 2026-08-20T00:30:55Z
generated_by: tene-workflow
---

# design — Workflow Approval Loop Completion

<!-- tene:section:purpose -->
## Purpose

Define the executable state machines for profile approvals, loop budgets and gap dispositions.

<!-- tene:section:scope -->
## Scope

Domain/store/workflow/app/schema/skills. Approval remains independent of waiver: approval crosses an automation boundary; waiver accepts a specific unresolved, waivable gap.

<!-- tene:section:layers -->
## Layers

CLI adapters; pure RequiredApproval/ValidApproval guards; Project maps and Sprint loop counters; journal events and schema validation.

<!-- tene:section:six-questions -->
## Six questions

`RequiredApproval(profile,from,to)` and `ValidApproval(project,sprint,id,from,to,now)` live in workflow and are called by phase transition. Approval CRUD and loop/gap actions live in app. They consume bounded flags and return records/findings; Store.Mutate is the sole state writer.

<!-- tene:section:traceability -->
## Traceability

AC-WF-01 profile guard, AC-WF-02 approval domain, AC-WF-03 loop state, AC-WF-04 gap lifecycle, AC-WF-05 status/schema/skills/tests.

<!-- tene:section:decisions -->
## Decisions

Approved records are reusable only for their exact scope until consumed; successful transition marks the approval consumed atomically. Dry-run validates but never consumes. Quality findings are evaluated before approval and remain blockers.

<!-- tene:section:freeform -->
## Freeform

Default loop maximum is five per Sprint and can be set at Sprint creation within 1..20.

<!-- tene:section:components -->
## Components

Domain Approval and expanded Gap/Sprint; state ensureMaps; workflow approval policy; app approval/phase/loop/status; project schema and CLI/skill docs; unit/integration tests.

<!-- tene:section:interfaces -->
## Interfaces

`sprint create --max-iterations N`; `approval request --from --to --reason --requester --expires`, `approval approve ID --approver`, `approval list`; `phase transition PHASE --approval ID [--dry-run]`; `loop iterate --outcome repair|passed|blocked --summary`; `loop resolve-gap ID --resolution --evidence`; `loop defer-gap ID --reason --owner --target-sprint`.

<!-- tene:section:data -->
## Data

Approval fields: id/sprint/from/to/reason/requester/approver/status/requested/approved/expires/consumed. Sprint: loop_iteration/max_loop_iterations/last_loop_outcome. Gap: deferred reason/owner/target/time and resolution.

<!-- tene:section:state-transitions -->
## State transitions

Approval requested→approved→consumed or expired-by-validation. Gap open→resolved or deferred. Loop iteration 0..max; repair increments, passed records completion, blocked/exhausted prevents forward false-pass until gaps are handled.

<!-- tene:section:failures -->
## Failures

Stable errors for approval required/invalid/expired/scope/consumed, iteration exhausted, invalid outcome/category, missing resolution, and forbidden defer.

<!-- tene:section:security -->
## Security

Approver identity and reason are metadata only; no secret input. Approval cannot waive security, evidence, intent or QA guards.

<!-- tene:section:tests -->
## Tests

Profile matrix, approval scope/expiry/consumption/dry-run, loop boundary, gap category/status/defer restrictions, old-state defaulting, CLI status and full regressions.
