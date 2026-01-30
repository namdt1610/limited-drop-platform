package integrations_test

import (
	"testing"

	integrations "ecommerce-backend/internal/integrations"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// =============================================================================
// HEALTH CHECK TESTS (Table-Driven)
// =============================================================================

func TestCheckHealth_TableDriven(t *testing.T) {
	// Setup in-memory database for "ok" cases
	validDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	tests := []struct {
		name           string
		db             *gorm.DB
		wantStatus     string
		wantDBContains string
	}{
		{
			name:           "success - database ok",
			db:             validDB,
			wantStatus:     "ok",
			wantDBContains: "ok",
		},
		{
			name:           "degraded - nil database",
			db:             nil,
			wantStatus:     "degraded",
			wantDBContains: "error:",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			status := integrations.CheckHealth(tc.db)

			assert.Equal(t, tc.wantStatus, status.Status)
			assert.NotEmpty(t, status.Timestamp)
			assert.Contains(t, status.Checks["database"], tc.wantDBContains)
		})
	}
}

func TestIsHealthy_TableDriven(t *testing.T) {
	// Setup in-memory database for "ok" cases
	validDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	tests := []struct {
		name        string
		db          *gorm.DB
		wantHealthy bool
	}{
		{
			name:        "healthy - valid database",
			db:          validDB,
			wantHealthy: true,
		},
		{
			name:        "unhealthy - nil database",
			db:          nil,
			wantHealthy: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			healthy := integrations.IsHealthy(tc.db)
			assert.Equal(t, tc.wantHealthy, healthy)
		})
	}
}
