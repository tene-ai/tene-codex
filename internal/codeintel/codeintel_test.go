package codeintel

import (
	"context"
	"os"
	"path/filepath"
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
