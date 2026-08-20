---
schema_version: 1.0.0
document_type: analysis
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wh5p6fktd02fj5fzjcm
phase: loop-check
status: draft
revision: 197
intent_ids: []
generated_at: 2026-08-20T00:19:23Z
generated_by: tene-workflow
---

# loop-check — Graph Context Freshness

<!-- tene:section:purpose -->
## Purpose

Compare the implementation against the Sprint PRD, plan, design, and four blocking ACs before QA.

<!-- tene:section:scope -->
## Scope

Graph impact, invariant validation, context selection/budgeting, provenance freshness, CLI wiring, tests, and user-facing command references.

<!-- tene:section:layers -->
## Layers

Interface commands are thin adapters; `tracecontext` owns business rules; explicit pack output is derived persistence; provider discovery is infrastructure.

<!-- tene:section:six-questions -->
## Six questions

`Impact` and `ValidateGraph` are defined in `internal/tracecontext/graph.go`; `BuildContextPack` and `ValidateContextPack` are in `context.go`; `internal/app/app.go` imports and calls them. Inputs are domain project/graph, root, phase and budget; outputs are immutable path/finding/pack/freshness DTOs. Only `writeDerivedJSON` mutates an explicit repository-relative output.

<!-- tene:section:traceability -->
## Traceability

AC-GC-01: impact unit+CLI tests. AC-GC-02: dangling/orphan unit test and live graph validate. AC-GC-03: phase/budget/overflow tests. AC-GC-04: file and state drift tests. AC-GC-05: full checks below.

<!-- tene:section:decisions -->
## Decisions

Unresolved static call targets are explicit low-confidence `ExternalSymbol` nodes. This retains graph integrity without claiming resolution.

<!-- tene:section:freeform -->
## Freeform

The repository has no `.codegraph/`; per policy no index was created and Go AST/filesystem capability was reported.

<!-- tene:section:baseline -->
## Baseline

Before this Sprint, `graph trace` expanded a whole connected component, `impact` did not exist, context ignored phase/budget, and no reusable pack freshness check existed.

<!-- tene:section:changed-artifacts -->
## Changed artifacts

New `internal/tracecontext/{graph,context}.go` and tests; updated app integration tests, CLI references, README examples, and `tene-status` guidance.

<!-- tene:section:gap-matrix -->
## Gap matrix

All four ACs match implementation. The first live validation found dangling `call:` destinations and a gap→waiver source with no Gap node; both builders were repaired and graph validation now returns only expected pre-QA unverified warnings.

<!-- tene:section:iterations -->
## Iterations

1. Implemented initial traversal/context contracts.
2. Strengthened validation exposed missing derived call nodes; added explicit placeholders.
3. Waiver edge exposed absent Gap nodes; added Gap materialization.
4. Symbol locators were normalized to file hashes for clean provenance.

<!-- tene:section:regression -->
## Regression

Focused unit and integration tests pass. Full race/vet/make/plugin/evidence checks are scheduled in QA.
