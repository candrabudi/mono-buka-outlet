package entity

import (
	"time"

	"github.com/google/uuid"
)

type AffiliatorWithdrawal struct {
	ID            uuid.UUID  `json:"id"`
	AffiliatorID  uuid.UUID  `json:"affiliator_id"`
	Amount        float64    `json:"amount"`
	BankName      string     `json:"bank_name"`
	AccountNumber string     `json:"account_number"`
	AccountHolder string     `json:"account_holder"`
	Status        string     `json:"status"`
	AdminNotes    string     `json:"admin_notes,omitempty"`
	ProcessedBy   *uuid.UUID `json:"processed_by,omitempty"`
	ProcessedAt   *time.Time `json:"processed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Relations
	Affiliator    *User `json:"affiliator,omitempty"`
	ProcessedUser *User `json:"processed_user,omitempty"`
}

const (
	WithdrawalStatusPending     = "PENDING"
	WithdrawalStatusApproved    = "APPROVED"
	WithdrawalStatusTransferred = "TRANSFERRED"
	WithdrawalStatusRejected    = "REJECTED"
)
