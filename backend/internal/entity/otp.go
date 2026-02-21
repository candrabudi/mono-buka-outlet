package entity

import (
	"time"

	"github.com/google/uuid"
)

type OTPCode struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Email     string     `json:"email"`
	Code      string     `json:"-"`
	Purpose   string     `json:"purpose"` // admin_login, mitra_login, reset_password
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

const (
	OTPPurposeAdminLogin = "admin_login"
	OTPPurposeMitraLogin = "mitra_login"
	OTPPurposeResetPwd   = "reset_password"
	OTPExpiryMinutes     = 5
	OTPLength            = 6
)

func (o *OTPCode) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

func (o *OTPCode) IsUsed() bool {
	return o.UsedAt != nil
}

// AdminRoles returns roles allowed to access admin panel
func AdminRoles() []string {
	return []string{
		RoleMaster,
		RoleAdmin,
		RoleFinance,
	}
}

// IsAdminRole checks if a role is allowed for admin panel
func IsAdminRole(role string) bool {
	for _, r := range AdminRoles() {
		if r == role {
			return true
		}
	}
	return false
}
