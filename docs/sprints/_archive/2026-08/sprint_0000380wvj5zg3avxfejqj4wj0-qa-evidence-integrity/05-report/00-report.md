---
schema_version: 1.0.0
document_type: report
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wvj5zg3avxfejqj4wj0
phase: report
status: complete
revision: 1
intent_ids: [intent_0000380wvk7q507samchfwpvxw]
generated_at: 2026-08-20T01:50:11Z
generated_by: tene-workflow
---

# report — QA Evidence Integrity and Independent Evaluation

<!-- tene:section:purpose -->
## Purpose

Record how the QA false-pass boundary was closed and proven.

<!-- tene:section:scope -->
## Scope

Evidence/run schemas, ingestion, evaluator, layer disposition, tests and operator documentation.

<!-- tene:section:layers -->
## Layers

Interface: strengthened QA commands. Business Logic: independent root-aware gate. Persistence: bound evidence and QA state. Infrastructure: native and external observation adapters.

<!-- tene:section:six-questions -->
## Six questions

`EvidenceAssertion`, `QASpecHash`, `EvaluateQAGateAtRoot`, `ReadObservation`, and QA CLI handlers are defined in domain/workflow/qaadapter/app; invoked by `qa plan|execute|observe|evaluate`; accept canonical intent, cases and artifacts; return evidence records, findings and derived verdicts.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380wvk7q507samchfwpvxw`, AC `ac_0000380wvk7q5mkdwqj2q1c1gc`, tasks `task_0000380ww94ys4jxvcmvvcsmj4` and `task_0000380ww96c0s257b3msrw78r`.

<!-- tene:section:decisions -->
## Decisions

Case pass cannot be manually assigned. Content integrity and semantic coverage are both mandatory. Non-applicable layers are explicit decisions.

<!-- tene:section:freeform -->
## Freeform

This corrects the earlier audit's mistaken assumption that a linked, redaction-safe hash necessarily proves an AC.

<!-- tene:section:previous-sprints -->
## Previous sprints

Extends final-traceability-audit by replacing its shallow evidence semantics; preserves all archived state and artifacts.

<!-- tene:section:changed-files -->
## Changed files

`internal/domain/types.go`: evidence/run contract. `internal/workflow/workflow.go`: spec hash and gate. `internal/app/app.go`: ingestion/disposition/manual-pass behavior. `internal/qaadapter/*`: structured observations. `schemas/*`: wire contracts. Tests, QA reference and Sprint documents were updated.

<!-- tene:section:intent-fulfillment -->
## Intent fulfillment

Wrong identity, stale specification, incomplete layers/references, failed assertions, missing provenance, redaction failure, content tampering and manual pass all block deterministically.

<!-- tene:section:qa-verdict -->
## Qa verdict

Passed run `run_0000380wwpw1tf7nv0wxs525qm`; seven variants, actual Go execution artifacts, structured checkpoint/assertion artifacts, mutation tests, race, vet and E2E regression.

<!-- tene:section:deferred-work -->
## Deferred work

No item in this Sprint scope is deferred. Automatic spec/code gap generation remains the next audited blocking capability.

<!-- tene:section:next-sprint -->
## Next sprint

Implement semantic bidirectional Loop Check and automatic gap lifecycle.


<!-- tene:generated:summary:start -->
### Generated Sprint Summary

- Sprint: `sprint_0000380wvj5zg3avxfejqj4wj0`
- Previous sprints: `sprint_0000380wrs5bseqwkpkg84mdwg`
- Intent IDs: `intent_0000380wvk7q507samchfwpvxw`
- Tasks: 2
- QA verdict: `passed`
- Open gaps: 0
- State revision: 555

<!-- tene:generated:summary:end -->
