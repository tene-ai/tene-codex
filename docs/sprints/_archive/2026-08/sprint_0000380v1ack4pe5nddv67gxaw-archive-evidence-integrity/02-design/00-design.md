---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v1ack4pe5nddv67gxaw
phase: design
status: draft
revision: 59
intent_ids: []
generated_at: 2026-08-19T17:21:14Z
generated_by: tene-workflow
---

# design — Archive evidence integrity stabilization

<!-- tene:section:purpose -->
## Purpose

Specify archive relocation and historical evidence verification changes.

<!-- tene:section:scope -->
## Scope

Modify only the application use cases and full Sprint integration fixture.

<!-- tene:section:layers -->
## Layers

Interface: unchanged command names. Business: evidence validation helper. Persistence: locator and manifest mutation. Infrastructure: doctor after archive.

<!-- tene:section:six-questions -->
## Six questions

`runtime.transition` is defined and used in `internal/app/app.go` by phase/sprint archive; input is target phase and state, output is moved documents plus updated projection. `invalidEvidence` is defined in the same file and called by evidence verify/doctor; input is root/project, output is sorted invalid IDs without mutation.

<!-- tene:section:traceability -->
## Traceability

Archive AC maps to transition and integration assertions; verification AC maps to shared helper and post-archive CLI checks.

<!-- tene:section:decisions -->
## Decisions

Rewrite only locators prefixed by the archived Sprint root. Preserve hashes because content does not change. Manifest is written after move but before state commit, with rollback on failure.

<!-- tene:section:freeform -->
## Freeform

The archived Sprint's report/evidence remains readable through updated relative paths.

<!-- tene:section:components -->
## Components

Archive branch in `transition`, `invalidEvidence`, evidence command, doctor and `TestCLICompleteSprintArchivesDocuments`.

<!-- tene:section:interfaces -->
## Interfaces

No breaking CLI change. `evidence list [sprint-id]` and `evidence verify` now work without an active Sprint.

<!-- tene:section:data -->
## Data

Update `Sprint.DocumentRoot`, `Sprint.ReportPath`, matching `Evidence.URI` and graph node locator. Write schema version, Sprint ID, archive time, revision, QA run/status and report path to manifest.

<!-- tene:section:state-transitions -->
## State transitions

Report→archived moves the tree, creates manifest, commits locator updates and clears active Sprint; failure restores the old path.

<!-- tene:section:failures -->
## Failures

Move, manifest marshal/write or state mutation return an error. Doctor emits blocking `QA_EVIDENCE_INVALID` findings for missing, changed or unsafe artifacts.

<!-- tene:section:security -->
## Security

Historical evidence is rescanned for secret patterns as well as hash changes.

<!-- tene:section:tests -->
## Tests

Archive integration places evidence inside the Sprint tree, asserts manifest and relocated URI/file, then the real dogfood Sprint runs evidence verify and doctor after archive.
