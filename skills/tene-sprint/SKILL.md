---
name: tene-sprint
description: Create, resume, advance, and archive spec-driven tene sprints when a coding request needs durable multi-phase workflow management.
---

# tene Sprint

Manage the project-level Sprint lifecycle. Read [workflow](../../references/workflow.md) and [CLI reference](../../references/cli.md).

Run `status --json` first. If the project is not initialized, explain the created paths and run `init` only when the user wants workflow state in this repository. Create a Sprint for meaningful feature, bug, or refactor work; do not create one for a general explanation.

Scaffold documents on creation. Advance one legal phase at a time after a dry-run guard passes. Never force a transition. At archive, ensure blocking QA passed, the report is validated, and unresolved work is explicitly deferred.

Use subagents only for bounded independent phase work; the parent owns state transitions.

