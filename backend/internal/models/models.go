package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	OrderPending   uint8 = 0 // Chưa thanh toán
	OrderConfirmed uint8 = 1 // Đã xác nhận, chưa thanh toán
	OrderPaid      uint8 = 2 // Đã thanh toán
	OrderDelivered uint8 = 4 // Đã giao hàng
	OrderCancelled uint8 = 8 // Đã hủy
)

const (
	PaymentCod uint8 = iota // 0
	PaymentQR               // 1
)

// 1. USER SYSTEM (Total: ~142 bytes - optimized for 8-byte alignment)
type User struct {
	ID             uint64     `gorm:"primaryKey"`                                  // 8 bytes
	TotalSpent     uint64     `gorm:"default:0;index"`                             // 8 bytes (VND)
	CreatedAt      time.Time  `gorm:"index"`                                       // 24 bytes
	UpdatedAt      time.Time  `json:"updated_at"`                                  // 24 bytes
	LastPurchaseAt *time.Time `gorm:"index"`                                       // 8 bytes (pointer)
	Email          string     `gorm:"uniqueIndex:idx_users_email_unique;not null"` // 16 bytes
	Password       string     `gorm:"not null"`                                    // 16 bytes
	Name           string     // 16 bytes
	Phone          string     `gorm:"uniqueIndex:idx_users_phone_unique"` // 16 bytes
	TotalOrders    uint32     `gorm:"default:0"`                          // 4 bytes
	IsActive       uint8      `gorm:"default:1;index"`                    // 1 byte (bitwise: 1 = active)
}

// 2. PRODUCT SYSTEM (Simplified - no variants, no categories, use tags array) (Total: ~184 bytes - optimized for 8-byte alignment)
type Product struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`                               // 8 bytes
	Price       uint64         `gorm:"default:0;check:price >= 0" db:"price" json:"price"` // 8 bytes
	CreatedAt   time.Time      `json:"created_at"`                                         // 24 bytes
	UpdatedAt   time.Time      `json:"updated_at"`                                         // 24 bytes
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`                  // 24 bytes
	Name        string         `gorm:"index" db:"name" json:"name"`                        // 16 bytes
	Description string         `gorm:"type:text" db:"description" json:"description"`      // 16 bytes
	Thumbnail   string         `gorm:"" db:"thumbnail" json:"thumbnail"`                   // 16 bytes
	Images      datatypes.JSON `gorm:"type:jsonb" db:"images" json:"images"`               // 16 bytes
	Tags        datatypes.JSON `gorm:"type:jsonb;default:'[]'" db:"tags" json:"tags"`      // 16 bytes
	Stock       uint32         `gorm:"default:0;check:stock >= 0" db:"stock" json:"stock"` // 4 bytes
	IsActive    uint8          `gorm:"default:1" db:"is_active" json:"is_active"`          // 1 byte (bitwise: 1 = active)
	Status      uint8          `db:"status" json:"status"`                                 // 1 byte
}

// 3. ORDER SYSTEM (Total: ~90 bytes - optimized for 8-byte alignment, saved 16 bytes by removing OrderNumber)
type Order struct {
	ID              uint64         `gorm:"primaryKey"` // 8 bytes
	TotalAmount     uint64         // 8 bytes
	CreatedAt       time.Time      `gorm:"index"` // 24 bytes
	CustomerPhone   string         // 16 bytes
	ShippingAddress datatypes.JSON `gorm:"type:jsonb"`              // 16 bytes
	Items           datatypes.JSON `gorm:"type:jsonb;default:'[]'"` // 16 bytes
	PaymentMethod   uint8          `gorm:"default:0"`               // 1 byte
	Status          uint8          `gorm:"index"`                   // 1 byte (OrderStatus constants: includes payment state)
	PayOSOrderCode  *int64         `gorm:"uniqueIndex"`             // 8 bytes (nullable) - for webhook idempotency
}

// 4. LIMITED DROP - Chiến thuật thả hàng 1 đợt duy nhất, dùng Postgres lock thay vì Redis (Total: ~77 bytes - optimized for 8-byte alignment)
type LimitedDrop struct {
	ID         uint64     `gorm:"primaryKey" json:"id"`                                          // 8 bytes
	ProductID  uint64     `gorm:"index" db:"product_id" json:"product_id"`                       // 8 bytes
	StartTime  time.Time  `gorm:"index" db:"start_time" json:"starts_at"`                        // 24 bytes
	EndTime    *time.Time `gorm:"index" db:"end_time" json:"ends_at"`                            // 8 bytes (pointer)
	Name       string     `gorm:"index" db:"name" json:"name"`                                   // 16 bytes
	TotalStock uint32     `gorm:"check:total_stock >= 0" db:"total_stock" json:"total_stock"`    // 4 bytes
	DropSize   uint32     `gorm:"default:1;check:drop_size > 0" db:"drop_size" json:"drop_size"` // 4 bytes
	Sold       uint32     `gorm:"default:0;check:sold >= 0" db:"sold" json:"sold"`               // 4 bytes
	IsActive   uint8      `gorm:"default:0;index" db:"is_active" json:"is_active"`               // 1 byte (bitwise: 1 = active)
}

// 5. SYMBICODE - Anti-counterfeit system (1 symbicode per product sale) (Total: ~105 bytes - optimized for 8-byte alignment)
type Symbicode struct {
	ID          uint64     `gorm:"primaryKey"`                               // 8 bytes
	OrderID     uint64     `gorm:"index" db:"order_id"`                      // 8 bytes (Link to order for tracking)
	ProductID   uint64     `gorm:"index" db:"product_id"`                    // 8 bytes
	CreatedAt   time.Time  `gorm:"index"`                                    // 24 bytes
	ActivatedAt *time.Time `gorm:"index" db:"activated_at"`                  // 8 bytes (pointer)
	Code        []byte     `gorm:"type:uuid;uniqueIndex;not null" db:"code"` // 16 bytes (UUID v7 binary)
	SecretKey   string     `gorm:"not null" db:"secret_key"`                 // 16 bytes
	ActivatedIP string     `gorm:"index" db:"activated_ip"`                  // 16 bytes
	IsActivated uint8      `gorm:"default:0;index" db:"is_activated"`        // 1 byte (bitwise: 1 = activated)
}
