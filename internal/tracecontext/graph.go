// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package tracecontext

import (
	"fmt"
	"sort"

	"github.com/tene-ai/tene-codex/internal/domain"
)

type GraphPath struct {
	NodeIDs []string      `json:"node_ids"`
	Edges   []domain.Edge `json:"edges"`
}

type ImpactResult struct {
	StartID       string      `json:"start_id"`
	MaxDepth      int         `json:"max_depth"`
	CallDepth     int         `json:"call_depth"`
	Paths         []GraphPath `json:"paths"`
	ImpactedACIDs []string    `json:"impacted_ac_ids"`
	VisitedIDs    []string    `json:"visited_ids"`
	Truncated     bool        `json:"truncated"`
}

func Impact(g domain.Graph, start string, maxDepth, callDepth int) (ImpactResult, error) {
	if _, ok := g.Nodes[start]; !ok {
		return ImpactResult{}, fmt.Errorf("GRAPH_NODE_NOT_FOUND: %s", start)
	}
	if maxDepth < 1 || callDepth < 0 {
		return ImpactResult{}, fmt.Errorf("GRAPH_DEPTH_INVALID")
	}
	adj := map[string][]domain.Edge{}
	for _, edge := range g.Edges {
		adj[edge.From] = append(adj[edge.From], edge)
	}
	for id := range adj {
		sort.Slice(adj[id], func(i, j int) bool {
			if adj[id][i].To == adj[id][j].To {
				return adj[id][i].ID < adj[id][j].ID
			}
			return adj[id][i].To < adj[id][j].To
		})
	}
	result := ImpactResult{StartID: start, MaxDepth: maxDepth, CallDepth: callDepth}
	visited := map[string]bool{start: true}
	acs := map[string]bool{}
	type frame struct {
		node  string
		nodes []string
		edges []domain.Edge
		calls int
	}
	queue := []frame{{node: start, nodes: []string{start}}}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if len(current.edges) >= maxDepth {
			if len(adj[current.node]) > 0 {
				result.Truncated = true
			}
			continue
		}
		for _, edge := range adj[current.node] {
			nextCalls := current.calls
			if edge.Kind == "calls" || edge.Kind == "imports" {
				nextCalls++
			}
			if nextCalls > callDepth {
				result.Truncated = true
				continue
			}
			if contains(current.nodes, edge.To) {
				continue
			}
			nodes := appendCopy(current.nodes, edge.To)
			edges := appendEdge(current.edges, edge)
			result.Paths = append(result.Paths, GraphPath{NodeIDs: nodes, Edges: edges})
			visited[edge.To] = true
			if node, ok := g.Nodes[edge.To]; ok && node.Kind == "AcceptanceCriterion" {
				acs[edge.To] = true
			}
			queue = append(queue, frame{node: edge.To, nodes: nodes, edges: edges, calls: nextCalls})
		}
	}
	if g.Nodes[start].Kind == "AcceptanceCriterion" {
		acs[start] = true
	}
	result.ImpactedACIDs = sortedKeys(acs)
	result.VisitedIDs = sortedKeys(visited)
	sort.Slice(result.Paths, func(i, j int) bool { return fmt.Sprint(result.Paths[i].NodeIDs) < fmt.Sprint(result.Paths[j].NodeIDs) })
	return result, nil
}

func ValidateGraph(p *domain.Project) []domain.Finding {
	var findings []domain.Finding
	for _, edge := range sortedEdges(p.Graph.Edges) {
		missing := []string{}
		if _, ok := p.Graph.Nodes[edge.From]; !ok {
			missing = append(missing, edge.From)
		}
		if _, ok := p.Graph.Nodes[edge.To]; !ok {
			missing = append(missing, edge.To)
		}
		if len(missing) > 0 {
			findings = append(findings, domain.Finding{Code: "GRAPH_EDGE_DANGLING", Severity: "blocker", SubjectRefs: append([]string{edge.ID}, missing...), Message: "graph edge references a missing node", Remediation: "Rebuild the graph or repair the source traceability reference."})
		}
	}
	for id, in := range p.Intents {
		if in.Status != "confirmed" {
			continue
		}
		hasAC := false
		for _, ac := range p.Criteria {
			if ac.IntentID == id {
				hasAC = true
			}
		}
		if !hasAC {
			findings = append(findings, domain.Finding{Code: "GRAPH_INTENT_ORPHAN", Severity: "blocker", SubjectRefs: []string{id}, Message: "confirmed intent has no acceptance criterion", Remediation: "Add an acceptance criterion."})
		}
	}
	for id, task := range p.Tasks {
		if len(task.CriterionIDs) == 0 && len(task.IntentIDs) == 0 {
			findings = append(findings, domain.Finding{Code: "GRAPH_TASK_ORPHAN", Severity: "warning", SubjectRefs: []string{id}, Message: "task has no intent or acceptance criterion", Remediation: "Link the task to its requirement or document why it is infrastructure-only."})
		}
	}
	for id, ev := range p.Evidence {
		if len(ev.CriterionIDs) == 0 {
			findings = append(findings, domain.Finding{Code: "GRAPH_EVIDENCE_ORPHAN", Severity: "warning", SubjectRefs: []string{id}, Message: "evidence verifies no acceptance criterion", Remediation: "Link evidence to an AC or keep it outside the blocking evidence set."})
		}
	}
	for id, ac := range p.Criteria {
		if ac.Priority != "blocking" {
			continue
		}
		verified := false
		for _, ev := range p.Evidence {
			if contains(ev.CriterionIDs, id) {
				verified = true
				break
			}
		}
		if !verified {
			findings = append(findings, domain.Finding{Code: "GRAPH_AC_UNVERIFIED", Severity: "warning", SubjectRefs: []string{id}, Message: "blocking criterion has no evidence edge", Remediation: "Register QA evidence."})
		}
	}
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Code == findings[j].Code {
			return fmt.Sprint(findings[i].SubjectRefs) < fmt.Sprint(findings[j].SubjectRefs)
		}
		return findings[i].Code < findings[j].Code
	})
	return findings
}

func sortedEdges(in map[string]domain.Edge) []domain.Edge {
	out := make([]domain.Edge, 0, len(in))
	for _, e := range in {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
func sortedKeys[T any](in map[string]T) []string {
	out := make([]string, 0, len(in))
	for k := range in {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
func contains(in []string, value string) bool {
	for _, x := range in {
		if x == value {
			return true
		}
	}
	return false
}
func appendCopy(in []string, value string) []string {
	out := append([]string(nil), in...)
	return append(out, value)
}
func appendEdge(in []domain.Edge, value domain.Edge) []domain.Edge {
	out := append([]domain.Edge(nil), in...)
	return append(out, value)
}
