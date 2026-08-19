# Requirements Traceability Plan

## 1. 기능 요구사항

| PRD | 구현 work package | 핵심 설계 | 검증 |
|---|---|---|---|
| FR-01 Sprint Master Plan | WP-02,03,04 | state/workflow/document | master status, dependency tests |
| FR-02 Fixed Lifecycle | WP-03 | transition table/guards | invalid-skip matrix |
| FR-03 Intent Memory | WP-05,06 | intent ledger/spec graph | resume/drift fixtures |
| FR-04 Standard Documents | WP-04 | templates/markers/schema | golden validation |
| FR-05 Loop Check | WP-08,09 | code graph/gap evaluator | spec mutation tests |
| FR-06 Intent-driven QA | WP-10 | charter/observer/adapters | full journey projects |
| FR-07 Evidence Gate | WP-09,10 | manifest/evaluator/blocker | false-pass suite |
| FR-08 Durable State | WP-02,07 | journal/projection/context | restart/replay |
| FR-09 Compact/Clear | WP-02,07 | retention classes | semantic equivalence |
| FR-10 Explicit/Implicit | WP-12 | skill router/evals | routing corpus |
| FR-11 Secret-safe | WP-11 | tene adapter/redactor | canary corpus |

## 2. 제품 acceptance criteria

| AC | 완료 evidence |
|---|---|
| AC-PRODUCT-01 | phase guard system test와 archived Sprint |
| AC-PRODUCT-02 | 새 Codex session resume transcript, revision provenance |
| AC-PRODUCT-03 | 4 Layers/6 Questions가 채워진 polyglot report |
| AC-PRODUCT-04 | blocker 누락 mutation이 gate에서 거절된 기록 |
| AC-PRODUCT-05 | UI→API→data full-flow evidence chain |
| AC-PRODUCT-06 | secret canary 0건과 forbidden-command tests |
| AC-PRODUCT-07 | 3 reference project clean install 결과 |
| AC-PRODUCT-08 | 연속 두 Sprint report의 predecessor/intent 연결 |

## 3. 문서 요구사항 추적

| 요구 | design owner | validator |
|---|---|---|
| 공통 frontmatter | `03-document-template-contracts.md` | schema + marker validator |
| 4 Understanding Layers | `06-graph-and-context-engine.md` | layer coverage rule |
| 6 Questions | `06-graph-and-context-engine.md` | symbol answer completeness |
| QA gate result | `07-qa-evidence-and-gates.md` | evidence evaluator |
| 이월/정책 결정 | report template | deferred/waiver schema |
| 자유 관점 허용 | document merge contract | generated-region test |

## 4. 추적성 운영 규칙

- 모든 task는 최소 하나의 `intent_id` 또는 `ac_id`를 가져야 한다.
- 모든 blocking AC는 하나 이상의 verification과 evidence를 가져야 한다.
- code component edge는 자동 추론일 수 있으나 confidence와 source를 기록한다.
- report 생성 시 orphan intent/AC/task/evidence를 blocker로 검사한다.
- 요구가 바뀌면 기존 ID를 재사용해 의미를 바꾸지 않고 새 revision 또는 superseding ID를 만든다.

