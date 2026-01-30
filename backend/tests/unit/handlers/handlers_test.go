package handlers_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ecommerce-backend/internal/handlers"
	"ecommerce-backend/internal/integrations"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// MOCK SERVICE - Implements service.Service interface
// =============================================================================

type mockService struct {
	// Product
	products    map[uint64]*models.Product
	productErr  error
	productsErr error

	// Drop
	drops       map[uint64]*models.LimitedDrop
	dropErr     error
	purchaseRes *service.PurchaseResult // Configurable result
	purchaseErr error               // Configurable error
	processPaymentErr error         // Configurable error

	// Order
	orders        map[uint64]*models.Order
	ordersByPhone map[string][]models.Order
	orderErr      error

	// Symbicode
	symbicode      *models.Symbicode
	symbicodeValid bool
	symbicodeErr   error
}

func newMockService() *mockService {
	return &mockService{
		products:      make(map[uint64]*models.Product),
		drops:         make(map[uint64]*models.LimitedDrop),
		orders:        make(map[uint64]*models.Order),
		ordersByPhone: make(map[string][]models.Order),
	}
}

// Product methods
func (m *mockService) GetProduct(id uint64) (*models.Product, error) {
	if m.productErr != nil {
		return nil, m.productErr
	}
	p, ok := m.products[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return p, nil
}

func (m *mockService) ListProducts() ([]models.Product, error) {
	if m.productsErr != nil {
		return nil, m.productsErr
	}
	var result []models.Product
	for _, p := range m.products {
		result = append(result, *p)
	}
	return result, nil
}

// Order methods
func (m *mockService) CreateOrder(customerPhone string, shippingAddress []byte, items []byte, paymentMethod uint8, payOSOrderCode *int64) (*models.Order, error) {
	return nil, nil
}

func (m *mockService) GetOrderByID(id uint64) (*models.Order, error) {
	if m.orderErr != nil {
		return nil, m.orderErr
	}
	o, ok := m.orders[id]
	if !ok {
		return nil, errors.New("order not found")
	}
	return o, nil
}

func (m *mockService) GetOrdersByUserPhone(phone string) ([]models.Order, error) {
	return m.ordersByPhone[phone], nil
}


// Drop methods
func (m *mockService) GetActiveDrops() ([]models.LimitedDrop, error) {
	if m.dropErr != nil {
		return nil, m.dropErr
	}
	var result []models.LimitedDrop
	for _, d := range m.drops {
		result = append(result, *d)
	}
	return result, nil
}

func (m *mockService) GetDropStatus(id uint64) (*service.LimitedDropStatus, error) {
	if m.dropErr != nil {
		return nil, m.dropErr
	}
	d, ok := m.drops[id]
	if !ok {
		return nil, errors.New("drop not found")
	}
	return &service.LimitedDropStatus{
		DropID:     d.ID,
		Name:       d.Name,
		TotalStock: d.TotalStock,
		Sold:       d.Sold,
		Available:  d.TotalStock - d.Sold,
		StartsAt:   d.StartTime,
		IsActive:   d.IsActive == 1,
	}, nil
}

func (m *mockService) PurchaseDrop(dropID uint64, req *service.PurchaseRequest) (*service.PurchaseResult, error) {
	if m.purchaseErr != nil {
		return nil, m.purchaseErr
	}
	if m.purchaseRes != nil {
		return m.purchaseRes, nil
	}
	return &service.PurchaseResult{
		Message:   "success",
		OrderCode: 12345,
	}, nil
}

func (m *mockService) ProcessSuccessfulDropPayment(orderCode int64) error {
	return m.processPaymentErr
}

// Symbicode methods
func (m *mockService) GenerateSymbicode(productID uint64, orderID *uint64) (*models.Symbicode, error) {
	return nil, nil
}

func (m *mockService) VerifySymbicode(code string) (*models.Symbicode, bool, error) {
	if m.symbicodeErr != nil {
		return nil, false, m.symbicodeErr
	}
	return m.symbicode, m.symbicodeValid, nil
}


// =============================================================================
// PRODUCT HANDLER TESTS
// =============================================================================

func TestListProducts_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*mockService)
		wantStatus int
		wantCount  int
	}{
		{
			name: "success - returns products",
			setup: func(m *mockService) {
				m.products[1] = &models.Product{ID: 1, Name: "Watch 1", Price: 100000}
				m.products[2] = &models.Product{ID: 2, Name: "Watch 2", Price: 200000}
			},
			wantStatus: 200,
			wantCount:  2,
		},
		{
			name:       "success - empty list",
			setup:      func(m *mockService) {},
			wantStatus: 200,
			wantCount:  0,
		},
		{
			name: "error - service error",
			setup: func(m *mockService) {
				m.productsErr = errors.New("database error")
			},
			wantStatus: 500,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := newMockService()
			tc.setup(mockSvc)

			app := fiber.New()
			h := handlers.NewHandlers(mockSvc)
			h.RegisterRoutes(app)

			req := httptest.NewRequest("GET", "/api/products", nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.wantStatus, resp.StatusCode)

			if tc.wantStatus == 200 {
				body, _ := io.ReadAll(resp.Body)
				var products []models.Product
				json.Unmarshal(body, &products)
				assert.Len(t, products, tc.wantCount)
			}
		})
	}
}

