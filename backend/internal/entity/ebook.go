package entity

import (
	"time"

	"github.com/google/uuid"
)

// Ebook status constants
const (
	EbookOrderStatusPending  = "PENDING"
	EbookOrderStatusPaid     = "PAID"
	EbookOrderStatusExpired  = "EXPIRED"
	EbookOrderStatusFailed   = "FAILED"
	EbookOrderStatusCanceled = "CANCELED"

	DownloadStatusNone      = "NONE"
	DownloadStatusRequested = "REQUESTED"
	DownloadStatusApproved  = "APPROVED"
	DownloadStatusRejected  = "REJECTED"
)

type EbookCategory struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Ebook struct {
	ID          uuid.UUID        `json:"id"`
	Title       string           `json:"title"`
	Slug        string           `json:"slug"`
	Description string           `json:"description"`
	Author      string           `json:"author"`
	CoverURL    string           `json:"cover_url"`
	FileURL     string           `json:"file_url"`
	Price       int64            `json:"price"`
	IsActive    bool             `json:"is_active"`
	TotalSold   int              `json:"total_sold"`
	Categories  []*EbookCategory `json:"categories,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type EbookOrder struct {
	ID          uuid.UUID `json:"id"`
	EbookID     uuid.UUID `json:"ebook_id"`
	UserID      uuid.UUID `json:"user_id"`
	OrderNumber string    `json:"order_number"`
	Amount      int64     `json:"amount"`
	Status      string    `json:"status"`

	// Download approval workflow
	DownloadStatus      string     `json:"download_status"`
	DownloadRequestedAt *time.Time `json:"download_requested_at,omitempty"`
	DownloadApprovedAt  *time.Time `json:"download_approved_at,omitempty"`
	DownloadNote        string     `json:"download_note"`

	// Midtrans
	MidtransOrderID           string `json:"midtrans_order_id,omitempty"`
	MidtransSnapToken         string `json:"midtrans_snap_token,omitempty"`
	MidtransRedirectURL       string `json:"midtrans_redirect_url,omitempty"`
	MidtransPaymentType       string `json:"midtrans_payment_type,omitempty"`
	MidtransTransactionID     string `json:"midtrans_transaction_id,omitempty"`
	MidtransTransactionStatus string `json:"midtrans_transaction_status,omitempty"`

	PaidAt    *time.Time `json:"paid_at,omitempty"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`

	// Manual payment proof
	PaymentProofURL string `json:"payment_proof_url,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Ebook *Ebook `json:"ebook,omitempty"`
	User  *User  `json:"user,omitempty"`
}
