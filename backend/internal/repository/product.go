package repository

import (
	"database/sql"
	"ecommerce-backend/internal/models"
)

// Product repository operations for public API only
func (r *repository) GetProductByID(id uint64) (*models.Product, error) {
	query := `
		SELECT id, price, created_at, updated_at, deleted_at,
			   name, description, thumbnail, images, tags, stock, is_active, status
		FROM products WHERE id = ? AND deleted_at IS NULL AND is_active = 1`

	var product models.Product
	var deletedAt sql.NullTime
	var imagesStr, tagsStr string

	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
		&deletedAt,
		&product.Name,
		&product.Description,
		&product.Thumbnail,
		&imagesStr,
		&tagsStr,
		&product.Stock,
		&product.IsActive,
		&product.Status,
	)

	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	// Parse JSON fields
	unmarshalJSON([]byte(imagesStr), &product.Images)
	unmarshalJSON([]byte(tagsStr), &product.Tags)

	return &product, nil
}

func (r *repository) GetAllProducts() ([]models.Product, error) {
	query := `
		SELECT id, price, created_at, updated_at, deleted_at,
			   name, description, thumbnail, images, tags, stock, is_active, status
		FROM products WHERE deleted_at IS NULL AND is_active = 1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var deletedAt sql.NullTime
		var imagesStr, tagsStr string

		err := rows.Scan(
			&product.ID,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt,
			&deletedAt,
			&product.Name,
			&product.Description,
			&product.Thumbnail,
			&imagesStr,
			&tagsStr,
			&product.Stock,
			&product.IsActive,
			&product.Status,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON fields
		unmarshalJSON([]byte(imagesStr), &product.Images)
		unmarshalJSON([]byte(tagsStr), &product.Tags)

		products = append(products, product)
	}

	return products, rows.Err()
}
