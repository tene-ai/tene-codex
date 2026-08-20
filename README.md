# tene-codex

**tene-codex is an open-source Codex plugin for spec-driven, intent-aware agentic coding.**

It is designed to turn an informal coding conversation into a durable engineering workflow that Codex can resume, inspect, verify, and improve across sessions. Instead of treating generated code or passing test scripts as proof of completion, tene-codex keeps product intent, implementation decisions, code impact, user journeys, data flows, and QA evidence connected throughout the work.

> Project status: **0.1 public beta**. The complete local workflow, nine skills, hooks, subagent profiles, routing and security evals, recovery, reference matrix, and release packaging are implemented. Public APIs remain pre-1.0 and may change with documented migration.

## Quick Start from Source

Requirements: Go 1.24 or later, Python 3, Git, and a current Codex installation. The separate `tene` CLI is required only for commands that need secrets.

```bash
git clone https://github.com/tene-ai/tene-codex.git
cd tene-codex
make check
go build -o dist/tene-workflow ./cmd/tene-workflow

# In a target repository
/path/to/tene-codex/dist/tene-workflow init --name my-project --profile standard
/path/to/tene-codex/dist/tene-workflow sprint create --title "My feature"
/path/to/tene-codex/dist/tene-workflow status --json
```

During plugin development, `scripts/tene-workflow` finds an installed binary, a bundled platform binary, or runs the source command with Go. Tagged releases are designed to package macOS and Linux binaries with the plugin files and SHA-256 checksums.

After installing or linking the plugin in Codex, start with `$tene:sprint`, resume with `$tene:status`, and use `$tene:qa` for evidence-based verification. Codex requires users to review and trust plugin-bundled hooks before those hooks can run.

Common core commands:

```text
tene-workflow status --json
tene-workflow route --text "UX와 데이터 흐름 QA해줘" --phase qa
tene-workflow phase transition <phase> --dry-run
tene-workflow document validate <phase>
tene-workflow context build --phase design --budget 32768 --output .tene-workflow/cache/context.json
tene-workflow context validate --input .tene-workflow/cache/context.json
tene-workflow approval request --to do --reason "Design reviewed" --requester kay --expires 2026-08-21T00:00:00Z
tene-workflow approval approve <approval-id> --approver kay
tene-workflow phase transition do --approval <approval-id>
tene-workflow loop check
tene-workflow graph providers|build|understand|trace|impact|validate
tene-workflow qa capabilities|plan|execute|observe|case|evaluate|status
tene-workflow waiver create|list|revoke
tene-workflow migrate status|dry-run|apply
tene-workflow doctor [--repair]
tene-workflow evidence register|verify|list
tene-workflow report generate|validate
tene-workflow doctor|compact|clear
```

## Why tene-codex?

Vibe coding is fast, but long-running work can lose the original product intent. An agent may satisfy a local request while missing an upstream policy, a downstream data mutation, a failure path, or a user-facing transition. Unit tests and scripted E2E tests can pass while the complete feature still behaves incorrectly.

tene-codex is intended to address that problem by making the specification and its evidence part of the workflow itself:

- Preserve product intent and acceptance criteria across sessions.
- Require an explicit sprint lifecycle instead of jumping directly to implementation.
- Connect requirements, plans, designs, tasks, code, tests, evidence, and reports through a traceability graph.
- Examine both the whole system and the changed code using four Understanding Layers and six code-understanding questions.
- Verify user experience and data-processing flows, not only individual test scripts.
- Treat QA completion as an evidence-based gate rather than an agent's self-reported conclusion.
- Keep secret values outside the model context and delegate secret injection to the `tene` CLI.

## The Sprint Workflow

Every meaningful coding task is organized as a sprint:

```text
PRD → Plan → Design → Do ↔ Loop Check → QA → Report → Archive
```

- **PRD** captures the problem, actors, product intent, policies, user and data journeys, acceptance criteria, and non-goals.
- **Plan** breaks the requirement into tasks, dependencies, decisions, risks, and verification work.
- **Design** defines components, interfaces, data shapes, state transitions, failure handling, security, and test seams.
- **Do** performs the implementation in traceable work units.
- **Loop Check** compares the PRD, plan, design, and actual changes, then repeats improvement until blocking gaps are resolved.
- **QA** runs the applicable static, unit, contract, integration, system, UX, recovery, and regression checks and collects evidence.
- **Report** explains what changed, why it changed, how it connects to previous sprints, and what remains deferred.
- **Archive** creates a durable sprint record after all required gates and approvals have passed.

A sprint master plan connects multiple sprints into a project-level workflow and task-management system.

## How It Is Intended to Work with Codex

tene-codex combines four Codex-facing mechanisms with a deterministic local workflow engine.

### Skills

Skills provide the user-facing entry points for sprint management, PRD discovery, planning, design, loop checking, QA, reporting, status, and secret-safe execution. Users will be able to invoke them explicitly, such as `$tene:qa`, or trigger them from sufficiently clear natural-language requests.

### Subagents

Subagents are intended to act as specialized workers for product discovery, architecture analysis, code understanding, implementation, test execution, and independent evaluation. The builder and evaluator should use separate contexts whenever possible so that QA is based on evidence rather than the builder's claim that the work is complete.

### Hooks

Hooks provide lifecycle automation and defense in depth: restoring sprint context at session start, detecting changed artifacts after tool use, recording resumable state before compaction or session end, blocking unsafe secret operations, and checking gates before completion. Core correctness must not depend on hooks, because hook availability can vary by Codex version and project trust configuration.

### `tene-workflow` CLI

`tene-workflow` is the local source of truth for:

