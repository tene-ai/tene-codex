---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v3xjbejnjk8w65p9q2w
phase: prd
status: draft
revision: 127
intent_ids: []
generated_at: 2026-08-19T17:43:57Z
generated_by: tene-workflow
---

# prd — Waiver Schema Migration and Recovery

<!-- tene:section:purpose -->
## Purpose

Complete the carried policy and recovery contracts so exceptions cannot bypass gates silently and persisted state can be upgraded or repaired predictably.

<!-- tene:section:scope -->
## Scope

Waiver create/list/revoke with scope, approver, reason and expiry; gap-gate integration; schema migration status/dry-run/apply with backup; projection repair after journal verification; strict unknown-field and referential validation.

<!-- tene:section:layers -->
## Layers

Interface: waiver/migrate/doctor commands. Business Logic: waiver validity and migration planning. Persistence: project projection, journal, backup and derived projections. Infrastructure: atomic filesystem replacement and locking.

<!-- tene:section:six-questions -->
## Six questions

Every new command/type records its definition, caller, structured input/output and mutations. Repair never mutates the event journal; migration records a hash-chained event.

<!-- tene:section:traceability -->
## Traceability

Implements WP-03 waiver, WP-01 schema compatibility, WP-02 doctor repair and WP-14 migration/runbook carry items.

<!-- tene:section:decisions -->
## Decisions

Waivers require explicit user identity, reason and expiry and cannot cover security/evidence corruption. Migration supports only declared source versions and always writes a recoverable backup first.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:problem -->
## Problem

Current gaps can only be resolved, state schema has no upgrade command, and doctor can diagnose but not rebuild derived projections. These omissions make policy exceptions informal and upgrades operationally unsafe.

<!-- tene:section:actors -->
## Actors

Repository owner approving a bounded exception; maintainer upgrading state; Codex diagnosing and repairing projections.

<!-- tene:section:journeys -->
## Journeys

Approve a non-security gap temporarily and have the phase gate recognize only an active waiver; revoke it and restore blocking. Preview a migration, apply with backup, verify journal/projections, and repair derived files without changing source state.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

- AC1: waiver lifecycle validates gap, approver, reason and future expiry; expired/revoked/security waivers never unblock gates.
- AC2: migrate status/dry-run/apply supports declared legacy schema with deterministic changes, backup and event; unsupported versions fail closed.
- AC3: doctor repair regenerates active/master projections only after journal/project validation, and all regression/security/plugin gates pass.

<!-- tene:section:non-goals -->
## Non goals

Downgrade, cross-major migration, event-journal semantic replay from historical delta events, remote backups or automatic approval.
