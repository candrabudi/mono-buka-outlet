package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/google/uuid"
)

type RevenueUseCase struct {
	revenueRepo     repository.RevenueRepository
	partnershipRepo repository.PartnershipRepository
}

func NewRevenueUseCase(rr repository.RevenueRepository, pr repository.PartnershipRepository) *RevenueUseCase {
	return &RevenueUseCase{revenueRepo: rr, partnershipRepo: pr}
}

type CreateRevenueRequest struct {
	PartnershipID string  `json:"partnership_id" binding:"required"`
	Month         string  `json:"month" binding:"required"`
	Revenue       float64 `json:"revenue" binding:"required"`
	Expense       float64 `json:"expense"`
}

func (uc *RevenueUseCase) Create(ctx context.Context, req CreateRevenueRequest) (*entity.Revenue, error) {
	pid, _ := uuid.Parse(req.PartnershipID)
	partnership, err := uc.partnershipRepo.FindByID(ctx, pid)
	if err != nil {
		return nil, err
	}
	profit := req.Revenue - req.Expense
	companyShare := profit * 0.30
	mitraShare := profit - companyShare
	var brandID2 uuid.UUID
	if partnership.BrandID != nil {
		brandID2 = *partnership.BrandID
	}
	rev := &entity.Revenue{
		ID: uuid.New(), PartnershipID: pid, BrandID: brandID2,
		Month: req.Month, Revenue: req.Revenue, Expense: req.Expense,
		Profit: profit, CompanyShare: companyShare, MitraShare: mitraShare,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := uc.revenueRepo.Create(ctx, rev); err != nil {
		return nil, fmt.Errorf("failed to create revenue: %w", err)
	}
	return rev, nil
}

func (uc *RevenueUseCase) GetByPartnership(ctx context.Context, pid uuid.UUID) ([]*entity.Revenue, error) {
	return uc.revenueRepo.FindByPartnershipID(ctx, pid)
}

func (uc *RevenueUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Revenue, error) {
	return uc.revenueRepo.FindByID(ctx, id)
}
