package entity

import (
	"time"

	"github.com/google/uuid"
)

type Location struct {
	ID             uuid.UUID  `json:"id"`
	LeadID         uuid.UUID  `json:"lead_id"`
	BrandID        uuid.UUID  `json:"brand_id"`
	Latitude       float64    `json:"lat"`
	Longitude      float64    `json:"lng"`
	Address        string     `json:"address"`
	Photo          string     `json:"photo"`
	ApprovalStatus string     `json:"approval_status"`
	SurveyNotes    string     `json:"survey_notes"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

const (
	LocationStatusPending  = "PENDING"
	LocationStatusApproved = "APPROVED"
	LocationStatusRejected = "REJECTED"
)
