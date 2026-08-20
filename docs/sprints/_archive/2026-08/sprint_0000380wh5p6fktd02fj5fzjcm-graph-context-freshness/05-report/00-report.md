---
schema_version: 1.0.0
document_type: report
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wh5p6fktd02fj5fzjcm
phase: report
status: draft
revision: 197
intent_ids: []
generated_at: 2026-08-20T00:19:23Z
generated_by: tene-workflow
---

# report — Graph Context Freshness

<!-- tene:section:purpose -->
## Purpose

Record how this Sprint completes the graph-impact and context-freshness foundation needed by later workflow, QA, and routing Sprints.

<!-- tene:section:scope -->
## Scope

Directed impact, graph integrity, bounded phase context, provenance validation, CLI/plugin guidance, and executable regression evidence.

<!-- tene:section:layers -->
## Layers

- Interface: new `graph impact` and `context validate`; expanded `context build` flags.
- Business Logic: deterministic traversal, invariant findings, ranked/deduplicated selection, mandatory-overflow and freshness rules.
- Persistence: repository-relative optional context JSON written via sibling temp+rename; state remains source-only and unchanged by reads.
- Infrastructure: CodeGraph/Go AST/filesystem capability snapshot and SHA-256 file provenance.

<!-- tene:section:six-questions -->
## Six questions

- Names: `Impact`, `ValidateGraph`, `BuildContextPack`, `ValidateContextPack`, their DTOs, plus app adapters `safeRootPath` and `writeDerivedJSON`.
- Definitions: `internal/tracecontext/{graph,context}.go` and `internal/app/app.go`.
- Imports/references: app imports tracecontext; tests directly import its public package API; `tene-status`, README, and CLI reference document commands.
- Callers: `runtime.graph`, `runtime.context`, `validateGraph`/doctor, and unit/CLI integration tests.
- Inputs: immutable domain Graph/Project/Sprint, repository root, phase, UTF-8 byte budget, provider capabilities, or a serialized pack.
- Outputs/mutations: sorted impact paths/findings/context/freshness DTOs; only an explicit `--output` writes a derived JSON file, while graph build remains the existing state mutation.

<!-- tene:section:traceability -->
## Traceability

AC-GC-01..04 are passed by run `run_0000380whyn7cf2pp7wnqg86c8`; AC-GC-05 is supported by evidence `evidence_0000380wj24t0gpbxhh14dcecg`. This realizes WP-06 impact/invariants and WP-07 phase/budget/freshness.

<!-- tene:section:decisions -->
## Decisions

Use deterministic UTF-8 bytes rather than pretend tokenizer precision; fail mandatory overflow; represent unresolved calls as low-confidence nodes; never auto-index missing CodeGraph.

<!-- tene:section:freeform -->
## Freeform

Context output intentionally exposes excluded counts/bytes so an agent cannot silently mistake a partial pack for complete repository understanding.

<!-- tene:section:previous-sprints -->
## Previous sprints

Builds on Reference Web Journey (`sprint_0000380v4nmvxm0xghdmr739yw`): its QA/evidence graph can now be impact-traversed and safely packed for later phases. Existing waiver/schema recovery remains intact.

<!-- tene:section:changed-files -->
## Changed files

- Added `internal/tracecontext/graph.go`, `context.go`, and tests.
- Updated `internal/app/app.go` and `app_test.go` for CLI integration and graph materialization.
- Updated English/Korean README command examples, `references/cli.md`, and `skills/tene-status/SKILL.md`.
- Added this Sprint's PRD, plan, design, analysis, QA, report, state events, and verification evidence.

<!-- tene:section:intent-fulfillment -->
## Intent fulfillment

Codex can now ask “what ACs can this node affect?”, receive bounded paths, build only phase-relevant context under an explicit budget, see what was excluded, and prove the pack still matches state/files before acting.

<!-- tene:section:qa-verdict -->
## Qa verdict

PASS: 12/12 cases, four blocking ACs, one hash-valid consolidated evidence artifact. `make check`, race, vet, Playwright 3/3, plugin validator, skill validators 9/9, evidence verify, and doctor all passed.

<!-- tene:section:deferred-work -->
## Deferred work

- Full polyglot/LSP/CodeGraph machine parsing remains WP-08 and the future reference-matrix Sprint.
- The pre-existing live downgrade/cross-major migration gap is not owned by this Sprint and retains its bounded waiver.
- Context budget ratios and recent decision/history stores need richer domain data; current ordering and exclusions are deterministic but minimal.

<!-- tene:section:next-sprint -->
## Next sprint

Workflow, approval, and iterative loop completion: human approval records, meaningful strict/standard/light behavior, bounded repair iterations, failure classification, and deferred lifecycle.

<!-- tene:generated:summary:start -->
### Generated Sprint Summary

- Sprint: `sprint_0000380wh5p6fktd02fj5fzjcm`
- Previous sprints: `sprint_0000380v4nmvxm0xghdmr739yw`
- Intent IDs: `intent_0000380wh9zhz8jw25h22w5za8`
- Tasks: 4
- QA verdict: `passed`
- Open gaps: 0
- State revision: 239

<!-- tene:generated:summary:end -->
