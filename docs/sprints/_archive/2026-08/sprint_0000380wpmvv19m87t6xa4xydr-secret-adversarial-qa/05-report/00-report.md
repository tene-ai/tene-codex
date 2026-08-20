---
schema_version: 1.0.0
document_type: report
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380wpmvv19m87t6xa4xydr
phase: report
status: draft
revision: 377
intent_ids: []
generated_at: 2026-08-20T01:07:13Z
generated_by: tene-workflow
---

# report — Secret Boundary and Adversarial QA

<!-- tene:section:purpose -->
## Purpose

Record completion of the tene value boundary and adversarial false-pass defenses.

<!-- tene:section:scope -->
## Scope

Adapter, detector, evidence integration, hooks, attack tests and full QA.

<!-- tene:section:layers -->
## Layers

Interface: secret commands/hooks. Business: validation and quarantine. Persistence: evidence pre-mutation rejection and hashes. Infrastructure: fake tene/CI.

<!-- tene:section:six-questions -->
## Six questions

`Run`, `ListNames`, `DetectLeak`, `pre_tool`, `post_tool` are defined in `internal/secret` and `hooks`, imported/referenced by app and manifest, called at secret/tool/evidence boundaries, accept env/argv/bytes or hook JSON, and return names-only/sanitized results or safe errors without values.

<!-- tene:section:traceability -->
## Traceability

Intent `intent_0000380wpyxj67yt4xvdz94fcr`; ACs `ac_0000380wpyxj6wp3vc3x9bjn88`, `ac_0000380wpzdhvb233dyvmvz19w`; tasks `task_0000380wq20q15fpwfk6w4sp6g`, `task_0000380wq21za0gg7m8zdqp5jw`; evidence `evidence_0000380wq991kbxmvz4783fywm`.

<!-- tene:section:decisions -->
## Decisions

Preview is never returned; matched output is quarantined rather than guessed-redacted; hooks augment but do not replace core enforcement.

<!-- tene:section:freeform -->
## Freeform

Retrospective: a type named metadata was insufficient—the outgoing DTO needed an explicit names-only projection.

<!-- tene:section:previous-sprints -->
## Previous sprints

Builds on routing Sprint `sprint_0000380wney2n4hv4wct0skz4w`: natural secret requests select `$tene-secrets`, while this Sprint makes the invoked boundary safe and evidence-honest.

<!-- tene:section:changed-files -->
## Changed files

Reworked `internal/secret/runner.go`; added `runner_test.go`; connected detector in `internal/app/app.go`; added app poison tests; expanded `hooks/tene_hook.py`, `hooks/hooks.json`, `tests/hooks_test.py`; added Sprint/state/evidence artifacts.

<!-- tene:section:intent-fulfillment -->
## Intent fulfillment

Codex can use tene names and child execution without receiving preview/value output. Unsafe commands/literals fail before execution; leaked outputs are discarded; poisoned evidence never mutates accepted state.

<!-- tene:section:qa-verdict -->
## Qa verdict

PASS: run `run_0000380wq6m550ya9bcbyactpw`, 6/6 cases. Adapter attacks 6, poison patterns 3, hooks 5, accepted plaintext 0, previews 0. Check/race/vet, Playwright 3/3, plugin and nine skill validators, evidence verification passed.

<!-- tene:section:deferred-work -->
## Deferred work

None. Production authorization remains deliberately outside plugin scope, not deferred implementation.

<!-- tene:section:next-sprint -->
## Next sprint

Reference project matrix: prove the same public workflow in mature monolith and polyglot/service repositories, including degraded graph providers.


<!-- tene:generated:summary:start -->
### Generated Sprint Summary

- Sprint: `sprint_0000380wpmvv19m87t6xa4xydr`
- Previous sprints: `sprint_0000380wney2n4hv4wct0skz4w`
- Intent IDs: `intent_0000380wpyxj67yt4xvdz94fcr`
- Tasks: 2
- QA verdict: `passed`
- Open gaps: 0
- State revision: 406

<!-- tene:generated:summary:end -->
