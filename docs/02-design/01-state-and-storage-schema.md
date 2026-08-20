# State and Storage Schema

## 1. 저장 구조

```text
.tene-workflow/
├── project.json
├── active.json
├── events.ndjson
├── master-plan.json
├── policies.yaml
├── graph/{nodes.ndjson,edges.ndjson,index.json}
├── evidence/<run-id>/{manifest.json,artifacts...}
├── cache/                     # clear 가능
├── backups/
└── .lock
docs/sprints/
├── <sprint-id>-<slug>/...
└── _archive/YYYY-MM/<sprint-id>-<slug>/...
```

`.tene/`는 secret vault 영역이며 이 제품이 읽거나 이동하지 않는다.

## 2. ID와 공통 envelope

- ID: `project_`, `sprint_`, `task_`, `intent_`, `ac_`, `journey_`, `component_`, `run_`, `evidence_`, `waiver_` + ULID.
- timestamp: RFC 3339 UTC, 표시 단계에서 locale 변환.
- revision: uint64 monotonic.
- persisted JSON은 UTF-8, stable key ordering, trailing newline.

```json
{
  "schema_version": "1.0.0",
  "revision": 42,
  "updated_at": "2026-08-20T03:00:00Z",
  "updated_by": {"kind": "codex", "session_id": "opaque"}
}
```

## 3. Project/Sprint/Task

```ts
type Phase = "draft"|"prd"|"plan"|"design"|"do"|"loop-check"|"qa"|"report"|"archived";
type Profile = "strict"|"standard"|"light"|"off";
interface ProjectState {
  project_id: string; profile: Profile; active_sprint_id?: string;
  sprints: Record<string,SprintSummary>; revision: number;
}
interface Sprint {
  sprint_id: string; slug: string; title: string; phase: Phase;
  predecessor_ids: string[]; intent_ids: string[]; task_ids: string[];
  started_at?: string; archived_at?: string; document_root: string;
  approvals: ApprovalRef[]; open_gap_ids: string[];
}
interface Task {
  task_id: string; title: string; status: "todo"|"doing"|"blocked"|"done"|"deferred";
  intent_ids: string[]; ac_ids: string[]; depends_on: string[];
  layer: "interface"|"business"|"persistence"|"infrastructure";
  owner?: string; artifacts: ArtifactRef[];
}
```

## 4. Intent model

```ts
interface Intent {
  intent_id: string; revision: number;
  status: "candidate"|"confirmed"|"superseded"|"deprecated";
  statement: string; rationale: string; actors: string[];
  desired_outcomes: string[]; non_goals: string[]; policies: PolicyRef[];
  source: {kind:"conversation"|"document"|"user"; locator:string};
  confirmed_by?: string; confirmed_at?: string; supersedes?: string;
}
interface AcceptanceCriterion {
  ac_id:string; intent_id:string; statement:string;
  priority:"blocking"|"non-blocking"; observable:string;
  preconditions:string[]; expected:string[]; forbidden:string[];
}
```

Intent 의미 변경은 in-place rewrite가 아니라 revision event 또는 supersede로 표현한다.

## 5. Event journal

```ts
interface Event {
  sequence:number; event_id:string; event_type:string; aggregate_id:string;
  occurred_at:string; actor:Actor; expected_revision:number;
  payload:object; previous_hash:string; hash:string;
}
```

`hash = sha256(canonical(event without hash))`. reducer는 pure function이어야 하고 journal replay 결과가 `active.json`과 일치해야 한다.

주요 event: `ProjectInitialized`, `SprintCreated`, `PhaseTransitioned`, `TaskChanged`, `IntentCaptured`, `IntentConfirmed`, `GraphDeltaApplied`, `EvidenceRegistered`, `VerdictRecorded`, `WaiverGranted`, `SprintArchived`, `ProjectionCompacted`.

## 6. Retention classes

- **source**: project, events, authored docs, evidence manifest. clear 금지.
- **derived**: active projection, graph index, context packs. 재생성 가능.
- **ephemeral**: locks, temp, browser trace working copy. clear 가능.
- **large evidence**: policy에 따라 외부 durable store로 승격 가능하나 manifest와 hash는 유지.

`compact`는 journal snapshot을 만들고 기존 segment를 checksummed archive로 옮긴다. `clear`는 derived/ephemeral만 제거한다. `archive`는 Sprint 문서 경로를 이동하고 immutable flag/event를 남긴다.

구현된 replay contract에서 `compact`는 현재 revision의 full `ProjectionCheckpoint`와 checksummed snapshot을 남긴다. 이후 event는 원래 domain payload와 canonical projection merge patch를 함께 hash-chain에 포함한다. `doctor`는 최신 checkpoint부터 patch를 replay하여 세 projection을 비교하고, `--repair`는 timestamped backup 후 replay 결과만 기록한다. Journal 자체가 손상되었거나 checkpoint 이후 patch가 누락되면 fail closed 한다.

## 7. Atomicity

mutation은 repo-scoped advisory lock을 획득하고 expected revision을 비교한다. 파일은 같은 filesystem의 temp에 쓴 뒤 flush/fsync/rename한다. lock timeout 기본 5초, stale lock은 PID/host/time 확인 뒤 `doctor`만 정리한다.
