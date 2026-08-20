# Copyright 2026 Kay Kim (kay@agentkay.it)
# SPDX-License-Identifier: Apache-2.0

from pathlib import Path
import unittest


ROOT = Path(__file__).resolve().parents[1]


class LanguageContractTests(unittest.TestCase):
    def test_workflow_reference_defines_user_conversation_language(self):
        workflow = (ROOT / "references" / "workflow.md").read_text()
        self.assertIn("## Conversation language", workflow)
        self.assertIn("language currently used by the user", workflow)
        self.assertIn("explicit user language request overrides", workflow)

    def test_every_skill_carries_the_language_contract(self):
        skills = sorted((ROOT / "skills").glob("*/SKILL.md"))
        self.assertEqual(len(skills), 9)
        for skill in skills:
            with self.subTest(skill=skill.parent.name):
                body = skill.read_text()
                self.assertIn("user's current conversation language", body)
                self.assertIn("workflow language contract", body)


if __name__ == "__main__":
    unittest.main()
