---
schema_version: 1.0.0
document_type: report
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wney2n4hv4wct0skz4w
phase: report
status: draft
revision: 340
intent_ids: []
generated_at: 2026-08-20T00:56:52Z
generated_by: tene-workflow
---

# report — Skill Routing and Eval Completion

<!-- tene:section:purpose -->
## Purpose

Review and preserve the completed WP-12 routing/evaluation capability.

<!-- tene:section:scope -->
## Scope

Deterministic router, public CLI, bilingual corpus, evaluator, CI gate, documentation and complete validation evidence.

<!-- tene:section:layers -->
## Layers

- Interface: `tene-workflow route` returns standard JSON/human outcomes.
- Business Logic: explicit precedence, cue/phase scoring, ambiguity and multi-intent safety.
- Persistence: corpus is versioned; project state is read-only during routing.
- Infrastructure: routing evaluator is mandatory in `make check`.

<!-- tene:section:six-questions -->
## Six questions

Names are `router.Route`, `Decision`, `Candidate`, `runtime.route`, and evaluator `main`. They are defined in `internal/router/router.go`, `internal/app/app.go`, and `cmd/tene-routing-eval/main.go`; imported by app/evaluator and invoked by CLI/Makefile. Inputs are prompt, active flag, phase, or corpus path. Outputs are non-mutating decision JSON or per-skill metrics/process status.

<!-- tene:section:traceability -->
## Traceability

Confirmed intent `intent_0000380wp1g79k7a136ay9wx00`; three blocking ACs; tasks `task_0000380wp7q4t8mjxdp4ysr78r` and `task_0000380wp7rdq0awdk0amq65fm`; evidence `evidence_0000380wpecf38gdhhq2cvsbw8`.

<!-- tene:section:decisions -->
## Decisions

Keep native Codex skill discovery as the conversational entry and use this router as auditable companion/preflight. Never let implicit routing authorize approvals or skip phases.

<!-- tene:section:freeform -->
## Freeform

Retrospective: making multi-intent performance a separate metric caught failures hidden by precision/recall. Future new cues must update the negative and collision corpus.

<!-- tene:section:previous-sprints -->
## Previous sprints

Continues `sprint_0000380wmbmks792hdhe03j7bg` (state recovery): routing safely reads the durable active phase and inherits replay/drift guarantees. It completes the previously planned WP-12 surface atop lifecycle, graph/context and approval-loop work.

<!-- tene:section:changed-files -->
## Changed files

- Added `internal/router/router.go`, `internal/router/router_test.go`: routing contract and tests.
- Added `cmd/tene-routing-eval/main.go`, `evals/routing-corpus.json`: expanded corpus gate.
- Modified `internal/app/app.go`, `references/cli.md`: public read-only command.
- Modified `Makefile`: mandatory eval gate.
- Added/updated Sprint documents, evidence and `.tene-workflow` projections.

<!-- tene:section:intent-fulfillment -->
## Intent fulfillment

Users can name any `$tene-*` explicitly or ask naturally in Korean/English. Clear requests select; ambiguity and multi-intent requests propose; unrelated/wrong-phase requests cannot mutate. Every result explains candidates and safe next actions.

<!-- tene:section:qa-verdict -->
## Qa verdict

PASS: run `run_0000380wpc76mrbdegm10wvyj4`, 9/9 cases. Explicit 45/45; precision 1.0; recall ≥0.95; multi-intent ≥0.90; wrong-phase and unnecessary trigger 0. Full Go/race/vet, Playwright, plugin/skill validators, evidence integrity and doctor passed.

<!-- tene:section:deferred-work -->
## Deferred work

None. No policy choice, waiver, open gap or routing task was carried forward.

<!-- tene:section:next-sprint -->
## Next sprint

Secret and adversarial QA completion: prove value non-disclosure, canary handling, permission/child failure, evidence poisoning and evaluator false-pass resistance.


<!-- tene:generated:summary:start -->
### Generated Sprint Summary

- Sprint: `sprint_0000380wney2n4hv4wct0skz4w`
- Previous sprints: ``
- Intent IDs: `intent_0000380wp1g79k7a136ay9wx00`
- Tasks: 2
- QA verdict: `passed`
- Open gaps: 0
- State revision: 372

<!-- tene:generated:summary:end -->
