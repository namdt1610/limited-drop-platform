package repository

import (
	"database/sql"
	"time"
	"ecommerce-backend/internal/models"
)

// Symbicode repository operations for verification flow only
func (r *repository) CreateSymbicode(symbicode *models.Symbicode) error {
	query := `
		INSERT INTO symbicodes (
			order_id, product_id, created_at, activated_at, code, secret_key, activated_ip, is_activated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		symbicode.OrderID,
		symbicode.ProductID,
		symbicode.CreatedAt,
		ptrToNullTime(symbicode.ActivatedAt),
		symbicode.Code,
		symbicode.SecretKey,
		symbicode.ActivatedIP,
		symbicode.IsActivated,
	)
	return err
}

func (r *repository) GetSymbicodeByCode(code []byte) (*models.Symbicode, error) {
	query := `
		SELECT id, order_id, product_id, created_at, activated_at, code, secret_key, activated_ip, is_activated
		FROM symbicodes WHERE code = ?`

	var symbicode models.Symbicode
	var activatedAt sql.NullTime

	err := r.db.QueryRow(query, code).Scan(
		&symbicode.ID,
		&symbicode.OrderID,
		&symbicode.ProductID,
		&symbicode.CreatedAt,
		&activatedAt,
		&symbicode.Code,
		&symbicode.SecretKey,
		&symbicode.ActivatedIP,
		&symbicode.IsActivated,
	)

	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	symbicode.ActivatedAt = nullTimeToPtr(activatedAt)
	return &symbicode, nil
}

func (r *repository) ActivateSymbicode(id uint64, ip string) error {
	query := `UPDATE symbicodes SET is_activated = 1, activated_at = ?, activated_ip = ? WHERE id = ?`
	_, err := r.db.Exec(query, time.Now(), ip, id)
	return err
}
