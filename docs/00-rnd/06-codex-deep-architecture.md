# Codex 내부 아키텍처·프로토콜·확장 포인트 심층 분석

이 문서는 tene 개발자가 Codex를 단순 CLI가 아니라 **model + harness + execution environment + extension surfaces + clients**로 이해하기 위한 구현 기준서다. 공개 문서와 `openai/codex` 저장소에서 확인되는 사실을 중심으로 작성했으며, 추론은 별도로 표시한다.

## 1. Codex는 하나의 프로그램이 아니라 공통 harness를 사용하는 제품군이다

Codex CLI, IDE extension, desktop app, cloud/web는 사용자 경험은 다르지만 핵심 agent logic을 공유한다. OpenAI는 이를 Codex harness 또는 Codex core라고 설명한다.

```text
CLI TUI ───────────┐
IDE extension ─────┤
Desktop app ───────┼── App Server / client adapters ── Codex Core
Cloud/Web ─────────┤                                  ├─ thread runtime
Embedded product ──┘                                  ├─ agent loop
                                                       ├─ prompt/context assembly
                                                       ├─ tool execution
                                                       ├─ sandbox/approval policy
                                                       ├─ auth/config
                                                       └─ persistence/events
```

App Server는 Codex core를 client-friendly bidirectional JSON-RPC로 노출하는 장기 실행 프로세스다. 프로세스 내부에는 stdio reader, message processor, thread manager, thread별 core session이 있다. 이는 tene가 rich UI 또는 별도 orchestrator를 만들 때 CLI stdout parsing보다 App Server를 우선 검토해야 한다는 뜻이다.

