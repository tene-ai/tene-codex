# Secret Safety

Secret values must not enter prompts, workflow state, documents, graph nodes, logs, reports, or evidence.

Allowed:

- `tene version`, `tene whoami`
- metadata-only `tene list --env <name> --json`
- `tene run --env <name> -- <command>`
- `tene-workflow secret check|list|run`

Forbidden in agent automation:

- reading `.tene/**`
- `tene get`
- plaintext `tene export`
- `env`, `printenv`, shell wrappers, or debug dumps inside a secret-injected child
- asking the user to paste a secret into chat

If the user pasted a value, do not repeat or store it. Recommend rotation and ask the user to set it directly through tene. If tene is unavailable or denies access, fail closed.

