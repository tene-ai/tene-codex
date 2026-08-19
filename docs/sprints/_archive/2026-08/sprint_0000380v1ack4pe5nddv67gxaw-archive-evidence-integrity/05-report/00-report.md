---
schema_version: 1.0.0
document_type: report
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v1ack4pe5nddv67gxaw
phase: report
status: draft
revision: 59
intent_ids: []
generated_at: 2026-08-19T17:21:14Z
generated_by: tene-workflow
---

# report — Archive evidence integrity stabilization

<!-- tene:section:purpose -->
## Purpose

Report the archive integrity stabilization discovered while verifying the first MVP Sprint.

<!-- tene:section:scope -->
## Scope

Archive locator/manifest mutation, active-independent evidence verification, phase-aware doctor and regression tests.

<!-- tene:section:layers -->
## Layers

Interface behavior stays compatible; business verification is shared; persistence locators are consistent; infrastructure health checks work after archive.

<!-- tene:section:six-questions -->
## Six questions

`runtime.transition` in `internal/app/app.go` is called by archive, receives target/current state, and moves documents plus updates report/evidence/graph locators. `invalidEvidence` in the same file is called by evidence verify and doctor, receives root/project, reads and hashes artifacts, and returns sorted invalid IDs without mutation. `documentDue` is called by doctor and maps current/document phases to validation scope.

<!-- tene:section:traceability -->
## Traceability

The predecessor Sprint led directly to two confirmed ACs, two completed tasks, six passed QA cases and one new evidence artifact.

<!-- tene:section:decisions -->
## Decisions

Archive manifest is mandatory, historical verification needs no active Sprint, and doctor validates due documents only while validating all archived documents.

<!-- tene:section:freeform -->
## Freeform

This Sprint demonstrates cross-Sprint continuity and regression use of archived evidence.

<!-- tene:section:previous-sprints -->
## Previous sprints

Extends `sprint_0000380v060y01n452j1rk3xr4`, which implemented the MVP and exposed stale locator and future-document doctor risks during final verification.

<!-- tene:section:changed-files -->
## Changed files

`internal/app/app.go`, `internal/app/app_test.go`, `references/cli.md`, archive manifests and this Sprint document/evidence tree.

<!-- tene:section:intent-fulfillment -->
## Intent fulfillment

Reports, evidence URIs and graph locators move with the Sprint; archive emits a manifest; evidence verify and doctor work after the active pointer is cleared.

<!-- tene:section:qa-verdict -->
## Qa verdict

Passed: race/full tests, vet, Python hooks, six QA cases, two evidence hashes, graph validation and phase-aware doctor.

<!-- tene:section:deferred-work -->
## Deferred work

Cross-filesystem copy transactions, remote evidence retention and migration of pre-alpha external repositories remain deferred because no released schema requires them yet.

<!-- tene:section:next-sprint -->
## Next sprint

Proceed to code-intelligence providers and browser/data-flow QA adapters described in the MVP report.


<!-- tene:generated:summary:start -->
### Generated Sprint Summary

- Sprint: `sprint_0000380v1ack4pe5nddv67gxaw`
- Previous sprints: `sprint_0000380v060y01n452j1rk3xr4`
- Intent IDs: `intent_0000380v1bjszxx5fvyqj97dbc`
- Tasks: 2
- QA verdict: `passed`
- Open gaps: 0
- State revision: 85

<!-- tene:generated:summary:end -->
