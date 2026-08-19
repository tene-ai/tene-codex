---
schema_version: 1.0.0
document_type: design
project_id: project_0000380v060bf8678xh7wgm6y8
sprint_id: sprint_0000380v2b2q0668t7tajkj3m0
phase: design
status: draft
revision: 88
intent_ids: []
generated_at: 2026-08-19T17:30:10Z
generated_by: tene-workflow
---

# design — Code Intelligence and Intent QA Adapters

<!-- tene:section:purpose -->
## Purpose

Specify implementable boundaries for deterministic code understanding and whole-flow QA evidence.

<!-- tene:section:scope -->
## Scope

Additive Go packages and CLI actions; backward-compatible state fields; JSON artifact schemas; Codex skill routing.

<!-- tene:section:layers -->
## Layers

`app` is Interface/orchestration, `codeintel` and `qaadapter` are Business Logic ports plus local adapters, state/evidence are Persistence, and external binaries/browser tools are Infrastructure capabilities.

<!-- tene:section:six-questions -->
## Six questions

`codeintel.Report` contains Components. Each Component contains name, locator, references, calls, inputs, outputs/effects, primary/secondary layers, provider and confidence. Empty provider support is represented in Unknown, not omitted.

<!-- tene:section:traceability -->
## Traceability

AC1 → Provider/Report/graph commands. AC2 → Capability/Charter/Observation/qa commands. AC3 → package and CLI tests plus dogfood evidence.

<!-- tene:section:decisions -->
## Decisions

Provider selection is capability negotiation, not hard-coded success. Go AST is the deterministic baseline. CodeGraph is probed first only under repository policy, and its raw query is advisory until a stable machine schema exists.

<!-- tene:section:freeform -->
## Freeform

<!-- tene:section:components -->
## Components

- `codeintel.Discover(root)` returns provider capabilities.
- `codeintel.Analyze(root, options)` walks bounded Go source, parses declarations/imports/calls/signatures, classifies layers, and emits a graph delta plus report.
- `qaadapter.Discover(root)` probes Go/native, npm scripts, Playwright config and external browser-observation support.
- `qaadapter.WriteObservation` validates a structured observation and atomically writes the evidence artifact.
- app merges derived code nodes with authored graph state and registers observation evidence through the existing hash/redaction path.

<!-- tene:section:interfaces -->
## Interfaces

`graph providers`; `graph build`; `graph understand [--changed] [--path CSV]`; `qa capabilities`; `qa plan`; `qa observe <case-id> --input FILE`. Browser/Chrome skills create the input JSON after performing the journey; the CLI never assumes a particular MCP server name.

<!-- tene:section:data -->
## Data

QACase gains actor, preconditions, steps, forbidden outcomes, required layer dispositions and risk. Observation contains schema version, adapter, case/run IDs, checkpoints/assertions, timestamps and redaction status. Graph nodes use stable `file:<path>` and `symbol:<path>:<name>` identifiers.

<!-- tene:section:state-transitions -->
## State transitions

Analysis itself is read-only; graph build commits one `GraphRebuilt` event. Observation import validates active run/case, writes artifact atomically, commits EvidenceRegistered and QACaseRecorded semantics, then QA evaluate applies the unchanged strict gate.

<!-- tene:section:failures -->
## Failures

Malformed AST becomes a diagnostic; unreadable or oversized files are skipped. Invalid observation/run/case, failed assertion, secret pattern, stale evidence or unavailable provider fails closed with actionable codes.

<!-- tene:section:security -->
## Security

No shell, no environment dump, no secret values in capabilities or artifacts. Paths must resolve inside the project. Observation artifacts are scanned before state commit and content-hashed afterward.

<!-- tene:section:tests -->
## Tests

Fixture Go modules validate declarations/import/call/signature/effect/layer output. Capability fixtures cover Go/npm/Playwright. CLI tests prove graph persistence, observation import, case linkage, negative path traversal and QA gate behavior.
