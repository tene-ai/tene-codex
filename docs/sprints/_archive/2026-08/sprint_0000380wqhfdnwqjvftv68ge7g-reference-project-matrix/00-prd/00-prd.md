---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wqhfdnwqjvftv68ge7g
phase: prd
status: draft
revision: 411
intent_ids: []
generated_at: 2026-08-20T01:15:02Z
generated_by: tene-workflow
---

# prd — Reference Project Portability Matrix

<!-- tene:section:purpose -->
## Purpose

Prove AC-PRODUCT-07 portability across three repository archetypes.

<!-- tene:section:scope -->
## Scope

Static reference fixtures, polyglot filesystem fallback, layer rules, vault exclusion and matrix gate.

<!-- tene:section:layers -->
## Layers

Interface, Business Logic, Persistence and Infrastructure are represented and machine checked.

<!-- tene:section:six-questions -->
## Six questions

`codeintel.Analyze/sourcePaths/classify` are defined in `internal/codeintel`, called by graph commands/tests, accept root/path context, and return providers/files/components/edges/diagnostics with structured inputs, outputs, effects and unknowns.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380wqhjfnbjh3xxzpkhg1c`; AC `ac_0000380wqhjfq5y4p9jcbajkpm`; AC-PRODUCT-03/07, WP-08 and reference matrix plan.

<!-- tene:section:decisions -->
## Decisions

Unsupported languages produce filesystem components with explicit unknown Six Questions, never invented semantic edges. `.tene` is unconditionally excluded.

<!-- tene:section:freeform -->
## Freeform

Provider degradation is a supported outcome when honestly represented.

<!-- tene:section:problem -->
## Problem

Only the greenfield Go web journey existed; non-Go files were invisible and portable architecture claims were unproven.

<!-- tene:section:actors -->
## Actors

Plugin user, design/loop-check skill, code graph provider and evaluator.

<!-- tene:section:journeys -->
## Journeys

Analyze greenfield, mature and polyglot roots; compare layers/components; retain Go AST detail; surface non-Go uncertainty; verify vault exclusion.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

Three fixtures pass; mature project covers all four layers; polyglot has four source files with degraded unknowns; every unsupported component answers inputs/outputs/effects as unknown and has five missing semantic categories; `.tene` count is zero.

<!-- tene:section:non-goals -->
## Non goals

Implementing every language AST/LSP, running external services, or auto-indexing CodeGraph.
