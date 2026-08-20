// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/tene-ai/tene-codex/internal/domain"
)

type eventEnvelope struct {
	Data            any             `json:"data,omitempty"`
	ProjectionPatch any             `json:"projection_patch,omitempty"`
	Projection      *domain.Project `json:"projection,omitempty"`
}
type ProjectionDrift struct {
	Path        string   `json:"path"`
	Status      string   `json:"status"`
	Differences []string `json:"differences,omitempty"`
}

func mergePatch(before, after any) any {
	bm, bok := before.(map[string]any)
	am, aok := after.(map[string]any)
	if !bok || !aok {
		if equalJSON(before, after) {
			return map[string]any{}
		}
		return after
	}
	out := map[string]any{}
	for k := range bm {
		if _, ok := am[k]; !ok {
			out[k] = nil
		}
	}
	for k, av := range am {
		bv, ok := bm[k]
		if !ok {
			out[k] = av
			continue
		}
		if equalJSON(bv, av) {
			continue
		}
		if _, ok := bv.(map[string]any); ok {
			if _, ok := av.(map[string]any); ok {
				p := mergePatch(bv, av)
				if m, ok := p.(map[string]any); !ok || len(m) > 0 {
					out[k] = p
				}
				continue
			}
		}
		out[k] = av
	}
	return out
}
func applyMergePatch(target, patch any) any {
	pm, ok := patch.(map[string]any)
	if !ok {
		return patch
	}
	tm, ok := target.(map[string]any)
	if !ok {
		tm = map[string]any{}
	}
	out := map[string]any{}
	for k, v := range tm {
		out[k] = v
	}
	for k, v := range pm {
		if v == nil {
			delete(out, k)
		} else {
			out[k] = applyMergePatch(out[k], v)
		}
	}
	return out
}
func equalJSON(a, b any) bool {
	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)
	return bytes.Equal(ab, bb)
}
func jsonTree(v any) any {
	b, _ := json.Marshal(v)
	var out any
	_ = json.Unmarshal(b, &out)
	return out
}

func (s *Store) Replay() (*domain.Project, error) {
	events, err := s.VerifyJournal()
	if err != nil {
		return nil, err
	}
	var tree any
	start := -1
	for i := len(events) - 1; i >= 0; i-- {
		var env eventEnvelope
		b, _ := json.Marshal(events[i].Payload)
		if json.Unmarshal(b, &env) == nil && env.Projection != nil {
			tree = jsonTree(env.Projection)
			start = i
			break
		}
	}
	if start < 0 {
		return nil, fmt.Errorf("%w: no projection checkpoint; run compact before repair", ErrCorrupt)
	}
	for i := start + 1; i < len(events); i++ {
		var env eventEnvelope
		b, _ := json.Marshal(events[i].Payload)
		if err := json.Unmarshal(b, &env); err != nil || env.ProjectionPatch == nil {
			return nil, fmt.Errorf("%w: event %d has no projection patch", ErrCorrupt, events[i].Sequence)
		}
		tree = applyMergePatch(tree, env.ProjectionPatch)
	}
	b, err := json.Marshal(tree)
	if err != nil {
		return nil, err
	}
	var p domain.Project
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&p); err != nil {
		return nil, fmt.Errorf("%w: replay decode: %v", ErrCorrupt, err)
	}
	ensureMaps(&p)
	return &p, nil
}

