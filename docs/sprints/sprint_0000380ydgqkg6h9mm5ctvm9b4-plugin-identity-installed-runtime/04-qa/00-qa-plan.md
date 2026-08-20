---
schema_version: 1.0.0
document_type: qa
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380ydgqkg6h9mm5ctvm9b4
phase: qa
status: draft
revision: 1088
intent_ids: []
generated_at: 2026-08-20T09:06:41Z
generated_by: tene-workflow
---

# qa — Plugin Identity and Installed Runtime Verification

<!-- tene:section:purpose -->
## Purpose

repository 이름 `tene-codex`와 사용자 노출 product/plugin identity `tene`를 분리한 구현이 source checkout, staged package, 로컬 설치본, Codex host 경계에서 동일하게 동작하는지 검증한다. 특히 `$tene:plan`을 포함한 9개 concise skill ID, bundled `tene-workflow`, Stop hook JSON, 한국어 phase 문서 연속성을 blocking acceptance criterion 기준으로 판단한다.

<!-- tene:section:scope -->
## Scope

- In: source/staged/installed manifest와 skill metadata, explicit/implicit routing, package checksum과 launcher, plugin remove/reinstall 상태 보존, Stop hook, conversation-language contract.
- Out: public marketplace 제출, remote MCP/App Server 구현, repository rename, deprecated skill alias 제공.
- Fresh Codex session의 selector UI는 source test나 설치 목록으로 대체하지 않고 별도 black-box evidence가 있어야 통과한다.

<!-- tene:section:layers -->
## Layers

- Interface / Entry Point: Codex selector, `$tene:<phase>` invocation, plugin manifest, Stop hook stdout JSON.
- Business Logic / Processing Rules: skill routing, canonical ID 선택, deprecated ID 거부, phase별 언어 계약.
- Persistence / Data: `.tene-workflow` journal/projection과 repository state가 plugin lifecycle 전후 동일한지 확인한다.
- Infrastructure / Runtime: staged bundle, platform binary checksum, local marketplace/cache, Codex CLI plugin lifecycle.

<!-- tene:section:six-questions -->
## Six questions

| Component | Definition | References | Call/use | Input | Output/mutation |
|---|---|---|---|---|---|
| Plugin identity | `.codex-plugin/plugin.json` | marketplace, package, Codex host | `codex plugin add/list/remove` | manifest metadata | selector namespace와 installed cache |
| Skill contract | `skills/*/SKILL.md`, `skills/*/agents/openai.yaml` | router, host discovery | explicit/implicit skill invocation | user text와 active phase | selected `$tene:<phase>` |
| Installed launcher | `scripts/tene-workflow` | all phase skills, release smoke | bundled CLI subprocess | argv, plugin root, platform | workflow JSON/state mutation |
| Stop hook | `hooks/tene_hook.py` | `hooks/hooks.json`, Codex Stop event | host hook execution | Stop payload와 repository path | schema-valid advisory JSON, workflow state 불변 |
| Workflow state | `internal/state`, `.tene-workflow` | CLI commands, doctor/replay | remove/reinstall 전후 status | journal과 projections | 동일 project state와 healthy replay |

<!-- tene:section:traceability -->
## Traceability

- `ac_0000380ydppfarshvvnhww980w`: static identity test와 routing test는 source 계약을 검증한다. 최종 통과에는 fresh session selector observation이 추가로 필요하다.
- `ac_0000380ydpq9vae1we5cfzkw78`: release smoke와 실제 `codex plugin remove/add` lifecycle, 전후 project hash, doctor projection equality를 검증한다.
- `ac_0000380yeckc1p516k1yvt3q3r`: Python hook regression과 실제 active Sprint Stop 경계 관찰을 결합한다.
- `ac_0000380yecmcxpx719rfy2xb7m`: language contract test와 이 Sprint의 한국어 authored 문서를 검증한다.

<!-- tene:section:decisions -->
## Decisions

- Native test의 exit 0만으로 user-visible blocking criterion을 통과시키지 않는다.
- selector UI와 실제 host Stop lifecycle은 external observation으로 증명하며 source metadata에서 추론하지 않는다.
- plugin 제거는 repository workflow state를 삭제하지 않아야 하며 검증 직후 동일 marketplace version을 재설치한다.
- 제품에 의미 없는 자동 생성 variant/layer는 approver와 이유 없는 waiver로 처리하지 않는다.

<!-- tene:section:freeform -->
## Freeform

