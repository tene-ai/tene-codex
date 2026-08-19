# tene Codex Plugin — Product Requirements & Architecture

작성 기준일: 2026-08-20

이 문서 집합은 한 사용자의 작업 습관을 모든 사용자·모든 프로젝트에서 재사용 가능한 Codex plugin으로 일반화한 제품 기획 및 기술 설계다. 핵심 목표는 다음 두 가지다.

1. 모든 기능 개발을 `PRD → Plan → Design → Do → Loop Check → QA → Report → Archive` sprint로 수행한다.
2. 대화에서 합의된 기획 의도를 지속 가능한 spec으로 보존하고, 코드·UX·데이터 흐름을 그 의도에 맞춰 검증한다.

## 문서 목록

1. [제품 PRD](./00-product-prd.md)
2. [Sprint 및 Workflow/Task Management 요구사항](./01-sprint-workflow-requirements.md)
3. [문서 정보 아키텍처와 표준 양식](./02-document-system-and-templates.md)
4. [Plugin 기술 아키텍처](./03-plugin-architecture.md)
5. [Harness·Graph·Context Engineering 설계](./04-engineering-architecture.md)
6. [기획 의도 기반 QA 및 Gate 설계](./05-intent-driven-qa-gates.md)
7. [Skills·Commands·자연어 Trigger 명세](./06-skills-commands-triggers.md)
8. [tene CLI 통합과 Secret Security Boundary](./07-tene-cli-security-integration.md)
9. [MVP 범위·로드맵·Acceptance Criteria](./08-mvp-roadmap-and-acceptance.md)

## 핵심 제품 결정

- Plugin은 workflow를 안내하는 prompt 묶음이 아니라 **artifact-backed sprint control plane**이다.
- Codex thread/plan/memory는 실행 보조 수단이고 source of truth는 아니다.
- Canonical state는 Git으로 검토 가능한 Markdown/YAML/JSON이며, graph index는 재생성 가능한 파생물이다.
- 자연어 자동 trigger는 편의 기능이다. phase transition과 gate 판정은 결정론적 validator가 담당한다.
- 모든 프로젝트에 같은 기술 스택을 강요하지 않는다. Understanding Layer와 6 Questions라는 공통 관찰 모델만 강제한다.
- `.tene/`는 encrypted secret vault 전용이다. Plugin은 읽지 않으며 workflow state는 `.tene-workflow/`에 분리한다.
- QA는 test count가 아니라 requirement·journey·business invariant·evidence coverage로 판정한다.
- “100% 달성”은 LLM의 자기평가 점수가 아니라, sprint 시작 시 합의된 acceptance criteria가 모두 evidence로 증명된 상태다.

