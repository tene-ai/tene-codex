# Work Breakdown and Dependencies

## 1. 목표 repository layout

```text
tene-codex/
├── .codex-plugin/plugin.json
├── skills/tene-*/SKILL.md
├── skills/tene-*/agents/openai.yaml
├── references/
├── scripts/
├── cmd/tene-workflow/main.go
├── internal/{app,config,state,workflow,document,intent,graph,context,qa,secret,evidence}/
├── schemas/
├── templates/
├── testdata/
├── evals/
└── docs/
```

## 2. Work package

| ID | 패키지 | 주요 산출물 | 선행 | 검증 |
|---|---|---|---|---|
| WP-01 | Domain contract | IDs, enums, schema, errors | 없음 | schema fixtures |
| WP-02 | State store | journal, projection, lock, atomic write | WP-01 | crash/concurrency tests |
| WP-03 | Workflow | transition guards, approvals, profiles | WP-01,02 | transition table tests |
| WP-04 | Documents | scaffold, frontmatter/section validation | WP-01,03 | golden templates |
| WP-05 | Intent | capture, confirm, revise, supersede | WP-01,02 | conversation fixtures |
| WP-06 | Graph | node/edge store, trace, impact, orphan | WP-01,05 | graph invariant tests |
| WP-07 | Context | phase pack builder, budgets, freshness | WP-02,06 | snapshot/token tests |
| WP-08 | Code intel | git/static/CodeGraph adapters, 4L/6Q | WP-06 | polyglot fixtures |
| WP-09 | Loop | spec/code comparison, gaps, iteration | WP-04,06,08 | mutation/gap tests |
| WP-10 | QA | charters, adapters, observers, evaluator | WP-05,06,08 | reference journeys |
| WP-11 | Secret | metadata/run adapters, redaction | WP-02,10 | canary leak tests |
| WP-12 | Skills | 9 skill contracts/router/references | WP-03~11 | skill eval suite |
| WP-13 | Plugin | manifest, install, compatibility | WP-12 | plugin validator |
| WP-14 | Release | binaries, migration, docs, marketplace | WP-13 | clean-room matrix |

## 3. 세부 task와 Definition of Done

### WP-01 Domain contract

- 모든 ID는 `<type>_<ulid>` 형식이며 생성 후 바뀌지 않는다.
- schema version은 각 persisted document에 기록한다.
- enum과 required field를 Go type과 JSON Schema 양쪽에서 생성/검증한다.
- unknown field 정책은 persisted state에는 reject, user-authored Markdown frontmatter에는 preserve+warn이다.

### WP-02 State store

- append-only `events.ndjson`가 source of truth, `active.json`은 재생 가능한 projection이다.
- write는 lock → revision 확인 → temp write → fsync → rename → directory fsync 순서다.
- event hash chain과 monotonic sequence로 truncation/tampering을 탐지한다.
- `doctor --repair`는 원본을 백업한 뒤 projection만 재생성한다.

### WP-03 Workflow

- transition table을 data-driven rule로 구현한다.
- 각 guard는 stable code, 설명, remediation을 반환한다.
- `--dry-run`은 동일한 guard를 실행하되 mutation하지 않는다.
- force는 없고, waiver는 이유·scope·approver·expiry를 갖는 별도 domain event다.

### WP-04 Documents

- template은 모든 공통 section과 문서별 section을 포함한다.
- section marker는 heading text가 아니라 stable HTML comment ID로 검증한다.
- 자유 섹션은 보존한다.
- generated block만 명시적 marker 사이에서 갱신하며 사용자 서술을 덮어쓰지 않는다.

### WP-05~07 Intent, graph, context

- intent 후보는 저장 가능하지만 `confirmed` 전에는 build gate input이 아니다.
- edge는 source locator와 confidence를 가진다.
- context pack은 phase allowlist, relevance, recency, token budget 순으로 구성한다.
- summary에는 원문 artifact ID와 revision을 포함해 검증 가능하게 한다.

### WP-08~10 Code, loop, QA

- provider가 불완전해도 capability를 보고하고 fallback을 사용한다.
- loop evaluator는 missing, mismatch, unverified, regression, debt로 gap을 분류한다.
- QA planner는 AC마다 happy/alternate/error/recovery와 observer를 연결한다.
- evaluator는 구현 에이전트의 서술이 아니라 evidence manifest를 읽는다.

### WP-11 Secret boundary

- 허용: `tene list --json`, `tene whoami`, `tene run --env ... -- ...`.
- 금지: `.tene/**` read, `tene get`, plaintext export, secret stdout capture.
- child process log는 명령 template과 secret name만 남기며 값은 남기지 않는다.
- leak 발견 시 evidence를 격리하고 QA를 즉시 fail한다.

### WP-12~14 Skills, plugin, release

- skill description은 겹치지 않는 trigger와 negative trigger를 포함한다.
- scripts는 structured JSON을 입출력하고 사람이 해석할 메시지는 skill이 만든다.
- plugin manifest에는 지원되는 field만 둔다. hook은 manifest에 임의 field로 넣지 않는다.
- schema/CLI/plugin compatibility matrix와 rollback path를 릴리스마다 갱신한다.

## 4. 병렬화 가능한 흐름

WP-01 계약 확정 뒤 다음을 병렬 진행할 수 있다.

- state/workflow/document
- intent/graph/context
- QA adapter spike와 secret threat-test fixture
- plugin/skill validator 및 packaging spike

단, 통합 branch에서는 core contract test가 통과한 결과만 합친다. graph나 QA가 자체 ID/schema를 만들지 않도록 WP-01 변경은 architecture decision record를 요구한다.

## 5. 구현 task 공통 체크리스트

- 어떤 intent/AC를 충족하는가?
- 4개 Understanding Layer 중 어디에 속하는가?
- 이름, 정의 파일, import/reference, caller, input, output/mutation이 문서화됐는가?
- 실패, retry, rollback, idempotency가 설계됐는가?
- secret/PII/logging 영향이 있는가?
- unit 외에 어떤 integration/journey evidence가 필요한가?
- migration과 이전 schema 독자가 필요한가?

