---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380xznb6mpm5bx0xpxex6w
phase: qa
status: active
revision: 961
intent_ids: [intent_0000380y046tkxxteemz74gg58]
generated_at: 2026-08-20T07:05:37Z
generated_by: tene-workflow
---

# qa — Durable Journal Compaction

<!-- tene:section:purpose -->
## Purpose

Prove durable compaction without trusting one happy-path assertion.

<!-- tene:section:scope -->
## Scope

Unit, race, static, reference/E2E, package, semantic audit, dogfood shrink, evidence and doctor.

<!-- tene:section:layers -->
## Layers

L1 build, L2 state units, L3 app, L4 filesystem/CLI, L5 compact-resume journey, L6 regression, L7 audit/doctor.

<!-- tene:section:six-questions -->
## Six questions

Observe compact/doctor callers, archive and active files, sequence/checksum inputs, replayed output and corruption errors.

<!-- tene:section:traceability -->
## Traceability

All variants bind to `ac_0000380y046tkpdm8ey9panj70` with spec/state hashes.

<!-- tene:section:decisions -->
## Decisions

No manual pass; any mismatch, tamper acceptance, race, release, audit or doctor failure blocks archive.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:environment -->
## Environment

Repository toolchain plus isolated temporary state fixtures and the real dogfood state.

<!-- tene:section:capabilities -->
## Capabilities

Go, Node/Playwright, Python, filesystem hashing, plugin validator and Codex probe.

<!-- tene:section:charters -->
## Charters

Happy shrink; alternate second compact; initialized-only boundary; malformed anchor; filesystem failure; tampered archive; failed compact preserves active journal.

<!-- tene:section:ux-data-flow -->
## Ux data flow

User compact → archive proof → bounded checkpoint → resume/replay → doctor; tampering → explicit blocker and restore guidance.

<!-- tene:section:evidence -->
## Evidence

Seven sanitized, run/case-bound observations and the bounded command transcript are under `04-qa/observations`. They record archive SHA-256 `6c0642b93e12794ebd47a44c42a889c657ba7c7f0c174fae68130147488dbeec`, 952 archived events, and active journal size change from 7,542,532 to 2,134,981 bytes.

<!-- tene:section:verdict -->
## Verdict

PASS. All seven variants and L1–L7 assertions passed. `make check`, full race tests, vet, reference flows, release smoke, contracts audit, evidence verification and doctor passed; doctor verified one archived event segment and matching project/active/master projections.

<!-- tene:generated:traceability:start -->
### Generated Traceability

- State revision: `961`
- Sprint: `sprint_0000380xznb6mpm5bx0xpxex6w`
- Intents: `intent_0000380y046tkxxteemz74gg58`
- Tasks: `task_0000380y0cqhsx7fjvbj0zwk94`, `task_0000380y0crx5qtznz6hsj9pw8`

<!-- tene:generated:traceability:end -->
