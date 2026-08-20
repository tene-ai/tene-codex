# tene 구현 완성도 리뷰

검토 기준일: 2026-08-20
검토 대상 Sprint: `sprint_0000380ydgqkg6h9mm5ctvm9b4`
현재 workflow phase: `qa`
현재 QA run: `run_0000380z202tx0dwz90g22c0x0`
현재 판정: **조건부 미완료 — 구현과 핵심 happy path는 확인됐지만 전체 blocking AC의 실제 variant 증거가 아직 부족함**

## 1. 리뷰 결론

현재 구현은 plugin identity, bundled CLI, workflow lifecycle, Stop hook, 문서 언어 계약의 핵심 코드와 회귀 테스트를 갖췄다. Fresh Codex session의 설치본 skill discovery와 `$tene:status` invocation, release smoke 7단계, 실제 host turn 정상 종료도 관찰했다.

그러나 현재 `qa evaluate`는 4개 blocking AC를 모두 `AC_UNVERIFIED`로 유지한다. Native test 결과만으로 host/user-visible layer를 통과시키지 않으며, alternate·validation·failure·recovery 중 일부에 실제 variant별 host/runtime observation이 부족하기 때문이다. 따라서 현재 결과물을 “구현 완료” 또는 “출시 준비 완료”로 판정하지 않는다.

## 2. 기획 의도 추적성

| 기획 의도 | Acceptance Criterion | 주요 구현 | 실행 검증 | 현재 판정 |
|---|---|---|---|---|
| 사용자는 제품을 `tene` plugin과 `$tene:<phase>` 이름으로 발견한다 | `ac_0000380ydppfarshvvnhww980w` | `.codex-plugin/plugin.json`, `.agents/plugins/marketplace.json`, `skills/*`, `internal/router/router.go` | Fresh Codex에서 9개 installed skill과 `$tene:status` invocation 관찰; Python/Go regression 실행 | Happy path 입증, 나머지 host variant 보강 필요 |
| 로컬 설치본이 source와 동일하게 실행되고 update/remove가 repository state를 보존한다 | `ac_0000380ydpq9vae1we5cfzkw78` | `scripts/tene-workflow`, `scripts/release-smoke.sh`, `scripts/portable-workflow-smoke.py`, package metadata | Release smoke 7단계 exit 0; tampered checksum, portable workflow, update/remove, state preservation stage 통과 | 핵심 lifecycle 입증, case별 관찰 연결 보강 필요 |
| Active Sprint의 Stop hook이 host lifecycle을 깨지 않는다 | `ac_0000380yeckc1p516k1yvt3q3r` | `hooks/tene_hook.py`, `hooks/hooks.json`, `tests/hooks_test.py` | 실제 ephemeral Codex process가 final response 후 exit 0; invalid Stop output과 unintended continuation 없음 | Happy path 입증, malformed/failure host variant 보강 필요 |
| Workflow 문서는 현재 사용자 대화 언어를 따른다 | `ac_0000380yecmcxpx719rfy2xb7m` | `references/workflow.md`, 9개 `skills/*/SKILL.md`, `tests/language_contract_test.py`, active Sprint 문서 | 독립 evaluator가 PRD·Plan·Design 전 섹션과 machine marker 보존 확인; language tests 통과 | 한국어 happy/recovery 입증, override/failure case 연결 보강 필요 |

## 3. 실제로 구현·수정된 핵심 사항

### Plugin identity와 routing

- Plugin namespace를 `tene`로 통일했다.
- 9개 skill suffix를 `sprint`, `prd`, `plan`, `design`, `loop-check`, `qa`, `report`, `status`, `secrets`로 통일했다.
- `$tene-codex:*`와 `$tene:tene-*`가 active discovery/runtime metadata에 남지 않게 했다.
- Explicit router와 eval oracle이 짧은 canonical suffix를 반환하도록 정렬했다.

### Installed runtime과 lifecycle

- Source, staged bundle, installed launcher의 manifest·inventory·checksum 검증을 구성했다.
- Bundled binary checksum mismatch는 실행 전에 실패하도록 검증한다.
- Release smoke가 package manifest/SBOM, bundled CLI, routing, tampered binary, portable workflow, update/remove, repository-state preservation을 실제 실행한다.
- QA 중 semantic evidence guard가 portable fixture의 복제 assertion을 탐지했고, layer별 actual/expected로 수정한 뒤 release smoke 전체를 재통과했다.

### Workflow와 QA 무결성

- Deprecated/superseded intent가 소유한 criterion을 active context와 QA plan에서 제외한다.
- 반복 실행 evidence가 같은 파일을 덮어쓰지 않도록 artifact identity를 분리했다.
- 과거 overwrite evidence를 검증 가능한 현재 bytes와 조정하는 reconciliation command를 추가했다.
- 외부 observation에서 같은 actual/expected를 여러 layer에 복제하면 `QA_OBSERVATION_DUPLICATE_LAYER_ASSERTION`으로 거부한다.
- 완전히 해당하지 않는 QA case는 독립 evaluator가 승인한 layer별 `not-applicable` disposition으로 처리하며 synthetic evidence를 요구하지 않는다.

### Stop hook과 언어 계약

