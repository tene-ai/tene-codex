---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wh5p6fktd02fj5fzjcm
phase: plan
status: draft
revision: 197
intent_ids: []
generated_at: 2026-08-20T00:19:23Z
generated_by: tene-workflow
---

# plan — Graph Context Freshness

<!-- tene:section:purpose -->
## Purpose

Deliver AC-GC-01 through AC-GC-05 as one backward-compatible vertical slice.

<!-- tene:section:scope -->
## Scope

Implement a focused domain package, connect CLI contracts, add fixtures/tests, update skill guidance if the public command contract changes, then run the complete Sprint gates.

<!-- tene:section:layers -->
## Layers

Interface CLI parsing; business graph/context algorithms; persistence explicit pack file; infrastructure provider capability snapshot.

<!-- tene:section:six-questions -->
## Six questions

Each exported DTO/function will be recorded in the report with definition, caller, inputs, outputs, and mutations. CLI writing is isolated from pure pack construction.

<!-- tene:section:traceability -->
## Traceability

Tasks map one-to-one to AC-GC-01..05 and WP-06/07. No task can be completed without its named test evidence.

<!-- tene:section:decisions -->
## Decisions

Create `internal/tracecontext` instead of expanding the large app dispatcher. Use stable sorting everywhere to make snapshots reproducible.

<!-- tene:section:freeform -->
## Freeform

The existing open migration debt remains outside this Sprint and stays governed by its active waiver.

<!-- tene:section:work-packages -->
## Work packages

1. Define impact/path, finding, context item, provenance, exclusion, and freshness DTOs.
2. Implement directed impact traversal and graph invariant validation.
3. Implement phase-aware priority selection, byte budgeting, hashing, JSON persistence, and freshness validation.
4. Wire CLI commands and human/JSON outputs; update usage/skill references.
5. Add unit and CLI integration tests, execute loop/QA gates, report, archive, commit, and push.

<!-- tene:section:dependencies -->
## Dependencies

Existing domain graph/project types, state loader, codeintel capabilities, and the archived Reference Web Journey Sprint.

<!-- tene:section:verification -->
## Verification

`go test ./...`, `go test -race ./...`, `go vet ./...`, focused CLI tests, `make check`, plugin/skill validators, `evidence verify`, and `doctor`.

<!-- tene:section:risks -->
## Risks

- Very small budgets could force unsafe omission: fail with a stable error.
- Directed edges do not always encode reverse change impact: report traversal direction and paths explicitly.
- File churn after pack creation: freshness validation reports exact changed/missing locators.

<!-- tene:section:yagni -->
## Yagni

Do not introduce embeddings, databases, daemons, network calls, or automatic semantic indexing.
