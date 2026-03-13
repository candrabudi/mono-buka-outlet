package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type AffiliatorWithdrawalRepo struct {
	db *sql.DB
}

func NewAffiliatorWithdrawalRepo(db *sql.DB) *AffiliatorWithdrawalRepo {
	return &AffiliatorWithdrawalRepo{db: db}
}

func (r *AffiliatorWithdrawalRepo) Create(ctx context.Context, w *entity.AffiliatorWithdrawal) error {
	query := `
		INSERT INTO affiliator_withdrawals (id, affiliator_id, amount, bank_name, account_number, account_holder, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		w.ID, w.AffiliatorID, w.Amount, w.BankName, w.AccountNumber, w.AccountHolder, w.Status, w.CreatedAt, w.UpdatedAt,
	)
	return err
}

func (r *AffiliatorWithdrawalRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.AffiliatorWithdrawal, error) {
	query := `
		SELECT w.id, w.affiliator_id, w.amount, w.bank_name, w.account_number, w.account_holder,
		       w.status, w.admin_notes, w.processed_by, w.processed_at, w.created_at, w.updated_at,
		       u.id, u.name, u.email
		FROM affiliator_withdrawals w
		JOIN users u ON u.id = w.affiliator_id
		WHERE w.id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	w := &entity.AffiliatorWithdrawal{}
	w.Affiliator = &entity.User{}
	var adminNotes, processedBy sql.NullString
	var processedAt sql.NullTime

	err := row.Scan(
		&w.ID, &w.AffiliatorID, &w.Amount, &w.BankName, &w.AccountNumber, &w.AccountHolder,
		&w.Status, &adminNotes, &processedBy, &processedAt, &w.CreatedAt, &w.UpdatedAt,
		&w.Affiliator.ID, &w.Affiliator.Name, &w.Affiliator.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("permintaan penarikan tidak ditemukan")
		}
		return nil, err
	}
	if adminNotes.Valid {
		w.AdminNotes = adminNotes.String
	}
	if processedBy.Valid {
		uid, _ := uuid.Parse(processedBy.String)
		w.ProcessedBy = &uid
	}
	if processedAt.Valid {
		w.ProcessedAt = &processedAt.Time
	}
	return w, nil
}

func (r *AffiliatorWithdrawalRepo) FindByAffiliatorID(ctx context.Context, affiliatorID uuid.UUID, page, limit int) ([]*entity.AffiliatorWithdrawal, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM affiliator_withdrawals WHERE affiliator_id = $1`
	if err := r.db.QueryRowContext(ctx, countQuery, affiliatorID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, affiliator_id, amount, bank_name, account_number, account_holder,
		       status, admin_notes, processed_at, created_at, updated_at
		FROM affiliator_withdrawals
		WHERE affiliator_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	offset := (page - 1) * limit
	rows, err := r.db.QueryContext(ctx, query, affiliatorID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []*entity.AffiliatorWithdrawal
	for rows.Next() {
		w := &entity.AffiliatorWithdrawal{}
		var adminNotes sql.NullString
		var processedAt sql.NullTime

		err := rows.Scan(
			&w.ID, &w.AffiliatorID, &w.Amount, &w.BankName, &w.AccountNumber, &w.AccountHolder,
			&w.Status, &adminNotes, &processedAt, &w.CreatedAt, &w.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		if adminNotes.Valid {
			w.AdminNotes = adminNotes.String
		}
		if processedAt.Valid {
			w.ProcessedAt = &processedAt.Time
		}
		results = append(results, w)
	}
	return results, total, nil
}

func (r *AffiliatorWithdrawalRepo) FindAll(ctx context.Context, status string, page, limit int) ([]*entity.AffiliatorWithdrawal, int, error) {
	countQuery := `SELECT COUNT(*) FROM affiliator_withdrawals WHERE 1=1`
	dataQuery := `
		SELECT w.id, w.affiliator_id, w.amount, w.bank_name, w.account_number, w.account_holder,
		       w.status, w.admin_notes, w.processed_at, w.created_at, w.updated_at,
		       u.id, u.name, u.email
		FROM affiliator_withdrawals w
		JOIN users u ON u.id = w.affiliator_id
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if status != "" {
		filter := fmt.Sprintf(" AND w.status = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		dataQuery += filter
		args = append(args, status)
		argIdx++
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataQuery += fmt.Sprintf(" ORDER BY w.created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	offset := (page - 1) * limit
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []*entity.AffiliatorWithdrawal
	for rows.Next() {
		w := &entity.AffiliatorWithdrawal{}
		w.Affiliator = &entity.User{}
		var adminNotes sql.NullString
		var processedAt sql.NullTime

		err := rows.Scan(
			&w.ID, &w.AffiliatorID, &w.Amount, &w.BankName, &w.AccountNumber, &w.AccountHolder,
			&w.Status, &adminNotes, &processedAt, &w.CreatedAt, &w.UpdatedAt,
			&w.Affiliator.ID, &w.Affiliator.Name, &w.Affiliator.Email,
		)
		if err != nil {
			return nil, 0, err
		}
		if adminNotes.Valid {
			w.AdminNotes = adminNotes.String
		}
		if processedAt.Valid {
			w.ProcessedAt = &processedAt.Time
		}
		results = append(results, w)
	}
	return results, total, nil
}

func (r *AffiliatorWithdrawalRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status, adminNotes string, processedBy uuid.UUID) error {
	query := `
		UPDATE affiliator_withdrawals
		SET status = $1, admin_notes = $2, processed_by = $3, processed_at = $4, updated_at = $5
		WHERE id = $6
	`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, status, adminNotes, processedBy, now, now, id)
	return err
}

func (r *AffiliatorWithdrawalRepo) GetTotalWithdrawn(ctx context.Context, affiliatorID uuid.UUID) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM affiliator_withdrawals WHERE affiliator_id = $1 AND status IN ('APPROVED', 'TRANSFERRED')`
	var total float64
	err := r.db.QueryRowContext(ctx, query, affiliatorID).Scan(&total)
	return total, err
}

func (r *AffiliatorWithdrawalRepo) GetTotalPending(ctx context.Context, affiliatorID uuid.UUID) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM affiliator_withdrawals WHERE affiliator_id = $1 AND status = 'PENDING'`
	var total float64
	err := r.db.QueryRowContext(ctx, query, affiliatorID).Scan(&total)
	return total, err
}
