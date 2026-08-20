---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wvj5zg3avxfejqj4wj0
phase: design
status: approved
revision: 1
intent_ids: [intent_0000380wvk7q507samchfwpvxw]
generated_at: 2026-08-20T01:50:11Z
generated_by: tene-workflow
---

# design — QA Evidence Integrity and Independent Evaluation

<!-- tene:section:purpose -->
## Purpose

Specify the executable evidence eligibility contract that prevents false completion.

<!-- tene:section:scope -->
## Scope

QA run fingerprinting, evidence assertions/layers/tool identity, observation validation, case aggregation and immutable artifact verification.

<!-- tene:section:layers -->
## Layers

Interface: QA CLI. Business Logic: evaluator. Persistence: project state and evidence files. Infrastructure: native/browser adapters.

<!-- tene:section:six-questions -->
## Six questions

Definitions live in `internal/domain`; app imports adapter output; `workflow.EvaluateQAGateAtRoot` is called by `qa evaluate`; inputs are run/case/evidence; output is findings and derived case/run status.

<!-- tene:section:traceability -->
## Traceability

Every evidence record points to run, case, AC, spec hash and QA-plan revision. Assertions point to layer and requirement references.

<!-- tene:section:decisions -->
## Decisions

Manual pass is invalid. Case status is derived during evaluation. Content bytes must still match registered size and SHA-256.

<!-- tene:section:freeform -->
## Freeform

`not-applicable` is valid only as a structured approver-and-reason disposition; the initial planner requires all layers.

<!-- tene:section:components -->
## Components

`domain.EvidenceAssertion`, extended `Evidence/QARun`, qaadapter structured observation, spec hasher and root-aware evaluator.

<!-- tene:section:interfaces -->
## Interfaces

Observation JSON requires `spec_hash`, `state_revision`, `layers`, `tool_version` and assertions with `layer` plus `requirement_refs`.

<!-- tene:section:data -->
## Data

Required references are `observable`, `variant:<name>`, every `expected:<index>` and `forbidden:<index>`. Passing assertions collectively cover required layers and references.

<!-- tene:section:state-transitions -->
## State transitions

planned → evidenced while artifacts arrive; evaluation derives every case as passed/failed and then derives the run. There is no manual transition to passed.

<!-- tene:section:failures -->
## Failures

Stale spec, mismatched identity, absent metadata, failed assertion, missing layer/reference, redaction failure or changed bytes produce blockers.

<!-- tene:section:security -->
## Security

Existing secret scanning remains before persistence; evaluator also requires redaction status and artifact hash integrity.

<!-- tene:section:tests -->
## Tests

Mutation table for wrong run/case/spec/layer/assertion/tool/redaction plus a complete seven-layer structured-observation CLI journey.
