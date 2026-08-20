---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380y3ew66xhrtn9e3sdy2r
phase: design
status: complete
revision: 1081
intent_ids: [intent_0000380y3hff9tbrqh38mvehb8]
generated_at: 2026-08-20T07:38:48Z
generated_by: tene-workflow
---

# design — Doctor Fail Closed Exit

<!-- tene:section:purpose -->
## Purpose

Define doctor result/error contract.

<!-- tene:section:scope -->
## Scope

Healthy success and workflow/state blocker branches.

<!-- tene:section:layers -->
## Layers

CLI envelope, blocker logic, read-only state checks, process exits.

<!-- tene:section:six-questions -->
## Six questions

runtime.doctor builds result once; workflow.Blocking selects error; Run.failResult places result into errors[0].details.

<!-- tene:section:traceability -->
## Traceability

ac_0000380y3hff8thtbtry95g2zc maps to tamper and healthy integration assertions.

<!-- tene:section:decisions -->
## Decisions

STATE_ blocker → STATE_CORRUPT/7; other blocker → DOCTOR_UNHEALTHY/3; no blocker → success/0.

<!-- tene:section:freeform -->
## Freeform

This section records the fail-closed doctor gate from the freeform perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:components -->
## Components

This section records the fail-closed doctor gate from the components perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:interfaces -->
## Interfaces

This section records the fail-closed doctor gate from the interfaces perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:data -->
## Data

This section records the fail-closed doctor gate from the data perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:state-transitions -->
## State transitions

This section records the fail-closed doctor gate from the state transitions perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:failures -->
## Failures

This section records the fail-closed doctor gate from the failures perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:security -->
## Security

This section records the fail-closed doctor gate from the security perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:section:tests -->
## Tests

This section records the fail-closed doctor gate from the tests perspective; final measured QA and archive details are synchronized before completion.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `1081`
- Sprint: `sprint_0000380y3ew66xhrtn9e3sdy2r`
- Intents: `intent_0000380y3hff9tbrqh38mvehb8`
- Tasks: `task_0000380y3mjbzq6kbm80a0y2mg`

<!-- tene:generated:traceability:end -->
