package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Password     string     `json:"-"`
	Phone        string     `json:"phone"`
	Role         string     `json:"role"`
	ReferralCode string     `json:"referral_code,omitempty"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

const (
	RoleMaster     = "master"     // Full access to all data
	RoleAdmin      = "admin"      // Kemitraan, outlet, management mitra, invoice
	RoleFinance    = "finance"    // Keuangan
	RoleAffiliator = "affiliator" // Affiliator yang mengelola referral
	RoleMitra      = "mitra"      // Mitra portal (separate frontend)
)

func ValidRoles() []string {
	return []string{
		RoleMaster,
		RoleAdmin,
		RoleFinance,
		RoleAffiliator,
		RoleMitra,
	}
}

func IsValidRole(role string) bool {
	for _, r := range ValidRoles() {
		if r == role {
			return true
		}
	}
	return false
}
