# Intent-driven QA, Evidence and Gates

## 1. QA compiler

QA는 test file 목록이 아니라 confirmed intent와 observable AC를 실행 가능한 charter로 compile한다.

```text
Intent/AC + UX Journey + Data Journey + Policies
        ↓ normalize and coverage expand
Test Charters (happy/alternate/error/recovery)
        ↓ capability-aware planning
Adapters + Observers
        ↓ run
Evidence Manifest
        ↓ isolated evaluator
AC Verdict + Gaps + Gate
```

## 2. Seven layers

| Layer | 질문 | 대표 evidence |
|---|---|---|
| L1 Static | 구조/타입/security rule이 맞는가 | lint/type/scan |
| L2 Unit/Contract | 규칙과 boundary contract가 맞는가 | test/JUnit/schema |
| L3 Integration/Data | 실제 data flow와 side effect가 맞는가 | API/DB/queue observation |
| L4 System E2E | 시스템을 통해 완료되는가 | trace/video/network |
| L5 Intent/UX | 사용자가 목적을 이해하고 달성하는가 | journey assertions/screenshots |
| L6 Adversarial/Recovery | 실패·권한·retry·rollback이 안전한가 | fault/result/state diff |
| L7 Regression/Drift | 기존 의도와 동작을 깨지 않았는가 | baseline comparison |

모든 feature가 모든 layer를 실행할 필요는 없지만 planner는 각 layer에 `required`, `not-applicable(reason)`, `waived` 중 하나를 기록한다.

## 3. Charter schema

```ts
interface TestCharter {
  charter_id:string; ac_ids:string[]; title:string;
  actor:string; preconditions:string[]; starting_state:DataRef[];
  steps:{action:string; expected_ui?:string; expected_data?:string; observer_ids:string[]}[];
  variants:("happy"|"alternate"|"empty"|"error"|"permission"|"retry"|"recovery")[];
  forbidden_outcomes:string[]; required_layers:string[]; risk:"low"|"medium"|"high";
}
```

## 4. Adapter and observer contracts

```go
type QAAdapter interface {
  Probe(ctx context.Context, project Project) Capability
  Prepare(ctx context.Context, c Charter) (ExecutionPlan, error)
  Execute(ctx context.Context, p ExecutionPlan, sink EvidenceSink) RunResult
  Cleanup(ctx context.Context, p ExecutionPlan) error
}
type Observer interface {
  Before(ctx context.Context, checkpoint Checkpoint) Observation
  After(ctx context.Context, checkpoint Checkpoint) Observation
  Compare(before, after Observation) AssertionResult
}
```

Adapters: project-native test, HTTP, Playwright, interactive browser/Chrome MCP, CLI, DB read-only, log/trace, manual checkpoint. 기존 Playwright suite가 있으면 회귀에 우선 사용하고, browser/Computer Use는 exploratory/UX 확인에 사용한다.

Observer는 UI text/URL/DOM/accessibility, HTTP request/response shape, DB row/state, file, queue event, external stub, logs/traces를 지원한다. production 데이터 mutation은 기본 금지다.

## 5. Evidence manifest

```ts
interface EvidenceManifest {
  run_id:string; sprint_id:string; state_revision:number;
  environment:{name:string; fingerprint:string; secret_env_name?:string};
  tool_versions:Record<string,string>; started_at:string; finished_at:string;
  cases:{case_id:string; ac_ids:string[]; status:string; assertions:Assertion[]; artifacts:Artifact[]}[];
  redaction:{policy_version:string; scan_status:"passed"|"failed"};
}
interface Artifact { id:string; kind:string; uri:string; sha256:string; size:number; created_at:string; }
```

Artifact는 content hash와 생성 tool을 가진다. screenshot 단독은 data flow를 증명하지 않으므로 AC observable에 필요한 observer chain을 모두 연결한다.

## 6. Independent evaluator

Evaluator input은 PRD/AC revision, charter, capability snapshot, evidence manifest, allowed waiver다. builder summary는 제외한다. 판정:

- `passed`: observable/forbidden outcome을 충분한 evidence가 증명.
- `failed`: 반대 evidence 또는 assertion 실패.
- `insufficient`: evidence 누락/오염/stale.
- `not-applicable`: 승인된 근거와 scope 존재.

LLM judge를 쓰는 subjective UX 항목은 rubric, cited evidence, confidence를 요구한다. deterministic assertion을 LLM 판단으로 뒤집을 수 없다.

## 7. Gate algorithm

```text
for each blocking AC:
  reject if no charter
  reject if required layer unresolved
  reject if verdict != passed
  reject if evidence hash/redaction/freshness invalid
reject if open blocker gap or expired waiver
otherwise pass
```

“100%”는 blocking AC 전부가 위 조건을 통과한다는 뜻이다. non-blocking 결과는 점수와 debt로 표시할 수 있지만 blocker를 상쇄하지 않는다.

## 8. Loop policy

QA fail은 defect를 생성하고 root cause를 `requirements`, `design`, `implementation`, `test`, `environment`, `policy`로 분류한다. requirements/design이면 해당 문서를 revision하고 downstream impact를 무효화한다. implementation이면 do로, trace/evidence 누락이면 loop-check로 돌아간다. 기본 자동 반복 한도는 3회이며 이후 사람에게 blocker와 시도 evidence를 제공한다.

## 9. UX와 data-flow completeness

각 journey는 시작점, 사용자 행동, 화면/상태 전이, API/command boundary, data mutation, confirmation/error feedback, 최종 observable을 가진다. 최소 변형은 happy, validation, empty, permission, downstream failure, retry/recovery다. UI가 없어도 CLI/API actor journey로 같은 모델을 사용한다.

