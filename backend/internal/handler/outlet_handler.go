package handler

import (
	"net/http"
	"strconv"

	"github.com/franchise-system/backend/internal/repository"
	"github.com/franchise-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OutletHandler struct {
	outletUC *usecase.OutletUseCase
}

func NewOutletHandler(uc *usecase.OutletUseCase) *OutletHandler {
	return &OutletHandler{outletUC: uc}
}

func (h *OutletHandler) Create(c *gin.Context) {
	var req usecase.CreateOutletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid", "details": err.Error()})
		return
	}

	var createdBy *uuid.UUID
	if uid, exists := c.Get("user_id"); exists {
		id := uid.(uuid.UUID)
		createdBy = &id
	}

	outlet, validationErrs := h.outletUC.Create(c.Request.Context(), req, createdBy)
	if validationErrs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":  "Validasi gagal",
			"errors": validationErrs,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": outlet, "message": "Outlet berhasil dibuat"})
}

func (h *OutletHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := repository.OutletFilter{
		Search:   c.Query("search"),
		City:     c.Query("city"),
		Province: c.Query("province"),
		Page:     page,
		Limit:    limit,
	}

	if catStr := c.Query("category_id"); catStr != "" {
		if catID, err := uuid.Parse(catStr); err == nil {
			filter.CategoryID = &catID
		}
	}
	if activeStr := c.Query("active"); activeStr != "" {
		active := activeStr == "true"
		filter.Active = &active
	}
	if featuredStr := c.Query("featured"); featuredStr != "" {
		featured := featuredStr == "true"
		filter.Featured = &featured
	}

	outlets, total, err := h.outletUC.GetAll(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    outlets,
		"meta":    gin.H{"total": total, "page": page, "limit": limit},
	})
}

func (h *OutletHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	outlet, err := h.outletUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Outlet tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": outlet})
}

func (h *OutletHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	var req usecase.CreateOutletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid", "details": err.Error()})
		return
	}
	outlet, validationErrs := h.outletUC.Update(c.Request.Context(), id, req)
	if validationErrs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":  "Validasi gagal",
			"errors": validationErrs,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": outlet, "message": "Outlet berhasil diupdate"})
}

func (h *OutletHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	if err := h.outletUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Outlet berhasil dihapus"})
}

func (h *OutletHandler) ToggleActive(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	outlet, err := h.outletUC.ToggleActive(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	msg := "Outlet dinonaktifkan"
	if outlet.IsActive {
		msg = "Outlet diaktifkan"
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": outlet, "message": msg})
}

func (h *OutletHandler) ToggleFeatured(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	outlet, err := h.outletUC.ToggleFeatured(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	msg := "Outlet dihapus dari featured"
	if outlet.IsFeatured {
		msg = "Outlet dijadikan featured"
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": outlet, "message": msg})
}

// Public endpoints (no auth)
func (h *OutletHandler) PublicList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	active := true

	filter := repository.OutletFilter{
		Search:   c.Query("search"),
		City:     c.Query("city"),
		Province: c.Query("province"),
		Active:   &active,
		Page:     page,
		Limit:    limit,
	}

	if catStr := c.Query("category_id"); catStr != "" {
		if catID, err := uuid.Parse(catStr); err == nil {
			filter.CategoryID = &catID
		}
	}

	outlets, total, err := h.outletUC.GetAll(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    outlets,
		"meta":    gin.H{"total": total, "page": page, "limit": limit},
	})
}

func (h *OutletHandler) PublicDetail(c *gin.Context) {
	// Try by slug first, then by ID
	param := c.Param("id")
	id, err := uuid.Parse(param)
	if err != nil {
		// Try slug
		outlet, err := h.outletUC.GetBySlug(c.Request.Context(), param)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Outlet tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": outlet})
		return
	}
	outlet, err := h.outletUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Outlet tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": outlet})
}
