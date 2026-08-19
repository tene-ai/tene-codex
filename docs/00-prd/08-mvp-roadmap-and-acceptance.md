# MVP 범위·로드맵·Acceptance Criteria

## 1. 개발 전략

한 번에 App Server, remote MCP, 모든 graph provider를 만들지 않는다. 먼저 어떤 repository에서도 작동하는 local skills-only plugin + deterministic core를 검증한다.

## 2. Sprint Master Plan

### Sprint 1 — Document and State Foundation

목표:

- `.tene-workflow` schema
- docs/sprints scaffold
- sprint state machine
- document validator
- status/resume/compact/archive

Acceptance:

- 새 project에서 1개 command로 구조 생성
- invalid phase skip 차단
- process 종료 후 active phase/task 복원
- terminal sprint compact 후 active state 크기 감소
- `.tene/` read/write 0건

### Sprint 2 — Intent and Spec Graph

목표:

- intent candidate extraction/confirmation
- stable IDs와 source refs
- PRD/plan/design templates
- explicit traceability graph

Acceptance:

- 대화에서 proposed intent 생성
- 사용자 확인 없이 confirmed 승격 불가
- requirement→AC→task→design link validation
- source 없는 requirement 탐지

### Sprint 3 — Code Understanding and Loop Check

목표:

- Understanding Layer classifier
- Six Questions extractor
- CodeGraph/LSP/AST/rg provider interface
- impact/drift/gap engine

Acceptance:

- 변경 symbol의 definition/reference/call/input/output mapping 생성
- dynamic/unknown relation을 명시
- seeded PRD/design/code drift 90% 이상 탐지
- blocking gap 0까지 loop iteration 기록

### Sprint 4 — Intent-driven QA

목표:

- QA charter와 evidence schema
- unit/integration/E2E runner adapters
- browser/API/data observers
- independent evaluator와 gate

Acceptance:

- 한 reference app에서 UI→API→DB journey 증거 연결
- product/spec/test/environment failure 분류
- blocking AC evidence coverage 100% 계산
- secret leakage seed를 모두 차단

### Sprint 5 — Codex Plugin UX

목표:

- 9개 skill
- `agents/openai.yaml`
- implicit/explicit routing eval
- optional hooks
- `.codex-plugin/plugin.json`
- local marketplace test

Acceptance:

- positive/negative trigger precision·recall 기준 통과
- 새 Codex thread에서 sprint resume
- hook이 없어도 core workflow 완료
- plugin validator 통과

### Sprint 6 — Advanced Orchestration

목표:

- subagent role profiles
- worktree orchestration
- App Server adapter
- scheduled/CI workflows
- optional remote MCP

Acceptance:

- thread/turn/item을 workflow run에 mapping
- interrupt/resume/fork recovery
- parallel sprint isolation
- approval request 처리

## 3. MVP 필수 기능

- Sprint master plan과 fixed lifecycle
- Canonical durable state
- PRD/plan/design/analysis/QA/report templates
- Understanding Layer + Six Questions
- intent confirmation과 traceability
- loop-check gap engine
- QA gate와 evidence manifest
- explicit/implicit skill invocation
- tene CLI secret-safe adapter
- archive/clear/compact

## 4. MVP 이후

- Graph visualization UI
- Central team spec registry
- Linear/GitHub/Figma/log/DB MCP
- App Server desktop dashboard
- cross-project learning과 reusable policy packs
- marketplace public submission

## 5. 최상위 Acceptance Criteria

### AC-PRODUCT-01 Mandatory sprint discipline

Implementation 요청이 active mandatory project에서 들어오면 현재 sprint/phase를 확인하고 required predecessor artifact가 없으면 생성·확정 흐름으로 안내한다.

### AC-PRODUCT-02 Durable intent

새 session에서 원본 대화 없이 confirmed intent, open policy, current gap, next action을 복원한다.

### AC-PRODUCT-03 Whole-system understanding

변경 기능을 네 Understanding Layer와 Six Questions로 설명하고 unknown을 숨기지 않는다.

### AC-PRODUCT-04 Evidence-based 100%

모든 blocking AC가 evidence-backed pass일 때만 100%/complete를 선언한다.

### AC-PRODUCT-05 Full-flow QA

대표 reference app에서 UI transition, API response, persistence effect를 하나의 journey로 연결한다.

### AC-PRODUCT-06 Secret safety

Plugin의 모든 공식 workflow와 red-team test에서 secret plaintext가 prompt/tool output/document/evidence에 나타나지 않는다.

### AC-PRODUCT-07 Portable use

Go CLI, Node web app, Python API 등 최소 3개 stack에서 동일한 product workflow와 document contract가 작동한다.

### AC-PRODUCT-08 Report continuity

Report가 이전 sprint 연결, changed files/functions, intent 충족, layer, Six Questions, policy/deferred work를 빠짐없이 포함한다.

## 6. 검증용 Reference Projects

1. Go CLI: command→service→SQLite/file
2. Next.js full-stack: UI→route handler→DB
3. Python API + worker: endpoint→service→queue→worker→DB

각 project에 의도적으로 architecture drift, missing error UX, wrong data mutation, secret leakage를 주입해 baseline Codex와 tene-assisted Codex를 비교한다.

## 7. 주요 위험

- Process fatigue: 작은 변경에 과도한 문서 비용
  - 대응: sprint size profile(light/standard/strict), 필수 invariant는 유지
- LLM document inflation
  - 대응: schema와 evidence 중심, 중복 content lint
- False graph confidence
  - 대응: explicit/inferred/runtime edge 분리
- Self-validating agent
  - 대응: deterministic gate + independent evaluator
- Host API drift
  - 대응: Codex version contract test, core portability
- Secret leakage through test logs
  - 대응: tene runtime + pre-storage redaction + leakage gate

## 8. Definition of Done

MVP는 문서가 작성되었을 때가 아니라 다음이 모두 증명됐을 때 완료다.

- Plugin/core/tests가 실제 reference projects에서 작동
- phase transition과 resume가 crash/compaction scenario를 견딤
- intent→code→QA evidence graph가 query 가능
- seeded defect와 drift 기준 통과
- secret red-team 0 leak
- report가 source data에서 재생성 가능
- local marketplace install과 새 thread invocation 검증

