# tene CLI 통합과 Secret Security Boundary

## 1. 분석한 현재 tene CLI

경로: `/Users/kaykim/Documents/GitHub/agent-kay-it/tene`

현재 tene는 Go/Cobra 기반 local-first encrypted secret manager다.

- `.tene/vault.db`: encrypted SQLite vault
- Argon2id: master password key derivation
- XChaCha20-Poly1305: secret encryption
- OS keychain: master key cache
- 환경별 secret isolation
- `tene run -- <command>`: child process environment에만 secret 주입
- permission tiers: `metaread`, `secretwrite`, `secretread`
- fail-closed command registry: command에 tier가 없거나 stale tier가 있으면 startup panic
- audit: command tier/verb와 secret injection metadata 기록

분석한 주요 구현 위치:

| Concern | 실제 구현 |
|---|---|
| CLI entry/등록과 fail-closed dispatch | `internal/cli/root.go`, `cmd/tene/main.go` |
| permission tier source of truth | `internal/auth/permissions.go` |
| secret child-process injection | `internal/cli/run.go` |
| no-decrypt metadata listing | `internal/cli/list.go`, `internal/vault/vault.go` |
| encrypted vault/schema | `internal/vault/`, `internal/encfile/` |
| cryptography/key derivation | `pkg/crypto/` |
| keychain/fallback | `internal/keychain/` |
| agent usage contract | `skills/tene-cli/SKILL.md`, root `AGENTS.md` |

중요한 실제 호출 흐름:

```text
tene run
  → Cobra PersistentPreRunE
  → CommandTier[run] == secretread 검증
  → loadApp / master key resolution
  → encryption subkey derive
  → active env의 ciphertext 조회·복호화
  → child environment 구성
  → exec.Command 실행
  → plaintext는 child env에만 존재
```

`tene list`는 `encrypted_value`를 읽지 않고 metadata와 제한된 preview만 조회하는 metaread path다. `tene get`/plain export는 stdout에 secret을 노출할 수 있으므로 agent가 호출하면 안 된다.

## 2. Plugin Security Invariants

### SEC-01 Never read vault files

Plugin, skill, hook, graph indexer는 `.tene/**`를 열거나 search/index 대상에 포함하지 않는다. encrypted bytes도 context에 넣지 않는다.

### SEC-02 Never reveal plaintext

다음을 실행하지 않는다.

```text
tene get <KEY>
tene get <KEY> --json
tene export                  # plaintext mode
cat/read .tene/*
env | grep ...              # injected child 내부의 전체 env dump
printenv <SECRET_NAME>
```

### SEC-03 Discover by names only

필요한 key 존재 여부는 `tene list --json`으로 확인한다. preview는 secret value로 취급해 report/evidence에 복사하지 않는 것이 plugin 기본 정책이다.

### SEC-04 User enters secret

Secret set은 agent가 value를 받아 command argument로 넘기지 않는다. 사용자에게 별도 terminal에서 interactive prompt 또는 stdin을 사용하도록 안내한다.

```text
tene set OPENAI_API_KEY
# 또는 사용자가 직접
tene set OPENAI_API_KEY --stdin
```

### SEC-05 Execute through tene runtime

필요한 secret name이 존재하면 test/dev command를 다음처럼 실행한다.

```text
tene run --env <env> -- <command> <args...>
```

`--env`는 child command separator 앞에 둔다.

### SEC-06 Redact evidence

stdout/stderr/network/log/screenshot을 저장하기 전에 key-name-aware 및 high-entropy/token-pattern redaction을 수행한다. secret leakage detector가 match하면 QA gate를 fail하고 raw artifact를 보존하지 않는다.

## 3. tene CLI Adapter

```go
type SecretRuntime interface {
    IsInstalled(ctx context.Context) (bool, error)
    Status(ctx context.Context, projectDir string) (SafeStatus, error)
    ListNames(ctx context.Context, projectDir, env string) ([]SecretMeta, error)
    Run(ctx context.Context, projectDir, env string, argv []string) (RunResult, error)
}
```

Interface는 plaintext 반환 method를 제공하지 않는다. `GetSecret`이라는 함수 자체를 만들지 않는 것이 misuse 방지에 가장 효과적이다.

## 4. Command risk policy

| tene command | Plugin policy |
|---|---|
| `version`, `whoami`, `permissions` | 자동 허용 |
| `list --json`, `env list` | 자동 허용, preview 저장 금지 |
| `run -- ...` | project policy와 child command risk 평가 후 허용 |
| `init` | 사용자 요청/확인 필요 |
| `set`, `import`, `delete`, `passwd`, `recover` | 사용자 주도 interactive operation |
| `get`, plain `export` | plugin에서 금지 |
| encrypted export/import | 명시 요청과 target path 확인 |
| cloud-disabled commands | 제안하지 않음 |

## 5. Project initialization

tene secret vault가 필요한 project에서:

1. `tene version`으로 설치 확인
2. `tene whoami --json` 또는 safe status로 initialized 여부 확인
3. 미초기화이면 사용자에게 `tene init --codex` 제안
4. 생성되는 AGENTS.md와 기존 AGENTS.md가 충돌하면 overwrite 대신 merge/diff
5. `.tene/`가 Git과 graph/context scan에서 제외되는지 검사
6. workflow state는 `.tene-workflow/`에 생성

## 6. QA integration

QA environment manifest에는 secret 이름만 기록한다.

```yaml
secretRuntime:
  provider: tene
  environment: staging
  requiredNames: [DATABASE_URL, TEST_USER_TOKEN]
  availability: satisfied
```

실행 결과에는 injected count, environment, child command와 exit code만 기록하고 value는 기록하지 않는다. Tene의 audit log와 sprint evidence를 연결할 때도 secret key/value나 command arguments 중 민감 값은 포함하지 않는다.

## 7. Plugin hook safety

- PreToolUse 성격의 hook에서 `.env`, `tene get`, plain export, `.tene/` read를 차단
- shell command에 secret-looking literal이 있으면 사용자에게 tene set/run 흐름 제안
- PostToolUse에서 output leakage scan
- hook 미지원 surface에서도 skill/core CLI가 같은 정책을 적용

Hook은 defense in depth다. 핵심 안전은 SecretRuntime API가 plaintext method를 제공하지 않는 구조로 보장한다.

## 8. tene CLI에 제안할 향후 기능

Plugin 통합 신뢰성을 높이기 위해 기존 tene CLI에 다음을 검토한다.

- `tene status --json --safe`: initialized/env/count/keychain만 반환
- `tene has KEY --json`: preview 없이 존재 여부만 반환
- `tene run --audit-context <opaque-id>`: sprint/QA run correlation, secret 미포함
- `tene policy check --command-json ...`: child command risk policy 사전검사
- `tene agent doctor --json`: AGENTS rule, gitignore, unsafe `.env` 존재 확인
- structured execution manifest: injected names조차 노출하지 않고 count/env/exit만 반환

이 기능들은 현재 구현된 것으로 가정하지 않는다. MVP adapter는 기존 `list --json`, `whoami`, `permissions`, `run`만 사용한다.
