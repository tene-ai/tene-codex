---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wyd2xm0eey9n6e31px0
phase: design
status: complete
revision: 719
intent_ids: [intent_0000380wyd7d99vkcqmwbmz158]
generated_at: 2026-08-20T02:15:00Z
generated_by: tene-workflow
---

# design — Document Sync and CLI Contract Completion

<!-- tene:section:purpose -->
## Purpose

Specify managed document regions and exact command state machines.

<!-- tene:section:scope -->
## Scope

Preview/apply sync, retry cache, dedup keys, waiver lifecycle, aliases and global output flags.

<!-- tene:section:layers -->
## Layers

CLI dispatch → contract logic → project/event/document persistence → filesystem/tool runtime.

<!-- tene:section:six-questions -->
## Six questions

`document.Sync`, request cache helpers and waiver handlers are defined in document/app/domain, called by CLI, accept state and explicit args, and return hash diffs or persisted lifecycle results.

<!-- tene:section:traceability -->
## Traceability

AC `ac_0000380wyd7d8zhc4mwqd86wkm`; tasks `task_0000380wz9ygjx9rs8dkgkrtf0` and `task_0000380wza05c2zmfcjch8s92m`.

<!-- tene:section:decisions -->
## Decisions

Generated traceability has stable markers. Waivers remain inactive until separate approval. Quiet suppresses human success output but never JSON.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:contract path="internal/document/document.go" symbol="func Sync" -->
<!-- tene:contract path="internal/app/app.go" symbol="WaiverRequested" -->

<!-- tene:section:components -->
## Components

Document sync engine, request cache, evidence deduplicator, waiver controller and CLI aliases.

<!-- tene:section:interfaces -->
## Interfaces

`document|docs sync [--phase] [--apply]`; documented global flags; `waiver request|approve|expire`; `qa run` alias.

<!-- tene:section:data -->
## Data

Request results include command hash; waivers include requester/requested/approved/expired timestamps; sync returns before/after hashes.

<!-- tene:section:state-transitions -->
## State transitions

Waiver requested→active→expired/revoked. Request ID unused→completed; same hash returns cached; different hash conflicts.

<!-- tene:section:failures -->
## Failures

Invalid paths/expiry/approval, request reuse, stale revision and write errors fail closed with stable exits.

<!-- tene:section:security -->
## Security

No secret-bearing flags; managed regions only; approval identity and reason are retained.

<!-- tene:section:tests -->
## Tests

Golden preservation/idempotence, CLI aliases/quiet, request retry/conflict, evidence dedup and waiver inactive-before-approval.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `719`
- Sprint: `sprint_0000380wyd2xm0eey9n6e31px0`
- Intents: `intent_0000380wyd7d99vkcqmwbmz158`
- Tasks: `task_0000380wz9ygjx9rs8dkgkrtf0`, `task_0000380wza05c2zmfcjch8s92m`

<!-- tene:generated:traceability:end -->
