package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/google/uuid"
)

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

// GetByMitra returns all agreements across all partnerships belonging to a mitra
func (uc *AgreementUseCase) GetByMitra(ctx context.Context, mitraID uuid.UUID) ([]*entity.Agreement, error) {
	partnerships, err := uc.partnershipRepo.FindByMitraID(ctx, mitraID)
	if err != nil {
		return nil, err
	}
	var all []*entity.Agreement
	for _, p := range partnerships {
		agreements, err := uc.agreementRepo.FindByPartnershipID(ctx, p.ID)
		if err != nil {
			continue
		}
		all = append(all, agreements...)
	}
	if all == nil {
		all = []*entity.Agreement{}
	}
	return all, nil
}
