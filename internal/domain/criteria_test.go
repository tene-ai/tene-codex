// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package domain

import "testing"

func TestConfirmedCriteriaExcludesInactiveAndOrphanedOwners(t *testing.T) {
	project := &Project{
		Intents: map[string]*Intent{
			"confirmed":  {IntentID: "confirmed", Status: "confirmed"},
			"deprecated": {IntentID: "deprecated", Status: "deprecated"},
			"candidate":  {IntentID: "candidate", Status: "candidate"},
		},
		Criteria: map[string]*Criterion{
			"ac-confirmed":  {CriterionID: "ac-confirmed", IntentID: "confirmed"},
			"ac-deprecated": {CriterionID: "ac-deprecated", IntentID: "deprecated"},
			"ac-candidate":  {CriterionID: "ac-candidate", IntentID: "candidate"},
			"ac-orphan":     {CriterionID: "ac-orphan", IntentID: "missing"},
		},
	}
	sprint := &Sprint{IntentIDs: []string{"confirmed", "deprecated", "candidate", "missing"}}
	criteria := ConfirmedCriteria(project, sprint)
	if len(criteria) != 1 || criteria[0].CriterionID != "ac-confirmed" {
		t.Fatalf("unexpected active criteria: %#v", criteria)
	}
}
