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
