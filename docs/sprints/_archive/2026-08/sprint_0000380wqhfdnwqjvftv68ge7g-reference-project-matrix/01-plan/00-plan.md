---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wqhfdnwqjvftv68ge7g
phase: plan
status: draft
revision: 411
intent_ids: []
generated_at: 2026-08-20T01:15:02Z
generated_by: tene-workflow
---

# plan — Reference Project Portability Matrix

<!-- tene:section:purpose -->
## Purpose

Add portable fixtures and a release-blocking matrix test.

<!-- tene:section:scope -->
## Scope

Source discovery, fallback components, classification and reference repositories.

<!-- tene:section:layers -->
## Layers

All four Understanding Layers plus provider infrastructure.

<!-- tene:section:six-questions -->
## Six questions

Analyze is imported by CLI/matrix, reads eligible sources, and returns AST or filesystem components with provenance/confidence/unknowns.

<!-- tene:section:traceability -->
## Traceability

One vertical work package satisfies the blocking portability AC.

<!-- tene:section:decisions -->
## Decisions

Keep compact static fixtures and test semantics rather than fixture build systems.

<!-- tene:section:freeform -->
## Freeform

Rollback is code/fixture only.

<!-- tene:section:work-packages -->
## Work packages

Expand discovery/exclusion; add non-Go DTO; refine layers; add mature/polyglot fixtures; add matrix assertions; run full regression.

<!-- tene:section:dependencies -->
## Dependencies

Existing reference web, Go AST analyzer and graph command.

<!-- tene:section:verification -->
## Verification

Focused matrix, full check/race/vet, Playwright and validators.

<!-- tene:section:risks -->
## Risks

Path heuristics can overstate confidence; fallback confidence stays 0.4 and unknowns are mandatory.

<!-- tene:section:yagni -->
## Yagni

No language servers or external queues.
