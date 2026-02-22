package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/google/uuid"
)

type PaymentUseCase struct {
	paymentRepo     repository.PaymentRepository
	partnershipRepo repository.PartnershipRepository
	logRepo         repository.ActivityLogRepository
}

func NewPaymentUseCase(pr repository.PaymentRepository, par repository.PartnershipRepository, lr repository.ActivityLogRepository) *PaymentUseCase {
	return &PaymentUseCase{paymentRepo: pr, partnershipRepo: par, logRepo: lr}
}

type CreatePaymentRequest struct {
	PartnershipID string  `json:"partnership_id" binding:"required"`
	Type          string  `json:"type" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	ProofURL      string  `json:"proof_url"`
	Notes         string  `json:"notes"`
}

func (uc *PaymentUseCase) Create(ctx context.Context, req CreatePaymentRequest) (*entity.Payment, error) {
	pid, err := uuid.Parse(req.PartnershipID)
	if err != nil {
		return nil, fmt.Errorf("invalid partnership_id")
	}
	partnership, err := uc.partnershipRepo.FindByID(ctx, pid)
	if err != nil {
		return nil, err
	}

	var brandID uuid.UUID
	if partnership.BrandID != nil {
		brandID = *partnership.BrandID
	}
	payment := &entity.Payment{
		ID: uuid.New(), PartnershipID: pid, BrandID: brandID,
		Type: req.Type, Amount: req.Amount, ProofURL: req.ProofURL,
		VerifiedStatus: entity.PaymentStatusPending, Notes: req.Notes,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := uc.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}
	return payment, nil
}

func (uc *PaymentUseCase) Verify(ctx context.Context, id uuid.UUID, status string, verifiedBy uuid.UUID) error {
	payment, err := uc.paymentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := uc.paymentRepo.Verify(ctx, id, status, verifiedBy); err != nil {
		return fmt.Errorf("failed to verify payment: %w", err)
	}
	if status == entity.PaymentStatusVerified && payment.Type == entity.PaymentTypeDP {
		uc.partnershipRepo.UpdateProgress(ctx, payment.PartnershipID, 25, entity.PartnershipStatusDPVerified)
	}
	uc.logRepo.Create(ctx, &entity.ActivityLog{
		ID: uuid.New(), EntityType: "payment", EntityID: id,
		Action: "PAYMENT_VERIFIED", Description: fmt.Sprintf("Payment %s verified as %s", payment.Type, status),
		NewValue: status, PerformedBy: verifiedBy, CreatedAt: time.Now(),
	})
	return nil
}

func (uc *PaymentUseCase) GetByPartnership(ctx context.Context, pid uuid.UUID) ([]*entity.Payment, error) {
	return uc.paymentRepo.FindByPartnershipID(ctx, pid)
}

func (uc *PaymentUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	return uc.paymentRepo.FindByID(ctx, id)
}
