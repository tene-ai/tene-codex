// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/tene-ai/tene-codex/internal/domain"
	"github.com/tene-ai/tene-codex/internal/router"
)

type skillSet struct {
	Phase domain.Phase `json:"phase"`
	Stems []string     `json:"stems"`
}
type corpus struct {
	SchemaVersion      string              `json:"schema_version"`
	PositiveSuffixes   []string            `json:"positive_suffixes"`
	Skills             map[string]skillSet `json:"skills"`
	NegativePrompts    []string            `json:"negative_prompts"`
	MultiIntentPrompts []string            `json:"multi_intent_prompts"`
}
type metric struct {
	ExplicitTotal, ExplicitPassed, PositiveTotal, TruePositive, NegativeTotal, FalsePositive, MultiTotal, MultiPassed, ConflictTotal, WrongPhase int
	Precision, Recall, WrongPhaseRate, UnnecessaryTriggerRate, MultiIntentRate                                                                   float64
	Passed                                                                                                                                       bool
}
type output struct {
	SchemaVersion string            `json:"schema_version"`
	Metrics       map[string]metric `json:"metrics"`
	OverallPassed bool              `json:"overall_passed"`
}

func main() {
	path := "evals/routing-corpus.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	b, e := os.ReadFile(path)
	if e != nil {
		fatal(e)
	}
	var c corpus
	if e = json.Unmarshal(b, &c); e != nil {
		fatal(e)
	}
	out := output{c.SchemaVersion, map[string]metric{}, true}
	names := make([]string, 0, len(c.Skills))
	for n := range c.Skills {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, name := range names {
		s := c.Skills[name]
		m := metric{}
		for i := 0; i < 5; i++ {
			m.ExplicitTotal++
			d := router.Route("use $tene:"+name+" now", true, s.Phase)
			if d.SelectedSkill == name && d.Explicit && d.Mode == "selected" {
				m.ExplicitPassed++
			}
		}
		for _, stem := range s.Stems {
			for _, suffix := range c.PositiveSuffixes {
				m.PositiveTotal++
				d := router.Route(stem+" "+suffix, true, s.Phase)
				if d.SelectedSkill == name && (d.Mode == "selected" || d.Mode == "proposed") {
					m.TruePositive++
				}
			}
		}
		for _, p := range c.NegativePrompts {
			m.NegativeTotal++
			d := router.Route(p, false, "")
			if d.SelectedSkill == name {
				m.FalsePositive++
			}
		}
		if s.Phase != "" && name != "sprint" && name != "status" && name != "secrets" {
			for _, stem := range s.Stems {
				for _, wrong := range []domain.Phase{domain.PhasePRD, domain.PhasePlan, domain.PhaseDesign, domain.PhaseQA, domain.PhaseReport} {
					if wrong == s.Phase {
						continue
					}
					m.ConflictTotal++
					d := router.Route(stem+" 해줘", true, wrong)
					if d.SelectedSkill == name && d.Mode == "selected" {
						m.WrongPhase++
					}
				}
			}
		}
		for _, p := range c.MultiIntentPrompts {
			m.MultiTotal++
			d := router.Route(p, true, s.Phase)
			if d.SelectedSkill == "sprint" && d.Mode == "proposed" {
				m.MultiPassed++
			}
		}
		den := m.TruePositive + m.FalsePositive
		if den > 0 {
			m.Precision = float64(m.TruePositive) / float64(den)
		}
		m.Recall = float64(m.TruePositive) / float64(m.PositiveTotal)
		if m.ConflictTotal > 0 {
			m.WrongPhaseRate = float64(m.WrongPhase) / float64(m.ConflictTotal)
		}
		m.UnnecessaryTriggerRate = float64(m.FalsePositive) / float64(m.NegativeTotal)
		m.MultiIntentRate = float64(m.MultiPassed) / float64(m.MultiTotal)
		m.Passed = m.ExplicitPassed == m.ExplicitTotal && m.Precision >= .90 && m.Recall >= .90 && m.MultiIntentRate >= .90 && m.WrongPhaseRate <= .10 && m.UnnecessaryTriggerRate <= .10
		if !m.Passed {
			out.OverallPassed = false
		}
		out.Metrics[name] = m
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
	if !out.OverallPassed {
		os.Exit(1)
	}
}
func fatal(e error) { fmt.Fprintln(os.Stderr, e); os.Exit(2) }
