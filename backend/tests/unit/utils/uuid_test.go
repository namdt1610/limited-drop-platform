package utils_test

import (
	"ecommerce-backend/internal/utils/uuid"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatUUIDToString(t *testing.T) {
	// Known UUID bytes for "123e4567-e89b-12d3-a456-426614174000"
	uuidBytes := []byte{
		0x12, 0x3e, 0x45, 0x67,
		0xe8, 0x9b,
		0x12, 0xd3,
		0xa4, 0x56,
		0x42, 0x66, 0x14, 0x17, 0x40, 0x00,
	}
	expected := "123e4567-e89b-12d3-a456-426614174000"

	result := uuid.FormatUUIDToString(uuidBytes)
	assert.Equal(t, expected, result)
}

func TestFormatUUIDToString_Empty(t *testing.T) {
	result := uuid.FormatUUIDToString(nil)
	assert.Equal(t, "", result)

	result = uuid.FormatUUIDToString([]byte{})
	assert.Equal(t, "", result)
}
