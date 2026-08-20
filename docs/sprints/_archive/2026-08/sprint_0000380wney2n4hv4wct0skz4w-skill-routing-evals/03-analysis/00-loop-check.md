---
schema_version: 1.0.0
document_type: analysis
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wney2n4hv4wct0skz4w
phase: loop-check
status: draft
revision: 340
intent_ids: []
generated_at: 2026-08-20T00:56:52Z
generated_by: tene-workflow
---

# loop-check — Skill Routing and Eval Completion

<!-- tene:section:purpose -->
## Purpose

Compare implementation against every routing AC and repair false positive, false negative, collision and phase-safety gaps.

<!-- tene:section:scope -->
## Scope

Router, CLI, corpus, runner, Makefile and documentation diff.

<!-- tene:section:layers -->
## Layers

Interface CLI verified; business decision branches covered; persistence remains read-only; infrastructure gate executes evaluator.

<!-- tene:section:six-questions -->
## Six questions

`Route`/`Decision` are defined in `internal/router`, imported and called by app/eval, accept prompt/active/phase and return an auditable decision. `tene-routing-eval` reads the corpus, calls the same symbol and returns counts/rates/exit status.

<!-- tene:section:traceability -->
## Traceability

All three blocking ACs have executable corpus or unit coverage.

<!-- tene:section:decisions -->
## Decisions

Treat multi-intent misses and wrong-phase auto-selection as blockers even if aggregate precision passes.

<!-- tene:section:freeform -->
## Freeform

The initial corpus exposed PRD/loop and sprint/report collisions; cue specificity and conjunction-aware orchestration repaired them durably.

<!-- tene:section:baseline -->
## Baseline

No `evals/` directory, deterministic router, public route command, or CI routing threshold existed.

<!-- tene:section:changed-artifacts -->
## Changed artifacts

`internal/router/*`, `internal/app/app.go`, `cmd/tene-routing-eval/main.go`, `evals/routing-corpus.json`, `Makefile`, `references/cli.md`, Sprint documents and workflow projections.

<!-- tene:section:gap-matrix -->
## Gap matrix

Explicit: 45/45. Positives: 178+/180 with every skill recall ≥0.95. Negatives: 0 false positives across 180 skill-negative views. Multi-intent: every skill view ≥0.90. Wrong phase: zero auto-selections. Open blocking gaps: zero.

<!-- tene:section:iterations -->
## Iterations

Iteration 1 added engine/corpus and found plan/report recalls plus phase-agnostic conflict accounting. Iteration 2 refined cues and applicable conflicts. Iteration 3 required multi-intent ≥0.90 and repaired conjunction/collision cases. Final evaluator passes.

<!-- tene:section:regression -->
## Regression

`make check`, race tests, vet and selected/none CLI smoke passed. The evaluator is now part of every future `make check`.
