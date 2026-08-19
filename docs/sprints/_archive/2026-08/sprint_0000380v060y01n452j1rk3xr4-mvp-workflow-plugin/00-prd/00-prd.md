---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v060y01n452j1rk3xr4
phase: prd
status: draft
revision: 1
intent_ids: []
generated_at: 2026-08-19T17:11:18Z
generated_by: tene-workflow
---

# prd — MVP workflow engine and Codex plugin

<!-- tene:section:purpose -->
## Purpose

Deliver a runnable pre-alpha vertical slice of tene-codex that turns Codex work into durable, gated sprints and demonstrates its own workflow on this repository.

<!-- tene:section:scope -->
## Scope

Implement the Go `tene-workflow` core, fixed phase guards, document scaffolding, intent/task/graph/context state, loop gaps, evidence-based QA, tene secret execution, nine Codex skills, lifecycle hooks, subagent profiles, tests, CI, and release packaging.

<!-- tene:section:layers -->
## Layers

- Interface: CLI commands, Codex skills, hook events, and generated agent profiles.
- Business Logic: workflow guards, graph construction, QA evaluator, routing instructions.
- Persistence: `.tene-workflow` journal/projections, Sprint Markdown, evidence hashes.
- Infrastructure: Go builds, GitHub Actions, plugin bundle and cross-platform binaries.

<!-- tene:section:six-questions -->
## Six questions

The primary names are `app.Run`, `state.Store`, `workflow.CanTransition`, `workflow.EvaluateQAGate`, `document.Scaffold`, and the plugin skills/hooks. They are defined under `cmd/`, `internal/`, `skills/`, and `hooks/`; imported by the CLI composition root and called by user commands or Codex lifecycle events. Inputs are typed command/domain values and JSON hook events. Outputs are JSON envelopes, atomic state changes, Markdown artifacts, evidence verdicts, or hook decisions.

<!-- tene:section:traceability -->
## Traceability

This Sprint implements FR-01 through FR-11 and creates executable evidence for AC-PRODUCT-01 through AC-PRODUCT-08 from `docs/00-prd`.

<!-- tene:section:decisions -->
## Decisions

Use a standalone Go binary, append-only hash-chained events plus projections, strict document content gates, default plugin hook discovery, project-scoped custom agents, and Apache-2.0. No remote MCP or hosted service is required for the MVP.

<!-- tene:section:freeform -->
## Freeform

The implementation remains pre-alpha; browser-driven product QA adapters and full schema migration automation remain follow-up work.

<!-- tene:section:problem -->
## Problem

Codex can implement and test code but does not by itself provide a project-portable source of truth that preserves product intent, forces a full Sprint cycle, and refuses completion without UX/data-flow evidence.

<!-- tene:section:actors -->
## Actors

- Developer using Codex in any Git repository.
- Product owner confirming intent and policy.
- Builder subagent implementing a bounded task.
- Independent evaluator reviewing evidence.

<!-- tene:section:journeys -->
## Journeys

The developer initializes tene, creates a Sprint, confirms intent and observable acceptance criteria, plans and designs work, implements it, closes spec/code gaps, executes applicable QA, reviews the report, and archives a durable record. Secret-dependent commands cross only the `tene run` boundary.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

1. An invalid phase skip is rejected with a stable code and remediation.
2. A new session can recover the active Sprint, phase, revision, intent, task and gap context.
3. The plugin exposes nine valid skills, trusted-review hooks and five specialized project agent profiles.
4. Every blocking QA case requires matching, hash-valid, redaction-safe evidence before report/archive.
5. Direct secret retrieval, vault reads, exports and environment dumps are blocked; approved commands use `tene run` without a shell.
6. Go tests, race tests, vet, hook tests, validators and four-platform packaging pass.

<!-- tene:section:non-goals -->
## Non goals

No cloud control plane, GUI dashboard, remote MCP server, universal language parser, production deployment automation, or plaintext secret API is included in this Sprint.
