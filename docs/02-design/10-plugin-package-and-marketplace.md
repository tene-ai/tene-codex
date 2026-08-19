# Plugin Package and Marketplace Design

## 1. Package layout

```text
tene/
├── .codex-plugin/plugin.json
├── README.md
├── LICENSE
├── CHANGELOG.md
├── skills/
│   ├── tene-sprint/{SKILL.md,agents/openai.yaml,references/}
│   └── ...
├── references/{schemas,security,workflow}/
├── scripts/{invoke-core,doctor}
├── bin/{darwin-arm64,linux-amd64,...}/tene-workflow
└── checksums.txt
```

개발 repository의 Go source/templates/evals와 사용자에게 배포하는 plugin artifact를 분리한다. bundle은 실행 binary 또는 명확한 installer를 포함해야 하며 network download를 skill 실행 중 암묵적으로 수행하지 않는다.

## 2. Manifest

```json
{
  "name": "tene",
  "version": "0.1.0",
  "description": "Spec-driven sprint workflow and intent-based QA for Codex",
  "author": {"name": "agent-kay-it"},
  "license": "LICENSE",
  "repository": "https://github.com/agent-kay-it/tene-codex"
}
```

실제 작성 시 current official schema와 plugin validator를 기준으로 field를 확정한다. 지원되지 않는 `hooks`, arbitrary commands, binary metadata를 manifest에 임의로 넣지 않는다. runtime 정보는 plugin-owned config/reference에서 관리한다.

## 3. Installation contract

1. Marketplace/plugin install이 bundle을 plugin directory에 배치한다.
2. Codex가 skill을 discover한다.
3. 최초 `$tene-sprint`가 bundled binary platform/sha256을 검사한다.
4. 프로젝트 mutation 전에 `tene-workflow doctor --json` 결과를 보여준다.
5. init은 `.tene-workflow`, docs skeleton, optional AGENTS managed block과 `.gitignore` patch를 각각 명시한다.

기존 파일이 있으면 overwrite하지 않고 diff/merge를 제공한다. uninstall은 plugin bundle만 제거하며 project docs/state는 보존한다. 별도 purge는 core의 명시적 destructive command로도 MVP에는 제공하지 않는다.

## 4. Marketplace release

- tag 기반 immutable artifact, SemVer, checksum, provenance/SBOM.
- release note에 Codex minimum version, core/schema version range, migration, known limitations.
- clean Codex profile에서 install → explicit skill → implicit routing → update → uninstall 검증.
- marketplace 제출 전 plugin validator와 skill validator 결과를 release evidence로 첨부.
- beta channel에서 reference projects와 routing telemetry를 확인한 뒤 stable 승격.

Marketplace의 실제 제출 UI/API와 심사 규칙은 변할 수 있으므로 release 당시 공식 Codex 문서를 다시 확인하는 release checklist 항목으로 둔다. 저장소 문서는 특정 비공식 Claude marketplace 절차를 Codex 절차로 오인하지 않는다.

## 5. Compatibility

```yaml
compatibility:
  plugin: 0.1.x
  core: ">=0.1.0 <0.2.0"
  schema: ">=1.0.0 <2.0.0"
  codex: "probe-required"
  tene_cli: ">=current-tested-minimum"
```

Codex는 capability probe 결과를 우선하며, minimum version은 known-bad를 차단하는 보조 수단이다. `tene` secret CLI는 secret-required workflow에서만 hard dependency다.

## 6. Supply-chain and permissions

- bundled binary checksum을 실행 전에 검증한다.
- plugin scripts는 repo root 밖에 쓰지 않는다.
- remote MCP는 MVP 기본 bundle에 포함하지 않는다.
- child tools는 최소 권한과 project policy allowlist로 실행한다.
- telemetry는 opt-in이며 secret/path/content 원문을 전송하지 않는다.

## 7. Local development verification

```text
go test ./...
tene-workflow doctor --root <fixture> --json
plugin validator <staged-plugin>
skill validator <each-skill>
skill eval runner <routing-corpus>
clean-profile smoke test
```

정확한 validator command는 설치된 Codex tooling을 discovery해 release automation에 고정한다.

