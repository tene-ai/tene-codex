---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wrs5bseqwkpkg84mdwg
phase: prd
status: draft
revision: 472
intent_ids: []
generated_at: 2026-08-20T01:25:52Z
generated_by: tene-workflow
---

# prd — Final Requirements Traceability Audit

<!-- tene:section:purpose -->
## Purpose

Prove complete implementation coverage against authored FR/AC/WP contracts and repair residual contradictions.

<!-- tene:section:scope -->
## Scope

Machine traceability manifest/auditor, Codex capability probes, all state debt/tasks, documentation accuracy and final gates.

<!-- tene:section:layers -->
## Layers

Interface audit/doctor; business coverage rules; persistence state/debt scan; infrastructure CI/full suite.

<!-- tene:section:six-questions -->
## Six questions

`requirements-audit.py` and `ProbeCodex` are defined in scripts/projectconfig, called by make/doctor/final audit, accept repository/state, and return structured coverage, missing locators, unfinished work and live Codex/plugin capabilities.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380wrs7a7t859epzgtpes4`; AC `ac_0000380wrs7a76scmce852d2ar`; all 11 FRs, 8 product ACs and 14 WPs.

<!-- tene:section:decisions -->
## Decisions

Required scope is sourced from PRD/plan/design. Optional post-MVP remote services remain explicit non-goals; capability-dependent Codex surfaces are probed, not assumed.

<!-- tene:section:freeform -->
## Freeform

Completion claim is deferred until this Sprint itself is archived and `--final` passes.

<!-- tene:section:problem -->
## Problem

Many passing Sprints do not alone prove global coverage or absence of stale debt and documentation contradictions.

<!-- tene:section:actors -->
## Actors

User/reviewer, Codex, deterministic auditor and release maintainer.

<!-- tene:section:journeys -->
## Journeys

Extract IDs → map locators → validate existence → scan state → probe host → repair gaps → full QA → archive → final inactive-state audit.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

11/11 FR, 8/8 product AC, 14/14 WP; missing locators zero; open/deferred gaps zero; unfinished tasks zero; final active Sprint empty; doctor capabilities and all gates pass.

<!-- tene:section:non-goals -->
## Non goals

Remote MCP/App Server product implementation, public portal action, tagging without user release decision, or claiming optional future roadmap as MVP.
