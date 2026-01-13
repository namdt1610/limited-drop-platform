/**
 * SEED SCRIPT
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/utils/uuid"

	googleuuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	buf := make([]byte, 4096)
	n, _ := f.Read(buf)
	data := string(buf[:n])
	for _, line := range splitLines(data) {
		line = trimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}
		if k, v, ok := parseKeyValue(line); ok {
			if os.Getenv(k) == "" {
				os.Setenv(k, v)
			}
		}
	}
}

func splitLines(s string) []string {
	var out []string
	cur := ""
	for _, r := range s {
		if r == '\n' || r == '\r' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			continue
		}
		cur += string(r)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

func trimSpace(s string) string {
	i := 0
	j := len(s) - 1
	for i <= j && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	for j >= i && (s[j] == ' ' || s[j] == '\t') {
		j--
	}
	if i > j {
		return ""
	}
	return s[i : j+1]
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func parseKeyValue(s string) (string, string, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			k := trimSpace(s[:i])
			v := trimSpace(s[i+1:])
			if len(v) >= 2 && ((v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'')) {
				v = v[1 : len(v)-1]
			}
			return k, v, k != ""
		}
	}
	return "", "", false
}

func main() {
	// Load .env
	loadDotEnv(".env")

	// Build DSN
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "./database.db"
	}

	// Connect to database
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Check if already seeded (skip check if FORCE_SEED=true)
	forceSeed := os.Getenv("FORCE_SEED") == "true"
	if !forceSeed {
		var userCount int64
		db.Model(&models.User{}).Count(&userCount)
		if userCount > 0 {
			fmt.Println("Database already seeded, skipping...")
			fmt.Println("To force re-seed, set FORCE_SEED=true")
			return
		}
	} else {
		fmt.Println("FORCE_SEED=true: Clearing existing data...")
		// Clear data in reverse order of dependencies
		db.Exec("DROP TABLE IF EXISTS symbicode")
		db.Exec("DROP TABLE IF EXISTS limited_drops")
		db.Exec("DROP TABLE IF EXISTS orders")
		db.Exec("DROP TABLE IF EXISTS users")
		db.Exec("DROP TABLE IF EXISTS products")
		fmt.Println("Cleared existing data")

		// Re-migrate tables after dropping
		if err := db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.LimitedDrop{}, &models.Symbicode{}); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		fmt.Println("Database migrated successfully")
	}

	rand.Seed(time.Now().UnixNano())
	fmt.Println("Seeding database with random data...")

	// List of cities for random address generation
	cities := []string{
		"Hà Nội", "Hồ Chí Minh", "Đà Nẵng", "Cần Thơ", "Hải Phòng",
		"Biên Hòa", "Nha Trang", "Vũng Tàu", "Quy Nhơn", "Huế",
	}

	// 1. Create 50 demo users (regular users only)
	fmt.Println("Creating 50 demo users...")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	users := make([]models.User, 50)
	for i := 0; i < 50; i++ {
		users[i] = models.User{
			Email:       fmt.Sprintf("user%d@example.com", i+1),
			Password:    string(hashedPassword),
			Name:        fmt.Sprintf("User %d", i+1),
			Phone:       fmt.Sprintf("090%08d", i+1),
			TotalSpent:  0, // uint64
			TotalOrders: 0, // uint32
			IsActive:    1, // uint8 (1 = active)
		}
	}
	if err := db.Create(&users).Error; err != nil {
		log.Fatalf("Failed to create users: %v", err)
	}
	fmt.Printf("Created %d demo users\n", len(users))

	// 3. Create 1 Test Product (10k for payment testing)
	fmt.Println("Creating 1 test product...")
	imagesJSON, _ := json.Marshal([]string{"/placeholder-product.svg"})
	testTagsJSON, _ := json.Marshal([]string{"test", "payment"})
	testProduct := models.Product{
		Name:        "Test Payment 10k",
		Description: "Sản phẩm test để kiểm tra thanh toán với giá 10.000 VND",
		Images:      datatypes.JSON(imagesJSON),
		Thumbnail:   "/placeholder-product.svg",
		Tags:        datatypes.JSON(testTagsJSON),
		Stock:       100,   // uint32
		Status:      1,     // uint8 (1 = active)
		Price:       10000, // uint64 (10k VND for payment testing)
		IsActive:    1,     // uint16 (1 = active)
	}
	if err := db.Create(&testProduct).Error; err != nil {
		log.Fatalf("Failed to create test product 10k: %v", err)
	}
	fmt.Println("Created test product 10k (for payment testing)")

	// 5. Create 50 Orders (linked to users) - with items stored as JSONB
	fmt.Println("Creating 50 orders...")
	orders := make([]models.Order, 50)
	orderStatuses := []uint8{
		models.OrderPending,
		models.OrderConfirmed,
		models.OrderPaid,
		models.OrderDelivered,
		models.OrderCancelled,
	}
	for i := 0; i < 50; i++ {
		user := users[rand.Intn(len(users))]
		status := orderStatuses[rand.Intn(len(orderStatuses))]
		shippingAddr, _ := json.Marshal(map[string]string{
			"name":    user.Name,
			"phone":   user.Phone,
			"address": fmt.Sprintf("Số %d, Đường ABC, %s", rand.Intn(999)+1, cities[rand.Intn(len(cities))]),
		})

		// Create items as JSONB - 1-3 items per order
		itemCount := rand.Intn(3) + 1
		items := make([]map[string]interface{}, itemCount)
		totalAmount := uint64(0)
		for j := 0; j < itemCount; j++ {
			quantity := rand.Intn(5) + 1
			price := uint64(10000) // Test product price (10k VND)
			items[j] = map[string]interface{}{
				"product_id":   testProduct.ID,
				"product_name": testProduct.Name,
				"price":        price,
				"quantity":     quantity,
				"subtotal":     price * uint64(quantity),
			}
			totalAmount += price * uint64(quantity)
		}
		itemsJSON, _ := json.Marshal(items)

		orders[i] = models.Order{
			CustomerPhone:   user.Phone, // Store customer phone directly
			ShippingAddress: datatypes.JSON(shippingAddr),
			Items:           datatypes.JSON(itemsJSON),
			TotalAmount:     totalAmount,       // uint64 VND
			PaymentMethod:   models.PaymentCod, // uint8 constant
			Status:          status,            // uint8 OrderStatus
		}
	}
	if err := db.Create(&orders).Error; err != nil {
		log.Fatalf("Failed to create orders: %v", err)
	}
	fmt.Printf("Created %d orders\n", len(orders))

	// 7. Create 1 Limited Drop (linked to the test product)
	fmt.Println("Creating 1 limited drop...")
	now := time.Now().UTC()
	limitedDrops := []models.LimitedDrop{
		{
			Name:       "Drop Test Product 10k - Single Unit",
			ProductID:  testProduct.ID,
			TotalStock: 1,                               // Single unit only
			DropSize:   1,                               // Max 1 per purchase
			StartTime:  now,                             // Starts right now
			EndTime:    timePtr(now.Add(1 * time.Hour)), // Ends in 1 hour
			Sold:       0,                               // Not sold yet
			IsActive:   1,                               // uint8: 1 = active
		},
	}
	if err := db.Create(&limitedDrops).Error; err != nil {
		log.Fatalf("Failed to create limited drops: %v", err)
	}
	fmt.Printf("Created %d limited drops\n", len(limitedDrops))

	// 8. Create 1 Symbicode per product for testing anti-counterfeit system
	fmt.Println("Creating 1 symbicode record...")

	// Generate real UUID v7 binary using the utils function
	testUUID, err := uuid.GenerateUUIDv7()
	if err != nil {
		log.Fatalf("Failed to generate UUID v7: %v", err)
	}

	symbicodeRecords := []models.Symbicode{
		{
			Code:        testUUID, // UUID v7 binary
			SecretKey:   googleuuid.New().String(),
			ProductID:   testProduct.ID,
			IsActivated: 0, // uint8: not activated
		},
	}
	if err := db.Create(&symbicodeRecords).Error; err != nil {
		log.Fatalf("Failed to create symbicode records: %v", err)
	}
	fmt.Printf("Created %d symbicode records\n", len(symbicodeRecords))

	fmt.Println("Database seeded successfully!")
	fmt.Println("")
	fmt.Println("Summary:")
	fmt.Printf("  - Users: 50 (demo users only)\n")
	fmt.Printf("  - Products: 1 (test product 10k VND)\n")
	fmt.Printf("  - Orders: 50 (with items stored as JSONB, uint64 prices, Base32 order numbers)\n")
	fmt.Printf("  - Limited Drops: %d (Postgres pessimistic locking)\n", len(limitedDrops))
	fmt.Printf("  - Symbicodes: %d (UUID v7 anti-counterfeit system)\n", len(symbicodeRecords))
}
