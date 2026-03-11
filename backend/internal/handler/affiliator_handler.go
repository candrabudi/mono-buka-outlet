package handler

import (
	"net/http"
	"strconv"

	"github.com/franchise-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AffiliatorHandler struct {
	affiliatorUC *usecase.AffiliatorUseCase
}

func NewAffiliatorHandler(uc *usecase.AffiliatorUseCase) *AffiliatorHandler {
	return &AffiliatorHandler{affiliatorUC: uc}
}

func (h *AffiliatorHandler) GetDashboard(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	stats, err := h.affiliatorUC.GetDashboard(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

func (h *AffiliatorHandler) GetMyPartnerships(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	list, total, err := h.affiliatorUC.GetMyPartnerships(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    list,
		"meta":    gin.H{"total": total, "page": page, "limit": limit},
	})
}

func (h *AffiliatorHandler) GetReferralCode(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	code, err := h.affiliatorUC.GetReferralCode(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"referral_code": code}})
}
