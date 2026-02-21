package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	uploadDir string
	maxSize   int64
	baseURL   string
}

func NewUploadHandler(uploadDir string, maxSize int64, baseURL string) *UploadHandler {
	return &UploadHandler{
		uploadDir: uploadDir,
		maxSize:   maxSize,
		baseURL:   baseURL,
	}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File wajib diunggah"})
		return
	}

	if file.Size > h.maxSize {
		maxMB := h.maxSize / (1024 * 1024)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Ukuran file maksimal %d MB", maxMB)})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true, ".svg": true, ".pdf": true, ".doc": true, ".docx": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak didukung. Gunakan: jpg, png, webp, gif, svg, pdf, docx"})
		return
	}

	// Create date-based subdirectory
	subDir := time.Now().Format("2006/01")
	fullDir := filepath.Join(h.uploadDir, subDir)
	os.MkdirAll(fullDir, os.ModePerm)

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String()[:12], ext)
	savePath := filepath.Join(fullDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
		return
	}

	// Build public URL
	urlPath := fmt.Sprintf("/uploads/%s/%s", subDir, filename)
	fullURL := urlPath
	if h.baseURL != "" {
		fullURL = strings.TrimRight(h.baseURL, "/") + urlPath
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"url":      fullURL,
			"path":     urlPath,
			"filename": filename,
			"size":     file.Size,
		},
	})
}
