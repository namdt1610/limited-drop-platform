package service

import (
	"ecommerce-backend/internal/models"
	"time"
)

// LimitedDropStatus represents the current status of a drop for the frontend
type LimitedDropStatus struct {
	DropID      uint64     `json:"drop_id"`
	Name        string     `json:"name"`
	ProductID   uint64     `json:"product_id"`
	ProductName string     `json:"product_name"`
	Price       uint64     `json:"price"`
	TotalStock  uint32     `json:"total_stock"`
	Sold        uint32     `json:"sold"`
	Available   uint32     `json:"available"`
	DropSize    uint32     `json:"drop_size"`
	IsActive    bool       `json:"is_active"`
	StartsAt    time.Time  `json:"starts_at"`
	EndsAt      *time.Time `json:"ends_at"`
	Now         time.Time  `json:"now"`
}

// GetActiveDrops returns all active drops
func (s *service) GetActiveDrops() ([]models.LimitedDrop, error) {
	return s.repo.GetActiveDrops()
}

// GetDropStatus returns the status of a specific drop
func (s *service) GetDropStatus(id uint64) (*LimitedDropStatus, error) {
	drop, err := s.repo.GetDropByID(id)
	if err != nil {
		return nil, err
	}

	product, err := s.repo.GetProductByID(drop.ProductID)
	if err != nil {
		return nil, err
	}

	available := uint32(0)
	if drop.TotalStock > drop.Sold {
		available = drop.TotalStock - drop.Sold
	}

	return &LimitedDropStatus{
		DropID:      drop.ID,
		Name:        drop.Name,
		ProductID:   drop.ProductID,
		ProductName: product.Name,
		Price:       product.Price,
		TotalStock:  drop.TotalStock,
		Sold:        drop.Sold,
		Available:   available,
		DropSize:    drop.DropSize,
		IsActive:    drop.IsActive == 1,
		StartsAt:    drop.StartTime,
		EndsAt:      drop.EndTime,
		Now:         time.Now().UTC(),
	}, nil
}

// PurchaseDrop has been moved to drop_purchase.go
// ProcessSuccessfulDropPayment has been moved to drop_payment.go

