#!/usr/bin/env python3
"""Codex lifecycle adapter for tene-codex.

The hook is intentionally advisory except for a small, fail-closed set of
secret-exposure operations. Canonical workflow mutations remain in the CLI.
"""

from __future__ import annotations

import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys


DENIED = (
    re.compile(r"(?<![\w-])tene\s+get(?:\s|$|[\"'])", re.I),
    re.compile(r"(?<![\w-])tene\s+export(?:\s|$|[\"'])", re.I),
    re.compile(r"(?<![\w-])(?:env|printenv)(?:\s|$|[\"'])", re.I),
    re.compile(r"(?:^|[\s'\"/])\.tene(?:/|\\)", re.I),
)


def read_input() -> dict:
    try:
        value = json.load(sys.stdin)
        return value if isinstance(value, dict) else {}
    except (json.JSONDecodeError, OSError):
        return {}


def workflow_command() -> list[str] | None:
    installed = shutil.which("tene-workflow")
    if installed:
        return [installed]
    root = Path(os.environ.get("PLUGIN_ROOT", Path(__file__).resolve().parent.parent))
    wrapper = root / "scripts" / "tene-workflow"
    if wrapper.exists():
        return [str(wrapper)]
    return None


def workflow_json(cwd: str, *args: str) -> dict | None:
    command = workflow_command()
    if not command or not (Path(cwd) / ".tene-workflow").exists():
        return None
    try:
        completed = subprocess.run(
            [*command, "--root", cwd, "--json", *args],
            text=True,
            capture_output=True,
            timeout=8,
            check=False,
            env={**os.environ, "TENE_HOOK": "1"},
        )
        return json.loads(completed.stdout) if completed.stdout.strip() else None
    except (OSError, subprocess.TimeoutExpired, json.JSONDecodeError):
        return None


def additional(event: str, text: str) -> None:
    print(json.dumps({
        "hookSpecificOutput": {
            "hookEventName": event,
            "additionalContext": text,
        }
    }))


def session_start(data: dict) -> None:
    cwd = data.get("cwd") or os.getcwd()
    status = workflow_json(cwd, "status")
    if not status or not status.get("ok"):
        return
    result = status.get("result", {})
    sprint = result.get("active_sprint")
    if not sprint:
        return
    additional("SessionStart", (
        "tene-codex has an active sprint. Resume it before implementation. "
        f"Sprint={sprint.get('sprint_id')} phase={sprint.get('phase')} "
        f"revision={status.get('revision')}. Run $tene-status for the phase context."
    ))


def pre_tool(data: dict) -> None:
    raw = json.dumps(data.get("tool_input", {}), ensure_ascii=False)
    if any(pattern.search(raw) for pattern in DENIED):
        print(json.dumps({
            "hookSpecificOutput": {
                "hookEventName": "PreToolUse",
                "permissionDecision": "deny",
                "permissionDecisionReason": (
                    "tene-codex blocked direct secret retrieval, export, environment dumping, "
                    "or .tene vault access. Use `tene run --env <name> -- <command>`."
                ),
            }
        }))


def pre_compact(data: dict) -> None:
    cwd = data.get("cwd") or os.getcwd()
    snapshot = workflow_json(cwd, "compact")
    if snapshot and snapshot.get("ok"):
        additional("PreCompact", "tene-codex saved a workflow snapshot. Reload $tene-status after compaction.")


def subagent_start(data: dict) -> None:
    cwd = data.get("cwd") or os.getcwd()
    status = workflow_json(cwd, "status")
    if status and status.get("ok") and status.get("result", {}).get("active_sprint"):
        additional("SubagentStart", (
            "Work only on the bounded delegated task. Read the current tene context pack, "
            "cite evidence and file locators, and return findings to the parent. Do not mutate "
            "workflow state directly or claim a gate passed without evidence."
        ))


def stop(data: dict) -> None:
    cwd = data.get("cwd") or os.getcwd()
    status = workflow_json(cwd, "status")
    if not status or not status.get("ok"):
        return
    sprint = status.get("result", {}).get("active_sprint")
    if sprint and sprint.get("phase") not in ("archived", None):
        additional("Stop", (
            f"Active tene sprint {sprint.get('sprint_id')} remains in {sprint.get('phase')}. "
            "The state is resumable; report the current gate and next action rather than implying the sprint is archived."
        ))


def main() -> int:
    action = sys.argv[1] if len(sys.argv) > 1 else ""
    data = read_input()
    actions = {
        "session-start": session_start,
        "pre-tool": pre_tool,
        "pre-compact": pre_compact,
        "subagent-start": subagent_start,
        "stop": stop,
    }
    handler = actions.get(action)
    if not handler:
        return 2
    handler(data)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
