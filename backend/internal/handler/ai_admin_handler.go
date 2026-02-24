package handler

import (
	"net/http"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository/postgres"
	"github.com/franchise-system/backend/internal/service/chatbot"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// AIAdminHandler handles AI settings management from admin panel
type AIAdminHandler struct {
	repo    *postgres.AIKnowledgeRepo
	service *chatbot.Service
}

func NewAIAdminHandler(repo *postgres.AIKnowledgeRepo, service *chatbot.Service) *AIAdminHandler {
	return &AIAdminHandler{repo: repo, service: service}
}

// ══════════════════════════════════════════════════════════════
// KNOWLEDGE BASE
// ══════════════════════════════════════════════════════════════

func (h *AIAdminHandler) ListKnowledge(c *gin.Context) {
	items, err := h.repo.GetAllActiveKnowledge()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat knowledge base"})
		return
	}
	if items == nil {
		items = []entity.AIKnowledgeBase{}
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *AIAdminHandler) CreateKnowledge(c *gin.Context) {
	var req struct {
		CategoryID string   `json:"category_id"`
		Title      string   `json:"title" binding:"required"`
		Slug       string   `json:"slug" binding:"required"`
		Content    string   `json:"content" binding:"required"`
		Keywords   []string `json:"keywords"`
		Priority   int      `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	kb := &entity.AIKnowledgeBase{
		Title:    req.Title,
		Slug:     req.Slug,
		Content:  req.Content,
		Keywords: pq.StringArray(req.Keywords),
		Priority: req.Priority,
		IsActive: true,
	}

	if req.CategoryID != "" {
		catID, err := uuid.Parse(req.CategoryID)
		if err == nil {
			kb.CategoryID = &catID
		}
	}

	if err := h.repo.CreateKnowledge(kb); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan knowledge"})
		return
	}

	h.service.InvalidateCache()
	c.JSON(http.StatusCreated, gin.H{"message": "Knowledge berhasil ditambahkan", "item": kb})
}

func (h *AIAdminHandler) UpdateKnowledge(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req struct {
		CategoryID string   `json:"category_id"`
		Title      string   `json:"title" binding:"required"`
		Slug       string   `json:"slug" binding:"required"`
		Content    string   `json:"content" binding:"required"`
		Keywords   []string `json:"keywords"`
		Priority   int      `json:"priority"`
		IsActive   bool     `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	kb := &entity.AIKnowledgeBase{
		ID:       id,
		Title:    req.Title,
		Slug:     req.Slug,
		Content:  req.Content,
		Keywords: pq.StringArray(req.Keywords),
		Priority: req.Priority,
		IsActive: req.IsActive,
	}

	if req.CategoryID != "" {
		catID, err := uuid.Parse(req.CategoryID)
		if err == nil {
			kb.CategoryID = &catID
		}
	}

	if err := h.repo.UpdateKnowledge(kb); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate knowledge"})
		return
	}

	h.service.InvalidateCache()
	c.JSON(http.StatusOK, gin.H{"message": "Knowledge berhasil diupdate"})
}

func (h *AIAdminHandler) DeleteKnowledge(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.repo.DeleteKnowledge(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus knowledge"})
		return
	}

	h.service.InvalidateCache()
	c.JSON(http.StatusOK, gin.H{"message": "Knowledge berhasil dihapus"})
}

// ══════════════════════════════════════════════════════════════
// CATEGORIES
// ══════════════════════════════════════════════════════════════

func (h *AIAdminHandler) ListCategories(c *gin.Context) {
	cats, err := h.repo.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat kategori"})
		return
	}
	if cats == nil {
		cats = []entity.AIKBCategory{}
	}
	c.JSON(http.StatusOK, gin.H{"categories": cats})
}

func (h *AIAdminHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Slug        string `json:"slug" binding:"required"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	cat := &entity.AIKBCategory{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		IsActive:    true,
	}
	if err := h.repo.CreateCategory(cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan kategori"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Kategori berhasil ditambahkan", "category": cat})
}

// ══════════════════════════════════════════════════════════════
// SYSTEM PROMPTS
// ══════════════════════════════════════════════════════════════

func (h *AIAdminHandler) ListPrompts(c *gin.Context) {
	prompts, err := h.repo.GetAllSystemPrompts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat system prompts"})
		return
	}
	if prompts == nil {
		prompts = []entity.AISystemPrompt{}
	}
	c.JSON(http.StatusOK, gin.H{"prompts": prompts})
}

func (h *AIAdminHandler) CreatePrompt(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Prompt   string `json:"prompt" binding:"required"`
		IsActive bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	sp := &entity.AISystemPrompt{
		Name:     req.Name,
		Prompt:   req.Prompt,
		IsActive: req.IsActive,
	}
	if err := h.repo.CreateSystemPrompt(sp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan prompt"})
		return
	}

	h.service.InvalidateCache()
	c.JSON(http.StatusCreated, gin.H{"message": "System prompt berhasil ditambahkan"})
}

func (h *AIAdminHandler) UpdatePrompt(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req struct {
		Name     string `json:"name" binding:"required"`
		Prompt   string `json:"prompt" binding:"required"`
		IsActive bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	sp := &entity.AISystemPrompt{
		ID:       id,
		Name:     req.Name,
		Prompt:   req.Prompt,
		IsActive: req.IsActive,
	}
	if err := h.repo.UpdateSystemPrompt(sp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate prompt"})
		return
	}

	h.service.InvalidateCache()
	c.JSON(http.StatusOK, gin.H{"message": "System prompt berhasil diupdate"})
}

// ══════════════════════════════════════════════════════════════
// AI CONFIG
// ══════════════════════════════════════════════════════════════

func (h *AIAdminHandler) GetConfig(c *gin.Context) {
	configs, err := h.repo.GetAllConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat konfigurasi"})
		return
	}
	if configs == nil {
		configs = []entity.AIConfig{}
	}
	c.JSON(http.StatusOK, gin.H{"configs": configs})
}

func (h *AIAdminHandler) UpdateConfig(c *gin.Context) {
	var req struct {
		Configs []struct {
			Key         string `json:"key"`
			Value       string `json:"value"`
			Description string `json:"description"`
		} `json:"configs" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	for _, cfg := range req.Configs {
		if err := h.repo.SetConfig(cfg.Key, cfg.Value, cfg.Description); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan konfigurasi " + cfg.Key})
			return
		}
	}

	h.service.InvalidateCache()
	c.JSON(http.StatusOK, gin.H{"message": "Konfigurasi AI berhasil diperbarui"})
}

// ══════════════════════════════════════════════════════════════
// CACHE
// ══════════════════════════════════════════════════════════════

func (h *AIAdminHandler) InvalidateCache(c *gin.Context) {
	h.service.InvalidateCache()
	c.JSON(http.StatusOK, gin.H{"message": "Cache AI berhasil di-refresh"})
}
