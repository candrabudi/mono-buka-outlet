package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/google/uuid"
)

type PartnershipUseCase struct {
	partnershipRepo repository.PartnershipRepository
	logRepo         repository.ActivityLogRepository
}

func NewPartnershipUseCase(pr repository.PartnershipRepository, lr repository.ActivityLogRepository) *PartnershipUseCase {
	return &PartnershipUseCase{partnershipRepo: pr, logRepo: lr}
}

type CreatePartnershipRequest struct {
	LeadID    string `json:"lead_id"`
	BrandID   string `json:"brand_id"`
	MitraID   string `json:"mitra_id" binding:"required"`
	LeaderID  string `json:"leader_id"`
	OutletID  string `json:"outlet_id"`
	PackageID string `json:"package_id"`
}

func (uc *PartnershipUseCase) Create(ctx context.Context, req CreatePartnershipRequest) (*entity.Partnership, error) {
	mitraID, err := uuid.Parse(req.MitraID)
	if err != nil {
		return nil, fmt.Errorf("invalid mitra_id")
	}

	p := &entity.Partnership{
		ID: uuid.New(), MitraID: mitraID,
		ProgressPercentage: 0, Status: entity.PartnershipStatusPending,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	if req.LeadID != "" {
		leadID, _ := uuid.Parse(req.LeadID)
		p.LeadID = &leadID
	}
	if req.BrandID != "" {
		brandID, _ := uuid.Parse(req.BrandID)
		p.BrandID = &brandID
	}
	if req.LeaderID != "" {
		leaderID, err := uuid.Parse(req.LeaderID)
		if err == nil {
			p.LeaderID = &leaderID
		}
	}
	if req.OutletID != "" {
		outletID, err := uuid.Parse(req.OutletID)
		if err == nil {
			p.OutletID = &outletID
		}
	}
	if req.PackageID != "" {
		packageID, err := uuid.Parse(req.PackageID)
		if err == nil {
			p.PackageID = &packageID
		}
	}

	if err := uc.partnershipRepo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("failed to create partnership: %w", err)
	}
	return p, nil
}

func (uc *PartnershipUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Partnership, error) {
	return uc.partnershipRepo.FindByID(ctx, id)
}

func (uc *PartnershipUseCase) GetAll(ctx context.Context, brandID, mitraID *uuid.UUID, page, limit int) ([]*entity.Partnership, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return uc.partnershipRepo.FindAll(ctx, brandID, mitraID, page, limit)
}

func (uc *PartnershipUseCase) GetByMitra(ctx context.Context, mitraID uuid.UUID) ([]*entity.Partnership, error) {
	return uc.partnershipRepo.FindByMitraID(ctx, mitraID)
}

// Agreement Use Case
type AgreementUseCase struct {
	agreementRepo   repository.AgreementRepository
	partnershipRepo repository.PartnershipRepository
	logRepo         repository.ActivityLogRepository
}

func NewAgreementUseCase(ar repository.AgreementRepository, pr repository.PartnershipRepository, lr repository.ActivityLogRepository) *AgreementUseCase {
	return &AgreementUseCase{agreementRepo: ar, partnershipRepo: pr, logRepo: lr}
}

type CreateAgreementRequest struct {
	PartnershipID string `json:"partnership_id" binding:"required"`
	Title         string `json:"title"`
	Type          string `json:"type"`
	FileURL       string `json:"file_url"`
}

func (uc *AgreementUseCase) Create(ctx context.Context, req CreateAgreementRequest) (*entity.Agreement, error) {
	pid, _ := uuid.Parse(req.PartnershipID)
	partnership, err := uc.partnershipRepo.FindByID(ctx, pid)
	if err != nil {
		return nil, err
	}
	existing, _ := uc.agreementRepo.FindByPartnershipID(ctx, pid)
	version := 1
	if len(existing) > 0 {
		version = existing[0].Version + 1
	}
	var brandID uuid.UUID
	if partnership.BrandID != nil {
		brandID = *partnership.BrandID
	}
	agrType := req.Type
	if agrType == "" {
		agrType = entity.AgreementTypeContract
	}
	a := &entity.Agreement{
		ID: uuid.New(), PartnershipID: pid, BrandID: brandID,
		Title: req.Title, Type: agrType, FileURL: req.FileURL, Version: version, Status: entity.AgreementStatusDraft,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := uc.agreementRepo.Create(ctx, a); err != nil {
		return nil, fmt.Errorf("failed to create agreement: %w", err)
	}
	return a, nil
}

func (uc *AgreementUseCase) Sign(ctx context.Context, id uuid.UUID, performedBy uuid.UUID) error {
	a, err := uc.agreementRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := uc.agreementRepo.Sign(ctx, id); err != nil {
		return fmt.Errorf("failed to sign agreement: %w", err)
	}
	// Auto update partnership progress: 50% agreement signed
	uc.partnershipRepo.UpdateProgress(ctx, a.PartnershipID, 50, entity.PartnershipStatusAgreementSigned)
	uc.logRepo.Create(ctx, &entity.ActivityLog{
		ID: uuid.New(), EntityType: "agreement", EntityID: id,
		Action: "AGREEMENT_SIGNED", Description: "Agreement signed",
		PerformedBy: performedBy, CreatedAt: time.Now(),
	})
	return nil
}

func (uc *AgreementUseCase) GetByPartnership(ctx context.Context, pid uuid.UUID) ([]*entity.Agreement, error) {
	return uc.agreementRepo.FindByPartnershipID(ctx, pid)
}

func (uc *AgreementUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Agreement, error) {
	return uc.agreementRepo.FindByID(ctx, id)
}

// Dashboard Use Case
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
