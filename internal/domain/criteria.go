// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package domain

import "sort"

// ConfirmedCriteria returns only criteria owned by confirmed intents in the
// selected sprint. Deprecated, superseded, candidate, and orphaned records
// remain auditable but must not become active context or QA obligations.
func ConfirmedCriteria(project *Project, sprint *Sprint) []*Criterion {
	if project == nil || sprint == nil {
		return nil
	}
	owned := make(map[string]bool, len(sprint.IntentIDs))
	for _, id := range sprint.IntentIDs {
		if intent := project.Intents[id]; intent != nil && intent.Status == "confirmed" {
			owned[id] = true
		}
	}
	criteria := make([]*Criterion, 0)
	for _, criterion := range project.Criteria {
		if criterion != nil && owned[criterion.IntentID] {
			criteria = append(criteria, criterion)
		}
	}
	sort.Slice(criteria, func(i, j int) bool {
		return criteria[i].CriterionID < criteria[j].CriterionID
	})
	return criteria
}
