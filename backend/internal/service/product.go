package service

import (
	"ecommerce-backend/internal/models"
)

// GetProduct retrieves a product by ID
func (s *service) GetProduct(id uint64) (*models.Product, error) {
	return s.repo.GetProductByID(id)
}

// ListProducts retrieves all active products
func (s *service) ListProducts() ([]models.Product, error) {
	products, err := s.repo.GetAllProducts()
	if err != nil {
		return nil, err
	}

	// Filter only active products
	var activeProducts []models.Product
	for _, product := range products {
		if product.IsActive == 1 {
			activeProducts = append(activeProducts, product)
		}
	}

	return activeProducts, nil
}
