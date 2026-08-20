---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wpmvv19m87t6xa4xydr
phase: prd
status: draft
revision: 377
intent_ids: []
generated_at: 2026-08-20T01:07:13Z
generated_by: tene-workflow
---

# prd — Secret Boundary and Adversarial QA

<!-- tene:section:purpose -->
## Purpose

Complete WP-11 and false-pass defenses so secret values never enter Codex-visible state or accepted QA evidence.

<!-- tene:section:scope -->
## Scope

Names-only metadata, command policy, child-output quarantine, shared evidence leak detection, pre/post tool hooks, adversarial tests and gates.

<!-- tene:section:layers -->
## Layers

Interface: secret CLI/hooks. Business: allow/deny/redact/quarantine. Persistence: evidence hash/redaction validation. Infrastructure: fake tene and CI.

<!-- tene:section:six-questions -->
## Six questions

`secret.Run`, `ListNames`, `DetectLeak`, hook `pre_tool/post_tool` are defined in `internal/secret` and `hooks`, referenced by app/evidence/QA and hook manifest, called for secret execution or tool boundaries, accept env/argv or tool payload, and return names-only/sanitized results or fail closed without plaintext.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380wpyxj67yt4xvdz94fcr`; ACs `ac_0000380wpyxj6wp3vc3x9bjn88`, `ac_0000380wpzdhvb233dyvmvz19w`; FR-11/WP-11/AC-PRODUCT-06.

<!-- tene:section:decisions -->
## Decisions

Preview is treated as secret. Unknown leakage cannot be safely redacted, so matching output is discarded and marked quarantined. Secret-looking argv is rejected before process start.

<!-- tene:section:freeform -->
## Freeform

Tene remains a separate security product; this plugin consumes only metadata and child execution surfaces.

<!-- tene:section:problem -->
## Problem

The former adapter forwarded preview fields and could retain an unlabeled canary in child output.

<!-- tene:section:actors -->
## Actors

Project user, Codex, tene CLI, child test process and independent QA evaluator.

<!-- tene:section:journeys -->
## Journeys

List returns sorted names only. Run validates argv, invokes tene without a shell, captures output, quarantines leaks, preserves safe exit status. Hooks block unsafe inputs and flag leaked outputs. Evidence refuses leaked/tampered artifacts.

<!-- tene:section:acceptance-criteria -->
## Acceptance criteria

Preview/value never returned; env dump/shell/sensitive args denied; canary/token output quarantined; missing/permission/child errors safe; hook coverage includes `.env`, literals and post-output; poisoned evidence cannot pass QA/doctor.

<!-- tene:section:non-goals -->
## Non goals

Vault cryptography, secret creation, plaintext retrieval, production authorization or replacing tene permissions.
