---
schema_version: 1.0.0
document_type: analysis
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380ydgqkg6h9mm5ctvm9b4
phase: loop-check
status: draft
revision: 1088
intent_ids: []
generated_at: 2026-08-20T09:06:41Z
generated_by: tene-workflow
---

# loop-check — Plugin Identity and Installed Runtime Verification

<!-- tene:section:purpose -->
## Purpose

확정된 PRD, Plan, Design과 실제 변경 파일·테스트·설치 관찰을 대조해 구현 편차를 찾고, blocker가 0이 될 때까지 Do와 Loop Check를 반복한다.

<!-- tene:section:scope -->
## Scope

plugin namespace와 9개 skill identity, router, packaging·설치본, confirmed-only criterion, Stop hook JSON, 사용자 대화 언어 문서 계약을 검증한다. 외부 Codex selector UI는 실제 새 session 관찰 전까지 unknown으로 유지한다.

<!-- tene:section:layers -->
## Layers

- Interface: `$tene:<phase>` metadata, Stop hook 응답, 한국어 authored 문서.
- Business Logic: router와 confirmed-only criterion 선택, loop analyzer.
- Persistence: source/stage/install inventory와 workflow state 보존.
- Infrastructure: launcher, checksum, local marketplace 설치, Codex host 경계.

<!-- tene:section:six-questions -->
## Six questions

정의 위치·참조·호출·입력·출력·mutation은 Design의 Six Questions 표를 기준으로 하며, 변경 symbol은 테스트 및 artifact link와 함께 재검증한다.

<!-- tene:section:traceability -->
## Traceability

4개 blocking AC와 6개 task의 정확한 ID를 PRD·Plan·Design·변경 artifact·테스트 결과에 연결한다. Deprecated intent는 audit history로만 보존하고 active blocker나 QA case로 취급하지 않는다.

<!-- tene:section:decisions -->
## Decisions

- 자동 analyzer 결과는 후보 gap이며 구현·문서·도구 오탐으로 분류한 뒤 수정한다.
- 삭제된 rename 원본을 무시하는 방식으로 추적성을 약화하지 않고, status-aware changed artifact 모델로 처리한다.
- 민감정보로 오인된 raw tool output은 evidence에 복사하지 않는다.

<!-- tene:section:freeform -->
## Freeform

이 문서는 사용자의 현재 대화 언어인 한국어로 유지하며 machine marker, ID, path와 code symbol은 원문을 보존한다.

<!-- tene:section:baseline -->
## Baseline

첫 자동 실행은 32개 blocker를 보고했다. 주된 원인은 deprecated intent 오분류, exact ID 문서 누락, directory rename의 삭제 원본을 일반 변경 파일로 취급한 것이었다.

<!-- tene:section:changed-artifacts -->
## Changed artifacts

각 실제 변경 artifact는 T1~T6에 연결한다. rename은 새 경로와 삭제 tombstone을 한 변경 단위로 추적하며, 존재하지 않는 구 경로를 신규 artifact처럼 요구하지 않는다.

<!-- tene:section:gap-matrix -->
## Gap matrix

| 분류 | 판정 | 조치 |
|---|---|---|
| deprecated intent blocker | analyzer 결함 | confirmed intent만 분석하도록 수정 |
| PRD/Plan/Design exact ID 누락 | 문서 추적성 결함 | active 문서에 ID 보강 |
| rename 삭제 원본 unlinked | analyzer 모델 결함 | porcelain status를 보존해 삭제를 별도 처리 |
| 외부 selector UI | 아직 unknown | QA에서 실제 새 session 관찰 필요 |

<!-- tene:section:iterations -->
## Iterations

- Iteration 0: 32 blockers 발견, Do repair로 회귀한다.
- Iteration 1: deprecated intent filtering, exact-ID traceability, deletion tombstone 판정을 수정했다. 자동 재분석 결과 신규 gap 0, 기존 gap 32개 해소, blocker 0으로 수렴했다.

<!-- tene:section:regression -->
## Regression

Go 전체 테스트, Python hook·identity·language tests, routing eval, release smoke, installed-cache doctor를 재실행한다. Loop pass 후 독립 QA evidence gate로 이동한다.

- Go 전체 package: 통과.
- Python plugin/hook/language suite: 14개 통과.
- Routing corpus: 전체 기준 통과; 9개 skill explicit invocation은 각 5/5, wrong-phase와 unnecessary trigger는 0.
- Release smoke: staged package, identity와 launcher 검증 통과.
