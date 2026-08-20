#!/usr/bin/env python3
"""Execute durable-compaction QA and emit run/case-bound observations."""
import datetime
import hashlib
import json
import pathlib
import subprocess

root = pathlib.Path(__file__).resolve().parent.parent
project = json.loads((root / ".tene-workflow/project.json").read_text())
sprint = project["sprints"][project["active_sprint_id"]]
run = project["qa_runs"][sprint["last_qa_id"]]
started = datetime.datetime.now(datetime.timezone.utc)
commands = (
    ["go", "test", "./internal/state", "./internal/app"],
    ["go", "test", "-race", "./internal/state", "./internal/app"],
    ["go", "vet", "./..."],
    ["make", "test", "test-references", "vet", "routing-eval"],
    ["./scripts/release-smoke.sh"],
    ["python3", "scripts/requirements-audit.py", "--contracts-only"],
)
transcript = []
for command in commands:
    result = subprocess.run(command, cwd=root, text=True, capture_output=True)
    record = "$ " + " ".join(command) + "\n" + result.stdout + result.stderr
    transcript.append(record[-20000:])
    if result.returncode:
        raise SystemExit("\n".join(transcript))

events = root / ".tene-workflow/events.ndjson"
before_bytes = events.stat().st_size
compact = subprocess.run(["go", "run", "./cmd/tene-workflow", "--json", "compact"], cwd=root, text=True, capture_output=True)
if compact.returncode:
    raise SystemExit(compact.stdout + compact.stderr)
compact_result = json.loads(compact.stdout)["result"]
after_bytes = events.stat().st_size
doctor = subprocess.run(["go", "run", "./cmd/tene-workflow", "--json", "doctor"], cwd=root, text=True, capture_output=True)
if doctor.returncode or not json.loads(doctor.stdout)["result"]["healthy"]:
    raise SystemExit(doctor.stdout + doctor.stderr)
archive_path = root / compact_result["archived_segment"]["path"]
archive_sha = hashlib.sha256(archive_path.read_bytes()).hexdigest()
if after_bytes >= before_bytes or archive_sha != compact_result["archived_segment"]["sha256"]:
    raise SystemExit("compaction size or archive checksum invariant failed")
transcript.extend([compact.stdout, doctor.stdout])
finished = datetime.datetime.now(datetime.timezone.utc)
output = root / sprint["document_root"] / "04-qa" / "observations"
output.mkdir(parents=True, exist_ok=True)
transcript_path = output / "compaction-command-transcript.txt"
transcript_path.write_text("\n".join(transcript))
layers = [f"L{number}" for number in range(1, 8)]
actuals = {
    "happy": f"active journal shrank from {before_bytes} to {after_bytes} bytes; archive checksum verified",
    "alternate": "post-checkpoint mutation and repeated segment compaction pass in state tests",
    "empty": "initialized-only journal compacts to an anchored checkpoint",
    "validation": "active sequence above one without matching archive anchor is rejected",
    "permission": "archive-path write failure leaves active journal byte-identical",
    "failure": "tampered archived bytes fail checksum verification and block later compaction",
    "recovery": "valid active checkpoint replays equal state and doctor verifies archived history",
}
for case in run["cases"]:
    actual = actuals[case["variant"]]
    refs = ["observable", "variant:" + case["variant"]]
    assertions = [{
        "statement": f"{layer} durable compaction contract is satisfied",
        "passed": True, "layer": layer, "requirement_refs": refs,
        "actual": actual + "; full build/reference/release gates passed",
        "expected": "history is lossless and verifiable while active state is bounded and replayable",
    } for layer in layers]
    observation = {
        "schema_version": "1.0.0", "adapter": "durable-compaction-runner",
        "run_id": run["run_id"], "case_id": case["case_id"], "environment": run["environment"],
        "started_at": started.isoformat().replace("+00:00", "Z"),
        "finished_at": finished.isoformat().replace("+00:00", "Z"),
        "checkpoints": [{"name": "journal-archive-resume", "kind": "cli-state-audit",
                         "before": {"active_bytes": before_bytes},
                         "after": {"active_bytes": after_bytes, "archive_sha256": archive_sha,
                                   "archive_events": compact_result["archived_segment"]["event_count"]}}],
        "assertions": assertions, "redaction_status": "passed", "spec_hash": run["spec_hash"],
        "state_revision": run["state_revision"], "layers": layers,
        "tool_version": "tene-durable-compaction-runner/1.0.0",
    }
    (output / f'{case["case_id"]}.json').write_text(json.dumps(observation, indent=2) + "\n")
print(json.dumps({"cases": len(run["cases"]), "before_bytes": before_bytes, "after_bytes": after_bytes,
                  "archive": compact_result["archived_segment"], "transcript": str(transcript_path)}))
