package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type AgreementRepo struct {
	db *sql.DB
}

func NewAgreementRepo(db *sql.DB) *AgreementRepo {
	return &AgreementRepo{db: db}
}

func (r *AgreementRepo) Create(ctx context.Context, a *entity.Agreement) error {
	query := `INSERT INTO agreements (id, partnership_id, brand_id, title, type, file_url, version, status, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.ExecContext(ctx, query, a.ID, a.PartnershipID, a.BrandID, a.Title, a.Type, a.FileURL, a.Version, a.Status, a.CreatedAt, a.UpdatedAt)
	return err
}

func (r *AgreementRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Agreement, error) {
	a := &entity.Agreement{}
	err := r.db.QueryRowContext(ctx, `SELECT id, partnership_id, brand_id, title, type, file_url, version, status, signed_at, created_at, updated_at FROM agreements WHERE id=$1 AND deleted_at IS NULL`, id).Scan(&a.ID, &a.PartnershipID, &a.BrandID, &a.Title, &a.Type, &a.FileURL, &a.Version, &a.Status, &a.SignedAt, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("agreement not found")
	}
	return a, err
}

func (r *AgreementRepo) FindByPartnershipID(ctx context.Context, pid uuid.UUID) ([]*entity.Agreement, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, partnership_id, brand_id, title, type, file_url, version, status, signed_at, created_at, updated_at FROM agreements WHERE partnership_id=$1 AND deleted_at IS NULL ORDER BY version DESC`, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*entity.Agreement
	for rows.Next() {
		a := &entity.Agreement{}
		if err := rows.Scan(&a.ID, &a.PartnershipID, &a.BrandID, &a.Title, &a.Type, &a.FileURL, &a.Version, &a.Status, &a.SignedAt, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *AgreementRepo) Update(ctx context.Context, a *entity.Agreement) error {
	_, err := r.db.ExecContext(ctx, `UPDATE agreements SET title=$1, type=$2, file_url=$3, status=$4, updated_at=$5 WHERE id=$6 AND deleted_at IS NULL`, a.Title, a.Type, a.FileURL, a.Status, time.Now(), a.ID)
	return err
}

func (r *AgreementRepo) Sign(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE agreements SET status='SIGNED', signed_at=$1, updated_at=$1 WHERE id=$2 AND deleted_at IS NULL`, time.Now(), id)
	return err
}
