package integrations_test

import (
	"testing"

	integrations "ecommerce-backend/internal/integrations"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCheckHealth_DatabaseOK(t *testing.T) {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Test
	status := integrations.CheckHealth(db)

	// Assertions
	assert.Equal(t, "ok", status.Status)
	assert.NotEmpty(t, status.Timestamp)
	assert.Equal(t, "ok", status.Checks["database"])
}

func TestCheckHealth_DatabaseError(t *testing.T) {
	// Setup invalid database connection - use a different approach
	// Since SQLite in-memory always succeeds, we'll test with a nil db
	var db *gorm.DB // nil database

	// Test
	status := integrations.CheckHealth(db)

	// Assertions
	assert.Equal(t, "degraded", status.Status)
	assert.NotEmpty(t, status.Timestamp)
	assert.Contains(t, status.Checks["database"], "error:")
}

func TestIsHealthy_OK(t *testing.T) {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Test
	healthy := integrations.IsHealthy(db)

	// Assertions
	assert.True(t, healthy)
}

func TestIsHealthy_Degraded(t *testing.T) {
	// Test with nil database
	var db *gorm.DB // nil database

	// Test
	healthy := integrations.IsHealthy(db)

	// Assertions
	assert.False(t, healthy)
}
