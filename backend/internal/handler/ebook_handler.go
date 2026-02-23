package handler

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/franchise-system/backend/internal/service/midtrans"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EbookHandler struct {
	ebookRepo      repository.EbookRepository
	ebookOrderRepo repository.EbookOrderRepository
	settingRepo    repository.SystemSettingRepository
	midtransSvc    *midtrans.Service
	uploadDir      string
}

func NewEbookHandler(
	ebookRepo repository.EbookRepository,
	ebookOrderRepo repository.EbookOrderRepository,
	settingRepo repository.SystemSettingRepository,
	midtransSvc *midtrans.Service,
	uploadDir string,
) *EbookHandler {
	return &EbookHandler{
		ebookRepo:      ebookRepo,
		ebookOrderRepo: ebookOrderRepo,
		settingRepo:    settingRepo,
		midtransSvc:    midtransSvc,
		uploadDir:      uploadDir,
	}
}

// ══════════════════════════════════════
// Admin CRUD
// ══════════════════════════════════════

type CreateEbookRequest struct {
	Title       string      `json:"title" binding:"required"`
	CategoryIDs []uuid.UUID `json:"category_ids"`
	Description string      `json:"description"`
	Author      string      `json:"author"`
	CoverURL    string      `json:"cover_url"`
	FileURL     string      `json:"file_url" binding:"required"`
	Price       int64       `json:"price"`
	IsActive    bool        `json:"is_active"`
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func (h *EbookHandler) Create(c *gin.Context) {
	var req CreateEbookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slug := slugify(req.Title)
	// Ensure unique slug
	if _, err := h.ebookRepo.FindBySlug(c.Request.Context(), slug); err == nil {
		slug = fmt.Sprintf("%s-%d", slug, time.Now().Unix())
	}

	ebook := &entity.Ebook{
		Title:       req.Title,
		Slug:        slug,
		Description: req.Description,
		Author:      req.Author,
		CoverURL:    req.CoverURL,
		FileURL:     req.FileURL,
		Price:       req.Price,
		IsActive:    req.IsActive,
	}

	if err := h.ebookRepo.Create(c.Request.Context(), ebook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat ebook"})
		return
	}

	// Sync categories
	if len(req.CategoryIDs) > 0 {
		h.ebookRepo.SyncCategories(c.Request.Context(), ebook.ID, req.CategoryIDs)
	}
	// Reload to include categories in response
	ebook, _ = h.ebookRepo.FindByID(c.Request.Context(), ebook.ID)

	c.JSON(http.StatusCreated, gin.H{"data": ebook})
}

func (h *EbookHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ebook, err := h.ebookRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ebook tidak ditemukan"})
		return
	}

	var req CreateEbookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newSlug := slugify(req.Title)
	if newSlug != ebook.Slug {
		if existing, err := h.ebookRepo.FindBySlug(c.Request.Context(), newSlug); err == nil && existing.ID != ebook.ID {
			newSlug = fmt.Sprintf("%s-%d", newSlug, time.Now().Unix())
		}
	}

	ebook.Title = req.Title
	ebook.Slug = newSlug
	ebook.Description = req.Description
	ebook.Author = req.Author
	ebook.CoverURL = req.CoverURL
	ebook.FileURL = req.FileURL
	ebook.Price = req.Price
	ebook.IsActive = req.IsActive

	if err := h.ebookRepo.Update(c.Request.Context(), ebook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update ebook"})
		return
	}

	// Sync categories
	h.ebookRepo.SyncCategories(c.Request.Context(), ebook.ID, req.CategoryIDs)
	// Reload to include categories
	ebook, _ = h.ebookRepo.FindByID(c.Request.Context(), ebook.ID)

	c.JSON(http.StatusOK, gin.H{"data": ebook})
}

func (h *EbookHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ebookRepo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus ebook"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ebook berhasil dihapus"})
}

func (h *EbookHandler) GetAll(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	ebooks, total, err := h.ebookRepo.FindAll(c.Request.Context(), false, search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if ebooks == nil {
		ebooks = []*entity.Ebook{}
	}
	c.JSON(http.StatusOK, gin.H{"data": ebooks, "total": total, "page": page, "limit": limit})
}

func (h *EbookHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	ebook, err := h.ebookRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ebook tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ebook})
}

