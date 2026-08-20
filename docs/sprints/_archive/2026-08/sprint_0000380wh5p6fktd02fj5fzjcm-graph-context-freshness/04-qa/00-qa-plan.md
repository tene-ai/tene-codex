---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wh5p6fktd02fj5fzjcm
phase: qa
status: draft
revision: 197
intent_ids: []
generated_at: 2026-08-20T00:19:23Z
generated_by: tene-workflow
---

# qa — Graph Context Freshness

<!-- tene:section:purpose -->
## Purpose

Independently verify every blocking AC using executable output and hash-addressed evidence.

<!-- tene:section:scope -->
## Scope

Unit/contract, CLI integration, race, vet, repository checks, plugin/skill validation, live graph/context journeys, evidence verification, and doctor health.

<!-- tene:section:layers -->
## Layers

Exercise CLI interface, graph/context business logic, saved derived pack persistence, and provider capability infrastructure.

<!-- tene:section:six-questions -->
## Six questions

Verify exported functions via tests and inspect app call sites; capture input/output shapes in Go AST graph output and this report. Confirm only the explicit context output path changes.

<!-- tene:section:traceability -->
## Traceability

Each AC receives happy, error, and recovery cases from `qa plan`; one consolidated test artifact may prove multiple cases only when its commands and exit statuses are recorded.

<!-- tene:section:decisions -->
## Decisions

Graph warnings for not-yet-registered QA evidence are expected before evidence registration; dangling references remain blocking.

<!-- tene:section:freeform -->
## Freeform

Browser UX is not applicable because this Sprint changes a CLI, but the CLI journey includes save, immediate validate, mutate provenance, fail stale, rebuild, and pass recovery.

<!-- tene:section:environment -->
## Environment

Local macOS workspace, Go toolchain, Node/Playwright dependencies already installed, no CodeGraph index.

<!-- tene:section:capabilities -->
## Capabilities

Go unit/race/vet, CLI integration, plugin and skill validators, filesystem hashing. CodeGraph unavailable by repository policy; browser unnecessary for this CLI-only surface.

<!-- tene:section:charters -->
## Charters

- Impact: valid directed path, call cutoff, cycle, unknown node.
- Graph: valid current build, dangling edge, orphan task/evidence.
- Context: PRD filtering, design inclusion, optional exclusion, mandatory overflow.
- Freshness: immediate pass, state drift, file drift, rebuild recovery, path escape.

<!-- tene:section:ux-data-flow -->
## Ux data flow

User command → app flag parser → state/project and provider snapshot → pure tracecontext operation → JSON response/optional derived file → validator re-reads file hashes and current revision.

<!-- tene:section:evidence -->
## Evidence

`evidence_0000380wj24t0gpbxhh14dcecg` points to `04-qa/evidence/verification.txt`, SHA-256 `a6f2fde205d3a827d346871a6c3d1ca65484ca2744301d7f079ca3eb0d58e0cd`. It records core/race/vet, Playwright regression, plugin and 9 skill validators, live CLI journeys, evidence integrity, and doctor health.

<!-- tene:section:verdict -->
## Verdict

PASS. Run `run_0000380whyn7cf2pp7wnqg86c8` passed all 12 happy/error/recovery cases for four blocking ACs with hash-valid evidence. No blocker or secret leak remains.
