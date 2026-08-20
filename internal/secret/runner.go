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
	"sort"
	"strings"
)

var forbidden = []*regexp.Regexp{regexp.MustCompile(`(?i)^(env|printenv)$`), regexp.MustCompile(`(?i)(^|/)(sh|bash|zsh)$`), regexp.MustCompile(`(?i)(dump.?env|show.?env|process\.env|/proc/.*/environ)`)}
var sensitiveArgument = []*regexp.Regexp{regexp.MustCompile(`(?i)^(?:sk|ghp|github_pat|xox[baprs]|AKIA)[-_A-Za-z0-9]{12,}$`), regexp.MustCompile(`(?i)^(?:api[_-]?key|token|password|secret)=.+$`)}
var leakage = []*regexp.Regexp{regexp.MustCompile(`(?i)(api[_-]?key|token|password|secret)\s*[:=]\s*[^\s]+`), regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9._~+/-]+=*`), regexp.MustCompile(`(?i)TENE_(?:TEST_)?CANARY[-_A-Za-z0-9]{8,}`), regexp.MustCompile(`(?i)(?:sk|ghp|github_pat|xox[baprs])[-_A-Za-z0-9]{16,}`)}

type Result struct {
	ExitCode     int      `json:"exit_code"`
	Stdout       string   `json:"stdout,omitempty"`
	Stderr       string   `json:"stderr,omitempty"`
	Environment  string   `json:"environment"`
	Command      []string `json:"command"`
	LeakDetected bool     `json:"leak_detected"`
	Quarantined  bool     `json:"quarantined"`
}
type SecretMeta struct {
	Name string `json:"name"`
}
type Names struct {
	Environment string       `json:"environment"`
	Names       []SecretMeta `json:"secrets"`
	Count       int          `json:"count"`
}

func Check() (string, error) { return checkPath(exec.LookPath) }
func checkPath(lookup func(string) (string, error)) (string, error) {
	p, e := lookup("tene")
	if e != nil {
		return "", fmt.Errorf("SEC_TENE_MISSING: install the tene CLI before secret-required execution")
	}
	return p, nil
}
func Run(ctx context.Context, environment string, command []string) (Result, error) {
	p, e := Check()
	if e != nil {
		return Result{}, e
	}
	return runPath(ctx, p, environment, command, nil)
}
func runPath(ctx context.Context, path, environment string, command, canaries []string) (Result, error) {
	if environment == "" || len(command) == 0 {
		return Result{}, fmt.Errorf("SEC_INVALID_REQUEST: environment and command are required")
	}
	joined := strings.Join(command, " ")
	for _, re := range forbidden {
		if re.MatchString(command[0]) || re.MatchString(joined) {
			return Result{}, fmt.Errorf("SEC_FORBIDDEN_COMMAND: commands that expose the child environment or invoke a shell are not allowed")
		}
	}
	for _, arg := range command {
		for _, re := range sensitiveArgument {
			if re.MatchString(arg) {
				return Result{}, fmt.Errorf("SEC_SENSITIVE_ARGUMENT: pass secret values only through the tene child environment")
			}
		}
	}
	cmd := exec.CommandContext(ctx, path, append([]string{"run", "--env", environment, "--"}, command...)...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	runErr := cmd.Run()
	code := 0
	if runErr != nil {
		if ex, ok := runErr.(*exec.ExitError); ok {
			code = ex.ExitCode()
		} else {
			return Result{}, fmt.Errorf("SEC_TENE_EXEC_FAILED: %w", runErr)
		}
	}
	stdout, stderr := out.String(), errOut.String()
	if detectsLeak(stdout+"\n"+stderr, canaries) {
		return Result{ExitCode: code, Environment: environment, Command: append([]string(nil), command...), LeakDetected: true, Quarantined: true}, fmt.Errorf("SEC_OUTPUT_LEAK: child output was quarantined because it matched secret leakage policy")
	}
	return Result{ExitCode: code, Stdout: sanitize(stdout), Stderr: sanitize(stderr), Environment: environment, Command: append([]string(nil), command...)}, runErr
}
func ListNames(ctx context.Context, environment string) (Names, error) {
	p, e := Check()
	if e != nil {
		return Names{}, e
	}
	return listNamesPath(ctx, p, environment)
}
func listNamesPath(ctx context.Context, path, environment string) (Names, error) {
	cmd := exec.CommandContext(ctx, path, "list", "--env", environment, "--json")
	b, e := cmd.Output()
	if e != nil {
		return Names{}, fmt.Errorf("SEC_TENE_LIST_FAILED: metadata permission denied or unavailable")
	}
	var raw struct {
		Environment string `json:"environment"`
		Secrets     []struct {
			Name string `json:"name"`
		} `json:"secrets"`
	}
	if e = json.Unmarshal(b, &raw); e != nil {
		return Names{}, fmt.Errorf("SEC_TENE_OUTPUT_INVALID: %w", e)
	}
	n := Names{Environment: raw.Environment}
	if n.Environment == "" {
		n.Environment = environment
	}
	for _, x := range raw.Secrets {
		if x.Name != "" {
			n.Names = append(n.Names, SecretMeta{Name: x.Name})
		}
	}
	sort.Slice(n.Names, func(i, j int) bool { return n.Names[i].Name < n.Names[j].Name })
	n.Count = len(n.Names)
	return n, nil
}
func detectsLeak(v string, canaries []string) bool {
	for _, c := range canaries {
		if c != "" && strings.Contains(v, c) {
			return true
		}
	}
	for _, re := range leakage {
		if re.MatchString(v) {
			return true
		}
	}
	return false
}

// DetectLeak applies the same fail-closed artifact policy used by the runner.
func DetectLeak(value []byte) bool { return detectsLeak(string(value), nil) }
func sanitize(v string) string {
	for _, re := range leakage {
		v = re.ReplaceAllString(v, "[REDACTED]")
	}
	return v
}
