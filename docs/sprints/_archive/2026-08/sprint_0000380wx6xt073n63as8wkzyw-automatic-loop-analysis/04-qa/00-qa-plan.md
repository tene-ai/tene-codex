---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wx6xt073n63as8wkzyw
phase: qa
status: passed
revision: 560
intent_ids: [intent_0000380wx710enjk9h81emvgf8]
generated_at: 2026-08-20T02:04:35Z
generated_by: tene-workflow
---

# qa — Automatic Bidirectional Loop Analysis

<!-- tene:section:purpose -->
## Purpose

Verify detection accuracy, fail-closed behavior and gap convergence.

<!-- tene:section:scope -->
## Scope

Six mutation classes, seven QA variants, native tests, self-analysis and regression.

<!-- tene:section:layers -->
## Layers

L1/L2/L3/L6/L7 required; L4/L5 explicitly not applicable for a local non-UI analyzer.

<!-- tene:section:six-questions -->
## Six questions

Cases invoke native tests and structured observations; evidence binds the loop AC and produces the deterministic QA verdict.

<!-- tene:section:traceability -->
## Traceability

Run `run_0000380wxtrexn4d1f6wqcte4w`, AC `ac_0000380wx710e3a7bw0g1kd6dc`.

<!-- tene:section:decisions -->
## Decisions

Detection rate is measured from seeded independent fixtures rather than self-description.

<!-- tene:section:freeform -->
## Freeform

Manual gap entry was not used to obtain the loop pass.

<!-- tene:section:environment -->
## Environment

Local Go repository fingerprint and spec hash `baa51dbc...` at QA plan revision 582.

<!-- tene:section:capabilities -->
## Capabilities

Go test, git status, filesystem design-contract analyzer and canonical gap reconciler.

<!-- tene:section:charters -->
## Charters

Happy, alternate, empty, validation, permission, failure and recovery; each cites the six-class mutation matrix and convergence test.

<!-- tene:section:ux-data-flow -->
## Ux data flow

CLI → document/state/git analysis → fingerprinted candidates → canonical gap lifecycle → repair rerun result.

<!-- tene:section:evidence -->
## Evidence

Each case includes actual `go-test` output and structured before/after analyzer observations under `04-qa/evidence`.

<!-- tene:section:verdict -->
## Verdict

Passed: 6/6 mutations detected, clean fixture and real Sprint converge to zero, duplicate/resolution/reopen lifecycle passes.
