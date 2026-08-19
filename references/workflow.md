# Workflow Contract

Use the plugin wrapper as `${PLUGIN_ROOT}/scripts/tene-workflow`; in a source checkout, `scripts/tene-workflow` is equivalent.

The lifecycle is fixed:

`draft → prd → plan → design → do ↔ loop-check → qa → report → archived`

Before any phase work:

1. Run `status --json`.
2. Run `context build --json` when a Sprint is active.
3. Use the active Sprint and current revision; do not create a second Sprint to hide unfinished work.
4. Perform the phase work and validate its artifact.
5. Run `phase transition <next> --dry-run --json` before committing the transition.

Never edit `.tene-workflow/*.json` or `events.ndjson` directly. Markdown under `docs/sprints/` is user-editable, but validate it through the CLI.

## Understanding contract

Every analysis, design, QA result, and report must cover or explicitly mark N/A for:

- Interface / Entry Point
- Business Logic / Processing Rules
- Persistence / Data
- Infrastructure / Runtime

For important changed components answer: name, definition file, import/reference sites, call/use sites, input shape, output or mutation.

## Delegation

Delegate only bounded, independent work when subagents are available. Good roles are product-intent discovery, read-only code exploration, builder, QA executor, and independent evaluator. Give each subagent the phase context and required output, not the full conversation. Subagents return evidence; the parent commits canonical state through the CLI.