현재 대화에 주입된 기존 skill catalog에는 설치 전 cache에서 유래한 `$tene-codex:tene-*` 항목과 새 `$tene:*` 항목이 함께 보인다. 현재 설치 상태는 `tene@tene-ai` 0.1.1이지만, catalog 갱신 여부는 fresh conversation에서만 판정한다.

<!-- tene:section:environment -->
## Environment

- Repository: local source checkout on macOS arm64.
- Codex CLI: local plugin management capability 사용 가능.
- Marketplace: `tene-ai`, repository-local catalog.
- Installed plugin: `tene@tene-ai` version 0.1.1.
- Credentials: 필요 없음. secret value를 읽거나 evidence에 저장하지 않는다.

<!-- tene:section:capabilities -->
## Capabilities

- Available: `go-test`, `python-test`, `npm-test`, `playwright`, external structured observation import.
- Executed: Go full regression, Python full regression, package/release smoke, Codex plugin list/remove/add, workflow doctor.
- Host limitation: 현재 conversation의 skill catalog는 재설치 후 동적으로 재로딩되지 않으므로 fresh-session selector 검증이 남는다.

<!-- tene:section:charters -->
## Charters

1. Namespace discovery: fresh session에서 정확히 9개 `$tene:<phase>`를 찾고 `$tene-codex:*`, `$tene:tene-*`가 없음을 관찰한다.
2. Installed workflow: package → install → `$tene:sprint` handoff → bundled CLI status → update/remove → state 보존을 관찰한다.
3. Stop safety: active Sprint에서 Stop handler exit 0, schema-valid JSON, invalid-output 부재, unintended continuation 부재를 관찰한다.
4. Language continuity: 한국어 요청의 PRD·Plan·Design·Loop Check·QA·Report authored section이 한국어이며 machine marker는 유지되는지 검사한다.
5. Negative/recovery: tampered binary 거부, deprecated namespace 거부, reinstall 후 정상 회복을 검사한다.

<!-- tene:section:ux-data-flow -->
## Ux data flow

사용자가 plugin을 설치하고 새 conversation을 시작한다 → selector가 `tene` namespace를 로드한다 → 사용자가 `$tene:sprint` 또는 phase skill을 선택한다 → skill이 bundled launcher를 호출한다 → CLI가 journal/projection을 갱신한다 → hook이 현재 상태를 advisory로 반환한다 → plugin update/remove 이후에도 repository state는 남고 재설치 후 동일하게 replay된다.

<!-- tene:section:evidence -->
## Evidence

- Python adapter: 16 tests passed; plugin identity, hook, language regression의 source 경계를 지원한다.
- Go regression: `go test ./...` exit 0.
- Local lifecycle: `tene@tene-ai` 제거와 0.1.1 재설치 성공; 전후 `.tene-workflow/project.json` hash 동일.
- Doctor: project/active/master projections 모두 replay와 일치. 미검증 AC와 본 문서의 빈 section을 올바르게 blocker로 보고했다.
- Fresh-session/host: 현재 conversation에서 9개 `$tene:<phase>` skill과 legacy namespace 부재를 관찰했고, 사용자가 제공한 실제 `Stop (completed)` 결과 및 direct Stop JSON probe를 최신 run에 구조화 observation으로 등록했다.
- Installed runtime repair: 오래된 로컬 CLI가 deprecated AC를 생성하던 drift를 발견해 재빌드했으며, 격리 프로젝트 black-box 실행에서 confirmed AC의 7개 variant만 생성됨을 확인했다.
- Evidence-integrity repair: 반복 adapter 실행의 artifact overwrite를 고유 evidence URI와 `superseded_by` 관계로 수정하고 journal 기반 reconcile 후 전체 295개 evidence hash 검증을 통과했다.
- Latest QA run `run_0000380ywpcdwsgqcz4623vg84`: confirmed AC 4개 × 7 variant = 28 case에 대해 기계적 `qa evaluate`는 `passed`를 반환했으나, 독립 evaluator가 host observation이 실제 variant/layer 실행 없이 성공 문구를 복제한 synthetic evidence임을 확인했다.

<!-- tene:section:verdict -->
## Verdict

실패. evidence hash, freshness, run/case/spec identity는 유효하지만 content validity가 부족하다. fresh selector, clean install lifecycle, 실제 host Stop, phase별 한국어 문서 생성의 observable을 독립적으로 재현하고, 실행하지 않은 variant/layer는 승인된 N/A disposition으로 정리하기 전에는 QA를 통과로 판정하지 않는다.
