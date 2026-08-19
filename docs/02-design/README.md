# tene Codex Plugin — 상세 설계

이 폴더는 [구현 계획](../01-plan/README.md)을 코드로 옮길 때 지켜야 할 데이터·명령·컴포넌트 계약을 정의한다. `MUST`, `MUST NOT`, `SHOULD`는 각각 필수, 금지, 권고를 뜻한다.

## 문서 지도

| 문서 | 설계 범위 |
|---|---|
| [00-system-architecture.md](00-system-architecture.md) | 전체 구조와 실행 흐름 |
| [01-state-and-storage-schema.md](01-state-and-storage-schema.md) | 영속 상태, ID, event, 동시성 |
| [02-workflow-cli-contract.md](02-workflow-cli-contract.md) | 상태 머신과 CLI/API |
| [03-document-template-contracts.md](03-document-template-contracts.md) | 폴더·문서 template·검증 |
| [04-skill-contracts-and-routing.md](04-skill-contracts-and-routing.md) | Codex skills와 자연어 routing |
| [05-hooks-codex-integration.md](05-hooks-codex-integration.md) | Codex surface와 hook 경계 |
| [06-graph-and-context-engine.md](06-graph-and-context-engine.md) | spec/code graph, 4 Layers, 6 Questions, context |
| [07-qa-evidence-and-gates.md](07-qa-evidence-and-gates.md) | 기획 의도 기반 QA와 판정 |
| [08-tene-secret-runtime.md](08-tene-secret-runtime.md) | tene CLI 보안 통합 |
| [09-errors-migrations-concurrency.md](09-errors-migrations-concurrency.md) | 오류·복구·migration |
| [10-plugin-package-and-marketplace.md](10-plugin-package-and-marketplace.md) | plugin package와 배포 |
| [11-testing-evals-and-acceptance.md](11-testing-evals-and-acceptance.md) | 테스트/eval/최종 수용 |

## 설계 우선순위

`사용자 기획 의도 > 안전 invariant > 상태 일관성 > 재현 가능한 evidence > 자동화 편의성` 순서로 충돌을 해결한다.

