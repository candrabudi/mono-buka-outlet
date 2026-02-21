package entity

import (
	"time"

	"github.com/google/uuid"
)

type Partnership struct {
	ID                 uuid.UUID  `json:"id"`
	LeadID             *uuid.UUID `json:"lead_id,omitempty"`
	BrandID            *uuid.UUID `json:"brand_id,omitempty"`
	MitraID            uuid.UUID  `json:"mitra_id"`
	LeaderID           *uuid.UUID `json:"leader_id,omitempty"`
	OutletID           *uuid.UUID `json:"outlet_id,omitempty"`
	PackageID          *uuid.UUID `json:"package_id,omitempty"`
	ProgressPercentage int        `json:"progress_percentage"`
	Status             string     `json:"status"`
	StartDate          *time.Time `json:"start_date,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Lead    *Lead          `json:"lead,omitempty"`
	Brand   *Brand         `json:"brand,omitempty"`
	Mitra   *User          `json:"mitra,omitempty"`
	Leader  *User          `json:"leader,omitempty"`
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
