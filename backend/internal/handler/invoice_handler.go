package handler

import (
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
