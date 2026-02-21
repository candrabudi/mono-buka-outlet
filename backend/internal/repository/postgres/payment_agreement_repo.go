package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type PaymentRepo struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) Create(ctx context.Context, p *entity.Payment) error {
	query := `INSERT INTO payments (id, partnership_id, brand_id, type, amount, proof_url, verified_status, notes, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.PartnershipID, p.BrandID, p.Type, p.Amount, p.ProofURL, p.VerifiedStatus, p.Notes, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *PaymentRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	p := &entity.Payment{}
	query := `SELECT id, partnership_id, brand_id, type, amount, proof_url, verified_status, verified_by, verified_at, notes, created_at, updated_at FROM payments WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.PartnershipID, &p.BrandID, &p.Type, &p.Amount, &p.ProofURL, &p.VerifiedStatus, &p.VerifiedBy, &p.VerifiedAt, &p.Notes, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}
	return p, err
}

func (r *PaymentRepo) FindByPartnershipID(ctx context.Context, pid uuid.UUID) ([]*entity.Payment, error) {
	query := `SELECT id, partnership_id, brand_id, type, amount, proof_url, verified_status, verified_by, verified_at, notes, created_at, updated_at FROM payments WHERE partnership_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []*entity.Payment
	for rows.Next() {
		p := &entity.Payment{}
		if err := rows.Scan(&p.ID, &p.PartnershipID, &p.BrandID, &p.Type, &p.Amount, &p.ProofURL, &p.VerifiedStatus, &p.VerifiedBy, &p.VerifiedAt, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *PaymentRepo) Update(ctx context.Context, p *entity.Payment) error {
	_, err := r.db.ExecContext(ctx, `UPDATE payments SET amount=$1, proof_url=$2, notes=$3, updated_at=$4 WHERE id=$5 AND deleted_at IS NULL`, p.Amount, p.ProofURL, p.Notes, time.Now(), p.ID)
	return err
}

func (r *PaymentRepo) Verify(ctx context.Context, id uuid.UUID, status string, verifiedBy uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE payments SET verified_status=$1, verified_by=$2, verified_at=$3, updated_at=$3 WHERE id=$4 AND deleted_at IS NULL`, status, verifiedBy, time.Now(), id)
	return err
}
