// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package domain

import "time"

const SchemaVersion = "1.0.0"

type Phase string

const (
	PhaseDraft     Phase = "draft"
	PhasePRD       Phase = "prd"
	PhasePlan      Phase = "plan"
	PhaseDesign    Phase = "design"
	PhaseDo        Phase = "do"
	PhaseLoopCheck Phase = "loop-check"
	PhaseQA        Phase = "qa"
	PhaseReport    Phase = "report"
	PhaseArchived  Phase = "archived"
)

var PhaseOrder = []Phase{PhaseDraft, PhasePRD, PhasePlan, PhaseDesign, PhaseDo, PhaseLoopCheck, PhaseQA, PhaseReport, PhaseArchived}

type Actor struct {
	Kind      string `json:"kind"`
	ID        string `json:"id,omitempty"`
	SessionID string `json:"session_id,omitempty"`
}

type Project struct {
	SchemaVersion  string                   `json:"schema_version"`
	ProjectID      string                   `json:"project_id"`
	Name           string                   `json:"name"`
	Profile        string                   `json:"profile"`
	MasterPlan     MasterPlan               `json:"master_plan"`
	ActiveSprintID string                   `json:"active_sprint_id,omitempty"`
	Revision       uint64                   `json:"revision"`
	UpdatedAt      time.Time                `json:"updated_at"`
	Sprints        map[string]*Sprint       `json:"sprints"`
	Tasks          map[string]*Task         `json:"tasks"`
	Intents        map[string]*Intent       `json:"intents"`
	Criteria       map[string]*Criterion    `json:"acceptance_criteria"`
	Gaps           map[string]*Gap          `json:"gaps"`
	Waivers        map[string]*Waiver       `json:"waivers"`
	Approvals      map[string]*Approval     `json:"approvals"`
	Evidence       map[string]*Evidence     `json:"evidence"`
	QARuns         map[string]*QARun        `json:"qa_runs"`
	Graph          Graph                    `json:"graph"`
	RequestResults map[string]RequestResult `json:"request_results,omitempty"`
}

type MasterPlan struct {
	Objective             string   `json:"objective"`
	Milestones            []string `json:"milestones"`
	Releases              []string `json:"releases"`
	CommonRisks           []string `json:"common_risks"`
	CrossSprintInvariants []string `json:"cross_sprint_invariants"`
}

type Sprint struct {
	SprintID          string     `json:"sprint_id"`
	Slug              string     `json:"slug"`
	Title             string     `json:"title"`
	Milestone         string     `json:"milestone,omitempty"`
	Release           string     `json:"release,omitempty"`
	Phase             Phase      `json:"phase"`
	Predecessors      []string   `json:"predecessor_ids,omitempty"`
	IntentIDs         []string   `json:"intent_ids,omitempty"`
	TaskIDs           []string   `json:"task_ids,omitempty"`
	OpenGapIDs        []string   `json:"open_gap_ids,omitempty"`
	DocumentRoot      string     `json:"document_root"`
	StartedAt         *time.Time `json:"started_at,omitempty"`
	ArchivedAt        *time.Time `json:"archived_at,omitempty"`
	ApprovalRefs      []string   `json:"approval_refs,omitempty"`
	LastQAStatus      string     `json:"last_qa_status,omitempty"`
	LastQAID          string     `json:"last_qa_id,omitempty"`
	ReportPath        string     `json:"report_path,omitempty"`
	LoopIteration     int        `json:"loop_iteration"`
	MaxLoopIterations int        `json:"max_loop_iterations"`
	LastLoopOutcome   string     `json:"last_loop_outcome,omitempty"`
	LastLoopSummary   string     `json:"last_loop_summary,omitempty"`
	LastLoopAt        *time.Time `json:"last_loop_at,omitempty"`
}

type Task struct {
	TaskID       string   `json:"task_id"`
	SprintID     string   `json:"sprint_id"`
	Title        string   `json:"title"`
	Status       string   `json:"status"`
	IntentIDs    []string `json:"intent_ids,omitempty"`
	CriterionIDs []string `json:"ac_ids,omitempty"`
	DependsOn    []string `json:"depends_on,omitempty"`
	Layer        string   `json:"layer"`
	Artifacts    []string `json:"artifacts,omitempty"`
}

