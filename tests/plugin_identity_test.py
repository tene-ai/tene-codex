# Copyright 2026 Kay Kim (kay@agentkay.it)
# SPDX-License-Identifier: Apache-2.0

import json
from pathlib import Path
import re
import unittest


ROOT = Path(__file__).resolve().parents[1]
SUFFIXES = ("sprint", "prd", "plan", "design", "loop-check", "qa", "report", "status", "secrets")


class PluginIdentityTests(unittest.TestCase):
    def test_plugin_and_marketplace_namespace_is_tene(self):
        manifest = json.loads((ROOT / ".codex-plugin" / "plugin.json").read_text())
        marketplace = json.loads((ROOT / ".agents" / "plugins" / "marketplace.json").read_text())
        self.assertEqual(manifest["name"], "tene")
        self.assertEqual(manifest["interface"]["displayName"], "tene")
        self.assertEqual(marketplace["plugins"][0]["name"], "tene")
        self.assertEqual(marketplace["plugins"][0]["interface"]["displayName"], "tene")

    def test_nine_skills_have_concise_matching_identities(self):
        skill_dirs = sorted(path.name for path in (ROOT / "skills").iterdir() if path.is_dir())
        self.assertEqual(skill_dirs, sorted(SUFFIXES))
        for suffix in SUFFIXES:
            with self.subTest(suffix=suffix):
                skill = (ROOT / "skills" / suffix / "SKILL.md").read_text()
                prompt = (ROOT / "skills" / suffix / "agents" / "openai.yaml").read_text()
                match = re.search(r"^name:\s*(.+)$", skill, re.MULTILINE)
                self.assertIsNotNone(match)
                self.assertEqual(match.group(1).strip(), suffix)
                self.assertIn(f"$tene:{suffix}", prompt)

    def test_release_smoke_exposes_each_runtime_stage(self):
        smoke = (ROOT / "scripts" / "release-smoke.sh").read_text()
        stages = (
            "package-manifest-sbom",
            "bundled-cli",
            "routing-explicit-implicit",
            "tampered-binary-rejection",
            "portable-workflow-matrix",
            "update-remove-simulation",
            "project-state-preservation",
        )
        for stage in stages:
            with self.subTest(stage=stage):
                self.assertIn(f'stage {stage} passed', smoke)


if __name__ == "__main__":
    unittest.main()
