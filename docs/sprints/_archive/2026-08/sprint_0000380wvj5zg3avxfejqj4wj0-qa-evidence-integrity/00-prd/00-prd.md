---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wvj5zg3avxfejqj4wj0
phase: prd
status: confirmed
revision: 1
intent_ids: [intent_0000380wvk7q507samchfwpvxw]
generated_by: tene-workflow
---

# QA Evidence Integrity and Independent Evaluation

<!-- tene:section:purpose -->
## Purpose

Eliminate every path that can declare a blocking acceptance criterion complete without evidence proving the intended behavior.

<!-- tene:section:scope -->
## Scope

Bind evidence to run, case, specification hash and state revision; require layer and requirement coverage, passing assertions, tool/environment metadata, redaction and content integrity. Manual pass is forbidden. Automatic deterministic evaluation is authoritative. App Server orchestration remains outside this Sprint.

<!-- tene:section:layers -->
## Understanding Layers

Interface: `qa plan|execute|observe|case|evaluate`. Business Logic: evidence eligibility and gate evaluation. Persistence: QA run/evidence state and immutable artifacts. Infrastructure: native tools, Playwright and external browser observers.

<!-- tene:section:six-questions -->
## Six Questions

Names: `QASpecHash`, `EvaluateQAGateAtRoot`, evidence assertions. Definitions: workflow, app, domain and qaadapter packages. Imports: CLI runtime calls evaluator and adapters. Callers: QA commands and tests. Inputs: confirmed intent, AC, case, structured evidence. Outputs: deterministic findings and pass/fail state.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380wvk7q507samchfwpvxw`; AC `ac_0000380wvk7q5mkdwqj2q1c1gc`; closes FR-06, FR-07 and AC-PRODUCT-04/05 false-pass gaps.

<!-- tene:section:decisions -->
## Decisions

All seven layers start required. A future explicit disposition command may mark a layer not applicable only with approver and reason. An adapter exit code is evidence for its declared layers but never substitutes for observable/variant/expected/forbidden assertions.

<!-- tene:section:freeform -->
## Additional Perspective

The evaluator ignores builder completion prose and consumes only canonical intent, case contracts and registered artifacts.

<!-- tene:section:problem -->
## Problem

The former gate accepted any hash-bearing evidence linked to an AC and allowed manual `passed`, so unit-test success could falsely prove UX and data-flow intent.

<!-- tene:section:actors -->
## Actors

Project user, builder agent, isolated evaluator, QA adapter and external observer.

<!-- tene:section:journeys -->
## Journeys

Plan charters, execute native checks, import structured journey observations, evaluate all required layers and requirement references, then either block with stable findings or pass.

<!-- tene:section:acceptance-criteria -->
## Acceptance Criteria

Wrong-run, wrong-case, stale-spec, wrong-revision, missing-layer, failed-assertion, missing-tool metadata, redaction failure and tampered content all block. Manual pass blocks. Complete structured evidence passes.

<!-- tene:section:non-goals -->
## Non-goals

No LLM-only override and no remote QA service in this Sprint.
