---
name: sprint
description: Create, resume, advance, and archive spec-driven tene sprints when a coding request needs durable multi-phase workflow management.
---

# tene Sprint

Author all new workflow documents in the user's current conversation language, following the workflow language contract.

Manage the project-level Sprint lifecycle. Read [workflow](../../references/workflow.md) and [CLI reference](../../references/cli.md).

Run `status --json` first. If the project is not initialized, explain the created paths and run `init` only when the user wants workflow state in this repository. Create a Sprint for meaningful feature, bug, or refactor work; do not create one for a general explanation.

Scaffold documents on creation. Advance one legal phase at a time after a dry-run guard passes. Never force a transition. Respect profile approval boundaries: strict requires a scoped human approval for design→do and report→archive; standard requires one for report→archive. Request and approve the exact transition with a reason and expiry, then pass its ID to the transition. Approval never overrides a quality finding. At archive, ensure blocking QA passed, the report is validated, and unresolved work is explicitly deferred.

Use subagents only for bounded independent phase work; the parent owns state transitions.

For a deliberately accepted non-security gap, use the explicit `waiver` lifecycle with a named approver, bounded scope, reason, and future expiry; never rewrite a failed result as passed. Before state upgrades run `migrate dry-run`, and use `doctor --repair` only for derived projections after journal verification.
