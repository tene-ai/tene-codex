#!/usr/bin/env python3
"""Fail-closed semantic and executable audit for documented MVP contracts."""
import argparse, hashlib, json, pathlib, re, subprocess, sys

EXPECTED={"functional_requirements":[f"FR-{i:02d}" for i in range(1,12)],"product_acceptance":[f"AC-PRODUCT-{i:02d}" for i in range(1,9)],"work_packages":[f"WP-{i:02d}" for i in range(1,15)]}
COMMANDS={
 "app-master":[["go","test","./internal/app","-run","TestMaster"]], "app-intent":[["go","test","./internal/app","-run","TestMasterMetadataAndStructuredIntentLifecycle"]], "app-report":[["go","test","./internal/app","-run","TestCLIHappyPath"]],
 "workflow":[["go","test","./internal/workflow"]], "document":[["go","test","./internal/document"]], "loopcheck":[["go","test","./internal/loopcheck"]], "qa":[["go","test","./internal/qaadapter","./internal/workflow","./internal/app","-run","QA|Observation|Evidence"]],
 "state":[["go","test","./internal/state"]], "context":[["go","test","./internal/tracecontext","./internal/app","-run","Context|Graph"]], "graph":[["go","test","./internal/tracecontext","-run","Graph|Impact"]], "codeintel":[["go","test","./internal/codeintel"]],
 "routing":[["go","run","./cmd/tene-routing-eval","evals/routing-corpus.json"]], "secret":[["go","test","./internal/secret"]], "hooks":[["python3","-m","unittest","tests/hooks_test.py"]], "references":[["npm","run","test:references"]],
 "playwright":[["npx","playwright","test","tests/e2e/reference-web.spec.ts"]], "release":[["./scripts/release-smoke.sh"]], "app-compact":[["go","test","./internal/app","-run","TestCLICompactBoundsJournalAndDoctorDetectsArchiveTamper"]], "schema":[["@json-schemas"]], "plugin":[["@plugin"]]}

def audit_contract(root,identifier,contract):
    failures=[]; source=root/contract.get("source","")
    if not source.is_file(): failures.append("source-missing:"+contract.get("source",""))
    elif identifier not in source.read_text(errors="replace"): failures.append("source-id-missing:"+contract["source"])
    if not contract.get("symbols"): failures.append("symbols-empty")
    for locator in contract.get("symbols",[]):
        if not isinstance(locator,list) or len(locator)!=2: failures.append("symbol-invalid:"+repr(locator)); continue
        path,pattern=locator; target=root/path
        if not target.is_file(): failures.append("symbol-file-missing:"+path); continue
        try: matched=re.search(pattern,target.read_text(errors="replace"),re.MULTILINE) is not None
        except re.error as exc: failures.append("symbol-regex-invalid:"+str(exc)); continue
        if not matched: failures.append("symbol-missing:"+path+":"+pattern)
    if not contract.get("commands"): failures.append("commands-empty")
    failures += ["command-unknown:"+name for name in contract.get("commands",[]) if name not in COMMANDS]
    return failures

def run_command(root,name,timeout):
    logs=[]
    for argv in COMMANDS[name]:
        if argv[0]=="@json-schemas":
            try:
                for path in sorted((root/"schemas").glob("*.json")): json.loads(path.read_text())
                logs.append("validated schemas/*.json")
            except Exception as exc: return False,str(exc)
            continue
        if argv[0]=="@plugin":
            try:
                for path in (root/".codex-plugin/plugin.json",root/".agents/plugins/marketplace.json",root/"hooks/hooks.json"): json.loads(path.read_text())
                logs.append("validated plugin manifests")
            except Exception as exc: return False,str(exc)
            continue
        try: result=subprocess.run(argv,cwd=root,capture_output=True,text=True,timeout=timeout)
        except (OSError,subprocess.TimeoutExpired) as exc: return False,str(exc)
        logs.append("$ "+" ".join(argv)+"\n"+(result.stdout+result.stderr)[-2000:])
        if result.returncode: return False,"\n".join(logs)
    return True,"\n".join(logs)

