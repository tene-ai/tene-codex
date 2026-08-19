# Skill Contracts and Routing

## 1. Skill catalog

| Skill | 명시 호출 | 책임 | 주요 core command |
|---|---|---|---|
| `tene-sprint` | `$tene-sprint` | init/create/resume/master/archive | `sprint`, `phase` |
| `tene-prd` | `$tene-prd` | 대화에서 intent/AC/journey 확정 | `intent`, `document` |
| `tene-plan` | `$tene-plan` | 작업·의존성·검증 계획 | `task`, `document` |
| `tene-design` | `$tene-design` | 코드 조사와 상세 설계 | `graph`, `document` |
| `tene-loop-check` | `$tene-loop-check` | spec↔code gap 반복 개선 | `loop`, `graph` |
| `tene-qa` | `$tene-qa` | 7-layer QA/evidence/evaluate | `qa`, `evidence` |
| `tene-report` | `$tene-report` | review/회고/report | `report` |
| `tene-status` | `$tene-status` | 상태와 next action | `status`, `context` |
| `tene-secrets` | `$tene-secrets` | secret-safe 실행 안내/위임 | `tene list/run` |

Codex에는 임의의 slash-command 시스템이 있다고 가정하지 않는다. 공식적인 explicit surface는 skill 이름 호출이며, UX 문서에서 `$tene-*`로 표기한다.

## 2. 공통 SKILL contract

각 `SKILL.md`는 다음 순서를 MUST 따른다.

1. `status --json`으로 project/Sprint/phase 확인.
2. phase가 다르면 자동 mutation하지 말고 가능한 transition 또는 올바른 skill을 제안.
3. `context build`로 해당 phase pack만 로드.
4. 필요한 코드/문서/도구를 조사하고 structured draft 생성.
5. 사용자 정책/intent confirmation이 필요한 지점만 승인 요청.
6. core CLI로 validate/register/transition.
7. changed artifacts, verdict, next action을 요약.

상세 schema/template/tool matrix는 `references/`에 두고 SKILL 본문은 routing과 workflow에 집중한다. 결정론적 파싱/검증은 `scripts/` 또는 core에 둔다.

## 3. Skill input/output

```ts
interface SkillInvocation {
  user_request:string; explicit:boolean; repo_root:string;
  requested_sprint?:string; requested_phase?:string;
}
interface SkillOutcome {
  selected_skill:string; actions_taken:string[]; artifact_paths:string[];
  state_revision:number; gate?:"passed"|"failed"|"needs-approval";
  next_actions:string[]; warnings:string[];
}
```

Skill은 내부 추론을 state에 기록하지 않고 결정, 근거 locator, 산출물, evidence만 기록한다.

## 4. Implicit router

Router score:

`0.40 intent match + 0.25 phase compatibility + 0.20 artifact cue + 0.15 action cue`.

- score ≥ 0.80이고 단일 winner: 해당 skill 실행.
- 0.60~0.79 또는 top-2 차이 < 0.10: 상태 조회만 하고 후보/이유를 제안.
- < 0.60: 자동 호출하지 않는다.
- mutation은 implicit match만으로 phase를 건너뛰거나 승인 gate를 우회하지 않는다.

### Positive cues

- PRD: 요구사항, 기능 의도, 사용자 문제, acceptance criteria
- Plan: 구현 계획, task 분해, 의존성, 일정
- Design: 아키텍처, 데이터 구조, 인터페이스, 처리 로직
- Loop: 빠진 요구, 설계대로, 반복 검증, gap
- QA: UX 흐름, 종합 테스트, Playwright, 데이터 흐름
- Report: 회고, 작업 결과, 변경 파일, 이월
- Status: 어디까지, 이어서, 현재 단계
- Secrets: 비밀정보, API key, credential, env injection

### Hard negatives

일반 지식 질문, 코드 변경 없는 설명, 단순 “계획이 뭐야?”에는 Sprint mutation을 자동 실행하지 않는다. secret value를 보여달라는 요청은 routing 성공 여부와 무관하게 거절하고 안전한 `tene run` 방식을 제시한다.

## 5. `agents/openai.yaml`

각 skill은 `display_name`, 짧은 `description`, `$skill-name`을 포함한 `default_prompt`를 갖는다. implicit invocation은 기본 활성화하되 `tene-secrets`는 명확한 secret/runtime cue가 있을 때만 선택되도록 description과 eval을 강화한다.

## 6. Cross-skill handoff

handoff는 자연어 요약이 아니라 `{from_phase,to_phase,revision,context_pack_id,open_items}`를 사용한다. 다음 skill은 revision mismatch면 context를 다시 생성한다. QA fail은 원인에 따라 `do` 또는 `loop-check`로 돌아가며 새 Sprint를 만들지 않는다.

