package tracecontext

import (
	"github.com/tene-ai/tene-codex/internal/domain"
	"testing"
)

func TestImpactIsDirectedBoundedAndFindsAC(t *testing.T) {
	g := domain.Graph{Nodes: map[string]domain.Node{"symbol:a": {ID: "symbol:a", Kind: "Symbol"}, "symbol:b": {ID: "symbol:b", Kind: "Symbol"}, "ac:x": {ID: "ac:x", Kind: "AcceptanceCriterion"}}, Edges: map[string]domain.Edge{"e1": {ID: "e1", From: "symbol:a", To: "symbol:b", Kind: "calls"}, "e2": {ID: "e2", From: "symbol:b", To: "ac:x", Kind: "realizes"}, "e3": {ID: "e3", From: "ac:x", To: "symbol:a", Kind: "realizes"}}}
	r, err := Impact(g, "symbol:a", 4, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.ImpactedACIDs) != 1 || r.ImpactedACIDs[0] != "ac:x" || len(r.Paths) != 2 {
		t.Fatalf("%#v", r)
	}
	r, err = Impact(g, "symbol:a", 4, 0)
	if err != nil || !r.Truncated || len(r.ImpactedACIDs) != 0 {
		t.Fatalf("%#v %v", r, err)
	}
}

func TestValidateGraphFindsDanglingAndOrphans(t *testing.T) {
	p := domain.NewProject("project_x", "x", "strict", testTime())
	p.Intents["intent_x"] = &domain.Intent{IntentID: "intent_x", Status: "confirmed"}
	p.Tasks["task_x"] = &domain.Task{TaskID: "task_x"}
	p.Evidence["evidence_x"] = &domain.Evidence{EvidenceID: "evidence_x"}
	p.Graph.Edges["edge_x"] = domain.Edge{ID: "edge_x", From: "missing", To: "also-missing"}
	f := ValidateGraph(p)
	codes := map[string]bool{}
	for _, x := range f {
		codes[x.Code] = true
	}
	for _, code := range []string{"GRAPH_EDGE_DANGLING", "GRAPH_INTENT_ORPHAN", "GRAPH_TASK_ORPHAN", "GRAPH_EVIDENCE_ORPHAN"} {
		if !codes[code] {
			t.Fatalf("missing %s in %#v", code, f)
		}
	}
}
