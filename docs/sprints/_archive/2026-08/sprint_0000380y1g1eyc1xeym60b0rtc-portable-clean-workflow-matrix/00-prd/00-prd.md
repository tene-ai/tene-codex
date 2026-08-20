---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380y1g1eyc1xeym60b0rtc
phase: prd
status: complete
revision: 1005
intent_ids: [intent_0000380y1xbkhac1273cdc3284]
generated_at: 2026-08-20T07:21:39Z
generated_by: tene-workflow
---

# prd — Portable Clean Workflow Matrix

<!-- tene:section:purpose -->
## Purpose

Close AC-PRODUCT-07 with black-box proof that a staged plugin completes the entire governed lifecycle in clean Go, Next.js and Python repositories.

<!-- tene:section:scope -->
## Scope

Staged package, isolated repository fixtures, every workflow phase, document sync/validation, automatic loop, seven-variant QA, report/archive, evidence verify and doctor.

<!-- tene:section:layers -->
## Layers

Interface: public wrapper CLI. Business: identical phase/gate semantics. Persistence: per-project journal/docs/evidence/archive. Infrastructure: packaged binary, clean temp repositories and Git baselines.

<!-- tene:section:six-questions -->
## Six questions

Names: portable-workflow-smoke runner, STACKS, run_stack and release-smoke caller. Defined in scripts; called by release and semantic gates; input staged CLI/clean roots; output three archived verified workflow summaries.

<!-- tene:section:traceability -->
## Traceability

AC-PRODUCT-07, WP-14, FR-02/04/06/08/09 and blocking AC ac_0000380y1xbkgdqhqk45j21axg.

<!-- tene:section:decisions -->
## Decisions

Use the staged package wrapper as a black box. Each CLI call is a new process, also proving persisted resume across sessions. Light profile removes human approval boundaries but retains all logical phases and quality gates.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:problem -->
## Problem

Prior reference tests analyzed five project shapes and exercised web/worker journeys, but never completed PRD→archive independently in three clean stacks.

<!-- tene:section:actors -->
## Actors

Plugin adopter, Codex agent and release maintainer.

<!-- tene:section:journeys -->
## Journeys

Package → create clean stack repo → init/spec/task/design/do/loop → seven observations/evaluate → report/archive → doctor/evidence/status.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

Three named stacks each finish at revision 32 with 33 events, seven evidence records, one archive manifest, no active Sprint and healthy doctor.

<!-- tene:section:non-goals -->
## Non goals

Framework runtime correctness, deployment, exhaustive language coverage or bypassing approvals in non-light profiles.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `1005`
- Sprint: `sprint_0000380y1g1eyc1xeym60b0rtc`
- Intents: `intent_0000380y1xbkhac1273cdc3284`
- Tasks: `task_0000380y21dtymdms0ys2mcqyw`, `task_0000380y21f450p879vr1qgpv4`

<!-- tene:generated:traceability:end -->
