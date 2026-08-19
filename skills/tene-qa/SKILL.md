---
name: tene-qa
description: Plan, execute, and independently evaluate intent-driven QA across UX and data flows for an active tene sprint.
---

# tene QA

Require the `qa` phase. Read [QA contract](../../references/qa.md), [workflow](../../references/workflow.md), and [security](../../references/security.md) when credentials are involved.

Run `qa capabilities` and `qa plan`, then refine generated cases using the confirmed user and data journeys. Prefer allowlisted `qa execute` native adapters and Playwright. Use browser or Chrome tools for exploratory UX checks, persist their structured result using `schemas/qa-observation.schema.json`, and import it through `qa observe`. Observe interface, business, persistence, and infrastructure boundaries where applicable.

Register only sanitized reproducible artifacts, attach them to the matching criteria, and record every case result. If available, delegate evaluation to a separate read-only subagent with only the criteria, charters, diff summary, and evidence. Run `evidence verify` and `qa evaluate`; never override a failed blocker with a score or assertion of completion.
