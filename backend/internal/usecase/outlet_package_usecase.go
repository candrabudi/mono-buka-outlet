package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type OutletPackageUseCase struct {
	pkgRepo repository.OutletPackageRepository
}

func NewOutletPackageUseCase(pr repository.OutletPackageRepository) *OutletPackageUseCase {
	return &OutletPackageUseCase{pkgRepo: pr}
}

type CreateOutletPackageRequest struct {
	OutletID     uuid.UUID `json:"outlet_id"`
	Name         string    `json:"name"`
	Price        int64     `json:"price"`
	MinimumDP    int64     `json:"minimum_dp"`
	Duration     string    `json:"duration"`
	Image        string    `json:"image"`
	EstimatedBEP string    `json:"estimated_bep"`
	NetProfit    string    `json:"net_profit"`
	Description  string    `json:"description"`
	Benefits     []string  `json:"benefits"`
	SortOrder    int       `json:"sort_order"`
}

func (req *CreateOutletPackageRequest) Validate() []ValidationError {
	var errs []ValidationError
	if strings.TrimSpace(req.Name) == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "Nama paket wajib diisi"})
	}
	if req.Price <= 0 {
		errs = append(errs, ValidationError{Field: "price", Message: "Harga harus lebih dari 0"})
	}
	if req.OutletID == uuid.Nil {
		errs = append(errs, ValidationError{Field: "outlet_id", Message: "Outlet wajib dipilih"})
	}
	return errs
}

func (uc *OutletPackageUseCase) Create(ctx context.Context, req CreateOutletPackageRequest) (*entity.OutletPackage, []ValidationError) {
	if errs := req.Validate(); len(errs) > 0 {
		return nil, errs
	}

	slug := slugify(req.Name)

	pkg := &entity.OutletPackage{
		ID:           uuid.New(),
		OutletID:     req.OutletID,
		Name:         strings.TrimSpace(req.Name),
		Slug:         slug,
		Price:        req.Price,
		MinimumDP:    req.MinimumDP,
		Duration:     req.Duration,
		Image:        req.Image,
		EstimatedBEP: req.EstimatedBEP,
		NetProfit:    req.NetProfit,
		Description:  req.Description,
		Benefits:     pq.StringArray(req.Benefits),
		SortOrder:    req.SortOrder,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.pkgRepo.Create(ctx, pkg); err != nil {
		return nil, []ValidationError{{Field: "_general", Message: fmt.Sprintf("Gagal menyimpan paket: %v", err)}}
	}
	return pkg, nil
}

func (uc *OutletPackageUseCase) GetByOutletID(ctx context.Context, outletID uuid.UUID) ([]*entity.OutletPackage, error) {
	return uc.pkgRepo.FindByOutletID(ctx, outletID)
}

func (uc *OutletPackageUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.OutletPackage, error) {
	return uc.pkgRepo.FindByID(ctx, id)
}

func (uc *OutletPackageUseCase) Update(ctx context.Context, id uuid.UUID, req CreateOutletPackageRequest) (*entity.OutletPackage, []ValidationError) {
	if errs := req.Validate(); len(errs) > 0 {
		return nil, errs
	}

	pkg, err := uc.pkgRepo.FindByID(ctx, id)
	if err != nil || pkg == nil {
		return nil, []ValidationError{{Field: "_general", Message: "Paket tidak ditemukan"}}
	}

	if pkg.Name != req.Name {
		pkg.Slug = slugify(req.Name)
	}

	pkg.Name = strings.TrimSpace(req.Name)
	pkg.Price = req.Price
	pkg.MinimumDP = req.MinimumDP
	pkg.Duration = req.Duration
	pkg.Image = req.Image
	pkg.EstimatedBEP = req.EstimatedBEP
	pkg.NetProfit = req.NetProfit
	pkg.Description = req.Description
	pkg.Benefits = pq.StringArray(req.Benefits)
	pkg.SortOrder = req.SortOrder
	pkg.UpdatedAt = time.Now()

	if err := uc.pkgRepo.Update(ctx, pkg); err != nil {
		return nil, []ValidationError{{Field: "_general", Message: fmt.Sprintf("Gagal update paket: %v", err)}}
	}
	return pkg, nil
}

func (uc *OutletPackageUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.pkgRepo.Delete(ctx, id)
}
