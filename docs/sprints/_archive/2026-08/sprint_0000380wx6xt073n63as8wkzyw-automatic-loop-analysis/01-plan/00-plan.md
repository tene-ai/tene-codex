---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wx6xt073n63as8wkzyw
phase: plan
status: approved
revision: 560
intent_ids: [intent_0000380wx710enjk9h81emvgf8]
generated_at: 2026-08-20T02:04:35Z
generated_by: tene-workflow
---

# plan — Automatic Bidirectional Loop Analysis

<!-- tene:section:purpose -->
## Purpose

Implement automatic comparison and convergent gap synchronization.

<!-- tene:section:scope -->
## Scope

Analyzer package, task artifact command, loop integration, contract parser and mutation tests.

<!-- tene:section:layers -->
## Layers

Interface `task artifact`/`loop check`; Business Logic analyzer; Persistence gap ledger; Infrastructure git and filesystem.

<!-- tene:section:six-questions -->
## Six questions

Analyzer definitions are imported by app runtime, called in loop-check, accept canonical state plus repository files, and return stable candidate gaps and reconciliation counts.

<!-- tene:section:traceability -->
## Traceability

AC `ac_0000380wx710e3a7bw0g1kd6dc`; tasks `task_0000380wxj3hep0x4eaqm3bkdm` and `task_0000380wxj4kryss8dx8k5wdrm`.

<!-- tene:section:decisions -->
## Decisions

Gap fingerprints are content-derived; analyzer-owned gaps may auto-resolve but manually authored gaps are never altered.

<!-- tene:section:freeform -->
## Freeform

Changed implementation files require an explicit task artifact owner.

<!-- tene:section:work-packages -->
## Work packages

Analyzer rules; canonical reconciliation; CLI artifact links; mutation and convergence tests.

<!-- tene:section:dependencies -->
## Dependencies

Domain fields → analyzer → app integration → tests and skills.

<!-- tene:section:verification -->
## Verification

Six mutation classes, clean convergence, duplicate-run stability, full regression and an actual Sprint self-analysis.

<!-- tene:section:risks -->
## Risks

Text contracts can overclaim; they are explicit high-confidence assertions whose paths and symbols are still verified.

<!-- tene:section:yagni -->
## Yagni

No new graph database or automatic CodeGraph indexing.
