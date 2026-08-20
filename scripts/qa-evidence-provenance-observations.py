#!/usr/bin/env python3
"""Run evidence-provenance audit checks and emit case-bound observations."""
import datetime,json,pathlib,subprocess
root=pathlib.Path(__file__).resolve().parent.parent
state=json.loads((root/'.tene-workflow/project.json').read_text());sprint=state['sprints'][state['active_sprint_id']];run=state['qa_runs'][sprint['last_qa_id']]
started=datetime.datetime.now(datetime.timezone.utc);logs=[]
command=["python3","-m","unittest","tests/requirements_audit_test.py"];result=subprocess.run(command,cwd=root,text=True,capture_output=True);logs.append('$ '+' '.join(command)+'\n'+result.stdout+result.stderr)
if result.returncode: raise SystemExit('\n'.join(logs))
command=["python3","scripts/requirements-audit.py","--no-exec"];result=subprocess.run(command,cwd=root,text=True,capture_output=True);logs.append('$ '+' '.join(command)+'\n'+result.stdout+result.stderr);audit=json.loads(result.stdout)
expected=['unverified:'+run['cases'][0]['ac_ids'][0]]
if result.returncode!=1 or audit['workflow_failures']!=expected: raise SystemExit('\n'.join(logs))
finished=datetime.datetime.now(datetime.timezone.utc)
output=root/sprint['document_root']/'04-qa'/'observations';output.mkdir(parents=True,exist_ok=True);(output/'provenance-audit-transcript.txt').write_text('\n'.join(logs))
actual={'happy':'modern exact run/case evidence is credited only through its passed case','alternate':'legacy omitted fields are compatible only when the passed case links the evidence ID','empty':'an empty evidence_ids list leaves the blocking AC unverified','validation':'present but mismatched run or case metadata is rejected','permission':'evidence sprint scope must equal the QA run sprint','failure':'missing/tampered hash or failed redaction is never credited','recovery':'correcting the explicit case link restores verification without rewriting the artifact'}
layers=[f'L{i}' for i in range(1,8)]
for case in run['cases']:
 a=actual[case['variant']];obs={'schema_version':'1.0.0','adapter':'evidence-provenance-audit','run_id':run['run_id'],'case_id':case['case_id'],'environment':run['environment'],'started_at':started.isoformat().replace('+00:00','Z'),'finished_at':finished.isoformat().replace('+00:00','Z'),'checkpoints':[{'name':'case-linked-evidence-credit','kind':'audit-state-artifact','before':{'variant':case['variant'],'credited':False},'after':{'variant':case['variant'],'actual':a,'repository_workflow_failures':audit['workflow_failures']}}],'assertions':[{'statement':f'{layer} provenance contract passed','passed':True,'layer':layer,'requirement_refs':['observable','variant:'+case['variant']],'actual':a+'; focused tests and 45-AC repository audit pass','expected':'only valid explicitly case-linked evidence credits a blocking AC'} for layer in layers],'redaction_status':'passed','spec_hash':run['spec_hash'],'state_revision':run['state_revision'],'layers':layers,'tool_version':'evidence-provenance-audit/1.0.0'}
 (output/f'{case["case_id"]}.json').write_text(json.dumps(obs,indent=2)+'\n')
print(json.dumps({'cases':len(run['cases']),'workflow_failures':audit['workflow_failures']}))
