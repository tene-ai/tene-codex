---
schema_version: 1.0.0
document_type: plan
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380ydgqkg6h9mm5ctvm9b4
phase: plan
status: ready
revision: 1102
intent_ids: [intent_0000380ydppfay2szegwvaqsqm, intent_0000380ydpq9tm54tgwpqnfd30, intent_0000380yeckc0mm5dsr9j59h6m, intent_0000380yecmcwb7v4ttjm9bbrm]
generated_at: 2026-08-20T09:06:41Z
generated_by: tene-workflow
---

# plan — Plugin Identity and Installed Runtime Verification

<!-- tene:section:purpose -->
## 목적

Plugin identity migration과 실제 설치본 검증을 하나의 vertical release slice로 수행한다. 먼저 canonical identifiers를 바꾸고, router·문서·eval을 같은 계약으로 정렬한 뒤, staged/local-installed artifact에서 검증한다. 이 과정에서 발견된 deprecated-intent criterion pollution도 QA false-pass/false-block을 막기 위해 함께 제거한다.

<!-- tene:section:scope -->
## 범위

- In: manifest/marketplace identity, 9 skill directory/frontmatter/agent metadata, router IDs, explicit syntax, docs/evals/package/install smoke, versioned local reinstall, deprecated intent filtering.
- Out: repository/module rename, public submission, remote integrations, V2 이후 core semantic audit items.
- Stop condition: 두 blocking AC가 executable evidence로 관찰되고 open blocker가 0이거나, 8회 loop budget을 소진해 blocker를 명시적으로 보고할 때.

<!-- tene:section:layers -->
## Understanding Layers

| Layer | Planned work | Expected artifacts | Verification |
|---|---|---|---|
| Interface | manifest/marketplace/skill identities와 prompts migration | `.codex-plugin/plugin.json`, `.agents/plugins/marketplace.json`, `skills/*` | static identity validator + fresh selector observation |
| Business | router canonical IDs, explicit `$tene:<phase>` parser, deprecated criterion filtering | `internal/router/*`, context/QA selection code, tests/evals | Go unit/integration + routing corpus + QA plan inventory |
| Persistence | source/stage/cache identity equality, project state preservation | staged bundle manifest, local cache inventory, before/after workflow state evidence | hash/inventory comparison + doctor |
| Infrastructure | package/version/checksum/install/update/uninstall orchestration | package/release smoke scripts and sanitized transcripts | clean-profile/local marketplace journey |

<!-- tene:section:six-questions -->
## Six Questions

| Name | Defined at | Referenced by | Called by | Input | Output/mutation |
|---|---|---|---|---|---|
| Plugin manifest identity | `.codex-plugin/plugin.json` | marketplace/package/Codex | host discovery | metadata JSON | user-visible namespace |
| Skill metadata | `skills/<phase>/SKILL.md`, `agents/openai.yaml` | package and host | explicit/implicit invocation | request + state | skill outcome/core commands |
| Router | `internal/router/router.go` | route CLI/tests/evals | `tene-workflow route` | text/phase/active | canonical suffix, confidence, reason |
| Active criterion selector | context builder and QA planner | phase context/QA gate | `context build`, `qa plan` | sprint intent IDs + intent statuses + criteria | confirmed-only spec/cases |
| Package lifecycle | `scripts/package-plugin.sh`, release smoke | release/local install | maintainer/QA | version/platform/source tree | staged/installable bundle |
| Installed launcher | `scripts/tene-workflow` | all skills | installed invocation | argv/platform/checksum | core process and repository mutations |

<!-- tene:section:traceability -->
## 추적성

- T1 `task_0000380ydvzc0pjsatbyeff3yc` → AC namespace → manifest/skills/marketplace.
- T2 `task_0000380ydwcs0gd3wr3vy2wcdr` → AC namespace → router/prompts/docs/evals; depends on T1.
- T3 `task_0000380ydwdvfe6vc2ajcpd7g8` → AC installed runtime → package/install/update/uninstall; depends on T1 and consumes T2 for invocation.
- T4 `task_0000380ydwexgf4pewweg3pj8g` → AC installed runtime integrity → confirmed-only context/QA selection; independent code repair, required before QA planning.
- T5 `task_0000380yed20cf2ayn6vfmkn9w` → Stop hook AC → valid common JSON output, no mutation regression.
- T6 `task_0000380yed2vtrspn07jwg36sr` → language AC → common workflow reference, 9개 skill 계약과 한국어 fixture.

<!-- tene:section:decisions -->
## 결정 사항

- Canonical internal skill ID는 UI alias만 바꾸는 것이 아니라 짧은 suffix 자체로 통일한다. 이로써 selector, router 결과, agent prompt와 eval oracle을 동일하게 유지한다.
- Validator와 discovery의 기대를 충족하도록 directory 이름은 frontmatter 이름을 따른다.
- Compatibility alias는 두지 않는다. Migration 문서에서 과거 이름을 설명할 수는 있지만 executable discovery에는 노출하지 않는다.
- 의미 있는 upgrade 검증을 위해 설치 cache의 identity/version이 불변이므로 patch version을 올린다. 정확한 번호는 Do에서 현재 package version을 기준으로 결정한다.
- Host selector 동작은 black box로 검증하며, repository code가 외부 selector 조합을 소유한다고 가정하지 않는다.

<!-- tene:section:freeform -->
## 추가 관점

구현 순서는 모호한 중간 상태를 줄이기 위해 identity, 내부 reference, package/install 순으로 진행한다. T4는 deprecated candidate artifact가 필수 QA case를 생성하지 못하게 한다. 이는 workflow 자체를 사용하는 과정에서 드러난 결함이므로 설치본 handoff를 신뢰할 수 있게 만드는 범위에 포함한다.

