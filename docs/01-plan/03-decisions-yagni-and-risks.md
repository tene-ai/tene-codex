# Decisions, YAGNI and Risks

## 1. 확정할 초기 결정

| ID | 결정 | 근거 | 재검토 시점 |
|---|---|---|---|
| DEC-01 | core는 Go `tene-workflow` 별도 binary | portability, 기존 tene 생태계, secret 책임 분리 | 안정화 후 통합 요구가 클 때 |
| DEC-02 | event journal + projection | 복구/감사와 빠른 resume 동시 달성 | event 규모가 병목일 때 |
| DEC-03 | JSON Schema 2020-12 | 언어 중립·tooling·YAML frontmatter 검증 | 복잡한 semantic rule 한계 시 |
| DEC-04 | Markdown은 사람, JSON/NDJSON은 기계 canonical | 편집성과 결정성 분리 | 없음 |
| DEC-05 | default profile `standard` | 보편 사용성과 규율 균형 | beta telemetry 후 |
| DEC-06 | skill + local CLI가 MVP | Codex에서 이식 가능하고 검증 가능 | durable remote orchestration 필요 시 |
| DEC-07 | graph는 embedded adjacency/index | 외부 DB 없이 설치 | 대규모 성능 한계 시 |
| DEC-08 | hooks는 보조 안전망 | hook availability에 core correctness를 의존하지 않음 | 공식 surface 안정화 시 |
| DEC-09 | QA adapter는 capability negotiation | 프로젝트마다 도구가 다름 | adapter ecosystem 성숙 후 |
| DEC-10 | secret 값은 domain에서 표현 불가 | 구조적 유출 방지 | 없음 |

## 2. 대안 검토

### 기존 tene binary에 workflow를 즉시 합치기

단일 설치 장점은 있지만 secret manager의 attack surface와 release coupling이 커진다. MVP에서는 분리하고 `tene-workflow doctor`가 `tene` 호환 버전을 검사한다.

### SQLite를 canonical store로 사용

transaction과 query는 좋지만 diff/review/portable copy가 약하다. MVP는 Git 친화적 파일 journal을 쓰고, SQLite는 대형 graph의 derived cache 후보로 남긴다.

### 모든 작업에 strict full Sprint 강제

품질은 높지만 사소한 작업의 adoption을 해친다. 조직 정책은 strict를 선택할 수 있고, standard/light에서도 intent와 evidence의 핵심은 유지한다.

### 원격 orchestration/MCP/App Server부터 구현

durability와 UI 확장에는 좋지만 인증·운영·호환 부담이 크다. local core contract를 먼저 안정화하며 App Server는 향후 UI/remote controller의 transport로만 추가한다.

## 3. MVP에서 의도적으로 하지 않는 것

- 자체 IDE, dashboard, cloud sync, multi-user approval server
- 범용 AST parser를 모든 언어에 직접 구현
- 자연어만으로 무승인 production/deployment
- secret 조회/복호화/마스킹 값을 모델 context에 전달
- 점수 하나로 모든 품질을 대표
- 문서의 자유 서술을 자동으로 덮어쓰기

이 항목은 핵심인 spec discipline, durable intent, full-flow QA를 검증한 뒤 채택한다.

## 4. 주요 위험과 완화

| 위험 | 신호 | 완화 | owner |
|---|---|---|---|
| 문서 의식화 | 문서는 pass지만 코드와 불일치 | IDs, graph, diff-derived evidence, independent evaluator | Loop/QA |
| context 비대화 | budget 초과·stale 판단 | phase pack, summary provenance, compact | Context |
| implicit trigger 오작동 | 불필요한 Sprint 생성 | confidence threshold, propose-only, negative eval | Skills |
| graph 부정확 | call/data edge 누락 | provider confidence, unknown 표기, runtime observer | Graph |
| QA flaky | 같은 run 결과 변동 | deterministic adapter 우선, retry 분리, quarantine | QA |
| secret 유출 | canary artifact 검출 | no-value model, tene run, redaction, immediate fail | Security |
| 동시 세션 충돌 | stale overwrite | lock+revision+event journal | State |
| Codex surface 변화 | plugin/hook break | official-doc compatibility matrix, core 독립성 | Plugin |
| 과도한 강제 | 사용자 bypass 증가 | profiles, 명시 waiver, 소형 작업 축약 | Product |

## 5. 사람에게 남겨둘 정책

다음 값은 project config default를 제공하되 사용자가 확정한다.

- 어떤 경로/작업 유형이 `strict` 대상인지
- blocker severity 기준과 최종 승인자
- evidence 보존 위치/기간
- 브라우저 자동화 계정과 test-data reset 정책
- production-like 환경 실행 허용 범위
- non-blocking debt의 최대 이월 Sprint 수

이 결정이 없으면 안전한 기본값(로컬/비파괴/strict blocker/secret 미사용)을 적용하고 report의 “정책 결정 필요”에 기록한다.

