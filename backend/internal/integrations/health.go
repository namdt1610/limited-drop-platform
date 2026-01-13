package integrations

import (
	"context"
	"fmt"
	"log"
	"time"

	"ecommerce-backend/internal/models"

	"gorm.io/gorm"
)

// HealthStatus represents the health check result
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

// CheckHealth performs comprehensive health checks
func CheckHealth(db *gorm.DB) HealthStatus {
	status := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    make(map[string]string),
	}

	// Check database connection
	if db == nil {
		status.Status = "degraded"
		status.Checks["database"] = "error: database connection is nil"
	} else {
		sqlDB, err := db.DB()
		if err != nil {
			status.Status = "degraded"
			status.Checks["database"] = fmt.Sprintf("error: %v", err)
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := sqlDB.PingContext(ctx); err != nil {
				status.Status = "degraded"
				status.Checks["database"] = fmt.Sprintf("ping failed: %v", err)
			} else {
				status.Checks["database"] = "ok"
			}
		}
	}

	// Auto-activate expired symbicode (run during health checks)
	if db != nil {
		if err := AutoActivateExpiredSymbicode(db); err != nil {
			// Log error but don't fail health check
			status.Checks["symbicode_auto_activate"] = fmt.Sprintf("warning: %v", err)
		} else {
			status.Checks["symbicode_auto_activate"] = "ok"
		}
	}

	return status
}

// AutoActivateExpiredSymbicode activates symbicode that are older than 3 days and not yet activated
func AutoActivateExpiredSymbicode(db *gorm.DB) error {
	threeDaysAgo := time.Now().AddDate(0, 0, -3)

	result := db.Model(&models.Symbicode{}).
		Where("is_activated = ? AND created_at < ?", 0, threeDaysAgo).
		Updates(map[string]interface{}{
			"is_activated": true,
			"activated_at": time.Now(),
			"activated_ip": "AUTO_ACTIVATED",
		})

	if result.Error != nil {
		return fmt.Errorf("failed to auto-activate expired symbicode: %w", result.Error)
	}

	log.Printf("[symbicode] auto-activated %d expired symbicode", result.RowsAffected)
	return nil
}

// IsHealthy returns true if all critical services are healthy
func IsHealthy(db *gorm.DB) bool {
	status := CheckHealth(db)
	return status.Status == "ok"
}
