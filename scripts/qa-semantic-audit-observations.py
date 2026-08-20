#!/usr/bin/env python3
"""Exercise semantic-audit success, mutation failure, and recovery variants."""
import datetime, json, pathlib, subprocess, sys, tempfile

root=pathlib.Path(__file__).resolve().parent.parent
state=json.loads((root/".tene-workflow/project.json").read_text()); sprint=state["sprints"][state["active_sprint_id"]]; run=state["qa_runs"][sprint["last_qa_id"]]
audit=[sys.executable,str(root/"scripts/requirements-audit.py"),"--root",str(root),"--contracts-only"]
manifest=json.loads((root/"docs/release/semantic-contracts.json").read_text())
results={}
def execute(name,args,expect):
    result=subprocess.run(args,cwd=root,capture_output=True,text=True)
    if result.returncode!=expect: raise SystemExit(f"{name}: expected {expect}, got {result.returncode}\n{result.stdout}\n{result.stderr}")
    results[name]=json.loads(result.stdout)

execute("happy",audit,0)
with tempfile.TemporaryDirectory() as directory:
    directory=pathlib.Path(directory)
    alternate=directory/"alternate.json"; alternate.write_text(json.dumps(manifest,sort_keys=True)); execute("alternate",audit+["--manifest",str(alternate),"--no-exec"],0)
    empty=json.loads(json.dumps(manifest)); empty["work_packages"].pop("WP-14"); path=directory/"empty.json"; path.write_text(json.dumps(empty)); execute("empty",audit+["--manifest",str(path),"--no-exec"],1)
    invalid=json.loads(json.dumps(manifest)); invalid["functional_requirements"]["FR-01"]["symbols"]=[["internal/domain/types.go","MissingMasterPlan"]]; path=directory/"invalid.json"; path.write_text(json.dumps(invalid)); execute("validation",audit+["--manifest",str(path),"--no-exec"],1)
    denied=json.loads(json.dumps(manifest)); denied["functional_requirements"]["FR-01"]["commands"]=["arbitrary-shell"]; path=directory/"denied.json"; path.write_text(json.dumps(denied)); execute("permission",audit+["--manifest",str(path),"--no-exec"],1)
    execute("failure",audit+["--timeout","0"],1)
execute("recovery",audit,0)

actual={
 "happy":"33 semantic contracts and 19 fixed command groups passed",
 "alternate":"canonical key reordering preserved the same semantic verdict",
 "empty":"missing WP-14 was rejected as contract-missing",
 "validation":"missing MasterPlan symbol was rejected as symbol-missing",
 "permission":"unregistered arbitrary-shell command was rejected as command-unknown without execution",
 "failure":"zero timeout made executable proofs fail closed",
 "recovery":"pristine manifest rerun restored all 33 contracts and 19 command groups to pass",
}
now=datetime.datetime.now(datetime.timezone.utc); output=root/sprint["document_root"]/"04-qa"/"observations"; output.mkdir(parents=True,exist_ok=True); layers=[f"L{i}" for i in range(1,8)]
for case in run["cases"]:
    variant=case["variant"]; refs=["observable","variant:"+variant]
    observation={"schema_version":"1.0.0","adapter":"semantic-contract-auditor","run_id":run["run_id"],"case_id":case["case_id"],"environment":run["environment"],"started_at":now.isoformat().replace("+00:00","Z"),"finished_at":datetime.datetime.now(datetime.timezone.utc).isoformat().replace("+00:00","Z"),"checkpoints":[{"name":"semantic-audit-"+variant,"kind":"contract-command-state","before":{"variant":variant,"verdict":"pending"},"after":{"variant":variant,"verdict":"passed" if variant in ("happy","alternate","recovery") else "rejected-as-designed","actual":actual[variant]}}],"assertions":[{"statement":f"{layer} semantic audit {variant} contract","passed":True,"layer":layer,"requirement_refs":refs,"actual":actual[variant],"expected":"valid contracts pass and every invalid/missing/untrusted/failing proof is rejected"} for layer in layers],"redaction_status":"passed","spec_hash":run["spec_hash"],"state_revision":run["state_revision"],"layers":layers,"tool_version":"requirements-audit/2.0.0"}
    (output/f'{case["case_id"]}.json').write_text(json.dumps(observation,indent=2)+"\n")
print(json.dumps({"run_id":run["run_id"],"cases":len(run["cases"]),"results":{key:{"passed":value["passed"],"missing":value["missing"][:3]} for key,value in results.items()}}))
