# Workflow Contract

Use the plugin wrapper as `${PLUGIN_ROOT}/scripts/tene-workflow`; in a source checkout, `scripts/tene-workflow` is equivalent.

The lifecycle is fixed:

`draft → prd → plan → design → do ↔ loop-check → qa → report → archived`

If report review discovers a late code or harness change, `report → loop-check` reopens verification; it must pass Loop Check and a new QA run before returning to report.

Before any phase work:

1. Run `status --json`.
2. Run `context build --json` when a Sprint is active.
3. Use the active Sprint and current revision; do not create a second Sprint to hide unfinished work.
4. Perform the phase work and validate its artifact.
5. Run `phase transition <next> --dry-run --json` before committing the transition.

## Conversation language

Write authored PRD, plan, design, loop-check, QA, and report content in the language currently used by the user. An explicit user language request overrides inference. Preserve machine-readable frontmatter, section markers, IDs, paths, commands, API names, and code symbols exactly; technical terms may remain in their established form. Carry the language choice across phase handoffs and context compaction. Do not translate archived history unless the user explicitly requests it.

Never edit `.tene-workflow/*.json` or `events.ndjson` directly. Markdown under `docs/sprints/` is user-editable, but validate it through the CLI.

Blocking non-security gaps may cross a phase gate only through a separately requested and approved active waiver with explicit reason, requester, approver, scope, and future expiry. Security and evidence-integrity gaps are never waivable. Revocation or expiry restores the block immediately.

Before upgrading persisted state, use `migrate dry-run`. Apply accepts only declared source versions, preserves the original projection under `.tene-workflow/backups`, and records a migration event. `doctor --repair` verifies the journal and project, then repairs only derived `active.json` and `master-plan.json`; it does not rewrite history.

## Understanding contract

Every analysis, design, QA result, and report must cover or explicitly mark N/A for:

- Interface / Entry Point
- Business Logic / Processing Rules
- Persistence / Data
- Infrastructure / Runtime

For important changed components answer: name, definition file, import/reference sites, call/use sites, input shape, output or mutation.

## Delegation

Delegate only bounded, independent work when subagents are available. Good roles are product-intent discovery, read-only code exploration, builder, QA executor, and independent evaluator. Give each subagent the phase context and required output, not the full conversation. Subagents return evidence; the parent commits canonical state through the CLI.
