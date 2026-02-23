package handler

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EbookCategoryHandler struct {
	repo repository.EbookCategoryRepository
}

func NewEbookCategoryHandler(repo repository.EbookCategoryRepository) *EbookCategoryHandler {
	return &EbookCategoryHandler{repo: repo}
}

type EbookCategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	IsActive *bool  `json:"is_active"`
}

func (h *EbookCategoryHandler) Create(c *gin.Context) {
	var req EbookCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slug := ebookCatSlugify(req.Name)
	// Ensure unique slug
	if _, err := h.repo.FindBySlug(c.Request.Context(), slug); err == nil {
		slug = fmt.Sprintf("%s-%d", slug, time.Now().Unix())
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	cat := &entity.EbookCategory{
		Name:     req.Name,
		Slug:     slug,
		IsActive: isActive,
	}

	if err := h.repo.Create(c.Request.Context(), cat); err != nil {
		log.Printf("[EbookCategory] Create error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Gagal membuat kategori: %v", err)})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": cat})
}

func (h *EbookCategoryHandler) GetAll(c *gin.Context) {
	activeOnly := c.Query("active_only") == "true"
	cats, err := h.repo.FindAll(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cats == nil {
		cats = []*entity.EbookCategory{}
	}
	c.JSON(http.StatusOK, gin.H{"data": cats})
}

func (h *EbookCategoryHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	cat, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cat})
}

func (h *EbookCategoryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	cat, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	var req EbookCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat.Name = req.Name
	cat.Slug = ebookCatSlugify(req.Name)
	if req.IsActive != nil {
		cat.IsActive = *req.IsActive
	}

	if err := h.repo.Update(c.Request.Context(), cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cat})
}

func (h *EbookCategoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
}

func (h *EbookCategoryHandler) ToggleActive(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	cat, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}
	cat.IsActive = !cat.IsActive
	if err := h.repo.Update(c.Request.Context(), cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cat})
}

func ebookCatSlugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
