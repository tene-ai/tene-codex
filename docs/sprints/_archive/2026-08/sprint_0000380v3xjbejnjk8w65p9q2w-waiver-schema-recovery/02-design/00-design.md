---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v3xjbejnjk8w65p9q2w
phase: design
status: draft
revision: 127
intent_ids: []
generated_at: 2026-08-19T17:43:57Z
generated_by: tene-workflow
---

# design — Waiver Schema Migration and Recovery

<!-- tene:section:purpose -->
## Purpose

Define fail-closed policy and storage recovery contracts.

<!-- tene:section:scope -->
## Scope

Additive schema 1.0 fields and commands plus explicit legacy 0.9→1.0 transform.

<!-- tene:section:layers -->
## Layers

App exposes commands, workflow evaluates active waivers, state validates/migrates/repairs, filesystem supplies atomic persistence.

<!-- tene:section:six-questions -->
## Six questions

`Waiver` is stored by ID and referenced from gaps. `MigrationPlan` is derived from raw project schema. `Migrate` and `RepairDerived` accept no arbitrary paths and return structured plan/backup/projection results.

<!-- tene:section:traceability -->
## Traceability

AC1 owns waiver fields/gate tests; AC2 owns migration APIs/fixtures; AC3 owns repair and validation tests.

<!-- tene:section:decisions -->
## Decisions

Waivable categories exclude security and evidence integrity. Expiry comparison uses UTC. Apply is lock-protected, preserves original bytes in backups, appends an event, and writes all projections.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:components -->
## Components

Domain Waiver; workflow ActiveWaiver; state PlanMigration/Migrate/RepairDerived/ValidateProject; app waiver/migrate/doctor routing.

<!-- tene:section:interfaces -->
## Interfaces

`waiver create --gap ID --reason TEXT --approver ID --expires RFC3339`; `waiver list`; `waiver revoke ID`; `migrate status|dry-run|apply`; `doctor [--repair]`.

<!-- tene:section:data -->
## Data

Waiver includes ID, sprint/gap, reason, scope, approver, timestamps, status and optional revocation. Migration plan includes from/to, required, supported and changes.

<!-- tene:section:state-transitions -->
## State transitions

Active→revoked or naturally expired. Legacy→backup→event append→1.0 projections. Repair verifies journal→validates project→rewrites active/master only.

<!-- tene:section:failures -->
## Failures

Missing gap, forbidden category, bad/expired time, unsupported schema, corrupt journal or invalid project returns stable failure without partial projection updates.

<!-- tene:section:security -->
## Security

No waiver for secret leak, evidence redaction/hash or corrupt journal. Reason/approver are metadata and scanned as normal state.

<!-- tene:section:tests -->
## Tests

Workflow time-bound waiver tests, state legacy/unsupported/repair tests, CLI lifecycle and regression suite.
