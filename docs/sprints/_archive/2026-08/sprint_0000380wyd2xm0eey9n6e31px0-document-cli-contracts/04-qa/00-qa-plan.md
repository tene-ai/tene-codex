---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wyd2xm0eey9n6e31px0
phase: qa
status: complete
revision: 719
intent_ids: [intent_0000380wyd7d99vkcqmwbmz158]
generated_at: 2026-08-20T02:15:00Z
generated_by: tene-workflow
---

# qa — Document Sync and CLI Contract Completion

<!-- tene:section:purpose -->
## Purpose

Prove document preservation and exact CLI state/ retry contracts.

<!-- tene:section:scope -->
## Scope

Seven variants, golden sync, flags/aliases, request and evidence dedup, waiver approval and regression.

<!-- tene:section:layers -->
## Layers

L1/L2/L3/L6/L7 required; L4/L5 explicitly not applicable for a local non-UI CLI contract.

<!-- tene:section:six-questions -->
## Six questions

Cases invoke public commands and native tests, persist structured evidence, and feed deterministic QA evaluation.

<!-- tene:section:traceability -->
## Traceability

Final run `run_0000380x0446gnmd4c3gbf6xhg`; AC `ac_0000380wyd7d8zhc4mwqd86wkm`. The earlier run was superseded after report review added atomic crash-window request markers.

<!-- tene:section:decisions -->
## Decisions

Preview must leave bytes unchanged; repeated apply and request must be idempotent.

<!-- tene:section:freeform -->
## Freeform

Legacy compatibility paths are tested without becoming the recommended workflow.

<!-- tene:section:environment -->
## Environment

Local repository fingerprint with spec hash `3122cb13...` at QA plan revision 648.

<!-- tene:section:capabilities -->
## Capabilities

Go native tests, CLI integration, filesystem hash comparison and event/project state inspection.

<!-- tene:section:charters -->
## Charters

Happy, alternate, empty, validation, permission, failure and recovery cover sync/retry/approval behavior.

<!-- tene:section:ux-data-flow -->
## Ux data flow

CLI input → validation/controller → document or canonical state mutation → stable result/error and retry response.

<!-- tene:section:evidence -->
## Evidence

Every case has actual `go-test` output and structured before/after contract observations.

<!-- tene:section:verdict -->
## Verdict

Passed: authored sentinel preserved, sync converges, retry deduplicates/conflicts correctly, a committed-but-incomplete retry never re-executes, waiver needs approval, aliases and flags work.


<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `719`
- Sprint: `sprint_0000380wyd2xm0eey9n6e31px0`
- Intents: `intent_0000380wyd7d99vkcqmwbmz158`
- Tasks: `task_0000380wz9ygjx9rs8dkgkrtf0`, `task_0000380wza05c2zmfcjch8s92m`

<!-- tene:generated:traceability:end -->
