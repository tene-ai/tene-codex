---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wr84t6p1hb89mf2wbkc
phase: plan
status: draft
revision: 437
intent_ids: []
generated_at: 2026-08-20T01:21:13Z
generated_by: tene-workflow
---

# plan — Release and Marketplace Completion

<!-- tene:section:purpose -->
## Purpose

Close WP-13 and the inherited migration debt.

<!-- tene:section:scope -->
## Scope

Scripts, legal/listing docs, catalog, CI, compatibility tests and release QA.

<!-- tene:section:layers -->
## Layers

Marketplace interface, package guards, preserved project state, CI supply chain.

<!-- tene:section:six-questions -->
## Six questions

Package/smoke/SBOM scripts are invoked by make/CI, take explicit paths/version, and output staged files or fail. Migration planner consumes schema header and returns supported plan without mutation until apply.

<!-- tene:section:traceability -->
## Traceability

Two tasks map to two ACs; inherited gap closes only after release evidence exists.

<!-- tene:section:decisions -->
## Decisions

Test current platform in smoke while tag builds four; stage overwrite is refused; uninstall removes bundle only.

<!-- tene:section:freeform -->
## Freeform

Portal submission itself remains an external organization action after code readiness.

<!-- tene:section:work-packages -->
## Work packages

Package/SBOM/checksum; smoke marketplace lifecycle; CI provenance/docs; migration window and cross-sprint debt resolution; release QA.

<!-- tene:section:dependencies -->
## Dependencies

Plugin validator, nine skills, routing eval, reference matrix, GitHub Actions and official OpenAI portal.

<!-- tene:section:verification -->
## Verification

Make/race/e2e/validators, package inventory, checksum tamper, SBOM JSON, migration tests, doctor.

<!-- tene:section:risks -->
## Risks

Portal requirements change; refresh official manual per release. Shell cleanup is confined to mktemp root.

<!-- tene:section:yagni -->
## Yagni

No MCP, telemetry backend or auto-submission.
