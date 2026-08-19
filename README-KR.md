# tene-codex

**tene-codex는 문서 주도(spec-driven), 기획 의도 기반(intent-aware) 에이전틱 코딩을 위한 오픈소스 Codex 플러그인입니다.**

자연스러운 코딩 대화를 Codex가 여러 세션에 걸쳐 재개하고 조사하고 검증하고 개선할 수 있는 지속 가능한 엔지니어링 workflow로 전환하는 것이 목적입니다. 생성된 코드나 통과한 test script만으로 완료를 판단하지 않고, 제품의 기획 의도, 구현 결정, 코드 영향, 사용자 여정, 데이터 흐름과 QA evidence를 작업 전체에서 연결해 관리합니다.

> 프로젝트 상태: **동작하는 pre-alpha**. 첫 번째 실행 가능한 `tene-workflow` vertical slice, Codex plugin manifest, 9개 skill, lifecycle hook, 전문 subagent profile, test와 release packaging이 포함되어 있습니다. Public API와 영속 schema는 아직 변경될 수 있으므로 유일한 production 통제로 사용하면 안 됩니다.

## Source에서 빠르게 시작하기

Go 1.24 이상, Python 3, Git과 최신 Codex가 필요합니다. 별도의 `tene` CLI는 secret이 필요한 command에서만 필수입니다.

```bash
git clone https://github.com/tene-ai/tene-codex.git
cd tene-codex
make check
go build -o dist/tene-workflow ./cmd/tene-workflow

# 적용할 repository에서
/path/to/tene-codex/dist/tene-workflow init --name my-project --profile standard
/path/to/tene-codex/dist/tene-workflow sprint create --title "My feature"
/path/to/tene-codex/dist/tene-workflow status --json
```

Plugin 개발 중에는 `scripts/tene-workflow`가 설치된 binary, bundle에 포함된 platform binary 또는 Go source command를 순서대로 찾습니다. Tag release는 macOS와 Linux binary, plugin 파일과 SHA-256 checksum을 함께 package하도록 설계했습니다.

Codex에 plugin을 설치하거나 link한 다음 `$tene-sprint`로 시작하고, `$tene-status`로 재개하며, `$tene-qa`로 evidence 기반 검증을 수행합니다. Plugin에 포함된 hook은 Codex에서 사용자가 내용을 검토하고 trust해야 실행됩니다.

주요 core command:

```text
tene-workflow status --json
tene-workflow phase transition <phase> --dry-run
tene-workflow document validate <phase>
tene-workflow context build
tene-workflow loop check
tene-workflow graph providers|build|understand|trace|validate
tene-workflow qa capabilities|plan|execute|observe|case|evaluate|status
tene-workflow evidence register|verify|list
tene-workflow report generate|validate
tene-workflow doctor|compact|clear
```

## 왜 tene-codex인가?

바이브 코딩은 빠르지만 장기간 작업에서는 최초 기획 의도를 잃기 쉽습니다. Agent가 국소적인 요청은 해결하면서 상위 정책, 하위 데이터 변경, 실패 경로 또는 사용자 화면 전이를 누락할 수 있습니다. Unit test와 scripted E2E가 통과해도 전체 기능은 잘못 동작할 수 있습니다.

tene-codex는 specification과 그 evidence를 workflow 자체의 일부로 만들어 이 문제를 해결하고자 합니다.

- 제품 기획 의도와 acceptance criteria를 여러 세션에 걸쳐 보존합니다.
- 구현부터 바로 시작하지 않고 명시적인 Sprint lifecycle을 적용합니다.
- 요구사항, 계획, 설계, task, 코드, test, evidence와 report를 traceability graph로 연결합니다.
- 네 가지 Understanding Layer와 여섯 가지 코드 이해 질문으로 전체 시스템과 변경 코드를 함께 조사합니다.
- 개별 test script뿐 아니라 UX와 데이터 처리 흐름을 검증합니다.
- Agent의 자체 완료 선언이 아니라 evidence 기반 gate로 QA 완료를 판단합니다.
- Secret 값을 model context 밖에 유지하고 값 주입을 `tene` CLI에 위임합니다.

## Sprint Workflow

의미 있는 모든 코딩 작업은 하나의 Sprint로 관리합니다.

```text
PRD → Plan → Design → Do ↔ Loop Check → QA → Report → Archive
```

