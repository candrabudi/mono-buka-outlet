package entity

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID             uuid.UUID  `json:"id"`
	PartnershipID  uuid.UUID  `json:"partnership_id"`
	BrandID        uuid.UUID  `json:"brand_id"`
	Type           string     `json:"type"`
	Amount         float64    `json:"amount"`
	ProofURL       string     `json:"proof_url"`
	VerifiedStatus string     `json:"verified_status"`
	VerifiedBy     *uuid.UUID `json:"verified_by,omitempty"`
	VerifiedAt     *time.Time `json:"verified_at,omitempty"`
	Notes          string     `json:"notes"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Partnership *Partnership `json:"partnership,omitempty"`
	Verifier    *User        `json:"verifier,omitempty"`
}

const (
	PaymentTypeDP        = "DP"
	PaymentTypePelunasan = "PELUNASAN"

	PaymentStatusPending  = "PENDING"
	PaymentStatusVerified = "VERIFIED"
	PaymentStatusRejected = "REJECTED"
)
