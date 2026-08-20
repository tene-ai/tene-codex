package codeintel

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAnalyzeGoSixQuestions(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "internal", "service"), 0755); err != nil {
		t.Fatal(err)
	}
	src := `package service
import "os"
func Save(name string, values []int) (int, error) { b := []byte(name); err := os.WriteFile(name,b,0600); return len(values),err }
`
	if err := os.WriteFile(filepath.Join(root, "internal", "service", "save.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	r, err := Analyze(context.Background(), root, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Components) != 1 {
		t.Fatalf("%#v", r)
	}
	c := r.Components[0]
	if c.Name != "Save" || c.PrimaryLayer != "Business Logic" || len(c.Inputs) != 2 || len(c.Outputs) != 2 || len(c.Calls) == 0 || len(c.Effects) == 0 || len(c.Unknown) == 0 {
		t.Fatalf("incomplete six questions: %#v", c)
	}
}
func TestRejectsOutsidePath(t *testing.T) {
	_, err := Analyze(context.Background(), t.TempDir(), []string{"../secret"}, false)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPolyglotProvidersKeepBoundedUncertainty(t *testing.T) {
	root := t.TempDir()
	for path, body := range map[string]string{"web/ui.ts": "export const submit = () => fetch('/api')", "worker/job.py": "def consume(message): pass", ".tene/vault.py": "SECRET='never index'"} {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	r, err := Analyze(context.Background(), root, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Files) != 2 || len(r.Components) != 2 {
		t.Fatalf("files=%+v components=%+v", r.Files, r.Components)
	}
	for _, c := range r.Components {
		if (c.Provider != "typescript-static" && c.Provider != "python-static") || len(c.Unknown) == 0 || len(c.Inputs) == 0 || len(c.Outputs) == 0 || len(c.Effects) == 0 {
			t.Fatalf("not uncertainty-honest: %+v", c)
		}
		if strings.Contains(c.File, ".tene") {
			t.Fatal("vault indexed")
		}
	}
}
