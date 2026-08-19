# 기획 의도 기반 QA와 Harness Gate 설계

## 1. QA의 목표

QA는 “test script가 통과했는가”가 아니라 “사용자가 합의한 목적이 실제 시스템 전체에서 충족되는가”를 판단한다.

```text
Intent
  → Acceptance Criterion
  → User Journey / Business Invariant
  → Observable UI/API/Data outcome
  → Evidence
  → Verdict
```

## 2. Intent 모델

다음 항목은 QA oracle 생성에 사용한다.

- Actor/persona와 권한
- Preconditions
- Trigger/action
- Expected UX state sequence
- Business rule와 forbidden outcome
- Data mutation/invariant
- External side effect/event
- Error/retry/recovery behavior
- Non-functional threshold
- Explicit out-of-scope

각 항목은 source와 confirmation status를 가져야 한다. `proposed` intent는 blocking verdict의 근거로 사용할 수 없고 clarification gap을 만든다.

## 3. QA 레이어

### L1 Static and structural

- compile/type/lint
- forbidden dependency
- schema/link/template validation
- secret pattern and unsafe command scan

### L2 Unit and contract

- function/domain rule
- API/event/schema contract
- mutation/negative case

### L3 Integration and data flow

- Interface → Business Logic → Persistence
- transaction, cache, queue, external adapter
- before/after state와 idempotency

### L4 System E2E

- 실제 process/services
- browser/CLI/API entry point
- representative environment/seed

### L5 Intent and UX journey

- persona별 goal completion
- 화면/route/modal/loading/error/retry/success state transition
- accessibility와 user-visible feedback
- backend effect와 UI 결과 correlation

### L6 Adversarial and recovery

- invalid/duplicate/out-of-order action
- network failure, timeout, partial success
- permission escalation, stale state, concurrent action

### L7 Regression and drift

- 이전 sprint journey
- impacted requirement/test 재실행
- snapshot/schema/version compatibility

## 4. Observer architecture

| Observer | 수집 evidence |
|---|---|
| Browser/Chrome/Playwright | DOM/accessibility tree, screenshot, route, console, network |
| CLI | exit code, stdout/stderr redacted, prompt sequence |
| API | request/response schema, status, correlation ID |
| Persistence | read-only before/after query, transaction outcome |
| Queue/Event | topic, event type, idempotency/correlation key |
| Observability | structured log, trace/span, metric threshold |
| Git/Code | diff, symbol/call graph, generated artifact hash |

Secret value와 raw credential은 evidence에 저장하지 않는다. redactor는 저장 전 적용되고 누출 pattern 발견 시 gate를 즉시 fail한다.

## 5. Test charter 생성

Confirmed AC마다 최소 다음 case를 검토한다.

1. Happy path
2. Alternate path
3. Negative/forbidden path
4. Error and recovery
5. Role/permission variation
6. Data side effect and invariant
7. Previous sprint regression

모든 조합을 폭발시키지 않는다. risk, state edge, changed graph neighborhood를 기준으로 pairwise/impact selection을 한다.

## 6. Evaluator 분리

- Builder agent: 구현과 self-check
- Deterministic evaluator: schema, test, invariant, threshold
- Independent QA agent: UX, intent coverage, missing scenario
- Human: policy ambiguity, taste, waiver, final product acceptance

Independent QA agent에는 builder의 chain-of-thought나 자기평가를 주지 않고, confirmed intent, artifact, observable evidence만 제공한다.

## 7. Gap 분류

| Type | 의미 | 다음 경로 |
|---|---|---|
| implementation-defect | spec은 명확하나 코드가 불일치 | do |
| design-defect | design이 intent/constraint를 충족 못함 | design |
| spec-ambiguity | 기대 결과가 확정되지 않음 | prd + user decision |
| test-defect | oracle/fixture/driver 오류 | qa harness fix |
| environment-failure | service/seed/tool 문제 | environment repair |
| observation-insufficient | 판정 evidence 부족 | observer 추가 |
| accepted-risk | 명시적 waiver | report + expiry |

## 8. Gate 점수와 100% 정의

### Blocking coverage

```text
blockingCoverage = evidencedPassedBlockingAC / totalConfirmedBlockingAC
```

다음 조건을 모두 만족해야 QA pass다.

- Blocking coverage = 100%
- Critical/high gap = 0
- Mandatory layer gate pass
- Required evidence completeness = 100%
- Secret leakage = 0
- Waiver 없는 failed invariant = 0
- Regression selection 완료

Low-risk optional AC는 별도 점수로 표시할 수 있지만 blocking score와 섞지 않는다. LLM confidence나 “대략 100%” 표현은 허용하지 않는다.

## 9. Evidence manifest

```json
{
  "runId": "QA-RUN-20260820-01",
  "sprintId": "SPR-004-checkout-recovery",
  "environment": {"name":"ephemeral-test","fingerprint":"sha256:..."},
  "cases": [{
    "id":"QA-PAY-009",
    "intentRefs":["REQ-PAY-014","RULE-ORDER-003"],
    "journeyEdges":["pending->payment_failed"],
    "observers":["browser","api","db"],
    "artifacts":["screens/009.png","network/009.json","db/009.json"],
    "verdict":"passed"
  }],
  "redactionVerified": true
}
```

## 10. Loop Check와 QA 차이

- Loop Check: 문서·설계·코드가 서로 일치하는가
- QA: 실행된 제품이 실제 사용자·business intent를 충족하는가

Loop Check가 통과해도 잘못 정의된 spec일 수 있다. QA 중 spec ambiguity가 발견되면 PRD로 되돌아가 사용자 결정을 받는다.

