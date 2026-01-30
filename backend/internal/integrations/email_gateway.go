package integrations

import (
	"ecommerce-backend/internal/models"
	"time"
)

// =============================================================================
// RESEND EMAIL SENDER IMPLEMENTATION
// =============================================================================

// resendEmailer implements EmailSender interface
type resendEmailer struct{}

// NewResendEmailer creates a new Resend email sender
func NewResendEmailer() EmailSender {
	return &resendEmailer{}
}

func (r *resendEmailer) SendOrderConfirmation(email, orderNumber string, amount float64) error {
	return SendOrderConfirmationEmail(email, orderNumber, amount)
}

func (r *resendEmailer) SendSymbioteReceipt(email, phone, status, elapsed string) error {
	return SendSymbioteReceipt(email, phone, status, elapsed)
}

func (r *resendEmailer) SendOrderDetails(email string, order interface{}) error {
	if o, ok := order.(*models.Order); ok {
		return SendOrderDetailsEmail(email, o)
	}
	return nil
}

// =============================================================================
// GOOGLE SHEETS SUBMITTER IMPLEMENTATION
// =============================================================================

// sheetsSubmitter implements SheetSubmitter interface
type sheetsSubmitter struct{}

// NewSheetsSubmitter creates a new Google Sheets submitter
func NewSheetsSubmitter() SheetSubmitter {
	return &sheetsSubmitter{}
}

func (s *sheetsSubmitter) SubmitOrder(name, phone, email, address, notes string, amount float64, timestamp interface{}) error {
	if t, ok := timestamp.(time.Time); ok {
		return SubmitOrderToGoogleSheet(name, phone, email, address, notes, amount, t)
	}
	return nil
}
