---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wyd2xm0eey9n6e31px0
phase: prd
status: complete
revision: 719
intent_ids: [intent_0000380wyd7d99vkcqmwbmz158]
generated_at: 2026-08-20T02:15:00Z
generated_by: tene-workflow
---

# prd — Document Sync and CLI Contract Completion

<!-- tene:section:purpose -->
## Purpose

Make the documented CLI and document synchronization contracts executable and consistent.

<!-- tene:section:scope -->
## Scope

Managed Markdown synchronization, global flags, request retry deduplication, evidence deduplication, QA aliases and waiver approval lifecycle.

<!-- tene:section:layers -->
## Layers

CLI interface, command contract logic, canonical request/waiver state, and atomic filesystem/event infrastructure.

<!-- tene:section:six-questions -->
## Six questions

Sync/request/waiver components are defined in document/domain/app/state, invoked by public commands, consume canonical state and user documents, and return stable previews, events or lifecycle objects.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380wyd7d99vkcqmwbmz158`; AC `ac_0000380wyd7d8zhc4mwqd86wkm`.

<!-- tene:section:decisions -->
## Decisions

Sync defaults to preview; `--apply` is explicit. User prose is never rewritten. Reused request IDs with different commands conflict.

<!-- tene:section:freeform -->
## Freeform

Legacy aliases remain compatible while documented names become authoritative.

<!-- tene:section:problem -->
## Problem

Several commands existed only in design, global flags were rejected, and request IDs were parsed without affecting execution.

<!-- tene:section:actors -->
## Actors

CLI user, Codex skill, workflow controller and human waiver approver.

<!-- tene:section:journeys -->
## Journeys

Preview/apply document sync; retry a mutation safely; request then approve/expire a waiver; run QA through documented aliases.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

`ac_0000380wyd7d8zhc4mwqd86wkm`: golden manual-edit preservation, idempotent sync/retry/evidence, conflicting retry rejection and waiver inactive-before-approval.

<!-- tene:section:non-goals -->
## Non goals

No remote transaction coordinator or hosted document editor.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `719`
- Sprint: `sprint_0000380wyd2xm0eey9n6e31px0`
- Intents: `intent_0000380wyd7d99vkcqmwbmz158`
- Tasks: `task_0000380wz9ygjx9rs8dkgkrtf0`, `task_0000380wza05c2zmfcjch8s92m`

<!-- tene:generated:traceability:end -->
