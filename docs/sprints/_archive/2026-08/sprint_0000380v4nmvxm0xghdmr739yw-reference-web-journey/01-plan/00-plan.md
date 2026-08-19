---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v4nmvxm0xghdmr739yw
phase: plan
status: draft
revision: 163
intent_ids: []
generated_at: 2026-08-19T17:50:32Z
generated_by: tene-workflow
---

# plan — Reference Web Journey Harness

<!-- tene:section:purpose -->
## Purpose

Deliver the first live reference journey and evidence pipeline.

<!-- tene:section:scope -->
## Scope

Reference app, Playwright tests/reporter, scripts, CI, docs and dogfood evidence.

<!-- tene:section:layers -->
## Layers

UI/API, processing rules, durable test file and browser runtime are all exercised.

<!-- tene:section:six-questions -->
## Six questions

Every handler/test/reporter records definition, references, calls, request/response and side effects.

<!-- tene:section:traceability -->
## Traceability

App+happy test→AC1; negative/retry tests→AC2; reporter/CI/import→AC3.

<!-- tene:section:decisions -->
## Decisions

Pin Playwright dependency and run one Chromium worker deterministically.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:work-packages -->
## Work packages

1. Build reference app. 2. Add Playwright journeys and data observer endpoint. 3. Emit observation artifact. 4. Integrate CI/docs. 5. QA/archive/push.

<!-- tene:section:dependencies -->
## Dependencies

Server before browser tests; tests before reporter evidence; evidence before gate.

<!-- tene:section:verification -->
## Verification

Go tests, Playwright suite, observation schema/import, race/vet/plugin/skill/doctor.

<!-- tene:section:risks -->
## Risks

Browser availability: install pinned Chromium in CI/local. Port collisions: Playwright webServer owns a fixed test port with reuse disabled in CI.

<!-- tene:section:yagni -->
## Yagni

No React/database container/auth service; behavior matters, not framework complexity.
