#!/usr/bin/env python3
"""Complete the public Sprint lifecycle in three clean repository stacks."""
import argparse
import datetime
import json
import os
import pathlib
import re
import subprocess
import tempfile

LAYERS = [f"L{n}" for n in range(1, 8)]
STACKS = {
    "go-cli": {
        "artifact": "cmd/portable/main.go",
        "files": {"go.mod": "module portable.example/go\n\ngo 1.24\n", "cmd/portable/main.go": "package main\nfunc main() {}\n"},
    },
    "nextjs-web": {
        "artifact": "app/page.tsx",
        "files": {"package.json": '{"name":"portable-next","private":true}\n', "app/page.tsx": "export default function Page(){return <main>portable</main>}\n", "app/api/orders/route.ts": "export async function POST(){return Response.json({ok:true})}\n"},
    },
    "python-api-worker": {
        "artifact": "app/api.py",
        "files": {"pyproject.toml": '[project]\nname="portable-python"\nversion="0.1.0"\n', "app/api.py": "def create_order(data): return {'ok': True, 'data': data}\n", "worker.py": "def run(job): return job\n"},
    },
}

parser = argparse.ArgumentParser()
parser.add_argument("--cli", required=True)
parser.add_argument("--workspace")
args = parser.parse_args()
cli = str(pathlib.Path(args.cli).resolve())
workspace_context = tempfile.TemporaryDirectory(prefix="tene-portable-workflow-") if not args.workspace else None
workspace = pathlib.Path(args.workspace or workspace_context.name)
workspace.mkdir(parents=True, exist_ok=True)
env = os.environ.copy()

def call(project, *command):
    result = subprocess.run([cli, "--root", str(project), "--json", *command], text=True, capture_output=True, env=env)
    try:
        payload = json.loads(result.stdout)
    except json.JSONDecodeError as exc:
        raise RuntimeError(f"invalid JSON for {command}: {result.stdout}\n{result.stderr}") from exc
    if result.returncode or not payload.get("ok"):
        raise RuntimeError(f"command failed {command}: {payload}\n{result.stderr}")
    return payload["result"]

def complete_documents(project, sprint):
    root = project / sprint["document_root"]
    for path in root.glob("*/*.md"):
        text = path.read_text()
        text = re.sub(r"(?m)^(## [^\n]+)\n(?=\s*(?:<!--|##|$))", r"\1\n\nPortable public workflow contract verified in an isolated stack repository.\n", text)
        path.write_text(text)