func (h *EbookHandler) ToggleActive(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	ebook, err := h.ebookRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ebook tidak ditemukan"})
		return
	}
	ebook.IsActive = !ebook.IsActive
	if err := h.ebookRepo.Update(c.Request.Context(), ebook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ebook})
}

// ══════════════════════════════════════
// Admin — Download Approval
// ══════════════════════════════════════

func (h *EbookHandler) ListDownloadRequests(c *gin.Context) {
	orders, err := h.ebookOrderRepo.FindPendingDownloads(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if orders == nil {
		orders = []*entity.EbookOrder{}
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func (h *EbookHandler) ApproveDownload(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Note string `json:"note"`
	}
	c.ShouldBindJSON(&body)

	if err := h.ebookOrderRepo.ApproveDownload(c.Request.Context(), id, body.Note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Download disetujui"})
}

func (h *EbookHandler) RejectDownload(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Note string `json:"note" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alasan penolakan wajib diisi"})
		return
	}

	if err := h.ebookOrderRepo.RejectDownload(c.Request.Context(), id, body.Note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Download ditolak"})
}

func (h *EbookHandler) ListAllOrders(c *gin.Context) {
	status := c.Query("status")
	downloadStatus := c.Query("download_status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	orders, total, err := h.ebookOrderRepo.FindAllOrders(c.Request.Context(), status, downloadStatus, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if orders == nil {
		orders = []*entity.EbookOrder{}
	}
	c.JSON(http.StatusOK, gin.H{"data": orders, "total": total, "page": page, "limit": limit})
}

// ══════════════════════════════════════
// Mitra — Browse & Purchase
// ══════════════════════════════════════

func (h *EbookHandler) ListForMitra(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	ebooks, total, err := h.ebookRepo.FindAll(c.Request.Context(), true, search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type EbookWithPurchase struct {
		*entity.Ebook
		AlreadyPurchased bool   `json:"already_purchased"`
		DownloadStatus   string `json:"download_status,omitempty"`
	}

	result := make([]EbookWithPurchase, 0, len(ebooks))
	for _, e := range ebooks {
		item := EbookWithPurchase{Ebook: e}
		// Hide file_url from response
		item.FileURL = ""
		purchased, _ := h.ebookOrderRepo.HasUserPurchased(c.Request.Context(), userID, e.ID)
		item.AlreadyPurchased = purchased
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{"data": result, "total": total, "page": page, "limit": limit})
}

func (h *EbookHandler) GetForMitra(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ebook, err := h.ebookRepo.FindByID(c.Request.Context(), id)
	if err != nil || !ebook.IsActive {
		c.JSON(http.StatusNotFound, gin.H{"error": "ebook tidak ditemukan"})
		return
	}

	purchased, _ := h.ebookOrderRepo.HasUserPurchased(c.Request.Context(), userID, id)

	// Hide file URL
	ebook.FileURL = ""

	c.JSON(http.StatusOK, gin.H{
		"data":              ebook,
		"already_purchased": purchased,
	})
}

func (h *EbookHandler) Purchase(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ebook, err := h.ebookRepo.FindByID(c.Request.Context(), id)
	if err != nil || !ebook.IsActive {
		c.JSON(http.StatusNotFound, gin.H{"error": "ebook tidak ditemukan"})
		return
	}

	// Check if already purchased
	purchased, _ := h.ebookOrderRepo.HasUserPurchased(c.Request.Context(), userID, id)
	if purchased {
		c.JSON(http.StatusConflict, gin.H{"error": "Anda sudah membeli ebook ini"})
		return
	}

	orderNumber, err := h.ebookOrderRepo.GenerateOrderNumber(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate order number"})
		return
	}

	order := &entity.EbookOrder{
		EbookID:     id,
		UserID:      userID,
		OrderNumber: orderNumber,
		Amount:      ebook.Price,
		Status:      entity.EbookOrderStatusPending,
	}

	// Free ebook — auto-mark as paid
	if ebook.Price == 0 {
		now := time.Now()
		order.Status = entity.EbookOrderStatusPaid
		order.PaidAt = &now
		order.MidtransPaymentType = "free"

		if err := h.ebookOrderRepo.Create(c.Request.Context(), order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat order"})
			return
		}
		h.ebookRepo.IncrementSold(c.Request.Context(), id)
		c.JSON(http.StatusCreated, gin.H{"data": order, "message": "Ebook gratis berhasil diklaim"})
		return
	}

	// Paid ebook — create Midtrans transaction
	midtransOrderID := fmt.Sprintf("EB-%s-%d", orderNumber, time.Now().Unix())

	snapReq := midtrans.SnapRequest{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  midtransOrderID,
			GrossAmt: ebook.Price,
		},
		ItemDetails: []midtrans.ItemDetail{
			{
				ID:       ebook.ID.String(),
				Name:     fmt.Sprintf("Ebook: %s", ebook.Title),
				Price:    ebook.Price,
				Quantity: 1,
			},
		},
	}

	snapResp, err := h.midtransSvc.CreateSnapTransaction(c.Request.Context(), snapReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Midtrans error: %v", err)})
		return
	}

	expiry := time.Now().Add(24 * time.Hour)
	order.MidtransOrderID = midtransOrderID
	order.MidtransSnapToken = snapResp.Token
	order.MidtransRedirectURL = snapResp.RedirectURL
	order.ExpiredAt = &expiry

	if err := h.ebookOrderRepo.Create(c.Request.Context(), order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": order})
}

func (h *EbookHandler) MyOrders(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	orders, err := h.ebookOrderRepo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if orders == nil {
		orders = []*entity.EbookOrder{}
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func (h *EbookHandler) CancelOrder(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Verify ownership
	order, err := h.ebookOrderRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}
	if order.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses"})
		return
	}
	if order.Status != entity.EbookOrderStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya pesanan pending yang bisa dibatalkan"})
		return
	}

	if err := h.ebookOrderRepo.CancelOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pesanan berhasil dibatalkan"})
}

func (h *EbookHandler) UploadPaymentProof(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	order, err := h.ebookOrderRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("[UploadPaymentProof] FindByID error for %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Pesanan tidak ditemukan: %v", err)})
		return
	}
	if order.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses"})
		return
	}
	if order.Status != entity.EbookOrderStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya pesanan pending yang bisa diupload bukti"})
		return
	}

	var req struct {
		PaymentProofURL string `json:"payment_proof_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL bukti pembayaran wajib diisi"})
		return
	}

	if err := h.ebookOrderRepo.UploadPaymentProof(c.Request.Context(), id, req.PaymentProofURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload bukti"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bukti pembayaran berhasil diupload"})
}

