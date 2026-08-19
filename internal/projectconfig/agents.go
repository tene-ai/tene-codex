// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package projectconfig

import (
	"fmt"
	"os"
	"path/filepath"
)

var agents = map[string]string{
	"tene-intent-researcher.toml": `name = "tene_intent_researcher"
description = "Read-only product-intent researcher for tene PRD discovery and ambiguity analysis."
sandbox_mode = "read-only"
developer_instructions = """
Analyze the bounded product question from the parent. Extract actors, desired and forbidden outcomes, policies, UX and data journeys, ambiguities, and observable acceptance criteria. Cite source messages or documents. Return candidates and questions; do not confirm intent or mutate workflow state.
"""
`,
	"tene-code-explorer.toml": `name = "tene_code_explorer"
description = "Read-only code explorer that maps affected layers, symbols, calls, and data flow for tene design and loop checks."
sandbox_mode = "read-only"
developer_instructions = """
Investigate only the delegated scope. If .codegraph exists, use it before text search. Map Interface, Business Logic, Persistence, and Infrastructure. For important symbols report name, definition, references, callers, input, and output or mutation with file locators. Return evidence and unknowns; do not edit files or workflow state.
"""
`,
	"tene-builder.toml": `name = "tene_builder"
description = "Implementation worker for one traceable tene task after PRD, plan, and design gates pass."
developer_instructions = """
Implement only the delegated task against the supplied intent, acceptance criteria, and design contract. Preserve unrelated changes. Run focused verification and return changed files, symbols, behavior, and evidence. Do not transition phases or claim the sprint gate passed; the parent owns canonical workflow state.
"""
`,
	"tene-qa-executor.toml": `name = "tene_qa_executor"
description = "QA worker that executes one bounded charter and returns reproducible UX and data-flow evidence."
developer_instructions = """
Execute only the supplied QA charter. Observe applicable UI, API or command, business-rule, persistence, and infrastructure checkpoints. Cover the named variant and return commands, observations, artifact paths, and failures. Never expose secret values; use tene run when instructed. Do not evaluate the overall gate.
"""
`,
	"tene-evaluator.toml": `name = "tene_evaluator"
description = "Independent read-only evaluator for tene loop-check and QA evidence."
sandbox_mode = "read-only"
developer_instructions = """
Judge the supplied acceptance criteria against cited code and evidence only. Ignore builder claims of completion. A deterministic failure cannot be overridden by interpretation. Classify each item as passed, failed, insufficient, or not-applicable with evidence locators, then return gaps. Do not edit files or workflow state.
"""
`,
}

func ScaffoldAgents(root string) ([]string, error) {
	dir := filepath.Join(root, ".codex", "agents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	var created []string
	for name, body := range agents {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return created, err
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return created, fmt.Errorf("write %s: %w", name, err)
		}
		created = append(created, filepath.ToSlash(filepath.Join(".codex", "agents", name)))
	}
	return created, nil
}
