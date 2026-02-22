package handler

import (
	"net/http"

	"github.com/franchise-system/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	repo repository.SystemSettingRepository
}

func NewSettingHandler(repo repository.SystemSettingRepository) *SettingHandler {
	return &SettingHandler{repo: repo}
}

func (h *SettingHandler) GetAll(c *gin.Context) {
	group := c.Query("group")
	var err error
	var data interface{}
	if group != "" {
		data, err = h.repo.FindByGroup(c.Request.Context(), group)
	} else {
		data, err = h.repo.FindAll(c.Request.Context())
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *SettingHandler) GetByKey(c *gin.Context) {
	key := c.Param("key")
	s, err := h.repo.FindByKey(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": s})
}

type BulkUpdateRequest struct {
	Settings map[string]string `json:"settings" binding:"required"`
}

func (h *SettingHandler) BulkUpdate(c *gin.Context) {
	var req BulkUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.BulkUpsert(c.Request.Context(), req.Settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Settings updated"})
}
