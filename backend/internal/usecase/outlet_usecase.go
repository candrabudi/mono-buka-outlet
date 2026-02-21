package usecase

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/google/uuid"
)

type OutletUseCase struct {
	outletRepo   repository.OutletRepository
	categoryRepo repository.OutletCategoryRepository
}

func NewOutletUseCase(or repository.OutletRepository, cr repository.OutletCategoryRepository) *OutletUseCase {
	return &OutletUseCase{outletRepo: or, categoryRepo: cr}
}

// Request structs with server-side validation
type CreateOutletRequest struct {
	Name                    string     `json:"name"`
	Logo                    string     `json:"logo"`
	Banner                  string     `json:"banner"`
	CategoryID              *uuid.UUID `json:"category_id"`
	Description             string     `json:"description"`
	ShortDescription        string     `json:"short_description"`
	MinimumInvestment       float64    `json:"minimum_investment"`
	MaximumInvestment       *float64   `json:"maximum_investment"`
	ProfitSharingPercentage float64    `json:"profit_sharing_percentage"`
	EstimatedROI            string     `json:"estimated_roi"`
	LocationRequirement     string     `json:"location_requirement"`
	Address                 string     `json:"address"`
	City                    string     `json:"city"`
	Province                string     `json:"province"`
	Latitude                *float64   `json:"latitude"`
	Longitude               *float64   `json:"longitude"`
	ContactPhone            string     `json:"contact_phone"`
	ContactEmail            string     `json:"contact_email"`
	ContactWhatsapp         string     `json:"contact_whatsapp"`
	Website                 string     `json:"website"`
	TotalOutlets            int        `json:"total_outlets"`
	YearEstablished         *int       `json:"year_established"`
}

// ValidationError holds per-field validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (req *CreateOutletRequest) Validate() []ValidationError {
	var errs []ValidationError

	if strings.TrimSpace(req.Name) == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "Nama outlet wajib diisi"})
	} else if len(req.Name) < 3 {
		errs = append(errs, ValidationError{Field: "name", Message: "Nama outlet minimal 3 karakter"})
	} else if len(req.Name) > 255 {
		errs = append(errs, ValidationError{Field: "name", Message: "Nama outlet maksimal 255 karakter"})
	}

	if req.CategoryID == nil {
		errs = append(errs, ValidationError{Field: "category_id", Message: "Kategori wajib dipilih"})
	}

	if req.MinimumInvestment <= 0 {
		errs = append(errs, ValidationError{Field: "minimum_investment", Message: "Investasi minimum harus lebih dari 0"})
	}

	if req.MaximumInvestment != nil && *req.MaximumInvestment < req.MinimumInvestment {
		errs = append(errs, ValidationError{Field: "maximum_investment", Message: "Investasi maksimum harus lebih dari minimum"})
	}

	if req.ProfitSharingPercentage < 0 || req.ProfitSharingPercentage > 100 {
		errs = append(errs, ValidationError{Field: "profit_sharing_percentage", Message: "Profit sharing harus antara 0 - 100%"})
	}

	if strings.TrimSpace(req.Description) == "" {
		errs = append(errs, ValidationError{Field: "description", Message: "Deskripsi outlet wajib diisi"})
	}

	if req.ContactEmail != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(req.ContactEmail) {
			errs = append(errs, ValidationError{Field: "contact_email", Message: "Format email tidak valid"})
		}
	}

	if req.Website != "" && !strings.HasPrefix(req.Website, "http") {
		errs = append(errs, ValidationError{Field: "website", Message: "Website harus diawali dengan http:// atau https://"})
	}

	if strings.TrimSpace(req.City) == "" {
		errs = append(errs, ValidationError{Field: "city", Message: "Kota wajib diisi"})
	}

	if strings.TrimSpace(req.Province) == "" {
		errs = append(errs, ValidationError{Field: "province", Message: "Provinsi wajib diisi"})
	}

	return errs
}

// slugify creates a URL-friendly slug from a string
func slugify(s string) string {
	s = strings.ToLower(s)
	reg := regexp.MustCompile(`[^a-z0-9\s-]`)
	s = reg.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	reg = regexp.MustCompile(`[\s-]+`)
	s = reg.ReplaceAllString(s, "-")
	return s
}

