package codeintel

import (
	"context"
	"path/filepath"
	"testing"
)

func TestReferenceProjectMatrix(t *testing.T) {
	root := filepath.Join("..", "..", "testdata")
	tests := []struct {
		name, path string
		minFiles   int
		degraded   bool
	}{{"greenfield-web", "reference-web", 1, false}, {"mature-monolith", "reference-mature", 4, false}, {"polyglot-services", "reference-polyglot", 4, true}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, e := Analyze(context.Background(), filepath.Join(root, tt.path), nil, false)
			if e != nil {
				t.Fatal(e)
			}
			if len(r.Files) < tt.minFiles {
				t.Fatalf("files=%d report=%+v", len(r.Files), r)
			}
			layers := map[string]bool{}
			degraded := false
			for _, f := range r.Files {
				layers[f.Layer] = true
			}
			for _, c := range r.Components {
				if c.Provider == "filesystem" && len(c.Unknown) > 0 {
					degraded = true
				}
			}
			if tt.name == "mature-monolith" {
				for _, l := range []string{"Interface", "Business Logic", "Persistence", "Infrastructure"} {
					if !layers[l] {
						t.Errorf("missing layer %s: %+v", l, r.Files)
					}
				}
			}
			if degraded != tt.degraded {
				t.Fatalf("degraded=%v want %v", degraded, tt.degraded)
			}
		})
	}
}
