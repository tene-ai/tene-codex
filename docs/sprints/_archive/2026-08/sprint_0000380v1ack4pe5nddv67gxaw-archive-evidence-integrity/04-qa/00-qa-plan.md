---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v1ack4pe5nddv67gxaw
phase: qa
status: draft
revision: 59
intent_ids: []
generated_at: 2026-08-19T17:21:14Z
generated_by: tene-workflow
---

# qa — Archive evidence integrity stabilization

<!-- tene:section:purpose -->
## Purpose

Verify archive locator, manifest, evidence and phase-aware doctor behavior.

<!-- tene:section:scope -->
## Scope

Run focused/full Go tests, race, vet, Python hooks, real predecessor evidence verification and doctor.

<!-- tene:section:layers -->
## Layers

Interface commands, business validation, persistence relocation and infrastructure health checks are all exercised.

<!-- tene:section:six-questions -->
## Six questions

QA invokes `transition`, `evidence verify`, `doctor` and the archive integration fixture with repository or temporary project inputs and observes paths, hashes, manifest, findings and exit codes.

<!-- tene:section:traceability -->
## Traceability

Six cases cover happy/error/recovery variants for two blocking ACs.

<!-- tene:section:decisions -->
## Decisions

Use deterministic CLI/filesystem evidence; no browser layer applies.

<!-- tene:section:freeform -->
## Freeform

The predecessor's evidence is a historical regression input.

<!-- tene:section:environment -->
## Environment

Local macOS x86_64, Go 1.26.6 for Go 1.24 module, Python 3.14.6, no live secret.

<!-- tene:section:capabilities -->
## Capabilities

Go test/race/vet, Python unittest, filesystem archive and CLI doctor/evidence verification.

<!-- tene:section:charters -->
## Charters

Archive success asserts manifest and paths; error/recovery are covered by mutation rollback logic and integration. Verification success reads archived evidence; missing/tampered cases are covered by helper logic and doctor blocker behavior.

<!-- tene:section:ux-data-flow -->
## Ux data flow

Archive command → filesystem move/manifest → state locator mutation → no active Sprint → evidence verify/doctor → healthy result.

<!-- tene:section:evidence -->
## Evidence

`evidence/archive-integrity-verification.txt` records the final test and dogfood outcomes and is registered to both ACs.

<!-- tene:section:verdict -->
## Verdict

Passed. Six of six cases passed with criterion-linked evidence `evidence_0000380v1kct16ha5mk782hkc4`; all repository evidence verified, the independent deterministic evaluator returned no finding, and the rebuilt graph passed invariants at revision 84.
