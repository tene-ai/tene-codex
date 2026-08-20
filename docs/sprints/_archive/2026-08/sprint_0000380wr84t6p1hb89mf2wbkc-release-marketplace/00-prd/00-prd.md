---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wr84t6p1hb89mf2wbkc
phase: prd
status: draft
revision: 437
intent_ids: []
generated_at: 2026-08-20T01:21:13Z
generated_by: tene-workflow
---

# prd — Release and Marketplace Completion

<!-- tene:section:purpose -->
## Purpose

Make 0.1.0 reproducibly packageable, installable and ready for public skills-only submission.

<!-- tene:section:scope -->
## Scope

Bundle verification, SBOM/checksums/provenance, smoke lifecycle, marketplace catalog/materials, CI and migration support closure.

<!-- tene:section:layers -->
## Layers

Interface docs/catalog; business release guards; persistence project-state preservation; infrastructure CI/artifacts.

<!-- tene:section:six-questions -->
## Six questions

`package-plugin.sh`, `release-smoke.sh`, `generate-sbom.py`, wrapper checksum path and CI workflows are defined in scripts/GitHub, called by make/tag release, accept version/stage, and return a verified bundle/evidence or nonzero failure without deleting project state.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380wr877dd8bq79gacg7kc`; ACs `ac_0000380wr877cqrzpmgq248dwc`, `ac_0000380wr8jmksakse7nx1mv4w`; WP-13 and release plan.

<!-- tene:section:decisions -->
## Decisions

Skills-only public submission; local marketplace catalog for testing; binaries verified before execution; official portal checklist frozen per release.

<!-- tene:section:freeform -->
## Freeform

OpenAI official manual refreshed 2026-08-20.

<!-- tene:section:problem -->
## Problem

Packaging existed but lacked runtime checksum verification, SBOM, provenance, clean update/uninstall smoke and complete portal materials.

<!-- tene:section:actors -->
## Actors

Maintainer, GitHub Actions, plugin reviewer, installer and project user.

<!-- tene:section:journeys -->
## Journeys

Check → package → validate/checksum/SBOM → install/routing → update → uninstall preserving state → tag/attest/release → portal submission.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

Both blocking AC observables pass; tampered/missing checksum fails; supported migration succeeds with backup; unsupported versions do not mutate; old debt closes with evidence.

<!-- tene:section:non-goals -->
## Non goals

Creating an MCP server, submitting on behalf of the organization without portal credentials, or tagging before final audit.
