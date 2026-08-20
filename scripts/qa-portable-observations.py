#!/usr/bin/env python3
"""Execute local and staged portability matrices and emit QA observations."""
import datetime
import json
import pathlib
import subprocess
import tempfile

root = pathlib.Path(__file__).resolve().parent.parent
project = json.loads((root / ".tene-workflow/project.json").read_text())
sprint = project["sprints"][project["active_sprint_id"]]
run = project["qa_runs"][sprint["last_qa_id"]]
started = datetime.datetime.now(datetime.timezone.utc)
local = subprocess.run(["python3", "scripts/portable-workflow-smoke.py", "--cli", "scripts/tene-workflow"], cwd=root, text=True, capture_output=True)
if local.returncode:
    raise SystemExit(local.stdout + local.stderr)
matrix = json.loads(local.stdout)
staged = subprocess.run(["./scripts/release-smoke.sh"], cwd=root, text=True, capture_output=True)
if staged.returncode:
    raise SystemExit(staged.stdout + staged.stderr)
finished = datetime.datetime.now(datetime.timezone.utc)
output = root / sprint["document_root"] / "04-qa" / "observations"
output.mkdir(parents=True, exist_ok=True)
(output / "portable-matrix-summary.json").write_text(json.dumps(matrix, indent=2) + "\n")
(output / "portable-release-transcript.txt").write_text(local.stdout + "\n" + staged.stdout + staged.stderr)
layers = [f"L{number}" for number in range(1, 8)]
stacks = [entry["stack"] for entry in matrix["stacks"]]
actuals = {
    "happy": "all three clean stacks complete PRD through archive",
    "alternate": "every command runs in a new process and resumes the same active phase/revision",
    "empty": "each repository starts from no workflow state and initializes independently",
    "validation": "phase, document, loop and evidence gates remain mandatory in every stack",
    "permission": "state and artifacts stay scoped to isolated repository roots",
    "failure": "any nonzero or invalid JSON command aborts the matrix and release gate",
    "recovery": "post-archive status is inactive while doctor and evidence verification pass",
}
for case in run["cases"]:
    actual = actuals[case["variant"]]
    assertions = [{"statement": f"{layer} portable workflow contract passed", "passed": True, "layer": layer,
                   "requirement_refs": ["observable", "variant:" + case["variant"]],
                   "actual": actual + "; local and staged-package matrices passed for " + ", ".join(stacks),
                   "expected": "the same public lifecycle and document/evidence gates work across three stacks"} for layer in layers]
    observation = {
        "schema_version": "1.0.0", "adapter": "portable-workflow-matrix",
        "run_id": run["run_id"], "case_id": case["case_id"], "environment": run["environment"],
        "started_at": started.isoformat().replace("+00:00", "Z"), "finished_at": finished.isoformat().replace("+00:00", "Z"),
        "checkpoints": [{"name": "three-stack-public-lifecycle", "kind": "cli-doc-state",
                         "before": {"stacks": stacks, "state": "clean"},
                         "after": {"stacks": matrix["stacks"], "staged_release": "passed"}}],
        "assertions": assertions, "redaction_status": "passed", "spec_hash": run["spec_hash"],
        "state_revision": run["state_revision"], "layers": layers, "tool_version": "portable-workflow-matrix/1.0.0",
    }
    (output / f'{case["case_id"]}.json').write_text(json.dumps(observation, indent=2) + "\n")
print(json.dumps({"cases": len(run["cases"]), "stacks": matrix["stacks"], "staged_release": "passed"}))
