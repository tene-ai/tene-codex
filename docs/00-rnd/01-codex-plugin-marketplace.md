# Codex 플러그인 개발 및 Marketplace 배포 연구

## 1. 현재 Codex 플러그인의 정의

OpenAI 공식 문서에서 플러그인은 ChatGPT와 Codex에 설치할 수 있는 패키지다. 최소 하나 이상의 skill을 포함하며, 필요하면 remote MCP server와 UI를 결합한다. 공개 플러그인은 ChatGPT와 Codex가 공유하는 universal plugin directory에 한 번 게시한다. 개인 워크플로를 반복 검증하는 단계는 standalone skill, 팀·외부 배포가 필요해진 단계는 plugin이 적합하다.

근거: [OpenAI Build plugins](https://developers.openai.com/codex/build-plugins), [OpenAI Plugin architecture](https://developers.openai.com/plugins/concepts/plugin-architecture), [OpenAI Build skills](https://developers.openai.com/codex/build-skills)

### 권장 기본 구조

```text
tene/
├── .codex-plugin/
│   └── plugin.json
├── skills/
│   ├── intent-interview/SKILL.md
│   ├── spec-maintain/SKILL.md
│   ├── workflow-run/SKILL.md
│   └── qa-journey/SKILL.md
├── scripts/                 # 결정론적 validator, graph/index builder
├── references/              # schema, rubric, templates
├── assets/                  # 문서 템플릿
├── hooks/                   # Codex에서만 쓰는 보조 enforcement
├── .mcp.json                # 필요할 때만; 공개 배포는 remote HTTP 우선
├── README.md
└── LICENSE
```

Skill은 `SKILL.md`의 `name`, `description`, 본문 지침과 선택적 `scripts/`, `references/`, `assets/`로 구성된다. 호스트는 먼저 name/description만 노출하고 선택된 skill의 전체 내용을 나중에 읽는 progressive disclosure 방식을 쓴다. 따라서 description은 단순 소개가 아니라 라우팅 규칙이다.

## 2. 개발 절차

### 2.1 skill-first로 검증

1. 대표 사용자 시나리오를 2~3개 정한다.
2. `.agents/skills/<skill-name>/SKILL.md`에서 로컬 반복 검증한다.
3. 입력·출력·실패 조건·완료 조건을 명시한다.
4. 자연어 판단은 skill에, 재현성이 필요한 검증·변환은 script에 둔다.
5. 외부 실시간 데이터나 권한 있는 action만 MCP tool로 둔다.
6. 안정화 후 `.codex-plugin/plugin.json`으로 묶는다.

Codex는 명시적 `$skill-name` 호출과 description 기반 암시적 호출을 모두 지원한다. 변경 감지는 자동이지만 반영되지 않으면 Codex 재시작이 필요할 수 있다.

### 2.2 plugin manifest

정확한 schema는 제출 시점의 [Package your plugin](https://developers.openai.com/plugins/build/package)을 따라야 한다. 최소한 이름, 설명, 버전, 구성 요소 경로가 명확해야 한다. 버전은 SemVer를 권장하고 manifest, changelog, release tag를 일치시킨다.

### 2.3 로컬 Marketplace 테스트

OpenAI는 개발 중 local marketplace로 패키지를 설치해 새 대화에서 대표 요청을 시험하도록 안내한다. `$plugin-creator`가 manifest, 폴더 구조, local marketplace entry 생성을 지원한다. 검증 항목은 다음과 같다.

- 명시적 호출과 암시적 호출 모두 올바른 skill을 선택하는가
- 비관련 요청에는 오작동하지 않는가
- 깨끗한 새 thread에서도 숨은 대화 문맥 없이 동작하는가
- script 경로가 package-relative인가
- 네트워크·권한 실패가 안전하게 종료되는가
- core workflow가 hook 없이도 가능한가
- 결과가 ChatGPT와 Codex 양쪽 표면에서 이해 가능한가

근거: [Connect and test your plugin](https://developers.openai.com/plugins/test/connect), [Plugin guidelines](https://developers.openai.com/plugins/resources/guidelines)

## 3. Marketplace 제출

일반적인 흐름은 다음과 같다.

1. package와 README, privacy/security 설명, support URL을 준비한다.
2. Skills-only 또는 Skills + remote MCP 경로를 선택한다.
3. [OpenAI plugin submission portal](https://platform.openai.com/)에서 archive 또는 요구 형식으로 제출한다.
4. portal이 생성·정규화한 `.codex-plugin/plugin.json`을 검토한다.
5. 자동 validation, 보안·정책 review, metadata review를 통과한다.
6. 승인 후 universal directory에 publish하고 버전별 회귀 테스트를 운영한다.

MCP가 있으면 public Streamable HTTP endpoint, 인증·개인정보 처리, tool metadata, destructive action confirmation이 review의 핵심이다. 로컬 stdio-only MCP는 공개 제출의 기본 경로가 아니다.

## 4. Claude Code 플러그인과의 호환성

Claude 공식 plugin은 `.claude-plugin/plugin.json`과 함께 skills, agents, hooks, `.mcp.json`, `.lsp.json`, monitors, bin, settings를 가질 수 있다. Marketplace는 `.claude-plugin/marketplace.json` 카탈로그를 GitHub/GitLab/URL/local path로 배포하며, 사용자가 `/plugin marketplace add` 후 설치한다.

근거: [Claude Code Create plugins](https://code.claude.com/docs/en/plugins), [Create and distribute a plugin marketplace](https://code.claude.com/docs/en/plugin-marketplaces), [Plugins reference](https://code.claude.com/docs/en/plugins-reference)

OpenAI 공식 변환 가이드의 핵심 차이는 아래와 같다.

| Claude 구성 | OpenAI 제출 전략 |
|---|---|
| `skills/` | 대체로 유지; provider-neutral 문구 사용 |
| `commands/`, `agents/` | 재사용 절차를 skill로 변환 |
| command hooks | Codex hook runtime에 맞게 적응하되 core workflow 필수조건으로 삼지 않음 |
| prompt/agent hooks | Codex에서 미지원인 경우가 있으므로 skill/MCP/결정론적 script로 재설계 |
| local stdio MCP | 공개 HTTP MCP로 노출하거나 공개 배포에서 제외 |
| `userConfig` | 설치 프롬프트·변수 치환에 의존하지 말고 MCP 인증/설정으로 이동 |
| Claude live artifacts | 일반 대화 출력 또는 portable file artifact로 대체 |
| `.claude-plugin/plugin.json` | 직접 archive 제출 시 portal이 `.codex-plugin/plugin.json`으로 변환 가능 |

근거: [Submit your Claude Code plugin to OpenAI](https://developers.openai.com/plugins/guides/submit-claude-plugin)

## 5. tene에 대한 결론

tene는 처음부터 이중 호스트를 고려하되 공통 core를 Agent Skills 표준에 맞춰야 한다.

- 공통 계층: skills, Markdown schemas, templates, deterministic scripts
- Codex 계층: `.codex-plugin/plugin.json`, Codex hooks, optional remote MCP
- Claude 계층: `.claude-plugin/plugin.json`, agents/hooks/marketplace metadata
- 절대 원칙: 의도·spec·QA의 핵심 상태 전이는 특정 host hook에만 의존하지 않는다.

이 구조가 중요한 이유는 OpenAI 플러그인이 Claude plugin 전체 runtime을 그대로 재현하는 컨테이너가 아니기 때문이다. “Claude plugin을 먼저 만들고 자동 변환”보다 “portable core + host adapters”가 유지보수 비용이 낮다.

