package entity

import (
	"time"

	"github.com/google/uuid"
)

type AffiliatorCommission struct {
	ID            uuid.UUID  `json:"id"`
	AffiliatorID  uuid.UUID  `json:"affiliator_id"`
	PartnershipID *uuid.UUID `json:"partnership_id,omitempty"`
	Amount        float64    `json:"amount"`
	Type          string     `json:"type"`
	Description   string     `json:"description"`
	GivenBy       uuid.UUID  `json:"given_by"`
	CreatedAt     time.Time  `json:"created_at"`

	// Relations
	Affiliator  *User        `json:"affiliator,omitempty"`
	Partnership *Partnership `json:"partnership,omitempty"`
	GivenByUser *User        `json:"given_by_user,omitempty"`
}

const (
	CommissionTypeCommission  = "COMMISSION"
	CommissionTypeBonus       = "BONUS"
	CommissionTypeAdjustment  = "ADJUSTMENT"
)

func ValidCommissionTypes() []string {
	return []string{
		CommissionTypeCommission,
		CommissionTypeBonus,
		CommissionTypeAdjustment,
	}
}

func IsValidCommissionType(t string) bool {
	for _, v := range ValidCommissionTypes() {
		if v == t {
			return true
		}
	}
	return false
}
