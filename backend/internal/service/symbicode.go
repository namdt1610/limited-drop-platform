package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/utils/uuid"

	googleuuid "github.com/google/uuid"
)

const VerifyBaseURL = "/verify"

// GenerateSymbicode creates and persists a new symbicode for a product/order
func (s *service) GenerateSymbicode(productID uint64, orderID *uint64) (*models.Symbicode, error) {
	code, err := uuid.GenerateUUIDv7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate uuid: %w", err)
	}
	secret := generateSecretKey()

	sym := &models.Symbicode{
		Code:        code,
		SecretKey:   secret,
		ProductID:   productID,
		IsActivated: 0,
	}
	if orderID != nil {
		sym.OrderID = *orderID
	}

	if err := s.repo.CreateSymbicode(sym); err != nil {
		return nil, fmt.Errorf("failed to create symbicode: %w", err)
	}
	return sym, nil
}

// VerifySymbicode verifies and activates a symbicode. Returns (symbicode, isFirstActivation, error)
func (s *service) VerifySymbicode(codeStr string) (*models.Symbicode, bool, error) {
	codeStr = strings.TrimSpace(codeStr)

	parsed, err := parseUUID(codeStr)
	if err != nil {
		return nil, false, fmt.Errorf("invalid symbicode format: %w", err)
	}

	symbicode, err := s.repo.GetSymbicodeByCode(parsed)
	if err != nil {
		return nil, false, err
	}

	isFirst := symbicode.IsActivated&1 == 0
	if isFirst {
		if err := s.repo.ActivateSymbicode(symbicode.ID, ""); err != nil {
			return nil, false, fmt.Errorf("failed to activate symbicode: %w", err)
		}
		// refresh the record
		symbicode, err = s.repo.GetSymbicodeByCode(parsed)
		if err != nil {
			return nil, false, err
		}
	}

	return symbicode, isFirst, nil
}

// GenerateQRCodeData formats verification URL with the UUID string
func GenerateQRCodeData(code []byte) string {
	return fmt.Sprintf("%s?code=%s", VerifyBaseURL, uuid.FormatUUIDToString(code))
}

func generateSecretKey() string {
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	hash := sha256.Sum256(randomBytes)
	return hex.EncodeToString(hash[:])
}

// parseUUID parses a UUID string to binary format (16 bytes)
func parseUUID(uuidStr string) ([]byte, error) {
	parsed, err := googleuuid.Parse(uuidStr)
	if err != nil {
		return nil, err
	}
	return parsed[:], nil
}
