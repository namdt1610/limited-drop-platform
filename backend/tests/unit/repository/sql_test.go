package repository_test

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	order := &models.Order{
		TotalAmount:   100000,
		CreatedAt:     time.Now(),
		CustomerPhone: "0909123456",
		ShippingAddress: []byte(`{"city":"HCM"}`),
		Items:           []byte(`[{"id":1}]`),
		PaymentMethod: 1, // models.PaymentQR
		Status:        1,
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO orders")).
		WithArgs(
			order.TotalAmount,
			order.CreatedAt,
			order.CustomerPhone,
			string(order.ShippingAddress),
			string(order.Items),
			order.PaymentMethod,
			order.Status,
			nil, // PayOSCode defaults to nil
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateOrder(order)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), order.ID)
}

func TestGetOrderByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "total_amount", "created_at", "customer_phone", "shipping_address", "items", "payment_method", "status", "pay_os_order_code"}).
		AddRow(1, 100000, time.Now(), "0909123456", `{"city":"HCM"}`, `[{"id":1}]`, 1, 1, nil) // PaymentMethod=1

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, total_amount")).
		WithArgs(1).
		WillReturnRows(rows)

	order, err := repo.GetOrderByID(1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), order.ID)
	assert.Equal(t, uint8(1), order.PaymentMethod)
	assert.Nil(t, order.PayOSOrderCode)
}

func TestGetOrdersByUserPhone(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "total_amount", "created_at", "customer_phone", "shipping_address", "items", "payment_method", "status", "pay_os_order_code"}).
		AddRow(1, 100000, time.Now(), "0909123456", `{}`, `[]`, 1, 1, 12345).
		AddRow(2, 200000, time.Now(), "0909123456", `{}`, `[]`, 1, 1, 67890)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, total_amount")).
		WithArgs("0909123456").
		WillReturnRows(rows)

	orders, err := repo.GetOrdersByUserPhone("0909123456")
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
	assert.Equal(t, int64(12345), *orders[0].PayOSOrderCode)
}

func TestGetActiveDrops(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "product_id", "start_time", "end_time", "name", "total_stock", "drop_size", "sold", "is_active"}).
		AddRow(1, 101, now, now.Add(time.Hour), "Drop 1", 100, 1, 0, 1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT") + " .* " + regexp.QuoteMeta("FROM limited_drops")).
		WillReturnRows(rows)

	drops, err := repo.GetActiveDrops()
	assert.NoError(t, err)
	if assert.Len(t, drops, 1) {
		assert.Equal(t, "Drop 1", drops[0].Name)
	}
}

func TestIncrementSoldCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE limited_drops")).
		WithArgs(5, 1, 5). // increment, id, increment
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.IncrementSoldCount(1, 5)
	assert.NoError(t, err)
}

func TestWithTransaction_Commit(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectCommit()

	err = repo.WithTransaction(func(r repository.Repository) error {
		return nil
	})
	assert.NoError(t, err)
}

func TestWithTransaction_Rollback(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectRollback()

	err = repo.WithTransaction(func(r repository.Repository) error {
		return sql.ErrTxDone 
	})
	assert.Error(t, err)
}

func TestGetOrderByPayOSOrderCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "total_amount", "created_at", "customer_phone", "shipping_address", "items", "payment_method", "status", "pay_os_order_code"}).
		AddRow(1, 100000, time.Now(), "0909123456", `{}`, `[]`, 1, 1, 12345)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, total_amount")).
		WithArgs(int64(12345)).
		WillReturnRows(rows)

	order, err := repo.GetOrderByPayOSOrderCode(12345)
	assert.NoError(t, err)
	assert.Equal(t, int64(12345), *order.PayOSOrderCode)
}

func TestUpdateOrderStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE orders SET status = ? WHERE id = ?")).
		WithArgs(uint8(2), uint64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateOrderStatus(1, 2)
	assert.NoError(t, err)
}

func TestDecrementSoldCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE limited_drops")).
		WithArgs(1, 1, 1). // decrement, id, decrement
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DecrementSoldCount(1, 1)
	assert.NoError(t, err)
}
