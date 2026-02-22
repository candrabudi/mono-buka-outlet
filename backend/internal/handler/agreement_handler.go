package handler

import (
	"net/http"

	"github.com/franchise-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

func (h *AgreementHandler) GetByMitra(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	agreements, err := h.agreementUC.GetByMitra(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": agreements})
}
