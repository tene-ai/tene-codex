package router

import (
	"github.com/tene-ai/tene-codex/internal/domain"
	"testing"
)

func TestExplicitAlwaysWins(t *testing.T) {
	d := Route("please use $tene:qa and explain planning", true, domain.PhasePlan)
	if d.SelectedSkill != "qa" || !d.Explicit || d.Mode != "selected" {
		t.Fatalf("%+v", d)
	}
}

func TestObsoleteExplicitNamesAreNotSelected(t *testing.T) {
	for _, text := range []string{"$tene-codex:tene-qa", "$tene:tene-qa", "$tene-qa"} {
		d := Route(text, true, domain.PhaseQA)
		if d.Explicit {
			t.Fatalf("obsolete invocation selected: %q => %+v", text, d)
		}
	}
}
func TestHighConfidencePhaseRoute(t *testing.T) {
	d := Route("Playwright로 UX 흐름 종합 테스트하고 evidence 검증해줘", true, domain.PhaseQA)
	if d.SelectedSkill != "qa" || d.Mode != "selected" {
		t.Fatalf("%+v", d)
	}
}
func TestMultiIntentOrchestrates(t *testing.T) {
	d := Route("요구사항 정리하고 구현 계획과 아키텍처 설계해줘", true, domain.PhasePRD)
	if d.SelectedSkill != "sprint" || d.Mode != "proposed" || d.MutationAllowed {
		t.Fatalf("%+v", d)
	}
}
func TestHardNegative(t *testing.T) {
	d := Route("코드 변경 없이 일반적인 계획이 뭐야?", false, "")
	if d.Mode != "none" {
		t.Fatalf("%+v", d)
	}
}
func TestWrongPhaseOnlyProposes(t *testing.T) {
	d := Route("Playwright QA 종합 테스트해줘", true, domain.PhaseDesign)
	if d.Mode == "selected" || d.MutationAllowed {
		t.Fatalf("%+v", d)
	}
}