func TestGetProduct_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		productID  string
		setup      func(*mockService)
		wantStatus int
		wantName   string
	}{
		{
			name:      "success - product exists",
			productID: "1",
			setup: func(m *mockService) {
				m.products[1] = &models.Product{ID: 1, Name: "Test Watch", Price: 500000}
			},
			wantStatus: 200,
			wantName:   "Test Watch",
		},
		{
			name:       "error - product not found",
			productID:  "999",
			setup:      func(m *mockService) {},
			wantStatus: 404,
		},
		{
			name:       "error - invalid ID format",
			productID:  "abc",
			setup:      func(m *mockService) {},
			wantStatus: 400,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := newMockService()
			tc.setup(mockSvc)

			app := fiber.New()
			h := handlers.NewHandlers(mockSvc)
			h.RegisterRoutes(app)

			req := httptest.NewRequest("GET", "/api/products/"+tc.productID, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.wantStatus, resp.StatusCode)

			if tc.wantStatus == 200 {
				body, _ := io.ReadAll(resp.Body)
				var product models.Product
				json.Unmarshal(body, &product)
				assert.Equal(t, tc.wantName, product.Name)
			}
		})
	}
}

// =============================================================================
// HEALTH HANDLER TESTS
// =============================================================================

func TestHealthEndpoint(t *testing.T) {
	mockSvc := newMockService()
	app := fiber.New()
	h := handlers.NewHandlers(mockSvc)
	h.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	assert.Equal(t, "ok", result["status"])
}

// =============================================================================
// DROP HANDLER TESTS
// =============================================================================

func TestGetActiveDrops_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*mockService)
		wantStatus int
		wantCount  int
	}{
		{
			name: "success - returns active drops",
			setup: func(m *mockService) {
				m.drops[1] = &models.LimitedDrop{ID: 1, Name: "Drop 1", IsActive: 1, StartTime: time.Now()}
				m.drops[2] = &models.LimitedDrop{ID: 2, Name: "Drop 2", IsActive: 1, StartTime: time.Now()}
			},
			wantStatus: 200,
			wantCount:  2,
		},
		{
			name:       "success - no active drops",
			setup:      func(m *mockService) {},
			wantStatus: 200,
			wantCount:  0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := newMockService()
			tc.setup(mockSvc)

			app := fiber.New()
			h := handlers.NewHandlers(mockSvc)
			h.RegisterRoutes(app)

			req := httptest.NewRequest("GET", "/api/drops", nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.wantStatus, resp.StatusCode)
		})
	}
}

func TestGetDropStatus_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		dropID     string
		setup      func(*mockService)
		wantStatus int
	}{
		{
			name:   "success - drop exists",
			dropID: "1",
			setup: func(m *mockService) {
				m.drops[1] = &models.LimitedDrop{ID: 1, Name: "Test Drop", IsActive: 1, TotalStock: 100, Sold: 50, StartTime: time.Now()}
			},
			wantStatus: 200,
		},
		{
			name:       "error - drop not found",
			dropID:     "999",
			setup:      func(m *mockService) {},
			wantStatus: 404,
		},
		{
			name:       "error - invalid ID",
			dropID:     "abc",
			setup:      func(m *mockService) {},
			wantStatus: 400,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := newMockService()
			tc.setup(mockSvc)

			app := fiber.New()
			h := handlers.NewHandlers(mockSvc)
			h.RegisterRoutes(app)

			req := httptest.NewRequest("GET", "/api/drops/"+tc.dropID+"/status", nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.wantStatus, resp.StatusCode)
		})
	}
}

