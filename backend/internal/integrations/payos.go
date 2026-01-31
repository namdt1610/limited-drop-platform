/**
 * PAYOS SERVICE
 *
 * PayOS payment gateway integration
 */

package integrations

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type PayOSCheckoutRequest struct {
	OrderCode   int64             `json:"orderCode"`
	Amount      int64             `json:"amount"` // VND
	Description string            `json:"description"`
	ReturnURL   string            `json:"returnUrl"`
	CancelURL   string            `json:"cancelUrl"`
	Items       []PayOSItem       `json:"items"`
	ExpiredAt   *int64            `json:"expiredAt,omitempty"` // Unix timestamp
	Metadata    map[string]string `json:"metadata,omitempty"`
	Signature   string            `json:"signature"` // HMAC_SHA256 signature
}

type PayOSItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int64  `json:"price"` // VND
}

type PayOSCheckoutResponse struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
	Data struct {
		Bin           string `json:"bin"`
		AccountNumber string `json:"accountNumber"`
		AccountName   string `json:"accountName"`
		Amount        int64  `json:"amount"`
		Description   string `json:"description"`
		OrderCode     int64  `json:"orderCode"`
		Currency      string `json:"currency"`
		PaymentLinkID string `json:"paymentLinkId"`
		QRCode        string `json:"qrCode"`
		CheckoutURL   string `json:"checkoutUrl"`
	} `json:"data"`
}

type PayOSVerifyRequest struct {
	OrderCode int64 `json:"orderCode"`
}

type PayOSVerifyResponse struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
	Data    struct {
		OrderCode   int64  `json:"orderCode"`
		Amount      int64  `json:"amount"`
		Status      string `json:"status"`
		Description string `json:"description"`
	} `json:"data"`
}

type PayOSRefundRequest struct {
	OrderCode   int64  `json:"orderCode"`
	Description string `json:"description"`
}

// CreatePayOSCheckout: Create PayOS checkout session
func CreatePayOSCheckout(req PayOSCheckoutRequest) (*PayOSCheckoutResponse, error) {
	clientID := os.Getenv("PAYOS_CLIENT_ID")
	apiKey := os.Getenv("PAYOS_API_KEY")
	checkoutURL := os.Getenv("PAYOS_CHECKOUT_URL")

	if clientID == "" || apiKey == "" {
		// Mock Mode for System Testing
		// Allows running load tests without real PayOS credentials
		return &PayOSCheckoutResponse{
			Code: "00",
			Desc: "Success (Mock)",
			Data: struct {
				Bin           string `json:"bin"`
				AccountNumber string `json:"accountNumber"`
				AccountName   string `json:"accountName"`
				Amount        int64  `json:"amount"`
				Description   string `json:"description"`
				OrderCode     int64  `json:"orderCode"`
				Currency      string `json:"currency"`
				PaymentLinkID string `json:"paymentLinkId"`
				QRCode        string `json:"qrCode"`
				CheckoutURL   string `json:"checkoutUrl"`
			}{
				CheckoutURL: "http://localhost:3000/mock-checkout",
			},
		}, nil
	}

	if checkoutURL == "" {
		// PayOS API v2 endpoint - correct domain from documentation
		checkoutURL = "https://api-merchant.payos.vn/v2/payment-requests"
	}

	// Generate order code if not provided
	if req.OrderCode == 0 {
		req.OrderCode = time.Now().Unix()
	}

	// Set default URLs
	if req.ReturnURL == "" {
		returnURL := os.Getenv("PAYOS_RETURN_URL")
		if returnURL == "" {
			returnURL = "http://localhost:5173/orders?payment=success"
		}
		req.ReturnURL = returnURL
	}
	if req.CancelURL == "" {
		cancelURL := os.Getenv("PAYOS_CANCEL_URL")
		if cancelURL == "" {
			cancelURL = "http://localhost:5173/checkout?payment=cancelled"
		}
		req.CancelURL = cancelURL
	}

	// Generate signature using checksum key
	checksumKey := os.Getenv("PAYOS_CHECKSUM_KEY")
	if checksumKey == "" {
		return nil, fmt.Errorf("PAYOS_CHECKSUM_KEY not configured")
	}

	// Create signature data: sort alphabetically
	signatureData := fmt.Sprintf("amount=%d&cancelUrl=%s&description=%s&orderCode=%d&returnUrl=%s",
		req.Amount, req.CancelURL, req.Description, req.OrderCode, req.ReturnURL)

	// Generate HMAC_SHA256 signature
	h := hmac.New(sha256.New, []byte(checksumKey))
	h.Write([]byte(signatureData))
	req.Signature = hex.EncodeToString(h.Sum(nil))

	// Marshal request
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	httpReq, err := http.NewRequest("POST", checkoutURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-client-id", clientID)
	httpReq.Header.Set("x-api-key", apiKey)

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PayOS error: %s", string(respBody))
	}

	var result PayOSCheckoutResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check response code (PayOS uses "code" field, "00" means success)
	if result.Code != "00" {
		return nil, fmt.Errorf("PayOS error: %s", result.Desc)
	}

	return &result, nil
}

