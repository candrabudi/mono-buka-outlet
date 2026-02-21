package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/google/uuid"
)

type OutletCategoryUseCase struct {
	catRepo repository.OutletCategoryRepository
}

func NewOutletCategoryUseCase(cr repository.OutletCategoryRepository) *OutletCategoryUseCase {
	return &OutletCategoryUseCase{catRepo: cr}
}

type CreateOutletCategoryRequest struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

func (req *CreateOutletCategoryRequest) Validate() []ValidationError {
	var errs []ValidationError
	if strings.TrimSpace(req.Name) == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "Nama kategori wajib diisi"})
	} else if len(req.Name) > 100 {
		errs = append(errs, ValidationError{Field: "name", Message: "Nama kategori maksimal 100 karakter"})
	}
	return errs
}

func (uc *OutletCategoryUseCase) Create(ctx context.Context, req CreateOutletCategoryRequest) (*entity.OutletCategory, []ValidationError) {
	if errs := req.Validate(); len(errs) > 0 {
		return nil, errs
	}

	slug := slugify(req.Name)
	// Check slug uniqueness
	if existing, _ := uc.catRepo.FindBySlug(ctx, slug); existing != nil {
		slug = slug + "-" + uuid.New().String()[:4]
	}

	cat := &entity.OutletCategory{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(req.Name),
		Slug:        slug,
		Icon:        req.Icon,
		Description: req.Description,
		IsActive:    true,
		SortOrder:   req.SortOrder,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.catRepo.Create(ctx, cat); err != nil {
		return nil, []ValidationError{{Field: "_general", Message: fmt.Sprintf("Gagal menyimpan kategori: %v", err)}}
	}
	return cat, nil
}

func (uc *OutletCategoryUseCase) GetAll(ctx context.Context, activeOnly bool) ([]*entity.OutletCategory, error) {
	return uc.catRepo.FindAll(ctx, activeOnly)
}

func (uc *OutletCategoryUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.OutletCategory, error) {
	return uc.catRepo.FindByID(ctx, id)
}

func (uc *OutletCategoryUseCase) Update(ctx context.Context, id uuid.UUID, req CreateOutletCategoryRequest) (*entity.OutletCategory, []ValidationError) {
	if errs := req.Validate(); len(errs) > 0 {
		return nil, errs
	}

	cat, err := uc.catRepo.FindByID(ctx, id)
	if err != nil {
		return nil, []ValidationError{{Field: "_general", Message: "Kategori tidak ditemukan"}}
	}

	if cat.Name != req.Name {
		slug := slugify(req.Name)
		if existing, _ := uc.catRepo.FindBySlug(ctx, slug); existing != nil && existing.ID != id {
			slug = slug + "-" + uuid.New().String()[:4]
		}
		cat.Slug = slug
	}

	cat.Name = strings.TrimSpace(req.Name)
	cat.Icon = req.Icon
	cat.Description = req.Description
	cat.SortOrder = req.SortOrder
	cat.UpdatedAt = time.Now()

	if err := uc.catRepo.Update(ctx, cat); err != nil {
		return nil, []ValidationError{{Field: "_general", Message: fmt.Sprintf("Gagal update kategori: %v", err)}}
	}
	return cat, nil
}

func (uc *OutletCategoryUseCase) ToggleActive(ctx context.Context, id uuid.UUID) (*entity.OutletCategory, error) {
	cat, err := uc.catRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	cat.IsActive = !cat.IsActive
	if err := uc.catRepo.Update(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (uc *OutletCategoryUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.catRepo.Delete(ctx, id)
}
