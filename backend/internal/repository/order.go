package repository

import (
	"database/sql"
	"ecommerce-backend/internal/models"
)

// Order repository operations for purchase flow and order tracking
func (r *repository) CreateOrder(order *models.Order) error {
	query := `
		INSERT INTO orders (
			total_amount, created_at, customer_phone, shipping_address, items, payment_method, status, pay_os_order_code
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	// Don't marshal existing JSON types, just convert to string
	// shippingAddrJSON, _ := marshalJSON(order.ShippingAddress)
	// itemsJSON, _ := marshalJSON(order.Items)

	// Convert datatypes.JSON to string directly
	shippingAddrStr := string(order.ShippingAddress)
	itemsStr := string(order.Items)

	var payosOrderCode interface{}
	if order.PayOSOrderCode != nil {
		payosOrderCode = *order.PayOSOrderCode
	} else {
		payosOrderCode = nil
	}

	result, err := r.db.Exec(query,
		order.TotalAmount,
		order.CreatedAt,
		order.CustomerPhone,
		shippingAddrStr,
		itemsStr,
		order.PaymentMethod,
		order.Status,
		payosOrderCode,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	order.ID = uint64(id)
	return nil
}

func (r *repository) GetOrderByID(id uint64) (*models.Order, error) {
	query := `
		SELECT id, total_amount, created_at, customer_phone, shipping_address, items, payment_method, status, pay_os_order_code
		FROM orders WHERE id = ?`

	var order models.Order
	var shippingAddrStr, itemsStr string
	var payosOrderCode sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.TotalAmount,
		&order.CreatedAt,
		&order.CustomerPhone,
		&shippingAddrStr,
		&itemsStr,
		&order.PaymentMethod,
		&order.Status,
		&payosOrderCode,
	)

	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	// Handle nullable PayOSOrderCode
	if payosOrderCode.Valid {
		order.PayOSOrderCode = &payosOrderCode.Int64
	}

	// Parse JSON fields
	unmarshalJSON([]byte(shippingAddrStr), &order.ShippingAddress)
	unmarshalJSON([]byte(itemsStr), &order.Items)

	return &order, nil
}

func (r *repository) GetOrdersByUserPhone(phone string) ([]models.Order, error) {
	query := `
		SELECT id, total_amount, created_at, customer_phone, shipping_address, items, payment_method, status, pay_os_order_code
		FROM orders WHERE customer_phone = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, phone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var shippingAddrStr, itemsStr string
		var payosOrderCode sql.NullInt64

		err := rows.Scan(
			&order.ID,
			&order.TotalAmount,
			&order.CreatedAt,
			&order.CustomerPhone,
			&shippingAddrStr,
			&itemsStr,
			&order.PaymentMethod,
			&order.Status,
			&payosOrderCode,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable PayOSOrderCode
		if payosOrderCode.Valid {
			order.PayOSOrderCode = &payosOrderCode.Int64
		}

		// Parse JSON fields
		unmarshalJSON([]byte(shippingAddrStr), &order.ShippingAddress)
		unmarshalJSON([]byte(itemsStr), &order.Items)

		orders = append(orders, order)
	}

	return orders, rows.Err()
}

func (r *repository) GetOrderByPayOSOrderCode(orderCode int64) (*models.Order, error) {
	query := `
		SELECT id, total_amount, created_at, customer_phone, shipping_address, items, payment_method, status, pay_os_order_code
		FROM orders WHERE pay_os_order_code = ?`

	var order models.Order
	var shippingAddrStr, itemsStr string
	var payosOrderCode sql.NullInt64

	err := r.db.QueryRow(query, orderCode).Scan(
		&order.ID,
		&order.TotalAmount,
		&order.CreatedAt,
		&order.CustomerPhone,
		&shippingAddrStr,
		&itemsStr,
		&order.PaymentMethod,
		&order.Status,
		&payosOrderCode,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Return nil instead of error for not found
	}
	if err != nil {
		return nil, err
	}

	// Set PayOSOrderCode since this method is for PayOS orders
	if payosOrderCode.Valid {
		order.PayOSOrderCode = &payosOrderCode.Int64
	}

	// Parse JSON fields
	unmarshalJSON([]byte(shippingAddrStr), &order.ShippingAddress)
	unmarshalJSON([]byte(itemsStr), &order.Items)

	return &order, nil
}

func (r *repository) UpdateOrderStatus(id uint64, status uint8) error {
	query := `UPDATE orders SET status = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, id)
	return err
}
