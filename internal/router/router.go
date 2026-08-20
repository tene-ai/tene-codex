// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

// Package router provides a deterministic companion to Codex implicit skill
// discovery. It never mutates workflow state; it selects, proposes, or declines.
package router

import (
	"regexp"
	"sort"
	"strings"

	"github.com/tene-ai/tene-codex/internal/domain"
)

type Candidate struct {
	Skill   string   `json:"skill"`
	Score   float64  `json:"score"`
	Reasons []string `json:"reasons"`
}
type Decision struct {
	Mode              string      `json:"mode"`
	SelectedSkill     string      `json:"selected_skill,omitempty"`
	Candidates        []Candidate `json:"candidates"`
	Explicit          bool        `json:"explicit"`
	NeedsConfirmation bool        `json:"needs_confirmation"`
	MutationAllowed   bool        `json:"mutation_allowed"`
	RequiredActions   []string    `json:"required_actions"`
	ForbiddenActions  []string    `json:"forbidden_actions"`
	Reason            string      `json:"reason"`
}

type rule struct {
	skill                    string
	phase                    domain.Phase
	intent, artifact, action []string
}

var explicitRE = regexp.MustCompile(`(?i)\$?(tene-(?:sprint|prd|plan|design|loop-check|qa|report|status|secrets))\b`)
var rules = []rule{
	{"tene-sprint", "", []string{"sprint 시작", "스프린트 시작", "workflow master plan", "워크플로 마스터 플랜", "master plan", "마스터 플랜", "여러 스프린트", "sprint archive", "스프린트 아카이브", "resume sprint workflow"}, []string{"master plan"}, []string{"시작", "생성", "나눠", "archive", "resume"}},
	{"tene-prd", domain.PhasePRD, []string{"요구사항", "기획 의도", "기획", "사용자 문제", "아이디어", "acceptance criteria", "requirement", "product intent", "policy", "정책", "예외 케이스"}, []string{"prd", "acceptance criteria", "ac"}, []string{"정리", "구체화", "discover", "define", "질문"}},
	{"tene-plan", domain.PhasePlan, []string{"구현 계획", "작업 계획", "계획", "plan", "task 분해", "dependency", "의존성", "implementation plan", "work breakdown", "병렬"}, []string{"plan", "task", "dependency"}, []string{"세워", "나눠", "순서", "plan", "schedule"}},
	{"tene-design", domain.PhaseDesign, []string{"처리 로직", "아키텍처", "architecture", "design", "data flow", "데이터 흐름", "interface", "인터페이스", "symbol", "api contract", "db schema"}, []string{"design", "adr", "schema", "contract"}, []string{"설계", "design", "정해", "trace"}},
	{"tene-loop-check", domain.PhaseLoopCheck, []string{"문서대로", "설계대로", "빠진 요구", "반복 검증", "gap", "불일치", "100%", "spec code", "regression"}, []string{"prd", "design", "diff", "gap"}, []string{"확인", "비교", "수정", "iterate", "검증"}},
	{"tene-qa", domain.PhaseQA, []string{"qa", "종합 테스트", "ux 흐름", "사용자 journey", "user journey", "playwright", "chrome", "e2e", "데이터 처리 흐름", "data journey"}, []string{"test", "evidence", "qa", "playwright"}, []string{"테스트", "검증", "evaluate", "observe"}},
	{"tene-report", domain.PhaseReport, []string{"작업 결과", "회고", "report", "리포트", "변경 파일", "이월", "retrospective", "구현 파일"}, []string{"report", "evidence", "changed files"}, []string{"작성", "정리", "설명", "generate"}},
	{"tene-status", "", []string{"어디까지", "이어서", "현재 단계", "진행 상황", "next action", "status", "blocker", "다음 할 일", "resume"}, []string{"status", "state", "context"}, []string{"보여", "알려", "계속", "resume"}},
	{"tene-secrets", "", []string{"secret", "비밀정보", "api key", "credential", "token", ".env", "환경 변수", "environment variable", "credentials"}, []string{"tene", ".env", "secret"}, []string{"실행", "run", "inject", "주입", "관리"}},
}

