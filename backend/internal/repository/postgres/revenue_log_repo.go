package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type RevenueRepo struct {
	db *sql.DB
}

func NewRevenueRepo(db *sql.DB) *RevenueRepo {
	return &RevenueRepo{db: db}
}

func (r *RevenueRepo) Create(ctx context.Context, rev *entity.Revenue) error {
	query := `INSERT INTO revenues (id, partnership_id, brand_id, month, revenue, expense, profit, company_share, mitra_share, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.db.ExecContext(ctx, query, rev.ID, rev.PartnershipID, rev.BrandID, rev.Month, rev.Revenue, rev.Expense, rev.Profit, rev.CompanyShare, rev.MitraShare, rev.CreatedAt, rev.UpdatedAt)
	return err
}

func (r *RevenueRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Revenue, error) {
	rev := &entity.Revenue{}
	err := r.db.QueryRowContext(ctx, `SELECT id, partnership_id, brand_id, month, revenue, expense, profit, company_share, mitra_share, created_at, updated_at FROM revenues WHERE id=$1 AND deleted_at IS NULL`, id).Scan(&rev.ID, &rev.PartnershipID, &rev.BrandID, &rev.Month, &rev.Revenue, &rev.Expense, &rev.Profit, &rev.CompanyShare, &rev.MitraShare, &rev.CreatedAt, &rev.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("revenue not found")
	}
	return rev, err
}

func (r *RevenueRepo) FindByPartnershipID(ctx context.Context, pid uuid.UUID) ([]*entity.Revenue, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, partnership_id, brand_id, month, revenue, expense, profit, company_share, mitra_share, created_at, updated_at FROM revenues WHERE partnership_id=$1 AND deleted_at IS NULL ORDER BY month DESC`, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*entity.Revenue
	for rows.Next() {
		rev := &entity.Revenue{}
		if err := rows.Scan(&rev.ID, &rev.PartnershipID, &rev.BrandID, &rev.Month, &rev.Revenue, &rev.Expense, &rev.Profit, &rev.CompanyShare, &rev.MitraShare, &rev.CreatedAt, &rev.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, rev)
	}
	return list, nil
}

func (r *RevenueRepo) FindByBrandID(ctx context.Context, bid uuid.UUID, month string) ([]*entity.Revenue, error) {
	query := `SELECT id, partnership_id, brand_id, month, revenue, expense, profit, company_share, mitra_share, created_at, updated_at FROM revenues WHERE brand_id=$1 AND deleted_at IS NULL`
	args := []interface{}{bid}
	if month != "" {
		query += " AND month=$2"
		args = append(args, month)
	}
	query += " ORDER BY month DESC"
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*entity.Revenue
	for rows.Next() {
		rev := &entity.Revenue{}
		if err := rows.Scan(&rev.ID, &rev.PartnershipID, &rev.BrandID, &rev.Month, &rev.Revenue, &rev.Expense, &rev.Profit, &rev.CompanyShare, &rev.MitraShare, &rev.CreatedAt, &rev.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, rev)
	}
	return list, nil
}

func (r *RevenueRepo) Update(ctx context.Context, rev *entity.Revenue) error {
	_, err := r.db.ExecContext(ctx, `UPDATE revenues SET revenue=$1, expense=$2, profit=$3, company_share=$4, mitra_share=$5, updated_at=$6 WHERE id=$7`, rev.Revenue, rev.Expense, rev.Profit, rev.CompanyShare, rev.MitraShare, time.Now(), rev.ID)
	return err
}

// Activity Log Repository
type ActivityLogRepo struct {
	db *sql.DB
}

func NewActivityLogRepo(db *sql.DB) *ActivityLogRepo {
	return &ActivityLogRepo{db: db}
}

func (r *ActivityLogRepo) Create(ctx context.Context, l *entity.ActivityLog) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO activity_logs (id, entity_type, entity_id, action, description, old_value, new_value, performed_by, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, l.ID, l.EntityType, l.EntityID, l.Action, l.Description, l.OldValue, l.NewValue, l.PerformedBy, l.CreatedAt)
	return err
}

func (r *ActivityLogRepo) FindByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*entity.ActivityLog, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT al.id, al.entity_type, al.entity_id, al.action, al.description, al.old_value, al.new_value, al.performed_by, al.created_at, u.name FROM activity_logs al JOIN users u ON al.performed_by = u.id WHERE al.entity_type=$1 AND al.entity_id=$2 ORDER BY al.created_at DESC`, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*entity.ActivityLog
	for rows.Next() {
		l := &entity.ActivityLog{Performer: &entity.User{}}
		if err := rows.Scan(&l.ID, &l.EntityType, &l.EntityID, &l.Action, &l.Description, &l.OldValue, &l.NewValue, &l.PerformedBy, &l.CreatedAt, &l.Performer.Name); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, nil
}
