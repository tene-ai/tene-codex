// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/tene-ai/tene-codex/internal/domain"
)

const DirName = ".tene-workflow"

var (
	ErrNotInitialized = errors.New("STATE_NOT_INITIALIZED")
	ErrConflict       = errors.New("STATE_REVISION_CONFLICT")
	ErrLocked         = errors.New("STATE_LOCKED")
	ErrCorrupt        = errors.New("STATE_CORRUPT")
)

type Store struct {
	Root string
	Now  func() time.Time
}

func New(root string) *Store { return &Store{Root: root, Now: time.Now} }

func (s *Store) Dir() string            { return filepath.Join(s.Root, DirName) }
func (s *Store) ProjectPath() string    { return filepath.Join(s.Dir(), "project.json") }
func (s *Store) ActivePath() string     { return filepath.Join(s.Dir(), "active.json") }
func (s *Store) EventsPath() string     { return filepath.Join(s.Dir(), "events.ndjson") }
func (s *Store) LockPath() string       { return filepath.Join(s.Dir(), ".lock") }
func (s *Store) MasterPlanPath() string { return filepath.Join(s.Dir(), "master-plan.json") }

func (s *Store) Exists() bool {
	_, err := os.Stat(s.ProjectPath())
	return err == nil
}

func (s *Store) Initialize(project *domain.Project) error {
	if s.Exists() {
		return fmt.Errorf("project already initialized")
	}
	for _, dir := range []string{s.Dir(), filepath.Join(s.Dir(), "graph"), filepath.Join(s.Dir(), "evidence"), filepath.Join(s.Dir(), "cache"), filepath.Join(s.Dir(), "backups")} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	if err := s.withLock(func() error {
		if err := atomicJSON(s.ProjectPath(), project); err != nil {
			return err
		}
		if err := atomicJSON(s.ActivePath(), project); err != nil {
			return err
		}
		if err := atomicJSON(s.MasterPlanPath(), masterPlan(project)); err != nil {
			return err
		}
		if err := atomicBytes(filepath.Join(s.Dir(), "policies.yaml"), []byte("schema_version: 1.0.0\nworkflow_profile: "+project.Profile+"\nqa:\n  blocking_coverage: 100\nsecrets:\n  provider: tene\n")); err != nil {
			return err
		}
		e := domain.Event{Sequence: 1, EventID: domain.NewID("event"), EventType: "ProjectInitialized", AggregateID: project.ProjectID, OccurredAt: s.Now().UTC(), Actor: domain.Actor{Kind: "user"}, ExpectedRevision: 0, Payload: eventEnvelope{Projection: project}}
		e.Hash = eventHash(e)
		return appendEvent(s.EventsPath(), e)
	}); err != nil {
		return err
	}
	return nil
}

func (s *Store) Load() (*domain.Project, error) {
	b, err := os.ReadFile(s.ProjectPath())
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotInitialized
	}
	if err != nil {
		return nil, err
	}
	var p domain.Project
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&p); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCorrupt, err)
	}
	ensureMaps(&p)
	return &p, nil
}

type MigrationPlan struct {
	From      string   `json:"from"`
	To        string   `json:"to"`
	Required  bool     `json:"required"`
	Supported bool     `json:"supported"`
	Changes   []string `json:"changes"`
	Backup    string   `json:"backup,omitempty"`
}

func (s *Store) PlanMigration() (MigrationPlan, error) {
	b, err := os.ReadFile(s.ProjectPath())
	if errors.Is(err, os.ErrNotExist) {
		return MigrationPlan{}, ErrNotInitialized
	}
	if err != nil {
		return MigrationPlan{}, err
	}
	var header struct {
		SchemaVersion string `json:"schema_version"`
	}
	if err := json.Unmarshal(b, &header); err != nil {
		return MigrationPlan{}, fmt.Errorf("%w: %v", ErrCorrupt, err)
	}
	plan := MigrationPlan{From: header.SchemaVersion, To: domain.SchemaVersion, Changes: []string{}}
	switch header.SchemaVersion {
	case domain.SchemaVersion:
		plan.Supported = true
	case "0.9.0":
		plan.Required = true
		plan.Supported = true
		plan.Changes = []string{"set schema_version to 1.0.0", "initialize waivers and request_results maps", "regenerate active and master-plan projections"}
	default:
		plan.Required = header.SchemaVersion != domain.SchemaVersion
		plan.Supported = false
	}
	return plan, nil
}

