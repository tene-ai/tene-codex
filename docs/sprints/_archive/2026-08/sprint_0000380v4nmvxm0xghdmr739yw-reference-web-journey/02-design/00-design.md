---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v4nmvxm0xghdmr739yw
phase: design
status: draft
revision: 163
intent_ids: []
generated_at: 2026-08-19T17:50:32Z
generated_by: tene-workflow
---

# design — Reference Web Journey Harness

<!-- tene:section:purpose -->
## Purpose

Specify a portable full-flow QA fixture.

<!-- tene:section:scope -->
## Scope

Go HTTP binary, HTML/JS client, JSON persistence, Playwright test/reporter and CLI evidence import.

<!-- tene:section:layers -->
## Layers

Interface `/` and `/api/items`; Business name validation and one-shot failure; Persistence append/read JSON; Infrastructure localhost server and Playwright.

<!-- tene:section:six-questions -->
## Six questions

Handlers accept JSON/HTTP and return status/result; store reads/writes arrays; tests call page/request and observe UI/network/state; reporter outputs qa-observation JSON.

<!-- tene:section:traceability -->
## Traceability

AC1 success chain, AC2 negative/recovery, AC3 artifact/CI.

<!-- tene:section:decisions -->
## Decisions

`X-Fail-Once: true` simulates downstream failure in test only. `/api/state` is a read-only observer endpoint enabled by the reference fixture.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:components -->
## Components

`testdata/reference-web/main.go`, embedded page, Playwright spec/config and observation reporter.

<!-- tene:section:interfaces -->
## Interfaces

GET `/`; POST `/api/items`; GET `/api/state`; npm `test:e2e`.

<!-- tene:section:data -->
## Data

Input `{name}`; item `{id,name}`; state `[{id,name}]`; observation checkpoints/assertions matching schema 1.0.0.

<!-- tene:section:state-transitions -->
## State transitions

empty→400/no write; fail-once→503/no write; retry→201/one write; page displays corresponding feedback.

<!-- tene:section:failures -->
## Failures

Malformed JSON 400, simulated dependency 503, file failure 500, browser assertion fail produces test failure and no passed observation.

<!-- tene:section:security -->
## Security

Loopback only, temp state, bounded body, HTML textContent, no secrets.

<!-- tene:section:tests -->
## Tests

Go handler/unit plus Playwright success/validation/recovery/accessibility and observation schema/import.
