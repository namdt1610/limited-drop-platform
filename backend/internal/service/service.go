package service

import (
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"
)

// Service defines the interface for business logic operations
type Service interface {
	// Product services
	GetProduct(id uint64) (*models.Product, error)
	ListProducts() ([]models.Product, error)

	// Order services
	CreateOrder(customerPhone string, shippingAddress []byte, items []byte, paymentMethod uint8, payOSOrderCode *int64) (*models.Order, error)
	GetOrderByID(id uint64) (*models.Order, error)
	GetOrdersByUserPhone(phone string) ([]models.Order, error)

	// Drop services
	GetActiveDrops() ([]models.LimitedDrop, error)
	GetDropStatus(id uint64) (*LimitedDropStatus, error)
	PurchaseDrop(dropID uint64, req *PurchaseRequest) (*PurchaseResult, error)
	ProcessSuccessfulDropPayment(orderCode int64) error

	// Symbicode services
	GenerateSymbicode(productID uint64, orderID *uint64) (*models.Symbicode, error)
	VerifySymbicode(code string) (*models.Symbicode, bool, error)
}

// PurchaseRequest represents a limited drop purchase request
type PurchaseRequest struct {
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Province string `json:"province"`
	District string `json:"district"`
	Ward     string `json:"ward"`
}

// PurchaseResult represents the result of a purchase attempt
type PurchaseResult struct {
	Message    string `json:"message"`
	PaymentURL string `json:"payment_url,omitempty"`
	OrderCode  int64  `json:"order_code,omitempty"`
}

// service implements Service interface
type service struct {
	repo repository.Repository
}

// NewService creates a new service instance
func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}
