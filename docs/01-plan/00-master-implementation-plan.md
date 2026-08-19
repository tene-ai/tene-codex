# Master Implementation Plan

## 1. 목표와 완료 정의

어떤 저장소에서도 기획 의도를 지속적으로 기억하며, 고정 Sprint 사이클과 evidence 기반 QA를 강제하는 Codex 플러그인을 만든다. MVP 완료는 다음을 동시에 뜻한다.

- 새 프로젝트에서 1개 명령으로 초기화하고 Sprint를 생성·재개할 수 있다.
- PRD의 intent와 acceptance criterion이 설계, 코드 영향, QA evidence, report까지 ID로 추적된다.
- 허용되지 않은 단계 전이는 실패하며 이유와 다음 행동을 반환한다.
- Unit/E2E뿐 아니라 UX journey, 데이터 흐름, 실패·복구 흐름을 QA charter로 실행한다.
- Secret이 필요한 명령은 평문을 읽거나 기록하지 않고 `tene run --env <name> -- <command>`로만 실행한다.
- plugin validator, skill eval, core test, 세 개 reference project 검증을 통과한다.

## 2. 구현 전략

### 2.1 두 계층

- **Deterministic core**: Go binary `tene-workflow`. 상태 전이, 스키마 검증, 문서 scaffold, traceability 계산, evidence 판정을 담당한다.
- **Codex orchestration**: plugin skills, references, scripts. 자연어 의도 추출, 코드 조사, QA 도구 선택, 문서 서술과 review를 담당한다.

Go를 선택하는 이유는 기존 tene CLI와 배포 방식·보안 관행을 공유하고 단일 바이너리로 여러 프로젝트에 이식하기 쉽기 때문이다. core를 기존 secret CLI에 즉시 병합하지 않고 별도 바이너리로 시작해 두 책임의 보안 경계를 유지한다.

### 2.2 기본 정책

- 기본 workflow profile은 `standard`다. feature/bug/refactor는 전체 사이클, 작은 문서·주석 작업은 축약 사이클을 허용한다.
- `strict`는 모든 코드 변경에 전체 사이클과 최종 사람 승인을 요구한다.
- `light`는 PRD/Plan/Design을 하나의 spec 문서로 축약하지만 intent와 QA evidence는 생략하지 않는다.
- `off`는 자동 제안만 끄며 명시적 skill 호출은 유지한다.
- blocking acceptance criterion의 pass 비율만 100%여야 archive할 수 있다. 점수 평균으로 blocker를 상쇄하지 않는다.

## 3. 제품 Sprint

### Sprint 0 — 계약과 개발 기반

산출물:

- Go module, package layout, CI, lint/test 명령
- JSON Schema 2020-12 기반 project/sprint/task/intent/evidence schema
- error code, exit code, atomic-write 및 lock contract
- fixture 기반 schema compatibility test

Gate:

- 모든 schema의 valid/invalid fixture가 존재한다.
- 같은 입력에 같은 canonical JSON과 transition result가 나온다.
- crash simulation에서 기존 상태 파일이 손상되지 않는다.

### Sprint 1 — 상태·문서 기반

산출물:

- `init`, `sprint create/start/status/transition`, `task`, `compact`, `clear`, `archive`, `doctor`
- 표준 Sprint 디렉터리와 PRD/Plan/Design/Analysis/QA/Report template
- event journal과 active projection, session resume context
- transition guard와 human-approval 기록

Gate:

- 새 저장소에서 archived까지 전 단계 happy path가 재현된다.
- 필수 문서/section/evidence 누락 시 transition이 거절된다.
- clear는 active context를 제거하지 않고 derived cache만 제거한다.

### Sprint 2 — Intent·Spec graph

산출물:

- 대화에서 intent 후보 생성 후 사용자 confirm/replace/deprecate 흐름
- Intent, AC, journey, policy, task, design component 간 graph
- `context build`, `trace`, `impact`, `validate`
- requirement drift와 orphan 탐지

Gate:

- 한 intent에서 관련 설계·task·QA·evidence·report를 양방향 탐색한다.
- 변경된 AC의 영향 범위를 안정적으로 계산한다.
- 미확정 intent는 구현 gate를 통과하지 못한다.

