---
schema_version: 1.0.0
document_type: analysis
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wpmvv19m87t6xa4xydr
phase: loop-check
status: draft
revision: 377
intent_ids: []
generated_at: 2026-08-20T01:07:13Z
generated_by: tene-workflow
---

# loop-check — Secret Boundary and Adversarial QA

<!-- tene:section:purpose -->
## Purpose

Verify the secret boundary and evidence gate against documented attacks.

<!-- tene:section:scope -->
## Scope

Adapter, hooks, shared detector, app evidence mutation and tests.

<!-- tene:section:layers -->
## Layers

Interface and hooks deny unsafe use; business quarantines leaks; persistence rejects poisoned evidence; CI executes attacks.

<!-- tene:section:six-questions -->
## Six questions

`Run/ListNames/DetectLeak` are defined in `internal/secret`, imported by app, called by secret/evidence commands, accept env/argv/bytes, return names-only/sanitized DTOs or stable errors with no leaked value.

<!-- tene:section:traceability -->
## Traceability

Both ACs and WP-11 map to direct tests.

<!-- tene:section:decisions -->
## Decisions

Unknown matching output is discarded, not partially redacted.

<!-- tene:section:freeform -->
## Freeform

Implementation review found and fixed preview forwarding and missing bare-canary detection.

<!-- tene:section:baseline -->
## Baseline

Preview fields and unlabeled canary output could cross the boundary; hook lacked `.env`, literal and post-output coverage.

<!-- tene:section:changed-artifacts -->
## Changed artifacts

`internal/secret/*`, `internal/app/app.go`, `internal/app/app_test.go`, `hooks/*`, `tests/hooks_test.py`, Sprint/state files.

<!-- tene:section:gap-matrix -->
## Gap matrix

All known gaps resolved: names-only, shell/env/literal denial, canary quarantine, sanitized errors, safe child status, hook pre/post protection, poisoned evidence pre-mutation rejection.

<!-- tene:section:iterations -->
## Iterations

One implementation pass plus adversarial review; no open blocker remains.

<!-- tene:section:regression -->
## Regression

`make check` passes including Go, routing and five hook tests.
