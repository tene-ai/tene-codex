---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v3xjbejnjk8w65p9q2w
phase: qa
status: draft
revision: 127
intent_ids: []
generated_at: 2026-08-19T17:43:57Z
generated_by: tene-workflow
---

# qa — Waiver Schema Migration and Recovery

<!-- tene:section:purpose -->
## Purpose

Verify exception safety, migration recoverability and derived-only repair.

<!-- tene:section:scope -->
## Scope

Happy/error/recovery variants for all three ACs plus regression and artifact integrity.

<!-- tene:section:layers -->
## Layers

L1 schema/vet; L2 unit/race; L3 state/CLI; L4 N/A no UI; L5 operator feedback; L6 expiry/forbidden/unsupported/corruption; L7 prior suite.

<!-- tene:section:six-questions -->
## Six questions

Inspect Waiver and state API definitions/callers, structured flag/project inputs, plan/backup/event/projection outputs and explicit mutations.

<!-- tene:section:traceability -->
## Traceability

Three blocking ACs compile to nine cases with criterion-linked evidence.

<!-- tene:section:decisions -->
## Decisions

Migration uses isolated legacy fixture; the live repository only runs status and repair, avoiding artificial downgrade.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:environment -->
## Environment

Local Go/macOS repository, schema 1.0.0, valid journal, no product UI.

<!-- tene:section:capabilities -->
## Capabilities

Go native adapter available; external observation import available; browser not applicable.

<!-- tene:section:charters -->
## Charters

Waiver active/expired/security/revoked; migration current/legacy/unsupported; repair missing-derived/corrupt-source/healthy repeat.

<!-- tene:section:ux-data-flow -->
## Ux data flow

CLI metadata → validation → state/workflow rule → backup/journal/projection → operator result and doctor verification.

<!-- tene:section:evidence -->
## Evidence

Three native adapter artifacts cover the waiver, migration and repair ACs. `make check`, race, vet, schema parsing, hook tests, evidence verification and doctor passed; legacy/unknown/corrupt cases run in isolated unit fixtures.

<!-- tene:section:verdict -->
## Verdict

Passed: all nine cases and all three blocking ACs have criterion-linked, redaction-safe evidence; `qa evaluate` returned no findings.
