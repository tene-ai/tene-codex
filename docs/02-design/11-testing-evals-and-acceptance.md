# Testing, Evals and Acceptance Design

## 1. Test architecture

```text
internal/* unit tests
      ↓
schema/transition/golden contract tests
      ↓
temp-repo integration tests
      ↓
staged-plugin skill/system evals
      ↓
reference-project full Sprint journeys
```

테스트는 구현 함수가 아니라 public behavior와 invariant를 우선 고정한다.

## 2. Fixture sets

```text
testdata/
├── schemas/{valid,invalid,previous-version}/
├── transitions/{allowed,blocked,waived}/
├── documents/{minimal,expanded,malformed,manual-edits}/
├── graphs/{webapp,monolith,polyglot,partial-provider}/
├── evidence/{valid,missing,stale,tampered,secret-canary}/
├── conversations/{ko,en,ambiguous,negative,collision}/
└── repositories/{greenfield,monolith,polyglot}/
```

## 3. Core contract tests

- 모든 상태 전이 pair를 table-driven test로 검사한다.
- event replay, projection snapshot, compact 전후 semantic equality를 property test한다.
- random process interruption 시 atomicity와 journal prefix validity를 확인한다.
- 두 writer의 stale revision 중 정확히 하나만 commit되는지 검사한다.
- document generator가 freeform/manual content를 보존하는지 golden diff로 검사한다.
- graph traversal이 orphan, cycle, supersede, confidence를 올바르게 처리하는지 검사한다.

## 4. QA/evidence mutation tests

정상 reference feature에서 다음을 하나씩 변형하고 gate가 fail하는지 본다.

- AC 하나 삭제
- UI만 바꾸고 API contract 미반영
- API success지만 DB write 제거
- permission check 제거
- stale screenshot 재사용
- evidence hash 변경
- test output에 secret canary 삽입
- builder가 report에서 pass라고 주장하나 evidence 없음

이 suite가 false-pass 방지의 핵심 회귀 테스트다.

## 5. Skill eval schema

```yaml
- id: route-ko-qa-001
  user: "회원가입 기능이 기획대로 화면과 DB까지 동작하는지 종합 테스트해줘"
  state: {phase: qa, active_sprint: true}
  expect:
    skill: qa
    commands: ["status --json", "context build", "qa plan"]
    forbidden: ["tene get", "phase skip"]
    needs_confirmation: false
```

평가 항목: routing, state awareness, required step, forbidden action, artifact path, gate honesty. 한국어/영어 paraphrase와 negative example을 균형 있게 둔다.

## 6. End-to-end acceptance scenarios

### Scenario A — New feature

자연어 요청 → intent candidate → user confirmation → plan/design → implementation → loop gap repair → Playwright+API+DB QA → report → archive. 새 세션 두 번을 사이에 넣어 resume를 검증한다.

### Scenario B — Changed policy

archived Sprint의 정책을 변경 → new intent revision → impact graph가 기존 design/test를 표시 → regression QA → predecessor 연결 report.

### Scenario C — Secret-required E2E

key name metadata 확인 → `tene run` child test → secret canary scan → value 없는 evidence/report. 권한 거부/미설치/child failure도 검증한다.

### Scenario D — Incomplete provider

CodeGraph/AST가 없는 언어 → fallback과 unknown Six Questions → 사람 확인 또는 waiver → 과신하지 않는 report.

## 7. Product AC executable mapping

| Product AC | System assertion |
|---|---|
| 01 Mandatory cycle | phase skip returns exit 3 and guard codes |
| 02 Durable intent | new session produces same active intent/revision |
| 03 Whole-system | changed component has 4L/6Q coverage or explicit unknown gap |
| 04 Evidence 100% | every blocking AC has valid passed evidence |
| 05 Full-flow QA | journey has UI/API/data observations and recovery variant |
| 06 Secret safety | canary appears zero times outside child boundary |
| 07 Portable | 3 repository types complete same public workflow |
| 08 Report continuity | report links predecessor, intent, files, evidence, deferred work |

## 8. Performance and reliability budgets

- local `status`: p95 < 200 ms on 10k events with projection.
- transition guard: p95 < 1 s excluding external tools.
- context build: p95 < 3 s with warm graph index.
- journal recovery: 100k events < 10 s target.
- routing eval deterministic variance: 0 for explicit, tracked for implicit model versions.
- QA duration is project-dependent; every adapter supports timeout/cancel/cleanup.

## 9. Definition of Done for implementation

- code, schema, migration, docs, fixtures, observability가 함께 변경되었다.
- 관련 PRD ID와 AC가 task/commit/report에 연결되었다.
- 4 Layers와 6 Questions가 채워졌거나 evidence 있는 N/A/unknown이다.
- unit/contract/integration과 영향받은 system/skill eval이 통과했다.
- security/canary와 backward compatibility 결과가 있다.
- open blocker 0, non-blocking debt는 owner/reason/target Sprint가 있다.

