package service_test

import (
	"errors"
	"time"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"
)

// =============================================================================
// MOCK REPOSITORY
// =============================================================================

// mockRepository is a configurable mock for all repository operations
type mockRepository struct {
	// Products
	products       map[uint64]*models.Product
	productErr     error
	allProductsErr error

	// Orders
	orders         map[uint64]*models.Order
	ordersByPhone  map[string][]models.Order
	orderByPayOS   map[int64]*models.Order
	createOrderErr error
	getOrderErr    error
	getOrdersErr   error

	// Drops
	drops             map[uint64]*models.LimitedDrop
	activeDrops       []models.LimitedDrop
	getDropErr        error
	getActiveDropsErr error
	incrementErr      error
	decrementErr      error
	allowIncrement    bool
	allowDecrement    bool

	// Symbicode
	symbicodes     map[string]*models.Symbicode // key = hex of code
	createSymErr   error
	getSymErr      error
	activateSymErr error

	// Transaction
	txErr error
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		products:       make(map[uint64]*models.Product),
		orders:         make(map[uint64]*models.Order),
		ordersByPhone:  make(map[string][]models.Order),
		orderByPayOS:   make(map[int64]*models.Order),
		drops:          make(map[uint64]*models.LimitedDrop),
		activeDrops:    []models.LimitedDrop{},
		symbicodes:     make(map[string]*models.Symbicode),
		allowIncrement: true,
		allowDecrement: true,
	}
}

// Product operations
func (m *mockRepository) GetProductByID(id uint64) (*models.Product, error) {
	if m.productErr != nil {
		return nil, m.productErr
	}
	if p, ok := m.products[id]; ok {
		return p, nil
	}
	return nil, errors.New("product not found")
}

func (m *mockRepository) GetAllProducts() ([]models.Product, error) {
	if m.allProductsErr != nil {
		return nil, m.allProductsErr
	}
	var products []models.Product
	for _, p := range m.products {
		products = append(products, *p)
	}
	return products, nil
}

// Order operations
func (m *mockRepository) CreateOrder(order *models.Order) error {
	if m.createOrderErr != nil {
		return m.createOrderErr
	}
	order.ID = uint64(len(m.orders) + 1)
	m.orders[order.ID] = order
	return nil
}

func (m *mockRepository) GetOrderByID(id uint64) (*models.Order, error) {
	if m.getOrderErr != nil {
		return nil, m.getOrderErr
	}
	if o, ok := m.orders[id]; ok {
		return o, nil
	}
	return nil, errors.New("order not found")
}

func (m *mockRepository) GetOrdersByUserPhone(phone string) ([]models.Order, error) {
	if m.getOrdersErr != nil {
		return nil, m.getOrdersErr
	}
	return m.ordersByPhone[phone], nil
}

func (m *mockRepository) GetOrderByPayOSOrderCode(orderCode int64) (*models.Order, error) {
	if o, ok := m.orderByPayOS[orderCode]; ok {
		return o, nil
	}
	return nil, errors.New("order not found")
}

// Drop operations
func (m *mockRepository) GetActiveDrops() ([]models.LimitedDrop, error) {
	if m.getActiveDropsErr != nil {
		return nil, m.getActiveDropsErr
	}
	return m.activeDrops, nil
}

func (m *mockRepository) GetDropByID(id uint64) (*models.LimitedDrop, error) {
	if m.getDropErr != nil {
		return nil, m.getDropErr
	}
	if d, ok := m.drops[id]; ok {
		return d, nil
	}
	return nil, errors.New("drop not found")
}

func (m *mockRepository) IncrementSoldCount(id uint64, increment uint32) error {
	if m.incrementErr != nil {
		return m.incrementErr
	}
	if !m.allowIncrement {
		return repository.ErrSoldOut
	}
	if d, ok := m.drops[id]; ok {
		if d.Sold+increment > d.TotalStock {
			return repository.ErrSoldOut
		}
		d.Sold += increment
		return nil
	}
	return errors.New("drop not found")
}

func (m *mockRepository) DecrementSoldCount(id uint64, decrement uint32) error {
	if m.decrementErr != nil {
		return m.decrementErr
	}
	if !m.allowDecrement {
		return errors.New("cannot decrement")
	}
	if d, ok := m.drops[id]; ok {
		if d.Sold >= decrement {
			d.Sold -= decrement
			return nil
		}
		return errors.New("insufficient stock")
	}
	return errors.New("drop not found")
}

// Transaction support
func (m *mockRepository) WithTransaction(fn func(repository.Repository) error) error {
	if m.txErr != nil {
		return m.txErr
	}
	return fn(m)
}

// Symbicode operations
func (m *mockRepository) CreateSymbicode(symbicode *models.Symbicode) error {
	if m.createSymErr != nil {
		return m.createSymErr
	}
	symbicode.ID = uint64(len(m.symbicodes) + 1)
	key := string(symbicode.Code)
	m.symbicodes[key] = symbicode
	return nil
}

func (m *mockRepository) GetSymbicodeByCode(code []byte) (*models.Symbicode, error) {
	if m.getSymErr != nil {
		return nil, m.getSymErr
	}
	key := string(code)
	if s, ok := m.symbicodes[key]; ok {
		return s, nil
	}
	return nil, errors.New("symbicode not found")
}

func (m *mockRepository) ActivateSymbicode(id uint64, ip string) error {
	if m.activateSymErr != nil {
		return m.activateSymErr
	}
	for _, s := range m.symbicodes {
		if s.ID == id {
			s.IsActivated = 1
			now := time.Now()
			s.ActivatedAt = &now
			s.ActivatedIP = ip
			return nil
		}
	}
	return errors.New("symbicode not found")
}

func (m *mockRepository) UpdateOrderStatus(id uint64, status uint8) error {
	if order, ok := m.orders[id]; ok {
		order.Status = status
		return nil
	}
	return errors.New("order not found")
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func ptrTime(t time.Time) *time.Time { return &t }
