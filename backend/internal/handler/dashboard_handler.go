package handler

import (
	"net/http"

	"github.com/franchise-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
