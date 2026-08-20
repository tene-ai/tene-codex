---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wr84t6p1hb89mf2wbkc
phase: design
status: draft
revision: 437
intent_ids: []
generated_at: 2026-08-20T01:21:13Z
generated_by: tene-workflow
---

# design — Release and Marketplace Completion

<!-- tene:section:purpose -->
## Purpose

Specify trustworthy artifact assembly and lifecycle smoke.

<!-- tene:section:scope -->
## Scope

Package layout, checksums, SBOM, catalog, smoke, CI provenance and migration guard.

<!-- tene:section:layers -->
## Layers

Interface manifest/docs; business checksum/migration rules; persistence uninstall preservation; infrastructure cross-build/release.

<!-- tene:section:six-questions -->
## Six questions

Scripts live under `scripts`, are referenced by Makefile/workflows, called with version/stage, read source/binaries, and return self-contained stage/checksums/SBOM. Wrapper reads its own checksum entry before exec. Migration functions read schema header and mutate only supported apply.

<!-- tene:section:traceability -->
## Traceability

WP-13, design 10, two Sprint ACs and inherited gap.

<!-- tene:section:decisions -->
## Decisions

SPDX 2.3 deterministic JSON; SHA-256 per file and zip; GitHub OIDC attestation; local catalog path `./`; update is replacement; uninstall must leave `.tene-workflow`.

<!-- tene:section:freeform -->
## Freeform

README links to public submission requirements; no claim that portal submission is automated.

<!-- tene:section:components -->
## Components

Packager, SBOM generator, wrapper verifier, smoke harness, marketplace catalog, CI/release workflows and migration planner.

<!-- tene:section:interfaces -->
## Interfaces

`package-plugin.sh VERSION [STAGE]`, `release-smoke.sh`, wrapper CLI, tag `v*`, portal checklist.

<!-- tene:section:data -->
## Data

Manifest, skills/hooks/references/legal docs, four binaries, SPDX JSON, file/zip checksums, attestation and release notes.

<!-- tene:section:state-transitions -->
## State transitions

Install/update/uninstall affect plugin files only. Project workflow state persists. Supported migration writes backup+journal+projections; unsupported returns before write.

<!-- tene:section:failures -->
## Failures

Existing stage, missing checksum, mismatch, unsupported platform/schema, invalid JSON, routing smoke or state loss all fail nonzero.

<!-- tene:section:security -->
## Security

No network download during skill invocation; checksums precede binary exec; provenance/SBOM; Apache/NOTICE and security contact included.

<!-- tene:section:tests -->
## Tests

Release smoke plus full test pyramid and unsupported migration byte-equality.
