---
schema_version: 1.0.0
document_type: report
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wx6xt073n63as8wkzyw
phase: report
status: complete
revision: 560
intent_ids: [intent_0000380wx710enjk9h81emvgf8]
generated_at: 2026-08-20T02:04:35Z
generated_by: tene-workflow
---

# report — Automatic Bidirectional Loop Analysis

<!-- tene:section:purpose -->
## Purpose

Record delivery of the automatic bidirectional Loop Check foundation.

<!-- tene:section:scope -->
## Scope

Analyzer, task artifact ownership, executable design contracts, gap reconciliation and mutation QA.

<!-- tene:section:layers -->
## Layers

Interface: task/loop CLI. Business Logic: analyzer/reconciler. Persistence: fingerprinted gap lifecycle. Infrastructure: git/filesystem providers.

<!-- tene:section:six-questions -->
## Six questions

`Analyze`, `Reconcile`, `Candidate`, task artifacts and gap fingerprints are defined in loopcheck/domain, imported and called by app, consume canonical state/docs/files, and return persisted gaps plus convergence metrics.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380wx710enjk9h81emvgf8`; AC `ac_0000380wx710e3a7bw0g1kd6dc`; tasks `task_0000380wxj3hep0x4eaqm3bkdm`, `task_0000380wxj4kryss8dx8k5wdrm`.

<!-- tene:section:decisions -->
## Decisions

Explicit IDs/artifact links/contracts are proof; prose alone is not. Only analyzer-owned gaps auto-resolve.

<!-- tene:section:freeform -->
## Freeform

This replaces the previous false implementation where loop check merely listed pre-existing gaps.

<!-- tene:section:previous-sprints -->
## Previous sprints

Builds on QA evidence integrity so loop results and later QA use independent proof instead of completion claims.

<!-- tene:section:changed-files -->
## Changed files

`internal/loopcheck/analyzer.go` and tests; `internal/app/app.go` task artifact/loop integration; `internal/domain/types.go` gap provenance; Sprint evidence and documentation.

<!-- tene:section:intent-fulfillment -->
## Intent fulfillment

Automatically detects specification coverage, task edges, changed-file ownership, missing paths/symbols and forbidden dependencies, then converges stable gaps after repair.

<!-- tene:section:qa-verdict -->
## Qa verdict

Passed run `run_0000380wxtrexn4d1f6wqcte4w`; 6/6 seeded classes, seven variants, native Go evidence and full integrity gate.

<!-- tene:section:deferred-work -->
## Deferred work

Universal natural-language semantic equivalence remains intentionally unsupported; CodeGraph/native providers and explicit contracts supply bounded semantics. No required Sprint item deferred.

<!-- tene:section:next-sprint -->
## Next sprint

Complete document sync/generated regions and exact CLI/idempotency contracts.


<!-- tene:generated:summary:start -->
### Generated Sprint Summary

- Sprint: `sprint_0000380wx6xt073n63as8wkzyw`
- Previous sprints: `sprint_0000380wvj5zg3avxfejqj4wj0`
- Intent IDs: `intent_0000380wx710enjk9h81emvgf8`
- Tasks: 2
- QA verdict: `passed`
- Open gaps: 0
- State revision: 613

<!-- tene:generated:summary:end -->

<!-- tene:generated:summary:start -->
### Generated Sprint Summary

- Sprint: `sprint_0000380wx6xt073n63as8wkzyw`
- Previous sprints: `sprint_0000380wvj5zg3avxfejqj4wj0`
- Intent IDs: `intent_0000380wx710enjk9h81emvgf8`
- Tasks: 2
- QA verdict: `passed`
- Open gaps: 0
- State revision: 616

<!-- tene:generated:summary:end -->
