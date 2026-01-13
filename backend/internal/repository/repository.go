package repository

import (
	"database/sql"
	"ecommerce-backend/internal/models"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/datatypes"
)

// Repository defines the interface for business logic data operations only
// CRUD operations are handled by NocoDB - we only need business flow operations
type Repository interface {
	// Product operations for public API
	GetProductByID(id uint64) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)

	// Order operations for purchase completion and tracking
	CreateOrder(order *models.Order) error
	GetOrderByID(id uint64) (*models.Order, error)
	GetOrdersByUserPhone(phone string) ([]models.Order, error)
	GetOrderByPayOSOrderCode(orderCode int64) (*models.Order, error)
	UpdateOrderStatus(id uint64, status uint8) error

	// Drop operations for drop flow
	GetActiveDrops() ([]models.LimitedDrop, error)
	GetDropByID(id uint64) (*models.LimitedDrop, error)
	IncrementSoldCount(id uint64, increment uint32) error
	DecrementSoldCount(id uint64, decrement uint32) error

	// Transaction support
	WithTransaction(fn func(Repository) error) error

	// Symbicode operations for verification flow
	CreateSymbicode(symbicode *models.Symbicode) error
	GetSymbicodeByCode(code []byte) (*models.Symbicode, error)
	ActivateSymbicode(id uint64, ip string) error
}

// DBExecutor interface that both *sql.DB and *sql.Tx implement
type DBExecutor interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// repository implements Repository interface
type repository struct {
	db DBExecutor
}

// NewRepository creates a new repository instance
func NewRepository(db DBExecutor) Repository {
	return &repository{db: db}
}

// Helper functions for JSON marshaling/unmarshaling
func marshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func unmarshalJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Helper to convert datatypes.JSON to []byte for scanning
func jsonToBytes(j datatypes.JSON) []byte {
	return []byte(j)
}

// Helper to convert []byte to datatypes.JSON
func bytesToJSON(data []byte) datatypes.JSON {
	return datatypes.JSON(data)
}

// Helper to convert sql.NullTime to *time.Time
func nullTimeToPtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

// WithTransaction executes a function within a database transaction
func (r *repository) WithTransaction(fn func(Repository) error) error {
	// Try to use SmartExecutor's Begin method first
	if executor, ok := r.db.(interface{ Begin() (*sql.Tx, error) }); ok {
		tx, err := executor.Begin()
		if err != nil {
			return err
		}

		// Create a temporary repository with the transaction
		txRepo := &repository{db: tx}

		// Execute the function
		err = fn(txRepo)
		if err != nil {
			tx.Rollback()
			return err
		}

		return tx.Commit()
	}

	// Fallback: Type assert to *sql.DB to access Begin method
	db, ok := r.db.(*sql.DB)
	if !ok {
		return fmt.Errorf("WithTransaction can only be called on repository with *sql.DB or SmartExecutor")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Create a temporary repository with the transaction
	txRepo := &repository{db: tx}

	// Execute the function
	err = fn(txRepo)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Helper to convert *time.Time to sql.NullTime
func ptrToNullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{Valid: false}
}
