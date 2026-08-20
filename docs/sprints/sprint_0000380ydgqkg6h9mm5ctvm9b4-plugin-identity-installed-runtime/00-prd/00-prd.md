---
schema_version: 1.0.0
document_type: prd
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380ydgqkg6h9mm5ctvm9b4
phase: prd
status: confirmed
revision: 1097
intent_ids: [intent_0000380ydppfay2szegwvaqsqm, intent_0000380ydpq9tm54tgwpqnfd30, intent_0000380yeckc0mm5dsr9j59h6m, intent_0000380yecmcwb7v4ttjm9bbrm]
generated_at: 2026-08-20T09:06:41Z
generated_by: tene-workflow
---

# prd — Plugin Identity and Installed Runtime Verification

<!-- tene:section:purpose -->
## 목적

`docs/00-prd`와 `docs/02-design`이 정의한 제품 identity와 설치 계약을 실제 Codex discovery 경계에서 검증한다. 저장소·배포 소스의 이름은 `tene-codex`로 유지하지만 사용자가 보는 plugin namespace는 `tene`, 각 phase skill suffix는 짧은 이름이어야 한다. source checkout 성공과 설치본 성공을 별개로 검증해 installed-source drift를 차단한다.

<!-- tene:section:scope -->
## 범위

포함:

- `.codex-plugin/plugin.json`, local marketplace metadata, 9개 skill frontmatter/directory/agent metadata의 namespace migration
- `$tene:sprint`, `$tene:prd`, `$tene:plan`, `$tene:design`, `$tene:loop-check`, `$tene:qa`, `$tene:report`, `$tene:status`, `$tene:secrets` discovery와 invocation
- router/eval/docs/scripts/package metadata의 canonical skill ID 변경
- staged package와 로컬 설치 cache의 content/metadata/checksum 비교
- clean install, update, restart/new-session, uninstall 후 project workflow state 보존 검증

제외:

- Git repository와 provider organization 이름 변경
- public marketplace 제출, remote MCP, App Server 구현
- obsolete `$tene-codex:*` 또는 `$tene:tene-*` compatibility alias 유지

<!-- tene:section:layers -->
## Understanding Layers

| Layer | Entry/components | Intended change | Upstream/downstream | Required evidence |
|---|---|---|---|---|
| Interface | plugin selector, explicit `$...` invocation, implicit routing | namespace `tene`, suffix는 phase 이름 | manifest/frontmatter → Codex discovery → user selection | fresh session selector와 invocation observation |
| Business Logic | router explicit parser, canonical selected skill IDs, release smoke | `tene-*` 내부 ID를 짧은 canonical ID로 migration | prompt → router → phase skill → CLI | routing corpus와 handoff assertions |
| Persistence | marketplace entry, staged bundle, plugin cache, repository workflow state | source/staged/installed metadata 일치; update/uninstall이 project state를 변경하지 않음 | package/install lifecycle | hashes, installed manifest inventory, before/after doctor |
| Infrastructure | package script, bundled launcher/binary, checksum, Codex capability probe | 새 plugin identity로 package/install되고 launcher가 검증 후 실행 | OS/arch bundle → Codex host → repo CLI | clean-profile/package/doctor execution transcript |

<!-- tene:section:six-questions -->
## Six Questions

| Name | Defined at | Imported/referenced at | Called/used at | Input shape | Output/mutation |
|---|---|---|---|---|---|
| Plugin namespace manifest | `.codex-plugin/plugin.json` top-level `name` | marketplace entry, Codex plugin discovery | selector namespace composition | plugin metadata JSON | discovered namespace; no project mutation |
| Skill identity | `skills/<phase>/SKILL.md` frontmatter `name` | skill agent metadata, prompts, router evals | explicit/implicit invocation | user request + active phase | selected phase skill/action |
| Explicit router | `internal/router/router.go` | route command/tests | `tene-workflow route` | text, phase, active flag | canonical skill ID/confidence/reason |
| Portable launcher | `scripts/tene-workflow` | every installed phase skill | wrapper execution | CLI argv, platform, checksum file | execs PATH/bundled/source core; repository state via core |
| Package builder | `scripts/package-plugin.sh` | release smoke/local install | release workflow | version, stage path, platform targets | staged plugin files and binaries |
| Codex host selector | external host implementation (provider unknown) | fresh Codex session | user skill search/select | installed plugin+skill metadata | observable `$tene:<phase>` entries |

<!-- tene:section:traceability -->
## 추적성

- `intent_0000380ydppfay2szegwvaqsqm` → `ac_0000380ydppfarshvvnhww980w` → namespace/selector/router/package tasks → fresh-session discovery evidence.
- `intent_0000380ydpq9tm54tgwpqnfd30` → `ac_0000380ydpq9vae1we5cfzkw78` → staged/install/update/uninstall tasks → CLI health and state-preservation evidence.
- `intent_0000380yeckc0mm5dsr9j59h6m` → `ac_0000380yeckc1p516k1yvt3q3r` → Stop hook schema repair와 active-Sprint 회귀 테스트.
- `intent_0000380yecmcwb7v4ttjm9bbrm` → `ac_0000380yecmcxpx719rfy2xb7m` → 공통 언어 계약과 9개 phase skill 검증.
- Source intent: `docs/00-prd/03-plugin-architecture.md` §8의 Marketplace Name `tene`; user request dated 2026-08-20 supersedes the older `$tene-*` invocation spelling in `docs/00-prd/06-skills-commands-triggers.md`.

<!-- tene:section:decisions -->
## 결정 사항

