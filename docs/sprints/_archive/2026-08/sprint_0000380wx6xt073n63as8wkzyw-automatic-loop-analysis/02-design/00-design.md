---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wx6xt073n63as8wkzyw
phase: design
status: approved
revision: 560
intent_ids: [intent_0000380wx710enjk9h81emvgf8]
generated_at: 2026-08-20T02:04:35Z
generated_by: tene-workflow
---

# design — Automatic Bidirectional Loop Analysis

<!-- tene:section:purpose -->
## Purpose

Define deterministic loop inputs, gap identity and convergence.

<!-- tene:section:scope -->
## Scope

Document graph checks, changed-file ownership and executable design markers.

<!-- tene:section:layers -->
## Layers

Interface commands call Business Logic without embedding rules; results persist as gaps; git/filesystem are bounded Infrastructure providers.

<!-- tene:section:six-questions -->
## Six questions

`Analyze`, `Candidate`, `Task.Artifacts` and `Gap.Fingerprint` are defined in loopcheck/domain, referenced by app, invoked by loop check, consume documents/state/files, and produce reconciled gap mutations.

<!-- tene:section:traceability -->
## Traceability

AC `ac_0000380wx710e3a7bw0g1kd6dc`; tasks `task_0000380wxj3hep0x4eaqm3bkdm` and `task_0000380wxj4kryss8dx8k5wdrm`.

<!-- tene:section:decisions -->
## Decisions

Machine contracts use stable comments and do not parse natural-language promises as executable truth.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:contract path="internal/loopcheck/analyzer.go" symbol="func Analyze" -->
<!-- tene:contract path="internal/app/app.go" symbol="TaskArtifactLinked" -->
<!-- tene:forbid path="internal/loopcheck/analyzer.go" contains="os.RemoveAll" -->

<!-- tene:section:components -->
## Components

Analyzer, candidate fingerprinting, task artifact linker and gap reconciler.

<!-- tene:section:interfaces -->
## Interfaces

`task artifact ID --path FILE`; `loop check` returns detected/created/resolved/reopened and open gaps.

<!-- tene:section:data -->
## Data

Analyzer gaps store fingerprint, detector and detection revision; task artifacts are repository-relative paths.

<!-- tene:section:state-transitions -->
## State transitions

Undetected → open; same fingerprint remains one gap; absent on rerun → resolved; reappearing condition → reopened.

<!-- tene:section:failures -->
## Failures

Git provider errors fail analysis; missing files and symbols become gaps; invalid artifact paths are rejected.

<!-- tene:section:security -->
## Security

Repository-relative paths only; `.tene` and workflow projections are excluded; analyzer never deletes files.

<!-- tene:section:tests -->
## Tests

Seed missing ID, missing task edge, untraced file, absent symbol, forbidden dependency and missing linked artifact; clean fixture returns zero.
