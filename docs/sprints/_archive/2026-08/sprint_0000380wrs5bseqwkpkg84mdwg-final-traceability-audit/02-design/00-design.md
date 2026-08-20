---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wrs5bseqwkpkg84mdwg
phase: design
status: draft
revision: 472
intent_ids: []
generated_at: 2026-08-20T01:25:52Z
generated_by: tene-workflow
---

# design — Final Requirements Traceability Audit

<!-- tene:section:purpose -->
## Purpose

Design an executable global completeness proof.

<!-- tene:section:scope -->
## Scope

JSON locator map, state policy, audit exit contract and capability snapshot.

<!-- tene:section:layers -->
## Layers

CLI/script output; identifier/range rules; project JSON scan; CI/doctor probes.

<!-- tene:section:six-questions -->
## Six questions

`requirements-audit.py` lives in scripts and is referenced by Makefile; it accepts root/final, reads trace manifest/state, and returns counts/lists/pass. `ProbeCodex` lives in projectconfig, imported by app doctor, probes CLI/version/app-server plus local manifest/skills/hooks/agents/MCP and returns safe DTO.

<!-- tene:section:traceability -->
## Traceability

All identifier sets have exact generated ranges and one or more existing locators.

<!-- tene:section:decisions -->
## Decisions

Filesystem locators prove implementation/test presence; full runtime gates provide behavior evidence. Missing optional MCP config is reported false, not unhealthy.

<!-- tene:section:freeform -->
## Freeform

Probe errors are diagnostic and capability absence does not break core portability.

<!-- tene:section:components -->
## Components

Trace manifest, auditor, Codex probe, Makefile gate and final archived report/evidence.

<!-- tene:section:interfaces -->
## Interfaces

`requirements-audit.py [--root PATH] [--final]`; doctor JSON `capabilities.codex`.

<!-- tene:section:data -->
## Data

Three maps (FR/AC/WP) to locator arrays; result counts, missing, gaps, unfinished, active ID and pass.

<!-- tene:section:state-transitions -->
## State transitions

Auditor read-only. Final flag passes only after archive clears active pointer.

<!-- tene:section:failures -->
## Failures

Missing ID/locator, open/deferred gap, unfinished task or final active Sprint exits 1. Probe timeouts are reported without hanging.

<!-- tene:section:security -->
## Security

No content/secret scanning beyond named repository files; `.tene` remains excluded; probe captures version only.

<!-- tene:section:tests -->
## Tests

Make audit plus full product gate; post-archive `--final` is terminal proof.
