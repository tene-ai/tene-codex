# Codex Integration and Hook Boundaries

## 1. Codex surface 배치

| Codex surface | 사용 | 의존성 수준 |
|---|---|---|
| Plugin | skills/references/scripts 배포 단위 | 필수 |
| Skills | 자연어/명시 호출과 phase orchestration | 필수 |
| AGENTS.md | 프로젝트별 규칙·strict path·명령 | 선택/권고 |
| `codex exec` | CI/headless gate 실행 | 선택/지원 |
| MCP | browser/CodeGraph/외부 observer adapter | 선택 |
| App Server | 향후 UI·durable controller transport | post-MVP |
| Hooks | stop/command safety defense-in-depth | 선택 |
| Subagents | independent evaluator/병렬 조사 | capability 있을 때 |

## 2. Plugin discovery contract

`.codex-plugin/plugin.json`은 manifest source이며 skills는 plugin root의 `skills/`에 둔다. manifest에 공식 schema가 지원하지 않는 `hooks` 같은 임의 field를 추가하지 않는다. 설치 환경에서 validator와 실제 Codex discovery를 모두 확인한다.

## 3. AGENTS.md integration

`init`은 사용자의 명시 동의가 있을 때 managed block을 제안한다.

```markdown
<!-- TENE_WORKFLOW_START -->
For feature, bug, and refactor work, resume the active tene sprint first.
Do not skip PRD/Plan/Design/Loop Check/QA/Report gates.
Never read .tene/** or expose secrets; use tene run for secret-required commands.
<!-- TENE_WORKFLOW_END -->
```

기존 지시를 덮어쓰지 않으며 conflict가 있으면 doctor warning과 사용자 결정을 요구한다.

## 4. Hook design

Hook을 제공할 수 있는 Codex 버전에서는 다음만 수행한다.

- **pre-command**: `.tene/**`, `tene get`, plaintext export 같은 금지 패턴 차단.
- **post-tool/command**: changed file metadata와 exit code만 기록; stdout 원문은 evidence policy를 통과한 경우에만 저장.
- **stop/turn-complete**: active task/phase가 dirty면 compact context와 next action을 생성.

Hook 실패는 workflow state를 직접 수정하지 않으며 core command를 호출한다. hook 미지원 환경에서도 모든 invariant가 core/skill preflight로 유지되어야 한다.

## 5. App Server future adapter

App Server를 채택할 경우 thread/turn/item event를 domain event와 직접 동일시하지 않는다. adapter가 Codex event를 normalize하고, 승인된 action만 workflow event로 commit한다. reconnect 시 `last_seen_sequence` 이후를 replay하고 중복 request ID를 제거한다.

## 6. Subagent/evaluator isolation

가능하면 builder와 evaluator를 다른 agent/context로 실행한다. evaluator context에는 PRD/AC, diff summary, evidence manifest만 제공하고 builder의 “완료” 서술은 제외한다. subagent가 없으면 새 isolated `codex exec` 또는 deterministic evaluator로 fallback하며 capability를 report에 기록한다.

## 7. Version capability negotiation

`doctor`는 Codex version과 plugin discovery, skill invocation, hook/MCP/App Server availability를 probe해 capability snapshot을 저장한다. 설계는 version string 분기보다 실제 probe 결과를 우선한다. 공식 문서와 validator 변경은 compatibility test fixture로 관리한다.