func Route(text string, active bool, phase domain.Phase) Decision {
	lower := strings.ToLower(strings.TrimSpace(text))
	if m := explicitRE.FindStringSubmatch(lower); len(m) > 1 {
		return Decision{Mode: "selected", SelectedSkill: m[1], Explicit: true, MutationAllowed: true, Candidates: []Candidate{{Skill: m[1], Score: 1, Reasons: []string{"explicit skill invocation"}}}, RequiredActions: []string{"status --json", "context build"}, ForbiddenActions: []string{"phase skip"}, Reason: "explicit invocation has priority"}
	}
	if hardNegative(lower) {
		return Decision{Mode: "none", Candidates: []Candidate{}, RequiredActions: []string{}, ForbiddenActions: []string{"workflow mutation"}, Reason: "general explanation or explicit opt-out"}
	}
	var cs []Candidate
	for _, r := range rules {
		i, ir := matches(lower, r.intent)
		a, ar := matches(lower, r.artifact)
		x, xr := matches(lower, r.action)
		if !i && !a {
			continue
		}
		score := 0.0
		reasons := []string{}
		if i {
			score += .40
			reasons = append(reasons, "intent:"+ir)
		}
		if a {
			score += .20
			reasons = append(reasons, "artifact:"+ar)
		}
		if x {
			score += .15
			reasons = append(reasons, "action:"+xr)
		}
		if active && (r.phase == "" || r.phase == phase) {
			score += .25
			reasons = append(reasons, "phase compatible")
		}
		if !active && r.skill == "tene-sprint" {
			score += .25
			reasons = append(reasons, "no active sprint")
		}
		cs = append(cs, Candidate{r.skill, score, reasons})
	}
	sort.SliceStable(cs, func(i, j int) bool { return cs[i].Score > cs[j].Score })
	if len(cs) == 0 {
		return Decision{Mode: "none", Candidates: cs, ForbiddenActions: []string{"workflow mutation"}, Reason: "no skill reached the routing floor"}
	}
	strong := []Candidate{}
	for _, c := range cs {
		if c.Score >= .55 {
			strong = append(strong, c)
		}
	}
	intentful := 0
	for _, c := range cs {
		if c.Score >= .40 {
			intentful++
		}
	}
	if intentful > 1 && multiCue(lower) {
		return Decision{Mode: "proposed", SelectedSkill: "tene-sprint", Candidates: cs, NeedsConfirmation: true, MutationAllowed: false, RequiredActions: []string{"status --json", "coordinate phase skills"}, ForbiddenActions: []string{"phase skip", "implicit approval"}, Reason: "multiple intents require sprint orchestration"}
	}
	if len(strong) > 1 {
		return Decision{Mode: "proposed", SelectedSkill: "tene-sprint", Candidates: cs, NeedsConfirmation: true, MutationAllowed: false, RequiredActions: []string{"status --json", "coordinate phase skills"}, ForbiddenActions: []string{"phase skip", "implicit approval"}, Reason: "multiple intents require sprint orchestration"}
	}
	margin := 1.0
	if len(cs) > 1 {
		margin = cs[0].Score - cs[1].Score
	}
	if cs[0].Score >= .80 && margin >= .10 {
		return Decision{Mode: "selected", SelectedSkill: cs[0].Skill, Candidates: cs, MutationAllowed: true, RequiredActions: []string{"status --json", "context build"}, ForbiddenActions: []string{"phase skip", "implicit approval"}, Reason: "single high-confidence route"}
	}
	if cs[0].Score >= .60 {
		return Decision{Mode: "proposed", SelectedSkill: cs[0].Skill, Candidates: cs, NeedsConfirmation: true, MutationAllowed: false, RequiredActions: []string{"status --json"}, ForbiddenActions: []string{"workflow mutation", "phase skip"}, Reason: "ambiguous or medium-confidence route"}
	}
	return Decision{Mode: "none", Candidates: cs, MutationAllowed: false, ForbiddenActions: []string{"workflow mutation"}, Reason: "score below proposal threshold"}
}

func matches(text string, cues []string) (bool, string) {
	for _, c := range cues {
		if strings.Contains(text, strings.ToLower(c)) {
			return true, c
		}
	}
	return false, ""
}
func hardNegative(t string) bool {
	for _, x := range []string{"tene 없이", "without tene", "일반적인 설명", "general explanation", "개념만 설명", "코드 변경 없이", "what is agile", "애자일이 뭐", "계획이 뭐야"} {
		if strings.Contains(t, x) {
			return true
		}
	}
	return false
}

func multiCue(t string) bool {
	for _, x := range []string{"하고", " 한 후", "부터", "까지", " and ", " then ", ",", "모두", "같이"} {
		if strings.Contains(t, x) {
			return true
		}
	}
	return false
}
