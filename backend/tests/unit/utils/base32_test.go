package utils_test

import (
	"testing"

	"ecommerce-backend/internal/utils/base32"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// BASE32 ORDER NUMBER TESTS
// =============================================================================

func TestGenerateOrderNumber_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		id         uint64
		wantPrefix string
	}{
		{
			name:       "simple ID",
			id:         1,
			wantPrefix: "DV-",
		},
		{
			name:       "large ID",
			id:         12345,
			wantPrefix: "DV-",
		},
		{
			name:       "very large ID",
			id:         9999999999,
			wantPrefix: "DV-",
		},
		{
			name:       "zero ID",
			id:         0,
			wantPrefix: "DV-",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := base32.GenerateOrderNumber(tc.id)

			assert.True(t, len(result) > 3)
			assert.Equal(t, tc.wantPrefix, result[:3])
		})
	}
}

func TestDecodeOrderNumber_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		orderNumber string
		wantID      uint64
		wantErr     bool
	}{
		{
			name:        "valid - roundtrip ID 1",
			orderNumber: base32.GenerateOrderNumber(1),
			wantID:      1,
			wantErr:     false,
		},
		{
			name:        "valid - roundtrip ID 12345",
			orderNumber: base32.GenerateOrderNumber(12345),
			wantID:      12345,
			wantErr:     false,
		},
		{
			name:        "valid - roundtrip large ID",
			orderNumber: base32.GenerateOrderNumber(9999999999),
			wantID:      9999999999,
			wantErr:     false,
		},
		{
			name:        "error - invalid prefix",
			orderNumber: "XX-ABCD",
			wantID:      0,
			wantErr:     true,
		},
		{
			name:        "error - too short",
			orderNumber: "DV",
			wantID:      0,
			wantErr:     true,
		},
		{
			name:        "error - empty string",
			orderNumber: "",
			wantID:      0,
			wantErr:     true,
		},
		{
			name:        "error - invalid base32",
			orderNumber: "DV-!@#$%",
			wantID:      0,
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id, err := base32.DecodeOrderNumber(tc.orderNumber)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantID, id)
		})
	}
}

// =============================================================================
// ROUNDTRIP TEST
// =============================================================================

func TestOrderNumber_RoundTrip(t *testing.T) {
	testIDs := []uint64{0, 1, 100, 12345, 999999, 9999999999}

	for _, id := range testIDs {
		t.Run("", func(t *testing.T) {
			// Generate order number
			orderNumber := base32.GenerateOrderNumber(id)

			// Decode back to ID
			decodedID, err := base32.DecodeOrderNumber(orderNumber)

			assert.NoError(t, err)
			assert.Equal(t, id, decodedID, "roundtrip failed for ID %d", id)
		})
	}
}
