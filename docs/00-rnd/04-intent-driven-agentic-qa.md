# 기획 의도 기반 Agentic QA 시장·방법론·기술 아키텍처

## 1. 문제 재정의

Unit/E2E script 실행 능력의 부족이 핵심 문제가 아니다. 더 근본적인 문제는 **테스트 oracle**, 즉 무엇이 맞는지를 판단할 구조화된 기준이 대화와 사람의 머릿속에만 있다는 점이다. LLM agent가 browser를 잘 조작해도 올바른 UX, business rule, 데이터 side effect를 모르면 “클릭은 성공했다” 수준에서 멈춘다.

tene가 풀어야 할 문제는 다음 폐루프다.

```text
Conversation → Intent extraction → Versioned specification
      → Journey/state model → Executable oracle
      → UI/API/DB observation → Evidence + verdict
      → Human clarification / spec update / defect
```

## 2. 경쟁·인접 제품군

### Spec-driven development

- GitHub Spec Kit: Spec → Plan → Tasks → Implement, Markdown artifact와 cross-artifact analysis. Codex 포함 다수 agent 지원. [공식 문서](https://github.github.io/spec-kit/index.html), [GitHub 소개](https://github.blog/ai-and-ml/generative-ai/spec-driven-development-with-ai-get-started-with-a-new-open-source-toolkit/)
- Kiro: requirements/design/tasks와 EARS-style requirement 중심의 IDE workflow.
- OpenSpec: lightweight, tool-agnostic delta spec에 강점.
- BMAD: 역할 기반 multi-agent와 phase gate를 갖춘 무거운 방법론.
- Tessl 등: spec을 유지 가능한 engineering artifact로 다루는 계열.

이들은 “무엇을 만들지” 구조화하지만 실제 runtime evidence와 spec의 지속적 연결은 제품별 편차가 크다. tene는 이 빈 공간을 QA까지 연결해야 한다.

### AI browser/E2E automation

Momentic, mabl, Testim, QA Wolf, Autify, Functionize, KaneAI 등은 natural-language test authoring, self-healing locator, browser execution, visual regression, CI integration을 제공한다. 이 제품들은 실행 자동화와 유지보수 비용 절감에 강점이 있지만, repository-local 제품 의도 graph와 대화의 결정 이력을 source of truth로 삼는지는 별도 확인이 필요하다. 제품 비교 시 마케팅 주장보다 다음을 PoC로 검증해야 한다.

- business invariant를 UI 밖 API/DB까지 검사하는가
- 자동 수정이 false positive를 숨기지 않는가
- test가 어느 requirement version에서 파생됐는가
- ambiguity를 사람에게 되묻는가, 임의 가정하는가
- evidence가 audit/replay 가능한가

### Agent 평가·UX 연구

[UXAgent](https://arxiv.org/abs/2502.12561)는 LLM-simulated user와 browser connector로 usability study를 확장하는 연구다. Model-based testing은 상태 머신에서 path를 생성하고 coverage를 측정하는 오래된 기반 방법론이다. LLM은 persona·heuristic exploration에, 명시적 state graph는 coverage와 재현성에 적합하므로 둘을 결합하는 편이 좋다.

## 3. 기능 의도 추출 방법

대화를 그대로 저장하는 것만으로는 부족하다. 매 대화 turn에서 아래 타입의 candidate를 추출하고 사용자 확인을 거쳐 승격한다.

| 타입 | 예시 |
|---|---|
| Goal | 사용자는 2분 안에 가입을 완료한다 |
| Actor/Persona | 신규 개인 사용자, 조직 관리자 |
| Requirement | 이메일 인증 전 결제 불가 |
| Business rule | 쿠폰은 주문당 하나만 적용 |
| UX state | loading, empty, error, retry, success |
| Data invariant | 실패 결제는 주문을 paid로 바꾸지 않음 |
| Non-functional | p95 2초, WCAG 기준 |
| Assumption | 소셜 로그인은 이번 범위 제외 |
| Open question | 탈퇴 후 데이터 보존 기간 |
| Acceptance example | Given/When/Then concrete case |

각 항목은 stable ID, status(`proposed/confirmed/deprecated`), version, source, owner, confidence, effective range를 가진다. LLM confidence가 높아도 사용자 합의를 대체하지 않는다.

## 4. Spec 및 memory 모델

```yaml
id: REQ-CHECKOUT-012
type: business_rule
title: 결제 실패 시 주문 상태 보존
statement: 결제가 실패하면 주문은 paid가 될 수 없다.
status: confirmed
version: 3
sources:
  - kind: conversation
    ref: session-2026-08-20#turn-18
actors: [customer]
preconditions: [cart.not_empty]
invariants: [order.status != paid]
links:
  journeys: [JRN-CHECKOUT-001]
  code: [src/payments/confirm.ts]
  tests: [QA-CHECKOUT-FAIL-003]
supersedes: REQ-CHECKOUT-012@2
```

Canonical artifact는 Git의 Markdown/YAML이며, 검색과 영향 분석용 SQLite/graph index는 재생성 가능해야 한다. 원본 대화는 privacy 정책과 보존 기간을 적용하고, spec에는 필요한 근거 pointer와 합의 내용만 남긴다.

## 5. QA 실행 아키텍처

### Planner

현재 change와 관련된 requirement/journey neighborhood를 검색한다. happy path만 아니라 alternate/error/recovery path와 cross-role sequence를 만든다. pairwise/constraint 기반으로 case explosion을 제어한다.

### Executor

Codex의 browser/computer use 또는 Playwright 같은 결정론적 driver로 실제 UI를 조작한다. API client, log query, DB read-only probe를 함께 사용한다. write action은 격리된 test environment와 seed data에서만 수행한다.

### Observer

다중 채널 evidence를 수집한다.

- UI: screenshot, accessibility tree, visible text, route
- Network: request/response, status, correlation ID
- Backend: structured logs, trace, emitted event
- Data: before/after snapshot, invariant query
- Performance: duration, retries, console errors

### Evaluator

결정론적 assertion을 먼저 실행하고, UX rubric·문구 품질·인지 부하처럼 주관적 항목만 독립 LLM evaluator에 맡긴다. 구현 agent와 평가 agent를 분리하고 evaluator에는 expected outcome과 evidence를 주되 구현 reasoning은 최소화해 confirmation bias를 줄인다.

### Reporter/Updater

결과를 pass/fail만으로 줄이지 않는다.

- product defect
- spec ambiguity
- test/harness defect
- environment failure
- observation insufficient

으로 분류한다. spec 변경은 자동 overwrite하지 않고 proposed diff와 영향 범위를 사용자에게 제시한다.

## 6. Coverage 모델

전통 code coverage 외에 다음을 추적한다.

- Requirement coverage: confirmed requirement 중 oracle이 있는 비율
- Journey edge coverage: 상태 전이 edge 실행 비율
- Business-rule coverage: invariant가 적어도 한 번 검증된 비율
- Negative/recovery coverage
- Role/permission matrix coverage
- Evidence completeness: UI/API/data 중 필요한 채널 충족률
- Drift coverage: 변경 code와 연결된 spec/test의 최신성

## 7. Codex 기능 활용

- Plan mode/인터뷰: 모호한 의도 발굴
- AGENTS.md: 저장소 전체의 완료 정의와 QA 원칙
- Skills: intent interview, spec normalization, test planning, evidence review
- MCP: Linear/GitHub/Figma/analytics/log/DB 같은 외부 근거
- Subagents: planner, explorer, executor, evaluator 분리
- Worktrees: 구현과 QA 생성의 충돌 격리
- Browser/Computer Use: 실제 사용자 흐름
- `codex exec`/SDK/App Server: CI와 상위 orchestrator 연동
- Hooks: spec drift 또는 evidence 누락의 mechanical gate
- Scheduled tasks: nightly journey regression과 stale spec review

공식 OpenAI use case도 Computer Use로 실제 제품 흐름을 클릭하고 문제를 기록하는 QA를 제시한다. [Codex use cases](https://learn.chatgpt.com/use-cases), [Computer Use](https://learn.chatgpt.com/docs/computer-use)

## 8. 위험과 방어

- Hallucinated intent: source pointer 없는 requirement 생성 금지
- Self-confirming tests: independent evaluator와 mutation/negative case 사용
- Brittle browser actions: semantic locator 우선, healing은 기록·승인
- Privacy: conversation/raw evidence 최소 보존, secret redaction
- Non-determinism: seed, environment fingerprint, replay bundle 저장
- Spec bloat: active neighborhood만 context pack으로 구성
- False authority: LLM verdict는 evidence와 rubric 없는 최종 판정이 아님

