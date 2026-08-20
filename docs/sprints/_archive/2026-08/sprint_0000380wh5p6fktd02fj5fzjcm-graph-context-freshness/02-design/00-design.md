---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wh5p6fktd02fj5fzjcm
phase: design
status: draft
revision: 197
intent_ids: []
generated_at: 2026-08-20T00:19:23Z
generated_by: tene-workflow
---

# design — Graph Context Freshness

<!-- tene:section:purpose -->
## Purpose

Specify implementation contracts for deterministic graph impact and fresh, bounded context packs.

<!-- tene:section:scope -->
## Scope

Pure algorithms live in `internal/tracecontext`; `internal/app` handles flags, state loading, errors, and optional file output.

<!-- tene:section:layers -->
## Layers

- Interface: command parsing and response envelopes.
- Business Logic: graph adjacency traversal, invariant queries, candidate ranking, budget enforcement, freshness diff.
- Persistence: atomic-style explicit pack output through a temporary sibling and rename.
- Infrastructure: SHA-256 and provider discovery only.

<!-- tene:section:six-questions -->
## Six questions

`Impact(graph,start,maxDepth,callDepth)` is defined in the new package, called by graph CLI, accepts immutable graph and limits, returns paths/ACs/truncation. `ValidateGraph(project)` returns findings. `BuildContextPack(root,project,sprint,options,capabilities)` returns a pack. `ValidateContextPack(root,project,pack)` returns freshness diagnostics. Only the app-level save helper changes files.

<!-- tene:section:traceability -->
## Traceability

Impact covers AC-GC-01; validation covers AC-GC-02; pack selection covers AC-GC-03; provenance/freshness covers AC-GC-04; test matrix covers AC-GC-05.

<!-- tene:section:decisions -->
## Decisions

- Traverse outgoing edges by default because graph edge semantics encode realization/verification direction; include path edges so callers can judge direction.
- The default depth is 8 and calls/imports are separately capped at 4 hops.
- Candidate groups: policy/safety, intent, blocking AC, blockers, phase artifacts/tasks, graph slice, recent decisions/history.
- Mandatory overflow is a typed error; optional items are excluded with group/reason/count.

<!-- tene:section:freeform -->
## Freeform

The pack schema is versioned independently as `1.0.0` and records a deterministic content hash for consumer-side caching.

<!-- tene:section:components -->
## Components

- `internal/tracecontext/graph.go`: `ImpactResult`, `GraphPath`, `ValidateGraph`.
- `internal/tracecontext/context.go`: `ContextPack`, `ContextItem`, `SourceRef`, `BuildOptions`, `BuildContextPack`, `ValidateContextPack`.
- `internal/app/app.go`: command adapters and stable command errors.
- tests: synthetic graph unit tests and temp-repository CLI journeys.

<!-- tene:section:interfaces -->
## Interfaces

`graph impact <node-id> [--depth N] [--call-depth N]` and `context build [--phase PHASE] [--budget BYTES] [--output PATH]`; `context validate --input PATH`. JSON result includes explicit units and no hidden mutation.

<!-- tene:section:data -->
## Data

Pack fields: schema/id/phase/state revision/budget/budget unit/used/objective/items/capabilities/provenance/excluded summary/content hash. Each item has kind/ref/priority/mandatory/content/source revision/locator/content hash. Provenance has locator/revision/SHA-256/status.

<!-- tene:section:state-transitions -->
## State transitions

These read commands do not change workflow revision. Saving a pack writes only the requested derived file. Validation moves no state; consumers must rebuild before mutation when `fresh=false`.

<!-- tene:section:failures -->
## Failures

Unknown node, invalid depth/budget/phase, mandatory overflow, outside-root output/input, malformed pack, missing provenance, changed hash, and changed state revision receive stable diagnostic codes.

<!-- tene:section:security -->
## Security

Never read `.tene/**`; reject provenance/output path escape; hash bytes without including secret-vault files; pack content uses workflow state and authored Sprint documents only.

<!-- tene:section:tests -->
## Tests

Cycles, call-depth cutoff, multiple AC paths, dangling edges, orphan task/evidence, deterministic sorting, budget exclusion, mandatory overflow, phase filtering, save/load, state drift, file drift, path escape, race, and existing full suite.
