package midtrans

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/franchise-system/backend/internal/repository"
)

type Service struct {
	settingRepo repository.SystemSettingRepository
}

func NewService(settingRepo repository.SystemSettingRepository) *Service {
	return &Service{settingRepo: settingRepo}
}

// SnapRequest — full Midtrans Snap API request
type SnapRequest struct {
	TransactionDetails TransactionDetails `json:"transaction_details"`
	CustomerDetails    *CustomerDetails   `json:"customer_details,omitempty"`
	ItemDetails        []ItemDetail       `json:"item_details,omitempty"`
	EnabledPayments    []string           `json:"enabled_payments,omitempty"`
	CreditCard         *CreditCardConfig  `json:"credit_card,omitempty"`
	Expiry             *ExpiryConfig      `json:"expiry,omitempty"`
}

type TransactionDetails struct {
	OrderID  string `json:"order_id"`
	GrossAmt int64  `json:"gross_amount"`
}

type CustomerDetails struct {
	FirstName string `json:"first_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type ItemDetail struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Quantity int    `json:"quantity"`
}

// CreditCardConfig — credit card specific configuration
type CreditCardConfig struct {
	Secure      bool               `json:"secure"`            // Enable 3D Secure
	Bank        string             `json:"bank,omitempty"`    // Acquiring bank
	Channel     string             `json:"channel,omitempty"` // Payment channel
	Type        string             `json:"type,omitempty"`    // authorize or capture
	Installment *InstallmentConfig `json:"installment,omitempty"`
}

// InstallmentConfig — credit card installment config
type InstallmentConfig struct {
	Required bool             `json:"required"`
	Terms    map[string][]int `json:"terms"` // e.g. {"bca": [3, 6, 12], "mandiri": [3, 6, 12]}
}

// ExpiryConfig — transaction expiry
type ExpiryConfig struct {
	StartTime string `json:"start_time,omitempty"` // format: "2026-02-23 18:00:00 +0700"
	Unit      string `json:"unit"`                 // minute, hour, day
	Duration  int    `json:"duration"`
}

type SnapResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

// DefaultEnabledPayments — all supported payment methods
// Credit Card, PayLater (Kredivo, Akulaku), Bank Transfer, E-Wallet, QRIS, etc.
var DefaultEnabledPayments = []string{
	// Credit/Debit Card
	"credit_card",
	// Bank Transfer (VA)
	"bca_va", "bni_va", "bri_va", "permata_va", "other_va",
	// E-Channel (Mandiri Bill)
	"echannel",
	// E-Wallets
	"gopay", "shopeepay", "dana",
	// QRIS
	"qris",
	// Convenience Store
	"indomaret", "alfamart",
	// PayLater / Cardless Credit
	"kredivo", "akulaku",
}

func (s *Service) getConfig(ctx context.Context) (serverKey, baseURL string, err error) {
	serverKeySetting, err := s.settingRepo.FindByKey(ctx, "midtrans_server_key")
	if err != nil {
		return "", "", fmt.Errorf("midtrans server key not configured: %w", err)
	}
	envSetting, err := s.settingRepo.FindByKey(ctx, "midtrans_environment")
	if err != nil {
		return "", "", fmt.Errorf("midtrans environment not configured: %w", err)
	}

	serverKey = serverKeySetting.Value
	baseURL = "https://app.sandbox.midtrans.com"
	if envSetting.Value == "production" {
		baseURL = "https://app.midtrans.com"
	}
	return
}

// GetClientKey — returns the Midtrans client key for frontend Snap.js
func (s *Service) GetClientKey(ctx context.Context) (clientKey, snapURL string, err error) {
	clientKeySetting, err := s.settingRepo.FindByKey(ctx, "midtrans_client_key")
	if err != nil {
		return "", "", fmt.Errorf("midtrans client key not configured: %w", err)
	}
	envSetting, err := s.settingRepo.FindByKey(ctx, "midtrans_environment")
	if err != nil {
		return "", "", fmt.Errorf("midtrans environment not configured: %w", err)
	}

	clientKey = clientKeySetting.Value
	snapURL = "https://app.sandbox.midtrans.com/snap/snap.js"
	if envSetting.Value == "production" {
		snapURL = "https://app.midtrans.com/snap/snap.js"
	}
	return
}

func (s *Service) CreateSnapTransaction(ctx context.Context, req SnapRequest) (*SnapResponse, error) {
	serverKey, baseURL, err := s.getConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Apply defaults if not set
	if len(req.EnabledPayments) == 0 {
		req.EnabledPayments = DefaultEnabledPayments
	}
	if req.CreditCard == nil {
		req.CreditCard = &CreditCardConfig{
			Secure: true, // Always enable 3D Secure
		}
	}
	if req.Expiry == nil {
		req.Expiry = &ExpiryConfig{
			Unit:     "hour",
			Duration: 24,
		}
	}

	// Marshal request body
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/snap/v1/transactions", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Basic auth with server key
	authStr := base64.StdEncoding.EncodeToString([]byte(serverKey + ":"))
	httpReq.Header.Set("Authorization", "Basic "+authStr)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Execute
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("midtrans request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("midtrans error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var snapResp SnapResponse
	if err := json.Unmarshal(respBody, &snapResp); err != nil {
		return nil, fmt.Errorf("failed to parse midtrans response: %w", err)
	}

	return &snapResp, nil
}

// TransactionStatusResponse — response from Midtrans Status API
type TransactionStatusResponse struct {
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	PaymentType       string `json:"payment_type"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	FraudStatus       string `json:"fraud_status"`
	TransactionTime   string `json:"transaction_time"`
	SettlementTime    string `json:"settlement_time"`
	ExpiryTime        string `json:"expiry_time"`
	StatusMessage     string `json:"status_message"`
}

// GetTransactionStatus — check real-time status from Midtrans API
func (s *Service) GetTransactionStatus(ctx context.Context, orderID string) (*TransactionStatusResponse, error) {
	serverKey, _, err := s.getConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Midtrans status API uses api.sandbox.midtrans.com, not app.sandbox.midtrans.com
	envSetting, _ := s.settingRepo.FindByKey(ctx, "midtrans_environment")
	apiBase := "https://api.sandbox.midtrans.com"
	if envSetting != nil && envSetting.Value == "production" {
		apiBase = "https://api.midtrans.com"
	}

	url := fmt.Sprintf("%s/v2/%s/status", apiBase, orderID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	authStr := base64.StdEncoding.EncodeToString([]byte(serverKey + ":"))
	req.Header.Set("Authorization", "Basic "+authStr)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("midtrans status check failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var status TransactionStatusResponse
	if err := json.Unmarshal(respBody, &status); err != nil {
		return nil, fmt.Errorf("failed to parse midtrans status: %w", err)
	}

	return &status, nil
}
