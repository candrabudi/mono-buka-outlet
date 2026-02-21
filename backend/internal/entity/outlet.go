package entity

import (
	"time"

	"github.com/google/uuid"
)

type Outlet struct {
	ID                      uuid.UUID  `json:"id"`
	Name                    string     `json:"name"`
	Slug                    string     `json:"slug"`
	Logo                    string     `json:"logo"`
	Banner                  string     `json:"banner"`
	CategoryID              *uuid.UUID `json:"category_id"`
	Category                string     `json:"category"`      // kept for backward compat (legacy string col)
	CategoryName            string     `json:"category_name"` // populated via JOIN
	CategoryIcon            string     `json:"category_icon"` // populated via JOIN
	Description             string     `json:"description"`
	ShortDescription        string     `json:"short_description"`
	MinimumInvestment       float64    `json:"minimum_investment"`
	MaximumInvestment       *float64   `json:"maximum_investment,omitempty"`
	ProfitSharingPercentage float64    `json:"profit_sharing_percentage"`
	EstimatedROI            string     `json:"estimated_roi"`
	LocationRequirement     string     `json:"location_requirement"`
	Address                 string     `json:"address"`
	City                    string     `json:"city"`
	Province                string     `json:"province"`
	Latitude                *float64   `json:"latitude,omitempty"`
	Longitude               *float64   `json:"longitude,omitempty"`
	ContactPhone            string     `json:"contact_phone"`
	ContactEmail            string     `json:"contact_email"`
	ContactWhatsapp         string     `json:"contact_whatsapp"`
	Website                 string     `json:"website"`
	IsActive                bool       `json:"is_active"`
	IsFeatured              bool       `json:"is_featured"`
	TotalOutlets            int        `json:"total_outlets"`
	YearEstablished         *int       `json:"year_established,omitempty"`
	CreatedBy               *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
	DeletedAt               *time.Time `json:"deleted_at,omitempty"`
}
