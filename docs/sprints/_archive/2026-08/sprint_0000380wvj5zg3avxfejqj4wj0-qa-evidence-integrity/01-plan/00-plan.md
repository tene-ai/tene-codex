---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wvj5zg3avxfejqj4wj0
phase: plan
status: approved
revision: 1
intent_ids: [intent_0000380wvk7q507samchfwpvxw]
generated_at: 2026-08-20T01:50:11Z
generated_by: tene-workflow
---

# plan — QA Evidence Integrity and Independent Evaluation

<!-- tene:section:purpose -->
## Purpose

Replace hash-only QA completion with proof-oriented evaluation.

<!-- tene:section:scope -->
## Scope

Domain evidence metadata, structured observation schema, CLI ingestion, deterministic evaluator, mutation tests and compatibility updates.

<!-- tene:section:layers -->
## Layers

Interface commands feed Business Logic evaluation; artifacts persist in Sprint QA storage; adapters are Infrastructure boundaries.

<!-- tene:section:six-questions -->
## Six questions

The work defines evidence/run assertion fields in domain types, consumes them from app commands, evaluates them in workflow, and returns stable blocker findings without mutating source specifications.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380wvk7q507samchfwpvxw` and AC `ac_0000380wvk7q5mkdwqj2q1c1gc`.

<!-- tene:section:decisions -->
## Decisions

Specification freshness uses a canonical hash rather than project revision alone because registering evidence advances the project revision.

<!-- tene:section:freeform -->
## Freeform

Evidence may record a failed assertion; failure evidence must remain importable and force a failed verdict.

<!-- tene:section:work-packages -->
## Work packages

1. Extend canonical evidence and QA-run contracts. 2. Harden observation and adapter ingestion. 3. Replace evaluator eligibility rules. 4. Add adversarial mutation tests and schemas.

<!-- tene:section:dependencies -->
## Dependencies

Domain → adapter/workflow → app → tests and documentation.

<!-- tene:section:verification -->
## Verification

Unit mutation matrix, CLI happy path using structured observations, manual-pass rejection, content hash verification, race/vet/full check.

<!-- tene:section:risks -->
## Risks

Older evidence lacks new metadata; only newly evaluated active runs require the strengthened contract, preserving archive readability.

<!-- tene:section:yagni -->
## Yagni

No remote judge, hosted evidence service or probabilistic scoring.
