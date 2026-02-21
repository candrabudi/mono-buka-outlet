package handler

import (
	"net/http"

	"github.com/franchise-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OutletCategoryHandler struct {
	catUC *usecase.OutletCategoryUseCase
}

func NewOutletCategoryHandler(uc *usecase.OutletCategoryUseCase) *OutletCategoryHandler {
	return &OutletCategoryHandler{catUC: uc}
}

func (h *OutletCategoryHandler) Create(c *gin.Context) {
	var req usecase.CreateOutletCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid", "details": err.Error()})
		return
	}

	cat, validationErrs := h.catUC.Create(c.Request.Context(), req)
	if validationErrs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Validasi gagal", "errors": validationErrs})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": cat, "message": "Kategori berhasil dibuat"})
}

func (h *OutletCategoryHandler) GetAll(c *gin.Context) {
	activeOnly := c.Query("active_only") == "true"
	cats, err := h.catUC.GetAll(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": cats})
}

func (h *OutletCategoryHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	cat, err := h.catUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": cat})
}

func (h *OutletCategoryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	var req usecase.CreateOutletCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}
	cat, validationErrs := h.catUC.Update(c.Request.Context(), id, req)
	if validationErrs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Validasi gagal", "errors": validationErrs})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": cat, "message": "Kategori berhasil diupdate"})
}

func (h *OutletCategoryHandler) ToggleActive(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	cat, err := h.catUC.ToggleActive(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	msg := "Kategori dinonaktifkan"
	if cat.IsActive {
		msg = "Kategori diaktifkan"
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": cat, "message": msg})
}

func (h *OutletCategoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	if err := h.catUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Kategori berhasil dihapus"})
}

// Public endpoint — list active categories
func (h *OutletCategoryHandler) PublicList(c *gin.Context) {
	cats, err := h.catUC.GetAll(c.Request.Context(), true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": cats})
}
