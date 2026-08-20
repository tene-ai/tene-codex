---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wh5p6fktd02fj5fzjcm
phase: prd
status: draft
revision: 197
intent_ids: []
generated_at: 2026-08-20T00:19:23Z
generated_by: tene-workflow
---

# prd — Graph Context Freshness

<!-- tene:section:purpose -->
## Purpose

Complete WP-06 and WP-07 so Codex can calculate requirement impact and build a bounded, provenance-aware context pack that is rejected when repository state changes.

<!-- tene:section:scope -->
## Scope

- Add directed `graph impact` traversal with configurable call depth and impacted acceptance criteria.
- Detect dangling graph edges and orphaned traceability artifacts without turning provider limitations into false facts.
- Build phase-aware context packs with a deterministic character budget, priority ordering, exclusions, capabilities, provenance revisions, and file hashes.
- Add `context validate` so a previously saved pack can be checked before mutation.
- Preserve the existing graph build/trace and context build CLI behavior where compatible.

<!-- tene:section:layers -->
## Layers

- Interface: `tene-workflow graph impact`, `graph validate`, `context build`, and `context validate`.
- Business Logic: directed traversal, graph invariants, phase selection, priority/budget allocation, freshness comparison.
- Persistence: optional JSON context-pack output and hashes of project state and active phase documents; no new source-of-truth state.
- Infrastructure: provider capability snapshots and portable filesystem hashing.

<!-- tene:section:six-questions -->
## Six questions

The new public names are `Impact`, `Validate`, `BuildContextPack`, and `ValidateContextPack`; they are defined in `internal/tracecontext`, called by `internal/app`, accept project/graph/root/policy inputs, and return deterministic DTOs without mutating project state. The CLI may write only an explicitly requested context output file.

<!-- tene:section:traceability -->
## Traceability

This Sprint realizes WP-06, WP-07, AC-PRODUCT-03, and the context freshness parts of the architecture design. Evidence must include unit tests for traversal, invariants, budget, phase filtering, and stale hashes plus CLI integration tests.

<!-- tene:section:decisions -->
## Decisions

- Budget is measured in UTF-8 bytes as a deterministic local proxy, exposed as `budget_unit`; model tokenization is not silently guessed.
- Safety, confirmed intent, blocking AC, and open blockers are mandatory. If they cannot fit, the command fails instead of dropping them.
- Context packs are derived artifacts; validation compares recorded revision and SHA-256 hashes with current sources.

<!-- tene:section:freeform -->
## Freeform

Future semantic providers may enrich paths, but this Sprint does not auto-index CodeGraph and does not claim unsupported data-flow certainty.

<!-- tene:section:problem -->
## Problem

Current `trace` returns an entire connected component, there is no impact command, graph validation misses dangling/orphan links, and `context build` ignores phase, budget, provider capability, exclusion, and freshness. This allows unrelated or stale context to influence implementation and QA.

<!-- tene:section:actors -->
## Actors

- Developer or reviewer asking what requirements a change can affect.
- Codex skill assembling only the context needed for the active phase.
- QA/evaluator checking that evidence was produced from current specification and code inputs.

<!-- tene:section:journeys -->
## Journeys

1. Build the graph, request impact for a symbol or AC, inspect explicit paths and impacted blocking criteria.
2. Build a context pack for a phase and budget, persist it, then validate it immediately.
3. Modify a provenance document and validate again; the pack is reported stale with the changed locator.
4. Validate malformed graph links and receive stable findings with remediation.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

- AC-GC-01: impact traversal returns bounded directed paths, impacted AC IDs, truncation status, and rejects unknown nodes.
- AC-GC-02: graph validation reports dangling edges and orphan tasks/evidence as stable findings while provider unknowns remain warnings.
- AC-GC-03: context build honors requested phase and byte budget, never removes safety/blocking items, deduplicates content, and explains exclusions.
- AC-GC-04: every included source has locator, revision, and SHA-256 where file-backed; validation detects state or file drift.
- AC-GC-05: unit, CLI integration, race, vet, plugin, and evidence checks pass without regressions.

<!-- tene:section:non-goals -->
## Non goals

- Automatic CodeGraph indexing, full cross-language LSP data flow, remote vector memory, or model-specific token counting.
