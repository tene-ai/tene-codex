---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wjfxtbq8sgcngq0tmz4
phase: plan
status: draft
revision: 242
intent_ids: []
generated_at: 2026-08-20T00:30:55Z
generated_by: tene-workflow
---

# plan — Workflow Approval Loop Completion

<!-- tene:section:purpose -->
## Purpose

Implement AC-WF-01..05 as a backward-compatible state and CLI vertical slice.

<!-- tene:section:scope -->
## Scope

Domain fields, state defaulting/schema, workflow guards, CLI commands, status projection, skills/references, table tests, integration journey, QA/report/archive.

<!-- tene:section:layers -->
## Layers

Interface parsing and envelopes; business guard/state machine; persistence event journal/project projection; existing infrastructure unchanged.

<!-- tene:section:six-questions -->
## Six questions

Each public type/function and command caller will be enumerated in the report with inputs, outputs, references and mutations.

<!-- tene:section:traceability -->
## Traceability

Tasks map to profile/approval, loop budget, gap lifecycle, and integration/documentation ACs.

<!-- tene:section:decisions -->
## Decisions

Use optional fields under schema 1.0.0 with deterministic defaults so existing 1.0 state remains readable; no semantic migration is required.

<!-- tene:section:freeform -->
## Freeform

The live cross-major migration gap remains outside this Sprint.

<!-- tene:section:work-packages -->
## Work packages

1. Extend domain/state/schema with approvals and loop metadata.
2. Add profile boundary and scoped approval validation to workflow guards.
3. Add approval commands and phase `--approval` parsing.
4. Add bounded `loop iterate`, validated record/resolve/defer gap commands, and richer status.
5. Update skills/references and run all gates.

<!-- tene:section:dependencies -->
## Dependencies

Existing event journal, transition engine, waiver domain, document lifecycle, and Graph Context Freshness predecessor.

<!-- tene:section:verification -->
## Verification

Table-driven profile/approval tests, CLI integration for expiry/scope/iterations/defer, make check, race, vet, Playwright regression, validators, evidence verify and doctor.

<!-- tene:section:risks -->
## Risks

Existing strict repositories could become blocked: approval is required only at future guarded boundaries and remediation is explicit. Optional map initialization prevents nil-state failures.

<!-- tene:section:yagni -->
## Yagni

No web UI, signatures, organization directory, scheduler, or trust score.
