---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wqhfdnwqjvftv68ge7g
phase: qa
status: draft
revision: 411
intent_ids: []
generated_at: 2026-08-20T01:15:02Z
generated_by: tene-workflow
---

# qa — Reference Project Portability Matrix

<!-- tene:section:purpose -->
## Purpose

Gate repository portability and uncertainty honesty.

<!-- tene:section:scope -->
## Scope

Three fixtures, graph fallback and full product regression.

<!-- tene:section:layers -->
## Layers

L1 classification, L2 Report contract, L3 filesystem exclusion, L4 greenfield Playwright, L5 degraded provider, L6 deterministic/race, L7 full checks.

<!-- tene:section:six-questions -->
## Six questions

Matrix invokes Analyze on each root and asserts files/layers/providers/unknown shapes without relying on implementation prose.

<!-- tene:section:traceability -->
## Traceability

One blocking AC receives happy/error/recovery cases and shared matrix evidence.

<!-- tene:section:decisions -->
## Decisions

Greenfield UX journey remains Playwright; mature/polyglot validate architecture portability without artificial runnable stacks.

<!-- tene:section:freeform -->
## Freeform

No external network or credential.

<!-- tene:section:environment -->
## Environment

Local static fixtures and existing greenfield server.

<!-- tene:section:capabilities -->
## Capabilities

Go AST, filesystem fallback, Playwright, validators and doctor.

<!-- tene:section:charters -->
## Charters

Greenfield end-to-end; mature legacy four-layer; polyglot degradation; vault exclusion; outside-root rejection.

<!-- tene:section:ux-data-flow -->
## Ux data flow

Repository sources → provider discovery → per-file analysis → layers/Six Questions → graph/context consumers.

<!-- tene:section:evidence -->
## Evidence

`04-qa/evidence/reference-matrix.json` and reproducible matrix/Playwright tests.

<!-- tene:section:verdict -->
## Verdict

PASS. Three repository types, mature 4/4 layers, polyglot four files/five unknown categories, vault files zero, check/race and Playwright 3/3 passed.
