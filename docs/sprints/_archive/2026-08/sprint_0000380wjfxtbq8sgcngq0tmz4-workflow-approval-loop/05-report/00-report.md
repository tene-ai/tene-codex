---
schema_version: 1.0.0
document_type: report
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wjfxtbq8sgcngq0tmz4
phase: report
status: draft
revision: 242
intent_ids: []
generated_at: 2026-08-20T00:30:55Z
generated_by: tene-workflow
---

# report — Workflow Approval Loop Completion

<!-- tene:section:purpose -->
## Purpose

Record completion of enforceable workflow profiles, scoped human approvals, bounded loop iteration, and auditable gap disposition.

<!-- tene:section:scope -->
## Scope

Domain/state/schema, transition guards, CLI/status, tests, skills/references, QA and archive authorization.

<!-- tene:section:layers -->
## Layers

- Interface: approval request/approve/list, phase approval flag, loop iterate/defer/resolve, richer status.
- Business Logic: profile matrix, approval validity, consumption, iteration maximum, category and disposition rules.
- Persistence: Approval map, Sprint loop fields, Gap resolution/defer fields and append-only events.
- Infrastructure: existing atomic state store and validators; no new runtime dependency.

<!-- tene:section:six-questions -->
## Six questions

- Names: Approval, RequiredApproval, ApprovalValidity, CanTransitionWithApproval, loop and gap fields/commands.
- Definitions: `internal/domain/types.go`, `internal/workflow/workflow.go`, `internal/app/app.go`, `internal/state/store.go`.
- References: app imports workflow/domain/state; project schema and skills describe the public contract.
- Callers: phase transition, approval command, loop command, status, graph builder and tests.
- Inputs: profile/from/to/approval ID/time, requester/approver/reason/expiry, iteration outcome/summary, gap category/evidence/owner/target.
- Outputs/mutations: stable findings/records/status; successful mutations append events and update project/active/master projections atomically.

<!-- tene:section:traceability -->
## Traceability

AC-WF-01..04 passed in run `run_0000380wkft6h3n298wyrt1jdr`; AC-WF-05 is supported by evidence `evidence_0000380wkk3af0vb6ymwvb3fmr`. Implements WP-03 and WP-09 planned contracts.

<!-- tene:section:decisions -->
## Decisions

Approvals authorize only exact transitions and never bypass quality. Strict requires design→do and archive approvals; standard archive only; light/off retain logical phases without approval gates. Default loop maximum is five.

<!-- tene:section:freeform -->
## Freeform

Status now distinguishes raw open/deferred gaps from effective blockers and active waivers, avoiding the prior misleading single count.

<!-- tene:section:previous-sprints -->
## Previous sprints

Builds on Graph Context Freshness (`sprint_0000380wh5p6fktd02fj5fzjcm`): approval nodes and workflow state now participate in the rebuilt graph and context can expose the resulting current revision.

<!-- tene:section:changed-files -->
## Changed files

Updated domain, state store, workflow, app and tests; project schema; README English/Korean; CLI reference; tene-sprint and tene-loop-check skills; Sprint lifecycle documents and evidence.

<!-- tene:section:intent-fulfillment -->
## Intent fulfillment

An agent can no longer infer permission from profile labels, loop forever, resolve a gap without evidence, or silently defer safety work. Each action has a bounded state record and stable failure.

<!-- tene:section:qa-verdict -->
## Qa verdict

PASS: 12/12 QA cases, four blocking ACs, complete regression matrix, plugin/skills valid, evidence and journal healthy.

<!-- tene:section:deferred-work -->
## Deferred work

Trust-score automation and remote approval UI remain non-goals. Semantic journal replay/recovery is the next state Sprint. The pre-existing cross-major migration gap remains governed separately.

<!-- tene:section:next-sprint -->
## Next sprint

State and Recovery Completion: semantic event replay, crash/prefix simulation, deterministic graph IDs, projection comparison and repair from journal source of truth.

<!-- tene:generated:summary:start -->
### Generated Sprint Summary

- Sprint: `sprint_0000380wjfxtbq8sgcngq0tmz4`
- Previous sprints: `sprint_0000380wh5p6fktd02fj5fzjcm`
- Intent IDs: `intent_0000380wjjwe31ffek75hgbm0r`
- Tasks: 4
- QA verdict: `passed`
- Open gaps: 0
- State revision: 285

<!-- tene:generated:summary:end -->
