# Intent-driven QA Contract

QA begins from confirmed acceptance criteria, not existing test files. Cover the relevant happy, alternate, empty, validation, permission, downstream failure, retry, and recovery paths.

Prefer deterministic project-native tests and Playwright for regression. Add Codex browser, Chrome, or other interactive tools for exploratory UX checks. Observe the full chain where applicable: user action → interface state → API or command boundary → business rule → persistence or external side effect → user feedback.

Start with `qa capabilities` and `qa plan`. Use `qa execute <case-id> --adapter go-test|npm-test|playwright` for allowlisted native runners. For Codex browser/Chrome or external observers, create a `schemas/qa-observation.schema.json` artifact and import it with `qa observe <case-id> --input <path>`; this validates run/case identity, timestamps, assertions and redaction before hashing and linking evidence. Generic sanitized artifacts may still use `evidence register --path <path> --ac <id>`. Each remaining case is marked through `qa case <case-id> passed|failed --evidence <ids>`. `qa evaluate` is authoritative: all cases for every blocking criterion need valid, criterion-linked, redaction-safe evidence.

When possible, use a separate evaluator subagent. Give it the PRD/AC revision, charters, diff summary, and evidence manifest. Do not include the builder's completion claim.
