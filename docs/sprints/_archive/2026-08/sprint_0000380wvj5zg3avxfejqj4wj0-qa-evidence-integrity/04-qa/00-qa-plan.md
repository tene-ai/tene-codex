---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wvj5zg3avxfejqj4wj0
phase: qa
status: passed
revision: 1
intent_ids: [intent_0000380wvk7q507samchfwpvxw]
generated_at: 2026-08-20T01:50:11Z
generated_by: tene-workflow
---

# qa — QA Evidence Integrity and Independent Evaluation

<!-- tene:section:purpose -->
## Purpose

Prove that invalid QA evidence cannot produce a completed Sprint and valid evidence still can.

<!-- tene:section:scope -->
## Scope

Seven variants, evaluator mutations, CLI structured observation, native Go execution, regression and redaction/integrity checks.

<!-- tene:section:layers -->
## Layers

L1/L2/L3/L6/L7 required and evidenced. L4/L5 are explicitly not applicable because this Sprint modifies a local CLI evaluator without a deployed system or end-user UI; approver and rationale are recorded per case.

<!-- tene:section:six-questions -->
## Six questions

QA cases call adapter/observer ingestion, which creates evidence records consumed by the root-aware evaluator; JSON observations and native execution artifacts are the inputs and canonical case/run verdicts are the mutations.

<!-- tene:section:traceability -->
## Traceability

Run `run_0000380wwpw1tf7nv0wxs525qm`; AC `ac_0000380wvk7q5mkdwqj2q1c1gc`; every case has native execution plus structured evidence.

<!-- tene:section:decisions -->
## Decisions

No manual pass. Non-applicable layers require named approval and reason. Failed evidence remains preserved and forces failure.

<!-- tene:section:freeform -->
## Freeform

This Sprint tests the evaluator as a security boundary, not merely as a convenience command.

<!-- tene:section:environment -->
## Environment

Local fingerprint `local:1382c4208946a24a`, Go toolchain and repository state bound to QA plan revision 522 and specification hash `8b90ae3f...`.

<!-- tene:section:capabilities -->
## Capabilities

Go native test, schema validation, race detector, deterministic evaluator and structured external-observation import.

<!-- tene:section:charters -->
## Charters

Happy, alternate, empty, validation, permission, downstream failure and recovery variants. Mutation coverage includes wrong run/case/spec/revision, missing layers/references, failed assertion, absent tool identity, redaction failure, manual pass and changed artifact bytes.

<!-- tene:section:ux-data-flow -->
## Ux data flow

CLI request → QA case lookup → adapter/observation validation → evidence artifact/state → deterministic evaluator → stable blocker or derived pass. There is no UI transition in this component Sprint.

<!-- tene:section:evidence -->
## Evidence

Each case has a `codex-deterministic-test-observer` JSON artifact and an actual `go-test` execution artifact under `04-qa/evidence`; all hashes verify.

<!-- tene:section:verdict -->
## Verdict

Passed: all seven cases satisfy required layers and requirement references, and adversarial false-pass mutations are rejected.
