package integrations

// =============================================================================
// PAYOS GATEWAY IMPLEMENTATION
// =============================================================================

// payosGateway implements PaymentGateway interface
type payosGateway struct{}

// NewPayOSGateway creates a new PayOS payment gateway
func NewPayOSGateway() PaymentGateway {
	return &payosGateway{}
}

func (p *payosGateway) CreateCheckout(req PayOSCheckoutRequest) (*PayOSCheckoutResponse, error) {
	return CreatePayOSCheckout(req)
}

func (p *payosGateway) VerifyPayment(orderCode int64) (*PayOSVerifyResponse, error) {
	return VerifyPayOSPayment(orderCode)
}

func (p *payosGateway) RefundPayment(orderCode int64, reason string) error {
	return RefundPayOSPayment(orderCode, reason)
}

func (p *payosGateway) CancelPayment(orderCode int64) error {
	return CancelPayOSPayment(orderCode)
}

func (p *payosGateway) GenerateSignature(data string) string {
	return GeneratePayOSSignature(data)
}
