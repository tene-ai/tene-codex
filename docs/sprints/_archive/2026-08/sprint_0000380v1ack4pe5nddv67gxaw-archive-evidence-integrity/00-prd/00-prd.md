---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v1ack4pe5nddv67gxaw
phase: prd
status: draft
revision: 59
intent_ids: []
generated_at: 2026-08-19T17:21:14Z
generated_by: tene-workflow
---

# prd — Archive evidence integrity stabilization

<!-- tene:section:purpose -->
## Purpose

Preserve report and evidence integrity when a completed Sprint moves to its archive location.

<!-- tene:section:scope -->
## Scope

Update archive locators, write an archive manifest, verify evidence without an active Sprint, include evidence in doctor, and add regression tests.

<!-- tene:section:layers -->
## Layers

Interface: archive/evidence/doctor CLI. Business: archive invariants. Persistence: moved paths and manifest. Infrastructure: post-archive verification.

<!-- tene:section:six-questions -->
## Six questions

The changed names are `runtime.transition` and `invalidEvidence` in `internal/app/app.go`; they are called by archive, evidence verify and doctor, receive project/path state, and return locator mutations or findings.

<!-- tene:section:traceability -->
## Traceability

This Sprint follows the MVP Sprint and addresses a real archive-path gap found after its QA.

<!-- tene:section:decisions -->
## Decisions

Archive changes locators in the same canonical state mutation and writes `99-archive/archive-manifest.json`. Read-only evidence verification works without an active Sprint.

<!-- tene:section:freeform -->
## Freeform

No migration of older external projects is required in this pre-alpha repository.

<!-- tene:section:problem -->
## Problem

Moving a Sprint document tree could leave evidence URI, report path and graph locator values pointing to the pre-archive path, preventing future regression verification.

<!-- tene:section:actors -->
## Actors

Developers resuming archived work and auditors verifying historical evidence.

<!-- tene:section:journeys -->
## Journeys

Archive a passed Sprint, load its immutable metadata, verify evidence hashes with no active Sprint, and run doctor successfully.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

1. Archive relocates report/evidence/graph locators and produces a manifest.
2. `evidence verify` and `doctor` validate archived evidence when no Sprint is active.

<!-- tene:section:non-goals -->
## Non goals

No remote evidence store, retention scheduler or old-version migration is included.
