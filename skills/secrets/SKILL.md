---
name: secrets
description: Run tests or development commands that require credentials without exposing secret values to Codex or workflow artifacts.
---

# tene Secrets

Explain safe setup and sanitized failures in the user's current conversation language, following the workflow language contract.

Read [secret safety](../../references/security.md) before acting. Use `tene-workflow secret check`, metadata-only `secret list`, and `secret run <environment> -- <command>`.

Never read `.tene/**`, call `tene get`, export plaintext, dump a child environment, wrap the command in a shell, or ask the user to paste a secret. Ask the user to run `tene set` directly when a key is missing. Preserve only the sanitized command, environment alias, exit code, and redaction-safe evidence.

This skill does not authorize production access or destructive external actions. Obtain the same approval those operations normally require.
