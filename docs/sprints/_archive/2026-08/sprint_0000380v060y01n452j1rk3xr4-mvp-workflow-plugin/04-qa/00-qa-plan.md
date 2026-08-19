---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v060y01n452j1rk3xr4
phase: qa
status: draft
revision: 1
intent_ids: []
generated_at: 2026-08-19T17:11:18Z
generated_by: tene-workflow
---

# qa — MVP workflow engine and Codex plugin

<!-- tene:section:purpose -->
## Purpose

Verify the five blocking acceptance criteria with executable core, plugin, hook, security and distribution evidence.

<!-- tene:section:scope -->
## Scope

Run Go unit/integration/race tests, vet, Python hook tests, JSON checks, official plugin and skill validators, full Sprint archive fixture, package cross-build, file-format inspection and checksum verification.

<!-- tene:section:layers -->
## Layers

- Interface: CLI envelopes, skill metadata, hook decisions and generated agents.
- Business Logic: transition guards, content/task trace gates and all-case QA evaluator.
- Persistence: event hash chain, stale revision, compact/clear, document archive and evidence hash.
- Infrastructure: argv-only tene runner, CI definitions and four platform artifacts.

<!-- tene:section:six-questions -->
## Six questions

QA calls the public commands and hooks identified in Design. Inputs are CLI argv, temporary repositories, fixture domain state and Codex hook JSON. Outputs are exit codes, JSON envelopes, filesystem state, QA findings, validator results, platform binaries and SHA-256 verification.

<!-- tene:section:traceability -->
## Traceability

Each of 15 generated QA cases maps to one of five blocking AC IDs. The shared verification artifact is registered against all five IDs and attached to each case because it contains the complete test matrix.

<!-- tene:section:decisions -->
## Decisions

Use deterministic automated evidence for this CLI/plugin Sprint. Browser UX is not applicable because no browser UI is shipped; plugin interaction is verified at manifest, skill and hook contracts.

<!-- tene:section:freeform -->
## Freeform

ARM and Linux cross-built binaries cannot execute on the Intel macOS host; they are verified by successful Go cross-build, platform file headers and checksums. The host darwin-amd64 binary is executed.

<!-- tene:section:environment -->
## Environment

macOS x86_64 development host, Go 1.26.6 compiling a Go 1.24 module, Python 3.14.6, official local plugin/skill validator scripts, no production network or credentials.

<!-- tene:section:capabilities -->
## Capabilities

Required: Go compiler/race detector/vet, Python stdlib unittest, temporary PyYAML validator environment, shell packaging. Browser, Chrome, DB and external API are not applicable to this repository-only deliverable.

<!-- tene:section:charters -->
## Charters

For every AC, generated happy/error/recovery cases exercise success, invalid input or missing evidence, and recovery/retry behavior. Core integration includes the complete PRD→archive journey; hook tests cover deny/allow/unknown actions.

<!-- tene:section:ux-data-flow -->
## Ux data flow

CLI request → application parser → workflow/domain rule → atomic journal/projection/document mutation → JSON response is tested. Codex event → Python hook → optional status/compact CLI → supported hook JSON is tested. Secret command → policy preflight → argv-based tene child boundary is code-reviewed and hook-denial tested without live secrets.

<!-- tene:section:evidence -->
## Evidence

`evidence/mvp-verification.txt` records the verified commands and outcomes. The CLI stores its content hash and redaction status; the QA run links it to all cases and criteria.

<!-- tene:section:verdict -->
## Verdict

Passed. Evidence `evidence_0000380v0rhnsd6kfj7mghbvmg` was registered against all five blocking criteria, all 15 happy/error/recovery cases passed with that evidence, SHA-256 and redaction verification passed, and `qa evaluate` returned no findings at state revision 54. Graph rebuild and invariant validation also passed at revision 55.
