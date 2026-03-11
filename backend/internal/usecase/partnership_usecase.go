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
	BrandID      string `json:"brand_id"`
	MitraID      string `json:"mitra_id" binding:"required"`
	AffiliatorID string `json:"affiliator_id"`
	OutletID     string `json:"outlet_id"`
	PackageID    string `json:"package_id"`
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

	if req.BrandID != "" {
		brandID, _ := uuid.Parse(req.BrandID)
		p.BrandID = &brandID
	}
	if req.AffiliatorID != "" {
		affiliatorID, err := uuid.Parse(req.AffiliatorID)
		if err == nil {
			p.AffiliatorID = &affiliatorID
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

type UpdatePartnershipStatusRequest struct {
	Status   string `json:"status" binding:"required"`
	Progress *int   `json:"progress_percentage"`
}

var validStatuses = map[string]int{
	entity.PartnershipStatusPending:         0,
	entity.PartnershipStatusDPVerified:      25,
	entity.PartnershipStatusAgreementSigned: 50,
	entity.PartnershipStatusDevelopment:     75,
	entity.PartnershipStatusRunning:         90,
	entity.PartnershipStatusCompleted:       100,
}

func (uc *PartnershipUseCase) UpdateStatus(ctx context.Context, id uuid.UUID, req UpdatePartnershipStatusRequest) (*entity.Partnership, error) {
	// Validate status
	defaultProgress, ok := validStatuses[req.Status]
	if !ok {
		return nil, fmt.Errorf("status tidak valid: %s", req.Status)
	}

	// Use custom progress if provided, otherwise use default
	progress := defaultProgress
	if req.Progress != nil {
		if *req.Progress < 0 || *req.Progress > 100 {
			return nil, fmt.Errorf("progress harus antara 0-100")
		}
		progress = *req.Progress
	}

	// Check partnership exists
	p, err := uc.partnershipRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("partnership tidak ditemukan")
	}

	// Update
	if err := uc.partnershipRepo.UpdateProgress(ctx, id, progress, req.Status); err != nil {
		return nil, fmt.Errorf("gagal mengupdate status: %w", err)
	}

	p.Status = req.Status
	p.ProgressPercentage = progress
	return p, nil
}
