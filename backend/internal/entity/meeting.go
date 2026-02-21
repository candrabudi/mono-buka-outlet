package entity

import (
	"time"

	"github.com/google/uuid"
)

// ─── Meeting Types ───
const (
	MeetingTypeEdukasi      = "edukasi"
	MeetingTypeClosing      = "closing"
	MeetingTypeReviewLokasi = "review_lokasi"
	MeetingTypeOperasional  = "operasional"
)

// ─── Meeting Statuses ───
const (
	MeetingStatusScheduled       = "scheduled"
	MeetingStatusOngoing         = "ongoing"
	MeetingStatusCompleted       = "completed"
	MeetingStatusWaitingDecision = "waiting_decision"
	MeetingStatusApproved        = "approved"
	MeetingStatusCancelled       = "cancelled"
)

// ─── Decision Types ───
const (
	DecisionLanjutDP       = "lanjut_dp"
	DecisionRevisiProposal = "revisi_proposal"
	DecisionSurveiUlang    = "survei_ulang"
	DecisionTidakLanjut    = "tidak_lanjut"
)

type Meeting struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	MeetingType string     `json:"meeting_type"`
	MeetingDate time.Time  `json:"meeting_date"`
	Duration    int        `json:"duration"`
	MeetingLink string     `json:"meeting_link"`
	Location    string     `json:"location"`
	Status      string     `json:"status"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Participants []*MeetingParticipant `json:"participants,omitempty"`
	Notes        *MeetingNotes         `json:"notes,omitempty"`
	ActionPlans  []*MeetingActionPlan  `json:"action_plans,omitempty"`
	Files        []*MeetingFile        `json:"files,omitempty"`
	Creator      *User                 `json:"creator,omitempty"`
}

type MeetingParticipant struct {
	ID            uuid.UUID  `json:"id"`
	MeetingID     uuid.UUID  `json:"meeting_id"`
	UserID        *uuid.UUID `json:"user_id,omitempty"`
	ExternalName  string     `json:"external_name"`
	ExternalEmail string     `json:"external_email"`
	ExternalPhone string     `json:"external_phone"`
	Role          string     `json:"role"`
	CreatedAt     time.Time  `json:"created_at"`

	// Relation
	User *User `json:"user,omitempty"`
}

type MeetingNotes struct {
	ID               uuid.UUID `json:"id"`
	MeetingID        uuid.UUID `json:"meeting_id"`
	Purpose          string    `json:"purpose"`
	Summary          string    `json:"summary"`
	DiscussionPoints string    `json:"discussion_points"`
	Decision         string    `json:"decision"`
	NextStep         string    `json:"next_step"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type MeetingActionPlan struct {
	ID        uuid.UUID  `json:"id"`
	MeetingID uuid.UUID  `json:"meeting_id"`
	TaskName  string     `json:"task_name"`
	PIC       string     `json:"pic"`
	Deadline  *time.Time `json:"deadline,omitempty"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type MeetingFile struct {
	ID         uuid.UUID  `json:"id"`
	MeetingID  uuid.UUID  `json:"meeting_id"`
	FileName   string     `json:"file_name"`
	FilePath   string     `json:"file_path"`
	FileType   string     `json:"file_type"`
	UploadedBy *uuid.UUID `json:"uploaded_by,omitempty"`
	UploadedAt time.Time  `json:"uploaded_at"`
}
