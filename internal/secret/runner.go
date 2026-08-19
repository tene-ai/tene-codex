// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var forbidden = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^(env|printenv)$`),
	regexp.MustCompile(`(?i)(^|/)(sh|bash|zsh)$`),
	regexp.MustCompile(`(?i)(dump.?env|show.?env|process\.env)`),
}

type Result struct {
	ExitCode    int      `json:"exit_code"`
	Stdout      string   `json:"stdout,omitempty"`
	Stderr      string   `json:"stderr,omitempty"`
	Environment string   `json:"environment"`
	Command     []string `json:"command"`
}

func Check() (string, error) {
	path, err := exec.LookPath("tene")
	if err != nil {
		return "", fmt.Errorf("SEC_TENE_MISSING: install the tene CLI before secret-required execution")
	}
	return path, nil
}

func Run(ctx context.Context, environment string, command []string) (Result, error) {
	if environment == "" || len(command) == 0 {
		return Result{}, fmt.Errorf("SEC_INVALID_REQUEST: environment and command are required")
	}
	joined := strings.Join(command, " ")
	for _, re := range forbidden {
		if re.MatchString(command[0]) || re.MatchString(joined) {
			return Result{}, fmt.Errorf("SEC_FORBIDDEN_COMMAND: commands that expose the child environment or invoke a shell are not allowed")
		}
	}
	path, err := Check()
	if err != nil {
		return Result{}, err
	}
	args := append([]string{"run", "--env", environment, "--"}, command...)
	cmd := exec.CommandContext(ctx, path, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err = cmd.Run()
	code := 0
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			code = exit.ExitCode()
		} else {
			return Result{}, err
		}
	}
	return Result{ExitCode: code, Stdout: sanitize(stdout.String()), Stderr: sanitize(stderr.String()), Environment: environment, Command: command}, err
}

func ListNames(ctx context.Context, environment string) (any, error) {
	path, err := Check()
	if err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(ctx, path, "list", "--env", environment, "--json")
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("SEC_TENE_LIST_FAILED: %w", err)
	}
	var value any
	if err := json.Unmarshal(b, &value); err != nil {
		return nil, fmt.Errorf("SEC_TENE_OUTPUT_INVALID: %w", err)
	}
	return value, nil
}

func sanitize(value string) string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(api[_-]?key|token|password|secret)\s*[:=]\s*[^\s]+`),
		regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9._~+/-]+=*`),
	}
	for _, re := range patterns {
		value = re.ReplaceAllString(value, "$1=[REDACTED]")
	}
	return value
}
