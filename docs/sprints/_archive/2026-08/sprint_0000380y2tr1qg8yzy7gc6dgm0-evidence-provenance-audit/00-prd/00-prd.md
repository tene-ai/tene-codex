---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380y2tr1qg8yzy7gc6dgm0
phase: prd
status: complete
revision: 1043
intent_ids: [intent_0000380y2x0gjdsw44cv68h5d4]
generated_at: 2026-08-20T07:33:18Z
generated_by: tene-workflow
---

# prd — Evidence Provenance Audit

<!-- tene:section:purpose -->
## Purpose

Remove blanket archived-evidence grandfathering from the completion audit and require passed-case provenance for every blocking AC.

<!-- tene:section:scope -->
## Scope

Audit state algorithm, legacy compatibility rule, negative self-tests, full 45-AC verification and release regression.

<!-- tene:section:layers -->
## Layers

Interface: audit exit/JSON. Business: AC verification rule. Persistence: project run/case/evidence records and artifact hashes. Infrastructure: Python unittest and Make release gate.

<!-- tene:section:six-questions -->
## Six questions

Names: state_findings and provenance test; defined in requirements-audit.py/tests; called by Make/final audit; input project/evidence; output unverified failures without mutation.

<!-- tene:section:traceability -->
## Traceability

FR-07, AC-PRODUCT-04, WP-10 and ac_0000380y2x0gjmx70245j0j39r.

<!-- tene:section:decisions -->
## Decisions

Legacy omission is compatible only when the passed case explicitly lists the evidence ID. Present run/case metadata must match exactly; sprint, AC, hash and redaction must always match.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:problem -->
## Problem

The former archive-wide loop credited any valid evidence AC from a passed archived Sprint even if no passed QA case referenced it.

<!-- tene:section:actors -->
## Actors

Release auditor, maintainer and Codex agent.

<!-- tene:section:journeys -->
## Journeys

Load state → traverse passed run/case → validate linked evidence/provenance/hash → credit intersecting AC → fail every remaining blocker.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

Unlinked legacy evidence fails; linked field-absent legacy evidence passes; mismatched metadata fails; current 45 blocking ACs remain verified.

<!-- tene:section:non-goals -->
## Non goals

Mutating old evidence or weakening modern run/case bindings.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `1043`
- Sprint: `sprint_0000380y2tr1qg8yzy7gc6dgm0`
- Intents: `intent_0000380y2x0gjdsw44cv68h5d4`
- Tasks: `task_0000380y2zs85jqvd91fwzptj4`

<!-- tene:generated:traceability:end -->
