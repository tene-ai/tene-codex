---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v4nmvxm0xghdmr739yw
phase: prd
status: draft
revision: 163
intent_ids: []
generated_at: 2026-08-19T17:50:32Z
generated_by: tene-workflow
---

# prd — Reference Web Journey Harness

<!-- tene:section:purpose -->
## Purpose

Prove intent-driven QA on a real browser → API → persistent-data journey rather than only contract fixtures.

<!-- tene:section:scope -->
## Scope

Self-contained reference web app, Playwright happy/validation/failure/retry/accessibility journeys, API and file-state observers, reporter that emits tene QA observations, CI integration and documentation.

<!-- tene:section:layers -->
## Layers

Interface HTML/form/API; Business validation/failure/retry rules; Persistence JSON record file; Infrastructure HTTP server/Playwright/CI.

<!-- tene:section:six-questions -->
## Six questions

Document server handlers, browser tests, observer/reporter definitions, callers, request shapes, response/state mutations and artifacts.

<!-- tene:section:traceability -->
## Traceability

Implements AC-PRODUCT-05 and the remaining WP-10 reference journey gate.

<!-- tene:section:decisions -->
## Decisions

Use a dependency-light Go reference server and Playwright Chromium. Persistence is an isolated temp JSON file so tests observe real writes without external credentials.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:problem -->
## Problem

The adapter contracts are tested, but no live UI→API→data chain proves UX intent, error feedback, retry and persisted outcome together.

<!-- tene:section:actors -->
## Actors

End user, browser, API handler, file persistence observer and independent QA evaluator.

<!-- tene:section:journeys -->
## Journeys

Valid submission persists and confirms; empty input is blocked; simulated downstream failure shows retry feedback without persistence; retry succeeds exactly once; keyboard/label semantics remain accessible.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

- AC1: Playwright drives a real page and proves UI, network response and persisted file state for success.
- AC2: validation and downstream failure/retry paths prove forbidden writes and understandable recovery feedback.
- AC3: test output converts to schema-valid tene observation/evidence and runs in CI with full regression gates.

<!-- tene:section:non-goals -->
## Non goals

Production web framework, real database/cloud deployment, authentication provider or cross-browser matrix in this Sprint.
