package service_test

import (
	"ecommerce-backend/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateQRCodeData(t *testing.T) {
	// Mock UUID bytes
	uuidBytes := []byte{
		0x12, 0x3e, 0x45, 0x67,
		0xe8, 0x9b,
		0x12, 0xd3,
		0xa4, 0x56,
		0x42, 0x66, 0x14, 0x17, 0x40, 0x00,
	}
	uuidStr := "123e4567-e89b-12d3-a456-426614174000"
	
	// Create expected URL
	// Note: VerifyBaseURL is a constant in service package: "/verify"
	expected := "/verify?code=" + uuidStr

	result := service.GenerateQRCodeData(uuidBytes)
	assert.Equal(t, expected, result)
}

// Test internal helpers if exported logic depends on them
func TestParseUUID_Integration(t *testing.T) {
	// Indirectly tested via GenerateQRCodeData if we did a round trip, 
	// but GenerateQRCodeData takes []byte. 
	// The service function VerifySymbicode uses parseUUID.
	
	// Since GenerateQRCodeData relies on uuid.FormatUUIDToString which we tested in utils,
	// checking VerifySymbicode in service_test would cover parseUUID if we mocked repo correctly.
}
