# Harness·Graph·Context Engineering 설계

## 1. Harness Engineering

Harness는 모델에게 “잘해라”라고 요청하는 대신 성공 조건을 관찰·강제할 환경을 만든다.

### Planner harness

- PRD의 AC를 task와 test/evidence 계획에 연결
- dependency와 conflict 분석
- 계획 밖 파일/범위를 탐지할 expected touch set 생성

### Builder harness

- isolated worktree 또는 clean workspace
- project-defined build/lint/test command discovery
- small task 단위 checkpoint
- command exit/output와 diff를 evidence로 수집

### Loop-check harness

- 문서 간 coverage matrix
- actual changed symbol/call/data flow 분석
- design invariant와 forbidden dependency 검사
- blocking gap 생성 및 재실행

### QA harness

- deterministic test + agentic journey test
- browser/API/DB/log observation
- independent evaluator
- evidence completeness와 gate 판정

### Reporter harness

- Git diff, symbol graph, task/gate/evidence에서 report 자동 합성
- 사용자가 결정할 정책과 deferred item을 자동 누락 검사
- harness failure 자체를 회고 항목으로 승격

## 2. Graph Engineering

### Unified traceability graph

```text
Goal → Requirement → AcceptanceCriterion
  ├→ Journey → State → Transition
  ├→ BusinessRule → Invariant
  └→ DesignComponent → Symbol → File
                         ├→ Call/Import/DataFlow
                         └→ Test → Run → Evidence → Verdict
Sprint → Task ───────────────────────┘
Decision → supersedes → Decision
```

### Node 최소 schema

```yaml
id: REQ-CHECKOUT-014
kind: requirement
title: failed payment preserves order
status: confirmed
version: 2
sourceRefs: [conversation:abc#turn-18]
owners: [product]
```

### Edge 최소 schema

```yaml
from: REQ-CHECKOUT-014
type: verified_by
to: QA-CHECKOUT-022
confidence: explicit
source: docs/sprints/SPR-004/04-qa/qa-plan.md
```

### Graph source hierarchy

1. Explicit document IDs/links — highest authority
2. Code parser/LSP/CodeGraph symbol relations
3. Runtime traces/network/log correlation
4. LLM-inferred relation — `inferred` 표시, promotion 전 검증

Graph DB는 derived index다. Canonical documents를 parsing해 재생성 가능해야 한다.

### Impact analysis

파일 또는 symbol 변경 시 역방향 traversal로 연결된 design, requirement, journey, test를 찾는다. 연결이 없는 changed production symbol은 `orphan-change` gap으로 보고한다.

## 3. Understanding Layer를 graph에 반영

각 symbol/component node에 layer를 하나 이상 부여한다.

```text
Interface --calls--> BusinessLogic --reads/writes--> Persistence
    └---------------- runs-on / guarded-by ----------> Infrastructure
```

허용 관계는 project config로 정의한다. 예를 들어 Interface가 Persistence를 직접 호출하지 못하게 할 수 있지만, tene 기본값은 경고만 제공하고 stack-specific rule을 강제하지 않는다.

## 4. Six Questions 자동화

| 질문 | 주요 데이터원 |
|---|---|
| 이름 | AST/LSP/CodeGraph, framework registry |
| 정의 파일 | symbol definition |
| import/reference | import graph, reference search |
| call/use | call graph, route/event registration, runtime trace |
| input | type/schema/signature, request/event fixture |
| output/change | return type, mutation, DB/event/log side effect |

정적 분석만으로 알 수 없는 reflection, dependency injection, event bus는 runtime evidence 또는 명시적 document edge로 보완한다. `unknown`을 허용하되 숨기지 않는다.

## 5. Context Engineering

### Context Pack 구성

매 phase마다 다른 context를 조립한다.

| Phase | 필수 context |
|---|---|
| PRD | 이전 sprint summary, product goal, relevant intent, open policy |
| Plan | confirmed PRD, repo map, constraints, risk history |
| Design | PRD+plan, affected graph neighborhood, current contracts |
| Do | current task, relevant design slice, target symbols, doneWhen |
| Loop Check | PRD/plan/design hashes, actual diff/graph/test result |
| QA | AC, journey/invariant, environment, observer capabilities |
| Report | all verdict summaries, diff/symbol map, decisions/deferred |

### Budget 정책

1. Hard constraints와 active AC
2. Current task/design slice
3. Connected code/symbol context
4. Recent blocking gaps
5. Previous sprint summary
6. Optional background

budget 초과 시 낮은 순위를 pointer+summary로 바꾼다. source artifact를 삭제하거나 silently truncate하지 않는다.

### Context pollution 방지

- raw test log는 evidence file로 저장하고 failure excerpt만 main thread에 반환
- code exploration은 subagent가 수행하고 verified graph delta만 반환
- 전체 sprint 문서를 prompt에 반복 삽입하지 않음
- document hash로 stale context 탐지
- confirmed/proposed/deprecated intent를 명확히 구분

## 6. Dynamic Workflow

Workflow는 phase 순서를 강제하되 phase 내부 전략은 agent에게 맡긴다.

- Deterministic controller: state, transition, required artifact, gate
- Agent planner: 작업 분해, 탐색 전략, 도구 선택
- Human authority: policy/scope/waiver/destructive action

이 구조는 rigid graph의 취약성과 free-form agent의 drift를 모두 줄인다.

## 7. Feedback loop

```text
Observe failure
  → classify: product/spec/harness/environment
  → product: fix code
  → spec: ask user + version artifact
  → harness: improve validator/tool/context
  → environment: repair fixture/runtime
  → rerun narrow gate
  → rerun impacted full gate
```

같은 실패가 반복되면 일회성 prompt를 추가하지 않고 template, validator, AGENTS guidance, skill reference, test fixture 중 가장 작은 durable surface를 개선한다.

