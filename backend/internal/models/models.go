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
	CreatedAt      time.Time  `gorm:"index"`
	UpdatedAt      time.Time  `json:"updated_at"`
	LastPurchaseAt *time.Time `gorm:"index"`
	Email          string     `gorm:"uniqueIndex:idx_users_email_unique;not null"`
	Password       string     `gorm:"not null"`
	Name           string
	Phone          string `gorm:"uniqueIndex:idx_users_phone_unique"`
	ID             uint64 `gorm:"primaryKey"`
	TotalSpent     uint64 `gorm:"default:0;index"`
	TotalOrders    uint32 `gorm:"default:0"`
	IsActive       uint8  `gorm:"default:1;index"`
}

// 2. PRODUCT SYSTEM (Simplified - no variants, no categories, use tags array) (Total: ~184 bytes - optimized for 8-byte alignment)
type Product struct {
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Description string         `gorm:"type:text" db:"description" json:"description"`
	Name        string         `gorm:"index" db:"name" json:"name"`
	Thumbnail   string         `gorm:"" db:"thumbnail" json:"thumbnail"`
	Images      datatypes.JSON `gorm:"type:jsonb" db:"images" json:"images"`
	Tags        datatypes.JSON `gorm:"type:jsonb;default:'[]'" db:"tags" json:"tags"`
	Price       uint64         `gorm:"default:0;check:price >= 0" db:"price" json:"price"`
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Stock       uint32         `gorm:"default:0;check:stock >= 0" db:"stock" json:"stock"`
	IsActive    uint8          `gorm:"default:1" db:"is_active" json:"is_active"`
	Status      uint8          `db:"status" json:"status"`
}

// 3. ORDER SYSTEM (Total: ~90 bytes - optimized for 8-byte alignment, saved 16 bytes by removing OrderNumber)
type Order struct {
	CreatedAt       time.Time `gorm:"index"`
	PayOSOrderCode  *int64    `gorm:"uniqueIndex"`
	CustomerPhone   string
	ShippingAddress datatypes.JSON `gorm:"type:jsonb"`
	Items           datatypes.JSON `gorm:"type:jsonb;default:'[]'"`
	ID              uint64         `gorm:"primaryKey"`
	TotalAmount     uint64
	PaymentMethod   uint8 `gorm:"default:0"`
	Status          uint8 `gorm:"index"`
}

// 4. LIMITED DROP - Chiến thuật thả hàng 1 đợt duy nhất, dùng Postgres lock thay vì Redis (Total: ~77 bytes - optimized for 8-byte alignment)
type LimitedDrop struct {
	StartTime  time.Time  `gorm:"index" db:"start_time" json:"starts_at"`
	EndTime    *time.Time `gorm:"index" db:"end_time" json:"ends_at"`
	Name       string     `gorm:"index" db:"name" json:"name"`
	ID         uint64     `gorm:"primaryKey" json:"id"`
	ProductID  uint64     `gorm:"index" db:"product_id" json:"product_id"`
	TotalStock uint32     `gorm:"check:total_stock >= 0" db:"total_stock" json:"total_stock"`
	DropSize   uint32     `gorm:"default:1;check:drop_size > 0" db:"drop_size" json:"drop_size"`
	Sold       uint32     `gorm:"default:0;check:sold >= 0" db:"sold" json:"sold"`
	IsActive   uint8      `gorm:"default:0;index" db:"is_active" json:"is_active"`
}

// 5. SYMBICODE - Anti-counterfeit system (1 symbicode per product sale) (Total: ~105 bytes - optimized for 8-byte alignment)
type Symbicode struct {
	CreatedAt   time.Time  `gorm:"index"`
	ActivatedAt *time.Time `gorm:"index" db:"activated_at"`
	SecretKey   string     `gorm:"not null" db:"secret_key"`
	ActivatedIP string     `gorm:"index" db:"activated_ip"`
	Code        []byte     `gorm:"type:uuid;uniqueIndex;not null" db:"code"`
	ID          uint64     `gorm:"primaryKey"`
	OrderID     uint64     `gorm:"index" db:"order_id"`
	ProductID   uint64     `gorm:"index" db:"product_id"`
	IsActivated uint8      `gorm:"default:0;index" db:"is_activated"`
}
