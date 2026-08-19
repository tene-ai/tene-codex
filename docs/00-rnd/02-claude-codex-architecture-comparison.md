# Claude Code와 Codex의 워크플로·태스크·하네스·그래프·컨텍스트 비교

## 용어 주의

Dynamic Workflow, Task Management System, Harness Engineering, Graph Engineering, Context Engineering을 모두 Claude Code의 고유 기능명으로 보면 부정확하다. 공식 기능과 업계 아키텍처 패턴이 섞여 있다. 이 문서는 각 개념을 primitive 수준으로 분해해 비교한다.

## 1. Dynamic Workflow

### Claude Code

Claude Code는 저수준·비강제적 agent harness라는 철학 위에서 다음을 조합한다.

- skills/commands: 상황별 절차와 slash entry point
- hooks: SessionStart, UserPromptSubmit, PreToolUse, PostToolUse, Stop, SubagentStop, PreCompact 등 lifecycle event
- subagents/agent teams: 역할별 위임·병렬 실행·worktree isolation
- MCP/LSP/tools: 외부 action과 code intelligence
- plan mode와 task list: 실행 전 분해 및 진행 추적

Hook은 command, prompt, experimental agent hook 형태로 정책 검사·컨텍스트 주입·완료 gate를 만들 수 있다. 그러나 hook이 모든 내부 상태 전이를 제어하지는 않는다. 예를 들어 [PostCompact hook 요청 이슈](https://github.com/anthropics/claude-code/issues/32026)는 compaction 이후 구조화 상태 재주입의 공백을 보여준다.

근거: [Claude Hooks guide](https://code.claude.com/docs/en/hooks-guide), [Hooks reference](https://code.claude.com/docs/en/hooks), [Subagents](https://code.claude.com/docs/en/sub-agents), [Parallel agents](https://code.claude.com/docs/en/agents)

### Codex 대응 요소

Codex에는 Plan mode, skills, hooks, subagents, AGENTS.md, MCP, Git worktrees, scheduled tasks, non-interactive `codex exec`, SDK, App Server가 있다. 즉 primitive는 충분하지만 “하나의 선언형 workflow engine”이라기보다 조합 가능한 실행 표면이다.

tene는 workflow를 skill 본문의 암묵적 순서로만 두지 말고 machine-readable state로 소유해야 한다.

```text
DISCOVER → SPECIFY → DESIGN → PLAN → IMPLEMENT → VERIFY → ACCEPT
    ↑          │          │                    │          │
    └──────── change request / evidence-driven feedback ──┘
```

각 transition은 precondition, required artifacts, approver, evidence, rollback target을 가져야 한다. Codex/Claude는 이 state machine의 실행자이고 tene가 controller다.

## 2. Task Management System

Claude Code와 Codex 모두 모델이 작업을 분해하고 상태를 갱신할 수 있으나, 대화 내부 task list는 제품 요구사항의 영구 system of record로 쓰기 어렵다. session compaction, branch divergence, 병렬 agent 충돌, host 변경 때문이다.

권장 분리:

- Ephemeral execution tasks: 현재 run의 todo/plan
- Durable feature tasks: Git에 저장되는 `tasks.md` 또는 YAML
- Organizational work items: Linear/GitHub Issues/Jira 등 MCP 연결
- Traceability edges: requirement → design → task → code → test → evidence

Codex subagents는 독립 작업의 병렬화와 context pollution 감소에 유리하다. 공식 문서도 메인 thread는 요구·결정·최종 산출물에 집중하고 탐색·테스트·로그 분석을 subagent로 분리하라고 설명한다. 단, 병렬 write-heavy 작업은 충돌 비용이 커진다.

근거: [Codex Subagents](https://learn.chatgpt.com/docs/agent-configuration/subagents), [Codex Worktrees](https://learn.chatgpt.com/docs/environments/git-worktrees)

## 3. Harness Engineering

Harness는 모델 자체가 아니라 모델이 성공하도록 둘러싼 실행 환경이다.

| 계층 | 역할 |
|---|---|
| Instruction | system prompt, AGENTS/CLAUDE.md, skill |
| Context | retrieval, summaries, relevant specs, repo map |
| Tools | shell, patch, browser/computer, MCP, LSP |
| Control | workflow state, budget, retry, approval, checkpoint |
| Environment | sandbox, credentials, worktree, services, seed data |
| Feedback | compiler, tests, telemetry, screenshots, DB assertions |
| Evaluation | rubric, independent evaluator, regression suite |
| Memory | decisions, assumptions, unresolved questions, evidence |

Anthropic의 장기 앱 개발 연구는 planner–generator–evaluator 3-agent 구조, tractable chunk 분해, session 간 structured artifact handoff를 제안한다. OpenAI의 harness engineering 사례도 인간이 intent와 feedback loop를 설계하고 agent가 실행하는 방향을 강조한다.

근거: [Anthropic Harness design for long-running apps](https://www.anthropic.com/engineering/harness-design-long-running-apps), [OpenAI Harness engineering](https://openai.com/index/harness-engineering/), [Anthropic Effective harnesses for long-running agents](https://www.anthropic.com/engineering/effective-harnesses-for-long-running-agents)

## 4. Graph Engineering

공식 Claude Code의 독립 제품 기능이라기보다 다음 그래프들을 명시적으로 관리하는 설계 패턴이다.

- Code graph: symbol, call, import, data-flow
- Spec graph: goal, actor, requirement, rule, state, invariant
- Workflow graph: task dependency와 phase transition
- Evidence graph: test run, screenshot, log, DB snapshot, verdict
- Traceability graph: spec ↔ code ↔ test ↔ incident

Codex에서 이에 대응하는 내장 단일 graph database는 공식 핵심 primitive가 아니다. 대신 code search/MCP, AGENTS.md, skills, external indexes를 조합한다. 따라서 tene가 graph schema와 query layer를 제공하는 것이 실질적 차별점이다. 다만 graph DB를 source of truth로 삼으면 Git review가 어려워진다. Markdown/YAML의 stable IDs와 links를 canonical form으로 두고 SQLite/graph DB를 파생 materialized view로 만드는 것이 좋다.

최근의 [Spec Growth Engine 논문](https://arxiv.org/abs/2606.27045)은 machine-readable spec graph, ownership-path 기반 context assembler, drift gate를 제안해 이 방향과 유사하다. 논문 결과는 독립 재현 전이므로 설계 영감으로만 사용한다.

## 5. Context Engineering

Context engineering은 “긴 prompt 작성”이 아니라 매 turn에 어떤 정보가 들어오고, 무엇이 빠지며, 어떻게 회복되는지를 설계하는 일이다.

### Claude 쪽 패턴

- CLAUDE.md와 scoped memory
- skills의 progressive disclosure
- subagent로 context 격리
- compaction과 transcript/session artifact
- hooks를 통한 선택적 context 주입
- MCP tool search와 on-demand retrieval

Anthropic은 장기 agent에서 recoverable session log와 model context를 분리하고, harness가 event stream의 필요한 slice를 변환해 주입하는 구조를 설명한다. [Managed Agents](https://www.anthropic.com/engineering/managed-agents), [Effective context engineering](https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents)

### Codex 쪽 대응

- AGENTS.md 계층: durable repository guidance
- skills progressive disclosure
- project/global `config.toml`
- MCP/connectors: 외부 최신 context
- thread resume/fork/compact
- subagents: noisy intermediate context 격리
- memories와 Computer History(표면·계정별 가용성 차이)

### tene 설계 원칙

1. 모든 문서를 매번 주입하지 않는다.
2. active feature와 current journey의 neighborhood만 context pack으로 만든다.
3. 결정에는 source message/evidence와 timestamp를 붙인다.
4. 사실, 합의, 가정, 제안, 폐기 결정을 구분한다.
5. compaction 후에도 Git artifact에서 재구성 가능해야 한다.
6. agent summary는 원본의 대체물이 아니라 index다.

