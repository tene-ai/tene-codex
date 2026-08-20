---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380ydgqkg6h9mm5ctvm9b4
phase: design
status: ready
revision: 1103
intent_ids: [intent_0000380ydppfay2szegwvaqsqm, intent_0000380ydpq9tm54tgwpqnfd30, intent_0000380yeckc0mm5dsr9j59h6m, intent_0000380yecmcwb7v4ttjm9bbrm]
generated_at: 2026-08-20T09:06:41Z
generated_by: tene-workflow
---

# design — Plugin Identity and Installed Runtime Verification

<!-- tene:section:purpose -->
## 목적

`$tene-codex:tene-*`에서 `$tene:*`로의 atomic identity migration과 source→stage→installed 검증 경계를 정의한다. 또한 confirmed intent가 소유한 criterion만 context와 QA에 들어가게 하여 deprecated record는 audit용으로 보존하되 active gate가 되지 않게 한다.

<!-- tene:section:scope -->
## 범위

영향 범위는 plugin/marketplace JSON, 9개 skill directory와 metadata, router/eval/active 문서, package/release smoke, context·QA criterion 선택, focused test와 local-install 관찰이다. repository/module URL, archived historical document, remote service와 secret 값은 영향받지 않는다.

<!-- tene:section:layers -->
## Understanding Layers

| Layer | Components | Contract | Failure boundary |
|---|---|---|---|
| Interface | manifest, skill frontmatter, agent prompts, selector | plugin=`tene`; suffix in fixed nine-name set | partial identity produces forbidden mixed selector |
| Business | router and eligible-criterion selector | explicit `$tene:<suffix>` maps to the same suffix; only confirmed owner intents are active | stale aliases or deprecated AC cases fail tests/loop |
| Persistence | source tree, staged bundle, installed cache, workflow repo | identity/inventory equality; plugin lifecycle never deletes repo state | drift or state hash change blocks AC |
| Infrastructure | package builder, launcher, checksum, Codex host | verified executable selected per platform; host observation required | missing host capability stays unknown/blocking |

<!-- tene:section:six-questions -->
## Six Questions

| Name | Defined at | Imported/referenced at | Called/used at | Input shape | Output/mutation |
|---|---|---|---|---|---|
| Manifest identity | `.codex-plugin/plugin.json` | marketplace/package/host | Codex discovery | JSON metadata | namespace `tene` |
| Skill contract | `skills/<suffix>/SKILL.md` | agent YAML, package, host | skill invocation | user request/repo state | phase orchestration |
| `router.Route` | `internal/router/router.go` | app route command and evals | `tene-workflow route` | `{text,phase,active}` | `{skill,confidence,reason}` |
| `eligibleCriteria` (target helper) | `internal/app/app.go` or workflow package | context builder, QA planner and spec hash/gate selectors | context/qa commands | project+sprint | confirmed-owner criterion slice; no mutation |
| package builder | `scripts/package-plugin.sh` | release smoke | release/local QA | version/stage/platform | staged immutable bundle |
| launcher | `scripts/tene-workflow` | all installed skills | host child process | argv, platform, checksum | exec core; exit code |
| host selector | provider external/unknown | Codex UI | search/select | installed metadata | observable qualified labels |

<!-- tene:section:traceability -->
## 추적성

- AC `ac_0000380ydppfarshvvnhww980w`: manifest/skill/router contracts and selector QA; tasks T1/T2.
- AC `ac_0000380ydpq9vae1we5cfzkw78`: package/launcher/state contracts and confirmed-only QA selection; tasks T3/T4.
- AC `ac_0000380yeckc1p516k1yvt3q3r`: Codex Stop-event common JSON output와 workflow 무변경성; task T5.
- AC `ac_0000380yecmcxpx719rfy2xb7m`: 사용자 대화 언어를 따르는 authored 문서 계약; task T6.
- T1 `task_0000380ydvzc0pjsatbyeff3yc`, T2 `task_0000380ydwcs0gd3wr3vy2wcdr`, T3 `task_0000380ydwdvfe6vc2ajcpd7g8`, T4 `task_0000380ydwexgf4pewweg3pj8g`, T5 `task_0000380yed20cf2ayn6vfmkn9w`, T6 `task_0000380yed2vtrspn07jwg36sr`.
- Design source: `docs/02-design/04-skill-contracts-and-routing.md`, `05-hooks-codex-integration.md`, `10-plugin-package-and-marketplace.md`, `11-testing-evals-and-acceptance.md` plus the confirmed user naming policy.

