package usecase

import (
	"context"
	"database/sql"

	"github.com/franchise-system/backend/internal/repository"
	"github.com/google/uuid"
)

type DashboardUseCase struct {
	dashboardRepo repository.DashboardRepository
	db            *sql.DB
}

func NewDashboardUseCase(dr repository.DashboardRepository, db *sql.DB) *DashboardUseCase {
	return &DashboardUseCase{dashboardRepo: dr, db: db}
}

type DashboardStats struct {
	TotalLeads      int                      `json:"total_leads"`
	ActiveMitra     int                      `json:"active_mitra"`
	TotalInvestment float64                  `json:"total_investment"`
	MonthlyRevenue  float64                  `json:"monthly_revenue"`
	LeadsByStatus   map[string]int           `json:"leads_by_status"`
	RevenueChart    []map[string]interface{} `json:"revenue_chart"`

	TotalOutlets          int            `json:"total_outlets"`
	TotalPartnerships     int            `json:"total_partnerships"`
	PartnershipsByStatus  map[string]int `json:"partnerships_by_status"`
	PendingApplications   int            `json:"pending_applications"`
	TotalApplications     int            `json:"total_applications"`
	ApplicationsByStatus  map[string]int `json:"applications_by_status"`
	PendingInvoices       int            `json:"pending_invoices"`
	PaidInvoices          int            `json:"paid_invoices"`
	TotalIncomeAmount     float64        `json:"total_income_amount"`
	PendingInvoiceAmount  float64        `json:"pending_invoice_amount"`
	TotalPaymentsVerified float64        `json:"total_payments_verified"`
	TotalMitra            int            `json:"total_mitra"`
	TotalUsers            int            `json:"total_users"`
}

func (uc *DashboardUseCase) GetStats(ctx context.Context, brandID *uuid.UUID) (*DashboardStats, error) {
	stats := &DashboardStats{}
	var err error
	stats.TotalLeads, err = uc.dashboardRepo.GetTotalLeads(ctx, brandID)
	if err != nil {
		return nil, err
	}
	stats.ActiveMitra, err = uc.dashboardRepo.GetActiveMitra(ctx, brandID)
	if err != nil {
		return nil, err
	}
	stats.TotalInvestment, err = uc.dashboardRepo.GetTotalInvestment(ctx, brandID)
	if err != nil {
		return nil, err
	}
	stats.MonthlyRevenue, err = uc.dashboardRepo.GetMonthlyRevenue(ctx, brandID, "")
	if err != nil {
		return nil, err
	}
	stats.LeadsByStatus, err = uc.dashboardRepo.GetLeadsByStatus(ctx, brandID)
	if err != nil {
		return nil, err
	}
	stats.RevenueChart, err = uc.dashboardRepo.GetRevenueChart(ctx, brandID, 6)
	if err != nil {
		return nil, err
	}

	// Extra stats from tables (ignore errors, default to 0)
	uc.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM outlets WHERE deleted_at IS NULL").Scan(&stats.TotalOutlets)
	uc.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM partnerships WHERE deleted_at IS NULL").Scan(&stats.TotalPartnerships)
	uc.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM partnership_applications").Scan(&stats.TotalApplications)
	uc.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM partnership_applications WHERE status = 'PENDING'").Scan(&stats.PendingApplications)
	uc.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM invoices WHERE status = 'PENDING'").Scan(&stats.PendingInvoices)
	uc.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM invoices WHERE status = 'PAID'").Scan(&stats.PaidInvoices)
	uc.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE role = 'mitra' AND is_active = true").Scan(&stats.TotalMitra)
	uc.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE is_active = true").Scan(&stats.TotalUsers)

	// Monetary — cast BIGINT to NUMERIC for safe float64 scan
	var incomeAmt, pendingAmt, verifiedAmt int64
	uc.db.QueryRowContext(ctx, "SELECT COALESCE(SUM(amount),0)::BIGINT FROM invoices WHERE status = 'PAID'").Scan(&incomeAmt)
	uc.db.QueryRowContext(ctx, "SELECT COALESCE(SUM(amount),0)::BIGINT FROM invoices WHERE status = 'PENDING'").Scan(&pendingAmt)
	uc.db.QueryRowContext(ctx, "SELECT COALESCE(SUM(amount),0)::BIGINT FROM payments WHERE status = 'VERIFIED'").Scan(&verifiedAmt)
	stats.TotalIncomeAmount = float64(incomeAmt)
	stats.PendingInvoiceAmount = float64(pendingAmt)
	stats.TotalPaymentsVerified = float64(verifiedAmt)

	// Partnerships by status
	stats.PartnershipsByStatus = map[string]int{}
	pRows, pErr := uc.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM partnerships WHERE deleted_at IS NULL GROUP BY status")
	if pErr == nil {
		defer pRows.Close()
		for pRows.Next() {
			var s string
			var c int
			pRows.Scan(&s, &c)
			stats.PartnershipsByStatus[s] = c
		}
	}

	// Applications by status
	stats.ApplicationsByStatus = map[string]int{}
	aRows, aErr := uc.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM partnership_applications GROUP BY status")
	if aErr == nil {
		defer aRows.Close()
		for aRows.Next() {
			var s string
			var c int
			aRows.Scan(&s, &c)
			stats.ApplicationsByStatus[s] = c
		}
	}

	return stats, nil
}