type Intent struct {
	IntentID        string     `json:"intent_id"`
	SprintID        string     `json:"sprint_id"`
	Revision        uint64     `json:"revision"`
	Status          string     `json:"status"`
	Statement       string     `json:"statement"`
	Rationale       string     `json:"rationale,omitempty"`
	Actors          []string   `json:"actors,omitempty"`
	DesiredOutcomes []string   `json:"desired_outcomes,omitempty"`
	NonGoals        []string   `json:"non_goals,omitempty"`
	Policies        []string   `json:"policies,omitempty"`
	BusinessRules   []string   `json:"business_rules,omitempty"`
	UXStates        []string   `json:"ux_states,omitempty"`
	DataInvariants  []string   `json:"data_invariants,omitempty"`
	Constraints     []string   `json:"constraints,omitempty"`
	Assumptions     []string   `json:"assumptions,omitempty"`
	OpenQuestions   []string   `json:"open_questions,omitempty"`
	Source          string     `json:"source"`
	SourceLocator   string     `json:"source_locator,omitempty"`
	ConfirmedBy     string     `json:"confirmed_by,omitempty"`
	ConfirmedAt     *time.Time `json:"confirmed_at,omitempty"`
	Supersedes      string     `json:"supersedes,omitempty"`
}

type Criterion struct {
	CriterionID   string   `json:"ac_id"`
	IntentID      string   `json:"intent_id"`
	Statement     string   `json:"statement"`
	Priority      string   `json:"priority"`
	Observable    string   `json:"observable"`
	Preconditions []string `json:"preconditions,omitempty"`
	Expected      []string `json:"expected,omitempty"`
	Forbidden     []string `json:"forbidden,omitempty"`
}

type Gap struct {
	GapID                 string     `json:"gap_id"`
	SprintID              string     `json:"sprint_id"`
	Category              string     `json:"category"`
	Severity              string     `json:"severity"`
	Status                string     `json:"status"`
	Description           string     `json:"description"`
	SubjectRefs           []string   `json:"subject_refs,omitempty"`
	EvidenceRefs          []string   `json:"evidence_refs,omitempty"`
	Resolution            string     `json:"resolution,omitempty"`
	ResolutionEvidenceIDs []string   `json:"resolution_evidence_ids,omitempty"`
	DeferredReason        string     `json:"deferred_reason,omitempty"`
	DeferredOwner         string     `json:"deferred_owner,omitempty"`
	DeferredTargetSprint  string     `json:"deferred_target_sprint,omitempty"`
	DeferredAt            *time.Time `json:"deferred_at,omitempty"`
	Fingerprint           string     `json:"fingerprint,omitempty"`
	DetectedBy            string     `json:"detected_by,omitempty"`
	DetectedAtRevision    uint64     `json:"detected_at_revision,omitempty"`
}

type Approval struct {
	ApprovalID  string     `json:"approval_id"`
	SprintID    string     `json:"sprint_id"`
	From        Phase      `json:"from"`
	To          Phase      `json:"to"`
	Reason      string     `json:"reason"`
	Requester   string     `json:"requester"`
	Approver    string     `json:"approver,omitempty"`
	Status      string     `json:"status"`
	RequestedAt time.Time  `json:"requested_at"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
	ExpiresAt   time.Time  `json:"expires_at"`
	ConsumedAt  *time.Time `json:"consumed_at,omitempty"`
}

type Waiver struct {
	WaiverID    string     `json:"waiver_id"`
	SprintID    string     `json:"sprint_id"`
	GapID       string     `json:"gap_id"`
	Reason      string     `json:"reason"`
	Scope       string     `json:"scope"`
	Approver    string     `json:"approver"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	RevokedAt   *time.Time `json:"revoked_at,omitempty"`
	Requester   string     `json:"requester,omitempty"`
	RequestedAt *time.Time `json:"requested_at,omitempty"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
	ExpiredAt   *time.Time `json:"expired_at,omitempty"`
}

type Evidence struct {
	EvidenceID      string              `json:"evidence_id"`
	SprintID        string              `json:"sprint_id"`
	RunID           string              `json:"run_id,omitempty"`
	Kind            string              `json:"kind"`
	URI             string              `json:"uri"`
	SHA256          string              `json:"sha256"`
	Size            int64               `json:"size"`
	CriterionIDs    []string            `json:"ac_ids,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
	RedactionStatus string              `json:"redaction_status"`
	CaseID          string              `json:"case_id,omitempty"`
	SpecHash        string              `json:"spec_hash,omitempty"`
	StateRevision   uint64              `json:"state_revision,omitempty"`
	Layers          []string            `json:"layers,omitempty"`
	Assertions      []EvidenceAssertion `json:"assertions,omitempty"`
	Tool            string              `json:"tool,omitempty"`
	ToolVersion     string              `json:"tool_version,omitempty"`
	Environment     string              `json:"environment,omitempty"`
	StartedAt       *time.Time          `json:"started_at,omitempty"`
	FinishedAt      *time.Time          `json:"finished_at,omitempty"`
	SupersededBy    string              `json:"superseded_by,omitempty"`
}

