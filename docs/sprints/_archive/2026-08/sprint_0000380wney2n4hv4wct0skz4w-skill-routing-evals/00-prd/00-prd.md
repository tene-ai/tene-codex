---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wney2n4hv4wct0skz4w
phase: prd
status: draft
revision: 340
intent_ids: []
generated_at: 2026-08-20T00:56:52Z
generated_by: tene-workflow
---

# prd — Skill Routing and Eval Completion

<!-- tene:section:purpose -->
## Purpose

Close WP-12 with explicit and deterministic implicit routing for all nine tene skills.

<!-- tene:section:scope -->
## Scope

Public `route` command, shared routing engine, thresholds, hard negatives, phase conflicts, multi-intent orchestration, corpus runner, CI gate and CLI documentation. No automatic phase mutation.

<!-- tene:section:layers -->
## Layers

Interface: CLI and skill discovery. Business: scoring and ambiguity policy. Persistence: versioned corpus, read-only state. Infrastructure: `make check` gate.

<!-- tene:section:six-questions -->
## Six questions

`router.Route`, `Decision`, `Candidate`, and `tene-routing-eval` are defined in `internal/router` and `cmd/tene-routing-eval`, imported by `internal/app` and the Makefile path. CLI/evals call them with UTF-8 text plus active phase. They return selected/proposed/none with scores, reasons and action constraints; project state is unchanged.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380wp1g79k7a136ay9wx00`; ACs `ac_0000380wp1g78fnzykfkpq2f2r`, `ac_0000380wp235a561m2wv2004vw`, `ac_0000380wp24y32a4qxf3tt687w`; FR-10; WP-12.

<!-- tene:section:decisions -->
## Decisions

Explicit `$tene-*` wins. A high-confidence single match selects, ambiguity proposes, hard negatives decline, and multi-intent requests hand off to `tene-sprint`.

<!-- tene:section:freeform -->
## Freeform

Future model-routing telemetry may extend but cannot weaken deterministic safety guarantees.

<!-- tene:section:problem -->
## Problem

Skill descriptions alone do not prove bilingual routing quality, collision behavior, or wrong-phase safety.

<!-- tene:section:actors -->
## Actors

Codex user, skill dispatcher, plugin maintainer and release engineer.

<!-- tene:section:journeys -->
## Journeys

Explicit invocation selects the named skill. Natural language is scored with current phase. Clear winners return a preflight contract; uncertain/multi-intent requests propose without mutation; unrelated requests return none. CI expands and evaluates the corpus.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

Each skill has 5 explicit, 20 positive, 20 adjacent-negative, 10 multi-intent, and ≥10 wrong-phase cases where applicable. Explicit succeeds 100%; precision/recall/multi ≥90%; wrong-phase/unnecessary ≤10%. CLI emits auditable JSON and evaluation failure fails `make check`.

<!-- tene:section:non-goals -->
## Non goals

Replacing native Codex discovery, silently chaining phases, inferring authorization, or persisting raw prompts.
