# tene Secret Runtime Integration

## 1. Security boundary

Workflow plugin은 secret **값**을 데이터 타입, prompt, document, state, graph, evidence에 표현하지 않는다. 필요한 것은 환경 이름과 secret key 이름뿐이다. 복호화와 child environment injection은 기존 tene CLI 프로세스 내부에서만 일어난다.

## 2. 분석 기반 허용 surface

현재 tene CLI의 `internal/cli/root.go`는 command와 권한 검사를 구성하고, `internal/auth/permissions.go`는 metaread/secretwrite/secretread 권한을 fail-closed로 검증한다. `internal/cli/list.go`는 복호화 없이 metadata를 제공하며, `internal/cli/run.go`는 복호화한 값을 child process 환경에만 주입한다. plugin은 이 경계를 그대로 이용한다.

### Allow

- `tene version`, `tene whoami`
- `tene list --env <name> --json`처럼 값 없는 metadata 조회
- `tene run --env <name> -- <executable> <args...>`
- 사용자가 직접 수행하는 `tene init/set/import/passwd/recover`

### Deny from agent automation

- `tene get`
- plaintext `tene export`
- `.tene/**` 또는 vault/keychain 직접 read
- `env`, `printenv`, debug dump 등 child secret environment 출력
- command string을 shell로 재해석하는 `sh -c`/`eval` 형태

## 3. Secret reference

```ts
interface SecretRef {
  provider:"tene"; environment:string; key_name:string;
  purpose:string; required_by:string[];
}
```

value/preview/ciphertext field는 schema에서 `additionalProperties:false`로 거절한다. key name도 민감할 수 있어 report에는 policy에 따라 alias를 쓸 수 있다.

## 4. Runner algorithm

1. `exec.LookPath("tene")`와 minimum version 확인.
2. argv를 array로 구성하고 shell을 사용하지 않는다.
3. requested env/key name을 metadata 목록으로 preflight한다.
4. test command가 env dump/known exfiltration pattern인지 차단한다.
5. `tene run --env ENV -- COMMAND ARGS...`를 직접 실행한다.
6. stdout/stderr는 streaming redactor를 거친 뒤 필요한 최소분만 evidence로 저장한다.
7. child exit code, duration, sanitized command, env alias만 기록한다.
8. canary/secret leak scan 실패 시 artifact 격리, QA fail, 후속 단계 중단.

키 이름 목록으로 실제 값 redaction은 불가능하므로 1차 방어는 값을 parent/model에 전달하지 않는 구조다. redactor는 token/password 패턴과 test canary에 대한 보조 방어다.

## 5. Fail-closed behavior

| 상황 | 동작 |
|---|---|
| tene 미설치/버전 불일치 | dependency error, 설치 안내 |
| env/key 없음 | 이름만 표시, 사용자에게 `tene set` 직접 실행 요청 |
| 권한 거부 | 우회 금지, sanitized 오류 |
| child 실패 | exit code 보존, log 최소화, QA fail |
| leak 감지 | 즉시 중단·artifact quarantine·security blocker |

plugin은 사용자에게 secret 값을 채팅으로 입력하라고 요청하지 않는다. 대화에 값이 들어오면 재출력/저장하지 말고 rotation과 `tene set`을 안내한다.

## 6. Configuration

```yaml
secrets:
  provider: tene
  default_environment: test
  required:
    - alias: external-api
      key_name: EXTERNAL_API_KEY
      purposes: [qa-e2e]
```

config에는 값이 없다. QA manifest는 `secret_env_name`만 기록하며 command args에 secret 자체가 포함되는 도구는 adapter로 허용하지 않는다.

## 7. Future tene CLI requests

추후 tene에 `list --names-only --json`, `run --audit-json-fd`, `policy check <argv>`, non-interactive capability probe가 추가되면 adapter가 사용한다. 지원되기 전에는 존재한다고 가정하지 않는다.

