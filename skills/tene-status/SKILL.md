---
name: tene-status
description: Resume a tene workflow by showing the active sprint, current phase, blockers, revision, and next legal action.
---

# tene Status

Run `status --json`. If no project exists, explain that status is unavailable and suggest `$tene-sprint` only when the user wants to initialize this repository. If a Sprint is active, run `context build --phase <current-phase> --json` and summarize objective, current phase, included and excluded context, open tasks and gaps, gate status, and the next legal transition. Before reusing a saved pack, run `context validate --input <path>`.

Do not mutate state, create a Sprint, or claim work is complete. If the state revision changes during inspection, rebuild the context. Read [CLI reference](../../references/cli.md).
