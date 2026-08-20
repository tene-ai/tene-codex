---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wr84t6p1hb89mf2wbkc
phase: qa
status: draft
revision: 437
intent_ids: []
generated_at: 2026-08-20T01:21:13Z
generated_by: tene-workflow
---

# qa — Release and Marketplace Completion

<!-- tene:section:purpose -->
## Purpose

Gate source-to-public-artifact readiness and compatibility honesty.

<!-- tene:section:scope -->
## Scope

All package, install/update/uninstall, supply-chain, portal material and migration cases.

<!-- tene:section:layers -->
## Layers

L1 scripts/migration; L2 manifest/wrapper; L3 state preservation; L5 tamper/unsupported; L6 multi-platform/race; L7 CI/release.

<!-- tene:section:six-questions -->
## Six questions

QA invokes exact release commands and observes file inventory, hashes, SBOM, routing decisions, state path, exit codes and unchanged migration bytes.

<!-- tene:section:traceability -->
## Traceability

Six generated cases cover two blocking ACs; inherited debt resolves only using registered summary.

<!-- tene:section:decisions -->
## Decisions

Do not create tag/public submission until final 100% traceability audit.

<!-- tene:section:freeform -->
## Freeform

Portal credentials and verified organization identity are external prerequisites.

<!-- tene:section:environment -->
## Environment

Local clean mktemp profile plus GitHub workflow definitions.

<!-- tene:section:capabilities -->
## Capabilities

Cross-build, SHA-256, SPDX JSON, routing CLI, validators, Playwright and migration tests.

<!-- tene:section:charters -->
## Charters

Package success; tamper failure; explicit/implicit route; update; uninstall preserve; supported/unsupported migration; listing completeness.

<!-- tene:section:ux-data-flow -->
## Ux data flow

Source/tag → CI build → bundle verification → catalog/install → skill route/core → update/uninstall preserving repository.

<!-- tene:section:evidence -->
## Evidence

`04-qa/evidence/release-summary.json`, package smoke, migration tests and complete command suite.

<!-- tene:section:verdict -->
## Verdict

PASS. Package/checksum tamper/SBOM/route/update/uninstall/migration cases passed; full check/race/vet, Playwright 3/3, plugin and nine skill validators passed.
