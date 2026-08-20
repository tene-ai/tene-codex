import hashlib,importlib.util,json,pathlib,subprocess,sys,tempfile,unittest
ROOT=pathlib.Path(__file__).resolve().parents[1];SCRIPT=ROOT/'scripts/requirements-audit.py'
SPEC=importlib.util.spec_from_file_location('requirements_audit',SCRIPT);AUDIT=importlib.util.module_from_spec(SPEC);SPEC.loader.exec_module(AUDIT)
class RequirementsAuditTests(unittest.TestCase):
    def test_current_repository_has_complete_required_id_ranges(self):
        r=subprocess.run([sys.executable,str(SCRIPT),'--root',str(ROOT),'--no-exec'],capture_output=True,text=True)
        v=json.loads(r.stdout)
        self.assertEqual(v['coverage'],{'functional_requirements':11,'product_acceptance':8,'work_packages':14})
        self.assertEqual(v['missing'],[])
    def test_missing_symbol_and_unknown_command_fail_closed(self):
        manifest=json.loads((ROOT/'docs/release/semantic-contracts.json').read_text())
        manifest['functional_requirements']['FR-01']['symbols']=[['internal/domain/types.go','DefinitelyMissingSymbol']]
        manifest['functional_requirements']['FR-01']['commands']=['untrusted-shell-command']
        with tempfile.TemporaryDirectory() as directory:
            path=pathlib.Path(directory)/'manifest.json';path.write_text(json.dumps(manifest))
            r=subprocess.run([sys.executable,str(SCRIPT),'--root',str(ROOT),'--manifest',str(path),'--no-exec'],capture_output=True,text=True)
        self.assertNotEqual(r.returncode,0)
        missing=json.loads(r.stdout)['missing']
        self.assertTrue(any('symbol-missing' in value for value in missing),missing)
        self.assertTrue(any('command-unknown' in value for value in missing),missing)
    def test_final_mode_matches_active_state_contract(self):
        state=json.loads((ROOT/'.tene-workflow/project.json').read_text())
        r=subprocess.run([sys.executable,str(SCRIPT),'--root',str(ROOT),'--final','--no-exec'],capture_output=True,text=True)
        expected=1 if state.get('active_sprint_id') else 0
        self.assertEqual(r.returncode,expected,r.stdout+r.stderr)
        self.assertEqual(json.loads(r.stdout)['passed'],expected==0)
    def test_legacy_evidence_requires_passed_case_link_and_matching_provenance(self):
        with tempfile.TemporaryDirectory() as directory:
            root=pathlib.Path(directory);(root/'.tene-workflow').mkdir();artifact=root/'proof.txt';artifact.write_text('verified execution\n')
            evidence={'evidence_id':'ev','sprint_id':'sprint','kind':'legacy','uri':'proof.txt','sha256':hashlib.sha256(artifact.read_bytes()).hexdigest(),'ac_ids':['ac'],'redaction_status':'passed'}
            case={'case_id':'case','ac_ids':['ac'],'status':'passed','evidence_ids':[]}
            state={'active_sprint_id':'','gaps':{},'tasks':{},'sprints':{'sprint':{'sprint_id':'sprint','phase':'archived','last_qa_status':'passed'}},'acceptance_criteria':{'ac':{'ac_id':'ac','priority':'blocking'}},'evidence':{'ev':evidence},'qa_runs':{'run':{'run_id':'run','sprint_id':'sprint','status':'passed','cases':[case]}}}
            path=root/'.tene-workflow/project.json';path.write_text(json.dumps(state))
            failures,_=AUDIT.state_findings(root,False);self.assertIn('unverified:ac',failures)
            case['evidence_ids']=['ev'];path.write_text(json.dumps(state))
            failures,_=AUDIT.state_findings(root,False);self.assertNotIn('unverified:ac',failures)
            evidence['run_id']='different-run';path.write_text(json.dumps(state))
            failures,_=AUDIT.state_findings(root,False);self.assertIn('unverified:ac',failures)
if __name__=='__main__':unittest.main()
