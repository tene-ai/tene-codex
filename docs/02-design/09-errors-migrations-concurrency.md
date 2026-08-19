# Errors, Recovery, Migrations and Concurrency

## 1. Stable error namespace

- `CFG_*`: 설정/정책
- `STATE_*`: journal/projection/lock
- `WF_*`: phase/guard/approval
- `DOC_*`: template/frontmatter/marker
- `GRAPH_*`: provider/trace/invariant
- `QA_*`: capability/run/evidence/verdict
- `SEC_*`: forbidden action/leak/tene dependency
- `MIG_*`: schema migration
- `PLUGIN_*`: discovery/version

오류는 code, human message, structured details, remediation, retryable을 포함한다. path는 repo-relative로 표시하고 환경/secret 값은 넣지 않는다.

## 2. Failure matrix

| Failure | 자동 처리 | 사용자 개입 |
|---|---|---|
| stale revision | state reload/context rebuild, 1회 retry | semantic conflict 시 선택 |
| lock busy | bounded backoff | timeout 후 세션 확인 |
| corrupt projection | journal replay | journal도 손상 시 backup 선택 |
| partial artifact write | temp 제거, 이전 파일 유지 | 없음 |
| missing provider | fallback + capability warning | required quality 미달 시 waiver/설치 |
| flaky QA | 동일 조건 제한 retry, flaky 분리 | blocker waiver는 승인 필요 |
| expired evidence | rerun | 환경 불가 시 승인된 waiver |
| migration fail | backup restore | report와 수동 해결 |

## 3. Concurrency

- repo 당 mutation lock 하나로 시작한다. read는 lock-free이나 revision을 반환한다.
- 모든 mutation은 expected revision을 요구한다. interactive CLI가 생략하면 직전 load revision을 내부 사용한다.
- lock 안에서 외부 test/browser/LLM을 실행하지 않는다. `run planned` event 후 lock 해제, 실행, 결과 commit 시 revision conflict를 검사한다.
- 동일 request ID는 idempotency registry에서 최초 결과를 반환한다.
- 서로 다른 Sprint 병렬 작업은 post-MVP fine-grained lock 후보지만 shared graph/master plan commit은 여전히 직렬화한다.

## 4. Recovery

`doctor` 단계:

1. filesystem/layout/schema/version 검사
2. journal hash/sequence 검사
3. projection replay와 diff
4. document IDs/revisions/links 검사
5. graph orphan/index 검사
6. evidence path/hash/redaction status 검사
7. Codex/tene/tool capabilities 검사

`doctor --repair`는 dry-run이 기본이며 `--apply` 때 timestamped backup 후 derived data만 재생성한다. authored docs와 source journal의 자동 삭제/재작성은 금지한다.

## 5. Schema migration

각 migration은 `From`, `To`, `Plan`, `Apply`, `Verify`, `RollbackHint`를 구현한다. major version을 건너뛰지 않고 순차 적용한다.

```text
detect → compatibility check → migration plan JSON
       → dry-run semantic diff → backup
       → locked apply → replay/validate → MigrationApplied event
```

unknown future major schema는 read/write를 거절하되 raw 파일을 보존한다. migration 중 secret vault는 대상에 포함하지 않는다.

## 6. Waiver

```ts
interface Waiver {
  waiver_id:string; subject_refs:string[]; finding_codes:string[];
  reason:string; compensating_controls:string[];
  requested_by:string; approved_by:string; approved_at:string;
  expires_at?:string; expires_after_sprint?:string;
}
```

Security invariant, corrupted evidence, missing confirmed intent는 non-waivable이다. waiver는 문제를 pass로 바꾸지 않고 gate에서 승인된 예외로 별도 표시한다.

