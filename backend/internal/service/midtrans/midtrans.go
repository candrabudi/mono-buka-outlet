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

type SnapRequest struct {
	TransactionDetails TransactionDetails `json:"transaction_details"`
	CustomerDetails    *CustomerDetails   `json:"customer_details,omitempty"`
	ItemDetails        []ItemDetail       `json:"item_details,omitempty"`
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

type SnapResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

func (s *Service) CreateSnapTransaction(ctx context.Context, req SnapRequest) (*SnapResponse, error) {
	// Get settings from DB
	serverKeySetting, err := s.settingRepo.FindByKey(ctx, "midtrans_server_key")
	if err != nil {
		return nil, fmt.Errorf("midtrans server key not configured: %w", err)
	}
	envSetting, err := s.settingRepo.FindByKey(ctx, "midtrans_environment")
	if err != nil {
		return nil, fmt.Errorf("midtrans environment not configured: %w", err)
	}

	serverKey := serverKeySetting.Value
	env := envSetting.Value

	// Determine API URL
	baseURL := "https://app.sandbox.midtrans.com"
	if env == "production" {
		baseURL = "https://app.midtrans.com"
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
	auth := base64.StdEncoding.EncodeToString([]byte(serverKey + ":"))
	httpReq.Header.Set("Authorization", "Basic "+auth)
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