### Sprint 3 — Code understanding·Loop Check

산출물:

- git diff, 언어별 index adapter, optional CodeGraph adapter
- Understanding Layer 분류와 Six Questions materializer
- spec-to-code gap evaluator, iterative repair loop와 iteration budget
- waiver/deferred decision workflow

Gate:

- 변경 symbol별 정의/참조/호출/input/output·mutation evidence가 report에 채워진다.
- 테스트 통과만으로 spec gap이 닫히지 않는다.
- 반복 한도 도달 시 실패를 숨기지 않고 unresolved gap으로 중단한다.

### Sprint 4 — Intent-driven QA

산출물:

- 7-layer QA planner와 adapter interface
- unit/contract/integration/Playwright/browser/manual adapter
- UX journey와 데이터 lineage observer, evidence manifest, independent evaluator
- blocker gate와 regression baseline

Gate:

- 대표 웹 앱에서 UI → API → DB/외부 경계의 evidence chain을 생성한다.
- empty/error/loading/permission/retry/rollback 시나리오를 실행 또는 명시적으로 waive한다.
- secret redaction test corpus에서 값이 artifact/log에 나타나지 않는다.

### Sprint 5 — Codex plugin UX

산출물:

- `.codex-plugin/plugin.json`
- 9개 skills와 `agents/openai.yaml`, references, deterministic wrapper scripts
- implicit trigger eval, explicit `$tene-*` invocation
- optional stop/session hook는 지원 surface 확인 후 defense-in-depth로만 제공

Gate:

- plugin/skill validator 통과
- 명시 호출 성공률 100%, 자연어 routing precision/recall 목표 90% 이상
- 자연어 요청만으로 현재 Sprint를 재개하고 올바른 context pack을 구성한다.

### Sprint 6 — 안정화와 배포

산출물:

- macOS/Linux binaries, checksum, install/update 문서
- Marketplace metadata와 version compatibility matrix
- greenfield, mature monolith, polyglot 프로젝트 reference run
- migration/rollback/recovery runbook

Gate:

- clean environment 설치·제거·업그레이드 검증
- 이전 minor schema를 자동 migration하고 dry-run diff를 제공
- 세 reference project에서 AC-PRODUCT-01~08을 충족

## 4. 각 Sprint 내 의무 사이클

각 제품 Sprint 자체도 tene 방식으로 수행한다.

1. PRD: intent, actor, UX journey, data journey, policy, AC를 확정한다.
2. Plan: task, dependency, 검증, 정책 결정을 작성한다.
3. Design: component, schema, interface, failure/recovery를 정의한다.
4. Do: 작은 vertical slice로 구현하고 event/evidence를 남긴다.
5. Loop Check: PRD/Plan/Design과 diff를 비교하고 gap이 0이 될 때까지 수정한다.
6. QA: 7-layer 중 적용 가능한 레이어를 실행하고 불가 항목은 waiver 승인받는다.
7. Report: 이전 Sprint 연결, 파일/기능, intent 충족, 4 Layers, 6 Questions, 이월 정책을 기록한다.
8. Archive: gate와 승인 후 immutable snapshot으로 이동한다.

## 5. 작업 순서의 핵심 경로

`schemas → state store → transition engine → document validator → intent graph → context pack → code graph adapters → loop evaluator → QA adapters → skill router → packaging`

UI나 remote service보다 이 경로를 먼저 완성한다. 어느 단계도 선행 계약을 우회해 자체 상태를 만들지 않는다.

## 6. 운영 지표

- Traceability: confirmed intent 중 evidence까지 연결된 비율
- Gate integrity: blocker가 남은 채 통과한 false-pass 수(목표 0)
- Resume fidelity: 새 세션이 사람 개입 없이 올바른 Sprint/phase/task를 찾는 비율
- Routing: skill trigger precision/recall
- QA depth: AC별 실행된 layer, journey/data-flow coverage
- Secret safety: 평문 secret artifact/log incident 수(목표 0)
- Context efficiency: pack token 예산 초과율과 stale context 포함률

