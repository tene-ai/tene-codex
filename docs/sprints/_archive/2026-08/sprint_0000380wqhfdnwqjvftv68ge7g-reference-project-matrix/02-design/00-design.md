---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wqhfdnwqjvftv68ge7g
phase: design
status: draft
revision: 411
intent_ids: []
generated_at: 2026-08-20T01:15:02Z
generated_by: tene-workflow
---

# design — Reference Project Portability Matrix

<!-- tene:section:purpose -->
## Purpose

Design uncertainty-honest polyglot analysis and repeatable portability validation.

<!-- tene:section:scope -->
## Scope

Extension allowlist, skip rules, generic component shape, layer classifier and fixtures.

<!-- tene:section:layers -->
## Layers

Interface/frontend/API; Business service/worker; Persistence store/queue; Infrastructure infra/deploy.

<!-- tene:section:six-questions -->
## Six questions

`sourceExtension` selects known source files; `Analyze` emits Go symbols or generic file components; `classify` attaches layer/reason/confidence. Inputs are root/requested paths; outputs are a deterministic Report.

<!-- tene:section:traceability -->
## Traceability

Direct mapping to the Sprint AC, WP-08 and design 11 Scenario D.

<!-- tene:section:decisions -->
## Decisions

Non-Go generic component has empty imports/references/calls, explicit unknown input/output/effect markers, five Unknown entries, provider filesystem, confidence 0.4 and diagnostic.

<!-- tene:section:freeform -->
## Freeform

Future LSP providers replace unknown fields incrementally without changing the Report contract.

<!-- tene:section:components -->
## Components

Discovery, Go AST provider, filesystem fallback, classifier and reference matrix evaluator.

<!-- tene:section:interfaces -->
## Interfaces

Existing `Analyze(context, root, requested, changed) (Report,error)` remains backward compatible.

<!-- tene:section:data -->
## Data

Report preserves provider, confidence, locator, layers, Six Questions and diagnostics for every discovered source.

<!-- tene:section:state-transitions -->
## State transitions

None in analyzer; graph build may persist the resulting normalized nodes later.

<!-- tene:section:failures -->
## Failures

Outside-root path rejects; unreadable/oversized files diagnose; missing semantic provider degrades explicitly; secret directories skip.

<!-- tene:section:security -->
## Security

`.tene`, `.git`, workflow state, dependencies and build outputs are excluded.

<!-- tene:section:tests -->
## Tests

Three archetypes, four-layer mature assertion, polyglot unknowns, vault exclusion and full regression.
