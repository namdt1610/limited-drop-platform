package integrations

// =============================================================================
// PAYMENT GATEWAY INTERFACE
// =============================================================================

// PaymentGateway handles payment operations
type PaymentGateway interface {
	// CreateCheckout creates a checkout session for payment
	CreateCheckout(req PayOSCheckoutRequest) (*PayOSCheckoutResponse, error)

	// VerifyPayment verifies a payment status
	VerifyPayment(orderCode int64) (*PayOSVerifyResponse, error)

	// RefundPayment refunds a completed payment
	RefundPayment(orderCode int64, reason string) error

	// CancelPayment cancels a pending payment
	CancelPayment(orderCode int64) error

	// GenerateSignature generates webhook signature for verification
	GenerateSignature(data string) string
}

// =============================================================================
// EMAIL SENDER INTERFACE
// =============================================================================

// EmailSender handles email operations
type EmailSender interface {
	// SendOrderConfirmation sends order confirmation email
	SendOrderConfirmation(email, orderNumber string, amount float64) error

	// SendSymbioteReceipt sends ACCESS GRANTED/DENIED receipt
	SendSymbioteReceipt(email, phone, status, elapsed string) error

	// SendOrderDetails sends full order details (guest lookup)
	SendOrderDetails(email string, order interface{}) error
}

// =============================================================================
// GOOGLE SHEETS INTERFACE
// =============================================================================

// SheetSubmitter handles Google Sheets operations
type SheetSubmitter interface {
	// SubmitOrder submits order data to Google Sheet
	SubmitOrder(name, phone, email, address, notes string, amount float64, timestamp interface{}) error
}
