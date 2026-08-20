---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wpmvv19m87t6xa4xydr
phase: design
status: draft
revision: 377
intent_ids: []
generated_at: 2026-08-20T01:07:13Z
generated_by: tene-workflow
---

# design — Secret Boundary and Adversarial QA

<!-- tene:section:purpose -->
## Purpose

Define a value-unrepresentable parent boundary and fail-closed evidence gate.

<!-- tene:section:scope -->
## Scope

Secret runner/types, list projection, detector, hooks and app evidence integration.

<!-- tene:section:layers -->
## Layers

Interface names/run/hooks; business validation/detection; persistence evidence verifier; infrastructure subprocess isolation.

<!-- tene:section:six-questions -->
## Six questions

`runPath/listNamesPath` live in `internal/secret`, are called by exported adapter methods/tests, receive executable/env/argv or metadata request, and return safe DTOs/errors. `DetectLeak` is imported by app `looksSecret`; hook functions consume JSON events and emit deny/context without echoing matches.

<!-- tene:section:traceability -->
## Traceability

Every failure table row in design 08 maps to a test and the two Sprint ACs.

<!-- tene:section:decisions -->
## Decisions

Execute argv directly. Result has no value field. List parses only `environment` and `secrets[].name`. Leak match returns empty stdout/stderr with `leak_detected` and `quarantined`. Evidence uses the same detector.

<!-- tene:section:freeform -->
## Freeform

Post-tool hook is defense in depth; core/evidence enforcement remains authoritative when hooks are unavailable.

<!-- tene:section:components -->
## Components

Runner policy regexes, names-only DTO, detector, app evidence adapter, PreToolUse denial and PostToolUse warning.

<!-- tene:section:interfaces -->
## Interfaces

`Run(context, env, []string) (Result,error)`, `ListNames(context,env) (Names,error)`, `DetectLeak([]byte) bool`.

<!-- tene:section:data -->
## Data

Names contains environment/count/name array only. Result contains exit code, sanitized command/output and boolean leak/quarantine status; no secret or preview fields.

<!-- tene:section:state-transitions -->
## State transitions

None for secret execution. Evidence registration is rejected before mutation on leak; QA stays failed without valid evidence.

<!-- tene:section:failures -->
## Failures

Missing CLI, invalid request, forbidden command, sensitive argv, permission/list failure, child failure, leak quarantine and corrupt metadata JSON all return stable SEC codes.

<!-- tene:section:security -->
## Security

No shell, `.tene`, `.env`, env dump, get/export, preview, literal credential, or raw leaked output crosses the boundary.

<!-- tene:section:tests -->
## Tests

Fake executable attack matrix plus hook and full QA/evidence-integrity regression.
