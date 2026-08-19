# 문서 정보 아키텍처와 표준 양식

## 1. 프로젝트 문서 구조

```text
docs/
├── sprint-master/
│   ├── master-plan.md
│   └── decisions.md
└── sprints/
    ├── SPR-001-foundation/
    │   ├── sprint.yaml
    │   ├── 00-prd/
    │   │   └── prd.md
    │   ├── 01-plan/
    │   │   └── plan.md
    │   ├── 02-design/
    │   │   ├── design.md
    │   │   └── adr/
    │   ├── 03-analysis/
    │   │   ├── loop-check.md
    │   │   └── gap-register.md
    │   ├── 04-qa/
    │   │   ├── qa-plan.md
    │   │   ├── qa-result.md
    │   │   └── evidence-manifest.json
    │   ├── 05-report/
    │   │   └── report.md
    │   └── 99-archive/
    └── _archive/
```

사용자가 예시로 제시한 `00-prd, 01-plan, 02-design, 03-analysis, 04-report`를 확장해 QA를 독립 phase로 보존하기 위해 `04-qa`, `05-report`를 둔다. 폴더 번호는 생성 이후 바꾸지 않는다.

## 2. 모든 문서의 공통 frontmatter

```yaml
---
schema: tene.document/v1
documentId: DOC-SPR-001-PRD
sprintId: SPR-001-foundation
type: prd
status: confirmed
version: 3
owners: [product]
createdAt: 2026-08-20T10:00:00+09:00
updatedAt: 2026-08-20T12:00:00+09:00
sourceRefs:
  - conversation:session-id#turn-12
supersedes: DOC-SPR-001-PRD@2
tags: [foundation]
---
```

## 3. 공통 필수 섹션

PRD, plan, design, analysis, QA, report는 문서 목적에 맞게 깊이는 달라도 다음을 포함한다.

1. Purpose and relation to previous sprints
2. Scope / Out of scope
3. Intent and acceptance criteria
4. Understanding Layer impact
5. Six Questions mapping
6. Traceability links
7. Assumptions, risks, constraints
8. Decisions required from user
9. Deferred work and reason
10. Additional perspectives — 자유 섹션

## 4. Understanding Layer 표준

| Layer | 조사 질문 | 일반적인 대상 |
|---|---|---|
| Interface / Entry Point | 사용자는 어디서 진입하고 무엇을 관찰하는가? | UI, CLI, API controller, webhook, scheduler, command |
| Business Logic / Processing Rules | 어떤 규칙·분기·상태 전이가 적용되는가? | service, handler, usecase, reducer, domain policy |
| Persistence / Data | 무엇이 어디에 저장·조회·변경되는가? | DB, file, cache, queue, external API |
| Infrastructure / Runtime | 어떤 실행 환경·권한·배포 조건이 필요한가? | server, container, cloud, auth, CI/CD, observability |

각 layer는 `affected | not-affected | unknown`을 명시한다. `not-affected`도 조사 근거를 적어야 한다.

## 5. Six Questions 표준

모든 주요 processing component에 대해 다음을 답한다.

1. 선언·정의된 이름은 무엇인가?
2. 어떤 파일에 정의되어 있는가?
3. 어디에서 import 또는 참조되는가?
4. 어디에서 호출·사용되는가?
5. 어떤 형태·구조의 데이터를 입력받는가?
6. 어떤 형태·구조의 데이터를 반환하거나 변경하는가?

권장 표:

| Symbol | Defined at | Imported/referenced by | Called by | Input contract | Output/side effect | Layer |
|---|---|---|---|---|---|---|

이 표는 정적 code graph만 복사하는 문서가 아니다. dynamic dispatch, framework registration, CLI routing, event subscriber, database side effect도 포함해야 한다.

## 6. 문서별 추가 필수 항목

### PRD

- Problem, persona, user journey
- Functional/non-functional requirements
- Business rules, UX states, data invariants
- Acceptance criteria와 priority
- Policy decision/open question
- Failure/recovery scenarios

### Plan

- Work breakdown와 dependency graph
- Expected files/symbols/layers
- Test/evidence strategy
- Risks, pre-mortem, rollback
- Scope budget와 stop condition
- Parallelizable/conflicting tasks

### Design

- Current vs target architecture
- Interface/API/data/event contracts
- call/data flow sequence
- Understanding Layer별 상세 설계
- Six Questions 예상 mapping
- Security, migration, backward compatibility
- ADR와 rejected alternatives

### Analysis / Loop Check

- PRD↔plan↔design↔implementation coverage matrix
- changed files/symbols actual graph
- missing/extra/drift findings
- root cause와 severity
- fix iteration history
- remaining uncertainty

### QA

- Intent-based test charter
- Environment/seed/actor setup
- Unit/integration/E2E/journey/data test matrix
- UX transitions와 business invariant
- evidence manifest와 reproducibility
- defect/spec ambiguity/harness/environment classification
- gate verdict

### Report

- 이전 sprint와 기능 연결
- 생성·수정 파일과 구현 기능
- 어떤 기획 의도를 충족했는지
- Understanding Layer별 실제 변경
- Six Questions의 실제 구현 답변
- 정책 결정 필요 및 이월 작업/사유
- AC/gate 결과와 evidence
- 예상과 실제의 차이
- 회고: 잘된 점, 실패, harness 개선
- 다음 sprint recommendation

## 7. Template validation

Validator는 heading 문자열만 검사하지 않는다.

- frontmatter schema와 ID/link integrity
- required section의 non-empty content
- changed file과 Six Questions row 연결
- confirmed AC와 QA evidence 연결
- deferred task의 owner/reason/target 존재
- report의 previous/next sprint link
- secret-like pattern redaction

프로젝트는 `.tene-workflow/project.yaml`에서 추가 section과 custom layer를 선언할 수 있으나 공통 필수 항목을 제거할 수 없다.

