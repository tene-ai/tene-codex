---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wpmvv19m87t6xa4xydr
phase: qa
status: draft
revision: 377
intent_ids: []
generated_at: 2026-08-20T01:07:13Z
generated_by: tene-workflow
---

# qa — Secret Boundary and Adversarial QA

<!-- tene:section:purpose -->
## Purpose

Prove zero accepted plaintext and zero adversarial false-pass.

<!-- tene:section:scope -->
## Scope

Happy/error/recovery cases for both blocking ACs plus complete regression.

<!-- tene:section:layers -->
## Layers

L1 detector/runner; L2 CLI/hooks; L3 evidence state; L5 attack cases; L6 race/failure; L7 release checks. L4 N/A.

<!-- tene:section:six-questions -->
## Six questions

Tests call the same runner, detector, hook and evidence commands used in production with fake inputs and observe safe DTOs, exit codes, no mutations and no matched plaintext.

<!-- tene:section:traceability -->
## Traceability

QA plan generates six cases from two blocking ACs.

<!-- tene:section:decisions -->
## Decisions

Only sanitized aggregate results are evidence; raw fake canary output is intentionally never persisted.

<!-- tene:section:freeform -->
## Freeform

Real credentials and production access are forbidden in this QA.

<!-- tene:section:environment -->
## Environment

Local isolated temp directories and fake executable; no vault.

<!-- tene:section:capabilities -->
## Capabilities

Go tests/race/vet, Python hook tests, Playwright regression, plugin/skill validators and doctor.

<!-- tene:section:charters -->
## Charters

Names-only success; missing/permission/child failure; forbidden exfiltration; canary/token quarantine; hook input/output; evidence missing/tampered/leaked false-pass.

<!-- tene:section:ux-data-flow -->
## Ux data flow

User selects safe env/key names → tene child only → sanitized outcome → evidence hash/redaction → independent QA verdict.

<!-- tene:section:evidence -->
## Evidence

`04-qa/evidence/security-summary.json`, reproducible from checked-in fake-runner, hook and poisoned-evidence tests.

<!-- tene:section:verdict -->
## Verdict

PASS. Adapter attack cases 6, poisoned evidence patterns 3, hooks 5, accepted plaintext matches 0, previews returned 0. Full check/race/vet, Playwright 3/3, plugin and nine skill validators passed.
