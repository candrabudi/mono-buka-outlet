package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MitraHandler struct {
	appRepo    repository.PartnershipApplicationRepository
	outletRepo repository.OutletRepository
	pkgRepo    repository.OutletPackageRepository
	partRepo   repository.PartnershipRepository
	userRepo   repository.UserRepository
}

func NewMitraHandler(
	appRepo repository.PartnershipApplicationRepository,
	outletRepo repository.OutletRepository,
	pkgRepo repository.OutletPackageRepository,
	partRepo repository.PartnershipRepository,
	userRepo repository.UserRepository,
) *MitraHandler {
	return &MitraHandler{
		appRepo:    appRepo,
		outletRepo: outletRepo,
		pkgRepo:    pkgRepo,
		partRepo:   partRepo,
		userRepo:   userRepo,
	}
}

// ══════════════════════════════════
// OUTLET BROWSING (for mitra)
// ══════════════════════════════════

// ListOutlets — public list of active outlets
func (h *MitraHandler) ListOutlets(c *gin.Context) {
	search := c.Query("search")
	city := c.Query("city")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	active := true
	filter := repository.OutletFilter{
		Search: search,
		City:   city,
		Active: &active,
		Page:   page,
		Limit:  limit,
	}

	outlets, total, err := h.outletRepo.FindAll(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": outlets, "total": total, "page": page, "limit": limit})
}

// GetOutlet — single outlet detail
func (h *MitraHandler) GetOutlet(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	outlet, err := h.outletRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Outlet tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": outlet})
}

// GetOutletPackages — packages for a specific outlet
func (h *MitraHandler) GetOutletPackages(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	pkgs, err := h.pkgRepo.FindByOutletID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": pkgs})
}

// ══════════════════════════════════
// PARTNERSHIP APPLICATION
// ══════════════════════════════════

type ApplyRequest struct {
	OutletID         uuid.UUID `json:"outlet_id" binding:"required"`
	PackageID        uuid.UUID `json:"package_id" binding:"required"`
	Motivation       string    `json:"motivation"`
	Experience       string    `json:"experience"`
	ProposedLocation string    `json:"proposed_location"`
	InvestmentBudget float64   `json:"investment_budget"`
	ContactPhone     string    `json:"contact_phone" binding:"required"`
	ContactEmail     string    `json:"contact_email" binding:"required,email"`
}

// Apply — mitra submits a partnership application
func (h *MitraHandler) Apply(c *gin.Context) {
	mitraID := c.MustGet("user_id").(uuid.UUID)

	var req ApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Check outlet exists
	if _, err := h.outletRepo.FindByID(c.Request.Context(), req.OutletID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Outlet tidak ditemukan"})
		return
	}

	// Check package exists
	if _, err := h.pkgRepo.FindByID(c.Request.Context(), req.PackageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Paket tidak ditemukan"})
		return
	}

	// Check for duplicate active application
	hasActive, err := h.appRepo.HasActiveApplication(c.Request.Context(), mitraID, req.OutletID, req.PackageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Gagal mengecek pengajuan"})
		return
	}
	if hasActive {
		c.JSON(http.StatusConflict, gin.H{"success": false, "error": "Anda sudah memiliki pengajuan aktif untuk outlet dan paket ini"})
		return
	}

	now := time.Now()
	app := &entity.PartnershipApplication{
		ID:               uuid.New(),
		MitraID:          mitraID,
		OutletID:         req.OutletID,
		PackageID:        req.PackageID,
		Motivation:       req.Motivation,
		Experience:       req.Experience,
		ProposedLocation: req.ProposedLocation,
		InvestmentBudget: req.InvestmentBudget,
		ContactPhone:     req.ContactPhone,
		ContactEmail:     req.ContactEmail,
		Status:           entity.ApplicationStatusPending,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := h.appRepo.Create(c.Request.Context(), app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Gagal mengirim pengajuan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": app, "message": "Pengajuan berhasil dikirim"})
}

// CancelApplication — mitra cancels their own PENDING application
func (h *MitraHandler) CancelApplication(c *gin.Context) {
	mitraID := c.MustGet("user_id").(uuid.UUID)
	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	if err := h.appRepo.CancelByMitra(c.Request.Context(), appID, mitraID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Pengajuan berhasil dibatalkan"})
}

// MyApplications — list mitra's own applications
func (h *MitraHandler) MyApplications(c *gin.Context) {
	mitraID := c.MustGet("user_id").(uuid.UUID)

	apps, err := h.appRepo.FindByMitraID(c.Request.Context(), mitraID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": apps})
}

// GetApplication — single application detail
func (h *MitraHandler) GetApplication(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	app, err := h.appRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Pengajuan tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": app})
}

// ══════════════════════════════════
// ADMIN: Review applications
// ══════════════════════════════════

type ReviewRequest struct {
	Status     string `json:"status" binding:"required"`
	AdminNotes string `json:"admin_notes"`
}

// ListAllApplications — admin views all applications
func (h *MitraHandler) ListAllApplications(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	apps, total, err := h.appRepo.FindAll(c.Request.Context(), status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": apps, "total": total, "page": page, "limit": limit})
}

// ReviewApplication — admin approves/rejects an application
func (h *MitraHandler) ReviewApplication(c *gin.Context) {
	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	adminID := c.MustGet("user_id").(uuid.UUID)

	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if req.Status != entity.ApplicationStatusApproved && req.Status != entity.ApplicationStatusRejected && req.Status != entity.ApplicationStatusReviewed {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Status harus APPROVED, REJECTED, atau REVIEWED"})
		return
	}

	// Get app
	app, err := h.appRepo.FindByID(c.Request.Context(), appID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Pengajuan tidak ditemukan"})
		return
	}

	// Block changes on final statuses
	if app.Status == entity.ApplicationStatusRejected || app.Status == entity.ApplicationStatusCancelled || app.Status == entity.ApplicationStatusApproved {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Pengajuan dengan status " + app.Status + " tidak dapat diubah lagi"})
		return
	}

	// Update status
	if err := h.appRepo.UpdateStatus(c.Request.Context(), appID, req.Status, req.AdminNotes, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Gagal memperbarui pengajuan"})
		return
	}

	// If APPROVED → auto-create partnership
	if req.Status == entity.ApplicationStatusApproved {
		now := time.Now()
		partnership := &entity.Partnership{
			ID:                 uuid.New(),
			MitraID:            app.MitraID,
			OutletID:           &app.OutletID,
			PackageID:          &app.PackageID,
			ProgressPercentage: 0,
			Status:             entity.PartnershipStatusPending,
			CreatedAt:          now,
			UpdatedAt:          now,
		}

		// Link to affiliator if mitra was referred
		mitraUser, _ := h.userRepo.FindByID(c.Request.Context(), app.MitraID)
		if mitraUser != nil && mitraUser.ReferredBy != nil {
			partnership.AffiliatorID = mitraUser.ReferredBy
		}

		if err := h.partRepo.Create(c.Request.Context(), partnership); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Pengajuan disetujui tapi gagal membuat partnership: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Pengajuan disetujui, partnership berhasil dibuat", "partnership_id": partnership.ID})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Status pengajuan berhasil diperbarui"})
}