<!-- tene:section:work-packages -->
## 작업 패키지

### WP1 — Identity surface (T1)

- Manifest의 plugin identity/display name과 marketplace entry를 `tene`로 변경한다.
- 9개 skill directory와 frontmatter name을 간결한 suffix로 변경한다.
- 각 `agents/openai.yaml`의 default prompt/display metadata를 갱신한다.
- 예상 rollback: manifest/marketplace와 directory 이동을 하나의 atomic patch로 복구한다.

### WP2 — Router와 authored reference (T2)

- 명시적 `$tene:<phase>`를 parse하여 짧은 canonical skill ID를 반환한다.
- Executable invocation을 규정하는 unit test, eval fixture, release smoke와 현재 product/design 문서를 갱신한다.
- Executable fixture가 아닌 과거 archived Sprint 문구는 유지한다. Archive 이력은 불변이다.
- 예상 rollback: parser/oracle migration을 함께 되돌리며 절반만 migration된 identifier는 남기지 않는다.

### WP3 — 설치 lifecycle (T3)

- Patch version을 올리고 host platform용 plugin을 stage한다.
- Source/stage inventory와 checksum을 검증한다.
- Local marketplace entry를 갱신·재설치하고 restart/new session에서 모든 selector와 대표 skill handoff를 관찰한다.
- Upgrade/uninstall이 repository workflow docs/state를 변경하지 않는지 검증한다.
- 예상 rollback: 이전 cache plugin version을 재설치한다. Repository state는 rollback할 필요가 없다.

### WP4 — Confirmed-only criterion (T4)

- Context와 QA planning/gating에는 소유 intent가 `confirmed`인 criterion만 중앙에서 선택하도록 한다.
- Blocking criterion을 가진 deprecated/superseded intent 회귀 사례를 추가한다.
- 현재 Sprint가 confirmed AC만 생성하는지 검증한다.
- 예상 rollback: code만 되돌리며 저장된 deprecated record는 audit history로 유지한다.

<!-- tene:section:dependencies -->
## 의존성

의존 순서는 `T1 → T2 → T3`, `T1 → T3`, `T4 → QA plan`이다. T1과 T4는 서로 다른 code surface를 다루므로 독립 구현할 수 있지만 workflow transition은 parent가 소유한다. T3에는 외부 local plugin 변경과 host restart가 필요할 수 있으며, capability가 없으면 source test로 성공을 추론하지 않고 blocker로 유지한다.

<!-- tene:section:verification -->
## 검증

1. Static identity matrix: 9개 skill 모두에서 source/stage/install manifest name, marketplace entry, directory basename, frontmatter name과 agent prompt가 일치한다.
2. Router: 9개 명시적 `$tene:<phase>` 사례, 인접 negative, 이전 namespace 거부, active-phase implicit corpus를 검증한다.
3. Core regression: deprecated/superseded intent criterion은 context와 QA plan에서 제외하고 confirmed criterion은 유지한다.
4. Package: launcher checksum 정상·변조 사례, 필수 inventory, platform binary 실행을 검증한다.
5. Local lifecycle: clean install → restart/new session → selector inventory → 대표 `$tene:sprint`/handoff → update → uninstall을 수행하고 전후 doctor 및 state hash를 비교한다.
6. Loop Check: 모든 변경 artifact를 T1–T4에 연결하고 executable design contract를 통과시키며, 독립 evaluator가 semantic namespace drift 부재를 확인한다.
7. QA: blocking AC를 interface/business/persistence/infrastructure의 실제값-기대값 관찰로 검증한다. 관찰하지 않은 external host unknown을 pass로 바꾸지 않는다.

<!-- tene:section:risks -->
## 위험

- Codex가 이전 ID의 plugin을 cache할 수 있다. Version bump, 명시적 remove/disable, restart와 새 conversation으로 완화한다.
- 짧은 suffix가 충돌할 수 있다. Qualified `$tene:*`를 구분자로 사용하고 host selector에서 9개 모두를 검증한다.
- 대량 directory rename이 오래된 docs/tests를 남길 수 있다. Targeted identity audit와 package inventory로 탐지한다.
- Archived docs에는 과거 이름이 정당하게 남을 수 있다. Scanner는 불변 historical record와 active/runtime fixture를 구분한다.
- Local install 변경에는 host permission/restart가 필요할 수 있다. 해당 단계에서 정확히 필요한 외부 action만 요청한다.
- Security filter가 보안 예시에 반응할 수 있다. 민감한 raw output은 보존하지 않고 credential이 필요한 command에는 secrets workflow만 사용한다.

<!-- tene:section:yagni -->
## YAGNI

### WP5 — Stop hook 계약 복구 (T5)

- Stop 이벤트에서 허용되지 않는 `hookSpecificOutput.additionalContext`를 제거한다.
- advisory output은 `continue: true`와 `systemMessage`를 사용하며 `decision: block`으로 turn을 강제 연장하지 않는다.
- active Sprint fixture로 JSON schema·exit code·무변경성을 회귀 검증한다.

### WP6 — 문서 언어 연속성 (T6)

- 공통 workflow reference와 9개 phase skill에 현재 사용자 대화 언어 기본 규칙을 둔다.
- 한국어 요청 fixture가 모든 phase authored section에서 한국어를 유지하는지 eval한다.
- machine marker와 code identifier는 번역 대상에서 제외한다.

- Selector UX만을 위해 Go module, Git repository, release staging directory, provider organization 또는 documentation URL을 변경하지 않는다.
- Alias resolution, 자동 cache 삭제, public marketplace publication, App Server 또는 remote MCP를 구현하지 않는다.
- T4를 전체 V2 semantic guard audit로 확장하지 않는다. Confirmed-status criterion 선택만 복구하고 regression test로 검증한다.
