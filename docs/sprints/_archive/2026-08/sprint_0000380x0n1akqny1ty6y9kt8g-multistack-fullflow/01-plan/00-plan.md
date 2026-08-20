---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380x0n1akqny1ty6y9kt8g
phase: plan
status: complete
revision: 804
intent_ids: [intent_0000380x0n5an517gpwmw7rrcg]
generated_at: 2026-08-20T02:34:39Z
generated_by: tene-workflow
---

# plan — Multi-stack Intelligence and Full-flow References

<!-- tene:section:purpose -->
## Purpose

Deliver semantic multi-stack analysis and full-flow evidence as independently testable work packages.

<!-- tene:section:scope -->
## Scope

Code intelligence providers/tests, Next.js and Python references, Python QA adapter, root verification wiring and Sprint evidence.

<!-- tene:section:layers -->
## Layers

CLI/analyzer interface; extraction and QA discovery rules; fixture persistence state; Node/Python/Go runtimes and Make verification.

<!-- tene:section:six-questions -->
## Six questions

Provider functions live in `internal/codeintel`, are called by `Analyze`, consume bounded source text and return components/edges/diagnostics. Adapter discovery lives in `internal/qaadapter`, consumes repository capabilities and returns allowlisted execution commands.

<!-- tene:section:traceability -->
## Traceability

AC `ac_0000380x0n5anrn57q3bhyrjj8`; semantic task `task_0000380x15tjkzr77rj40s62wm`; journey task `task_0000380x15wqs0j3kqh8k7zps4`.

<!-- tene:section:decisions -->
## Decisions

First establish explicit provider contracts, then add realistic fixtures, then require their journeys in `make check`.

<!-- tene:section:freeform -->
## Freeform

Runtime tests compensate for intentionally bounded static analysis.

<!-- tene:section:work-packages -->
## Work packages

1. Extract TypeScript/Python declarations, imports, calls, inputs, outputs and effects with explicit dynamic uncertainty.
2. Model four layers in Next.js and Python fixture trees.
3. Execute checkout and API→queue→worker→database journeys.
4. Discover Python unittest as an allowlisted QA adapter and enforce all references in CI checks.

<!-- tene:section:dependencies -->
## Dependencies

Standard Go library, Node standard library and Python standard library; existing npm runner and workflow evidence engine.

<!-- tene:section:verification -->
## Verification

Focused codeintel/qaadapter tests, reference matrix, `npm run test:references`, `go test -race ./...`, `go vet ./...`, full `make check`, and structured seven-layer QA evidence.

<!-- tene:section:risks -->
## Risks

Regex syntax extraction can over/under-match; limit confidence, preserve unknown dynamic edges and prove behavior at runtime.

<!-- tene:section:yagni -->
## Yagni

No tree-sitter/compiler embedding or fixture framework dependency until measured repositories require it.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `804`
- Sprint: `sprint_0000380x0n1akqny1ty6y9kt8g`
- Intents: `intent_0000380x0n5an517gpwmw7rrcg`
- Tasks: `task_0000380x15tjkzr77rj40s62wm`, `task_0000380x15wqs0j3kqh8k7zps4`

<!-- tene:generated:traceability:end -->
