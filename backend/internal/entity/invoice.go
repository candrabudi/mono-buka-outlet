package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	InvoiceStatusPending  = "PENDING"
	InvoiceStatusPaid     = "PAID"
	InvoiceStatusExpired  = "EXPIRED"
	InvoiceStatusFailed   = "FAILED"
	InvoiceStatusCanceled = "CANCELED"

	InvoiceTypeDP        = "DP"
	InvoiceTypeCicilan   = "CICILAN"
	InvoiceTypePelunasan = "PELUNASAN"
	InvoiceTypeInvoice   = "INVOICE"
)

type Invoice struct {
	ID            uuid.UUID `json:"id"`
	PartnershipID uuid.UUID `json:"partnership_id"`
	InvoiceNumber string    `json:"invoice_number"`
	Type          string    `json:"type"`
	Amount        int64     `json:"amount"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`

	// Midtrans
	MidtransOrderID           string `json:"midtrans_order_id,omitempty"`
	MidtransSnapToken         string `json:"midtrans_snap_token,omitempty"`
	MidtransRedirectURL       string `json:"midtrans_redirect_url,omitempty"`
	MidtransPaymentType       string `json:"midtrans_payment_type,omitempty"`
	MidtransTransactionID     string `json:"midtrans_transaction_id,omitempty"`
	MidtransTransactionStatus string `json:"midtrans_transaction_status,omitempty"`

	ProofURL  string     `json:"proof_url,omitempty"`
	PaidAt    *time.Time `json:"paid_at,omitempty"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Relations
	Partnership *Partnership `json:"partnership,omitempty"`
}