// ══════════════════════════════════════
// Admin — Payment Approval
// ══════════════════════════════════════

func (h *EbookHandler) ApprovePayment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	order, err := h.ebookOrderRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}
	if order.Status != entity.EbookOrderStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya pesanan pending yang bisa diapprove"})
		return
	}
	if order.PaymentProofURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Belum ada bukti pembayaran"})
		return
	}

	if err := h.ebookOrderRepo.ApprovePayment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal approve pembayaran"})
		return
	}
	h.ebookRepo.IncrementSold(c.Request.Context(), order.EbookID)
	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran disetujui"})
}

func (h *EbookHandler) RejectPayment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req struct {
		Note string `json:"note"`
	}
	c.ShouldBindJSON(&req)

	if err := h.ebookOrderRepo.RejectPayment(c.Request.Context(), id, req.Note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal reject pembayaran"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran ditolak"})
}

// ══════════════════════════════════════
// Mitra — Read Online & Download
// ══════════════════════════════════════

func (h *EbookHandler) ReadOnline(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Verify purchase
	purchased, _ := h.ebookOrderRepo.HasUserPurchased(c.Request.Context(), userID, id)
	if !purchased {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda belum membeli ebook ini"})
		return
	}

	ebook, err := h.ebookRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ebook tidak ditemukan"})
		return
	}

	// Serve the PDF file for inline viewing
	filePath := ebook.FileURL
	// If file_url is a full URL (http://host/uploads/xxx.pdf), extract the /uploads/ path
	if strings.HasPrefix(filePath, "http") {
		if idx := strings.Index(filePath, "/uploads/"); idx != -1 {
			filePath = filePath[idx:]
		}
	}
	// If file_url is a relative path (from uploads), resolve it
	if !strings.HasPrefix(filePath, "/") {
		filePath = h.uploadDir + "/" + filePath
	}
	// If it's an absolute URL stored (like /uploads/xxx.pdf), strip the /uploads/ prefix
	if strings.HasPrefix(filePath, "/uploads/") {
		filePath = h.uploadDir + "/" + strings.TrimPrefix(filePath, "/uploads/")
	}

	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File ebook tidak ditemukan"})
		return
	}
	defer file.Close()

	stat, _ := file.Stat()
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline") // inline = view in browser, not download
	c.Header("Content-Length", fmt.Sprintf("%d", stat.Size()))
	c.Header("Cache-Control", "no-store") // prevent caching
	io.Copy(c.Writer, file)
}