- sprint, phase, task, and approval state;
- document scaffolding and validation;
- intent, specification, code, test, and evidence relationships;
- context-pack construction;
- loop-check gaps and iteration history;
- QA charters, evidence manifests, and gate verdicts;
- compaction, recovery, migration, and archive operations.
- checksummed immutable journal segments: `compact` bounds the active journal to a replayable checkpoint while `doctor` verifies the complete archived hash chain.

Skills, subagents, and hooks must use this CLI instead of editing workflow state independently. It is implemented as a standalone Go binary so it can be portable and deterministic while remaining separate from the security-sensitive secret manager.

## Engineering Model

### Understanding Layers

Every meaningful change is examined across these layers:

1. **Interface / Entry Point** — UI, CLI, API controller, webhook, scheduler, or command.
2. **Business Logic / Processing Rules** — service, handler, use case, reducer, or domain rule.
3. **Persistence / Data** — database, file, cache, queue, or external API.
4. **Infrastructure / Runtime** — server, container, cloud, authentication, and CI/CD.

### Six Questions

For each important changed component, tene-codex is designed to answer:

1. What is the declared or defined name?
2. In which file is it defined?
3. Where is it imported or referenced?
4. Where is it called or used?
5. What data shape does it receive?
6. What data shape does it return or mutate?

These checks are intended to reduce fragmented changes, hidden coupling, accidental duplication, and agent-generated technical debt.

## Intent-Driven QA

tene-codex treats product intent as executable QA input. Confirmed intent and acceptance criteria are compiled into test charters covering relevant happy, alternate, empty, validation, permission, failure, retry, and recovery paths.

QA can combine project-native tests, API checks, Playwright, Codex browser capabilities, Chrome integration, database or queue observers, logs, traces, screenshots, and manual checkpoints. An evidence manifest connects every blocking acceptance criterion to reproducible observations. A blocking criterion cannot be offset by an average score: all blocking criteria must pass with valid evidence before the sprint can be archived.

`graph understand` uses an existing CodeGraph index when explicitly queried, bounded Go AST for Go, and an uncertainty-honest filesystem fallback for supported non-Go source extensions. It materializes definitions, imports/references, calls/uses, input shape, output/side effects, Understanding Layer, provider and confidence; unavailable semantics remain explicit unknowns. `qa capabilities` discovers native and Playwright runners; `qa execute` permits only discovered allowlisted adapters, while `qa observe` imports schema-validated UX/API/data observations produced by Codex browser or Chrome tools.

The repository includes greenfield web, mature monolith, and polyglot service reference projects. `npm run test:e2e` drives Chromium through the greenfield UI/API/persistence journey; the reference matrix verifies four-layer coverage and uncertainty-honest fallback for unsupported languages.

## Packaging and Marketplace

`scripts/package-plugin.sh 0.1.0` creates self-contained macOS/Linux bundles with verified binaries, SPDX SBOM and SHA-256 checksums. `scripts/release-smoke.sh` exercises package, explicit and implicit routing, update and uninstall-state preservation. A repo marketplace catalog lives at `.agents/plugins/marketplace.json`; add it with `codex plugin marketplace add owner/repo` or link the local repository while developing. Public directory submission is completed through the OpenAI plugin submission portal and requires verified publisher identity, Apps Management write access, listing/support/privacy/terms URLs, five positive and three negative cases, release notes, and the final skills bundle.

## Secret-Safe Execution with tene

tene-codex and `tene-workflow` do not own, decrypt, or store secret values. Secret metadata and runtime injection are delegated to the separate [`tene`](https://github.com/agent-kay-it/tene) CLI.

The intended security boundary is:

- never read `.tene/**` vault files;
- never call `tene get` from agent automation;
- never request secret values in the conversation;
- inspect names or capability metadata only;
- run secret-dependent tests through `tene run --env <environment> -- <command>`;
- redact and scan QA artifacts before preserving them as evidence.

The plugin can operate without the tene secret CLI when a workflow does not require secrets. When secrets are required, failure is closed rather than bypassing the tene boundary.

## Documentation

- [`docs/00-rnd`](docs/00-rnd/README.md) — market and technical research
- [`docs/00-prd`](docs/00-prd/README.md) — product requirements and architecture requirements
- [`docs/01-plan`](docs/01-plan/README.md) — implementation roadmap and traceability
- [`docs/02-design`](docs/02-design/README.md) — implementation-level technical design

## Planned Repository Structure

```text
.codex-plugin/       Codex plugin manifest
skills/              Codex workflow skills
hooks/               Optional lifecycle enforcement
cmd/tene-workflow/   Go CLI entry point
internal/            Workflow engine implementation
schemas/             Persisted-state and document schemas
templates/           Sprint document templates
evals/               Skill routing and behavioral evaluations
docs/                Research, requirements, plans, and designs
```

## Contributing

The project is not yet accepting stability guarantees. Design review, threat modeling, workflow fixtures, QA scenarios, and implementation contributions are welcome as the public API is established. Contributions are expected to preserve the sprint-state invariants, evidence-based QA model, and secret boundary described in the design documents.

Unless explicitly stated otherwise, contributions submitted to this repository are licensed under the same Apache License 2.0 terms as the project.

## License

tene-codex, including the planned `tene-workflow` CLI and Codex plugin components in this repository, is licensed under the **Apache License, Version 2.0**. See [`LICENSE`](LICENSE) and [`NOTICE`](NOTICE).

Apache-2.0 permits commercial use, modification, private use, and redistribution. Distributions must preserve the applicable copyright, license, and NOTICE information, and modified files must carry notices describing significant changes. The license also includes an express patent grant and does not grant trademark rights.

This summary is informational and is not legal advice. The terms in `LICENSE` control.