- **PRD**는 문제, actor, 기획 의도, 정책, 사용자·데이터 여정, acceptance criteria와 non-goal을 기록합니다.
- **Plan**은 요구사항을 task, 의존성, 결정, 위험과 검증 작업으로 분해합니다.
- **Design**은 component, interface, data shape, 상태 전이, 실패 처리, 보안과 test seam을 정의합니다.
- **Do**는 추적 가능한 작업 단위로 구현합니다.
- **Loop Check**는 PRD, plan, design과 실제 변경을 비교하고 blocking gap이 없어질 때까지 반복 개선합니다.
- **QA**는 적용 가능한 static, unit, contract, integration, system, UX, recovery와 regression 검증을 수행하고 evidence를 수집합니다.
- **Report**는 무엇을 왜 변경했는지, 이전 Sprint와 어떻게 연결되는지, 무엇을 이월했는지 설명합니다.
- **Archive**는 모든 필수 gate와 승인을 통과한 뒤 지속 가능한 Sprint 기록을 만듭니다.

Sprint master plan은 여러 Sprint를 프로젝트 수준의 workflow와 task management system으로 연결합니다.

## Codex와 함께 동작하는 방식

tene-codex는 네 가지 Codex 연동 요소와 결정론적인 local workflow engine을 결합합니다.

### Skills

Skills는 Sprint 관리, PRD 탐색, plan, design, loop check, QA, report, status와 secret-safe 실행을 위한 사용자 진입점입니다. 사용자는 `$tene-qa`처럼 명시적으로 호출하거나, 충분히 명확한 자연어 요청으로 필요한 skill을 실행할 수 있습니다.

### Subagents

Subagent는 제품 탐색, architecture 분석, 코드 이해, 구현, test 실행과 독립 평가를 담당하는 전문 worker로 사용합니다. 가능하면 builder와 evaluator가 서로 다른 context를 사용하여, QA가 builder의 완료 주장이 아닌 evidence에 기반하도록 합니다.

### Hooks

Hooks는 lifecycle 자동화와 다중 안전장치를 제공합니다. Session 시작 시 Sprint context를 복원하고, tool 사용 후 변경 artifact를 감지하며, compaction이나 session 종료 전에 재개 가능한 상태를 기록하고, 위험한 secret 작업을 차단하며, 완료 전에 gate를 확인합니다. Codex version과 프로젝트 trust 설정에 따라 hook 지원 범위가 다를 수 있으므로 core correctness가 hook에 의존해서는 안 됩니다.

### `tene-workflow` CLI

`tene-workflow`는 다음 항목의 local source of truth입니다.

- Sprint, phase, task와 승인 상태
- 문서 scaffold와 validation
- intent, specification, code, test와 evidence의 관계
- Context Pack 생성
- Loop Check gap과 반복 이력
- QA charter, evidence manifest와 gate 판정
- compaction, recovery, migration과 archive

Skills, subagents와 hooks는 workflow 상태를 독립적으로 편집하지 않고 이 CLI를 사용해야 합니다. 이식 가능하고 결정론적으로 동작하면서 보안에 민감한 secret manager와 책임을 분리할 수 있도록 독립 Go binary로 구현했습니다.

## Engineering Model

### Understanding Layers

의미 있는 모든 변경을 다음 계층으로 조사합니다.

1. **Interface / Entry Point** — UI, CLI, API controller, webhook, scheduler 또는 command
2. **Business Logic / Processing Rules** — service, handler, use case, reducer 또는 domain rule
3. **Persistence / Data** — database, file, cache, queue 또는 external API
4. **Infrastructure / Runtime** — server, container, cloud, authentication과 CI/CD

### 여섯 가지 질문

중요한 변경 component마다 다음 질문에 답하도록 설계합니다.

1. 선언하거나 정의한 이름은 무엇인가?
2. 어떤 파일에 정의되어 있는가?
3. 어디에서 import하거나 참조하는가?
4. 어디에서 호출하거나 사용하는가?
5. 어떤 형태의 데이터를 입력받는가?
6. 어떤 형태의 데이터를 반환하거나 변경하는가?

이를 통해 단편적인 변경, 숨겨진 coupling, 불필요한 중복과 Agent가 만드는 기술 부채를 줄이고자 합니다.

## 기획 의도 기반 QA

tene-codex는 제품 기획 의도를 실행 가능한 QA 입력으로 취급합니다. Confirmed intent와 acceptance criteria를 관련 happy, alternate, empty, validation, permission, failure, retry와 recovery 경로를 포함하는 test charter로 변환합니다.

