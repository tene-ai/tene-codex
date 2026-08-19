---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v060y01n452j1rk3xr4
phase: plan
status: draft
revision: 1
intent_ids: []
generated_at: 2026-08-19T17:11:18Z
generated_by: tene-workflow
---

# plan — MVP workflow engine and Codex plugin

<!-- tene:section:purpose -->
## Purpose

Implement a tested vertical slice in dependency order: contracts and store, workflow/documents, intent/graph/QA, plugin integration, then packaging and dogfood verification.

<!-- tene:section:scope -->
## Scope

The work covers WP-01 through WP-14 at pre-alpha depth, with the deterministic local core and installable plugin surface complete enough for real iterative testing.

<!-- tene:section:layers -->
## Layers

Interface work follows core command contracts; business rules depend on domain types; persistence depends on atomic filesystem primitives; infrastructure packages the validated result.

<!-- tene:section:six-questions -->
## Six questions

Each work package records its public names and paths in the final report. Inputs and outputs are stabilized first through Go structs, JSON envelopes, hook event JSON, and Markdown section markers.

<!-- tene:section:traceability -->
## Traceability

Five blocking AC IDs are assigned to tasks. No implementation task is allowed to reach design without an AC link.

<!-- tene:section:decisions -->
## Decisions

Use standard-library Go without runtime dependencies. Use Python standard library for hooks. Validator-only PyYAML remains an external development tool. Preserve hook optionality and CLI correctness.

<!-- tene:section:freeform -->
## Freeform

The existing RND, PRD, plan, and design documents remain authoritative background and are not regenerated.

<!-- tene:section:work-packages -->
## Work packages

1. Domain, IDs, journal, projections, lock, atomic write, master plan and recovery.
2. Fixed lifecycle, authored document gates, intent/task/graph/context and loop gaps.
3. Evidence registry, QA planner/evaluator, secret adapter and full Sprint test.
4. Plugin manifest, nine skills, hooks, custom agents and validation.
5. CI, cross-platform package, README usage and dogfood report.

<!-- tene:section:dependencies -->
## Dependencies

Domain precedes state; state precedes workflow; workflow precedes QA; all core contracts precede skills/hooks and packaging. External QA tools remain adapters rather than core dependencies.

<!-- tene:section:verification -->
## Verification

Run `go test -race ./...`, `go vet ./...`, Python hook unit tests, JSON checks, official plugin/skill validators, complete-Sprint integration, secret denial tests, and four-target cross-build/package smoke tests.

<!-- tene:section:risks -->
## Risks

Pre-alpha schema evolution, Codex hook trust/version differences, false QA evidence, journal corruption, and accidental secret output. Mitigate with capability boundaries, revision locks, hash verification, strict blocker semantics, and fail-closed secret commands.

<!-- tene:section:yagni -->
## Yagni

Defer remote MCP, App Server control plane, GUI, cloud sync, universal AST/data-flow provider, production deployment and automatic plugin portal submission.
