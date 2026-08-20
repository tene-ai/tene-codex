import json,pathlib,subprocess,sys,tempfile,unittest
ROOT=pathlib.Path(__file__).resolve().parents[1];SCRIPT=ROOT/'scripts/requirements-audit.py'
class RequirementsAuditTests(unittest.TestCase):
    def test_current_repository_has_complete_required_id_ranges(self):
        r=subprocess.run([sys.executable,str(SCRIPT),'--root',str(ROOT)],capture_output=True,text=True);self.assertEqual(r.returncode,0,r.stdout+r.stderr);v=json.loads(r.stdout);self.assertEqual(v['coverage'],{'functional_requirements':11,'product_acceptance':8,'work_packages':14});self.assertEqual(v['missing'],[])
    def test_final_mode_matches_active_state_contract(self):
        state=json.loads((ROOT/'.tene-workflow/project.json').read_text())
        r=subprocess.run([sys.executable,str(SCRIPT),'--root',str(ROOT),'--final'],capture_output=True,text=True)
        expected=1 if state.get('active_sprint_id') else 0
        self.assertEqual(r.returncode,expected,r.stdout+r.stderr)
        self.assertEqual(json.loads(r.stdout)['passed'],expected==0)
if __name__=='__main__':unittest.main()
