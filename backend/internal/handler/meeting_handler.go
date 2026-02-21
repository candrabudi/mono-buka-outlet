package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MeetingHandler struct {
	meetingUC *usecase.MeetingUseCase
	uploadDir string
	baseURL   string
}

func NewMeetingHandler(uc *usecase.MeetingUseCase, uploadDir, baseURL string) *MeetingHandler {
	return &MeetingHandler{meetingUC: uc, uploadDir: uploadDir, baseURL: baseURL}
}

// ═══ MEETING CRUD ═══

func (h *MeetingHandler) Create(c *gin.Context) {
	var req usecase.CreateMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	userID := c.MustGet("user_id").(uuid.UUID)
	m, err := h.meetingUC.Create(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": m, "message": "Meeting berhasil dibuat"})
}

func (h *MeetingHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	meetingType := c.Query("type")

	list, total, err := h.meetingUC.GetAll(c.Request.Context(), status, meetingType, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": list, "total": total, "page": page, "limit": limit})
}

func (h *MeetingHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	m, err := h.meetingUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": m})
}

func (h *MeetingHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	var req usecase.UpdateMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	m, err := h.meetingUC.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": m, "message": "Meeting berhasil diupdate"})
}

func (h *MeetingHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	if err := h.meetingUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Meeting berhasil dihapus"})
}

// ═══ PARTICIPANTS ═══

func (h *MeetingHandler) AddParticipant(c *gin.Context) {
	meetingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	var req usecase.AddParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	p, err := h.meetingUC.AddParticipant(c.Request.Context(), meetingID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": p})
}

func (h *MeetingHandler) DeleteParticipant(c *gin.Context) {
	id, err := uuid.Parse(c.Param("participantId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	if err := h.meetingUC.DeleteParticipant(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Peserta berhasil dihapus"})
}

// ═══ NOTES ═══

func (h *MeetingHandler) SaveNotes(c *gin.Context) {
	meetingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	var req usecase.SaveNotesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	n, err := h.meetingUC.SaveNotes(c.Request.Context(), meetingID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": n, "message": "Notulensi berhasil disimpan"})
}

// ═══ ACTION PLANS ═══

func (h *MeetingHandler) AddActionPlan(c *gin.Context) {
	meetingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	var req usecase.AddActionPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	a, err := h.meetingUC.AddActionPlan(c.Request.Context(), meetingID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": a})
}

func (h *MeetingHandler) UpdateActionPlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("actionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	var req usecase.UpdateActionPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.meetingUC.UpdateActionPlan(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Action plan berhasil diupdate"})
}

func (h *MeetingHandler) DeleteActionPlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("actionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	if err := h.meetingUC.DeleteActionPlan(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Action plan berhasil dihapus"})
}

// ═══ FILE UPLOAD ═══

func (h *MeetingHandler) UploadFile(c *gin.Context) {
	meetingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "File wajib diunggah"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".pdf": true, ".docx": true, ".xlsx": true, ".jpg": true, ".jpeg": true, ".png": true, ".doc": true, ".xls": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Format file tidak didukung"})
		return
	}

	subDir := filepath.Join("meetings", time.Now().Format("2006/01"))
	fullDir := filepath.Join(h.uploadDir, subDir)
	os.MkdirAll(fullDir, os.ModePerm)

	filename := fmt.Sprintf("%s%s", uuid.New().String()[:12], ext)
	savePath := filepath.Join(fullDir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Gagal menyimpan file"})
		return
	}

	urlPath := fmt.Sprintf("/uploads/%s/%s", subDir, filename)
	userID := c.MustGet("user_id").(uuid.UUID)

	mf := &entity.MeetingFile{
		FileName:   file.Filename,
		FilePath:   urlPath,
		FileType:   ext[1:],
		UploadedBy: &userID,
	}
	if err := h.meetingUC.AddFile(c.Request.Context(), meetingID, mf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Gagal menyimpan data file"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": mf, "message": "File berhasil diunggah"})
}

func (h *MeetingHandler) DeleteFile(c *gin.Context) {
	id, err := uuid.Parse(c.Param("fileId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID tidak valid"})
		return
	}
	if err := h.meetingUC.DeleteFile(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "File berhasil dihapus"})
}
