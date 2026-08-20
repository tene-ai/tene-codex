---
schema_version: 1.0.0
document_type: analysis
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wx6xt073n63as8wkzyw
phase: loop-check
status: passed
revision: 560
intent_ids: [intent_0000380wx710enjk9h81emvgf8]
generated_at: 2026-08-20T02:04:35Z
generated_by: tene-workflow
---

# loop-check — Automatic Bidirectional Loop Analysis

<!-- tene:section:purpose -->
## Purpose

Run the new analyzer against its own Sprint and record convergence.

<!-- tene:section:scope -->
## Scope

Specification coverage, task traceability, changed artifacts and executable design contracts.

<!-- tene:section:layers -->
## Layers

CLI, analyzer rules, gap reconciliation state and git/filesystem providers were jointly verified.

<!-- tene:section:six-questions -->
## Six questions

`Analyze` and `Reconcile` are defined in loopcheck, imported by app, called by `loop check`, consume state/docs/files, and return candidates plus persisted gap lifecycle changes.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380wx710enjk9h81emvgf8`, AC `ac_0000380wx710e3a7bw0g1kd6dc`, tasks `task_0000380wxj3hep0x4eaqm3bkdm` and `task_0000380wxj4kryss8dx8k5wdrm`.

<!-- tene:section:decisions -->
## Decisions

No manual gaps were seeded for the self-analysis.

<!-- tene:section:freeform -->
## Freeform

Text not represented by an ID, artifact link or executable contract remains narrative rather than machine proof.

<!-- tene:section:baseline -->
## Baseline

Old `loop check` only returned `OpenGapIDs`.

<!-- tene:section:changed-artifacts -->
## Changed artifacts

`internal/loopcheck/*`, domain gap metadata and app task/loop commands; all production changes are task-linked.

<!-- tene:section:gap-matrix -->
## Gap matrix

Six seeded classes detected in unit fixtures (100%, above 90% target). Actual Sprint analysis detected 0 after repair. Duplicate reconciliation creates 0 additional gaps; disappearance resolves; recurrence reopens the same ID.

<!-- tene:section:iterations -->
## Iterations

Iteration 1 extracted the analyzer and reconciliation from CLI orchestration and added convergence tests.

<!-- tene:section:regression -->
## Regression

Focused analyzer/app suites pass; full QA follows.
