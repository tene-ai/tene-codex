---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v2b2q0668t7tajkj3m0
phase: qa
status: draft
revision: 88
intent_ids: []
generated_at: 2026-08-19T17:30:10Z
generated_by: tene-workflow
---

# qa — Code Intelligence and Intent QA Adapters

<!-- tene:section:purpose -->
## Purpose

Prove the code-intelligence and QA-adapter slice against observable intent, not implementation claims.

<!-- tene:section:scope -->
## Scope

Static/unit/integration, CLI journey, evidence integrity, security boundaries and documentation/plugin validation. No product UI exists, so L4 browser execution is not applicable; the external observation contract is tested with structured fixtures.

<!-- tene:section:layers -->
## Layers

L1 gofmt/vet/schema; L2 package/race tests; L3 CLI/state/evidence integration; L4 not applicable (no UI runtime); L5 CLI actor journeys and understandable output; L6 outside-path, unknown adapter, failed assertion and secret defenses; L7 all prior regression tests and doctor.

<!-- tene:section:six-questions -->
## Six questions

QA inspects public analyzer/adapter names and locators, app imports and call sites, root/path/case/run inputs, Report/ExecutionResult/Observation outputs, and graph/evidence state effects.

<!-- tene:section:traceability -->
## Traceability

AC1: codeintel unit + CLI understand/build. AC2: qaadapter unit + CLI capabilities/execute/observe. AC3: make/race/vet/hook/schema/plugin/skill/evidence/doctor evidence.

<!-- tene:section:decisions -->
## Decisions

One shared deterministic verification artifact may cover multiple variants only when every generated case is explicitly linked and the artifact directly contains the relevant assertions. External observation fixtures cannot attest a real browser session and are labeled contract evidence.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:environment -->
## Environment

Local macOS arm64 workspace, Go toolchain, Python hook tests, no `.codegraph/`, no Playwright configuration, tene CLI capability reported independently.

<!-- tene:section:capabilities -->
## Capabilities

Go native test available. Filesystem/static analysis available. External browser observation import available. CodeGraph and Playwright unavailable in this repository and explicitly reported rather than simulated.

<!-- tene:section:charters -->
## Charters

Happy: analyze code and execute verified native tests. Error: reject path traversal, dangling task references, unknown adapters and invalid observations. Recovery: repair task links and re-evaluate. Regression: run entire existing suite and post-evidence doctor.

<!-- tene:section:ux-data-flow -->
## Ux data flow

CLI request → app validation → analyzer/adapter → graph or artifact → state journal/projection → evidence hash → QA case → independent gate output. Confirm success and failure feedback at each boundary.

<!-- tene:section:evidence -->
## Evidence

Three `go-test` adapter executions are hash-bound to the three blocking ACs. A schema-valid external observation fixture proves checkpoint/assertion import and case linkage. `make check`, `go test -race ./...`, `go vet ./...`, all JSON parse checks, the official plugin validator and all nine official skill validations passed. Evidence verification reports six valid artifacts.

<!-- tene:section:verdict -->
## Verdict

Passed. All nine happy/error/recovery cases have criterion-linked, redaction-safe evidence and `qa evaluate` returned no findings. CodeGraph and Playwright are correctly marked not applicable for this repository environment, not reported as executed.
