package handler

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/franchise-system/backend/internal/service/midtrans"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InvoiceHandler struct {
	invoiceRepo     repository.InvoiceRepository
	partnershipRepo repository.PartnershipRepository
	settingRepo     repository.SystemSettingRepository
	midtransSvc     *midtrans.Service
}

func NewInvoiceHandler(
	invoiceRepo repository.InvoiceRepository,
	partnershipRepo repository.PartnershipRepository,
	settingRepo repository.SystemSettingRepository,
	midtransSvc *midtrans.Service,
) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceRepo:     invoiceRepo,
		partnershipRepo: partnershipRepo,
		settingRepo:     settingRepo,
		midtransSvc:     midtransSvc,
	}
}

type CreateInvoiceRequest struct {
	PartnershipID string `json:"partnership_id" binding:"required"`
	Type          string `json:"type"`
	Amount        int64  `json:"amount" binding:"required"`
	Description   string `json:"description"`
}

func (h *InvoiceHandler) Create(c *gin.Context) {
	var req CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	partnershipID, err := uuid.Parse(req.PartnershipID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid partnership_id"})
		return
	}

	// Get partnership with mitra info
	partnership, err := h.partnershipRepo.FindByID(c.Request.Context(), partnershipID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "partnership not found"})
		return
	}

	// Generate invoice number
	invoiceNumber, err := h.invoiceRepo.GenerateInvoiceNumber(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate invoice number"})
		return
	}

	orderID := fmt.Sprintf("BO-%s-%d", invoiceNumber, time.Now().Unix())

	// Call Midtrans Snap
	snapReq := midtrans.SnapRequest{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: req.Amount,
		},
		ItemDetails: []midtrans.ItemDetail{
			{
				ID:       invoiceNumber,
				Name:     req.Description,
				Price:    req.Amount,
				Quantity: 1,
			},
		},
	}

	// Add customer details if available
	if partnership.Mitra != nil {
		snapReq.CustomerDetails = &midtrans.CustomerDetails{
			FirstName: partnership.Mitra.Name,
			Email:     partnership.Mitra.Email,
			Phone:     partnership.Mitra.Phone,
		}
	}

	snapResp, err := h.midtransSvc.CreateSnapTransaction(c.Request.Context(), snapReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Midtrans error: %v", err)})
		return
	}

	// Set expiry 24 hours from now
	expiry := time.Now().Add(24 * time.Hour)

	invType := req.Type
	if invType == "" {
		invType = entity.InvoiceTypeInvoice
	}

	// Save invoice
	inv := &entity.Invoice{
		PartnershipID:       partnershipID,
		InvoiceNumber:       invoiceNumber,
		Type:                invType,
		Amount:              req.Amount,
		Description:         req.Description,
		Status:              entity.InvoiceStatusPending,
		MidtransOrderID:     orderID,
		MidtransSnapToken:   snapResp.Token,
		MidtransRedirectURL: snapResp.RedirectURL,
		ExpiredAt:           &expiry,
	}

	if err := h.invoiceRepo.Create(c.Request.Context(), inv); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invoice"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": inv,
	})
}

func (h *InvoiceHandler) GetByPartnership(c *gin.Context) {
	id, err := uuid.Parse(c.Param("partnership_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid partnership_id"})
		return
	}
	invoices, err := h.invoiceRepo.FindByPartnershipID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if invoices == nil {
		invoices = []*entity.Invoice{}
	}
	c.JSON(http.StatusOK, gin.H{"data": invoices})
}

func (h *InvoiceHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	inv, err := h.invoiceRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": inv})
}

// ManualApprove — admin manually marks invoice as paid
func (h *InvoiceHandler) ManualApprove(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		ProofURL string `json:"proof_url"`
	}
	c.ShouldBindJSON(&body)
	if body.ProofURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bukti pembayaran wajib diupload"})
		return
	}
	if err := h.invoiceRepo.ManualApprove(c.Request.Context(), id, body.ProofURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inv, _ := h.invoiceRepo.FindByID(c.Request.Context(), id)

	// Auto-update partnership progress
	if inv != nil {
		h.autoUpdatePartnershipProgress(c.Request.Context(), inv.PartnershipID)
	}

	c.JSON(http.StatusOK, gin.H{"data": inv, "message": "Invoice berhasil diverifikasi"})
}

