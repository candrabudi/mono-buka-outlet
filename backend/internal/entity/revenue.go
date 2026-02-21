package entity

import (
	"time"

	"github.com/google/uuid"
)

type Revenue struct {
	ID            uuid.UUID  `json:"id"`
	PartnershipID uuid.UUID  `json:"partnership_id"`
	BrandID       uuid.UUID  `json:"brand_id"`
	Month         string     `json:"month"`
	Revenue       float64    `json:"revenue"`
	Expense       float64    `json:"expense"`
	Profit        float64    `json:"profit"`
	CompanyShare  float64    `json:"company_share"`
	MitraShare    float64    `json:"mitra_share"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Partnership *Partnership `json:"partnership,omitempty"`
}