func TestPayOSWebhook_TableDriven(t *testing.T) {
	// Set API Key for signature generation (logic in integrations.GeneratePayOSSignature uses PAYOS_API_KEY)
	t.Setenv("PAYOS_API_KEY", "test-api-key")

	// Helper to generate signature
	sign := func(body string) string {
		return integrations.GeneratePayOSSignature(body)
	}

	tests := []struct {
		name           string
		body           string
		signature      string
		mockEnv        map[string]string
		setup          func(*mockService)
		wantStatus     int
		wantBody       string
	}{
		{
			name: "success - valid paid webhook",
			body: `{"code":"00","desc":"success","data":{"orderCode":123,"amount":100000,"status":"PAID"}}`,
			signature: sign(`{"code":"00","desc":"success","data":{"orderCode":123,"amount":100000,"status":"PAID"}}`),
			mockEnv: map[string]string{"PAYOS_CLIENT_ID": "client-id"},
			setup: func(m *mockService) {
				m.processPaymentErr = nil
			},
			wantStatus: 200,
			wantBody:   `{"message":"Payment processed successfully"}`,
		},
		{
			name: "ignored - payment not paid",
			body: `{"code":"00","desc":"success","data":{"orderCode":123,"status":"PENDING"}}`,
			signature: sign(`{"code":"00","desc":"success","data":{"orderCode":123,"status":"PENDING"}}`),
			mockEnv: map[string]string{"PAYOS_CLIENT_ID": "client-id"},
			setup: func(m *mockService) {},
			wantStatus: 200,
			wantBody:   `{"message":"Payment not completed"}`, 
		},
		{
			name: "error - invalid signature",
			body: `{"data":{"status":"PAID"}}`,
			signature: "invalid_sig",
			mockEnv: map[string]string{"PAYOS_CLIENT_ID": "client-id"},
			setup: func(m *mockService) {},
			wantStatus: 401,
			wantBody:   `{"error":"Invalid webhook signature"}`,
		},
		{
			name: "error - missing signature in production",
			body: `{"data":{"status":"PAID"}}`,
			signature: "",
			mockEnv: map[string]string{"PAYOS_CLIENT_ID": "client-id"},
			setup: func(m *mockService) {},
			wantStatus: 400,
			wantBody:   `{"error":"Missing webhook signature"}`,
		},
		{
			name: "success - dev mode (no signature)",
			body: `{"data":{"orderCode":123,"status":"PAID"}}`,
			signature: "",
			mockEnv: map[string]string{}, // Empty PAYOS_CLIENT_ID NOT setting it ensures it interprets as empty
			setup: func(m *mockService) {
				m.processPaymentErr = nil
			},
			wantStatus: 200,
		},
		{
			name: "success - service error (idempotency)",
			body: `{"data":{"orderCode":123,"status":"PAID"}}`,
			signature: sign(`{"data":{"orderCode":123,"status":"PAID"}}`),
			mockEnv: map[string]string{"PAYOS_CLIENT_ID": "client-id"},
			setup: func(m *mockService) {
				m.processPaymentErr = errors.New("already processed")
			},
			wantStatus: 200, // Still 200 to ack PayOS
			wantBody:   `{"message":"Payment processed with errors"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			
			// Set Env using t.Setenv for auto-cleanup
			for k, v := range tc.mockEnv {
				t.Setenv(k, v)
			}
			// Important: for the "dev mode" case, we rely on PAYOS_CLIENT_ID NOT being set.
			// t.Setenv only sets it for the rest of the test.
			// Since we start fresh in each t.Run, previous t.Setenv shouldn't leak if using t.Setenv properly?
			// Actually t.Setenv cleans up after the subtest finishes.
			
			mockSvc := newMockService()
			tc.setup(mockSvc)

			app := fiber.New()
			h := handlers.NewHandlers(mockSvc)
			h.RegisterRoutes(app)

			req := httptest.NewRequest("POST", "/api/limited-drops/webhook/payos", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			if tc.signature != "" {
				// Force lowercase header key
				req.Header["x-payos-signature"] = []string{tc.signature}
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			if resp.StatusCode != tc.wantStatus {
				body, _ := io.ReadAll(resp.Body)
				t.Logf("Response Body: %s", string(body))
			}
			assert.Equal(t, tc.wantStatus, resp.StatusCode)
			
			if tc.wantBody != "" {
				body, _ := io.ReadAll(resp.Body)
				assert.JSONEq(t, tc.wantBody, string(body))
			}
		})
	}
}

func TestPurchaseDrop_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		dropID     string
		body       string
		setup      func(*mockService)
		wantStatus int
	}{
		{
			name:   "success - valid purchase",
			dropID: "1",
			body:   `{"quantity":1,"name":"Test","phone":"0909","email":"test@test.com","address":"123","province":"HCM","district":"D1","ward":"W1"}`,
			setup: func(m *mockService) {
				m.purchaseRes = &service.PurchaseResult{PaymentURL: "http://pay", OrderCode: 123}
			},
			wantStatus: 200,
		},
		{
			name:   "error - missing required field (name)",
			dropID: "1",
			body:   `{"quantity":1,"phone":"0909","email":"test@test.com"}`,
			setup:  func(m *mockService) {},
			wantStatus: 400,
		},
		{
			name:   "error - service failed",
			dropID: "1",
			body:   `{"quantity":1,"name":"Test","phone":"0909","email":"test@test.com","address":"123","province":"HCM","district":"D1","ward":"W1"}`,
			setup: func(m *mockService) {
				m.purchaseErr = errors.New("sold out")
			},
			wantStatus: 400,
		},
		{
			name:   "error - invalid json",
			dropID: "1",
			body:   `{invalid}`,
			setup:  func(m *mockService) {},
			wantStatus: 400,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := newMockService()
			tc.setup(mockSvc)

			app := fiber.New()
			h := handlers.NewHandlers(mockSvc)
			h.RegisterRoutes(app)

			req := httptest.NewRequest("POST", "/api/drops/"+tc.dropID+"/purchase", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.wantStatus, resp.StatusCode)
		})
	}
}

// =============================================================================
// ORDER HANDLER TESTS
// =============================================================================

func TestGetOrderByID_TableDriven(t *testing.T) {
tests := []struct {
name       string
orderID    string
setup      func(*mockService)
wantStatus int
}{
{
name:    "success - order exists",
orderID: "1",
setup: func(m *mockService) {
m.orders[1] = &models.Order{ID: 1, TotalAmount: 500000, CustomerPhone: "0123456789"}
},
wantStatus: 200,
},
{
name:       "error - order not found",
orderID:    "999",
setup:      func(m *mockService) {},
wantStatus: 404,
},
{
name:       "error - invalid ID",
orderID:    "abc",
setup:      func(m *mockService) {},
wantStatus: 400,
},
}

for _, tc := range tests {
t.Run(tc.name, func(t *testing.T) {
mockSvc := newMockService()
tc.setup(mockSvc)

app := fiber.New()
h := handlers.NewHandlers(mockSvc)
h.RegisterRoutes(app)

req := httptest.NewRequest("GET", "/api/orders/"+tc.orderID, nil)
resp, err := app.Test(req)
require.NoError(t, err)
defer resp.Body.Close()

assert.Equal(t, tc.wantStatus, resp.StatusCode)
})
}
}

func TestGetOrdersByPhone_TableDriven(t *testing.T) {
tests := []struct {
name       string
phone      string
setup      func(*mockService)
wantStatus int
}{
{
name:  "success - orders found",
phone: "0123456789",
setup: func(m *mockService) {
m.ordersByPhone["0123456789"] = []models.Order{
{ID: 1, TotalAmount: 100000},
{ID: 2, TotalAmount: 200000},
}
},
wantStatus: 200,
},
{
name:       "success - no orders",
phone:      "0987654321",
setup:      func(m *mockService) {},
wantStatus: 200,
},
{
name:       "error - phone not provided",
phone:      "",
setup:      func(m *mockService) {},
wantStatus: 400,
},
}

for _, tc := range tests {
t.Run(tc.name, func(t *testing.T) {
mockSvc := newMockService()
tc.setup(mockSvc)

app := fiber.New()
h := handlers.NewHandlers(mockSvc)
h.RegisterRoutes(app)

url := "/api/orders"
if tc.phone != "" {
url += "?phone=" + tc.phone
}
req := httptest.NewRequest("GET", url, nil)
resp, err := app.Test(req)
require.NoError(t, err)
defer resp.Body.Close()

assert.Equal(t, tc.wantStatus, resp.StatusCode)
})
}
}

// =============================================================================
// SYMBICODE HANDLER TESTS
// =============================================================================

func TestVerifySymbicode_TableDriven(t *testing.T) {
tests := []struct {
name       string
body       string
setup      func(*mockService)
wantStatus int
}{
{
name: "success - valid symbicode",
body: `{"code":"abc123"}`,
setup: func(m *mockService) {
m.symbicodeValid = true
m.symbicode = &models.Symbicode{ID: 1, ProductID: 10, IsActivated: 0}
},
wantStatus: 200,
},
{
name: "error - invalid symbicode",
body: `{"code":"invalid"}`,
setup: func(m *mockService) {
m.symbicodeErr = errors.New("invalid symbicode")
},
wantStatus: 400,
},
{
name:       "error - invalid request body",
body:       "not json",
setup:      func(m *mockService) {},
wantStatus: 400,
},
}

for _, tc := range tests {
t.Run(tc.name, func(t *testing.T) {
mockSvc := newMockService()
tc.setup(mockSvc)

app := fiber.New()
h := handlers.NewHandlers(mockSvc)
h.RegisterRoutes(app)

req := httptest.NewRequest("POST", "/api/symbicode/verify", strings.NewReader(tc.body))
req.Header.Set("Content-Type", "application/json")
resp, err := app.Test(req)
require.NoError(t, err)
defer resp.Body.Close()

assert.Equal(t, tc.wantStatus, resp.StatusCode)
})
}
}
