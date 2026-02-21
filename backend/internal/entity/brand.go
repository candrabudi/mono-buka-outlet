package entity

import (
	"time"

	"github.com/google/uuid"
)

// Brand is kept as a minimal struct for backward compatibility with existing
// tables (leads, partnerships, payments, revenues) that still reference brand_id.
// The Brand CRUD feature has been removed from the admin panel.
type Brand struct {
	ID                      uuid.UUID  `json:"id"`
	Name                    string     `json:"name"`
	Logo                    string     `json:"logo"`
	Description             string     `json:"description"`
	MinimumInvestment       float64    `json:"minimum_investment"`
	ProfitSharingPercentage float64    `json:"profit_sharing_percentage"`
	EstimatedROI            string     `json:"estimated_roi"`
	LocationRequirement     string     `json:"location_requirement"`
	IsActive                bool       `json:"is_active"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
	DeletedAt               *time.Time `json:"deleted_at,omitempty"`
}
