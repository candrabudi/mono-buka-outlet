package entity

import (
	"time"

	"github.com/google/uuid"
)

type Partnership struct {
	ID                 uuid.UUID  `json:"id"`
	BrandID            *uuid.UUID `json:"brand_id,omitempty"`
	MitraID            uuid.UUID  `json:"mitra_id"`
	AffiliatorID       *uuid.UUID `json:"affiliator_id,omitempty"`
	OutletID           *uuid.UUID `json:"outlet_id,omitempty"`
	PackageID          *uuid.UUID `json:"package_id,omitempty"`
	ProgressPercentage int        `json:"progress_percentage"`
	Status             string     `json:"status"`
	StartDate          *time.Time `json:"start_date,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Brand   *Brand         `json:"brand,omitempty"`
	Mitra   *User          `json:"mitra,omitempty"`
	Affiliator *User          `json:"affiliator,omitempty"`
	Outlet  *Outlet        `json:"outlet,omitempty"`
	Package *OutletPackage `json:"package,omitempty"`
}

const (
	PartnershipStatusPending         = "PENDING"
	PartnershipStatusDPVerified      = "DP_VERIFIED"
	PartnershipStatusAgreementSigned = "AGREEMENT_SIGNED"
	PartnershipStatusDevelopment     = "DEVELOPMENT"
	PartnershipStatusRunning         = "RUNNING"
	PartnershipStatusCompleted       = "COMPLETED"
)
