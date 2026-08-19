# tene Codex Plugin — 구현 계획

이 폴더는 `docs/00-prd`의 제품 요구사항을 실제 개발 순서, 산출물, 검증 게이트로 변환한다. 계획의 기준 단위는 **제품 Sprint**이며, 각 Sprint 안에서 PRD → Plan → Design → Do → Loop Check → QA → Report → Archive를 수행한다.

## 문서 지도

| 문서 | 목적 |
|---|---|
| [00-master-implementation-plan.md](00-master-implementation-plan.md) | 구현 전략, 단계, 완료 조건 |
| [01-work-breakdown-and-dependencies.md](01-work-breakdown-and-dependencies.md) | 패키지·작업 단위·의존성 |
| [02-verification-and-release-plan.md](02-verification-and-release-plan.md) | 테스트, 평가, 릴리스 게이트 |
| [03-decisions-yagni-and-risks.md](03-decisions-yagni-and-risks.md) | 확정 설계 선택, 보류 범위, 위험 |
| [04-requirements-traceability.md](04-requirements-traceability.md) | PRD 요구사항에서 구현·검증까지의 추적표 |

## 실행 원칙

1. 상태 머신과 스키마를 먼저 고정하고, 그 위에 문서 생성기와 스킬을 올린다.
2. 사람에게 보이는 Markdown과 기계가 판정하는 구조화 상태를 분리하되 ID로 연결한다.
3. `loop-check`와 `qa`는 생성 에이전트의 자기 주장 대신 재현 가능한 evidence를 판정한다.
4. Secret 값은 플러그인의 데이터 모델에 존재하지 않는다. 이름만 보관하고 실행은 `tene run`에 위임한다.
5. Codex 기능이 없어도 core CLI가 결정론적으로 동작해야 하며, 스킬은 이를 오케스트레이션한다.

