---
name: tene-design
description: Investigate a codebase and produce implementation-level component, interface, data, failure, and test design for an active tene sprint.
---

# tene Design

Require the `design` phase and load the context pack. If `.codegraph/` exists, use CodeGraph before text search. Otherwise prefer language-native symbols and targeted search.

Trace all affected Understanding Layers and answer the Six Questions for important existing and proposed components. Define public interfaces, data shapes, state transitions, errors, idempotency, concurrency, security, observability, migration, and test seams. Link design elements back to criteria and tasks. Use a read-only explorer subagent for large independent code paths when available.

Validate the design document and dry-run the transition to `do`. Read [workflow](../../references/workflow.md).

