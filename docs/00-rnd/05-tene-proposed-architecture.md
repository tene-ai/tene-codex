# tene 제품·기술 아키텍처 제안

## 제품 정의

tene는 “문서 생성 플러그인”이 아니라 **intent-to-evidence control plane**이다. 사용자의 자연어 의도를 versioned spec으로 만들고, 구현 workflow와 QA oracle로 재사용하며, 실행 증거로 spec·code·test drift를 관리한다.

## 1. 제안 모듈

```text
┌──────────────── Codex / ChatGPT / Claude adapters ────────────────┐
│ skills · commands · hooks · MCP · App Server/SDK                 │
└───────────────────────────┬───────────────────────────────────────┘
                            │
┌──────────────────────── tene core ────────────────────────────────┐
│ Intent Interviewer  → Spec Registry → Context Assembler           │
│         │                  │                │                      │
│         └──────── Workflow Engine ← Impact/Graph Index            │
│                            │                                      │
│                 QA Planner/Executor/Evaluator                     │
│                            │                                      │
│                 Evidence Store + Drift Engine                     │
└───────────────────────────────────────────────────────────────────┘
```

### Intent Interviewer

Goal/actor/rule/state/invariant/open question을 추출한다. 모호성 점수와 contradiction을 계산하고, 개발 전에 영향이 큰 질문만 사용자에게 묻는다.

### Spec Registry

Markdown frontmatter 또는 YAML을 canonical storage로 사용한다. immutable ID, version, source, status, supersedes 관계를 강제한다.

### Graph Index

spec 문서에서 relation을 파싱해 SQLite 또는 embedded graph로 materialize한다. code symbol graph는 CodeGraph/LSP/AST adapter로 연결한다. 인덱스 삭제 후 Git artifact만으로 완전 재생성 가능해야 한다.

### Context Assembler

현재 task node에서 k-hop neighborhood를 가져오되 budget과 relevance로 절단한다. 항상 포함할 것은 active decisions, constraints, acceptance oracle, unresolved risks다. 과거 대화 전문과 전체 spec tree는 기본 제외한다.

### Workflow Engine

phase와 transition을 YAML state machine으로 관리한다. skill은 “어떻게 수행할지”, engine은 “현재 어디이며 다음에 무엇이 허용되는지”를 맡는다.

### QA Engine

journey graph에서 test charter를 만들고 UI/API/data observer를 실행한다. deterministic assertion과 LLM rubric을 분리하며 evidence bundle에 모든 판정 근거를 연결한다.

## 2. 저장소 제안

```text
.tene/
├── project.yaml
├── specs/
│   ├── goals/
│   ├── requirements/
│   ├── journeys/
│   ├── rules/
│   └── decisions/
├── workflows/
├── qa/
│   ├── charters/
│   ├── oracles/
│   └── rubrics/
├── evidence/              # 큰 binary는 외부 store pointer 가능
├── indexes/               # gitignore 가능한 파생 데이터
└── schemas/
```

## 3. Skill 분해

- `tene-discover`: 사용자 인터뷰와 의도 후보 추출
- `tene-spec`: 후보 확인, normalize, version/diff
- `tene-plan`: spec graph에서 implementation plan 생성
- `tene-workflow`: state transition과 checkpoint
- `tene-qa-plan`: journey/invariant 기반 charter 생성
- `tene-qa-run`: UI/API/data evidence 수집
- `tene-qa-review`: 독립 evaluator와 defect taxonomy
- `tene-drift`: spec-code-test 영향 분석

하나의 거대한 skill은 routing이 불안정하고 context가 커지므로 피한다. 공통 schema/rubric은 references에 두고 필요한 skill만 읽게 한다.

## 4. Hook과 MCP의 경계

### Hook에 적합

- 변경 파일과 연결된 spec ID 누락 검사
- phase transition 전 required artifact validation
- destructive command 전 policy gate
- 완료 시 evidence manifest 존재 확인

### Hook에 부적합

- 핵심 workflow state의 유일한 저장
- 긴 LLM reasoning
- remote business data의 source of truth
- ChatGPT에서도 반드시 필요한 core behavior

### MCP에 적합

- Linear/GitHub/Figma/analytics/log/trace/test DB 조회
- central team spec registry
- remote evidence upload와 query
- 조직 권한·OAuth가 필요한 action

MVP는 외부 서비스 없이 repository-local로 완결해야 한다. 향후 remote MCP를 추가해도 local mode가 깨지지 않아야 한다.

## 5. 단계별 로드맵

### Phase 0 — Portable core

- schema, ID/version/source 규칙
- Markdown templates
- validator와 index builder
- Codex skills-only plugin

성공 기준: 새 thread에서 대화 → confirmed spec → task → acceptance charter가 재현된다.

### Phase 1 — Intent memory와 drift

- conversation source pointer
- contradiction/open-question detection
- spec-code-test traceability
- impacted-spec report

성공 기준: 임의 code change에서 관련 requirement와 stale QA를 설명한다.

### Phase 2 — Full-stack QA

- browser/computer adapter
- API/log/DB observers
- evidence manifest
- independent evaluator와 defect taxonomy

성공 기준: 대표 journey의 UI 전이와 backend invariant를 한 report에서 판정한다.

### Phase 3 — Dynamic orchestration

- machine-readable workflow state
- checkpoints/retries/budgets
- subagent planner/executor/evaluator
- CI `codex exec`와 scheduled regression

성공 기준: session이 바뀌거나 compact돼도 artifact에서 안전하게 재개한다.

### Phase 4 — Marketplace/Team

- `.codex-plugin/plugin.json`
- local marketplace test matrix
- optional remote MCP/OAuth
- telemetry/privacy/admin policy
- universal directory submission

## 6. 핵심 ADR 제안

1. Git artifacts are canonical; DB/indexes are derived.
2. Core workflow must work without hooks.
3. Human confirmation is required to promote inferred intent to confirmed requirement.
4. LLM evaluator cannot override deterministic invariant failures.
5. Implementation and evaluation roles are isolated.
6. Every requirement and verdict must be source/evidence addressable.
7. Context is assembled per task, never dumped wholesale.
8. Provider adapters stay thin; skill content remains provider-neutral where possible.

## 7. 제품 차별화 가설

기존 SDD 도구가 spec 생성·task 분해에, AI QA 도구가 browser test 생성·유지에 집중한다면 tene는 두 시장 사이의 traceability에 집중한다.

> “왜 이 기능이 존재하는가”를 대화에서 보존하고, “어떤 UX·business 흐름이 맞는가”로 변환한 뒤, “실제로 그렇게 동작했다”는 UI/API/data 증거까지 같은 graph에 연결한다.

이 가설은 다음 PoC로 검증하는 것이 좋다.

- checkout 같은 상태·실패·데이터 side effect가 많은 한 feature 선택
- 사용자와 20분 대화 후 spec/journey/oracle 자동 생성
- 의도적으로 UI bug, API bug, DB invariant bug를 각각 주입
- unit/E2E-only baseline과 tene의 탐지율·오분류·수정 시간을 비교
- spec 변경 후 stale test 탐지와 영향 분석 정확도 측정

