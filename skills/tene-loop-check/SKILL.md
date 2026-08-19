---
name: tene-loop-check
description: Compare active sprint requirements, plan, design, and actual code changes; record and resolve implementation gaps before QA.
---

# tene Loop Check

Require `loop-check` or a QA repair transition. Load the context pack, git diff, changed symbols, and graph. Compare actual behavior against every confirmed intent, blocking criterion, planned task, and design contract.

Classify gaps as `missing`, `mismatch`, `unverified`, `regression`, or `debt`. Record each gap through the CLI. For changed components, cover all Understanding Layers and Six Questions; mark unknowns as gaps rather than inventing answers. Repair in `do`, rerun focused tests, and resolve a gap only with cited evidence.

Use a read-only evaluator subagent when available. Do not transition to QA while a blocking gap remains. Read [workflow](../../references/workflow.md).

