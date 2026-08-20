---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wjfxtbq8sgcngq0tmz4
phase: prd
status: draft
revision: 242
intent_ids: []
generated_at: 2026-08-20T00:30:55Z
generated_by: tene-workflow
---

# prd — Workflow Approval Loop Completion

<!-- tene:section:purpose -->
## Purpose

Complete WP-03 and WP-09 so workflow profiles have enforceable meaning, human approvals are durable and scoped, and loop repair cannot continue indefinitely or hide deferred work.

<!-- tene:section:scope -->
## Scope

- Durable approval request/approve/list lifecycle tied to one Sprint transition.
- Strict/standard/light/off profile policies applied by the same transition guard used for dry-run and mutation.
- Bounded loop iterations with explicit repair/passed/exhausted outcomes.
- Gap category validation, resolution rationale, and auditable defer lifecycle.
- Status output that separates raw open gaps, active waivers, effective blockers, deferred gaps, and iteration budget.

<!-- tene:section:layers -->
## Layers

Interface: approval/loop/status/phase CLI. Business: profile policy, approval validity, iteration and gap transitions. Persistence: approval/iteration/defer events in project+journal. Infrastructure: none beyond existing atomic store.

<!-- tene:section:six-questions -->
## Six questions

New domain names are Approval and LoopState fields; workflow exports RequiredApproval and ValidApproval; app commands create state events. Inputs are profile, transition, actor/reason/scope/expiry and gap disposition. Outputs are guard findings and durable records; mutation occurs only through Store.Mutate.

<!-- tene:section:traceability -->
## Traceability

Realizes WP-03 transition approvals/profiles and WP-09 iteration/deferred workflow. Evidence covers strict and standard boundaries, single-scope approval, expiry, dry-run equivalence, max iterations, and non-waivable defer rejection.

<!-- tene:section:decisions -->
## Decisions

Approvals never bypass quality guards. They authorize only a named from→to transition for one Sprint and expire. Strict requires approval for design→do and report→archive; standard requires report→archive; light/off add no approval boundary but never skip logical phases.

<!-- tene:section:freeform -->
## Freeform

This is deterministic harness policy, not autonomous permission escalation. The user remains the approver.

<!-- tene:section:problem -->
## Problem

Profiles currently validate as strings but behave identically. ApprovalRefs have no records or command. Loop check has no iteration budget, and gap resolution/defer lacks rationale and lifecycle metadata.

<!-- tene:section:actors -->
## Actors

Developer/owner requesting a transition, human approver accepting a bounded risk, Codex executing repairs, evaluator deciding whether gaps remain.

<!-- tene:section:journeys -->
## Journeys

1. Strict design→do dry-run fails with approval-required, a human creates/approves a scoped request, retry passes without bypassing other guards.
2. Codex records repair iterations; the sixth attempt is rejected after the default five-iteration budget.
3. A non-security debt is deferred with owner/reason/target Sprint and disappears from effective blockers while remaining visible.
4. Security/evidence-integrity blockers cannot be deferred or waived.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

- AC-WF-01: profile matrix enforces documented approval boundaries using stable guard codes in dry-run and actual transition.
- AC-WF-02: approvals have requester, approver, reason, exact transition scope, status, timestamps and expiry; invalid/expired/wrong-scope approvals fail closed.
- AC-WF-03: loop iteration is bounded, records outcome/summary, exposes remaining budget, and exhaustion blocks further repair claims.
- AC-WF-04: gaps validate categories, require resolution evidence/rationale, and may be deferred only with owner/reason/target while security/evidence-integrity remain non-deferrable.
- AC-WF-05: status, docs, skills, schemas and executable tests expose the new contract without regression.

<!-- tene:section:non-goals -->
## Non goals

Trust-score automation, remote approval UI, time-based background expiry jobs, or bypassing QA/evidence gates.
