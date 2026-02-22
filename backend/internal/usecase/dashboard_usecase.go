package usecase

import (
	"context"

	"github.com/franchise-system/backend/internal/repository"
	"github.com/google/uuid"
)

type DashboardUseCase struct {
	dashboardRepo repository.DashboardRepository
}

func NewDashboardUseCase(dr repository.DashboardRepository) *DashboardUseCase {
	return &DashboardUseCase{dashboardRepo: dr}
}

type DashboardStats struct {
	TotalLeads      int                      `json:"total_leads"`
	ActiveMitra     int                      `json:"active_mitra"`
	TotalInvestment float64                  `json:"total_investment"`
	MonthlyRevenue  float64                  `json:"monthly_revenue"`
	LeadsByStatus   map[string]int           `json:"leads_by_status"`
	RevenueChart    []map[string]interface{} `json:"revenue_chart"`
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
	return stats, nil
}
