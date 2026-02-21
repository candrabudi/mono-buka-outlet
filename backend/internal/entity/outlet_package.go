package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type OutletPackage struct {
	ID           uuid.UUID      `json:"id"`
	OutletID     uuid.UUID      `json:"outlet_id"`
	Name         string         `json:"name"`
	Slug         string         `json:"slug"`
	Price        int64          `json:"price"`
	MinimumDP    int64          `json:"minimum_dp"`
	Duration     string         `json:"duration"`
	Image        string         `json:"image"`
	EstimatedBEP string         `json:"estimated_bep"`
	NetProfit    string         `json:"net_profit"`
	Description  string         `json:"description"`
	Benefits     pq.StringArray `json:"benefits"`
	SortOrder    int            `json:"sort_order"`
	IsActive     bool           `json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