// RefundPayOSPayment attempts to refund a completed PayOS payment.
// This is used in limited-drop race-condition scenarios where multiple users paid but stock was already taken.
func RefundPayOSPayment(orderCode int64, reason string) error {
	clientID := os.Getenv("PAYOS_CLIENT_ID")
	apiKey := os.Getenv("PAYOS_API_KEY")

	if clientID == "" || apiKey == "" {
		return fmt.Errorf("PayOS not configured")
	}

	refundURL := os.Getenv("PAYOS_REFUND_URL")
	if refundURL == "" {
		// Default to generic refund endpoint; adjust via PAYOS_REFUND_URL if PayOS changes path.
		refundURL = fmt.Sprintf("https://api-merchant.payos.vn/v2/payment-requests/%d/refunds", orderCode)
	}

	reqBody := PayOSRefundRequest{
		OrderCode:   orderCode,
		Description: reason,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal refund request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", refundURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create refund request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-client-id", clientID)
	httpReq.Header.Set("x-api-key", apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to call refund API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read refund response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("PayOS refund error: %s", string(respBody))
	}

	return nil
}

// VerifyPayOSPayment: Verify PayOS payment
func VerifyPayOSPayment(orderCode int64) (*PayOSVerifyResponse, error) {
	clientID := os.Getenv("PAYOS_CLIENT_ID")
	apiKey := os.Getenv("PAYOS_API_KEY")

	if clientID == "" || apiKey == "" {
		return nil, fmt.Errorf("PayOS not configured")
	}

	// Allow override for testing
	baseURL := os.Getenv("PAYOS_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api-merchant.payos.vn/v2"
	}
	verifyURL := fmt.Sprintf("%s/payment-requests/%d", baseURL, orderCode)

	httpReq, err := http.NewRequest("GET", verifyURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("x-client-id", clientID)
	httpReq.Header.Set("x-api-key", apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PayOS error: %s", string(respBody))
	}

	var result PayOSVerifyResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// CancelPayOSPayment: Cancel pending PayOS payment
func CancelPayOSPayment(orderCode int64) error {
	clientID := os.Getenv("PAYOS_CLIENT_ID")
	apiKey := os.Getenv("PAYOS_API_KEY")

	if clientID == "" || apiKey == "" {
		return fmt.Errorf("PayOS not configured")
	}

	// Allow override for testing
	baseURL := os.Getenv("PAYOS_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api-merchant.payos.vn/v2"
	}
	cancelURL := fmt.Sprintf("%s/payment-requests/%d/cancel", baseURL, orderCode)

	httpReq, err := http.NewRequest("POST", cancelURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create cancel request: %w", err)
	}

	httpReq.Header.Set("x-client-id", clientID)
	httpReq.Header.Set("x-api-key", apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to call cancel API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read cancel response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("PayOS cancel error: %s", string(respBody))
	}

	return nil
}

// GeneratePayOSSignature: Generate PayOS webhook signature
func GeneratePayOSSignature(data string) string {
	apiKey := os.Getenv("PAYOS_API_KEY")
	if apiKey == "" {
		return ""
	}

	h := hmac.New(sha256.New, []byte(apiKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// Mock implementations for Redis-less operation
// These would normally use Redis for mapping order codes to request IDs
