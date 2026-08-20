---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380y2tr1qg8yzy7gc6dgm0
phase: design
status: complete
revision: 1043
intent_ids: [intent_0000380y2x0gjdsw44cv68h5d4]
generated_at: 2026-08-20T07:33:18Z
generated_by: tene-workflow
---

# design — Evidence Provenance Audit

<!-- tene:section:purpose -->
## Purpose

Specify precise evidence-credit predicates.

<!-- tene:section:scope -->
## Scope

Passed run/case traversal and modern/legacy binding compatibility.

<!-- tene:section:layers -->
## Layers

Audit JSON; predicate logic; state/artifact reads; test/release runner.

<!-- tene:section:six-questions -->
## Six questions

state_findings reads root/final, references project maps and artifact files, and returns failures/active ID without mutation. Unit test imports it directly.

<!-- tene:section:traceability -->
## Traceability

ac_0000380y2x0gjmx70245j0j39r maps to test_legacy_evidence_requires_passed_case_link_and_matching_provenance.

<!-- tene:section:decisions -->
## Decisions

valid = file+SHA+redaction+sprint; binding = omitted-or-equal run AND omitted-or-equal case; credit = case AC ∩ evidence AC, only for evidence_ids listed by passed case.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:components -->
## Components

state_findings and RequirementsAuditTests fixture.

<!-- tene:section:interfaces -->
## Interfaces

requirements-audit.py flags/output unchanged.

<!-- tene:section:data -->
## Data

Historical missing keys mean schema omission, never wildcard when a conflicting value exists.

<!-- tene:section:state-transitions -->
## State transitions

Read-only traversal; no transitions.

<!-- tene:section:failures -->
## Failures

Unlinked, failed run/case, missing/tampered file, wrong sprint, redaction failure, metadata mismatch or AC mismatch remain unverified.

<!-- tene:section:security -->
## Security

Hash and redaction remain mandatory before provenance credit.

<!-- tene:section:tests -->
## Tests

Unlinked fail → linked omission pass → wrong run fail; repository proves all blockers.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `1043`
- Sprint: `sprint_0000380y2tr1qg8yzy7gc6dgm0`
- Intents: `intent_0000380y2x0gjdsw44cv68h5d4`
- Tasks: `task_0000380y2zs85jqvd91fwzptj4`

<!-- tene:generated:traceability:end -->
