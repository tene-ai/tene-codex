# Verification and Release Plan

## 1. 검증 피라미드

| 층 | 대상 | 도구/방식 | merge blocker |
|---|---|---|---|
| Contract | schema, CLI JSON, transition | Go table test, JSON fixtures | 예 |
| Unit | reducer, rules, redactor, graph | `go test` | 예 |
| Integration | filesystem, git, tene subprocess | temp repo, fake binary | 예 |
| System | plugin → skill → core → documents | clean Codex project | 예 |
| Journey | UX/API/data 흐름 | Playwright/browser + observers | release |
| Adversarial | crash, race, prompt/secret leakage | fault injection, canary | release |
| Regression | archived Sprint 재검증 | golden/evidence replay | release |

## 2. Reference projects

1. **Greenfield web app**: UI form → API → DB, auth/validation/error/retry journey.
2. **Mature monolith**: 기존 문서와 기술 부채, 큰 call graph, partial test suite.
3. **Polyglot/service app**: frontend/backend/worker/queue/external API, secret-required E2E.

각 프로젝트에서 동일한 feature를 Sprint 전체 사이클로 수행해 portability와 이해 깊이를 비교한다.

## 3. 필수 test scenarios

### Workflow/state

- invalid phase skip, stale revision, double writer, interrupted write
- session restart 후 active Sprint와 next action 복원
- compact 전후 semantic state 동일
- clear가 source artifact/evidence를 지우지 않음
- archive 후 active pointer와 master plan 일치

### Intent/graph/context

- 모호한 요구가 candidate로 남고 확인 전 구현되지 않음
- policy 변경 시 impacted AC/design/test가 모두 표시됨
- orphan AC, stale evidence, conflicting intent 탐지
- token 예산 초과 시 priority 규칙대로 축약하고 provenance 유지

### Loop/QA

- 테스트는 pass하지만 AC가 누락된 mutation을 fail로 판정
- UI 성공 뒤 DB 반영 실패를 observer chain이 탐지
- loading/error/permission/retry/back-navigation/accessibility flow
- waiver 만료와 blocker 재개방
- evaluator가 없는/오염된 evidence를 pass로 인정하지 않음

### Secret

- fake secret canary가 stdout/stderr/Markdown/JSON/evidence에 0회 등장
- `tene` 미설치, env 미존재, 권한 거부 시 fail-closed
- 금지 명령과 vault path read가 preflight에서 거절됨
- child exit code와 sanitized diagnostics가 보존됨

## 4. Skill evaluation set

- 명시 호출: 각 `$tene-*` skill 5개씩 happy/error/resume fixture
- 자연어 positive: 한국어/영어, 직접 키워드/의미 기반 표현
- hard negative: 일반적인 “계획 알려줘”, “상태 보여줘”, secret과 무관한 환경변수 질문
- collision: prd vs plan, loop-check vs qa, status vs report
- expected output: selected skill, required core action, 금지 action, confirmation 필요 여부

목표는 명시 호출 100%, implicit precision/recall 각각 90% 이상, secret unsafe action 0건이다.

## 5. Gate 정책

### Pull request gate

- format/lint/unit/contract/integration
- schema backward compatibility
- generated artifact drift 없음
- 변경된 public contract의 문서와 fixture 갱신

### Release candidate gate

- plugin/skill validator
- clean install/update/uninstall
- reference project 3종
- secret canary, concurrency/crash test
- SBOM/checksum/license/security scan
- known limitation과 migration note

### Archive gate

- 모든 blocking AC `passed`
- evidence가 존재하고 hash 검증됨
- unresolved gap은 0; non-blocking은 승인된 waiver/deferred item으로 전환
- report 필수 섹션, 4 Layers, 6 Questions 완성
- profile이 요구하는 human approval 존재

## 6. 릴리스·호환성

- SemVer: schema/CLI/plugin 각각 version을 기록하고 compatible range를 manifest/doctor에서 확인한다.
- major: destructive migration 또는 public command 제거.
- minor: backward-compatible schema/command/skill 추가.
- patch: contract를 바꾸지 않는 수정.
- migration은 `plan → dry-run diff → backup → apply → verify` 순서이며 실패 시 backup과 journal로 복원한다.
- 배포 전 한 버전 이전 binary로 read-only status가 가능한지 검증한다.

## 7. Evidence 보존

- 작은 manifest와 요약은 Git 추적 대상이다.
- screenshot, trace, video 같은 큰 artifact는 `.tene-workflow/evidence/<run-id>/`에 두고 기본 gitignore한다.
- CI artifact URL을 쓸 경우 content hash, 만료일, 생성 도구를 manifest에 기록한다.
- archive 시 사라질 evidence는 허용하지 않으며, 만료 전 durable storage로 승격하거나 waiver를 받는다.