func (uc *OutletUseCase) Create(ctx context.Context, req CreateOutletRequest, createdBy *uuid.UUID) (*entity.Outlet, []ValidationError) {
	if errs := req.Validate(); len(errs) > 0 {
		return nil, errs
	}

	// Validate category exists
	if req.CategoryID != nil {
		cat, err := uc.categoryRepo.FindByID(ctx, *req.CategoryID)
		if err != nil || cat == nil {
			return nil, []ValidationError{{Field: "category_id", Message: "Kategori tidak ditemukan"}}
		}
	}

	slug := slugify(req.Name)
	if existing, _ := uc.outletRepo.FindBySlug(ctx, slug); existing != nil {
		slug = slug + "-" + uuid.New().String()[:8]
	}

	outlet := &entity.Outlet{
		ID:                      uuid.New(),
		Name:                    strings.TrimSpace(req.Name),
		Slug:                    slug,
		Logo:                    req.Logo,
		Banner:                  req.Banner,
		CategoryID:              req.CategoryID,
		Description:             req.Description,
		ShortDescription:        req.ShortDescription,
		MinimumInvestment:       req.MinimumInvestment,
		MaximumInvestment:       req.MaximumInvestment,
		ProfitSharingPercentage: req.ProfitSharingPercentage,
		EstimatedROI:            req.EstimatedROI,
		LocationRequirement:     req.LocationRequirement,
		Address:                 req.Address,
		City:                    strings.TrimSpace(req.City),
		Province:                strings.TrimSpace(req.Province),
		Latitude:                req.Latitude,
		Longitude:               req.Longitude,
		ContactPhone:            req.ContactPhone,
		ContactEmail:            req.ContactEmail,
		ContactWhatsapp:         req.ContactWhatsapp,
		Website:                 req.Website,
		IsActive:                true,
		IsFeatured:              false,
		TotalOutlets:            req.TotalOutlets,
		YearEstablished:         req.YearEstablished,
		CreatedBy:               createdBy,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}

	if err := uc.outletRepo.Create(ctx, outlet); err != nil {
		return nil, []ValidationError{{Field: "_general", Message: fmt.Sprintf("Gagal menyimpan outlet: %v", err)}}
	}
	return outlet, nil
}

func (uc *OutletUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Outlet, error) {
	return uc.outletRepo.FindByID(ctx, id)
}

func (uc *OutletUseCase) GetBySlug(ctx context.Context, slug string) (*entity.Outlet, error) {
	return uc.outletRepo.FindBySlug(ctx, slug)
}

func (uc *OutletUseCase) GetAll(ctx context.Context, filter repository.OutletFilter) ([]*entity.Outlet, int, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	return uc.outletRepo.FindAll(ctx, filter)
}

func (uc *OutletUseCase) Update(ctx context.Context, id uuid.UUID, req CreateOutletRequest) (*entity.Outlet, []ValidationError) {
	if errs := req.Validate(); len(errs) > 0 {
		return nil, errs
	}

	// Validate category exists
	if req.CategoryID != nil {
		cat, err := uc.categoryRepo.FindByID(ctx, *req.CategoryID)
		if err != nil || cat == nil {
			return nil, []ValidationError{{Field: "category_id", Message: "Kategori tidak ditemukan"}}
		}
	}

	outlet, err := uc.outletRepo.FindByID(ctx, id)
	if err != nil {
		return nil, []ValidationError{{Field: "_general", Message: "Outlet tidak ditemukan"}}
	}

	if outlet.Name != req.Name {
		slug := slugify(req.Name)
		if existing, _ := uc.outletRepo.FindBySlug(ctx, slug); existing != nil && existing.ID != id {
			slug = slug + "-" + uuid.New().String()[:8]
		}
		outlet.Slug = slug
	}

	outlet.Name = strings.TrimSpace(req.Name)
	if req.Logo != "" {
		outlet.Logo = req.Logo
	}
	if req.Banner != "" {
		outlet.Banner = req.Banner
	}
	outlet.CategoryID = req.CategoryID
	outlet.Description = req.Description
	outlet.ShortDescription = req.ShortDescription
	outlet.MinimumInvestment = req.MinimumInvestment
	outlet.MaximumInvestment = req.MaximumInvestment
	outlet.ProfitSharingPercentage = req.ProfitSharingPercentage
	outlet.EstimatedROI = req.EstimatedROI
	outlet.LocationRequirement = req.LocationRequirement
	outlet.Address = req.Address
	outlet.City = strings.TrimSpace(req.City)
	outlet.Province = strings.TrimSpace(req.Province)
	outlet.Latitude = req.Latitude
	outlet.Longitude = req.Longitude
	outlet.ContactPhone = req.ContactPhone
	outlet.ContactEmail = req.ContactEmail
	outlet.ContactWhatsapp = req.ContactWhatsapp
	outlet.Website = req.Website
	outlet.TotalOutlets = req.TotalOutlets
	outlet.YearEstablished = req.YearEstablished
	outlet.UpdatedAt = time.Now()

	if err := uc.outletRepo.Update(ctx, outlet); err != nil {
		return nil, []ValidationError{{Field: "_general", Message: fmt.Sprintf("Gagal update outlet: %v", err)}}
	}
	return outlet, nil
}

func (uc *OutletUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.outletRepo.Delete(ctx, id)
}

func (uc *OutletUseCase) ToggleActive(ctx context.Context, id uuid.UUID) (*entity.Outlet, error) {
	outlet, err := uc.outletRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	outlet.IsActive = !outlet.IsActive
	if err := uc.outletRepo.Update(ctx, outlet); err != nil {
		return nil, err
	}
	return outlet, nil
}

func (uc *OutletUseCase) ToggleFeatured(ctx context.Context, id uuid.UUID) (*entity.Outlet, error) {
	outlet, err := uc.outletRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	outlet.IsFeatured = !outlet.IsFeatured
	if err := uc.outletRepo.Update(ctx, outlet); err != nil {
		return nil, err
	}
	return outlet, nil
}
