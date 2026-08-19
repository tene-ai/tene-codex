# 제품 PRD: tene Codex Plugin

## 1. 제품 비전

tene는 바이브 코딩의 속도와 자유를 유지하면서도, 기능 의도·설계·구현·검증의 연결이 사라지지 않도록 만드는 spec-driven sprint plugin이다. 사용자는 자연어로 원하는 기능을 설명하고, tene는 이를 검토 가능한 문서와 실행 가능한 gate로 바꾸며, Codex는 그 범위 안에서 구현과 검증을 반복한다.

제품의 최종 산출물은 코드만이 아니다.

- 합의된 기획 의도
- 구현 계획과 설계 근거
- 코드 및 데이터 흐름의 이해 지도
- 검증 증거와 미충족 gap
- 다음 sprint가 이어받을 결정·부채·이월 작업

## 2. 문제 정의

일반적인 바이브 코딩은 다음 문제를 가진다.

1. 대화에서 합의한 요구사항과 정책이 session 종료·compaction 후 사라진다.
2. 구현 agent가 자신의 오해에 맞는 test를 작성하고 이를 성공으로 판정한다.
3. 파일 하나의 local correctness만 보고 전체 UX·business·data flow를 놓친다.
4. PRD, plan, design, code, test가 서로 drift한다.
5. 기능 간 연결과 이전 sprint의 결정 이유가 report에 남지 않는다.
6. LLM이 부분적인 code context만 보고 중복 abstraction과 우회 dependency를 만든다.
7. secret이 `.env`, command argument, tool output을 통해 model context에 노출될 수 있다.

## 3. 대상 사용자

### Primary

- Codex를 중심으로 개인 또는 소규모 팀 프로젝트를 반복 개발하는 사용자
- stack이 달라도 동일한 개발 discipline을 적용하고 싶은 사용자
- 문서 자체보다 문서와 구현의 일치 여부를 중요하게 보는 사용자

### Secondary

- 기존 brownfield codebase의 의도와 처리 흐름을 복원하려는 팀
- QA automation은 있지만 product intent coverage가 부족한 팀
- Claude Code 등 다른 agent와 portable workflow를 공유하려는 사용자

## 4. Jobs to be Done

- “내가 기능을 설명하면, 구현 전에 빠진 정책과 예외를 찾아 PRD로 확정하고 싶다.”
- “복잡한 작업을 sprint와 task로 쪼개고 session이 바뀌어도 이어가고 싶다.”
- “agent가 설계와 다르게 구현하면 자동으로 gap을 찾아 반복 수정하고 싶다.”
- “UI 클릭뿐 아니라 API, DB, event까지 제품 의도대로 동작하는지 검증하고 싶다.”
- “완료 후 어떤 의도를 어떤 파일·로직으로 충족했는지 다음 사람이 이해하고 싶다.”
- “API key가 필요해도 AI에게 값을 보여주지 않고 작업을 실행하고 싶다.”

## 5. 핵심 기능 요구사항

### FR-01 Sprint Master Plan

여러 sprint의 목표, 선후관계, 상태, release/milestone, 공통 위험과 cross-sprint invariant를 관리한다.

### FR-02 Fixed Sprint Lifecycle

모든 implementation sprint는 아래 상태를 따른다.

```text
draft → prd → plan → design → do ↔ loop-check → qa → report → archived
```

분석만 필요한 sprint는 `do`에서 코드 변경이 없을 수 있지만, 해당 사실과 산출물은 report에 남긴다.

### FR-03 Intent Memory

사용자 대화에서 goal, actor, policy, business rule, UX state, data invariant, constraint, assumption, open question, acceptance criterion을 추출한다. 추출 결과는 `proposed`로 저장하고 사용자 확인 뒤 `confirmed`로 승격한다.

### FR-04 Standardized Documents

PRD, plan, design, analysis, QA, report는 공통 metadata, Understanding Layer, 6 Questions, traceability, policy/deferred work 영역을 포함한다. 각 문서는 필수 섹션 외에 자유 섹션을 허용한다.

### FR-05 Loop Check

PRD→plan→design→code→test를 양방향 비교한다. 미충족, 무근거 구현, 문서 drift, 설계 위반을 gap으로 만들고 blocking gap이 0이 될 때까지 수정·재검증한다.

### FR-06 Intent-driven QA

Unit/E2E 외에도 journey state transition, role/permission, error/recovery, API side effect, persistence invariant, event/queue, observability evidence를 검증한다.

### FR-07 Evidence-based Gate

phase 완료는 agent 선언이 아니라 validator 결과와 evidence manifest로 판정한다. 100%는 confirmed acceptance criterion의 blocking coverage가 100%임을 의미한다.

### FR-08 Durable State and Resume

현재 master plan, active sprint, phase, task, gate, gap, decision, checkpoint를 파일에 저장한다. 새 session에서 이를 읽어 현재 위치와 다음 action을 복원한다.

### FR-09 Compaction/Clear

active state는 작게 유지하고 완료 event와 상세 evidence는 archive로 이동한다. clear는 기록 삭제가 아니라 snapshot+archive+active projection 재작성이다.

### FR-10 Explicit and Implicit Invocation

사용자는 `$tene-*` skill 또는 명령형 요청으로 직접 호출할 수 있다. 자연어 요청에서도 명확한 trigger에 따라 필요한 skill을 자동 선택한다.

### FR-11 Secret-safe Execution

secret 이름 조회는 `tene list --json`, secret이 필요한 command 실행은 `tene run -- ...`를 사용한다. Plugin은 `tene get`, 평문 `tene export`, `.tene/` 파일 읽기를 금지한다.

## 6. 비기능 요구사항

- Portability: 언어·framework·repository layout에 종속되지 않아야 한다.
- Auditability: 모든 phase transition과 gate verdict에 근거가 있어야 한다.
- Recoverability: thread transcript 없이 canonical files로 재개 가능해야 한다.
- Bounded context: active sprint에 필요한 정보만 model context로 조립해야 한다.
- Security: secret plaintext가 prompt, stdout, document, evidence에 들어가면 안 된다.
- Determinism: schema, state transition, required section, link integrity는 script로 검증한다.
- Extensibility: project별 추가 layer/question/gate/template을 허용한다.
- Human authority: 정책 결정과 scope 변경은 사용자가 최종 승인한다.

## 7. 범위 밖

- Jira/Linear를 대체하는 범용 조직 PM 제품
- 모든 언어의 완전한 static code graph engine 자체 구현
- LLM만으로 보안·법률·접근성 인증을 최종 보증
- secret 값의 저장·복호화 로직 재구현
- 모든 project에 동일한 architecture pattern 강제
- 문서가 많다는 이유만으로 품질이 높다고 평가하는 기능

## 8. 성공 지표

- Sprint resume success: 새 session에서 human 재설명 없이 현재 phase/action 복원 ≥ 95%
- Requirement evidence coverage: blocking AC 100%
- Traceability coverage: changed production file의 sprint/intent 연결 ≥ 95%
- Drift detection recall: seeded PRD/design/code drift 탐지 ≥ 90%
- QA defect classification precision: defect/spec ambiguity/environment 구분 ≥ 85%
- Secret leakage: automated red-team scenario에서 plaintext 노출 0건
- Context efficiency: active context pack이 전체 문서 대비 20% 이하를 기본 목표로 함

