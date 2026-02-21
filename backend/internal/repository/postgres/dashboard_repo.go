package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DashboardRepo struct {
	db *sql.DB
}

func NewDashboardRepo(db *sql.DB) *DashboardRepo {
	return &DashboardRepo{db: db}
}

func (r *DashboardRepo) GetTotalLeads(ctx context.Context, brandID *uuid.UUID) (int, error) {
	q := "SELECT COUNT(*) FROM leads WHERE deleted_at IS NULL"
	args := []interface{}{}
	if brandID != nil {
		q += " AND brand_id = $1"
		args = append(args, *brandID)
	}
	var total int
	err := r.db.QueryRowContext(ctx, q, args...).Scan(&total)
	return total, err
}

func (r *DashboardRepo) GetActiveMitra(ctx context.Context, brandID *uuid.UUID) (int, error) {
	q := "SELECT COUNT(DISTINCT mitra_id) FROM partnerships WHERE deleted_at IS NULL AND status IN ('RUNNING', 'DP_VERIFIED', 'AGREEMENT_SIGNED', 'DEVELOPMENT')"
	args := []interface{}{}
	if brandID != nil {
		q += " AND brand_id = $1"
		args = append(args, *brandID)
	}
	var total int
	err := r.db.QueryRowContext(ctx, q, args...).Scan(&total)
	return total, err
}

func (r *DashboardRepo) GetTotalInvestment(ctx context.Context, brandID *uuid.UUID) (float64, error) {
	q := "SELECT COALESCE(SUM(amount), 0) FROM payments WHERE deleted_at IS NULL AND verified_status = 'VERIFIED'"
	args := []interface{}{}
	if brandID != nil {
		q += " AND brand_id = $1"
		args = append(args, *brandID)
	}
	var total float64
	err := r.db.QueryRowContext(ctx, q, args...).Scan(&total)
	return total, err
}

func (r *DashboardRepo) GetMonthlyRevenue(ctx context.Context, brandID *uuid.UUID, month string) (float64, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}
	q := "SELECT COALESCE(SUM(revenue), 0) FROM revenues WHERE deleted_at IS NULL AND month = $1"
	args := []interface{}{month}
	if brandID != nil {
		q += fmt.Sprintf(" AND brand_id = $%d", len(args)+1)
		args = append(args, *brandID)
	}
	var total float64
	err := r.db.QueryRowContext(ctx, q, args...).Scan(&total)
	return total, err
}

func (r *DashboardRepo) GetLeadsByStatus(ctx context.Context, brandID *uuid.UUID) (map[string]int, error) {
	q := "SELECT status, COUNT(*) FROM leads WHERE deleted_at IS NULL"
	args := []interface{}{}
	if brandID != nil {
		q += " AND brand_id = $1"
		args = append(args, *brandID)
	}
	q += " GROUP BY status"
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		result[status] = count
	}
	return result, nil
}

func (r *DashboardRepo) GetRevenueChart(ctx context.Context, brandID *uuid.UUID, months int) ([]map[string]interface{}, error) {
	if months <= 0 {
		months = 6
	}
	q := `SELECT month, SUM(revenue) as total_revenue, SUM(profit) as total_profit FROM revenues WHERE deleted_at IS NULL AND month >= $1`
	startMonth := time.Now().AddDate(0, -months, 0).Format("2006-01")
	args := []interface{}{startMonth}
	if brandID != nil {
		q += fmt.Sprintf(" AND brand_id = $%d", len(args)+1)
		args = append(args, *brandID)
	}
	q += " GROUP BY month ORDER BY month ASC"
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []map[string]interface{}
	for rows.Next() {
		var month string
		var revenue, profit float64
		if err := rows.Scan(&month, &revenue, &profit); err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{
			"month":   month,
			"revenue": revenue,
			"profit":  profit,
		})
	}
	return result, nil
}
