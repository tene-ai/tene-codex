# tene Codex plugin R&D

조사 기준일: 2026-08-20 (Asia/Seoul)

이 디렉터리는 `tene`를 **Context Engineering + Harness Engineering + Dynamic Workflow + Spec-driven QA** 플러그인으로 설계하기 위한 기술·시장 조사다.

## 문서

1. [Codex 플러그인 개발 및 Marketplace 배포](./01-codex-plugin-marketplace.md)
2. [Claude Code와 Codex의 워크플로·태스크·하네스·그래프·컨텍스트 비교](./02-claude-codex-architecture-comparison.md)
3. [OpenAI Codex의 바이브 코딩·에이전틱 코딩 방향](./03-openai-agentic-coding-direction.md)
4. [기획 의도 기반 Agentic QA 시장·방법론·아키텍처](./04-intent-driven-agentic-qa.md)
5. [tene 제품·기술 아키텍처 제안](./05-tene-proposed-architecture.md)
6. [Codex 내부 아키텍처·프로토콜·확장 포인트 심층 분석](./06-codex-deep-architecture.md)
7. [근거 수준·미확정 사항·개발 전 검증 체크리스트](./07-evidence-and-open-questions.md)

## 핵심 결론

- Codex의 배포 단위는 `.codex-plugin/plugin.json`을 가진 플러그인이며, 핵심 구성은 **skills + 선택적 remote MCP server**다. ChatGPT와 Codex는 하나의 universal plugin directory를 공유한다.
- Claude Code의 plugin은 skills 외에 agents, hooks, MCP, LSP, monitors, bin, settings까지 폭넓게 담는다. OpenAI로 변환할 때 commands/agents는 skills로 흡수하고, Claude 전용 hook·설정·로컬 stdio MCP에 의존하면 안 된다.
- “Dynamic Workflow”, “Task Management”, “Graph Engineering”은 Claude Code의 단일 공식 제품명이 아니라 여러 primitive를 조합한 아키텍처 패턴으로 보는 편이 정확하다. Codex에도 대응 primitive가 대부분 있지만 동작과 배포 표면은 다르다.
- tene의 차별점은 문서를 많이 만드는 데 있지 않다. **대화에서 추출한 의도 → 버전·근거가 있는 spec graph → 실행 가능한 journey/contract oracle → UI·API·데이터 증거 → spec-code-test traceability**의 폐루프가 핵심이다.
- Markdown/Git을 source of truth로 두고 SQLite/graph index는 재생성 가능한 파생 인덱스로 두는 것이 이식성, 리뷰 가능성, 장애 복구에 유리하다.
- tene를 단순 skill bundle로 끝내면 장기 상태·증거·동적 orchestration이 취약하다. UX는 plugin skills로 제공하되, 실제 control plane은 repository artifacts와 필요 시 App Server client/remote MCP로 분리하는 것이 적합하다.

## 조사 방법과 한계

공식 OpenAI Codex manual과 OpenAI/Anthropic 공식 문서·엔지니어링 글을 1차 근거로 사용하고, GitHub 저장소·이슈와 논문·경쟁 제품 자료를 보조 근거로 교차 검증했다. 제품 기능과 제출 요건은 빠르게 바뀌므로 실제 배포 직전에 공식 제출 포털과 최신 문서를 다시 확인해야 한다. GitHub 이슈는 확정된 제품 사양이 아니라 운영상 한계와 사용자 요구를 보여주는 증거로만 사용했다.