func (s *Store) Migrate() (MigrationPlan, error) {
	plan, err := s.PlanMigration()
	if err != nil {
		return plan, err
	}
	if !plan.Supported {
		return plan, fmt.Errorf("STATE_MIGRATION_UNSUPPORTED: %s", plan.From)
	}
	if !plan.Required {
		return plan, nil
	}
	err = s.withLock(func() error {
		b, err := os.ReadFile(s.ProjectPath())
		if err != nil {
			return err
		}
		backup := filepath.Join(s.Dir(), "backups", fmt.Sprintf("pre-migration-%s-%d.json", plan.From, time.Now().UTC().UnixNano()))
		if err := atomicBytes(backup, b); err != nil {
			return err
		}
		var p domain.Project
		if err := json.Unmarshal(b, &p); err != nil {
			return err
		}
		ensureMaps(&p)
		before := jsonTree(&p)
		p.SchemaVersion = domain.SchemaVersion
		prev, seq, err := s.lastEvent()
		if err != nil {
			return err
		}
		old := p.Revision
		p.Revision++
		p.UpdatedAt = s.Now().UTC()
		e := domain.Event{Sequence: seq + 1, EventID: domain.NewID("event"), EventType: "SchemaMigrated", AggregateID: p.ProjectID, OccurredAt: p.UpdatedAt, Actor: domain.Actor{Kind: "system"}, ExpectedRevision: old, Payload: eventEnvelope{Data: plan, ProjectionPatch: mergePatch(before, jsonTree(&p))}, PreviousHash: prev}
		e.Hash = eventHash(e)
		if err := appendEvent(s.EventsPath(), e); err != nil {
			return err
		}
		if err := atomicJSON(s.ProjectPath(), &p); err != nil {
			return err
		}
		if err := atomicJSON(s.ActivePath(), &p); err != nil {
			return err
		}
		if err := atomicJSON(s.MasterPlanPath(), masterPlan(&p)); err != nil {
			return err
		}
		plan.Backup = filepath.ToSlash(backup)
		return nil
	})
	return plan, err
}

func (s *Store) RepairDerived() ([]string, error) {
	return s.RepairFromJournal()
}

func (s *Store) Mutate(expected *uint64, actor domain.Actor, eventType, aggregateID string, payload any, fn func(*domain.Project) error) (*domain.Project, error) {
	var result *domain.Project
	err := s.withLock(func() error {
		p, err := s.Load()
		if err != nil {
			return err
		}
		if expected != nil && p.Revision != *expected {
			return fmt.Errorf("%w: expected %d, current %d", ErrConflict, *expected, p.Revision)
		}
		if actor.SessionID != "" {
			if cached, ok := p.RequestResults[actor.SessionID]; ok {
				if cached.CommandHash != actor.ID {
					return fmt.Errorf("%w: request ID reused for a different command", ErrConflict)
				}
				return fmt.Errorf("%w: request mutation already committed", ErrConflict)
			}
		}
		before := jsonTree(p)
		if err := fn(p); err != nil {
			return err
		}
		if actor.SessionID != "" {
			p.RequestResults[actor.SessionID] = domain.RequestResult{Revision: p.Revision + 1, CommandHash: actor.ID, Completed: false}
		}
		previousHash, sequence, err := s.lastEvent()
		if err != nil {
			return err
		}
		oldRevision := p.Revision
		p.Revision++
		p.UpdatedAt = s.Now().UTC()
		patch := mergePatch(before, jsonTree(p))
		e := domain.Event{Sequence: sequence + 1, EventID: domain.NewID("event"), EventType: eventType, AggregateID: aggregateID, OccurredAt: p.UpdatedAt, Actor: actor, ExpectedRevision: oldRevision, Payload: eventEnvelope{Data: payload, ProjectionPatch: patch}, PreviousHash: previousHash}
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
		result = p
		return nil
	})
	return result, err
}

func masterPlan(p *domain.Project) map[string]any {
	type item struct {
		SprintID     string       `json:"sprint_id"`
		Title        string       `json:"title"`
		Milestone    string       `json:"milestone,omitempty"`
		Release      string       `json:"release,omitempty"`
		Phase        domain.Phase `json:"phase"`
		Predecessors []string     `json:"predecessor_ids,omitempty"`
	}
	items := make([]item, 0, len(p.Sprints))
	for _, sprint := range p.Sprints {
		items = append(items, item{SprintID: sprint.SprintID, Title: sprint.Title, Milestone: sprint.Milestone, Release: sprint.Release, Phase: sprint.Phase, Predecessors: sprint.Predecessors})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].SprintID < items[j].SprintID })
	return map[string]any{"schema_version": domain.SchemaVersion, "project_id": p.ProjectID, "revision": p.Revision, "active_sprint_id": p.ActiveSprintID, "plan": p.MasterPlan, "sprints": items}
}

