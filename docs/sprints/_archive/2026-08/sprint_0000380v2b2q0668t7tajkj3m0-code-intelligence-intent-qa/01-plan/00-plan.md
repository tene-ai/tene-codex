---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v2b2q0668t7tajkj3m0
phase: plan
status: draft
revision: 88
intent_ids: []
generated_at: 2026-08-19T17:30:10Z
generated_by: tene-workflow
---

# plan — Code Intelligence and Intent QA Adapters

<!-- tene:section:purpose -->
## Purpose

Deliver WP-08 and WP-10 as one vertical slice and prove it through the workflow itself.

<!-- tene:section:scope -->
## Scope

Provider contracts, Go AST fallback, CodeGraph capability probe, graph merge/materialization, QA capability probes, charter enrichment, observation import, schemas, skills, documentation and regression tests.

<!-- tene:section:layers -->
## Layers

Interface: CLI/skills. Business: codeintel and qaadapter packages. Persistence: projected graph/evidence. Infrastructure: filesystem, git and tool probes.

<!-- tene:section:six-questions -->
## Six questions

New public declarations must be discoverable by their own analyzer; tests assert definition, imports/calls, input/output and effects rather than only command exit codes.

<!-- tene:section:traceability -->
## Traceability

WP08-T1 provider contract → AC1; WP08-T2 Go analyzer/layers/6Q → AC1; WP10-T1 capabilities/charters → AC2; WP10-T2 observation evidence → AC2; integration/release checks → AC3.

<!-- tene:section:decisions -->
## Decisions

Keep analysis read-only and bounded to repository files. Persist graph deltas, but return the richer understanding report on demand. Model browser tools as external executors with a versioned artifact boundary.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:work-packages -->
## Work packages

1. Add `internal/codeintel` contracts, discovery, Go AST analysis and tests.
2. Merge code graph nodes/edges and add `graph providers|understand` commands.
3. Add `internal/qaadapter` discovery/observation contracts and tests.
4. Enrich QACase and implement `qa capabilities|observe`.
5. Update JSON schemas, skills/references and bilingual README.
6. Run loop check, full QA, report and archive.

<!-- tene:section:dependencies -->
## Dependencies

Domain types precede adapters; adapters precede CLI integration; tests precede documentation evidence. Existing state schema remains backward-compatible via optional fields.

<!-- tene:section:verification -->
## Verification

`go test ./...`, `go test -race ./...`, `go vet ./...`, hook tests, schema parsing, official plugin/skill validators, CLI fixture journeys, evidence hash verification and doctor.

<!-- tene:section:risks -->
## Risks

Static call analysis can overclaim dynamic behavior: assign confidence and unknowns. Browser observations can be forged: hash artifacts and preserve executor metadata, but do not claim attestation. Large repositories can be slow: ignore VCS/vendor/generated directories and bound file size.

<!-- tene:section:yagni -->
## Yagni

No automatic index creation, daemon, embedded browser, generalized LSP protocol, remote evidence store or arbitrary command execution.
