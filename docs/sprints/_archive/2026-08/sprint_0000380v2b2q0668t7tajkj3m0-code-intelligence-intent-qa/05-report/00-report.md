---
schema_version: 1.0.0
document_type: report
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v2b2q0668t7tajkj3m0
phase: report
status: draft
revision: 88
intent_ids: []
generated_at: 2026-08-19T17:30:10Z
generated_by: tene-workflow
---

# report — Code Intelligence and Intent QA Adapters

<!-- tene:section:purpose -->
## Purpose

Review delivery of the first-class code-intelligence and intent-QA adapter layer and preserve continuity for the next implementation Sprint.

<!-- tene:section:scope -->
## Scope

WP-08/WP-10 provider contracts, Go AST understanding, CodeGraph query boundary, capability-aware QA charters, allowlisted test execution, external observation import and traceability hardening.

<!-- tene:section:layers -->
## Layers

Interface: `graph providers|understand` and `qa capabilities|execute|observe`. Business Logic: codeintel and qaadapter. Persistence: derived graph nodes/edges, extended QA cases and hash-bound evidence. Infrastructure: safe argv probes/execution for git, Go, CodeGraph, npm/npx plus Codex-controlled external observers.

<!-- tene:section:six-questions -->
## Six questions

Definitions are in `internal/codeintel/codeintel.go`, `internal/qaadapter/qaadapter.go`, domain types and app routing. App imports and calls the packages from graph/QA commands. Inputs are repository roots, bounded paths/query strings, known adapter names and schema-valid observation files. Outputs are capabilities, Reports, ExecutionResults and Observations; graph/evidence/QA case state changes occur through hash-chained mutations.

<!-- tene:section:traceability -->
## Traceability

AC1 is satisfied by analyzer outputs/tests and the derived graph. AC2 is satisfied by capability discovery, native executions and observation import. AC3 is satisfied by full regression, race, vet, JSON, hook, plugin, skill, evidence and doctor gates.

<!-- tene:section:decisions -->
## Decisions

Semantic facts unsupported by AST remain explicit unknowns. Existing CodeGraph is opt-in via a bounded query and is never indexed automatically. External browser tools remain Codex capabilities behind a provider-neutral observation artifact.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:previous-sprints -->
## Previous sprints

The MVP Sprint created workflow graph/QA/evidence primitives. The archive-integrity Sprint made their locators durable. This Sprint extends those exact primitives with real code-derived nodes and runtime/observer evidence, without replacing the state or gate model.

<!-- tene:section:changed-files -->
## Changed files

Added `internal/codeintel/*`, `internal/qaadapter/*`, `schemas/qa-observation.schema.json`, and Sprint evidence/docs. Extended `internal/domain/types.go`, `internal/app/app.go`, app tests, QA skill/reference, CLI reference, README and README-KR. State/journal files record the dogfood lifecycle.

<!-- tene:section:intent-fulfillment -->
## Intent fulfillment

Codex can now answer the four layers and Six Questions from deterministic code facts with confidence/unknowns. QA begins from ACs, exposes seven-layer dispositions, runs only known adapters, and imports full UX/API/data checkpoint assertions into the same evidence gate.

<!-- tene:section:qa-verdict -->
## Qa verdict

Passed: all nine cases, six evidence artifacts, race/vet/unit/integration/hook/schema/plugin/skill validators, evidence verification and workflow QA gate.

<!-- tene:section:deferred-work -->
## Deferred work

Remaining master-plan work is tracked, not silently dropped: richer polyglot/LSP providers; stable machine parsing of CodeGraph results; first reference web application with live Playwright/Chrome and DB/API observers; waiver/migration commands; Windows packaging; marketplace portal assets/signing. These require distinct fixtures or release policy and are scheduled after the current contracts are stable.

<!-- tene:section:next-sprint -->
## Next sprint

Implement schema migration/waiver lifecycle and a reference web-journey harness that exercises Playwright plus API/data observers end-to-end. Then complete release reproducibility, signed artifacts and public marketplace readiness.

<!-- tene:generated:summary:start -->
### Generated Sprint Summary

- Sprint: `sprint_0000380v2b2q0668t7tajkj3m0`
- Previous sprints: ``
- Intent IDs: `intent_0000380v2edshbndyv7s262qx4`
- Tasks: 3
- QA verdict: `passed`
- Open gaps: 0
- State revision: 124

<!-- tene:generated:summary:end -->