func (s *Store) VerifyJournal() ([]domain.Event, error) {
	f, err := os.Open(s.EventsPath())
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var events []domain.Event
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 8*1024*1024)
	prev := ""
	var seq uint64
	for scanner.Scan() {
		var e domain.Event
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			return nil, fmt.Errorf("%w: invalid event: %v", ErrCorrupt, err)
		}
		if e.Sequence != seq+1 || e.PreviousHash != prev || e.Hash != eventHash(e) {
			return nil, fmt.Errorf("%w: invalid chain at sequence %d", ErrCorrupt, e.Sequence)
		}
		seq, prev = e.Sequence, e.Hash
		events = append(events, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Store) ClearDerived() error {
	return s.withLock(func() error {
		for _, path := range []string{filepath.Join(s.Dir(), "cache"), filepath.Join(s.Dir(), "graph", "index.json")} {
			info, err := os.Stat(path)
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			if err != nil {
				return err
			}
			if info.IsDir() {
				entries, err := os.ReadDir(path)
				if err != nil {
					return err
				}
				for _, entry := range entries {
					if err := os.RemoveAll(filepath.Join(path, entry.Name())); err != nil {
						return err
					}
				}
			} else if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Compact() (string, error) {
	var out string
	err := s.withLock(func() error {
		p, err := s.Load()
		if err != nil {
			return err
		}
		name := fmt.Sprintf("snapshot-%020d.json", p.Revision)
		out = filepath.Join(s.Dir(), "backups", name)
		return atomicJSON(out, p)
	})
	return out, err
}

func (s *Store) lastEvent() (string, uint64, error) {
	f, err := os.Open(s.EventsPath())
	if errors.Is(err, os.ErrNotExist) {
		return "", 0, nil
	}
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 8*1024*1024)
	var last domain.Event
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &last); err != nil {
			return "", 0, err
		}
	}
	return last.Hash, last.Sequence, scanner.Err()
}

func (s *Store) withLock(fn func() error) error {
	if err := os.MkdirAll(s.Dir(), 0o755); err != nil {
		return err
	}
	deadline := time.Now().Add(5 * time.Second)
	for {
		f, err := os.OpenFile(s.LockPath(), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err == nil {
			_, _ = fmt.Fprintf(f, "%d\n%s\n", os.Getpid(), time.Now().UTC().Format(time.RFC3339Nano))
			_ = f.Close()
			defer os.Remove(s.LockPath())
			return fn()
		}
		if !errors.Is(err, os.ErrExist) {
			return err
		}
		if time.Now().After(deadline) {
			return ErrLocked
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func atomicJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tene-tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0o644); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(b); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		return err
	}
	d, err := os.Open(filepath.Dir(path))
	if err == nil {
		_ = d.Sync()
		_ = d.Close()
	}
	return nil
}

func appendEvent(path string, e domain.Event) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(append(b, '\n')); err != nil {
		return err
	}
	return f.Sync()
}

func eventHash(e domain.Event) string {
	e.Hash = ""
	b, _ := json.Marshal(e)
	// Normalize typed payloads through JSON so verification after decoding into
	// `any` produces the same canonical bytes.
	var normalized any
	_ = json.Unmarshal(b, &normalized)
	b, _ = json.Marshal(normalized)
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func ensureMaps(p *domain.Project) {
	if p.Sprints == nil {
		p.Sprints = map[string]*domain.Sprint{}
	}
	if p.Tasks == nil {
		p.Tasks = map[string]*domain.Task{}
	}
	if p.Intents == nil {
		p.Intents = map[string]*domain.Intent{}
	}
	if p.Criteria == nil {
		p.Criteria = map[string]*domain.Criterion{}
	}
	if p.Gaps == nil {
		p.Gaps = map[string]*domain.Gap{}
	}
	if p.Waivers == nil {
		p.Waivers = map[string]*domain.Waiver{}
	}
	if p.Approvals == nil {
		p.Approvals = map[string]*domain.Approval{}
	}
	for _, sprint := range p.Sprints {
		if sprint.MaxLoopIterations == 0 {
			sprint.MaxLoopIterations = 5
		}
	}
	if p.Evidence == nil {
		p.Evidence = map[string]*domain.Evidence{}
	}
	if p.QARuns == nil {
		p.QARuns = map[string]*domain.QARun{}
	}
	if p.Graph.Nodes == nil {
		p.Graph.Nodes = map[string]domain.Node{}
	}
	if p.Graph.Edges == nil {
		p.Graph.Edges = map[string]domain.Edge{}
	}
	if p.RequestResults == nil {
		p.RequestResults = map[string]domain.RequestResult{}
	}
}

func CopyFile(dst string, src io.Reader) error {
	b, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	return atomicBytes(dst, b)
}

func atomicBytes(path string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tene-tmp-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err := io.Copy(tmp, bytes.NewReader(b)); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(name, path)
}
