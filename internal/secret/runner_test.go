package secret

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func fake(t *testing.T, body string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip()
	}
	p := filepath.Join(t.TempDir(), "tene")
	if err := os.WriteFile(p, []byte("#!/bin/sh\n"+body+"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	return p
}
func TestMissing(t *testing.T) {
	_, e := checkPath(func(string) (string, error) { return "", errors.New("no") })
	if e == nil || !strings.Contains(e.Error(), "SEC_TENE_MISSING") {
		t.Fatal(e)
	}
}
func TestListNamesStripsPreviewAndSorts(t *testing.T) {
	p := fake(t, `printf '%s' '{"environment":"test","project":"private","secrets":[{"name":"Z_KEY","preview":"canary"},{"name":"A_KEY","preview":"…1234"}]}'`)
	n, e := listNamesPath(context.Background(), p, "test")
	if e != nil || n.Count != 2 || n.Names[0].Name != "A_KEY" {
		t.Fatalf("%+v %v", n, e)
	}
	b := n.Names[0].Name + n.Names[1].Name
	if strings.Contains(b, "1234") || strings.Contains(b, "canary") {
		t.Fatal("preview escaped")
	}
}
func TestForbiddenAndSensitiveArgs(t *testing.T) {
	p := fake(t, "exit 0")
	for _, c := range [][]string{{"env"}, {"sh", "-c", "ok"}, {"tool", "token=plaintext"}, {"tool", "ghp_12345678901234567890"}} {
		if _, e := runPath(context.Background(), p, "test", c, nil); e == nil {
			t.Fatalf("allowed %v", c)
		}
	}
}
func TestCanaryLeakQuarantines(t *testing.T) {
	p := fake(t, `printf '%s' 'TENE_TEST_CANARY_0123456789'`)
	r, e := runPath(context.Background(), p, "test", []string{"go", "test"}, []string{"TENE_TEST_CANARY_0123456789"})
	if e == nil || !r.Quarantined || r.Stdout != "" {
		t.Fatalf("%+v %v", r, e)
	}
}
func TestChildFailurePreservesSafeCode(t *testing.T) {
	p := fake(t, `echo 'ordinary failure' >&2; exit 23`)
	r, e := runPath(context.Background(), p, "test", []string{"go", "test"}, nil)
	if e == nil || r.ExitCode != 23 || !strings.Contains(r.Stderr, "ordinary failure") {
		t.Fatalf("%+v %v", r, e)
	}
}
func TestPermissionFailureIsSanitized(t *testing.T) {
	p := fake(t, `echo 'permission denied: preview=should-not-escape' >&2; exit 7`)
	_, e := listNamesPath(context.Background(), p, "test")
	if e == nil || strings.Contains(e.Error(), "should-not-escape") {
		t.Fatal(e)
	}
}
