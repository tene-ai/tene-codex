# 근거 수준·미확정 사항·개발 전 검증 체크리스트

## 1. 근거 수준

| 등급 | 의미 | 사용 원칙 |
|---|---|---|
| A | OpenAI/Anthropic 공식 제품 문서·source | 구현 계약의 1차 근거 |
| B | 공식 engineering blog·공식 sample | 아키텍처 방향의 근거, API 계약은 문서 재확인 |
| C | upstream GitHub source/merged PR | 현재 구현 확인, release version과 비교 필요 |
| D | GitHub issue/discussion | 문제·수요의 증거, 제품 동작으로 단정 금지 |
| E | 논문·경쟁사·커뮤니티 | 설계 아이디어와 시장 비교, 독립 검증 필요 |

## 2. 핵심 source ledger

| 주제 | 출처 | 등급 |
|---|---|---|
| Codex agent loop/context | [Unrolling the Codex agent loop](https://openai.com/index/unrolling-the-codex-agent-loop/) | B |
| Core/App Server 구조 | [Unlocking the Codex harness](https://openai.com/index/unlocking-the-codex-harness/) | B |
| App Server protocol | [공식 문서](https://developers.openai.com/codex/app-server), [source README](https://github.com/openai/codex/blob/main/codex-rs/app-server/README.md) | A/C |
| Plugins/skills | [Build plugins](https://developers.openai.com/codex/build-plugins), [Build skills](https://developers.openai.com/codex/build-skills) | A |
| Claude plugin 전환 | [Submit Claude plugin](https://developers.openai.com/plugins/guides/submit-claude-plugin) | A |
| Agent-first repo | [Harness engineering](https://openai.com/index/harness-engineering/) | B |
| Task fleet orchestration | [Symphony](https://openai.com/index/open-source-codex-orchestration-symphony/) | B |
| Claude context/harness | [Context engineering](https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents), [Managed Agents](https://www.anthropic.com/engineering/managed-agents) | B |
| SDD baseline | [GitHub Spec Kit](https://github.github.io/spec-kit/index.html) | A |
| Memory/compaction gaps | [Codex #29816](https://github.com/openai/codex/issues/29816), [#22861](https://github.com/openai/codex/issues/22861) | D |

## 3. 아직 확정하면 안 되는 사항

- universal marketplace의 review SLA, ranking, monetization 정책
- plugin hook이 모든 ChatGPT/Codex surface에서 동일하게 실행되는지
- public plugin에서 local stdio MCP 지원 시점
- App Server protocol의 experimental event 안정성
- hook additional context의 현재 exact limit와 장기 보장
- Codex memory/Computer History의 계정·surface별 보장 범위
- proactive multi-agent delegation 정책과 한도
- 경쟁 QA 제품의 실제 탐지율·비용·data policy

## 4. 개발 전 spike

### Spike A — Skills-only plugin

- 4개 minimal skill을 package
- explicit/implicit routing confusion matrix 측정
- 새 thread, compact 후, repo root/subdir에서 반복
- Claude/OpenAI portable content 차이 기록

### Spike B — Durable resume

- 50-step workflow 중간에 process kill
- 새 Codex thread에서 `.tene/run.yaml`만으로 재개
- compact 전후 active constraints loss 측정
- duplicate/replayed action idempotency 확인

### Spike C — App Server

- initialize → thread/start → turn/start → event stream → approval → resume 구현
- completed/interrupted/failed 처리
- schema generation 및 pinned-version compatibility test
- process crash와 reconnect 처리

### Spike D — Intent QA

- checkout journey 하나 선정
- 대화에서 requirement/rule/invariant 추출
- UI/API/DB에 각기 다른 결함 삽입
- baseline E2E와 tene 탐지율·오분류 비교

## 5. Go/No-Go 기준

MVP 개발은 다음이 확인되면 Go다.

- skills-only plugin이 핵심 workflow를 hook 없이 수행
- confirmed spec과 run state가 thread 밖에서 복구 가능
- validator가 invalid transition과 source 없는 intent를 차단
- 한 journey에서 UI·API·DB evidence correlation 성공
- evaluator가 defect/spec ambiguity/environment failure를 구분

다음이면 설계를 수정한다.

- implicit skill routing이 핵심 상태 전이를 결정
- compaction 후 hook 없이는 재개 불가
- graph index 없이 canonical artifact를 읽을 수 없음
- 자동 spec update가 사용자 의도를 덮어씀
- browser healing이 실제 UX defect를 숨김

