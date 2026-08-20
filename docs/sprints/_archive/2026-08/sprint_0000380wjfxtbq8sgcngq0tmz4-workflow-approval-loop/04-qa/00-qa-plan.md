---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wjfxtbq8sgcngq0tmz4
phase: qa
status: draft
revision: 242
intent_ids: []
generated_at: 2026-08-20T00:30:55Z
generated_by: tene-workflow
---

# qa — Workflow Approval Loop Completion

<!-- tene:section:purpose -->
## Purpose

Verify authorization boundaries and bounded failure behavior with evidence independent of completion claims.

<!-- tene:section:scope -->
## Scope

Profile matrix, exact-scope approval, expiry/consumption/dry-run, iteration exhaustion, defer restrictions, schema/default compatibility and full repository regression.

<!-- tene:section:layers -->
## Layers

CLI inputs and error codes; pure guard logic; durable journal/projection; compatibility and validator tooling.

<!-- tene:section:six-questions -->
## Six questions

Inspect domain definitions and app/workflow call paths; assert Approval inputs/records and Sprint/Gap mutations; verify dry-run does not consume state.

<!-- tene:section:traceability -->
## Traceability

Four blocking ACs receive happy/error/recovery cases from qa plan and consolidated hash-valid command evidence.

<!-- tene:section:decisions -->
## Decisions

No approval test may bypass a deliberately failing quality guard. Security and evidence-integrity defer attempts must exit 3.

<!-- tene:section:freeform -->
## Freeform

Browser UX is unchanged; existing Playwright journeys serve regression rather than primary evidence.

<!-- tene:section:environment -->
## Environment

Local macOS, Go toolchain, Node/Chromium, Codex plugin/skill validators, strict active project.

<!-- tene:section:capabilities -->
## Capabilities

Unit, integration, race, vet, Playwright, schema JSON, plugin/skill validation, journal/evidence verification and doctor.

<!-- tene:section:charters -->
## Charters

- Strict/standard/light/off boundary table.
- Missing, requested, approved, wrong-scope, expired and consumed approval.
- Loop remaining/exhaustion at configurable maximum.
- Resolve requires evidence; defer requires owner/reason/target; security defer forbidden.
- Existing complete Sprint still archives through standard approval.

<!-- tene:section:ux-data-flow -->
## Ux data flow

User request→approval request metadata→human approve→shared transition guard→atomic approval consumption+phase mutation. Loop/gap commands similarly commit one event and status projects effective blockers separately.

<!-- tene:section:evidence -->
## Evidence

`evidence_0000380wkk3af0vb6ymwvb3fmr`, SHA-256 `69ed08f80efa9022e0c34ee73ccb56a9f5b957126399edcb48fb8ec60b82e1a4`, records profile/approval/loop/gap behavior, make/race/vet, Playwright 3/3, plugin, 9 skills, evidence integrity and doctor.

<!-- tene:section:verdict -->
## Verdict

PASS. QA run `run_0000380wkft6h3n298wyrt1jdr` passed 12/12 happy/error/recovery cases for all four blocking ACs with hash-valid evidence.
