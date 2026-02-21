package handler

import (
	"log"
	"net/http"

	"github.com/franchise-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OutletPackageHandler struct {
	uc *usecase.OutletPackageUseCase
}

func NewOutletPackageHandler(uc *usecase.OutletPackageUseCase) *OutletPackageHandler {
	return &OutletPackageHandler{uc: uc}
}

func (h *OutletPackageHandler) Create(c *gin.Context) {
	var req usecase.CreateOutletPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	pkg, errs := h.uc.Create(c.Request.Context(), req)
	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": pkg, "message": "Paket berhasil ditambahkan"})
}

func (h *OutletPackageHandler) GetByOutletID(c *gin.Context) {
	outletID, err := uuid.Parse(c.Param("outletId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID outlet tidak valid"})
		return
	}

	packages, err := h.uc.GetByOutletID(c.Request.Context(), outletID)
	if err != nil {
		log.Printf("❌ GetByOutletID error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat paket"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": packages})
}

func (h *OutletPackageHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	pkg, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil || pkg == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Paket tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pkg})
}

func (h *OutletPackageHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req usecase.CreateOutletPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	pkg, errs := h.uc.Update(c.Request.Context(), id, req)
	if errs != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pkg, "message": "Paket berhasil diupdate"})
}

func (h *OutletPackageHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus paket"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Paket berhasil dihapus"})
}
