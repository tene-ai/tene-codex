# Document and Template Contracts

## 1. Sprint directory

```text
docs/sprints/<sprint-id>-<slug>/
├── 00-prd/00-prd.md
├── 01-plan/00-plan.md
├── 02-design/00-design.md
├── 03-analysis/{00-loop-check.md,gaps.md}
├── 04-qa/{00-qa-plan.md,01-qa-result.md,evidence-manifest.json}
├── 05-report/00-report.md
└── 99-archive/archive-manifest.json
```

문서가 방대하면 같은 phase 폴더에 주제별 파일을 추가한다. `00-*`는 index/summary이며 링크 누락을 validator가 검사한다.

## 2. Frontmatter

```yaml
---
schema_version: 1.0.0
document_type: prd
project_id: project_...
sprint_id: sprint_...
phase: prd
status: draft
revision: 3
intent_ids: [intent_...]
predecessor_ids: []
generated_at: 2026-08-20T03:00:00Z
generated_by: tene-workflow
---
```

`status`, `revision`, ID reference는 machine state와 일치해야 한다. 사람 편집 후 `document sync`가 semantic diff를 보여주고 승인 후 event로 반영한다.

## 3. 모든 문서의 필수 section marker

```markdown
<!-- tene:section:purpose -->
## 목적과 기획 의도
<!-- tene:section:scope -->
## 범위와 비범위
<!-- tene:section:layers -->
## Understanding Layers
<!-- tene:section:six-questions -->
## Six Questions
<!-- tene:section:traceability -->
## 추적성
<!-- tene:section:decisions -->
## 결정·가정·미결정
<!-- tene:section:freeform -->
## 추가 관점
```

marker ID로 검증해 제목 언어/표현은 자유롭게 바꿀 수 있다. `freeform` 아래는 tool이 수정하지 않는다.

## 4. 문서별 추가 계약

### PRD

Problem, actors, current/desired journey, data journey, intent, AC, policy, non-goal, metric, open question.

### Plan

Work package, dependency, order, risk, verification, rollout/rollback, policy decision, YAGNI.

### Design

Component boundary, public interface, data schema, state transition, failure/recovery, security, observability, migration, test seam.

### Analysis / Loop Check

Baseline revision, changed files/symbols, PRD/Plan/Design gap matrix, 4 Layers, 6 Questions, regression/debt, iteration history. 각 gap은 `open/resolved/waived/deferred`와 evidence를 가진다.

### QA

AC coverage, capability snapshot, environment, test charter, UX/data journey, expected observers, run/evidence, evaluator verdict, defect, residual risk.

### Report

이전 Sprint 연결, 파일/기능 변경, 충족 intent, 4 Layers, 6 Questions, QA verdict, 정책 결정/이월 사유, operational note, next Sprint recommendation.

## 5. Understanding Layers table

| Layer | Entry/components | 변경 | upstream/downstream | evidence |
|---|---|---|---|---|
| Interface | | | | |
| Business Logic | | | | |
| Persistence | | | | |
| Infrastructure | | | | |

적용되지 않는 layer는 빈칸이 아니라 `N/A`와 조사 근거를 적는다.

## 6. Six Questions table

| Name | Defined at | Imported/referenced at | Called/used at | Input shape | Output/mutation |
|---|---|---|---|---|---|

파일 경로와 symbol/line 또는 stable locator를 쓴다. unknown은 추측하지 않고 provider/capability 이유를 기록해 gap으로 만든다.

## 7. Generated region

`<!-- tene:generated:<id>:start -->`와 `end` 사이만 자동 갱신한다. 갱신 전 content hash가 예상과 다르면 conflict를 반환한다. authored prose와 추가 section은 항상 보존한다.

