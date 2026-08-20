---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380y2tr1qg8yzy7gc6dgm0
phase: plan
status: complete
revision: 1043
intent_ids: [intent_0000380y2x0gjdsw44cv68h5d4]
generated_at: 2026-08-20T07:33:18Z
generated_by: tene-workflow
---

# plan — Evidence Provenance Audit

<!-- tene:section:purpose -->
## Purpose

Refactor the evidence-credit algorithm and prove backward-compatible fail-closed behavior.

<!-- tene:section:scope -->
## Scope

One audit implementation, one fixture self-test, semantic/full regression and Sprint evidence.

<!-- tene:section:layers -->
## Layers

Audit interface, provenance business rule, persisted state/artifact input, unittest/release infrastructure.

<!-- tene:section:six-questions -->
## Six questions

Plan state_findings and test callers/inputs/outputs; no product state mutation outside dogfood lifecycle.

<!-- tene:section:traceability -->
## Traceability

One task maps implementation+test to ac_0000380y2x0gjmx70245j0j39r.

<!-- tene:section:decisions -->
## Decisions

Delete archive-wide credit loop; credit only inside passed case traversal.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:work-packages -->
## Work packages

Implement rule; add three-state test; audit 45 ACs; full QA; report/archive.

<!-- tene:section:dependencies -->
## Dependencies

Existing state schema, historical evidence links and hash artifacts.

<!-- tene:section:verification -->
## Verification

Focused unittest, audit no-exec/full/final, make check, race/vet, evidence and doctor.

<!-- tene:section:risks -->
## Risks

Over-strict migration could invalidate honest history; explicit omission-only compatibility prevents that.

<!-- tene:section:yagni -->
## Yagni

No state migration or evidence rewriting.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `1043`
- Sprint: `sprint_0000380y2tr1qg8yzy7gc6dgm0`
- Intents: `intent_0000380y2x0gjdsw44cv68h5d4`
- Tasks: `task_0000380y2zs85jqvd91fwzptj4`

<!-- tene:generated:traceability:end -->
