#!/usr/bin/env python3
import argparse,json,pathlib,sys
p=argparse.ArgumentParser();p.add_argument("--root",default=".");p.add_argument("--final",action="store_true");a=p.parse_args();root=pathlib.Path(a.root).resolve()
manifest=json.loads((root/"docs/release/requirements-traceability.json").read_text());expected={"functional_requirements":[f"FR-{i:02d}" for i in range(1,12)],"product_acceptance":[f"AC-PRODUCT-{i:02d}" for i in range(1,9)],"work_packages":[f"WP-{i:02d}" for i in range(1,15)]};missing=[]
for group,ids in expected.items():
    got=manifest.get(group,{})
    for identifier in ids:
        locators=got.get(identifier,[])
        if not locators:missing.append(identifier+":no-locator")
        for locator in locators:
            if not (root/locator).exists():missing.append(identifier+":"+locator)
state=json.loads((root/".tene-workflow/project.json").read_text());open_gaps=[x["gap_id"] for x in state["gaps"].values() if x["status"] in ("open","deferred")];unfinished=[x["task_id"] for x in state["tasks"].values() if x["status"] not in ("done","deferred")]
bad_archives=[s["sprint_id"] for s in state["sprints"].values() if s["phase"]=="archived" and s.get("last_qa_status")!="passed"]
verified=set()
for run in state["qa_runs"].values():
    if run.get("status")!="passed":continue
    for case in run.get("cases",[]):
        if case.get("status")!="passed":continue
        for eid in case.get("evidence_ids",[]):
            ev=state["evidence"].get(eid,{});verified.update(ac for ac in case.get("ac_ids",[]) if ev.get("redaction_status")=="passed" and ac in ev.get("ac_ids",[]))
blocking={ac["ac_id"] for ac in state["acceptance_criteria"].values() if ac.get("priority")=="blocking"};unverified=sorted(blocking-verified)
result={"schema_version":manifest["schema_version"],"coverage":{"functional_requirements":len(expected["functional_requirements"]),"product_acceptance":len(expected["product_acceptance"]),"work_packages":len(expected["work_packages"])},"missing":missing,"open_gaps":open_gaps,"unfinished_tasks":unfinished,"unverified_blocking_criteria":unverified,"invalid_archived_sprints":bad_archives,"active_sprint_id":state.get("active_sprint_id","")}
result["passed"]=not missing and not open_gaps and not unfinished and not unverified and not bad_archives and (not a.final or not result["active_sprint_id"]);print(json.dumps(result,indent=2));sys.exit(0 if result["passed"] else 1)
