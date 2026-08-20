---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wrs5bseqwkpkg84mdwg
phase: plan
status: draft
revision: 472
intent_ids: []
generated_at: 2026-08-20T01:25:52Z
generated_by: tene-workflow
---

# plan — Final Requirements Traceability Audit

<!-- tene:section:purpose -->
## Purpose

Execute a global requirements-to-evidence closure audit.

<!-- tene:section:scope -->
## Scope

Trace manifest/auditor, host probes, doc reconciliation, complete regression and final inactive check.

<!-- tene:section:layers -->
## Layers

Audit interface, coverage logic, state scans and CI.

<!-- tene:section:six-questions -->
## Six questions

Auditor is defined in scripts, invoked by Makefile/final command, reads manifest/project paths, and returns JSON coverage or exit 1. Probe is defined in projectconfig, called by doctor, and returns capability facts without mutation.

<!-- tene:section:traceability -->
## Traceability

Single blocking AC spans every required identifier family.

<!-- tene:section:decisions -->
## Decisions

Regular audit permits its own active Sprint; `--final` additionally requires no active Sprint after archive.

<!-- tene:section:freeform -->
## Freeform

Any missing locator or state debt becomes a blocker and is fixed before QA.

<!-- tene:section:work-packages -->
## Work packages

Build manifest/auditor; implement capability probes; reconcile docs; audit residual MUST items; run full suite; archive and run final audit.

<!-- tene:section:dependencies -->
## Dependencies

All prior archived Sprints and current official Codex manual.

<!-- tene:section:verification -->
## Verification

Audit normal/final, make/race/vet/e2e, routing/security/reference/release, plugin/skills, evidence and doctor.

<!-- tene:section:risks -->
## Risks

Self-referential active Sprint; solve with two-stage normal then post-archive final check.

<!-- tene:section:yagni -->
## Yagni

No new optional remote service or UI.
