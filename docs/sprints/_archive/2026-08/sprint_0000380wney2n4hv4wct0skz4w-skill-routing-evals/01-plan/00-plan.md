---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wney2n4hv4wct0skz4w
phase: plan
status: draft
revision: 340
intent_ids: []
generated_at: 2026-08-20T00:56:52Z
generated_by: tene-workflow
---

# plan — Skill Routing and Eval Completion

<!-- tene:section:purpose -->
## Purpose

Deliver routing vertically from product contract through CLI and release evaluation.

<!-- tene:section:scope -->
## Scope

Rules, CLI integration, corpus, runner, tests, documentation, validators and regression suite.

<!-- tene:section:layers -->
## Layers

Interface command; routing policy service; corpus fixture; CI gate.

<!-- tene:section:six-questions -->
## Six questions

`Route` is defined in `internal/router/router.go`, imported by `internal/app`, called by CLI and eval runner, receives text/active/phase, and returns a non-mutating `Decision`. The runner reads JSON and returns metrics plus exit status.

<!-- tene:section:traceability -->
## Traceability

All work packages map to the confirmed intent and three blocking criteria.

<!-- tene:section:decisions -->
## Decisions

Product and evaluator share the Go package; 5 stems × 4 suffixes create 20 positives; 20 adjacent negatives apply to every skill.

<!-- tene:section:freeform -->
## Freeform

Rollback removes the command and check target; no state migration exists.

<!-- tene:section:work-packages -->
## Work packages

1. Decision types, cue scoring and safety policy. 2. Read-only CLI and unit tests. 3. Bilingual corpus and threshold runner. 4. Release gate, validators and documentation.

<!-- tene:section:dependencies -->
## Dependencies

Domain phases, project loader, CLI envelope and nine skill packages.

<!-- tene:section:verification -->
## Verification

Unit, corpus metrics, `make check`, race, vet, plugin/skill validators, JSON smoke, evidence verify and doctor.

<!-- tene:section:risks -->
## Risks

Cue collision can over-route; use hard negatives, phase compatibility, margin, proposal-only ambiguity and regression metrics.

<!-- tene:section:yagni -->
## Yagni

No embedding service, telemetry backend, model fine-tune, or remote router.
