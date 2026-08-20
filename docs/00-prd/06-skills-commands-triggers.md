# Skills·Commands·자연어 Trigger 명세

## 1. Invocation 원칙

모든 skill은 명시적 `$tene:*` 호출과 implicit invocation을 기본 지원한다. 단, `archive --purge`, production QA, secret mutation처럼 위험한 action은 skill discovery를 막는 것이 아니라 실행 직전에 authorization을 요구한다.

Skill description은 broad catch-all이 되지 않도록 실제 trigger와 boundary를 명확히 한다.

## 2. Skill catalog

| Skill | 책임 | 대표 명시 호출 |
|---|---|---|
| `sprint` | master plan, sprint 생성·시작·transition·archive | `$tene:sprint` |
| `prd` | intent interview, PRD, policy/open question | `$tene:prd` |
| `plan` | task/dependency/risk/verification plan | `$tene:plan` |
| `design` | layer/call/data contract와 ADR | `$tene:design` |
| `loop-check` | 문서↔code gap 반복 개선 | `$tene:loop-check` |
| `qa` | intent journey QA와 gate | `$tene:qa` |
| `report` | 구현 mapping, 회고, deferred/policy | `$tene:report` |
| `status` | resume/status/next action | `$tene:status` |
| `secrets` | tene CLI를 통한 secret-safe 실행 | `$tene:secrets` |

## 3. 자연어 Trigger

### `sprint`

긍정 trigger:

- “이 기능 sprint 시작해줘”
- “전체 개발 workflow/master plan 만들어줘”
- “작업들을 여러 sprint로 나눠줘”
- “현재 sprint를 archive해줘”

부정 boundary: 단순 일정 질문이나 일반 agile 설명에는 호출하지 않는다.

### `prd`

- “요구사항 정리해줘”, “기능 기획해줘”
- “내 아이디어를 구체화해줘”, “먼저 질문해줘”
- “정책과 예외 케이스를 정리해줘”

### `plan`

- “구현 계획 세워줘”, “task로 나눠줘”
- “어떤 순서로 작업해야 해?”
- “병렬 가능한 작업과 dependency 찾아줘”

### `design`

- “처리 로직 설계해줘”, “architecture/data flow 설계”
- “어떤 파일과 symbol을 바꿀지 설계”
- “API/DB/event contract 정해줘”

### `loop-check`

- “문서대로 구현됐는지 확인해줘”
- “100% 될 때까지 반복 수정해줘”
- “PRD/plan/design과 코드 gap 찾아줘”

### `qa`

- “기획 의도대로 동작하는지 QA해줘”
- “UX와 데이터 흐름까지 테스트해줘”
- “Chrome/Playwright로 전체 사용자 journey 검증해줘”

### `report`

- “작업 결과와 구현 파일 정리해줘”
- “sprint 회고/report 작성해줘”
- “이전 기능과 어떻게 이어지는지 설명해줘”

### `status`

- “어디까지 했지?”, “이어서 작업해줘”
- “현재 sprint 상태와 다음 할 일”
- “blocker와 이월 작업 보여줘”

### `secrets`

- secret, API key, credential, token, `.env`, environment variable
- secret이 필요한 dev/test/deploy command 실행

## 4. Router logic

```text
1. Explicit $skill mention → 해당 skill
2. Active sprint가 있고 요청이 현재 phase와 일치 → current phase skill
3. Natural language가 단일 phase trigger와 명확히 일치 → 해당 skill
4. 구현 요청인데 active sprint 없음 → sprint가 최소 상태 확인
5. 구현 요청인데 PRD/design 없음 → 현재 phase gate가 다음 필요 skill 안내
6. 복수 intent → sprint가 workflow를 조정하고 phase skill 순차 호출
7. 위험/불가역 action → 실행 전 사용자 확인
```

Plugin은 unrelated coding request를 무조건 sprint로 납치하면 안 된다. 사용자가 “tene 없이 빠르게 한 줄만 수정”처럼 명시적으로 opt-out하면 경고 후 one-off를 허용할 수 있다. 단, project policy가 mandatory mode라면 waiver를 기록한다.

## 5. Commands/Core CLI UX

```text
tene-workflow init
tene-workflow master create|status|validate
tene-workflow sprint create|start|status|transition|archive|clear
tene-workflow docs scaffold|validate
tene-workflow intent propose|confirm|deprecate|diff
tene-workflow graph build|query|impact|drift
tene-workflow context build --phase <phase>
tene-workflow gate check|waive
tene-workflow qa plan|record|judge
tene-workflow report generate|validate
tene-workflow doctor
```

모든 read command는 `--json`을 지원하고 write command는 `--dry-run`을 지원한다. transition과 archive는 idempotency key를 받는다.

## 6. Skill progressive disclosure

각 `SKILL.md`에는 다음만 둔다.

- 언제 trigger하는지
- 현재 state 확인 방법
- phase의 핵심 workflow와 stopping condition
- 반드시 지켜야 할 invariant
- 필요한 reference routing

상세 template, schema, QA rubric, stack adapter는 `references/`에서 필요할 때만 읽는다. 모든 phase를 하나의 giant skill에 넣지 않는다.

## 7. Suggested `agents/openai.yaml`

각 skill은 `allow_implicit_invocation: true`를 유지한다. UI metadata의 default prompt는 반드시 실제 `$skill-name`을 언급한다. 예:

```yaml
interface:
  display_name: "Tene QA"
  short_description: "기획 의도 기반 종합 QA와 gate 판정"
  default_prompt: "Use $tene:qa to verify this sprint's UX and data flow."
policy:
  allow_implicit_invocation: true
```

## 8. Trigger 품질 평가

각 skill마다 최소 다음 eval set을 유지한다.

- 20 positive natural-language cases
- 20 adjacent negative cases
- 10 multi-intent routing cases
- 10 active-state conflict cases
- Korean/English mixed prompts
- explicit invocation cases

측정치는 precision, recall, wrong-phase rate, unnecessary-trigger rate다.

