package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/google/uuid"
)

type AffiliatorUseCase struct {
	userRepo        repository.UserRepository
	partnershipRepo repository.PartnershipRepository
	commissionRepo  repository.AffiliatorCommissionRepository
	withdrawalRepo  repository.AffiliatorWithdrawalRepository
}

func NewAffiliatorUseCase(
	ur repository.UserRepository,
	pr repository.PartnershipRepository,
	cr repository.AffiliatorCommissionRepository,
	wr repository.AffiliatorWithdrawalRepository,
) *AffiliatorUseCase {
	return &AffiliatorUseCase{
		userRepo:        ur,
		partnershipRepo: pr,
		commissionRepo:  cr,
		withdrawalRepo:  wr,
	}
}

type AffiliatorDashboard struct {
	TotalReferrals      int     `json:"total_referrals"`
	ActivePartnerships  int     `json:"active_partnerships"`
	PendingPartnerships int     `json:"pending_partnerships"`
	ReferralCode        string  `json:"referral_code"`
	ConversionRate      float64 `json:"conversion_rate"`
	TotalEarned         float64 `json:"total_earned"`
	Balance             float64 `json:"balance"`
	PendingWithdrawal   float64 `json:"pending_withdrawal"`
	TotalWithdrawn      float64 `json:"total_withdrawn"`
}

func (uc *AffiliatorUseCase) GetDashboard(ctx context.Context, userID uuid.UUID) (*AffiliatorDashboard, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user tidak ditemukan")
	}
	if user.Role != entity.RoleAffiliator {
		return nil, fmt.Errorf("akses ditolak: bukan affiliator")
	}

	allPartnerships, total, err := uc.partnershipRepo.FindByAffiliatorID(ctx, userID, 1, 10000)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data partnership: %w", err)
	}

	activeCount := 0
	pendingCount := 0
	for _, p := range allPartnerships {
		switch p.Status {
		case entity.PartnershipStatusRunning, entity.PartnershipStatusDevelopment:
			activeCount++
		case entity.PartnershipStatusPending, entity.PartnershipStatusDPVerified:
			pendingCount++
		}
	}

	convRate := float64(0)
	if total > 0 {
		convRate = float64(activeCount) / float64(total) * 100
	}

	totalEarned, _ := uc.commissionRepo.GetTotalEarned(ctx, userID)
	balance, _ := uc.commissionRepo.GetBalance(ctx, userID)
	pendingWd, _ := uc.withdrawalRepo.GetTotalPending(ctx, userID)
	totalWithdrawn, _ := uc.withdrawalRepo.GetTotalWithdrawn(ctx, userID)

	return &AffiliatorDashboard{
		TotalReferrals:      total,
		ActivePartnerships:  activeCount,
		PendingPartnerships: pendingCount,
		ReferralCode:        user.ReferralCode,
		ConversionRate:      convRate,
		TotalEarned:         totalEarned,
		Balance:             balance,
		PendingWithdrawal:   pendingWd,
		TotalWithdrawn:      totalWithdrawn,
	}, nil
}

func (uc *AffiliatorUseCase) GetMyPartnerships(ctx context.Context, userID uuid.UUID, page, limit int) ([]*entity.Partnership, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return uc.partnershipRepo.FindByAffiliatorID(ctx, userID, page, limit)
}

func (uc *AffiliatorUseCase) GetReferralCode(ctx context.Context, userID uuid.UUID) (string, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("user tidak ditemukan")
	}
	if user.Role != entity.RoleAffiliator {
		return "", fmt.Errorf("akses ditolak: bukan affiliator")
	}
	return user.ReferralCode, nil
}

// ══════════════════════════════════
// COMMISSION — Admin gives commission
// ══════════════════════════════════

type GiveCommissionRequest struct {
	AffiliatorID  uuid.UUID  `json:"affiliator_id" binding:"required"`
	PartnershipID *uuid.UUID `json:"partnership_id"`
	Amount        float64    `json:"amount" binding:"required"`
	Type          string     `json:"type" binding:"required"`
	Description   string     `json:"description"`
}

func (uc *AffiliatorUseCase) GiveCommission(ctx context.Context, req GiveCommissionRequest, givenBy uuid.UUID) (*entity.AffiliatorCommission, error) {
	if req.Amount <= 0 {
		return nil, fmt.Errorf("jumlah komisi harus lebih dari 0")
	}
	if !entity.IsValidCommissionType(req.Type) {
		return nil, fmt.Errorf("tipe komisi tidak valid (gunakan: COMMISSION, BONUS, atau ADJUSTMENT)")
	}

	affiliator, err := uc.userRepo.FindByID(ctx, req.AffiliatorID)
	if err != nil {
		return nil, fmt.Errorf("affiliator tidak ditemukan")
	}
	if affiliator.Role != entity.RoleAffiliator {
		return nil, fmt.Errorf("user bukan affiliator")
	}

	commission := &entity.AffiliatorCommission{
		ID:            uuid.New(),
		AffiliatorID:  req.AffiliatorID,
		PartnershipID: req.PartnershipID,
		Amount:        req.Amount,
		Type:          req.Type,
		Description:   req.Description,
		GivenBy:       givenBy,
		CreatedAt:     time.Now(),
	}

	if err := uc.commissionRepo.Create(ctx, commission); err != nil {
		return nil, fmt.Errorf("gagal memberikan komisi: %w", err)
	}
	return commission, nil
}

