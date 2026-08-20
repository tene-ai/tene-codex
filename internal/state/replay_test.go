package state

import (
	"encoding/json"
	"github.com/tene-ai/tene-codex/internal/domain"
	"os"
	"sync"
	"testing"
	"time"
)

func TestReplayPatchAndRepairCorruptProjection(t *testing.T) {
	s := New(t.TempDir())
	p := domain.NewProject("project_test", "before", "strict", time.Now())
	if err := s.Initialize(p); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Mutate(nil, domain.Actor{Kind: "test"}, "NameChanged", p.ProjectID, map[string]string{"name": "after"}, func(p *domain.Project) error { p.Name = "after"; return nil }); err != nil {
		t.Fatal(err)
	}
	replayed, err := s.Replay()
	if err != nil || replayed.Name != "after" || replayed.Revision != 1 {
		t.Fatalf("%#v %v", replayed, err)
	}
	if err := os.WriteFile(s.ProjectPath(), []byte("broken"), 0644); err != nil {
		t.Fatal(err)
	}
	paths, err := s.RepairFromJournal()
	if err != nil || len(paths) == 0 {
		t.Fatalf("%v %v", paths, err)
	}
	got, err := s.Load()
	if err != nil || got.Name != "after" {
		t.Fatalf("%#v %v", got, err)
	}
}

func TestMergePatchRoundTrip(t *testing.T) {
	before := map[string]any{"a": map[string]any{"x": float64(1), "gone": true}, "list": []any{float64(1)}}
	after := map[string]any{"a": map[string]any{"x": float64(2)}, "list": []any{float64(2), float64(3)}}
	got := applyMergePatch(before, mergePatch(before, after))
	b, _ := json.Marshal(got)
	want, _ := json.Marshal(after)
	if string(b) != string(want) {
		t.Fatalf("%s != %s", b, want)
	}
}

func TestConcurrentExpectedRevisionOneCommit(t *testing.T) {
	s := New(t.TempDir())
	p := domain.NewProject("project_test", "x", "strict", time.Now())
	if err := s.Initialize(p); err != nil {
		t.Fatal(err)
	}
	expected := uint64(0)
	var wg sync.WaitGroup
	wg.Add(2)
	results := make(chan error, 2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			_, err := s.Mutate(&expected, domain.Actor{Kind: "test"}, "Race", p.ProjectID, nil, func(p *domain.Project) error { p.Name = "winner"; return nil })
			results <- err
		}()
	}
	wg.Wait()
	close(results)
	success := 0
	for err := range results {
		if err == nil {
			success++
		}
	}
	if success != 1 {
		t.Fatalf("success=%d", success)
	}
}
