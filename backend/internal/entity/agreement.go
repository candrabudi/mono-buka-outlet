package entity

import (
	"time"

	"github.com/google/uuid"
)

type Agreement struct {
	ID            uuid.UUID  `json:"id"`
	PartnershipID uuid.UUID  `json:"partnership_id"`
	BrandID       uuid.UUID  `json:"brand_id"`
	Title         string     `json:"title"`
	Type          string     `json:"type"`
	FileURL       string     `json:"file_url"`
	Version       int        `json:"version"`
	Status        string     `json:"status"`
	SignedAt      *time.Time `json:"signed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

const (
	AgreementStatusDraft   = "DRAFT"
	AgreementStatusReview  = "REVIEW"
	AgreementStatusSigned  = "SIGNED"
	AgreementStatusExpired = "EXPIRED"

	AgreementTypeContract = "CONTRACT"
	AgreementTypeDocument = "DOCUMENT"
)
