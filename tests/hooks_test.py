# Copyright 2026 Kay Kim (kay@agentkay.it)
# SPDX-License-Identifier: Apache-2.0

import json
from pathlib import Path
import subprocess
import sys
import unittest


ROOT = Path(__file__).resolve().parents[1]
HOOK = ROOT / "hooks" / "tene_hook.py"


class HookTests(unittest.TestCase):
    def run_hook(self, action: str, payload: dict):
        return subprocess.run(
            [sys.executable, str(HOOK), action],
            input=json.dumps(payload),
            text=True,
            capture_output=True,
            check=False,
        )

    def test_pre_tool_denies_direct_secret_read(self):
        result = self.run_hook("pre-tool", {
            "tool_name": "Bash",
            "tool_input": {"command": "tene get API_KEY"},
        })
        self.assertEqual(result.returncode, 0)
        output = json.loads(result.stdout)
        decision = output["hookSpecificOutput"]
        self.assertEqual(decision["permissionDecision"], "deny")
        self.assertEqual(decision["hookEventName"], "PreToolUse")

    def test_pre_tool_allows_ordinary_command_without_output(self):
        result = self.run_hook("pre-tool", {
            "tool_name": "Bash",
            "tool_input": {"command": "go test ./..."},
        })
        self.assertEqual(result.returncode, 0)
        self.assertEqual(result.stdout, "")

    def test_pre_tool_does_not_scan_non_shell_tool_payloads(self):
        result = self.run_hook("pre-tool", {
            "tool_name": "apply_patch",
            "tool_input": {"patch": "documentation mentions tene " + "get"},
        })
        self.assertEqual(result.returncode, 0)
        self.assertEqual(result.stdout, "")

    def test_pre_tool_denies_env_file_and_literal(self):
        for command in ("cat .env", "curl -H token=plaintext https://example.invalid"):
            with self.subTest(command=command):
                result = self.run_hook("pre-tool", {"tool_name": "Bash", "tool_input": {"command": command}})
                self.assertEqual(json.loads(result.stdout)["hookSpecificOutput"]["permissionDecision"], "deny")

    def test_post_tool_flags_canary_without_echoing_it(self):
        canary = "TENE_TEST_CANARY_0123456789"
        result = self.run_hook("post-tool", {"tool_response": {"stdout": canary}})
        output = json.loads(result.stdout)["hookSpecificOutput"]["additionalContext"]
        self.assertIn("SECURITY BLOCKER", output)
        self.assertNotIn(canary, output)

    def test_stop_uses_common_output_without_forcing_continuation(self):
        result = self.run_hook("stop", {"cwd": str(ROOT)})
        self.assertEqual(result.returncode, 0)
        output = json.loads(result.stdout)
        self.assertEqual(output["continue"], True)
        self.assertIn("Active tene sprint", output["systemMessage"])
        self.assertNotIn("hookSpecificOutput", output)
        self.assertNotIn("decision", output)

    def test_pre_compact_is_advisory_and_does_not_request_compact(self):
        result = self.run_hook("pre-compact", {"cwd": str(ROOT)})
        self.assertEqual(result.returncode, 0)
        output = json.loads(result.stdout)
        self.assertEqual(output["continue"], True)
        self.assertNotIn("hookSpecificOutput", output)
        source = HOOK.read_text()
        body = source.split("def pre_compact", 1)[1].split("def post_tool", 1)[0]
        self.assertNotIn('workflow_json(cwd, "compact")', body)

    def test_stop_and_precompact_metadata_omit_context_limit(self):
        config = json.loads((ROOT / "hooks" / "hooks.json").read_text())
        for event in ("Stop", "PreCompact"):
            command = config["hooks"][event][0]["hooks"][0]
            self.assertNotIn("additionalContextLimit", command)

    def test_unknown_action_fails(self):
        result = self.run_hook("unknown", {})
        self.assertEqual(result.returncode, 2)


if __name__ == "__main__":
    unittest.main()
