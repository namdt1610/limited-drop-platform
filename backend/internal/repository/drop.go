package repository

import (
	"database/sql"
	"ecommerce-backend/internal/models"
	"errors"
)

// ErrSoldOut is returned when the conditional update fails because stock is depleted
var ErrSoldOut = errors.New("limited drop is sold out")

// Drop repository operations for drop flow only
func (r *repository) GetDropByID(id uint64) (*models.LimitedDrop, error) {
	query := `
		SELECT id, product_id, start_time, end_time, name, total_stock, drop_size, sold, is_active
		FROM limited_drops WHERE id = ? AND is_active = 1`

	var drop models.LimitedDrop
	var endTime sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&drop.ID,
		&drop.ProductID,
		&drop.StartTime,
		&endTime,
		&drop.Name,
		&drop.TotalStock,
		&drop.DropSize,
		&drop.Sold,
		&drop.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	drop.EndTime = nullTimeToPtr(endTime)
	return &drop, nil
}

func (r *repository) GetActiveDrops() ([]models.LimitedDrop, error) {
	query := `
		SELECT id, product_id, start_time, end_time, name, total_stock, drop_size, sold, is_active
		FROM limited_drops WHERE is_active = 1
		ORDER BY start_time ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drops []models.LimitedDrop
	for rows.Next() {
		var drop models.LimitedDrop
		var endTime sql.NullTime

		err := rows.Scan(
			&drop.ID,
			&drop.ProductID,
			&drop.StartTime,
			&endTime,
			&drop.Name,
			&drop.TotalStock,
			&drop.DropSize,
			&drop.Sold,
			&drop.IsActive,
		)
		if err != nil {
			return nil, err
		}

		drop.EndTime = nullTimeToPtr(endTime)
		drops = append(drops, drop)
	}

	return drops, rows.Err()
}

func (r *repository) IncrementSoldCount(id uint64, increment uint32) error {
	// Atomic conditional update: only increment when resulting sold <= total_stock
	query := `UPDATE limited_drops SET sold = sold + ? WHERE id = ? AND is_active = 1 AND sold + ? <= total_stock`
	res, err := r.db.Exec(query, increment, id, increment)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrSoldOut
	}
	return nil
}

func (r *repository) DecrementSoldCount(id uint64, decrement uint32) error {
	// Atomic decrement: ensure sold doesn't go below 0
	query := `UPDATE limited_drops SET sold = sold - ? WHERE id = ? AND is_active = 1 AND sold >= ?`
	res, err := r.db.Exec(query, decrement, id, decrement)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("cannot decrement sold count")
	}
	return nil
}
