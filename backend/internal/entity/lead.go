package entity

import (
	"time"

	"github.com/google/uuid"
)

type Lead struct {
	ID                 uuid.UUID  `json:"id"`
	BrandID            uuid.UUID  `json:"brand_id"`
	SalesID            *uuid.UUID `json:"sales_id,omitempty"`
	FullName           string     `json:"full_name"`
	Email              string     `json:"email"`
	Phone              string     `json:"phone"`
	Status             string     `json:"status"`
	ProgressPercentage int        `json:"progress_percentage"`
	Notes              string     `json:"notes"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Brand *Brand `json:"brand,omitempty"`
	Sales *User  `json:"sales,omitempty"`
}

const (
	LeadStatusNew               = "NEW"
	LeadStatusConsultation      = "CONSULTATION"
	LeadStatusLocationSubmitted = "LOCATION_SUBMITTED"
	LeadStatusSurveyApproved    = "SURVEY_APPROVED"
	LeadStatusMeetingDone       = "MEETING_DONE"
	LeadStatusReadyForDP        = "READY_FOR_DP"
	LeadStatusDPPaid            = "DP_PAID"
	LeadStatusAgreementReview   = "AGREEMENT_REVIEW"
	LeadStatusFullyPaid         = "FULLY_PAID"
	LeadStatusActivePartnership = "ACTIVE_PARTNERSHIP"
	LeadStatusRunning           = "RUNNING"
	LeadStatusCompleted         = "COMPLETED"
)

func LeadStatuses() []string {
	return []string{
		LeadStatusNew,
		LeadStatusConsultation,
		LeadStatusLocationSubmitted,
		LeadStatusSurveyApproved,
		LeadStatusMeetingDone,
		LeadStatusReadyForDP,
		LeadStatusDPPaid,
		LeadStatusAgreementReview,
		LeadStatusFullyPaid,
		LeadStatusActivePartnership,
		LeadStatusRunning,
		LeadStatusCompleted,
	}
}

func IsValidLeadStatus(status string) bool {
	for _, s := range LeadStatuses() {
		if s == status {
			return true
		}
	}
	return false
}
