package integrations_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"ecommerce-backend/internal/integrations"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePayOSCheckout_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		req            integrations.PayOSCheckoutRequest
		mockHandler    func(w http.ResponseWriter, r *http.Request)
		mockEnv        map[string]string
		wantErr        bool
		wantCheckoutURL string
	}{
		{
			name: "success",
			req: integrations.PayOSCheckoutRequest{
				OrderCode: 123456,
				Amount:    100000,
				Items:     []integrations.PayOSItem{{Name: "Test", Quantity: 1, Price: 100000}},
				ReturnURL: "http://localhost/return",
				CancelURL: "http://localhost/cancel",
			},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/v2/payment-requests", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				json.NewEncoder(w).Encode(integrations.PayOSCheckoutResponse{
					Code: "00", Desc: "success",
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
					}{CheckoutURL: "https://payos.vn/checkout-mock"},
				})
			},
			mockEnv: map[string]string{
				"PAYOS_CLIENT_ID":    "test-client",
				"PAYOS_API_KEY":      "test-key",
				"PAYOS_CHECKSUM_KEY": "test-checksum",
			},
			wantErr:         false,
			wantCheckoutURL: "https://payos.vn/checkout-mock",
		},
		{
			name: "error - missing config",
			req:  integrations.PayOSCheckoutRequest{OrderCode: 123},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {},
			mockEnv:     map[string]string{}, // Empty env
			wantErr:     true,
		},
		{
			name: "error - payos api error",
			req: integrations.PayOSCheckoutRequest{
				OrderCode: 123456,
				Amount:    100000,
			},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"code":"10","desc":"bad request"}`))
			},
			mockEnv: map[string]string{
				"PAYOS_CLIENT_ID":    "test-client",
				"PAYOS_API_KEY":      "test-key",
				"PAYOS_CHECKSUM_KEY": "test-checksum",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tc.mockHandler))
			defer server.Close()

			for k, v := range tc.mockEnv {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}
			os.Setenv("PAYOS_CHECKOUT_URL", server.URL+"/v2/payment-requests")
			defer os.Unsetenv("PAYOS_CHECKOUT_URL")

			resp, err := integrations.CreatePayOSCheckout(tc.req)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantCheckoutURL, resp.Data.CheckoutURL)
			}
		})
	}
}

func TestVerifyPayOSPayment_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		orderCode   int64
		mockHandler func(w http.ResponseWriter, r *http.Request)
		mockEnv     map[string]string
		wantErr     bool
		wantStatus  string
	}{
		{
			name:      "success - paid",
			orderCode: 123456,
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/payment-requests/123456", r.URL.Path)
				json.NewEncoder(w).Encode(integrations.PayOSVerifyResponse{
					Error: 0, Message: "success",
					Data: struct {
						OrderCode   int64  `json:"orderCode"`
						Amount      int64  `json:"amount"`
						Status      string `json:"status"`
						Description string `json:"description"`
					}{Status: "PAID"},
				})
			},
			mockEnv: map[string]string{
				"PAYOS_CLIENT_ID": "test-client",
				"PAYOS_API_KEY":   "test-key",
			},
			wantErr:    false,
			wantStatus: "PAID",
		},
		{
			name:      "error - payos error",
			orderCode: 123456,
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal server error"))
			},
			mockEnv: map[string]string{
				"PAYOS_CLIENT_ID": "test-client",
				"PAYOS_API_KEY":   "test-key",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tc.mockHandler))
			defer server.Close()

			for k, v := range tc.mockEnv {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}
			os.Setenv("PAYOS_BASE_URL", server.URL)
			defer os.Unsetenv("PAYOS_BASE_URL")

			resp, err := integrations.VerifyPayOSPayment(tc.orderCode)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantStatus, resp.Data.Status)
			}
		})
	}
}

func TestRefundPayOSPayment_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		orderCode   int64
		mockHandler func(w http.ResponseWriter, r *http.Request)
		mockEnv     map[string]string
		wantErr     bool
	}{
		{
			name:      "success",
			orderCode: 123456,
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/v2/payment-requests/123456/refunds", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"code":"00","desc":"success"}`))
			},
			mockEnv: map[string]string{
				"PAYOS_CLIENT_ID": "test-client",
				"PAYOS_API_KEY":   "test-key",
			},
			wantErr: false,
		},
		{
			name:      "error - refund failed",
			orderCode: 123456,
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"code":"10","desc":"failed"}`))
			},
			mockEnv: map[string]string{
				"PAYOS_CLIENT_ID": "test-client",
				"PAYOS_API_KEY":   "test-key",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tc.mockHandler))
			defer server.Close()

			for k, v := range tc.mockEnv {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}
			os.Setenv("PAYOS_REFUND_URL", server.URL+"/v2/payment-requests/123456/refunds")
			defer os.Unsetenv("PAYOS_REFUND_URL")

			err := integrations.RefundPayOSPayment(tc.orderCode, "reason")

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCancelPayOSPayment_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		orderCode   int64
		mockHandler func(w http.ResponseWriter, r *http.Request)
		mockEnv     map[string]string
		wantErr     bool
	}{
		{
			name:      "success",
			orderCode: 123456,
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/payment-requests/123456/cancel", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"code":"00","desc":"success"}`))
			},
			mockEnv: map[string]string{
				"PAYOS_CLIENT_ID": "test-client",
				"PAYOS_API_KEY":   "test-key",
			},
			wantErr: false,
		},
		{
			name:      "error - cancel failed",
			orderCode: 123456,
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"code":"10","desc":"failed"}`))
			},
			mockEnv: map[string]string{
				"PAYOS_CLIENT_ID": "test-client",
				"PAYOS_API_KEY":   "test-key",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tc.mockHandler))
			defer server.Close()

			for k, v := range tc.mockEnv {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}
			os.Setenv("PAYOS_BASE_URL", server.URL)
			defer os.Unsetenv("PAYOS_BASE_URL")

			err := integrations.CancelPayOSPayment(tc.orderCode)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
