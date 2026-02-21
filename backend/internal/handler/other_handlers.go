package handler

import (
	"net/http"
	"strconv"

	"github.com/franchise-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Partnership Handler
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

// Payment Handler
type PaymentHandler struct {
	paymentUC *usecase.PaymentUseCase
}

func NewPaymentHandler(uc *usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{paymentUC: uc}
}

func (h *PaymentHandler) Create(c *gin.Context) {
	var req usecase.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := h.paymentUC.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": p})
}

func (h *PaymentHandler) Verify(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("user_id").(uuid.UUID)
	if err := h.paymentUC.Verify(c.Request.Context(), id, req.Status, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Payment verified"})
}

func (h *PaymentHandler) GetByPartnership(c *gin.Context) {
	pid, err := uuid.Parse(c.Param("partnership_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	list, err := h.paymentUC.GetByPartnership(c.Request.Context(), pid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
}

// Agreement Handler
type AgreementHandler struct {
	agreementUC *usecase.AgreementUseCase
}

func NewAgreementHandler(uc *usecase.AgreementUseCase) *AgreementHandler {
	return &AgreementHandler{agreementUC: uc}
}

func (h *AgreementHandler) Create(c *gin.Context) {
	var req usecase.CreateAgreementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a, err := h.agreementUC.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": a})
}

func (h *AgreementHandler) Sign(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	userID := c.MustGet("user_id").(uuid.UUID)
	if err := h.agreementUC.Sign(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Agreement signed"})
}

func (h *AgreementHandler) GetByPartnership(c *gin.Context) {
	pid, err := uuid.Parse(c.Param("partnership_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	list, err := h.agreementUC.GetByPartnership(c.Request.Context(), pid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
}

// Revenue Handler
type RevenueHandler struct {
	revenueUC *usecase.RevenueUseCase
}

func NewRevenueHandler(uc *usecase.RevenueUseCase) *RevenueHandler {
	return &RevenueHandler{revenueUC: uc}
}

func (h *RevenueHandler) Create(c *gin.Context) {
	var req usecase.CreateRevenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r, err := h.revenueUC.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": r})
}

func (h *RevenueHandler) GetByPartnership(c *gin.Context) {
	pid, err := uuid.Parse(c.Param("partnership_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	list, err := h.revenueUC.GetByPartnership(c.Request.Context(), pid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
}

// Dashboard Handler
type DashboardHandler struct {
	dashboardUC *usecase.DashboardUseCase
}

func NewDashboardHandler(uc *usecase.DashboardUseCase) *DashboardHandler {
	return &DashboardHandler{dashboardUC: uc}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	var brandID *uuid.UUID
	if bid := c.Query("brand_id"); bid != "" {
		id, _ := uuid.Parse(bid)
		brandID = &id
	}
	stats, err := h.dashboardUC.GetStats(c.Request.Context(), brandID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}
