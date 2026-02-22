package entity

import (
	"time"

	"github.com/google/uuid"
)

type PartnershipApplication struct {
	ID               uuid.UUID  `json:"id"`
	MitraID          uuid.UUID  `json:"mitra_id"`
	OutletID         uuid.UUID  `json:"outlet_id"`
	PackageID        uuid.UUID  `json:"package_id"`
	Motivation       string     `json:"motivation"`
	Experience       string     `json:"experience"`
	ProposedLocation string     `json:"proposed_location"`
	InvestmentBudget float64    `json:"investment_budget"`
	Status           string     `json:"status"`
	AdminNotes       string     `json:"admin_notes,omitempty"`
	ReviewedBy       *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Relations
	Mitra   *User          `json:"mitra,omitempty"`
	Outlet  *Outlet        `json:"outlet,omitempty"`
	Package *OutletPackage `json:"package,omitempty"`
}

const (
	ApplicationStatusPending  = "PENDING"
	ApplicationStatusReviewed = "REVIEWED"
	ApplicationStatusApproved = "APPROVED"
	ApplicationStatusRejected = "REJECTED"
)
