package handler

import (
	"net/http"

	"github.com/franchise-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminAuthHandler struct {
	adminAuthUC *usecase.AdminAuthUseCase
}

func NewAdminAuthHandler(uc *usecase.AdminAuthUseCase) *AdminAuthHandler {
	return &AdminAuthHandler{adminAuthUC: uc}
}

// POST /api/v1/admin/auth/login
// Step 1: validate email + password → send OTP to email
func (h *AdminAuthHandler) Login(c *gin.Context) {
	var req usecase.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	resp, err := h.adminAuthUC.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
		"message": resp.Message,
	})
}

// POST /api/v1/admin/auth/verify-otp
// Step 2: validate OTP → return JWT token
func (h *AdminAuthHandler) VerifyOTP(c *gin.Context) {
	var req usecase.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	resp, err := h.adminAuthUC.VerifyOTP(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// POST /api/v1/admin/auth/resend-otp
func (h *AdminAuthHandler) ResendOTP(c *gin.Context) {
	var req usecase.ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	resp, err := h.adminAuthUC.ResendOTP(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp, "message": resp.Message})
}

// GET /api/v1/admin/profile
func (h *AdminAuthHandler) Profile(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	user, err := h.adminAuthUC.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": user})
}