- Stop event는 event-specific invalid payload 대신 common Stop output을 반환한다.
- Non-shell tool payload를 shell command처럼 검사해 차단하던 pre-tool false positive를 수정했다.
- 한국어 요청에서 PRD, Plan, Design, Loop Check, QA, Report authored prose를 한국어로 유지하는 공통 계약을 추가했다.
- Machine-readable frontmatter, marker, ID, path, command, API와 code symbol은 번역하지 않는다.

## 4. 검증 결과

| 검증 | 결과 | 근거 |
|---|---|---|
| Python regression | PASS, 18 tests | `python3 -m unittest discover -s tests -p '*_test.py'` |
| Go focused regression | PASS | `go test ./internal/qaadapter ./internal/workflow ./internal/app` |
| Release smoke | PASS, 7 stages | `04-qa/release-smoke-observation-v3.txt` |
| Fresh installed selector/invocation | PASS for happy path | `04-qa/fresh-session-observation.txt`, `04-qa/host-session-observation-v2.txt` |
| Actual Stop lifecycle | PASS for active Sprint happy path | `04-qa/host-session-observation-v2.txt`, `04-qa/observations-v2/stop-happy.json` |
| Korean authored documents | PASS | `04-qa/language-observation-v2.txt`, independent evaluator result |
| Evidence byte/hash/redaction integrity | PASS, 349 active records at revision 1651 | `tene-workflow evidence verify --json` |
| Current blocking QA gate | FAIL by design | Four ACs remain `AC_UNVERIFIED` until remaining required host/runtime variants are evidenced |

## 5. 독립 평가에서 발견하고 수정한 false-pass

이전 QA run은 28개 host observation에 실제 variant 실행 없이 결론 문구를 생성하고 같은 내용을 L1–L7에 복제했다. 구조적 validator는 이를 통과시켰지만 독립 evaluator는 4개 AC 모두를 insufficient로 판정했다.

수정 사항:

1. Synthetic observation generator `scripts/qa-host-observe.py`를 제거했다.
2. 서로 다른 layer에 동일 actual/expected를 복제하는 observation을 importer가 거부하도록 했다.
3. Statement label만 `L1`, `L2`로 바꾸는 우회도 거부하는 회귀 테스트를 추가했다.
4. 새 QA run의 run ID, case ID, spec hash, state revision이 일치해야 evidence credit을 받을 수 있게 확인했다.
5. 기존 passed run을 완료 근거로 재사용하지 않고 새 QA run을 생성했다.

이 수정은 현재 QA가 쉽게 통과하지 못하게 만들었지만, 기획 의도와 실제 동작의 일치 여부를 더 신뢰할 수 있게 한다.

## 6. 사용자가 직접 확인할 리뷰 절차

Repository root에서 다음 명령을 실행한다.

```bash
tene-workflow status --json
tene-workflow master status --json
tene-workflow document validate prd --json
tene-workflow document validate plan --json
tene-workflow document validate design --json
tene-workflow evidence verify --json
tene-workflow qa status --json
tene-workflow qa evaluate --json
```

코드 회귀와 release 경계를 확인한다.

```bash
go test ./internal/qaadapter ./internal/workflow ./internal/app
python3 -m unittest discover -s tests -p '*_test.py'
./scripts/release-smoke.sh
```

Fresh Codex session에서는 다음을 직접 확인한다.

1. Skill selector에 `$tene:plan`을 포함한 정확히 9개의 `$tene:<phase>`가 보이는지 확인한다.
2. `$tene-codex:*`와 `$tene:tene-*`가 보이지 않는지 확인한다.
3. `$tene:status`를 호출하고 installed bundled CLI가 active Sprint 상태를 반환하는지 확인한다.
4. 응답 종료 시 Stop hook invalid-output 오류나 자동 continuation이 발생하지 않는지 확인한다.

## 7. 완료 판정 기준

아래 조건을 모두 충족할 때만 최종 완료로 판정한다.

- 모든 blocking AC에 실제 observable evidence가 존재한다.
- 의미 없는 variant/layer만 독립 evaluator 승인 N/A로 분류된다.
- `evidence verify`가 통과한다.
- 새 run의 `qa evaluate`가 통과한다.
- 독립 evaluator가 evidence content를 PASS로 판정한다.
- unresolved blocker가 0이다.
- Report가 PRD → task → design → changed file → QA evidence를 추적한다.
- Strict profile의 `report → archived` human approval을 사용자가 명시적으로 승인한다.

## 8. 현재 남은 작업

- Selector alternate·validation·stale-cache failure·recovery의 실제 host observation 연결
- Installed lifecycle alternate·uninstall·checksum failure·recovery의 case별 lifecycle observation 연결
- Stop empty·malformed·internal-failure variant의 실제 host/hook boundary observation 연결
- Language override·validation·drift detection·recovery의 실제 document fixture observation 연결
- 새 QA evaluate 및 독립 evaluator 재평가
- QA PASS 이후 Sprint report 생성·검증

현재 문서는 중간 상태를 완료로 포장하지 않는 검토 기준서다. 최종 report가 생성되면 이 문서의 미완료 항목과 최종 evidence verdict를 대조해야 한다.
