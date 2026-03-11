package usecase

import (
	"context"
	"fmt"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/google/uuid"
)

type AffiliatorUseCase struct {
	userRepo        repository.UserRepository
	partnershipRepo repository.PartnershipRepository
}

func NewAffiliatorUseCase(ur repository.UserRepository, pr repository.PartnershipRepository) *AffiliatorUseCase {
	return &AffiliatorUseCase{userRepo: ur, partnershipRepo: pr}
}

type AffiliatorDashboard struct {
	TotalReferrals      int     `json:"total_referrals"`
	ActivePartnerships  int     `json:"active_partnerships"`
	PendingPartnerships int     `json:"pending_partnerships"`
	ReferralCode        string  `json:"referral_code"`
	ConversionRate      float64 `json:"conversion_rate"`
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

	return &AffiliatorDashboard{
		TotalReferrals:      total,
		ActivePartnerships:  activeCount,
		PendingPartnerships: pendingCount,
		ReferralCode:        user.ReferralCode,
		ConversionRate:      convRate,
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
