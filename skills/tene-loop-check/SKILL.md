---
name: tene-loop-check
description: Compare active sprint requirements, plan, design, and actual code changes; record and resolve implementation gaps before QA.
---

# tene Loop Check

Require `loop-check` or a QA repair transition. Link every changed implementation/configuration artifact to its owning task with `task artifact ID --path FILE`, then run `loop check`. The deterministic analyzer compares confirmed intent and AC IDs across PRD/plan/design, AC→task coverage, task→changed-file ownership, linked-file existence, and executable design contracts. A design may assert `<!-- tene:contract path="FILE" symbol="SYMBOL" -->` or forbid a dependency with `<!-- tene:forbid path="FILE" contains="TEXT" -->`. Never manually rewrite analyzer output.

The analyzer creates fingerprinted gaps without duplicates, resolves only its own gaps when a rerun proves the condition disappeared, and reopens the same gap if drift returns. Add manual gaps only for semantic or runtime findings outside deterministic coverage. For changed components, cover all Understanding Layers and Six Questions; preserve provider unknowns. Record every evaluator pass with `loop iterate` and stop when its bounded budget is exhausted. Repair in `do`, rerun focused tests, and resolve manually authored gaps only with rationale plus registered evidence. Defer only eligible work with reason, owner, and target Sprint; never defer security or evidence-integrity gaps.

Use a read-only evaluator subagent when available. Do not transition to QA while a blocking gap remains. Read [workflow](../../references/workflow.md).
