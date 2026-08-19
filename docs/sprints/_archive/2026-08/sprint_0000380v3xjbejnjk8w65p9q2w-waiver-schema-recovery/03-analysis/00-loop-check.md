---
schema_version: 1.0.0
document_type: analysis
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v3xjbejnjk8w65p9q2w
phase: loop-check
status: draft
revision: 127
intent_ids: []
generated_at: 2026-08-19T17:43:57Z
generated_by: tene-workflow
---

# loop-check — Waiver Schema Migration and Recovery

<!-- tene:section:purpose -->
## Purpose

Compare waiver, migration and repair implementation with this Sprint's contracts.

<!-- tene:section:scope -->
## Scope

Domain, state, workflow, CLI, schema, tests and operator references.

<!-- tene:section:layers -->
## Layers

Interface commands; Business waiver rules; Persistence migration/repair; Infrastructure atomic files/locks.

<!-- tene:section:six-questions -->
## Six questions

`Waiver` is defined in domain and referenced by workflow/app/graph. `PlanMigration`, `Migrate`, and `RepairDerived` are defined in state and called by app; inputs are fixed actions/metadata and outputs are structured plans/paths with project/journal mutations only where specified.

<!-- tene:section:traceability -->
## Traceability

AC1 waiver tests and live gap; AC2 legacy migration fixture; AC3 unknown-field/repair tests and live doctor repair.

<!-- tene:section:decisions -->
## Decisions

The live cross-major/downgrade gap is explicitly waived until 2026-08-22, demonstrating rather than bypassing the policy mechanism.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:baseline -->
## Baseline

No waiver domain, migration command, strict unknown-field rejection or repair operation existed.

<!-- tene:section:changed-artifacts -->
## Changed artifacts

Domain/state/workflow/app/tests/schema, workflow/CLI references, Sprint skill and bilingual README.

<!-- tene:section:gap-matrix -->
## Gap matrix

AC1/AC2/AC3 matched. Open `gap_0000380v4d3aas9weqex1ajbb0` is non-security scope debt with active explicit waiver, not an implementation mismatch.

<!-- tene:section:iterations -->
## Iterations

Added contracts, corrected trace IDs through creation-time validation, ran current-schema status and live projection repair.

<!-- tene:section:regression -->
## Regression

Go tests and vet pass; full race/plugin/skill/evidence gates remain in QA.
