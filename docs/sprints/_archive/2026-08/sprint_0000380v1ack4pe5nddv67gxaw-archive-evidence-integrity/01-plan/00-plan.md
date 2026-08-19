---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v1ack4pe5nddv67gxaw
phase: plan
status: draft
revision: 59
intent_ids: []
generated_at: 2026-08-19T17:21:14Z
generated_by: tene-workflow
---

# plan — Archive evidence integrity stabilization

<!-- tene:section:purpose -->
## Purpose

Close the archive locator gap with a small, separately traceable stabilization Sprint.

<!-- tene:section:scope -->
## Scope

Two tasks: archive locator/manifest mutation and active-independent verification/doctor behavior, followed by regression QA.

<!-- tene:section:layers -->
## Layers

CLI interface calls business invariants that update persisted locators; automated tests and doctor provide infrastructure verification.

<!-- tene:section:six-questions -->
## Six questions

Plan targets `transition`, `evidence`, `doctor`, and `invalidEvidence`; definitions, callers, inputs and mutations are documented in design/report.

<!-- tene:section:traceability -->
## Traceability

Each task maps to exactly one of the two blocking ACs.

<!-- tene:section:decisions -->
## Decisions

Keep this fix inside existing packages and public commands; add no new service or dependency.

<!-- tene:section:freeform -->
## Freeform

The preceding archived Sprint remains immutable except for its required archive manifest created during the discovered compatibility repair.

<!-- tene:section:work-packages -->
## Work packages

1. Update paths and generate manifest atomically around archive.
2. Decouple read-only evidence verification from active Sprint and include it in doctor.

<!-- tene:section:dependencies -->
## Dependencies

Locator mutation precedes post-archive verification; tests precede final evidence registration.

<!-- tene:section:verification -->
## Verification

Full Sprint integration asserts manifest, relocated URI and existing file. After real archive, run evidence verify, doctor, Go tests and vet.

<!-- tene:section:risks -->
## Risks

Filesystem move can succeed while state mutation fails; existing rollback restores the original tree on mutation error.

<!-- tene:section:yagni -->
## Yagni

No generic transaction coordinator or cross-filesystem archive copy is introduced.
