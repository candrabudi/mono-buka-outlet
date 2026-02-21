package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/google/uuid"
)

type LeadUseCase struct {
	leadRepo repository.LeadRepository
	logRepo  repository.ActivityLogRepository
}

func NewLeadUseCase(lr repository.LeadRepository, alr repository.ActivityLogRepository) *LeadUseCase {
	return &LeadUseCase{leadRepo: lr, logRepo: alr}
}

type CreateLeadRequest struct {
	BrandID  string `json:"brand_id" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	Notes    string `json:"notes"`
}

type UpdateLeadStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (uc *LeadUseCase) Create(ctx context.Context, req CreateLeadRequest, salesID *uuid.UUID) (*entity.Lead, error) {
	brandID, err := uuid.Parse(req.BrandID)
	if err != nil {
		return nil, fmt.Errorf("invalid brand_id")
	}
	lead := &entity.Lead{
		ID: uuid.New(), BrandID: brandID, SalesID: salesID,
		FullName: req.FullName, Email: req.Email, Phone: req.Phone,
		Status: entity.LeadStatusNew, ProgressPercentage: 0,
		Notes: req.Notes, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := uc.leadRepo.Create(ctx, lead); err != nil {
		return nil, fmt.Errorf("failed to create lead: %w", err)
	}
	return lead, nil
}

func (uc *LeadUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Lead, error) {
	return uc.leadRepo.FindByID(ctx, id)
}

func (uc *LeadUseCase) GetAll(ctx context.Context, brandID *uuid.UUID, status string, page, limit int) ([]*entity.Lead, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return uc.leadRepo.FindAll(ctx, brandID, status, page, limit)
}

func (uc *LeadUseCase) GetKanban(ctx context.Context, brandID *uuid.UUID) (map[string][]*entity.Lead, error) {
	return uc.leadRepo.FindByBrandGrouped(ctx, brandID)
}

func (uc *LeadUseCase) UpdateStatus(ctx context.Context, id uuid.UUID, req UpdateLeadStatusRequest, performedBy uuid.UUID) (*entity.Lead, error) {
	if !entity.IsValidLeadStatus(req.Status) {
		return nil, fmt.Errorf("invalid status: %s", req.Status)
	}
	lead, err := uc.leadRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	oldStatus := lead.Status
	progress := uc.calculateProgress(req.Status)
	if err := uc.leadRepo.UpdateStatus(ctx, id, req.Status, progress); err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}
	// Log activity
	uc.logRepo.Create(ctx, &entity.ActivityLog{
		ID: uuid.New(), EntityType: "lead", EntityID: id,
		Action: "STATUS_CHANGE", Description: fmt.Sprintf("Status changed from %s to %s", oldStatus, req.Status),
		OldValue: oldStatus, NewValue: req.Status,
		PerformedBy: performedBy, CreatedAt: time.Now(),
	})
	lead.Status = req.Status
	lead.ProgressPercentage = progress
	return lead, nil
}

func (uc *LeadUseCase) Update(ctx context.Context, id uuid.UUID, req CreateLeadRequest) (*entity.Lead, error) {
	lead, err := uc.leadRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	brandID, _ := uuid.Parse(req.BrandID)
	lead.BrandID = brandID
	lead.FullName = req.FullName
	lead.Email = req.Email
	lead.Phone = req.Phone
	lead.Notes = req.Notes
	if err := uc.leadRepo.Update(ctx, lead); err != nil {
		return nil, err
	}
	return lead, nil
}

func (uc *LeadUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.leadRepo.Delete(ctx, id)
}

func (uc *LeadUseCase) calculateProgress(status string) int {
	progressMap := map[string]int{
		entity.LeadStatusNew: 0, entity.LeadStatusConsultation: 8,
		entity.LeadStatusLocationSubmitted: 15, entity.LeadStatusSurveyApproved: 25,
		entity.LeadStatusMeetingDone: 35, entity.LeadStatusReadyForDP: 40,
		entity.LeadStatusDPPaid: 50, entity.LeadStatusAgreementReview: 60,
		entity.LeadStatusFullyPaid: 75, entity.LeadStatusActivePartnership: 85,
		entity.LeadStatusRunning: 95, entity.LeadStatusCompleted: 100,
	}
	if p, ok := progressMap[status]; ok {
		return p
	}
	return 0
}
