# System Architecture

## 1. Context

Codex skill은 대화와 코드 조사를 잘하지만 장기 상태와 결정론적 gate의 source of truth가 되어서는 안 된다. `tene-workflow` core가 상태와 판정을 소유하고, Codex가 해석·조사·도구 실행을 맡는다.

```text
User / Natural language / $tene:*
             │
      Codex Skill Router
             │ Context Pack + Policy
     Phase Skill / QA Orchestrator
             │ structured JSON calls
        tene-workflow CLI
   ┌─────────┼─────────┬──────────┐
 State   Documents   Graph      Evidence
 Engine  Validator   Engine     Evaluator
   │         │          │           │
.tene-workflow/   docs/sprints/   Test/Browser/Git
                                      │
                         tene run (secret child only)
```

## 2. Components

| Component | 책임 | 입력 | 출력/변경 |
|---|---|---|---|
| `app` | command use case 조립 | parsed request | response DTO |
| `config` | profile/policy resolve | project/user defaults | effective config |
| `state` | event append/projection/lock | domain event | durable state |
| `workflow` | transition/guard | state + evidence | decision/events |
| `document` | scaffold/validate/generated region | template + domain | Markdown/issues |
| `intent` | candidate/confirm/revise | conversation extraction | intent events |
| `graph` | spec/code/evidence 연결 | nodes/edges/providers | trace/impact |
| `context` | phase별 최소 context | active state + graph | context pack |
| `qa` | charter/adapter/evaluate | AC + capabilities | runs/verdict |
| `secret` | tene metadata/run policy | secret refs + command | sanitized result |
| `evidence` | manifest/hash/redaction | artifacts | immutable metadata |

## 3. Runtime flows

### Session resume

1. skill이 repo root와 `.tene-workflow/project.json`을 확인한다.
2. `tene-workflow status --json`이 active Sprint, phase, open task/gap, revision을 반환한다.
3. `context build --phase <phase> --json`이 provenance가 있는 pack을 만든다.
4. Codex는 next actions를 설명하고, mutation 전 필요한 승인만 요청한다.

### Phase work

1. skill이 현재 phase와 user intent가 일치하는지 확인한다.
2. 문서/코드 변경 전 `task start --expected-revision`을 기록한다.
3. Codex가 작업하고 `artifact register`, `task complete`로 결과를 연결한다.
4. `transition --dry-run`에서 guard를 검사하고 pass 후 실제 전이한다.

### QA

1. graph에서 blocking AC와 journeys/data paths를 읽어 charter를 생성한다.
2. capability discovery로 deterministic adapter를 우선 선택한다.
3. 실행 결과를 observer가 artifact로 만들고 evidence manifest에 hash와 연결한다.
4. 별도 evaluator context가 manifest를 AC별로 판정한다.
5. fail/gap이면 loop-check로 되돌아가고 pass면 report로 전이한다.

## 4. Dependency rules

- domain package는 filesystem, Codex, browser를 import하지 않는다.
- adapters가 port interface를 구현한다. core에서 특정 MCP 이름을 hard-code하지 않는다.
- skill은 state 파일을 직접 수정하지 않고 CLI만 호출한다.
- 문서 수동 편집은 허용하되 다음 CLI 호출에서 parse/validate 후 event로 반영한다.
- plugin correctness는 hook, network, remote service 존재에 의존하지 않는다.

## 5. Package interfaces

```go
type StateStore interface {
  Load(ctx context.Context) (ProjectState, error)
  Append(ctx context.Context, expected Revision, events ...Event) (Revision, error)
}
type Guard interface { Evaluate(Context) []Finding }
type GraphProvider interface { Capabilities() CapabilitySet; Extract(ctx context.Context, Scope) (GraphDelta, error) }
type QAAdapter interface { Probe(context.Context) Capability; Plan(Charter) PlanResult; Run(context.Context, TestCase) RunResult }
type SecretRunner interface { ListNames(context.Context, string) ([]SecretRef, error); Run(context.Context, SecretRunRequest) SanitizedResult }
```

모든 public result는 `schema_version`, `request_id`, `revision`, `warnings`를 포함한다.