type EvidenceAssertion struct {
	Statement       string   `json:"statement"`
	Passed          bool     `json:"passed"`
	Layer           string   `json:"layer"`
	RequirementRefs []string `json:"requirement_refs"`
	Actual          string   `json:"actual,omitempty"`
	Expected        string   `json:"expected,omitempty"`
}

type QACase struct {
	CaseID            string            `json:"case_id"`
	CriterionIDs      []string          `json:"ac_ids"`
	Title             string            `json:"title"`
	Variant           string            `json:"variant"`
	Layers            []string          `json:"layers"`
	Status            string            `json:"status"`
	EvidenceIDs       []string          `json:"evidence_ids,omitempty"`
	Actor             string            `json:"actor,omitempty"`
	Preconditions     []string          `json:"preconditions,omitempty"`
	Steps             []QAStep          `json:"steps,omitempty"`
	ForbiddenOutcomes []string          `json:"forbidden_outcomes,omitempty"`
	RequiredLayers    map[string]string `json:"required_layers,omitempty"`
	Risk              string            `json:"risk,omitempty"`
}

type QAStep struct {
	Action       string   `json:"action"`
	ExpectedUI   string   `json:"expected_ui,omitempty"`
	ExpectedData string   `json:"expected_data,omitempty"`
	ObserverIDs  []string `json:"observer_ids,omitempty"`
}

type QARun struct {
	RunID         string     `json:"run_id"`
	SprintID      string     `json:"sprint_id"`
	Status        string     `json:"status"`
	Environment   string     `json:"environment"`
	Cases         []QACase   `json:"cases"`
	StartedAt     time.Time  `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at,omitempty"`
	StateRevision uint64     `json:"state_revision"`
	SpecHash      string     `json:"spec_hash"`
}

type Graph struct {
	Nodes map[string]Node `json:"nodes"`
	Edges map[string]Edge `json:"edges"`
}

type Node struct {
	ID         string         `json:"id"`
	Kind       string         `json:"kind"`
	Label      string         `json:"label"`
	Locator    string         `json:"locator,omitempty"`
	Source     string         `json:"source"`
	Confidence float64        `json:"confidence"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

type Edge struct {
	ID            string  `json:"id"`
	From          string  `json:"from"`
	To            string  `json:"to"`
	Kind          string  `json:"kind"`
	SourceLocator string  `json:"source_locator"`
	Provider      string  `json:"provider"`
	Confidence    float64 `json:"confidence"`
}

type Event struct {
	Sequence         uint64    `json:"sequence"`
	EventID          string    `json:"event_id"`
	EventType        string    `json:"event_type"`
	AggregateID      string    `json:"aggregate_id"`
	OccurredAt       time.Time `json:"occurred_at"`
	Actor            Actor     `json:"actor"`
	ExpectedRevision uint64    `json:"expected_revision"`
	Payload          any       `json:"payload"`
	PreviousHash     string    `json:"previous_hash"`
	Hash             string    `json:"hash"`
}

type RequestResult struct {
	Revision    uint64 `json:"revision"`
	CommandHash string `json:"command_hash"`
	Result      any    `json:"result"`
	Completed   bool   `json:"completed"`
}

type Finding struct {
	Code         string   `json:"code"`
	Severity     string   `json:"severity"`
	SubjectRefs  []string `json:"subject_refs,omitempty"`
	Message      string   `json:"message"`
	EvidenceRefs []string `json:"evidence_refs,omitempty"`
	Remediation  string   `json:"remediation"`
	Waivable     bool     `json:"waivable"`
}

func NewProject(id, name, profile string, now time.Time) *Project {
	return &Project{
		SchemaVersion: SchemaVersion, ProjectID: id, Name: name, Profile: profile,
		UpdatedAt: now.UTC(), MasterPlan: MasterPlan{Objective: name, Milestones: []string{}, Releases: []string{}, CommonRisks: []string{}, CrossSprintInvariants: []string{"blocking acceptance criteria require valid evidence"}}, Sprints: map[string]*Sprint{}, Tasks: map[string]*Task{},
		Intents: map[string]*Intent{}, Criteria: map[string]*Criterion{}, Gaps: map[string]*Gap{},
		Waivers: map[string]*Waiver{}, Approvals: map[string]*Approval{},
		Evidence: map[string]*Evidence{}, QARuns: map[string]*QARun{},
		Graph:          Graph{Nodes: map[string]Node{}, Edges: map[string]Edge{}},
		RequestResults: map[string]RequestResult{},
	}
}
