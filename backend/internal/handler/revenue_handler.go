package handler

import (
	"net/http"

	"github.com/franchise-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	list, err := h.revenueUC.GetByPartnership(c.Request.Context(), pid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
}
