# Codex Plugin 기술 아키텍처

## 1. 아키텍처 원칙

tene plugin은 세 층으로 나눈다.

```text
Experience Layer       Codex Skills / starter prompts / commands
Control Layer          tene-workflow CLI + validators + state machine
Integration Layer      Codex hooks, App Server adapter, MCP, tene CLI
```

Skill만으로 상태와 gate를 관리하면 LLM 비결정성에 취약하다. 반대로 모든 것을 CLI에 넣으면 intent interview와 문서 합성이 경직된다. 자연어 판단은 skill, 상태·검증·secret policy는 deterministic CLI가 맡는다.

## 2. 제안 Plugin 구조

```text
tene/
├── .codex-plugin/
│   └── plugin.json
├── skills/
│   ├── sprint/
│   ├── prd/
│   ├── plan/
│   ├── design/
│   ├── loop-check/
│   ├── qa/
│   ├── report/
│   ├── status/
│   └── secrets/
├── scripts/
│   └── tene-workflow             # portable core CLI launcher/binary
├── references/
│   ├── schemas/
│   ├── templates/
│   ├── qa-rubrics/
│   └── host-behavior/
├── hooks/
│   └── hooks.json                # optional mechanical checks
├── assets/
├── .mcp.json                     # optional remote integrations only
└── agents/
```

`plugin.json`에는 실제 존재하는 component만 선언한다. 현재 validator가 허용하지 않는 field를 임의로 넣지 않는다. hooks는 companion discovery/runtime 규칙을 실제 Codex version에서 검증한다.

## 3. Core CLI modules

```text
core/
├── project       init/config/discovery
├── sprint        create/start/status/archive/clear
├── state         transition/checkpoint/resume/event log
├── docs          scaffold/validate/link/hash
├── graph         extract/index/query/impact/drift
├── context       assemble/budget/redact
├── gates         evaluate/waive/report
├── qa            charter/evidence/verdict
├── secrets       tene CLI adapter/policy
└── adapters      git/codex/browser/test runners
```

## 4. Codex 기능 배치

| Codex 기능 | 활용 | 비고 |
|---|---|---|
| AGENTS.md | repo-level tene 사용 규칙과 완료 정의 | plugin 설치 시 무조건 overwrite하지 않음; 제안 diff |
| Skills | phase별 인터뷰·문서·판단 UX | progressive disclosure |
| Hooks | transition 전 validator, 종료 시 state/evidence check | core workflow 필수 의존 금지 |
| Subagents | code graph 조사, test 실행, independent evaluation | main thread는 intent/decision 유지 |
| Worktrees | 병렬 sprint/task 변경 격리 | write-heavy 충돌 방지 |
| MCP | tracker/Figma/log/DB 등 외부 context/action | secret value 전달 금지 |
| `codex exec` | CI drift/QA/report | JSONL/event result 보관 |
| App Server | rich dashboard, approval, thread/turn event mapping | Advanced edition |
| Scheduled tasks | nightly drift/QA/stale-doc scan | stable workflow에만 적용 |
| Memory | 사용자 선호와 문서 스타일 | sprint progress 저장 금지 |

## 5. Plugin과 core의 contract

Skill은 직접 state JSON을 편집하지 않는다. 다음 흐름을 사용한다.

```text
Skill asks/interprets
  → produces candidate artifact
  → core CLI validates candidate
  → user confirms policy/intent if needed
  → core CLI commits transition atomically
  → skill reports next action
```

예시:

```text
$tene:prd
  1. status --json
  2. intent interview
  3. write PRD draft
  4. docs validate --type prd
  5. intent diff --proposed
  6. user confirmation
  7. intent confirm + transition plan
```

## 6. App Server 사용 여부

MVP는 local plugin + core CLI로 충분하다. App Server는 다음 요구가 생길 때 도입한다.

- 여러 Codex thread를 dashboard에서 동시에 관리
- thread/turn/item event와 sprint task를 실시간 mapping
- approval request UI
- agent interrupt/retry/resume/fork
- aggregated diff/plan streaming

App Server event는 `.tene-workflow/events.ndjson`의 domain event로 변환하되 raw transient event를 canonical state로 사용하지 않는다.

## 7. Portability

- Core schemas와 templates는 provider-neutral
- Codex-specific invocation은 `agents/openai.yaml`, hook, App Server adapter에 격리
- Claude adapter는 commands/agents/hooks 차이를 별도 package로 처리
- repository discovery는 Git root가 없어도 current directory project로 동작
- code graph provider는 CodeGraph → LSP/AST → `rg` fallback interface를 제공

## 8. Plugin manifest 방향

초기 배포는 skills-only + bundled scripts가 적합하다. 외부 central registry가 필요해질 때 remote MCP를 추가한다. Marketplace용 metadata에는 다음이 필요하다.

- Name: `tene`
- Category: Productivity 또는 Developer Tools 가용 category
- Capabilities: spec workflow, project write, local command execution
- Authentication: local mode는 없음; remote MCP만 OAuth
- Starter prompts: sprint 시작, 현재 상태, 기획 의도 QA
- Privacy: `.tene/` never-read, no plaintext secret, repository-local artifacts

근거: [Codex Build plugins](https://developers.openai.com/codex/build-plugins), [Build skills](https://developers.openai.com/codex/build-skills), [Hooks](https://learn.chatgpt.com/docs/config-file/advanced#hooks), [App Server](https://developers.openai.com/codex/app-server)