<!-- tene:contract path=".codex-plugin/plugin.json" symbol="\"name\": \"tene\"" -->
<!-- tene:contract path="internal/router/router.go" symbol="tene:" -->
<!-- tene:forbid path=".codex-plugin/plugin.json" contains="\"name\": \"tene-codex\"" -->

<!-- tene:section:decisions -->
## 결정 사항

- 짧은 suffix는 표시용 alias에 그치지 않고 host metadata와 router 전반의 canonical 값이다.
- `eligibleCriteria(project,sprint)`는 공유 pure predicate다. Criterion priority는 caller가 평가하지만 owner intent는 반드시 존재하고 해당 Sprint에 속하며 status가 `confirmed`여야 한다.
- Deprecated/superseded intent와 criterion record는 journal에 보존하며 destructive migration을 사용하지 않는다.
- Package version은 한 번 올리고 지원되는 local marketplace/plugin lifecycle로 cache를 교체한다. Cache file을 직접 편집하지 않는다.
- 이전 invocation 문자열은 migration note와 archived history에만 허용하며 active prompt, test 또는 discoverable metadata에는 허용하지 않는다.

<!-- tene:section:freeform -->
## 추가 관점

이 저장소에는 `.codegraph/`가 없으므로 현재 static provider는 Go symbol discovery와 targeted reference search다. Codex selector 조합은 외부 dynamic behavior이므로 source assertion이 fresh-session 관찰을 대신할 수 없다.

<!-- tene:section:components -->
## 컴포넌트

### Identity matrix

9개 suffix `sprint`, `prd`, `plan`, `design`, `loop-check`, `qa`, `report`, `status`, `secrets`의 검증 table을 정의한다. 각 row는 directory basename, SKILL frontmatter, agent default prompt, 예상 qualified selector를 연결한다. Package smoke는 source와 stage에서 같은 table을 검증한다.

### Router migration

Explicit parser는 `$tene:<suffix>`를 받아 `<suffix>`를 반환한다. 자연어 scoring은 phase semantics를 유지하지만 모든 cue map과 oracle은 canonical suffix를 반환한다. 기존 qualified form은 hard negative다. Parsing은 deterministic하며 workflow state를 전이하지 않는다.

### Criterion eligibility

Pure helper는 Sprint의 intent ID와 현재 intent status에서 active criterion을 도출한다. Context build, QA planning, QA spec hashing과 gate coverage가 이 helper를 사용하므로 deprecated/superseded criterion은 mandatory context나 test case로 나타날 수 없다. Owner intent 누락은 active evidence credit이 아니라 data-integrity finding이다.

### Package/install verifier

명시적인 임시 stage에 build하고 JSON/inventory/checksum을 검증한 뒤 지원되는 Codex/local marketplace command로 설치한다. sanitized metadata, exit code, hash만 기록하며 설치 전과 update/uninstall 후 workflow health를 비교한다.

<!-- tene:section:interfaces -->
## 인터페이스

```text
route(text="$tene:plan", phase=?, active=?)
  -> { selected_skill: "plan", explicit: true, confidence: 1.0 }

eligibleCriteria(project, sprint)
  -> criteria where criterion.intent_id belongs to sprint
     and project.intents[criterion.intent_id].status == "confirmed"

verifyIdentity(root)
  -> rows[{suffix, directory, frontmatter, prompt, qualified_name}], findings[]
```

모든 실패는 기존 JSON envelope과 stable exit code를 사용한다. Identity validation은 read-only이며 package output은 명시적인 temporary/stage directory에 쓴다. 설치 lifecycle은 repository workflow state를 변경하지 않는다.

<!-- tene:section:data -->
## 데이터

- Plugin identity: top-level manifest `name="tene"`, display name `tene`, marketplace plugin entry의 name/display가 일치한다.
- Skill identity row: suffix, directory, frontmatter name, default prompt, implicit policy를 포함한다.
- Active criterion reference: 기존 immutable ID를 사용하며 schema migration은 없다.
- Install observation: package version, platform, inventory/hash 요약, selector label, command exit/health와 repository-state 전후 hash를 포함한다. Raw environment나 민감정보는 포함하지 않는다.
- Migration: 이전 installed ID는 지원되는 host operation으로 remove/disable하며 새 ID/version은 별도 cache entry를 만든다. Source repository 이름은 그대로 유지한다.