근거: [Unlocking the Codex harness](https://openai.com/index/unlocking-the-codex-harness/), [Codex App Server 문서](https://developers.openai.com/codex/app-server), [App Server source README](https://github.com/openai/codex/blob/main/codex-rs/app-server/README.md)

## 2. Codex core의 agent loop

공개된 loop는 다음과 같다.

1. 사용자 입력과 환경·정책·instructions를 prompt item으로 조립한다.
2. Responses API에 inference를 요청한다.
3. 모델이 final assistant message 또는 tool call을 생성한다.
4. tool call이면 harness가 policy를 평가하고 실행한다.
5. tool result를 conversation input에 추가해 다시 inference한다.
6. 모델이 assistant message를 반환하면 turn이 종료된다.

한 user message에서 assistant message까지가 turn이고, 여러 turn이 thread를 이룬다. 한 turn 내부에는 수십·수백 번의 model↔tool iteration이 있을 수 있다. 따라서 tene workflow의 “단계”를 Codex turn과 동일시하면 안 된다. 하나의 tene phase가 여러 thread/turn을 가질 수 있고, 반대로 한 Codex turn이 여러 workflow task를 수행할 수도 있다.

근거: [Unrolling the Codex agent loop](https://openai.com/index/unrolling-the-codex-agent-loop/)

### tene에 필요한 별도 식별자

```text
tene_project_id
  └─ feature_id
      └─ workflow_run_id
          ├─ phase_id
          ├─ codex_thread_id
          ├─ codex_turn_id[]
          └─ evidence_id[]
```

Codex thread ID만으로 제품 workflow를 식별하지 말아야 한다. thread fork, resume, retry, 다른 agent로의 이동을 흡수할 상위 run ID가 필요하다.

## 3. Prompt와 context 조립

Codex는 단순히 사용자 메시지만 모델에 보내지 않는다. 공개 설명상 initial context에는 다음 종류가 포함된다.

- model/base instructions 또는 `model_instructions_file`
- sandbox와 approval policy를 설명하는 developer message
- current working directory와 환경 정보
- `$CODEX_HOME` 및 repo 경로의 `AGENTS.override.md`/`AGENTS.md`
- model-visible skill의 name/description 목록
- 사용자가 명시한 skill의 전체 지침
- 대화 history와 tool calls/results
- 사용 가능한 tool definitions

AGENTS 파일은 기본 limit의 영향을 받고, skill catalog 역시 context budget 때문에 일부 description이 단축되거나 생략될 수 있다. MCP tool 목록의 순서나 변경도 prompt caching에 영향을 준다. 따라서 tene가 수백 개 skill을 쪼개 배포하거나 매 turn 동적으로 tool schema를 바꾸면 비용과 라우팅 신뢰성이 악화될 수 있다.

### 설계 결과

- tene의 public skill 수는 사용자 intent 기준 5~8개 정도로 제한한다.
- 세부 workflow step은 skill을 계속 늘리기보다 내부 reference와 state machine으로 라우팅한다.
- stable prefix를 유지하도록 정적 지침과 tool 정의 순서를 안정화한다.
- active spec neighborhood는 매번 full injection하지 않고 path/index를 제공한 뒤 필요한 artifact를 읽게 한다.
- 긴 hook `additionalContext`에 핵심 memory를 전부 밀어 넣지 않는다.

관련 운영 위험은 [hook additionalContext spill 이슈](https://github.com/openai/codex/issues/22861)와 [compaction guidance 이슈](https://github.com/openai/codex/issues/29816)에 나타난다. 이슈는 제품 계약이 아니므로 제한 수치 자체를 hard-code하면 안 된다.

## 4. Context compaction과 durable memory

긴 thread는 context window를 소모한다. Codex는 threshold를 넘으면 Responses API의 compaction 기능을 사용한다. 공개 설명에 따르면 compact response는 이전 input을 대체할 item 목록과 opaque compaction item을 반환해 latent understanding을 보존한다.

그러나 compaction은 application-level durable state와 동일하지 않다.

- 어떤 세부 결정이 보존되는지 tene가 검증할 수 없다.
- AGENTS 규칙이나 task progress가 요약에서 약화될 수 있다는 사용자 보고가 있다.
- hook 기반 재주입도 timing과 size에 의존할 수 있다.
- thread를 다른 host 또는 agent로 옮길 때 opaque state의 이식성을 기대할 수 없다.

따라서 tene의 memory 계층은 아래처럼 분리한다.

| 계층 | 수명 | 저장 대상 |
|---|---|---|
| Model context | inference/turn | 현재 판단에 필요한 최소 context |
| Codex thread | session/thread | transcript, tool events, compacted history |
| tene run state | workflow run | phase, task, checkpoint, retry, budget |
| Project memory | repository lifetime | confirmed specs, ADR, journeys, invariants |
| Evidence archive | policy-defined | screenshots, logs, traces, DB assertions |

핵심 원칙은 **“compaction이 실패해도 `.tene/`에서 다음 action을 재구성할 수 있어야 한다”**이다.

## 5. App Server 프로토콜 모델

App Server는 thread, turn, item의 계층을 이벤트로 전달한다.

```text
thread/started
  turn/started
    item/started
    item/.../delta        (0..n)
    item/completed
    turn/plan/updated     (0..n)
    turn/diff/updated     (0..n)
  turn/completed          completed | interrupted | failed
thread/archived | thread/closed
```

중요한 특성:

- `thread/read`, `resume`, `fork`가 지속 가능한 thread history를 다룬다.
- turn은 상태와 item 목록을 갖는다.
- item은 agent message, command/tool execution, file change 등 세부 lifecycle을 가진다.
- `turn/plan/updated`는 step과 pending/inProgress/completed 상태를 스트리밍한다.
- `turn/diff/updated`는 turn 전체 unified diff snapshot을 제공한다.
- approval은 server request/client response 흐름으로 다뤄질 수 있으므로 client가 응답하지 않으면 실행이 정지할 수 있다.
- 일부 realtime/raw/safety event는 transient이며 persisted ThreadItem으로 취급하면 안 된다.

### tene App Server adapter

tene가 desktop/dashboard/orchestrator를 제공할 경우 다음 adapter가 필요하다.

```ts
interface CodexRunAdapter {
  startThread(context: ContextPack): Promise<ThreadRef>;
  startTurn(thread: ThreadRef, task: TaskEnvelope): Promise<TurnRef>;
  streamEvents(turn: TurnRef): AsyncIterable<CodexEvent>;
  respondToApproval(request: ApprovalRequest): Promise<void>;
  interrupt(turn: TurnRef): Promise<void>;
  resume(threadId: string): Promise<ThreadRef>;
  fork(threadId: string): Promise<ThreadRef>;
}
```

App Server schema는 생성 가능한 TypeScript/JSON Schema를 제공하므로 문자열 event name을 임의 구현하기보다 pinned Codex version에서 schema를 생성하고 compatibility test를 둬야 한다.

## 6. 실행 안전성: sandbox, approval, policy

Codex의 tool execution은 모델의 요청을 바로 OS에서 수행하는 구조가 아니다. harness가 sandbox와 approval policy 아래서 shell/file/MCP tool을 실행한다.

tene는 세 가지 권한을 구분해야 한다.

1. Read/observe: source, logs, DOM, test DB 조회
2. Mutate workspace: code/spec/test/evidence manifest 변경
3. External/destructive: 배포, production data, 메시지 발송, 결제 등

QA executor가 production-like UI를 조작할 때 click 하나가 실제 mutation일 수 있다. tool 이름만 보고 위험을 분류하지 말고 environment, actor, target resource, reversibility를 포함한 action policy가 필요하다.

```yaml
action_policy:
  environment: ephemeral-test
  allowed_effects: [create_test_user, create_test_order]
  forbidden_effects: [real_payment, production_delete, external_email]
  evidence_redaction: [access_token, password, pii]
```

## 7. Extension surface 선택

| 필요 | Codex 표면 | tene 적용 |
|---|---|---|
| 저장소의 지속 규칙 | AGENTS.md | tene artifact 위치, 완료·QA 원칙 |
| 반복 가능한 사용자 workflow | Skill | discover/spec/plan/qa/drift UX |
| 설치·배포 | Plugin | skills와 optional MCP/hook 묶음 |
| 기계적 lifecycle gate | Hook | validation/evidence presence 보조 |
| 외부 시스템·중앙 서비스 | MCP | tracker, Figma, logs, DB, spec registry |
| CI one-shot | `codex exec` | drift/QA/report 자동화 |
| 코드에서 thread 실행 | Codex SDK | backend orchestration의 간단한 경로 |
| rich client/다중 thread | App Server | dashboard, approval, streaming, resume |
| 장기 task fleet | Symphony형 orchestrator | issue→workspace→agent supervision |

Plugin은 사용자에게 기능을 전달하는 포장 단위이지 durable workflow engine 자체가 아니다. tene의 상태 머신과 artifact validator는 plugin 밖에서도 실행 가능한 core library/CLI로 두는 편이 안전하다.

## 8. Multi-agent와 orchestration

Codex subagent는 현재 thread의 bounded delegation에 적합하다. 대규모 task fleet에는 상위 orchestrator가 필요하다. OpenAI의 Symphony 사례는 issue tracker를 control plane과 state machine으로 사용하고, issue마다 workspace/agent를 매핑하며 crash/stall 시 재시작한다.

그러나 OpenAI는 agent를 지나치게 고정된 state-machine node로 제한하는 것도 경고한다. 모델은 구현, review 반영, 여러 PR 생성 등 더 큰 단위를 수행할 수 있기 때문이다.

tene에는 **soft state machine + hard gates**가 적합하다.

- Soft: phase 안에서 agent가 탐색·분해·반복 방법을 선택
- Hard: confirmed spec 없이 irreversible implementation 금지
- Hard: invariant failure 시 acceptance 금지
- Hard: evidence/source 없는 requirement promotion 금지
- Soft: 한 phase에서 생성할 task/PR 수는 agent가 제안
- Human gate: scope, conflicting intent, destructive action

근거: [Open-source Codex orchestration: Symphony](https://openai.com/index/open-source-codex-orchestration-symphony/)

## 9. Harness engineering이 가리키는 미래 방향

공식 공개 자료에서 확인 가능한 방향은 다음과 같다.

### 확인된 방향

- 동일 core harness를 CLI/IDE/app/cloud와 embedded client에 재사용
- interactive session을 넘어 long-running, parallel, resumable work 확대
- agent가 repo, UI, logs, metrics, traces를 직접 읽을 수 있게 만드는 agent legibility
- 문서 권고보다 structural lint/test로 architecture invariant를 기계적으로 강제
- issue tracker와 worktree를 이용한 agent fleet orchestration
- model-native harness, sandbox, memory, filesystem tools를 Agents SDK에도 표준화
- 인간은 intent, architecture, taste, feedback loop와 예외 판단에 집중

### 합리적 추론이지만 확정 로드맵은 아님

- plugin과 App Server 생태계의 결합 강화
- richer lifecycle event와 memory restoration primitive
- agent-to-agent review의 일반화
- spec/issue/evidence가 orchestration control plane으로 수렴
- QA가 test script 실행에서 observable user journey 평가로 확장

공개 글의 “coming next”나 GitHub feature request는 출시 약속이 아니다. tene 로드맵은 미출시 기능에 의존하지 않아야 한다.

## 10. tene 구현에 대한 최종 기술 판단

### 하지 말아야 할 것

- 거대한 하나의 SKILL.md에 전체 방법론을 넣기
- hook output을 장기 memory의 유일한 저장소로 사용
- Codex plan status를 제품 task system으로 그대로 사용
- LLM summary를 confirmed spec으로 자동 승격
- plugin manifest에 orchestration logic을 과도하게 결합
- UI screenshot만으로 business flow 통과 판정

### 해야 할 것

- provider-neutral core schema와 CLI 구축
- `.tene/` Git artifacts를 canonical source로 사용
- plugin skills는 core command를 안전하게 호출하는 UX 계층으로 유지
- App Server adapter는 optional advanced runtime으로 분리
- workflow/agent/thread/turn/evidence ID를 별도로 추적
- deterministic validators와 policy gates를 우선 구현
- UI/API/log/data evidence를 correlation ID로 연결
- 실제 Codex 버전별 contract test 운영

