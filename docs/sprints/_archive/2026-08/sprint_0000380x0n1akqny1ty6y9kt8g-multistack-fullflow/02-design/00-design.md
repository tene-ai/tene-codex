---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380x0n1akqny1ty6y9kt8g
phase: design
status: complete
revision: 804
intent_ids: [intent_0000380x0n5an517gpwmw7rrcg]
generated_at: 2026-08-20T02:34:39Z
generated_by: tene-workflow
---

# design — Multi-stack Intelligence and Full-flow References

<!-- tene:section:purpose -->
## Purpose

Define bounded semantic-provider and executable reference contracts that cannot false-pass as filesystem coverage.

<!-- tene:section:scope -->
## Scope

Provider capability metadata, source extraction, component/edge projection, layer matrix, QA discovery and observable journeys.

<!-- tene:section:layers -->
## Layers

Source entry points map to Interface; services/workers to Business Logic; queue/repository/file database to Persistence; runtime configuration to Infrastructure.

<!-- tene:section:six-questions -->
## Six questions

`analyzeBoundedLanguage` is defined in `internal/codeintel/codeintel.go`, called per supported non-Go source by `Analyze`, accepts file/path/layer/extension, and returns semantic components, call edges and diagnostics. `qaadapter.Discover` returns an absolute allowlisted Python unittest command when `python3` and `tests/` exist.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380x0n5an517gpwmw7rrcg`; AC `ac_0000380x0n5anrn57q3bhyrjj8`; both tasks listed in Plan.

<!-- tene:section:decisions -->
## Decisions

Static providers publish confidence ≤0.75 and retain dynamic incoming-reference uncertainty. A component with empty input/output/effect fields or filesystem provider fails the reference contract.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:contract path="internal/codeintel/codeintel.go" symbol="func analyzeBoundedLanguage" -->
<!-- tene:contract path="internal/qaadapter/qaadapter.go" symbol="func Discover" -->

<!-- tene:section:components -->
## Components

TypeScript static provider, Python static provider, provider capability registry, reference matrix, Node checkout journey, Python worker journey and Python QA adapter.

<!-- tene:section:interfaces -->
## Interfaces

`codeintel.Analyze(ctx, root, paths, changed)`; `qaadapter.Discover(root)`/`Execute`; `npm run test:references`; `make test-references`.

<!-- tene:section:data -->
## Data

Components contain locator/provider/confidence/imports/calls/inputs/outputs/effects/unknown. Journeys return visible confirmation/API status and compare persisted JSON/NDJSON objects.

<!-- tene:section:state-transitions -->
## State transitions

Next.js: checkout request→API accepted→service creates→repository appends→confirmation. Python: API accepted→queue pending→worker consumes→repository stored.

<!-- tene:section:failures -->
## Failures

Unreadable/oversized sources diagnose and skip; unsupported sources degrade explicitly; missing runtime disables an adapter; failed assertion or persistence mismatch exits nonzero.

<!-- tene:section:security -->
## Security

Fixture state uses temporary directories; no credentials/network; analyzer excludes `.tene`, generated dependencies and oversized files.

<!-- tene:section:tests -->
## Tests

Provider selection and Six Questions, four-layer matrix, Python adapter discovery, Node visible/persisted outcome, Python response/worker/persisted identity, root race/vet/check.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `804`
- Sprint: `sprint_0000380x0n1akqny1ty6y9kt8g`
- Intents: `intent_0000380x0n5an517gpwmw7rrcg`
- Tasks: `task_0000380x15tjkzr77rj40s62wm`, `task_0000380x15wqs0j3kqh8k7zps4`

<!-- tene:generated:traceability:end -->
