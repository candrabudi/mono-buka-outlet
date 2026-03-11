package handler

import (
	"net/http"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository/postgres"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LocationSubmissionHandler struct {
	repo *postgres.LocationSubmissionRepo
}

func NewLocationSubmissionHandler(repo *postgres.LocationSubmissionRepo) *LocationSubmissionHandler {
	return &LocationSubmissionHandler{repo: repo}
}

// ─── CREATE ───
func (h *LocationSubmissionHandler) Create(c *gin.Context) {
	var req entity.LocationSubmission
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.PartnershipID != nil {
		existing, _ := h.repo.FindByPartnership(c.Request.Context(), *req.PartnershipID)
		for _, loc := range existing {
			if loc.Status != entity.LocSubStatusRejected {
				c.JSON(http.StatusConflict, gin.H{
					"error":  "Partnership ini sudah memiliki pengajuan lokasi yang masih aktif",
					"status": loc.Status,
				})
				return
			}
		}
	}

	req.Status = entity.LocSubStatusDraft
	h.calculateScore(&req)
	if err := h.repo.Create(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": req})
}

// ─── UPDATE ───
func (h *LocationSubmissionHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	existing, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(existing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existing.ID = id
	h.calculateScore(existing)
	if err := h.repo.Update(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updated, _ := h.repo.FindByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// ─── LIST ───
func (h *LocationSubmissionHandler) GetAll(c *gin.Context) {
	status := c.Query("status")
	kota := c.Query("kota")
	search := c.Query("search")
	list, err := h.repo.FindAll(c.Request.Context(), status, kota, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if list == nil {
		list = []*entity.LocationSubmission{}
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// ─── GET BY PARTNERSHIP ───
func (h *LocationSubmissionHandler) GetByPartnership(c *gin.Context) {
	partnershipID, err := uuid.Parse(c.Param("partnershipId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid partnership id"})
		return
	}
	list, err := h.repo.FindByPartnership(c.Request.Context(), partnershipID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if list == nil {
		list = []*entity.LocationSubmission{}
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// ─── GET BY ID ───
func (h *LocationSubmissionHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	ls, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ls})
}

// ─── UPDATE STATUS ───
func (h *LocationSubmissionHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.UpdateStatus(c.Request.Context(), id, body.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ls, _ := h.repo.FindByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, gin.H{"data": ls})
}

// ─── DELETE ───
func (h *LocationSubmissionHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ─── SURVEY ───
func (h *LocationSubmissionHandler) CreateSurvey(c *gin.Context) {
	locID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var s entity.LocationSurvey
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.LocationID = locID
	if err := h.repo.CreateSurvey(c.Request.Context(), &s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Auto-update status to SURVEYED
	_ = h.repo.UpdateStatus(c.Request.Context(), locID, entity.LocSubStatusSurveyed)
	ls, _ := h.repo.FindByID(c.Request.Context(), locID)
	c.JSON(http.StatusCreated, gin.H{"data": ls})
}

// ─── FILE UPLOAD ───
func (h *LocationSubmissionHandler) AddFile(c *gin.Context) {
	locID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var f entity.LocationFile
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	f.LocationID = locID
	if err := h.repo.CreateFile(c.Request.Context(), &f); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": f})
}

func (h *LocationSubmissionHandler) DeleteFile(c *gin.Context) {
	fileID, err := uuid.Parse(c.Param("fileId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file id"})
		return
	}
	if err := h.repo.DeleteFile(c.Request.Context(), fileID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
}

// ─── APPROVAL ───
func (h *LocationSubmissionHandler) Approve(c *gin.Context) {
	locID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Decision string `json:"decision" binding:"required"`
		Note     string `json:"note"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate decision
	if body.Decision != "approved" && body.Decision != "rejected" && body.Decision != "revision" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "decision must be approved, rejected, or revision"})
		return
	}
	if body.Decision == "rejected" && body.Note == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "alasan penolakan wajib diisi"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	approverID := userID

	approval := &entity.LocationApproval{
		LocationID: locID,
		ApprovedBy: approverID,
		Decision:   body.Decision,
		Note:       body.Note,
		ApprovedAt: time.Now(),
	}
	if err := h.repo.CreateApproval(c.Request.Context(), approval); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update status
	statusMap := map[string]string{
		"approved": entity.LocSubStatusApproved,
		"rejected": entity.LocSubStatusRejected,
		"revision": entity.LocSubStatusRevision,
	}
	_ = h.repo.UpdateStatus(c.Request.Context(), locID, statusMap[body.Decision])

	ls, _ := h.repo.FindByID(c.Request.Context(), locID)
	c.JSON(http.StatusOK, gin.H{"data": ls})
}

// ─── SCORING ENGINE ───
func (h *LocationSubmissionHandler) calculateScore(ls *entity.LocationSubmission) {
	// Traffic score (30%) — based on estimasi_lalu_lintas
	switch {
	case ls.EstimasiLaluLintas >= 5000:
		ls.ScoreTraffic = 100
	case ls.EstimasiLaluLintas >= 3000:
		ls.ScoreTraffic = 85
	case ls.EstimasiLaluLintas >= 1000:
		ls.ScoreTraffic = 70
	case ls.EstimasiLaluLintas >= 500:
		ls.ScoreTraffic = 55
	default:
		ls.ScoreTraffic = 30
	}

	// Rental score (20%) — lower is better
	switch {
	case ls.HargaSewa <= 30_000_000:
		ls.ScoreSewa = 100
	case ls.HargaSewa <= 60_000_000:
		ls.ScoreSewa = 80
	case ls.HargaSewa <= 100_000_000:
		ls.ScoreSewa = 60
	case ls.HargaSewa <= 150_000_000:
		ls.ScoreSewa = 40
	default:
		ls.ScoreSewa = 20
	}

	// Competitor score (20%) — fewer is better
	switch {
	case ls.JumlahKompetitor == 0:
		ls.ScoreKompetitor = 100
	case ls.JumlahKompetitor <= 2:
		ls.ScoreKompetitor = 80
	case ls.JumlahKompetitor <= 5:
		ls.ScoreKompetitor = 60
	case ls.JumlahKompetitor <= 10:
		ls.ScoreKompetitor = 40
	default:
		ls.ScoreKompetitor = 20
	}

	// Access score (15%) — based on lebar_jalan
	switch {
	case ls.LebarJalan >= 12:
		ls.ScoreAkses = 100
	case ls.LebarJalan >= 8:
		ls.ScoreAkses = 80
	case ls.LebarJalan >= 5:
		ls.ScoreAkses = 60
	case ls.LebarJalan >= 3:
		ls.ScoreAkses = 40
	default:
		ls.ScoreAkses = 20
	}

	// Market score (15%) — based on dekat_dengan content
	marketPoints := 0
	dekat := ls.DekatDengan
	if contains(dekat, "kampus") || contains(dekat, "universitas") {
		marketPoints += 25
	}
	if contains(dekat, "sekolah") {
		marketPoints += 20
	}
	if contains(dekat, "pasar") || contains(dekat, "mall") {
		marketPoints += 20
	}
	if contains(dekat, "perumahan") || contains(dekat, "apartemen") {
		marketPoints += 20
	}
	if contains(dekat, "perkantoran") || contains(dekat, "kantor") {
		marketPoints += 15
	}
	if marketPoints > 100 {
		marketPoints = 100
	}
	if marketPoints == 0 {
		marketPoints = 30
	}
	ls.ScoreMarket = marketPoints

	// Weighted total
	total := float64(ls.ScoreTraffic)*0.30 +
		float64(ls.ScoreSewa)*0.20 +
		float64(ls.ScoreKompetitor)*0.20 +
		float64(ls.ScoreAkses)*0.15 +
		float64(ls.ScoreMarket)*0.15
	ls.TotalScore = int(total)

	// Category
	switch {
	case ls.TotalScore >= 80:
		ls.ScoreCategory = "Sangat Layak"
	case ls.TotalScore >= 65:
		ls.ScoreCategory = "Layak"
	case ls.TotalScore >= 50:
		ls.ScoreCategory = "Pertimbangan"
	default:
		ls.ScoreCategory = "Tidak Layak"
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) >= len(substr) &&
			(func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			})())
}

// ─── RECALCULATE SCORE ───
func (h *LocationSubmissionHandler) RecalculateScore(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	ls, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	h.calculateScore(ls)
	_ = h.repo.UpdateScore(c.Request.Context(), id,
		ls.ScoreTraffic, ls.ScoreSewa, ls.ScoreKompetitor, ls.ScoreAkses, ls.ScoreMarket,
		ls.TotalScore, ls.ScoreCategory)
	updated, _ := h.repo.FindByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// ─── GET BY MITRA ───
func (h *LocationSubmissionHandler) GetByMitra(c *gin.Context) {
	mitraID := c.MustGet("user_id").(uuid.UUID)
	list, err := h.repo.FindByMitraID(c.Request.Context(), mitraID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
}
