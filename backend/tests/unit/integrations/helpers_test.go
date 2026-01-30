package integrations_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"ecommerce-backend/internal/integrations"
	"ecommerce-backend/internal/models"

	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func TestHelperFunctions(t *testing.T) {
	// Mock Brevo Server (since helpers call SendEmailBrevo)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"messageId":12345}`))
	}))
	defer server.Close()

	os.Setenv("BREVO_API_KEY", "test-key")
	os.Setenv("BREVO_BASE_URL", server.URL)
	defer os.Unsetenv("BREVO_BASE_URL")

	t.Run("SendWelcomeEmail", func(t *testing.T) {
		err := integrations.SendWelcomeEmail("test@example.com", "Test User")
		require.NoError(t, err)
	})

	t.Run("SendOrderConfirmationEmail", func(t *testing.T) {
		err := integrations.SendOrderConfirmationEmail("test@example.com", "ORDER-123", 100000)
		require.NoError(t, err)
	})

	t.Run("SendSymbioteReceipt", func(t *testing.T) {
		err := integrations.SendSymbioteReceipt("test@example.com", "09**99", "WINNER", "1s")
		require.NoError(t, err)
	})

	t.Run("SendPasswordResetEmail", func(t *testing.T) {
		err := integrations.SendPasswordResetEmail("test@example.com", "token123")
		require.NoError(t, err)
	})

	t.Run("SendOrderCreatedAdminEmail", func(t *testing.T) {
		// Mock order
		order := &models.Order{
			ID:              1,
			TotalAmount:     500000,
			Status:          models.OrderPending,
			ShippingAddress: datatypes.JSON(`{"name":"Admin Test","email":"test@admin.com"}`),
			Items:           datatypes.JSON(`[{"product_name":"Test Product","quantity":1,"price":500000}]`),
		}
		// Admin email recipients are fetched from env or default
		// Default is hardcoded, so it should work.
		err := integrations.SendOrderCreatedAdminEmail(order)
		require.NoError(t, err)
	})
}
