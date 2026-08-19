# OpenAI Codex의 바이브 코딩·에이전틱 코딩 방향과 아키텍처

## 1. 방향: 코드 생성기에서 소프트웨어 작업 수행자로

OpenAI의 공개 포지셔닝은 Codex를 autocomplete나 일회성 snippet 생성기가 아니라 end-to-end engineering agent로 둔다. 기능 구현, 복잡한 refactor, migration, test/review, cloud·local 실행을 하나의 작업 단위로 맡기는 방향이다.

공식 best practice가 권하는 입력도 단순 prompt가 아니라 Goal, Context, Constraints, Done when의 네 요소다. 복잡한 작업은 Plan mode 또는 사용자 인터뷰로 모호성을 줄이고, 반복 규칙은 AGENTS.md, 반복 절차는 skill, 외부 시스템은 MCP, 안정된 반복 실행은 scheduled task로 승격한다.

근거: [Codex](https://openai.com/codex/), [Codex best practices](https://learn.chatgpt.com/guides/best-practices), [Codex prompting](https://learn.chatgpt.com/docs/prompting)

## 2. “Humans steer, agents execute”

OpenAI의 harness engineering 사례는 사람이 코드를 직접 쓰는 비중보다 다음을 설계하는 역할이 커진다고 설명한다.

- 의도와 제약을 명확히 하기
- repo를 agent-legible하게 만들기
- 빠른 feedback loop 제공
- 테스트·관찰성·문서를 agent가 수정 가능한 형태로 유지
- 실패를 retrospection하여 harness에 규칙과 도구로 환류

이는 자유로운 vibe coding을 부정하기보다, 탐색 단계의 속도를 유지하면서 production 단계에는 검증 가능한 harness를 추가하는 방향이다. [Harness engineering: leveraging Codex in an agent-first world](https://openai.com/index/harness-engineering/)

## 3. 계층형 아키텍처

```text
User intent / product conversation
              ↓
Plan + durable instructions (AGENTS.md, specs, skills)
              ↓
Agent runtime (Codex CLI/App/IDE/Cloud)
              ↓
Orchestration (subagents, worktrees, SDK/App Server, exec)
              ↓
Tools (shell, patch, browser/computer use, MCP/connectors)
              ↓
Environment (sandbox, credentials, services, repo)
              ↓
Feedback (tests, review, screenshots, logs, telemetry, evals)
              ↺
```

### 주요 primitive

- AGENTS.md: 저장소 단위의 지속 규칙과 완료 정의
- Skills: task-specific reusable workflow와 progressive disclosure
- Plugins: skills와 MCP/hook 등을 배포하는 설치 단위
- MCP: live external context/action
- Subagents: 병렬 분해와 context isolation
- Worktrees: 병렬 변경 격리
- Hooks: tool/command/file-edit lifecycle의 기계적 enforcement
- `codex exec`: CI·script용 non-interactive execution
- Codex SDK/App Server: Codex를 상위 제품·workflow에 embedding
- Browser/Computer Use: UI의 실제 상태를 관찰하고 조작
- Evals/graders: 결과를 score하고 반복 개선

근거: [Codex build plugins](https://developers.openai.com/codex/build-plugins), [AGENTS.md](https://learn.chatgpt.com/docs/agent-configuration/agents-md), [Codex SDK](https://developers.openai.com/codex/sdk), [App Server](https://developers.openai.com/codex/app-server), [Non-interactive mode](https://developers.openai.com/codex/non-interactive)

## 4. 모델과 harness의 공동 최적화

Codex 계열 모델은 실제 software engineering과 장기 tool-use 작업을 겨냥하지만, 모델 성능만으로 제품 품질이 보장되지는 않는다. OpenAI 공식 자료의 반복되는 메시지는 환경·도구·feedback loop가 성능을 크게 좌우한다는 것이다. 따라서 tene는 특정 모델 prompt trick보다 다음 불변 계층에 투자해야 한다.

- 명확한 artifact contracts
- 재현 가능한 environment setup
- 결정론적 validation
- observable evidence
- model-independent state machine
- evaluator와 implementer의 역할 분리

## 5. 바이브 코딩의 구조적 약점

1. 대화 속 의도가 코드에 반영된 뒤 사라진다.
2. “보기에 작동함”과 business invariant 충족이 혼동된다.
3. agent가 만든 test가 같은 오해를 재확인할 수 있다.
4. 변경이 누적되면 spec-code drift가 생긴다.
5. context compaction과 session 교체로 결정 근거가 손실된다.
6. UI, API, DB를 나눠 검사하면 전체 journey 실패를 놓친다.

tene의 기회는 vibe를 무거운 waterfall로 바꾸는 것이 아니다. 짧은 대화에서 의도를 구조화하고, 필요한 부분만 현재 context에 투입하며, acceptance evidence를 자동 축적하는 “lightweight control plane”이 되는 것이다.