func (s *Store) CreateCheckpoint() (string, *domain.Project, error) {
	var path string
	var result *domain.Project
	err := s.withLock(func() error {
		p, err := s.Load()
		if err != nil {
			return err
		}
		prev, seq, err := s.lastEvent()
		if err != nil {
			return err
		}
		old := p.Revision
		p.Revision++
		p.UpdatedAt = s.Now().UTC()
		e := domain.Event{Sequence: seq + 1, EventID: domain.NewID("event"), EventType: "ProjectionCheckpoint", AggregateID: p.ProjectID, OccurredAt: p.UpdatedAt, Actor: domain.Actor{Kind: "system"}, ExpectedRevision: old, Payload: eventEnvelope{Projection: p}, PreviousHash: prev}
		e.Hash = eventHash(e)
		if err := appendEvent(s.EventsPath(), e); err != nil {
			return err
		}
		if err := atomicJSON(s.ProjectPath(), p); err != nil {
			return err
		}
		if err := atomicJSON(s.ActivePath(), p); err != nil {
			return err
		}
		if err := atomicJSON(s.MasterPlanPath(), masterPlan(p)); err != nil {
			return err
		}
		path = filepath.Join(s.Dir(), "backups", fmt.Sprintf("checkpoint-%020d.json", p.Revision))
		if err := atomicJSON(path, p); err != nil {
			return err
		}
		result = p
		return nil
	})
	return path, result, err
}

func (s *Store) ProjectionDrift() ([]ProjectionDrift, *domain.Project, error) {
	p, err := s.Replay()
	if err != nil {
		return nil, nil, err
	}
	checks := []struct {
		path  string
		value any
	}{{s.ProjectPath(), p}, {s.ActivePath(), p}, {s.MasterPlanPath(), masterPlan(p)}}
	var out []ProjectionDrift
	for _, c := range checks {
		expected, _ := json.MarshalIndent(c.value, "", "  ")
		expected = append(expected, '\n')
		actual, err := os.ReadFile(c.path)
		status := "match"
		if os.IsNotExist(err) {
			status = "missing"
		} else if err != nil {
			status = "unreadable"
		} else if !bytes.Equal(expected, actual) {
			var expectedTree, actualTree any
			if json.Unmarshal(expected, &expectedTree) != nil || json.Unmarshal(actual, &actualTree) != nil || !equalJSON(expectedTree, actualTree) {
				status = "drift"
			}
		}
		differences := []string{}
		if status == "drift" {
			var expectedTree, actualTree any
			_ = json.Unmarshal(expected, &expectedTree)
			_ = json.Unmarshal(actual, &actualTree)
			differences = diffPaths("$", expectedTree, actualTree, 20)
		}
		out = append(out, ProjectionDrift{Path: c.path, Status: status, Differences: differences})
	}
	return out, p, nil
}

func diffPaths(path string, expected, actual any, limit int) []string {
	if limit <= 0 || equalJSON(expected, actual) {
		return nil
	}
	em, eok := expected.(map[string]any)
	am, aok := actual.(map[string]any)
	if eok && aok {
		keys := map[string]bool{}
		for k := range em {
			keys[k] = true
		}
		for k := range am {
			keys[k] = true
		}
		names := make([]string, 0, len(keys))
		for k := range keys {
			names = append(names, k)
		}
		sort.Strings(names)
		var out []string
		for _, k := range names {
			out = append(out, diffPaths(path+"."+k, em[k], am[k], limit-len(out))...)
			if len(out) >= limit {
				break
			}
		}
		return out
	}
	return []string{path}
}

func (s *Store) RepairFromJournal() ([]string, error) {
	drift, p, err := s.ProjectionDrift()
	if err != nil {
		return nil, err
	}
	var repaired []string
	err = s.withLock(func() error {
		stamp := time.Now().UTC().Format("20060102T150405.000000000Z")
		for _, d := range drift {
			if d.Status == "match" {
				continue
			}
			if b, e := os.ReadFile(d.Path); e == nil {
				backup := filepath.Join(s.Dir(), "backups", filepath.Base(d.Path)+"."+stamp+".bak")
				if e := atomicBytes(backup, b); e != nil {
					return e
				}
			}
			var value any = p
			if d.Path == s.MasterPlanPath() {
				value = masterPlan(p)
			}
			if e := atomicJSON(d.Path, value); e != nil {
				return e
			}
			repaired = append(repaired, d.Path)
		}
		return nil
	})
	return repaired, err
}
