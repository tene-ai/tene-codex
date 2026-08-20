---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380x2dmqv18kkk1fr627v0
phase: qa
status: complete
revision: 871
intent_ids: [intent_0000380x2dqcf9ctzpr01bb7jr]
generated_at: 2026-08-20T02:50:07Z
generated_by: tene-workflow
---

# qa — Semantic Contract Completion Audit

<!-- tene:section:purpose -->
## Purpose

Independently prove every MVP contract and demonstrate that the auditor rejects incomplete, invalid, unauthorized and failed proofs.

<!-- tene:section:scope -->
## Scope

33 requirements, 19 fixed command groups, seven L1–L7 variants, repository state/evidence integrity and full race/vet/check/release/reference/security regression.

<!-- tene:section:layers -->
## Layers

L1 symbol/unit, L2 package integration, L3 schema/CLI contracts, L4 full repository system, L5 workflow user journey, L6 mutation/recovery regression and L7 security/release/operational integrity.

<!-- tene:section:six-questions -->
## Six questions

The manifest names definitions and files; tests prove callers and invocations; fixed argv lists define inputs; JSON audit results, command exits, state mutations and archived evidence define outputs/effects.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380x2dqcf9ctzpr01bb7jr`; AC `ac_0000380x2dqce7nbcsa2c3b3g8`; authoritative run `run_0000380x4x5bx76nykm72c0w18`; seven observation files.

<!-- tene:section:decisions -->
## Decisions

Negative variants pass only when the invalid condition is rejected with the expected nonzero result. `--contracts-only` breaks the pre-evidence circularity but still runs all semantic and command proofs.

<!-- tene:section:freeform -->
## Freeform

No manual pass, waiver or not-applicable layer disposition.

<!-- tene:section:environment -->
## Environment

Local fingerprint in the QA run; Go, Python, Node/npm, Playwright, packaging shell and Codex/plugin artifacts are exercised without secret values.

<!-- tene:section:capabilities -->
## Capabilities

Go unit/race/vet, Python unittest, semantic regex resolution, fixed subprocess registry, npm references, Playwright, routing eval, schema/plugin JSON validation, release packaging and journal/evidence verification.

<!-- tene:section:charters -->
## Charters

Happy full proof; canonical-order alternate; missing contract; missing symbol; unauthorized command; executable timeout failure; pristine recovery rerun.

<!-- tene:section:ux-data-flow -->
## Ux data flow

User completion claim→contract IDs→source symbols→fixed test commands→structured result→workflow/evidence gate. Invalid inputs stop before completion; recovery restores a clean all-pass result.

<!-- tene:section:evidence -->
## Evidence

`04-qa/observations/*.json` is run/case/spec/state bound with content hashes, L1–L7 assertions, actual mutation outcomes and no secret canary. Full repository commands also passed outside the observation generator.

<!-- tene:section:verdict -->
## Verdict

Passed at revision 865 with zero QA findings. Pre-archive semantic audit passed all 33 contracts and 19 command groups with zero workflow failures; race/vet/make check also passed.


<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `871`
- Sprint: `sprint_0000380x2dmqv18kkk1fr627v0`
- Intents: `intent_0000380x2dqcf9ctzpr01bb7jr`
- Tasks: `task_0000380x2g10m4b5hef6hemvb4`, `task_0000380x2g280ctd6py5fq0bk4`, `task_0000380x2g3eszpq5cxa1vtk8g`

<!-- tene:generated:traceability:end -->
