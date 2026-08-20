---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wpmvv19m87t6xa4xydr
phase: plan
status: draft
revision: 377
intent_ids: []
generated_at: 2026-08-20T01:07:13Z
generated_by: tene-workflow
---

# plan — Secret Boundary and Adversarial QA

<!-- tene:section:purpose -->
## Purpose

Implement secret and evidence defenses as executable attack cases.

<!-- tene:section:scope -->
## Scope

Runner hardening, hook hardening, shared detector, false-pass regression and release verification.

<!-- tene:section:layers -->
## Layers

CLI/hooks; security rules; immutable evidence; CI/fake executable.

<!-- tene:section:six-questions -->
## Six questions

Core names are defined in `internal/secret/runner.go`, imported by app, invoked by secret/evidence paths, accept arrays not shell strings, and return names-only or sanitized/quarantined results.

<!-- tene:section:traceability -->
## Traceability

Two work packages cover both blocking ACs and WP-11.

<!-- tene:section:decisions -->
## Decisions

Use executable fake tene fixtures, never a real vault or credential; error messages do not incorporate child stderr on metadata failure.

<!-- tene:section:freeform -->
## Freeform

Rollback is code-only with no state migration, but weak behavior must not be restored.

<!-- tene:section:work-packages -->
## Work packages

1. Harden adapter and add missing/preview/command/canary/permission/child tests. 2. Share leak detection with evidence and add pre/post hook cases. 3. Run adversarial QA and integrity tests.

<!-- tene:section:dependencies -->
## Dependencies

Existing tene CLI contract, evidence verifier, QA evaluator and hook runtime.

<!-- tene:section:verification -->
## Verification

Unit, Python hooks, app false-pass scenarios, race/vet, full check, validators, Playwright and doctor.

<!-- tene:section:risks -->
## Risks

False positives can block safe output; patterns remain scoped to credential/canary forms and quarantine is explicit.

<!-- tene:section:yagni -->
## Yagni

No vault access, decrypt API, remote scanner or telemetry.