def run_stack(name, config):
    project = workspace / name
    for relative, content in config["files"].items():
        path = project / relative
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content)
    subprocess.run(["git", "init", "-q"], cwd=project, check=True)
    subprocess.run(["git", "add", "."], cwd=project, check=True)
    subprocess.run(["git", "-c", "user.name=tene-smoke", "-c", "user.email=smoke@invalid.example", "commit", "-qm", "fixture baseline"], cwd=project, check=True)
    call(project, "init", "--name", name, "--profile", "light")
    created = call(project, "sprint", "create", "--title", f"{name} portable lifecycle")
    sprint = created["sprint"]
    complete_documents(project, sprint)
    call(project, "phase", "transition", "prd")
    captured = call(project, "intent", "capture", "--statement", f"{name} users retain spec-driven lifecycle semantics", "--rationale", "portable plugin behavior must not depend on repository stack", "--actors", "developer,Codex agent", "--outcomes", "archived sprint,verified evidence,resumable state", "--policies", "no skipped phases,all blocking evidence required", "--ac", f"{name} completes the public Sprint lifecycle", "--observable", "archive exists and doctor is healthy")
    intent_id, ac_id = captured["intent"]["intent_id"], captured["criterion"]["ac_id"]
    for relative in ("00-prd/00-prd.md", "02-design/00-design.md"):
        path = project / sprint["document_root"] / relative
        path.write_text(path.read_text() + f"\nPortable traceability: {intent_id} {ac_id}\n")
    call(project, "intent", "confirm", intent_id)
    call(project, "document", "sync", "--apply")
    call(project, "phase", "transition", "plan")
    task = call(project, "task", "add", "--title", f"Verify {name} lifecycle", "--layer", "interface", "--ac", ac_id)
    for relative in ("01-plan/00-plan.md", "02-design/00-design.md"):
        path = project / sprint["document_root"] / relative
        path.write_text(path.read_text() + f"\nPortable task traceability: {task['task_id']}\n")
    call(project, "document", "sync", "--apply")
    call(project, "phase", "transition", "design")
    call(project, "document", "sync", "--apply")
    call(project, "phase", "transition", "do")
    call(project, "task", "start", task["task_id"])
    call(project, "task", "artifact", task["task_id"], "--path", config["artifact"])
    call(project, "task", "complete", task["task_id"])
    call(project, "phase", "transition", "loop-check")
    call(project, "document", "sync", "--apply")
    checked = call(project, "loop", "check")
    if not checked["passed"]:
        raise RuntimeError(f"loop check failed for {name}: {checked}")
    call(project, "loop", "iterate", "--outcome", "passed", "--summary", "portable stack artifact and document traceability match")
    call(project, "phase", "transition", "qa")
    call(project, "document", "sync", "--apply")
    run = call(project, "qa", "plan")
    observation_dir = project / sprint["document_root"] / "04-qa" / "observations"
    observation_dir.mkdir(parents=True, exist_ok=True)
    now = datetime.datetime.now(datetime.timezone.utc)
    for case in run["cases"]:
        observation = {
            "schema_version": "1.0.0", "adapter": "portable-public-cli",
            "run_id": run["run_id"], "case_id": case["case_id"], "environment": run["environment"],
            "started_at": now.isoformat().replace("+00:00", "Z"),
            "finished_at": (now + datetime.timedelta(seconds=1)).isoformat().replace("+00:00", "Z"),
            "checkpoints": [{"name": "clean-stack-lifecycle", "kind": "cli-doc-state",
                             "before": {"stack": name, "phase": "draft"},
                             "after": {"stack": name, "phase": "qa", "artifact": config["artifact"]}}],
            "assertions": [{"statement": f"{layer} portable {case['variant']} contract passed", "passed": True,
                            "layer": layer, "requirement_refs": ["observable", "variant:" + case["variant"]],
                            "actual": f"isolated {name} {layer} boundary preserved its state, document or gate contract",
                            "expected": f"{layer} public workflow boundary behaves consistently across stacks"} for layer in LAYERS],
            "redaction_status": "passed", "spec_hash": run["spec_hash"], "state_revision": run["state_revision"],
            "layers": LAYERS, "tool_version": "portable-public-cli/1.0.0",
        }
        path = observation_dir / f"{case['case_id']}.json"
        path.write_text(json.dumps(observation, indent=2) + "\n")
        call(project, "qa", "observe", case["case_id"], "--input", str(path))
    evaluated = call(project, "qa", "evaluate")
    if evaluated["status"] != "passed":
        raise RuntimeError(f"QA failed for {name}: {evaluated}")
    call(project, "phase", "transition", "report")
    call(project, "report", "generate")
    call(project, "report", "validate")
    call(project, "sprint", "archive")
    doctor = call(project, "doctor")
    verified = call(project, "evidence", "verify")
    status = call(project, "status")
    archived = list((project / "docs/sprints/_archive").glob("*/*/99-archive/archive-manifest.json"))
    if not doctor["healthy"] or not verified["valid"] or status.get("active_sprint") is not None or len(archived) != 1:
        raise RuntimeError(f"post-archive verification failed for {name}")
    return {"stack": name, "revision": doctor["revision"], "events": doctor["events"], "evidence": verified["count"], "archive_manifest": str(archived[0].relative_to(project))}

results = [run_stack(name, config) for name, config in STACKS.items()]
print(json.dumps({"passed": True, "stacks": results, "workspace": str(workspace)}, indent=2))
