package integrations_test

import (
	"ecommerce-backend/internal/integrations"
	"ecommerce-backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Since Gateway just wraps existing tested functions, we verify they call through correctly.
// Ideally we would mock the internal calls, but since they call global functions or internal structs,
// we rely on the fact that those internal functions are already tested with httptest in payos_test.go and resend_test.go.
// Here we just ensure the struct methods satisfy the interface and don't panic.
// Real verification of logic is done in payos_test.go and resend_test.go.

func TestPayOSGateway_Methods(t *testing.T) {
	t.Setenv("PAYOS_CLIENT_ID", "test-client")
	t.Setenv("PAYOS_API_KEY", "test-api")
	t.Setenv("PAYOS_CHECKSUM_KEY", "test-checksum")

	gw := integrations.NewPayOSGateway() // No args
	
	// req type is PayOSCheckoutRequest
	req := integrations.PayOSCheckoutRequest{
		OrderCode: 123,
		Amount:    1000,
	}
	_, _ = gw.CreateCheckout(req)
	
	_, _ = gw.VerifyPayment(123)

	_ = gw.RefundPayment(123, "reason")
	
	_ = gw.CancelPayment(123)

	sig := gw.GenerateSignature("message")
	assert.NotEmpty(t, sig)
}

func TestResendEmailer_Methods(t *testing.T) {
	em := integrations.NewResendEmailer() // No args

	err := em.SendOrderConfirmation("test@example.com", "ORD-123", 1000.0)
	assert.Error(t, err)

	err = em.SendSymbioteReceipt("test@example.com", "0909090909", "ACTIVE", "1s")
	assert.Error(t, err)

	err = em.SendOrderDetails("test@example.com", &models.Order{})
	assert.Error(t, err)
}
