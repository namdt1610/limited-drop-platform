package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

// GenerateUUIDv7 generates a UUID v7 (timestamp-based) in binary format
func GenerateUUIDv7() ([]byte, error) {
	uuid := make([]byte, 16)

	// Generate 16 bytes of random data
	if _, err := rand.Read(uuid); err != nil {
		return nil, err
	}

	// Get current timestamp in milliseconds (48 bits)
	timestamp := uint64(time.Now().UnixMilli())

	// Set timestamp (first 6 bytes, big-endian)
	binary.BigEndian.PutUint64(uuid[0:8], timestamp)
	// Shift right to fit in 48 bits
	uuid[0] = byte(timestamp >> 40)
	uuid[1] = byte(timestamp >> 32)
	uuid[2] = byte(timestamp >> 24)
	uuid[3] = byte(timestamp >> 16)
	uuid[4] = byte(timestamp >> 8)
	uuid[5] = byte(timestamp)

	// Set version (4 bits: 0111 = 7)
	uuid[6] = (uuid[6] & 0x0F) | 0x70

	// Set variant (2 bits: 10)
	uuid[8] = (uuid[8] & 0x3F) | 0x80

	return uuid, nil
}

// FormatUUIDToString converts binary UUID to standard string format
func FormatUUIDToString(uuid []byte) string {
	if len(uuid) != 16 {
		return ""
	}
	// Convert binary UUID to string format
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}
