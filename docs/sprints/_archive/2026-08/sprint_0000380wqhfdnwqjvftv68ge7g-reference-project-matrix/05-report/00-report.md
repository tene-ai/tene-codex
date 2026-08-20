---
schema_version: 1.0.0
document_type: report
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wqhfdnwqjvftv68ge7g
phase: report
status: draft
revision: 411
intent_ids: []
generated_at: 2026-08-20T01:15:02Z
generated_by: tene-workflow
---

# report — Reference Project Portability Matrix

<!-- tene:section:purpose -->
## Purpose

Record the completed three-project portability proof.

<!-- tene:section:scope -->
## Scope

Polyglot fallback, secret exclusion, layer classification, fixtures and matrix QA.

<!-- tene:section:layers -->
## Layers

Interface, Business Logic, Persistence and Infrastructure are all represented and verified.

<!-- tene:section:six-questions -->
## Six questions

`Analyze`, `sourcePaths`, `sourceExtension`, `classify` are defined in codeintel, called by graph commands/tests, accept repository roots/options, and return provenance-rich Report components with known or explicit unknown inputs, outputs, calls, references and effects.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380wqhjfnbjh3xxzpkhg1c`, AC `ac_0000380wqhjfq5y4p9jcbajkpm`, task `task_0000380wqkxry81kry1arvgd1r`, evidence `evidence_0000380wqs9etpy3mtjrpy4dp8`.

<!-- tene:section:decisions -->
## Decisions

Filesystem fallback is first-class but low-confidence; missing semantics are structured unknowns, never blank certainty.

<!-- tene:section:freeform -->
## Freeform

Retrospective: portability testing caught a security exclusion defect that Go-only tests could not expose.

<!-- tene:section:previous-sprints -->
## Previous sprints

Follows secret Sprint `sprint_0000380wpmvv19m87t6xa4xydr`; the `.tene` exclusion proves that boundary also holds in cross-language graph discovery.

<!-- tene:section:changed-files -->
## Changed files

Modified `internal/codeintel/codeintel.go` and tests; added mature monolith and polyglot fixtures; added Sprint/state/evidence artifacts.

<!-- tene:section:intent-fulfillment -->
## Intent fulfillment

Users can apply the same design/loop workflow to greenfield, legacy and service repositories. Unsupported languages remain visible without hallucinated relationships.

<!-- tene:section:qa-verdict -->
## Qa verdict

PASS: run `run_0000380wqpcqrrqsk8bjf4j3p4`, 3/3 cases; three repo types, mature 4/4 layers, polyglot 4 files/5 unknown categories, vault zero; full check/race and Playwright 3/3 passed.

<!-- tene:section:deferred-work -->
## Deferred work

None. More semantic providers are optional future enhancements, not required scope.

<!-- tene:section:next-sprint -->
## Next sprint

Release and Marketplace completion, including package contents, install/update/uninstall smoke, SBOM/checksums/license and migration closure.


<!-- tene:generated:summary:start -->
### Generated Sprint Summary

- Sprint: `sprint_0000380wqhfdnwqjvftv68ge7g`
- Previous sprints: `sprint_0000380wpmvv19m87t6xa4xydr`
- Intent IDs: `intent_0000380wqhjfnbjh3xxzpkhg1c`
- Tasks: 1
- QA verdict: `passed`
- Open gaps: 0
- State revision: 432

<!-- tene:generated:summary:end -->