def state_findings(root,final):
    state=json.loads((root/".tene-workflow/project.json").read_text()); failures=[]
    failures += ["gap:"+x["gap_id"] for x in state["gaps"].values() if x["status"] in ("open","deferred")]
    failures += ["task:"+x["task_id"] for x in state["tasks"].values() if x["status"] not in ("done","deferred")]
    failures += ["archive:"+s["sprint_id"] for s in state["sprints"].values() if s["phase"]=="archived" and s.get("last_qa_status")!="passed"]
    verified=set()
    archived_passed={s["sprint_id"] for s in state["sprints"].values() if s["phase"]=="archived" and s.get("last_qa_status")=="passed"}
    for ev in state["evidence"].values():
        path=root/ev.get("uri",""); valid=path.is_file() and hashlib.sha256(path.read_bytes()).hexdigest()==ev.get("sha256")
        if valid and ev.get("redaction_status")=="passed" and ev.get("sprint_id") in archived_passed:
            verified.update(ev.get("ac_ids",[]))
    for run in state["qa_runs"].values():
        if run.get("status")!="passed": continue
        for case in run.get("cases",[]):
            if case.get("status")!="passed": continue
            for eid in case.get("evidence_ids",[]):
                ev=state["evidence"].get(eid,{}); path=root/ev.get("uri","")
                valid=path.is_file() and hashlib.sha256(path.read_bytes()).hexdigest()==ev.get("sha256")
                if valid and ev.get("redaction_status")=="passed" and ev.get("run_id")==run.get("run_id") and ev.get("case_id")==case.get("case_id"):
                    verified.update(ac for ac in case.get("ac_ids",[]) if ac in ev.get("ac_ids",[]))
    blocking={ac["ac_id"] for ac in state["acceptance_criteria"].values() if ac.get("priority")=="blocking"}
    failures += ["unverified:"+identifier for identifier in sorted(blocking-verified)]
    if final and state.get("active_sprint_id"): failures.append("active-sprint:"+state["active_sprint_id"])
    return failures,state.get("active_sprint_id","")

def main():
    p=argparse.ArgumentParser(); p.add_argument("--root",default="."); p.add_argument("--manifest",default="docs/release/semantic-contracts.json"); p.add_argument("--final",action="store_true"); p.add_argument("--no-exec",action="store_true"); p.add_argument("--contracts-only",action="store_true"); p.add_argument("--timeout",type=int,default=180); a=p.parse_args()
    root=pathlib.Path(a.root).resolve(); manifest_path=pathlib.Path(a.manifest); manifest_path=manifest_path if manifest_path.is_absolute() else root/manifest_path; manifest=json.loads(manifest_path.read_text()); contract_failures={}; used=set()
    for group,ids in EXPECTED.items():
        actual=manifest.get(group,{})
        for identifier in ids:
            if identifier not in actual: contract_failures[identifier]=["contract-missing"]; continue
            failures=audit_contract(root,identifier,actual[identifier]); used.update(actual[identifier].get("commands",[]))
            if failures: contract_failures[identifier]=failures
        for extra in sorted(set(actual)-set(ids)): contract_failures[extra]=["unexpected-contract"]
    command_results={}
    if not a.no_exec:
        for name in sorted(used):
            ok,log=run_command(root,name,a.timeout); command_results[name]={"passed":ok,"log_tail":log[-2000:]}
            if not ok: contract_failures.setdefault("command:"+name,[]).append("execution-failed")
    workflow_failures,active=([],"") if a.contracts_only else state_findings(root,a.final)
    result={"schema_version":manifest.get("schema_version"),"coverage":{group:len(ids) for group,ids in EXPECTED.items()},"missing":sorted(f"{i}:{f}" for i,fs in contract_failures.items() for f in fs),"contract_failures":contract_failures,"command_results":command_results,"workflow_failures":workflow_failures,"active_sprint_id":active}
    result["passed"]=not contract_failures and not workflow_failures; print(json.dumps(result,indent=2,ensure_ascii=False)); return 0 if result["passed"] else 1
if __name__=="__main__": sys.exit(main())
