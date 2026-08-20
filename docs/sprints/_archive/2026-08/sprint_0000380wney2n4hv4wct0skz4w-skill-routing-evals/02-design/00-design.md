---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wney2n4hv4wct0skz4w
phase: design
status: draft
revision: 340
intent_ids: []
generated_at: 2026-08-20T00:56:52Z
generated_by: tene-workflow
---

# design — Skill Routing and Eval Completion

<!-- tene:section:purpose -->
## Purpose

Make routing deterministic, state-aware, transparent, and unable to bypass gates.

<!-- tene:section:scope -->
## Scope

`internal/router`, `runtime.route`, `cmd/tene-routing-eval`, corpus and CI wiring.

<!-- tene:section:layers -->
## Layers

Interface: `route`. Business: normalize, explicit/hard-negative branches, scores and policy. Persistence: read-only project plus corpus. Infrastructure: evaluator exit code.

<!-- tene:section:six-questions -->
## Six questions

`Decision`/`Candidate` live in `internal/router`; app and evaluator import them. CLI dispatch calls `runtime.route` with text/overrides. It returns scores, reasons, selected skill, mode, confirmation/mutation flags and action constraints without state changes. Eval main reads a corpus and returns per-skill metrics JSON.

<!-- tene:section:traceability -->
## Traceability

FR-10/WP-12 and all three Sprint ACs are executable through the corpus gate.

<!-- tene:section:decisions -->
## Decisions

Score: 0.40 intent + 0.25 phase + 0.20 artifact + 0.15 action. A single ≥0.80 winner with ≥0.10 margin selects; ≥0.60 proposes. Multiple material intents propose `tene-sprint`. Explicit is 1.0. Unsafe outcomes disable mutation.

<!-- tene:section:freeform -->
## Freeform

Phase-agnostic sprint/status/secrets omit wrong-phase measurement and retain explicit action constraints.

<!-- tene:section:components -->
## Components

Pure concurrency-safe `router.Route`; read-only `runtime.route`; deterministic sorted evaluator.

<!-- tene:section:interfaces -->
## Interfaces

`Route(text string, active bool, phase domain.Phase) Decision`; `tene-workflow route --text TEXT [--phase PHASE] [--active auto|true|false]`; `tene-routing-eval [corpus.json]`.

<!-- tene:section:data -->
## Data

Decision contains mode, selected skill, candidates, explicit/confirmation/mutation flags, required/forbidden actions and reason. Corpus contains 9 skill fixtures, suffixes, 20 negatives and 10 multi-intent prompts.

<!-- tene:section:state-transitions -->
## State transitions

None. The command may read active phase and revision but cannot mutate any projection or journal.

<!-- tene:section:failures -->
## Failures

Missing text is usage error; absent project uses inactive revision zero; malformed corpus exits 2; threshold failure exits 1.

<!-- tene:section:security -->
## Security

Secret cues route only to safe runtime guidance. Routing never authorizes execution and ambiguous/implicit outcomes cannot skip approval.

<!-- tene:section:tests -->
## Tests

Explicit precedence, phase route, multi-intent, hard negative, wrong phase and 405+ expanded decisions plus full repository checks.
