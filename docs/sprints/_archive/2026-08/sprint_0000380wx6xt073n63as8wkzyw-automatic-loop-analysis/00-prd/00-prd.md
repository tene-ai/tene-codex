---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wx6xt073n63as8wkzyw
phase: prd
status: confirmed
revision: 560
intent_ids: [intent_0000380wx710enjk9h81emvgf8]
generated_at: 2026-08-20T02:04:35Z
generated_by: tene-workflow
---

# prd — Automatic Bidirectional Loop Analysis

<!-- tene:section:purpose -->
## Purpose

Turn Loop Check into an independent bidirectional specification/code analyzer.

<!-- tene:section:scope -->
## Scope

Intent/AC/task/document coverage, changed-artifact ownership, executable design contracts and stable automatic gap reconciliation.

<!-- tene:section:layers -->
## Layers

CLI analysis, deterministic comparison rules, canonical gap persistence and git/filesystem providers.

<!-- tene:section:six-questions -->
## Six questions

`loopcheck.Analyze` is defined in the loopcheck package, called by `loop check`, accepts project/Sprint/documents/changed files, and returns fingerprinted candidates that update canonical gaps.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380wx710enjk9h81emvgf8`; AC `ac_0000380wx710e3a7bw0g1kd6dc`.

<!-- tene:section:decisions -->
## Decisions

Only explicit document IDs, task artifact links and machine-readable design contracts are treated as high-confidence proof.

<!-- tene:section:freeform -->
## Freeform

Unknown semantic behavior remains visible rather than inferred as success.

<!-- tene:section:problem -->
## Problem

The previous command only counted gaps already entered by a person or agent.

<!-- tene:section:actors -->
## Actors

Builder, evaluator, reviewer and workflow controller.

<!-- tene:section:journeys -->
## Journeys

Run analysis, inspect generated gaps, repair documents/code/links, rerun, and converge without manual gap closure.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

`ac_0000380wx710e3a7bw0g1kd6dc`: six seeded drift classes are detected; duplicate runs do not duplicate gaps; repaired conditions auto-resolve.

<!-- tene:section:non-goals -->
## Non goals

No claim of universal semantic equivalence without explicit contracts or a semantic provider.
