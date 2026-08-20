# Intent-driven QA Contract

QA begins from confirmed acceptance criteria, not existing test files. Cover the relevant happy, alternate, empty, validation, permission, downstream failure, retry, and recovery paths.

Prefer deterministic project-native tests and Playwright for regression. Add Codex browser, Chrome, or other interactive tools for exploratory UX checks. Observe the full chain where applicable: user action → interface state → API or command boundary → business rule → persistence or external side effect → user feedback.

Start with `qa capabilities` and `qa plan`. Use `qa execute <case-id> --adapter go-test|npm-test|playwright` for allowlisted native runners. For Codex browser/Chrome or external observers, create a `schemas/qa-observation.schema.json` artifact and import it with `qa observe <case-id> --input <path>`. The observation must bind the run, case, QA-plan revision and specification hash; name its tool version and environment; include before/after checkpoints; and provide actual-versus-expected assertions linked to layers and requirement references. Generic sanitized artifacts registered by `evidence register` do not qualify as QA case proof by themselves.

Manual pass is forbidden. `qa case <case-id> failed` may record an explicit failure, while `qa evaluate` derives pass status from the evidence. Every required layer and the case's `observable`, `variant:<name>`, `expected:<index>` and `forbidden:<index>` references must be covered by passing assertions. A layer can be changed with `qa disposition <case-id> <layer> --status required|not-applicable|waived`; non-required states require an approver and reason. Artifact bytes, size, hash, redaction, freshness and identity are checked again at evaluation.

When possible, use a separate evaluator subagent. Give it the PRD/AC revision, charters, diff summary, and evidence manifest. Do not include the builder's completion claim.
