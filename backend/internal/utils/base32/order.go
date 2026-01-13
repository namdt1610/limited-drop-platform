package base32

import (
	"encoding/base32"
	"fmt"
	"strconv"
)

// GenerateOrderNumber generates order number from ID using Base32 encoding
// Format: DV-{Base32(ID)} (e.g., DV-7X9 for ID=12345)
func GenerateOrderNumber(id uint64) string {
	idStr := strconv.FormatUint(id, 10)
	encoded := base32.StdEncoding.EncodeToString([]byte(idStr))
	return "DV-" + encoded
}

// DecodeOrderNumber decodes order number back to ID
// Input: "DV-7X9" -> Output: 12345
func DecodeOrderNumber(orderNumber string) (uint64, error) {
	if len(orderNumber) < 3 || orderNumber[:3] != "DV-" {
		return 0, fmt.Errorf("invalid order number format")
	}

	encoded := orderNumber[3:] // Remove "DV-" prefix
	decoded, err := base32.StdEncoding.DecodeString(encoded)
	if err != nil {
		return 0, err
	}

	id, err := strconv.ParseUint(string(decoded), 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}
