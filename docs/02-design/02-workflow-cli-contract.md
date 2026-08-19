# Workflow and CLI Contract

## 1. 상태 전이

| From | To | 필수 guard |
|---|---|---|
| draft | prd | Sprint identity, scope 존재 |
| prd | plan | confirmed intent, actor/outcome/AC/policy/open question |
| plan | design | task/dependency/verification/decision 존재 |
| design | do | interface/data/error/test 설계, 승인 정책 충족 |
| do | loop-check | task 결과와 changed artifact 등록 |
| loop-check | do | repairable gap 존재 |
| loop-check | qa | blocking spec/code gap 0 |
| qa | do/loop-check | failed QA가 코드/설계 수정 요구 |
| qa | report | 모든 blocking AC pass, evidence hash valid |
| report | archived | report 완성, deferred/waiver 명시, 승인 정책 충족 |

단계 skip은 허용하지 않는다. light profile은 문서가 합쳐질 뿐 논리 phase와 guard event는 유지한다.

## 2. Command tree

```text
tene-workflow init [--profile]
tene-workflow status [--json]
tene-workflow sprint create|start|list|show|archive
tene-workflow phase show|transition <phase> [--dry-run]
tene-workflow task add|start|complete|block|defer|list
tene-workflow intent capture|confirm|revise|deprecate|list
tene-workflow document scaffold|validate|sync
tene-workflow graph build|trace|impact|validate
tene-workflow context build [--phase] [--budget]
tene-workflow loop check|record-gap|resolve-gap
tene-workflow qa plan|run|evaluate|status
tene-workflow evidence register|verify|list
tene-workflow waiver request|approve|expire|list
tene-workflow report generate|validate
tene-workflow compact|clear|doctor|migrate
```

## 3. Global flags

- `--root <path>`: repository root; default는 위로 탐색.
- `--json`: stdout에 machine JSON만 출력.
- `--expected-revision <n>`: mutation의 optimistic concurrency.
- `--request-id <id>`: retry deduplication.
- `--no-color`, `--quiet`, `--verbose`.

Secret 값과 raw child environment를 받는 flag는 정의하지 않는다.

## 4. Response envelope

```json
{
  "ok": false,
  "schema_version": "1.0.0",
  "request_id": "req_...",
  "revision": 17,
  "result": null,
  "warnings": [],
  "errors": [{"code":"WF_GUARD_FAILED","message":"...","details":{"guards":["AC_UNVERIFIED"]},"remediation":"Run qa plan"}]
}
```

stdout은 result 전용, diagnostics는 stderr다. `--json`에서도 secret/error raw payload를 포함하지 않는다.

## 5. Exit codes

| Code | 의미 |
|---:|---|
| 0 | 성공/pass |
| 2 | 사용법/validation 오류 |
| 3 | guard 미충족 또는 QA fail |
| 4 | conflict/lock/stale revision |
| 5 | dependency/capability 없음 |
| 6 | security policy 위반 |
| 7 | I/O/corruption/migration 실패 |
| 8 | child tool/test 실패 |
| 10 | 예상하지 못한 internal 오류 |

## 6. Idempotency

- read command는 무조건 idempotent다.
- mutation은 `request-id`가 같으면 최초 response를 돌려준다.
- `transition` target이 이미 현재 phase이면 no-op success다.
- artifact/evidence 등록은 content hash로 deduplicate한다.
- retry가 의미를 바꿀 경우 conflict를 반환하고 자동 덮어쓰지 않는다.

## 7. Guard finding

```ts
interface Finding {
  code:string; severity:"blocker"|"warning"|"info";
  subject_refs:string[]; message:string; evidence_refs:string[];
  remediation:{command?:string; description:string};
  waivable:boolean;
}
```

`--dry-run`과 실제 전이는 같은 guard list를 호출해야 한다. 시간/외부 상태가 필요한 guard는 evaluation timestamp와 capability snapshot을 고정한다.