- 확정: plugin namespace는 `tene`이며 repository 이름은 `tene-codex`로 유지한다.
- 확정: skill suffix는 `sprint|prd|plan|design|loop-check|qa|report|status|secrets`이다.
- 확정: 중복 discovery는 observable을 위반하므로 폐기된 selector alias를 유지하지 않는다.
- 검증할 가정: qualified name은 다른 plugin의 짧은 skill name과 충돌하지 않는다.
- 미확정 운영 세부사항: 현재 release 자동화를 확인한 뒤 Plan에서 다음 patch version을 선택하되, 그 결과가 namespace 계약을 바꾸지는 않는다.
- Host selector 생성은 외부 동작이며 restart/new Codex session에서 관찰하기 전까지 `unknown`으로 유지한다.

<!-- tene:section:freeform -->
## 추가 관점

selector 문자열은 plugin manifest 이름과 skill frontmatter 이름이라는 독립된 두 식별자를 조합한다. 하나만 변경하면 `$tene:tene-plan` 또는 `$tene-codex:plan`이 되어 확정 의도를 충족하지 못한다. Codex는 설치 cache를 읽으므로 source만 수정해서도 부족하다. packaging, 재설치, 재시작, 새 session 관찰이 하나의 acceptance boundary다.

<!-- tene:section:problem -->
## 문제

현재 source와 설치 cache는 plugin namespace를 `tene-codex`, skill name을 `tene-plan` 같은 형태로 선언한다. Codex는 이를 결합해 `$tene-codex:tene-plan`으로 표시한다. 이는 제품 이름 `tene`와 사용자가 원하는 간결한 phase UX에 어긋나며, source만 수정하고 설치본이 갱신되지 않으면 사용자는 계속 오래된 namespace를 보게 된다.

<!-- tene:section:actors -->
## 행위자

- Codex user: skill 검색·선택·명시 호출·자연어 호출을 수행하고 결과를 관찰한다.
- Plugin maintainer: manifest, skill metadata, router, package와 release compatibility를 유지한다.
- Codex host: 설치 cache를 discovery하고 selector를 구성하는 외부 runtime이다.
- `tene-workflow`: repository state와 phase gate를 소유하는 deterministic core다.

<!-- tene:section:journeys -->
## 사용자·데이터 여정

1. Maintainer가 source metadata를 migration하고 package를 staging한다.
2. Validator가 plugin/skill identities, agent prompts, bundle inventory와 checksum을 검사한다.
3. Clean local marketplace profile에 새 version을 설치하고 Codex를 restart/new session으로 연다.
4. User가 9개 skill을 검색해 정확한 `$tene:<phase>` 이름과 obsolete 이름 부재를 확인한다.
5. User가 `$tene:sprint`와 대표 handoff를 실행하고 installed bundled CLI health를 확인한다.
6. 기존 plugin에서 update한 repository와 uninstall한 repository에서 `.tene-workflow` state가 그대로 healthy인지 확인한다.
7. 실패 시 manifest/skill/router/package/install-cache/host 중 경계를 분류하고 evidence와 gap을 남긴다.

<!-- tene:section:acceptance-criteria -->
## Acceptance Criteria

### AC `ac_0000380ydppfarshvvnhww980w` — 간결한 namespace (blocking)

- Fresh Codex session에 정확히 9개의 `$tene:<phase>` skill이 discovery된다.
- `$tene:plan`을 선택·호출할 수 있다.
- `$tene-codex:*`와 `$tene:tene-*`는 discovery되지 않는다.
- explicit invocation은 100% 정확해야 하며 implicit routing은 기존 positive/negative/phase-conflict corpus에서 회귀하지 않는다.

### AC `ac_0000380ydpq9vae1we5cfzkw78` — 설치 runtime 동등성 (blocking)

- Source와 staged package, installed cache의 manifest/skill identity와 required inventory가 일치한다.
- 새 local install에서 `$tene:sprint`가 installed launcher를 통해 healthy CLI status를 반환한다.
- Bundled binary checksum mismatch는 실행 전에 fail closed 한다.
- Update와 uninstall 후 repository workflow state, docs, evidence가 보존되고 `doctor`가 healthy다.

금지 결과: 기존 namespace 중복, source-only pass, 검사하지 않은 설치 binary, plugin lifecycle에 의한 repository state 삭제, host 관찰 없이 selector 성공을 추정하는 것.

### AC `ac_0000380yeckc1p516k1yvt3q3r` — Stop hook 계약 (blocking)

- active Sprint에서 Stop handler는 Codex Stop-event schema에 맞는 JSON과 exit code 0을 반환한다.
- `invalid stop hook JSON output` 오류와 의도하지 않은 자동 continuation이 발생하지 않는다.
- hook은 workflow state를 변경하지 않고 advisory message만 제공한다.

### AC `ac_0000380yecmcxpx719rfy2xb7m` — 사용자 대화 언어 연속성 (blocking)

- 한국어 요청으로 시작한 workflow의 PRD, Plan, Design, Analysis, QA, Report authored section은 한국어로 작성한다.
- 사용자가 다른 언어를 명시하면 그 요청이 우선한다.
- frontmatter, section marker, ID, path, command, API와 code symbol은 변형하지 않는다.

<!-- tene:section:non-goals -->
## 비범위

- GitHub repository/module/package staging directory를 모두 `tene`로 rename하는 작업
- 기존 `$tene-*` 문서 표기 전체를 무비판적으로 보존하는 compatibility layer
- public marketplace 제출이나 외부 remote service 구현
- 이 Sprint에서 workflow core의 모든 semantic gap을 함께 수정하는 것; 발견 사항은 후속 V2~V7 Sprint로 연결한다.
