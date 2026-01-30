package integrations_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"ecommerce-backend/internal/integrations"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPayOSGatewayWrapper(t *testing.T) {
	// Mock PayOS Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":"00","desc":"success","data":{"checkoutUrl":"http://mock","status":"PAID"}}`))
	}))
	defer server.Close()

	os.Setenv("PAYOS_CLIENT_ID", "test")
	os.Setenv("PAYOS_API_KEY", "test")
	os.Setenv("PAYOS_CHECKSUM_KEY", "test")
	os.Setenv("PAYOS_BASE_URL", server.URL)
	os.Setenv("PAYOS_CHECKOUT_URL", server.URL+"/v2/payment-requests")
	// For refund
	os.Setenv("PAYOS_REFUND_URL", server.URL+"/refunds")
	defer os.Unsetenv("PAYOS_BASE_URL")
	defer os.Unsetenv("PAYOS_CHECKOUT_URL")

	gw := integrations.NewPayOSGateway()

	t.Run("CreateCheckout", func(t *testing.T) {
		resp, err := gw.CreateCheckout(integrations.PayOSCheckoutRequest{OrderCode: 123})
		require.NoError(t, err)
		assert.Equal(t, "http://mock", resp.Data.CheckoutURL)
	})

	t.Run("VerifyPayment", func(t *testing.T) {
		resp, err := gw.VerifyPayment(123)
		require.NoError(t, err)
		assert.Equal(t, "PAID", resp.Data.Status)
	})

	t.Run("RefundPayment", func(t *testing.T) {
		err := gw.RefundPayment(123, "reason")
		require.NoError(t, err)
	})

	// GenerateSignature just calls util, safe to skip mock server
	t.Run("GenerateSignature", func(t *testing.T) {
		sig := gw.GenerateSignature("test")
		assert.NotEmpty(t, sig)
	})
}

// Since ResendEmailer wrapper calls helper functions that call Brevo,
// we need to mock Brevo for it.
func TestEmailWrapper_WithBrevoMock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"messageId":123}`))
	}))
	defer server.Close()

	os.Setenv("BREVO_API_KEY", "test")
	os.Setenv("BREVO_BASE_URL", server.URL)
	defer os.Unsetenv("BREVO_BASE_URL")

	emailer := integrations.NewResendEmailer()

	t.Run("SendOrderConfirmation", func(t *testing.T) {
		err := emailer.SendOrderConfirmation("foo@bar.com", "123", 100)
		require.NoError(t, err)
	})

	t.Run("SendSymbioteReceipt", func(t *testing.T) {
		err := emailer.SendSymbioteReceipt("foo@bar.com", "123", "WINNER", "1s")
		require.NoError(t, err)
	})
}

func TestSheetsWrapper(t *testing.T) {
	// Sheets wrapper is hard to mock (Google SDK), so we just test the No-Op case
	// ensuring it doesn't crash when env vars are missing.
	os.Unsetenv("GSSHEET_SPREADSHEET_ID")
	
	sheets := integrations.NewSheetsSubmitter()
	err := sheets.SubmitOrder("Name", "Phone", "Email", "Addr", "Drop", 100, time.Now())
	require.NoError(t, err)
}
