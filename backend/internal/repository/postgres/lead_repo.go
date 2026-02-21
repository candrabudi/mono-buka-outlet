package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type LeadRepo struct {
	db *sql.DB
}

func NewLeadRepo(db *sql.DB) *LeadRepo {
	return &LeadRepo{db: db}
}

func (r *LeadRepo) Create(ctx context.Context, lead *entity.Lead) error {
	query := `INSERT INTO leads (id, brand_id, sales_id, full_name, email, phone, status, progress_percentage, notes, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query,
		lead.ID, lead.BrandID, lead.SalesID, lead.FullName, lead.Email, lead.Phone,
		lead.Status, lead.ProgressPercentage, lead.Notes, lead.CreatedAt, lead.UpdatedAt,
	)
	return err
}

func (r *LeadRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Lead, error) {
	lead := &entity.Lead{}

	query := `SELECT id, brand_id, sales_id, full_name, email, phone, status, progress_percentage, notes, created_at, updated_at
			  FROM leads WHERE id = $1 AND deleted_at IS NULL`

	var salesID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&lead.ID, &lead.BrandID, &salesID, &lead.FullName, &lead.Email, &lead.Phone,
		&lead.Status, &lead.ProgressPercentage, &lead.Notes, &lead.CreatedAt, &lead.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("lead not found")
	}
	if salesID.Valid {
		sid, _ := uuid.Parse(salesID.String)
		lead.SalesID = &sid
	}
	return lead, err
}

func (r *LeadRepo) FindAll(ctx context.Context, brandID *uuid.UUID, status string, page, limit int) ([]*entity.Lead, int, error) {
	var total int
	countQuery := "SELECT COUNT(*) FROM leads WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIdx := 1

	if brandID != nil {
		countQuery += fmt.Sprintf(" AND brand_id = $%d", argIdx)
		args = append(args, *brandID)
		argIdx++
	}
	if status != "" {
		countQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, brand_id, sales_id, full_name, email, phone, status, progress_percentage, notes, created_at, updated_at
			  FROM leads WHERE deleted_at IS NULL`
	queryArgs := []interface{}{}
	queryArgIdx := 1

	if brandID != nil {
		query += fmt.Sprintf(" AND brand_id = $%d", queryArgIdx)
		queryArgs = append(queryArgs, *brandID)
		queryArgIdx++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", queryArgIdx)
		queryArgs = append(queryArgs, status)
		queryArgIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", queryArgIdx, queryArgIdx+1)
	queryArgs = append(queryArgs, limit, (page-1)*limit)

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var leads []*entity.Lead
	for rows.Next() {
		lead := &entity.Lead{}
		var salesID sql.NullString
		err := rows.Scan(
			&lead.ID, &lead.BrandID, &salesID, &lead.FullName, &lead.Email, &lead.Phone,
			&lead.Status, &lead.ProgressPercentage, &lead.Notes, &lead.CreatedAt, &lead.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		if salesID.Valid {
			sid, _ := uuid.Parse(salesID.String)
			lead.SalesID = &sid
		}
		leads = append(leads, lead)
	}

	return leads, total, nil
}

func (r *LeadRepo) FindByBrandGrouped(ctx context.Context, brandID *uuid.UUID) (map[string][]*entity.Lead, error) {
	query := `SELECT id, brand_id, sales_id, full_name, email, phone, status, progress_percentage, notes, created_at, updated_at
			  FROM leads WHERE deleted_at IS NULL`
	args := []interface{}{}

	if brandID != nil {
		query += " AND brand_id = $1"
		args = append(args, *brandID)
	}

	query += " ORDER BY updated_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grouped := make(map[string][]*entity.Lead)
	for _, status := range entity.LeadStatuses() {
		grouped[status] = []*entity.Lead{}
	}

	for rows.Next() {
		lead := &entity.Lead{}
		var salesID sql.NullString
		err := rows.Scan(
			&lead.ID, &lead.BrandID, &salesID, &lead.FullName, &lead.Email, &lead.Phone,
			&lead.Status, &lead.ProgressPercentage, &lead.Notes, &lead.CreatedAt, &lead.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if salesID.Valid {
			sid, _ := uuid.Parse(salesID.String)
			lead.SalesID = &sid
		}
		grouped[lead.Status] = append(grouped[lead.Status], lead)
	}

	return grouped, nil
}

func (r *LeadRepo) Update(ctx context.Context, lead *entity.Lead) error {
	query := `UPDATE leads SET brand_id = $1, sales_id = $2, full_name = $3, email = $4, phone = $5, 
			  status = $6, progress_percentage = $7, notes = $8, updated_at = $9 
			  WHERE id = $10 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query,
		lead.BrandID, lead.SalesID, lead.FullName, lead.Email, lead.Phone,
		lead.Status, lead.ProgressPercentage, lead.Notes, time.Now(), lead.ID,
	)
	return err
}

func (r *LeadRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, progress int) error {
	query := `UPDATE leads SET status = $1, progress_percentage = $2, updated_at = $3 WHERE id = $4 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, status, progress, time.Now(), id)
	return err
}

func (r *LeadRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := "UPDATE leads SET deleted_at = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
