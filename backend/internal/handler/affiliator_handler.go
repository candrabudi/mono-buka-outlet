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

// ══════════════════════════════════
// AFFILIATOR SELF-SERVICE
// ══════════════════════════════════

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

func (h *AffiliatorHandler) GetMyCommissions(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	list, total, err := h.affiliatorUC.GetCommissions(c.Request.Context(), userID, page, limit)
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

func (h *AffiliatorHandler) GetMyBalance(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	balance, err := h.affiliatorUC.GetBalance(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"balance": balance}})
}

func (h *AffiliatorHandler) RequestWithdrawal(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	var req usecase.WithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	withdrawal, err := h.affiliatorUC.RequestWithdrawal(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": withdrawal, "message": "Permintaan penarikan berhasil dikirim"})
}

func (h *AffiliatorHandler) GetMyWithdrawals(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	list, total, err := h.affiliatorUC.GetMyWithdrawals(c.Request.Context(), userID, page, limit)
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

// ══════════════════════════════════
// ADMIN: Commission management
// ══════════════════════════════════

func (h *AffiliatorHandler) GiveCommission(c *gin.Context) {
	adminID := c.MustGet("user_id").(uuid.UUID)
	var req usecase.GiveCommissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	commission, err := h.affiliatorUC.GiveCommission(c.Request.Context(), req, adminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": commission, "message": "Komisi berhasil diberikan"})
}

func (h *AffiliatorHandler) GetAffiliatorCommissions(c *gin.Context) {
	affiliatorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	list, total, err := h.affiliatorUC.GetCommissions(c.Request.Context(), affiliatorID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	balance, _ := h.affiliatorUC.GetBalance(c.Request.Context(), affiliatorID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    list,
		"balance": balance,
		"meta":    gin.H{"total": total, "page": page, "limit": limit},
	})
}

// ══════════════════════════════════
// ADMIN: Withdrawal management
// ══════════════════════════════════

func (h *AffiliatorHandler) GetAllWithdrawals(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	list, total, err := h.affiliatorUC.GetAllWithdrawals(c.Request.Context(), status, page, limit)
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

func (h *AffiliatorHandler) GetWithdrawalByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	withdrawal, err := h.affiliatorUC.GetWithdrawalByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": withdrawal})
}

func (h *AffiliatorHandler) ProcessWithdrawal(c *gin.Context) {
	adminID := c.MustGet("user_id").(uuid.UUID)
	withdrawalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}

	var req usecase.ProcessWithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := h.affiliatorUC.ProcessWithdrawal(c.Request.Context(), withdrawalID, req, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Status penarikan berhasil diperbarui"})
}