func (uc *AffiliatorUseCase) GetCommissions(ctx context.Context, affiliatorID uuid.UUID, page, limit int) ([]*entity.AffiliatorCommission, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return uc.commissionRepo.FindByAffiliatorID(ctx, affiliatorID, page, limit)
}

func (uc *AffiliatorUseCase) GetBalance(ctx context.Context, affiliatorID uuid.UUID) (float64, error) {
	return uc.commissionRepo.GetBalance(ctx, affiliatorID)
}

// ══════════════════════════════════
// WITHDRAWAL — Affiliator requests withdrawal
// ══════════════════════════════════

type WithdrawalRequest struct {
	Amount        float64 `json:"amount" binding:"required"`
	BankName      string  `json:"bank_name" binding:"required"`
	AccountNumber string  `json:"account_number" binding:"required"`
	AccountHolder string  `json:"account_holder" binding:"required"`
}

func (uc *AffiliatorUseCase) RequestWithdrawal(ctx context.Context, userID uuid.UUID, req WithdrawalRequest) (*entity.AffiliatorWithdrawal, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user tidak ditemukan")
	}
	if user.Role != entity.RoleAffiliator {
		return nil, fmt.Errorf("akses ditolak: bukan affiliator")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("jumlah penarikan harus lebih dari 0")
	}

	balance, err := uc.commissionRepo.GetBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal cek saldo: %w", err)
	}

	pendingWd, _ := uc.withdrawalRepo.GetTotalPending(ctx, userID)
	availableBalance := balance - pendingWd
	if req.Amount > availableBalance {
		return nil, fmt.Errorf("saldo tidak cukup (tersedia: Rp %.0f)", availableBalance)
	}

	now := time.Now()
	withdrawal := &entity.AffiliatorWithdrawal{
		ID:            uuid.New(),
		AffiliatorID:  userID,
		Amount:        req.Amount,
		BankName:      req.BankName,
		AccountNumber: req.AccountNumber,
		AccountHolder: req.AccountHolder,
		Status:        entity.WithdrawalStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := uc.withdrawalRepo.Create(ctx, withdrawal); err != nil {
		return nil, fmt.Errorf("gagal mengajukan penarikan: %w", err)
	}
	return withdrawal, nil
}

func (uc *AffiliatorUseCase) GetMyWithdrawals(ctx context.Context, userID uuid.UUID, page, limit int) ([]*entity.AffiliatorWithdrawal, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return uc.withdrawalRepo.FindByAffiliatorID(ctx, userID, page, limit)
}

// ══════════════════════════════════
// ADMIN: Manage withdrawals
// ══════════════════════════════════

func (uc *AffiliatorUseCase) GetAllWithdrawals(ctx context.Context, status string, page, limit int) ([]*entity.AffiliatorWithdrawal, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return uc.withdrawalRepo.FindAll(ctx, status, page, limit)
}

type ProcessWithdrawalRequest struct {
	Status     string `json:"status" binding:"required"`
	AdminNotes string `json:"admin_notes"`
}

func (uc *AffiliatorUseCase) ProcessWithdrawal(ctx context.Context, withdrawalID uuid.UUID, req ProcessWithdrawalRequest, adminID uuid.UUID) error {
	withdrawal, err := uc.withdrawalRepo.FindByID(ctx, withdrawalID)
	if err != nil {
		return fmt.Errorf("permintaan penarikan tidak ditemukan")
	}

	if withdrawal.Status != entity.WithdrawalStatusPending && withdrawal.Status != entity.WithdrawalStatusApproved {
		return fmt.Errorf("permintaan penarikan dengan status %s tidak dapat diproses", withdrawal.Status)
	}

	validStatuses := map[string]bool{
		entity.WithdrawalStatusApproved:    true,
		entity.WithdrawalStatusTransferred: true,
		entity.WithdrawalStatusRejected:    true,
	}
	if !validStatuses[req.Status] {
		return fmt.Errorf("status harus APPROVED, TRANSFERRED, atau REJECTED")
	}

	if withdrawal.Status == entity.WithdrawalStatusApproved && req.Status == entity.WithdrawalStatusApproved {
		return fmt.Errorf("permintaan sudah disetujui")
	}

	return uc.withdrawalRepo.UpdateStatus(ctx, withdrawalID, req.Status, req.AdminNotes, adminID)
}

func (uc *AffiliatorUseCase) GetWithdrawalByID(ctx context.Context, id uuid.UUID) (*entity.AffiliatorWithdrawal, error) {
	return uc.withdrawalRepo.FindByID(ctx, id)
}
