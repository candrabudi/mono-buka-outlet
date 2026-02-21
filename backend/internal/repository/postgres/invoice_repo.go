package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type InvoiceRepo struct {
	db *sql.DB
}

func NewInvoiceRepo(db *sql.DB) *InvoiceRepo {
	return &InvoiceRepo{db: db}
}

func (r *InvoiceRepo) Create(ctx context.Context, inv *entity.Invoice) error {
	if inv.ID == uuid.Nil {
		inv.ID = uuid.New()
	}
	inv.CreatedAt = time.Now()
	inv.UpdatedAt = time.Now()

	query := `INSERT INTO invoices (id, partnership_id, invoice_number, type, amount, description, status,
			  midtrans_order_id, midtrans_snap_token, midtrans_redirect_url, expired_at, created_at, updated_at)
			  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := r.db.ExecContext(ctx, query,
		inv.ID, inv.PartnershipID, inv.InvoiceNumber, inv.Type, inv.Amount, inv.Description, inv.Status,
		inv.MidtransOrderID, inv.MidtransSnapToken, inv.MidtransRedirectURL, inv.ExpiredAt, inv.CreatedAt, inv.UpdatedAt)
	return err
}

func (r *InvoiceRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Invoice, error) {
	inv := &entity.Invoice{}
	query := `SELECT id, partnership_id, invoice_number, COALESCE(type,'INVOICE'), amount, description, status,
			  COALESCE(midtrans_order_id,''), COALESCE(midtrans_snap_token,''), COALESCE(midtrans_redirect_url,''),
			  COALESCE(midtrans_payment_type,''), COALESCE(midtrans_transaction_id,''), COALESCE(midtrans_transaction_status,''),
			  COALESCE(proof_url,''), paid_at, expired_at, created_at, updated_at
			  FROM invoices WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&inv.ID, &inv.PartnershipID, &inv.InvoiceNumber, &inv.Type, &inv.Amount, &inv.Description, &inv.Status,
		&inv.MidtransOrderID, &inv.MidtransSnapToken, &inv.MidtransRedirectURL,
		&inv.MidtransPaymentType, &inv.MidtransTransactionID, &inv.MidtransTransactionStatus,
		&inv.ProofURL, &inv.PaidAt, &inv.ExpiredAt, &inv.CreatedAt, &inv.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invoice not found")
	}
	return inv, err
}

func (r *InvoiceRepo) FindByPartnershipID(ctx context.Context, partnershipID uuid.UUID) ([]*entity.Invoice, error) {
	query := `SELECT id, partnership_id, invoice_number, COALESCE(type,'INVOICE'), amount, description, status,
			  COALESCE(midtrans_order_id,''), COALESCE(midtrans_snap_token,''), COALESCE(midtrans_redirect_url,''),
			  COALESCE(midtrans_payment_type,''), COALESCE(midtrans_transaction_id,''), COALESCE(midtrans_transaction_status,''),
			  COALESCE(proof_url,''), paid_at, expired_at, created_at, updated_at
			  FROM invoices WHERE partnership_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, partnershipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []*entity.Invoice
	for rows.Next() {
		inv := &entity.Invoice{}
		if err := rows.Scan(
			&inv.ID, &inv.PartnershipID, &inv.InvoiceNumber, &inv.Type, &inv.Amount, &inv.Description, &inv.Status,
			&inv.MidtransOrderID, &inv.MidtransSnapToken, &inv.MidtransRedirectURL,
			&inv.MidtransPaymentType, &inv.MidtransTransactionID, &inv.MidtransTransactionStatus,
			&inv.ProofURL, &inv.PaidAt, &inv.ExpiredAt, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
			return nil, err
		}
		invoices = append(invoices, inv)
	}
	return invoices, nil
}

func (r *InvoiceRepo) FindByOrderID(ctx context.Context, orderID string) (*entity.Invoice, error) {
	inv := &entity.Invoice{}
	query := `SELECT id, partnership_id, invoice_number, COALESCE(type,'INVOICE'), amount, description, status,
			  COALESCE(midtrans_order_id,''), COALESCE(midtrans_snap_token,''), COALESCE(midtrans_redirect_url,''),
			  COALESCE(midtrans_payment_type,''), COALESCE(midtrans_transaction_id,''), COALESCE(midtrans_transaction_status,''),
			  COALESCE(proof_url,''), paid_at, expired_at, created_at, updated_at
			  FROM invoices WHERE midtrans_order_id = $1`
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&inv.ID, &inv.PartnershipID, &inv.InvoiceNumber, &inv.Type, &inv.Amount, &inv.Description, &inv.Status,
		&inv.MidtransOrderID, &inv.MidtransSnapToken, &inv.MidtransRedirectURL,
		&inv.MidtransPaymentType, &inv.MidtransTransactionID, &inv.MidtransTransactionStatus,
		&inv.ProofURL, &inv.PaidAt, &inv.ExpiredAt, &inv.CreatedAt, &inv.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invoice not found for order: %s", orderID)
	}
	return inv, err
}

func (r *InvoiceRepo) UpdateMidtransStatus(ctx context.Context, orderID string, txnStatus, paymentType, txnID string) error {
	now := time.Now()
	status := entity.InvoiceStatusPending

	switch txnStatus {
	case "capture", "settlement":
		status = entity.InvoiceStatusPaid
	case "expire":
		status = entity.InvoiceStatusExpired
	case "cancel", "deny":
		status = entity.InvoiceStatusFailed
	case "pending":
		status = entity.InvoiceStatusPending
	}

	query := `UPDATE invoices SET
			  midtrans_transaction_status = $1,
			  midtrans_payment_type = $2,
			  midtrans_transaction_id = $3,
			  status = $4,
			  updated_at = $5`
	args := []interface{}{txnStatus, paymentType, txnID, status, now}

	if status == entity.InvoiceStatusPaid {
		query += `, paid_at = $6 WHERE midtrans_order_id = $7`
		args = append(args, now, orderID)
	} else {
		query += ` WHERE midtrans_order_id = $6`
		args = append(args, orderID)
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *InvoiceRepo) ManualApprove(ctx context.Context, id uuid.UUID, proofURL string) error {
	now := time.Now()
	query := `UPDATE invoices SET status = $1, paid_at = $2, midtrans_payment_type = 'manual', proof_url = $3, updated_at = $4 WHERE id = $5 AND status = $6`
	res, err := r.db.ExecContext(ctx, query, entity.InvoiceStatusPaid, now, proofURL, now, id, entity.InvoiceStatusPending)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("invoice not found or already processed")
	}
	return nil
}

func (r *InvoiceRepo) GenerateInvoiceNumber(ctx context.Context) (string, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM invoices`).Scan(&count)
	if err != nil {
		return "", err
	}
	now := time.Now()
	return fmt.Sprintf("INV-%s-%04d", now.Format("20060102"), count+1), nil
}