QA는 프로젝트의 기존 test, API 검사, Playwright, Codex browser 기능, Chrome 연동, database·queue observer, log, trace, screenshot과 manual checkpoint를 결합할 수 있습니다. Evidence manifest는 모든 blocking acceptance criterion을 재현 가능한 관찰 결과와 연결합니다. 평균 점수로 blocking criterion을 상쇄할 수 없으며, 모든 blocking criterion이 유효한 evidence와 함께 통과해야 Sprint를 archive할 수 있습니다.

`graph understand`는 명시적으로 요청할 때 기존 CodeGraph index를 사용하고, 그 외에는 범위가 제한된 Go AST 분석으로 fallback합니다. 각 선언의 정의 위치, import/reference, call/use, 입력 shape, 출력/side effect, Understanding Layer, provider와 confidence를 구체화합니다. `qa capabilities`는 native/Playwright runner를 발견하고, `qa execute`는 발견된 allowlist adapter만 허용합니다. Codex browser나 Chrome 도구가 생성한 UX/API/data 관찰은 `qa observe`가 schema 검증 후 evidence로 가져옵니다.

## tene를 이용한 Secret-Safe 실행

tene-codex와 `tene-workflow`는 secret 값을 소유하거나 복호화하거나 저장하지 않습니다. Secret metadata와 runtime injection은 별도의 [`tene`](https://github.com/agent-kay-it/tene) CLI에 위임합니다.

계획한 보안 경계는 다음과 같습니다.

- `.tene/**` vault 파일을 읽지 않습니다.
- Agent 자동화에서 `tene get`을 호출하지 않습니다.
- 대화에서 secret 값을 요청하지 않습니다.
- 이름 또는 capability metadata만 조사합니다.
- Secret이 필요한 test는 `tene run --env <environment> -- <command>`로 실행합니다.
- QA artifact를 evidence로 보관하기 전에 redaction과 leak scan을 수행합니다.

Workflow에 secret이 필요하지 않으면 tene secret CLI 없이도 plugin을 사용할 수 있습니다. Secret이 필요할 때는 tene 경계를 우회하지 않고 fail-closed 방식으로 중단합니다.

## 문서

- [`docs/00-rnd`](docs/00-rnd/README.md) — 시장 및 기술 조사
- [`docs/00-prd`](docs/00-prd/README.md) — 제품 요구사항 및 architecture 요구사항
- [`docs/01-plan`](docs/01-plan/README.md) — 구현 roadmap과 traceability
- [`docs/02-design`](docs/02-design/README.md) — 구현 수준의 상세 기술 설계

## 계획한 Repository 구조

```text
.codex-plugin/       Codex plugin manifest
skills/              Codex workflow skills
hooks/               선택적 lifecycle enforcement
cmd/tene-workflow/   Go CLI entry point
internal/            Workflow engine 구현
schemas/             영속 상태 및 문서 schema
templates/           Sprint 문서 template
evals/               Skill routing 및 행동 평가
docs/                조사, 요구사항, 계획과 설계
```

## 기여

아직 안정성 보장을 제공하는 단계는 아닙니다. Public API를 확립하는 동안 design review, threat modeling, workflow fixture, QA scenario와 구현 기여를 환영합니다. 모든 기여는 설계 문서에 정의된 Sprint 상태 invariant, evidence 기반 QA 모델과 secret boundary를 보존해야 합니다.

별도로 명시하지 않는 한 이 repository에 제출한 기여물은 프로젝트와 동일한 Apache License 2.0 조건으로 제공됩니다.

## 라이선스

이 repository에 포함되는 tene-codex, 계획 중인 `tene-workflow` CLI와 Codex plugin component는 **Apache License, Version 2.0**으로 제공됩니다. [`LICENSE`](LICENSE)와 [`NOTICE`](NOTICE)를 확인하십시오.

Apache-2.0은 상업적 사용, 수정, 사적 사용과 재배포를 허용합니다. 배포할 때는 적용되는 저작권, license와 NOTICE 정보를 보존하고 수정한 파일에는 중요한 변경사항을 표시해야 합니다. 명시적인 특허권 허여를 포함하며 trademark 사용 권한은 부여하지 않습니다.

이 설명은 정보 제공 목적이며 법률 자문이 아닙니다. 실제 조건은 `LICENSE`가 우선합니다.
