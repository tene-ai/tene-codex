# Graph and Context Engine

## 1. Unified graph

### Node kinds

`Intent`, `AcceptanceCriterion`, `Policy`, `Journey`, `Sprint`, `Task`, `DocumentSection`, `Layer`, `File`, `Symbol`, `DataShape`, `TestCase`, `Run`, `Evidence`, `Gap`, `Decision`, `Waiver`.

```ts
interface Node {
  id:string; kind:string; label:string; revision:number;
  locator?:string; attributes:Record<string,unknown>;
  source:"authored"|"derived"|"observed";
  confidence:number; content_hash?:string;
}
interface Edge {
  id:string; from:string; to:string;
  kind:"realizes"|"verifies"|"depends_on"|"calls"|"imports"|"reads"|"writes"|"transitions_to"|"observed_by"|"supersedes"|"waived_by"|"belongs_to";
  source_locator:string; confidence:number; provider:string; revision:number;
}
```

### Invariants

- confirmed intent에는 최소 1개 AC가 있다.
- blocking AC에는 design/task/test/evidence 경로가 있다.
- `passed` verdict에는 hash-valid evidence가 있다.
- archived Sprint의 authored node는 수정하지 않고 supersede한다.
- confidence가 낮은 inferred edge는 확정 사실처럼 보고하지 않는다.

## 2. Providers

Provider 우선순위:

1. project가 이미 가진 CodeGraph 등 semantic index
2. language-native compiler/LSP/test coverage
3. git diff + AST/static analyzer
4. text search
5. user-confirmed manual mapping

```go
type Capability struct { Languages []string; Calls, Imports, DataFlow, Runtime bool; Confidence float64 }
type GraphDelta struct { Nodes []Node; Edges []Edge; Diagnostics []Finding; BaseRevision uint64 }
```

`.codegraph/`가 있으면 CodeGraph provider를 먼저 사용한다. 없으면 indexing을 자동 강제하지 않고 다른 provider로 fallback한다.

## 3. Understanding Layer classifier

규칙과 semantic evidence를 결합한다.

- Interface: route/controller/CLI/UI/event consumer/scheduler entry.
- Business Logic: usecase/service/handler/reducer/domain rule.
- Persistence: repository/DB/file/cache/queue/external API adapter.
- Infrastructure: server/container/auth middleware/cloud/CI/CD/config runtime.

한 component가 여러 layer를 가질 수 있으나 primary layer와 근거를 둔다. 단순 경로명 추정은 confidence ≤ 0.5이며 report에 확인 필요로 표시한다.

## 4. Six Questions materializer

각 changed public symbol과 AC 핵심 component에 대해 다음 query를 실행한다.

1. declaration node/name
2. definition locator
3. incoming `imports`/`references`
4. incoming/outgoing `calls`/`uses`
5. input `DataShape`
6. output `DataShape`와 `writes`/side effects

답이 없다는 것도 결과다. provider가 지원하지 않아 unknown인지 실제 사용처가 없어 orphan인지 구분한다.

## 5. Impact and gap

`impact <node>`는 directed traversal을 하되 edge kind별 최대 깊이를 정책으로 둔다. 예: policy→AC→journey/design/task/test는 전부, call graph는 기본 4 hops. changed node에서 impacted AC까지 경로가 없으면 “untraced change” gap이다.

gap categories:

- `missing`: 필요한 artifact/component가 없음
- `mismatch`: spec과 구현 의미가 다름
- `unverified`: 구현은 있으나 evidence 없음
- `regression`: 이전 pass baseline이 실패
- `debt`: non-blocking 구조 품질 문제

## 6. Context pack

```ts
interface ContextPack {
  id:string; phase:Phase; state_revision:number; budget:number;
  objective:string; active_task?:Task; confirmed_intents:Intent[];
  relevant_ac:AcceptanceCriterion[]; policies:Policy[];
  graph_slices:GraphPath[]; open_gaps:Gap[]; recent_decisions:Decision[];
  tool_capabilities:Capability[]; provenance:SourceRef[];
  excluded_summary:{reason:string; count:number}[];
}
```

선택 순서: safety/policy → active intent/AC → open blockers → phase artifacts → impacted graph slice → recent decisions → optional history. 동일 내용은 content hash로 dedupe한다. summary마다 source revision/locator를 둔다.

기본 budget 비율은 instruction/policy 15%, intent/spec 25%, active work 20%, graph/code 25%, evidence/gaps 10%, reserve 5%다. 초과 시 오래된 narrative부터 축약하며 blocking AC와 safety rule은 제거하지 않는다.

## 7. Freshness

context pack의 revision과 file hashes가 작업 시작 시점과 다르면 mutation 전 rebuild한다. archived history는 predecessor link와 요약만 기본 포함하고 필요할 때 원문을 확장한다. `compact`는 source를 지우지 않고 context summary와 projection을 재생성한다.

