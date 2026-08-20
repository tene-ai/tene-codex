---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wney2n4hv4wct0skz4w
phase: qa
status: draft
revision: 340
intent_ids: []
generated_at: 2026-08-20T00:56:52Z
generated_by: tene-workflow
---

# qa — Skill Routing and Eval Completion

<!-- tene:section:purpose -->
## Purpose

Independently gate bilingual skill selection, ambiguity safety and release integration.

<!-- tene:section:scope -->
## Scope

All nine skill routes, public CLI, wrong-phase/hard-negative/multi-intent variants and regression suite.

<!-- tene:section:layers -->
## Layers

L1 pure rule tests; L2 CLI contract; L3 read-only project-state integration; L5 negative/adversarial; L6 race/reliability; L7 full `make check`. L4 is N/A because this feature has no visual UI.

<!-- tene:section:six-questions -->
## Six questions

The evaluator defined in `cmd/tene-routing-eval` imports `router.Route`, is invoked by Makefile/QA, accepts versioned JSON cases, and returns metrics with nonzero failure status. The CLI handler is invoked through `tene-workflow route` and returns the same decision in the standard envelope.

<!-- tene:section:traceability -->
## Traceability

Nine generated QA cases cover the three blocking ACs across happy/error/recovery variants.

<!-- tene:section:decisions -->
## Decisions

One immutable aggregate artifact may evidence multiple cases because it contains all variants and all criteria; every case explicitly references its evidence ID.

<!-- tene:section:freeform -->
## Freeform

No credential or external service is required.

<!-- tene:section:environment -->
## Environment

Local macOS workspace, Go toolchain, Python unittest, plugin and skill validator virtual environment.

<!-- tene:section:capabilities -->
## Capabilities

Go native tests, race detector, vet, deterministic corpus runner and structured CLI are available. Browser is not applicable.

<!-- tene:section:charters -->
## Charters

Positive/explicit selection; adjacent-negative non-trigger; wrong-phase proposal/no mutation; multi-intent Sprint orchestration; malformed/missing input; CI threshold failure contract.

<!-- tene:section:ux-data-flow -->
## Ux data flow

User prompt → CLI/skill route → candidate reasons → safe selection/proposal/none. Project projection is read to derive phase and verified unchanged by the route operation.

<!-- tene:section:evidence -->
## Evidence

`04-qa/evidence/routing-eval-summary.json`; full command logs were observed locally and summary values are reproducible from the checked-in runner/corpus.

<!-- tene:section:verdict -->
## Verdict

PASS. QA run `run_0000380wpc76mrbdegm10wvyj4` passed 9/9 happy/error/recovery cases with evidence `evidence_0000380wpecf38gdhhq2cvsbw8`. `make check`, race, vet, Playwright 3/3, plugin validator, nine skill validators, evidence verification and projection doctor all passed. No residual risk or waiver.
