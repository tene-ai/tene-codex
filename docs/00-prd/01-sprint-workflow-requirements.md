# Sprint, Workflow 및 Task Management 요구사항

## 1. 계층 모델

```text
Workflow / Product Initiative
  └─ Sprint Master Plan
      ├─ Sprint A
      │   ├─ Phase
      │   ├─ Task
      │   ├─ Gate
      │   └─ Evidence
      └─ Sprint B
```

- Workflow: 제품 목표 또는 release 단위의 장기 흐름
- Sprint Master Plan: 여러 sprint의 dependency와 진행 상태
- Sprint: 하나의 coherent outcome을 만드는 전체 lifecycle
- Phase: PRD/Plan/Design/Do/Loop Check/QA/Report/Archive
- Task: phase 안에서 수행하는 atomic work item
- Gate: 다음 phase로 이동하기 위한 검증 조건
- Gap: 기대와 실제의 불일치
- Evidence: gate 판정에 사용된 실행 결과

## 2. 상태 머신

### Sprint 상태

```yaml
states:
  - draft
  - prd
  - plan
  - design
  - do
  - loop-check
  - qa
  - report
  - archived
terminal_states: [archived, cancelled]
```

### Transition 규칙

| From | To | 필수 조건 |
|---|---|---|
| draft | prd | sprint ID, title, initial goal 존재 |
| prd | plan | confirmed requirement와 unresolved policy 목록 존재 |
| plan | design | task dependency, risk, verification plan 존재 |
| design | do | Understanding Layer 영향과 interface/data contract 존재 |
| do | loop-check | planned task가 done/deferred/blocked 중 하나로 분류 |
| loop-check | do | blocking gap 수정 필요 |
| loop-check | qa | PRD/plan/design/code gap 0 |
| qa | do | product defect 또는 implementation defect 발견 |
| qa | prd | intent ambiguity 또는 requirement 변경 필요 |
| qa | report | blocking AC 100%, mandatory QA gate pass |
| report | archived | report 완성, deferred/policy owner 지정, state snapshot 저장 |

phase skip은 기본 금지한다. 예외는 `waiver`에 이유, 승인자, 위험, 만료 조건을 기록해야 한다.

## 3. Canonical 상태 파일

`.tene/`는 secret vault이므로 workflow 상태를 저장하지 않는다.

```text
.tene-workflow/
├── project.yaml
├── master-plan.yaml
├── active.json
├── decisions.ndjson
├── events.ndjson
├── checkpoints/
├── indexes/
└── archive/
```

### `active.json`

```json
{
  "schemaVersion": 1,
  "workflowId": "WF-2026-Q3",
  "activeSprintId": "SPR-004-checkout-recovery",
  "phase": "loop-check",
  "phaseIteration": 2,
  "activeTaskIds": ["TASK-018"],
  "blockingGapIds": ["GAP-007"],
  "lastCheckpoint": "CP-20260820-143012",
  "updatedAt": "2026-08-20T14:30:12+09:00"
}
```

Active state는 pointer/projection만 가진다. full transcript, test log, 문서 본문을 넣지 않는다.

## 4. Sprint Master Plan 요구사항

Master plan은 다음을 포함한다.

- Product/workflow goal과 success metrics
- Sprint 목록, 목적, dependency, 상태, priority
- Cross-sprint architecture/business/security invariant
- Release/milestone mapping
- Global risks와 decision backlog
- 현재 critical path
- Capacity/budget 또는 timebox
- 완료/취소된 sprint archive pointer

Sprint 추가 시 dependency cycle과 ID 중복을 validator가 검사한다.

## 5. Task 모델

```yaml
id: TASK-018
sprint: SPR-004-checkout-recovery
phase: do
title: payment failure rollback 구현
status: in_progress
dependsOn: [TASK-012]
intentRefs: [REQ-PAY-014, RULE-ORDER-003]
designRefs: [DES-PAY-006]
layers: [BusinessLogic, Persistence]
filesExpected:
  - src/payment/confirm.ts
doneWhen:
  - order status remains pending on declined payment
  - payment failure event is emitted once
verificationRefs: [QA-PAY-009]
```

Task status는 `pending | ready | in_progress | blocked | done | deferred | cancelled`다. `done`에는 evidence가, `deferred`에는 reason/owner/target sprint가 필요하다.

## 6. Session Resume 로직

새 session 또는 `$tene-status` 호출 시:

1. `.tene-workflow/project.yaml` schema 확인
2. `active.json` 읽기
3. master plan에서 active sprint 확인
4. active sprint document manifest와 hash 확인
5. blocking gap, open policy, current tasks 검색
6. current task와 연결된 spec graph neighborhood 조립
7. “현재 위치 / 완료 / 미완료 / 위험 / 다음 action” 요약
8. user request와 현재 action이 충돌하면 확인 요청

Codex memory는 사용자 선호 기억에 보조적으로 사용할 수 있지만 progress source of truth로 사용하지 않는다.

## 7. Clear, Compact, Archive

### Compact

- 오래된 event를 sprint summary로 축약
- active sprint와 unresolved item은 유지
- original event는 archive file에 보존
- archive checksum과 pointer 기록

### Clear

- terminal sprint만 active projection에서 제거
- orphan checkpoint와 derived index 정리
- canonical docs/evidence를 삭제하지 않음
- `--purge-evidence` 같은 파괴적 옵션은 별도 explicit confirmation 필요

### Archive

- sprint directory를 `docs/sprints/_archive/YYYY-MM/<sprint>/`로 이동하거나 manifest status만 archived로 변경
- Git history가 있는 환경에서는 rename을 권장
- master plan에는 summary와 archive link만 유지

## 8. 동시성 및 충돌

- 한 sprint의 phase transition은 lock/version compare로 직렬화한다.
- 서로 다른 sprint는 병렬 실행 가능하다.
- 같은 파일을 수정할 수 있는 task는 worktree 또는 dependency로 격리한다.
- state write는 temp file + atomic rename 또는 transactional writer를 사용한다.
- event에는 monotonic sequence와 idempotency key를 둔다.
- merge 후 graph/index를 재생성하고 drift gate를 다시 실행한다.

