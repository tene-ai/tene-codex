---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v2b2q0668t7tajkj3m0
phase: prd
status: draft
revision: 88
intent_ids: []
generated_at: 2026-08-19T17:30:10Z
generated_by: tene-workflow
---

# prd — Code Intelligence and Intent QA Adapters

<!-- tene:section:purpose -->
## Purpose

Close the two highest-priority carry items from the MVP Sprint: semantic code understanding and capability-aware, intent-driven QA. The workflow must explain changed code through the four Understanding Layers and Six Questions, then plan and collect evidence across executable and externally observed UX/data journeys.

<!-- tene:section:scope -->
## Scope

- Discover CodeGraph without creating its index, and fall back to language-native/static providers.
- Analyze Go declarations, imports, calls, signatures and likely side effects; retain explicit unknowns.
- Classify files/symbols into Interface, Business Logic, Persistence and Infrastructure with evidence/confidence.
- Discover native test, Playwright and browser-observation capabilities.
- Extend QA charters with steps, observers, forbidden outcomes, required layers and risk.
- Import sanitized external browser/API/data observations as hash-bound evidence.

<!-- tene:section:layers -->
## Layers

Interface: CLI commands and QA observation import. Business Logic: provider selection, classification, Six Questions and capability-aware charter compilation. Persistence: graph projection, QA runs and immutable evidence artifacts. Infrastructure: git, Go parser/toolchain, CodeGraph and Playwright/browser capability probes.

<!-- tene:section:six-questions -->
## Six questions

The implementation must report every analyzed declaration's name, definition locator, import/reference sites, call/use sites, input shape and output/side effects. Unsupported facts are `unknown` with provider diagnostics, never fabricated.

<!-- tene:section:traceability -->
## Traceability

This Sprint implements WP-08 and WP-10 from `docs/01-plan`, and resolves the corresponding carry items in the MVP report. Every task links to the confirmed intent and both blocking acceptance criteria.

<!-- tene:section:decisions -->
## Decisions

CodeGraph is preferred only when `.codegraph/` already exists. External MCP/browser execution remains controlled by Codex; the CLI owns a stable observation import and evidence contract. Commands execute via argv without a shell.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:problem -->
## Problem

The current graph represents authored workflow objects but cannot explain code. QA creates generic variants but cannot describe observer checkpoints or identify available execution tools. This prevents trustworthy whole-flow review and encourages unsupported LLM claims.

<!-- tene:section:actors -->
## Actors

Developers, reviewers, Codex builder/evaluator agents, native test runners, Playwright, and interactive browser/Chrome tools.

<!-- tene:section:journeys -->
## Journeys

1. A developer builds the graph; tene-workflow probes providers, analyzes code, and materializes layer/Six Questions output.
2. QA planning compiles confirmed ACs into journey charters and records capability applicability for seven QA layers.
3. Codex or a browser tool performs a journey, writes a sanitized observation artifact, imports it, binds it to a case, and the independent gate evaluates hash-valid evidence.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

- AC1: `graph build` combines workflow state with deterministic Go code analysis, while `graph understand` returns all four layer categories and all Six Questions with provenance/confidence/unknowns.
- AC2: `qa capabilities`, richer `qa plan`, and observation import expose native/Playwright/browser capabilities and produce evidence-bound cases without shell execution or secret leakage.
- AC3: unit, integration, race, vet, schema/plugin/skill validation and a dogfood Sprint gate all pass.

<!-- tene:section:non-goals -->
## Non goals

Creating CodeGraph indexes, embedding an MCP browser inside the Go binary, implementing a full LSP client, or supporting every programming language in this Sprint.