// Webhook — called by Midtrans (no auth needed)
func (h *InvoiceHandler) MidtransWebhook(c *gin.Context) {
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

	log.Printf("[Midtrans Webhook] order_id=%s status=%s payment_type=%s", orderID, txnStatus, paymentType)

	// Verify signature
	serverKeySetting, err := h.settingRepo.FindByKey(c.Request.Context(), "midtrans_server_key")
	if err == nil {
		hash := sha512.New()
		hash.Write([]byte(orderID + statusCode + grossAmount + serverKeySetting.Value))
		expected := hex.EncodeToString(hash.Sum(nil))
		if expected != signatureKey {
			log.Printf("[Midtrans Webhook] signature mismatch for order %s", orderID)
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
			return
		}
	}

	// Update invoice status
	err = h.invoiceRepo.UpdateMidtransStatus(c.Request.Context(), orderID, txnStatus, paymentType, txnID)
	if err != nil {
		log.Printf("[Midtrans Webhook] failed to update order %s: %v", orderID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Auto-update partnership progress when payment settles
	if txnStatus == "settlement" || txnStatus == "capture" {
		inv, findErr := h.invoiceRepo.FindByOrderID(c.Request.Context(), orderID)
		if findErr == nil && inv != nil {
			h.autoUpdatePartnershipProgress(c.Request.Context(), inv.PartnershipID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetByMitra — returns all invoices across mitra's partnerships
func (h *InvoiceHandler) GetByMitra(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	partnerships, err := h.partnershipRepo.FindByMitraID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var all []*entity.Invoice
	for _, p := range partnerships {
		invoices, err := h.invoiceRepo.FindByPartnershipID(c.Request.Context(), p.ID)
		if err != nil {
			continue
		}
		all = append(all, invoices...)
	}
	if all == nil {
		all = []*entity.Invoice{}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": all})
}

// GetMidtransClientKey — returns client key & snap URL for frontend Snap.js popup
func (h *InvoiceHandler) GetMidtransClientKey(c *gin.Context) {
	clientKey, snapURL, err := h.midtransSvc.GetClientKey(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"client_key": clientKey,
		"snap_url":   snapURL,
	})
}

// CheckStatus — proactively check a single invoice status from Midtrans
func (h *InvoiceHandler) CheckStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	inv, err := h.invoiceRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
		return
	}

	// If no midtrans order, or already final status, return current
	if inv.MidtransOrderID == "" || inv.Status == entity.InvoiceStatusPaid || inv.Status == entity.InvoiceStatusExpired || inv.Status == entity.InvoiceStatusFailed {
		c.JSON(http.StatusOK, gin.H{"data": inv, "synced": false, "message": "Status sudah final"})
		return
	}

	// Check expired_at locally first
	if inv.ExpiredAt != nil && time.Now().After(*inv.ExpiredAt) && inv.Status == entity.InvoiceStatusPending {
		_ = h.invoiceRepo.UpdateMidtransStatus(c.Request.Context(), inv.MidtransOrderID, "expire", "", "")
		inv, _ = h.invoiceRepo.FindByID(c.Request.Context(), id)
		c.JSON(http.StatusOK, gin.H{"data": inv, "synced": true, "message": "Invoice telah kedaluwarsa"})
		return
	}

	// Call Midtrans API
	status, err := h.midtransSvc.GetTransactionStatus(c.Request.Context(), inv.MidtransOrderID)
	if err != nil {
		log.Printf("[CheckStatus] Midtrans API error for %s: %v", inv.MidtransOrderID, err)
		c.JSON(http.StatusOK, gin.H{"data": inv, "synced": false, "message": "Gagal cek status Midtrans"})
		return
	}

	// Update if status changed
	oldStatus := inv.Status
	if err := h.invoiceRepo.UpdateMidtransStatus(c.Request.Context(), inv.MidtransOrderID, status.TransactionStatus, status.PaymentType, status.TransactionID); err != nil {
		log.Printf("[CheckStatus] failed to update order %s: %v", inv.MidtransOrderID, err)
	}

	// Re-fetch updated invoice
	inv, _ = h.invoiceRepo.FindByID(c.Request.Context(), id)
	synced := inv.Status != oldStatus

	// Auto-update partnership progress if status changed to PAID
	if synced && inv.Status == entity.InvoiceStatusPaid {
		h.autoUpdatePartnershipProgress(c.Request.Context(), inv.PartnershipID)
	}

	message := "Status tidak berubah"
	if synced {
		message = fmt.Sprintf("Status berubah: %s → %s", oldStatus, inv.Status)
	}

	c.JSON(http.StatusOK, gin.H{"data": inv, "synced": synced, "message": message})
}

// SyncAllPending — batch sync all pending invoices with Midtrans
func (h *InvoiceHandler) SyncAllPending(c *gin.Context) {
	// First auto-expire past-due invoices
	expired, _ := h.invoiceRepo.ExpirePendingInvoices(c.Request.Context())

	// Then sync remaining pending invoices with Midtrans
	pendingInvoices, err := h.invoiceRepo.FindPendingWithMidtrans(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	synced := 0
	for _, inv := range pendingInvoices {
		status, err := h.midtransSvc.GetTransactionStatus(c.Request.Context(), inv.MidtransOrderID)
		if err != nil {
			log.Printf("[SyncPending] Midtrans error for %s: %v", inv.MidtransOrderID, err)
			continue
		}
		if err := h.invoiceRepo.UpdateMidtransStatus(c.Request.Context(), inv.MidtransOrderID, status.TransactionStatus, status.PaymentType, status.TransactionID); err != nil {
			log.Printf("[SyncPending] update error for %s: %v", inv.MidtransOrderID, err)
			continue
		}
		// Auto-update partnership progress on settlement
		if status.TransactionStatus == "settlement" || status.TransactionStatus == "capture" {
			h.autoUpdatePartnershipProgress(c.Request.Context(), inv.PartnershipID)
		}
		synced++
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"expired_count":  expired,
		"synced_count":   synced,
		"pending_checked": len(pendingInvoices),
		"message":        fmt.Sprintf("Expired: %d, Synced: %d/%d invoices", expired, synced, len(pendingInvoices)),
	})
}

// StartExpiryScheduler — background goroutine that auto-expires invoices every 5 minutes
func (h *InvoiceHandler) StartExpiryScheduler() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		// Run once immediately on startup
		h.runExpiryCheck()

		for range ticker.C {
			h.runExpiryCheck()
		}
	}()
	log.Println("🔄 Invoice expiry scheduler started (every 5 minutes)")
}

func (h *InvoiceHandler) runExpiryCheck() {
	ctx := context.Background()
	count, err := h.invoiceRepo.ExpirePendingInvoices(ctx)
	if err != nil {
		log.Printf("[ExpiryScheduler] error: %v", err)
		return
	}
	if count > 0 {
		log.Printf("[ExpiryScheduler] auto-expired %d invoice(s)", count)
	}
}

// autoUpdatePartnershipProgress calculates and updates partnership progress
// based on the current state of all invoices for that partnership.
func (h *InvoiceHandler) autoUpdatePartnershipProgress(ctx context.Context, partnershipID uuid.UUID) {
	invoices, err := h.invoiceRepo.FindByPartnershipID(ctx, partnershipID)
	if err != nil || len(invoices) == 0 {
		return
	}

	hasDPPaid := false
	allPaid := true
	for _, inv := range invoices {
		if inv.Status == entity.InvoiceStatusPaid {
			if inv.Type == entity.InvoiceTypeDP {
				hasDPPaid = true
			}
		} else if inv.Status == entity.InvoiceStatusPending {
			allPaid = false
		}
	}

	var progress int
	var status string

	if allPaid && hasDPPaid {
		// All invoices paid → at least DP_VERIFIED, check further
		progress = 50
		status = entity.PartnershipStatusAgreementSigned
	} else if hasDPPaid {
		// DP paid, other invoices pending
		progress = 25
		status = entity.PartnershipStatusDPVerified
	} else {
		// No DP paid yet
		return
	}

	// Only upgrade, never downgrade
	p, err := h.partnershipRepo.FindByID(ctx, partnershipID)
	if err != nil {
		return
	}

	if progress > p.ProgressPercentage {
		if err := h.partnershipRepo.UpdateProgress(ctx, partnershipID, progress, status); err != nil {
			log.Printf("[AutoProgress] failed to update partnership %s: %v", partnershipID, err)
		} else {
			log.Printf("[AutoProgress] partnership %s → %s (%d%%)", partnershipID, status, progress)
		}
	}
}
