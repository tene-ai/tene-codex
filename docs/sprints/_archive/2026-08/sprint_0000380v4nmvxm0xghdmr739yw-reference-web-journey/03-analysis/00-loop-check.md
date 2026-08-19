---
schema_version: 1.0.0
document_type: analysis
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v4nmvxm0xghdmr739yw
phase: loop-check
status: draft
revision: 163
intent_ids: []
generated_at: 2026-08-19T17:50:32Z
generated_by: tene-workflow
---

# loop-check — Reference Web Journey Harness

<!-- tene:section:purpose -->
## Purpose

Compare the live reference journey with AC-PRODUCT-05 and Sprint design.

<!-- tene:section:scope -->
## Scope

Browser UI, network boundary, validation/recovery rules, JSON persistence, Playwright output and CI.

<!-- tene:section:layers -->
## Layers

Interface page/API; Business validation/fail-once/retry; Persistence observed file-backed state; Infrastructure server/Chromium/CI.

<!-- tene:section:six-questions -->
## Six questions

Handlers are defined in `testdata/reference-web/main.go` and called by HTTP/JS/Playwright; tests/config are in `tests/e2e` and root config; inputs are `{name,fail_once}`, outputs are status/item/error, mutations are bounded JSON writes and test artifacts.

<!-- tene:section:traceability -->
## Traceability

AC1 success test; AC2 validation and recovery tests; AC3 pinned package/config/CI and tene evidence import.

<!-- tene:section:decisions -->
## Decisions

Real Chromium is mandatory; no mocked DOM or screenshot-only proof.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:baseline -->
## Baseline

Only adapter and observation contract fixtures existed.

<!-- tene:section:changed-artifacts -->
## Changed artifacts

Reference Go server, package lock/config, three Playwright journeys, CI, gitignore and Sprint docs/evidence.

<!-- tene:section:gap-matrix -->
## Gap matrix

All AC implementation paths match; no blocking gap.

<!-- tene:section:iterations -->
## Iterations

Initial implementation passed all three live journeys without repair.

<!-- tene:section:regression -->
## Regression

Playwright 3/3 passed; full Go/plugin gates run in QA.
