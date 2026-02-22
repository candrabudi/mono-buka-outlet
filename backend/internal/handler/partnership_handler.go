package handler

import (
	"net/http"
	"strconv"

	"github.com/franchise-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ══════════════════════════════════
// Partnership Handler
// ══════════════════════════════════

type PartnershipHandler struct {
	partnershipUC *usecase.PartnershipUseCase
}

func NewPartnershipHandler(uc *usecase.PartnershipUseCase) *PartnershipHandler {
	return &PartnershipHandler{partnershipUC: uc}
}

func (h *PartnershipHandler) Create(c *gin.Context) {
	var req usecase.CreatePartnershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := h.partnershipUC.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": p})
}

func (h *PartnershipHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	var brandID, mitraID *uuid.UUID
	if bid := c.Query("brand_id"); bid != "" {
		id, _ := uuid.Parse(bid)
		brandID = &id
	}
	if mid := c.Query("mitra_id"); mid != "" {
		id, _ := uuid.Parse(mid)
		mitraID = &id
	}
	list, total, err := h.partnershipUC.GetAll(c.Request.Context(), brandID, mitraID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": list, "meta": gin.H{"total": total, "page": page, "limit": limit}})
}

func (h *PartnershipHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	p, err := h.partnershipUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": p})
}

func (h *PartnershipHandler) GetMyPartnerships(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	list, err := h.partnershipUC.GetByMitra(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
}

// GetByMitra — alias for mitra portal
func (h *PartnershipHandler) GetByMitra(c *gin.Context) {
	h.GetMyPartnerships(c)
}
