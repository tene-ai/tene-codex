# Security Policy

tene-codex is pre-alpha software. Do not use it as the only control protecting production systems or credentials.

## Reporting a vulnerability

Do not open a public issue for a vulnerability or a suspected secret exposure. Report it privately to **Kay Kim at kay@agentkay.it** with the affected version, reproduction steps, impact, and any suggested remediation. Do not include live credentials; use test canaries.

## Secret boundary

The workflow engine must never read `.tene/**`, invoke `tene get`, request a secret value in a model conversation, or store a secret in state, documents, graphs, logs, or QA evidence. Secret-dependent child commands must run through `tene run --env <environment> -- <command>`.

If a release violates this boundary, treat it as a security incident: stop using the affected artifact, rotate exposed credentials, preserve sanitized diagnostic evidence, and report the issue privately.