func (h *EbookHandler) RequestDownload(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	order, err := h.ebookOrderRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order tidak ditemukan"})
		return
	}
	if order.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "bukan order milik Anda"})
		return
	}

	if err := h.ebookOrderRepo.RequestDownload(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request download berhasil dikirim, menunggu persetujuan admin"})
}

func (h *EbookHandler) Download(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Find order by ebook ID and verify download is approved
	orders, err := h.ebookOrderRepo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var approvedOrder *entity.EbookOrder
	for _, o := range orders {
		if o.EbookID == id && o.Status == entity.EbookOrderStatusPaid && o.DownloadStatus == entity.DownloadStatusApproved {
			approvedOrder = o
			break
		}
	}

	if approvedOrder == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Download belum disetujui atau ebook belum dibeli"})
		return
	}

	ebook, err := h.ebookRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ebook tidak ditemukan"})
		return
	}

	// Serve the PDF file for download
	filePath := ebook.FileURL
	if !strings.HasPrefix(filePath, "/") && !strings.HasPrefix(filePath, "http") {
		filePath = h.uploadDir + "/" + filePath
	}
	if strings.HasPrefix(filePath, "/uploads/") {
		filePath = h.uploadDir + "/" + strings.TrimPrefix(filePath, "/uploads/")
	}

	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File ebook tidak ditemukan"})
		return
	}
	defer file.Close()

	stat, _ := file.Stat()
	fileName := fmt.Sprintf("%s.pdf", slugify(ebook.Title))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("Content-Length", fmt.Sprintf("%d", stat.Size()))
	io.Copy(c.Writer, file)
}

// ══════════════════════════════════════
// Webhook — Midtrans notification for ebook orders
// ══════════════════════════════════════

func (h *EbookHandler) EbookWebhook(c *gin.Context) {
	var notification map[string]interface{}
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID, _ := notification["order_id"].(string)
	txnStatus, _ := notification["transaction_status"].(string)
	paymentType, _ := notification["payment_type"].(string)
	txnID, _ := notification["transaction_id"].(string)
	statusCode, _ := notification["status_code"].(string)
	grossAmount, _ := notification["gross_amount"].(string)
	signatureKey, _ := notification["signature_key"].(string)

	log.Printf("[Ebook Webhook] order_id=%s status=%s payment_type=%s", orderID, txnStatus, paymentType)

	// Verify signature
	serverKeySetting, err := h.settingRepo.FindByKey(c.Request.Context(), "midtrans_server_key")
	if err == nil {
		hash := sha512.New()
		hash.Write([]byte(orderID + statusCode + grossAmount + serverKeySetting.Value))
		expected := hex.EncodeToString(hash.Sum(nil))
		if expected != signatureKey {
			log.Printf("[Ebook Webhook] signature mismatch for order %s", orderID)
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
			return
		}
	}

	// Update ebook order status
	err = h.ebookOrderRepo.UpdateMidtransStatus(c.Request.Context(), orderID, txnStatus, paymentType, txnID)
	if err != nil {
		log.Printf("[Ebook Webhook] failed to update order %s: %v", orderID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If paid, increment ebook sold count
	if txnStatus == "capture" || txnStatus == "settlement" {
		order, err := h.ebookOrderRepo.FindByMidtransOrderID(c.Request.Context(), orderID)
		if err == nil {
			h.ebookRepo.IncrementSold(c.Request.Context(), order.EbookID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