<!-- tene:section:state-transitions -->
## 상태 전이

1. Source old identity → atomic source migration → static validation.
2. Source valid → staged bundle built → source/stage equality.
3. Old local install (optional) → supported update/remove → new `tene` install.
4. Host restart/new session → selector discovery → invocation.
5. Installed plugin invoke → launcher checksum → core status/doctor.
6. Uninstall → plugin bundle absent while repository workflow state remains healthy.

Validation/build retry는 idempotent하다. Install/update action은 version-scoped이며 넓은 directory를 recursive delete하면 안 된다. 동시 workflow mutation은 lifecycle hash 비교에서 제외하거나 revision mismatch로 탐지한다.

<!-- tene:section:failures -->
## 실패와 복구

- Mixed namespace: 정확한 metadata row를 지목하는 static blocker로 기록한다.
- 오래된 cache/중복 selector: installation blocker로 기록하고 이전 plugin을 remove/disable한 뒤 restart한다.
- Checksum mismatch 또는 binary 누락: exec 전에 dependency/security error로 실패한다.
- Host selector를 사용할 수 없음: capability를 `unknown`으로 두고 AC는 미입증 상태로 유지한다.
- Deprecated criterion이 context/QA에 나타남: regression failure와 loop gap으로 기록한다.
- Source/stage/install 불일치: installed-source drift blocker로 기록한다.
- Update/uninstall 중 repository state 변경: persistence blocker로 기록하고 재시도 전에 revision을 비교하여 외부 concurrent mutation을 진단한다.

<!-- tene:section:security -->
## 보안

- `.tene/**`를 읽거나 environment를 dump하거나 credential을 evidence에 저장하지 않는다.
- Plugin cache는 metadata/inventory 검증을 위해서만 읽고, 설치 변경에는 지원되는 host command와 정확한 plugin ID를 사용한다.
- Temp/stage path는 명시적이고 범위가 제한되어야 하며 넓은 recursive deletion을 금지한다.
- Bundled execution 전에 checksum을 검증한다.
- 추후 민감정보 의존 QA가 필요하면 `$tene:secrets`를 사용한다. 이 Sprint에는 비밀값이 필요하지 않다.
- Security와 evidence-integrity finding은 waiver할 수 없다.

<!-- tene:section:tests -->
## 테스트

### Stop hook output adapter

`Stop`은 active Sprint가 있을 때 common output `{continue: true, systemMessage: ...}`만 반환한다. `hookSpecificOutput.additionalContext`는 SessionStart 등 허용 event에만 사용한다. Stop은 advisory이므로 `decision: block`을 반환하지 않고, workflow mutation도 수행하지 않는다.

### 문서 언어 계약

각 phase skill은 공통 workflow reference에서 `document_language` 의미를 상속한다. 현재 사용자 요청의 dominant language를 기본으로 사용하고 명시적 언어 요청이 우선한다. authored prose만 해당하며 frontmatter, marker, ID, locator, command와 code symbol은 그대로 둔다. phase handoff와 compaction 후에도 같은 규칙을 다시 읽는다.

- Unit: router의 9개 explicit case, old/mixed negative, implicit phase case와 eligible criterion의 confirmed/candidate/deprecated/superseded/missing-owner table을 검증한다.
- Contract: source와 staged bundle의 identity matrix, manifest/marketplace 일치, directory/frontmatter 일치, default prompt의 예상 qualified name 포함 여부를 검증한다.
- Package: 지원 platform inventory, checksum 정상·변조, launcher source/bundled resolution을 검증한다.
- Integration: deprecated criterion이 있는 fixture의 QA plan/context에 confirmed AC만 포함되는지 검증한다.
- Host E2E: local clean install, restart/new conversation, 9개 selector, 대표 skill invocation, update/uninstall state 보존을 검증한다.
- Loop: 위 executable contract와 모든 changed file의 task 연결을 검증하고, 독립 evaluator가 이전 runtime identity 및 근거 없는 host claim의 부재를 확인한다.
