package service_test

import (
	"errors"
	"testing"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/service"
)

// =============================================================================
// SYMBICODE SERVICE TESTS
// =============================================================================

func TestGenerateSymbicode_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		productID uint64
		orderID   *uint64
		setup     func(*mockRepository)
		wantErr   bool
	}{
		{
			name:      "success - with order ID",
			productID: 10,
			orderID:   func() *uint64 { id := uint64(1); return &id }(),
			setup:     func(m *mockRepository) {},
			wantErr:   false,
		},
		{
			name:      "success - without order ID",
			productID: 10,
			orderID:   nil,
			setup:     func(m *mockRepository) {},
			wantErr:   false,
		},
		{
			name:      "error - repository create error",
			productID: 10,
			orderID:   nil,
			setup: func(m *mockRepository) {
				m.createSymErr = errors.New("database error")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := newMockRepository()
			tc.setup(repo)
			srv := service.NewService(repo, nil, nil, nil)

			sym, err := srv.GenerateSymbicode(tc.productID, tc.orderID)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
			if sym == nil {
				t.Fatal("expected symbicode, got nil")
			}
			if len(sym.Code) != 16 {
				t.Fatalf("expected 16-byte UUID, got %d bytes", len(sym.Code))
			}
		})
	}
}

func TestVerifySymbicode_TableDriven(t *testing.T) {
	// Pre-generate a valid UUID for testing
	validUUID := "01234567-89ab-cdef-0123-456789abcdef"
	validCode := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}

	tests := []struct {
		name           string
		codeStr        string
		setup          func(*mockRepository)
		wantErr        bool
		wantFirstCheck bool
	}{
		{
			name:    "success - first activation",
			codeStr: validUUID,
			setup: func(m *mockRepository) {
				m.symbicodes[string(validCode)] = &models.Symbicode{
					ID: 1, Code: validCode, ProductID: 10, IsActivated: 0,
				}
			},
			wantErr:        false,
			wantFirstCheck: true,
		},
		{
			name:    "success - already activated",
			codeStr: validUUID,
			setup: func(m *mockRepository) {
				m.symbicodes[string(validCode)] = &models.Symbicode{
					ID: 1, Code: validCode, ProductID: 10, IsActivated: 1,
				}
			},
			wantErr:        false,
			wantFirstCheck: false,
		},
		{
			name:    "error - invalid UUID format",
			codeStr: "not-a-valid-uuid",
			setup:   func(m *mockRepository) {},
			wantErr: true,
		},
		{
			name:    "error - symbicode not found",
			codeStr: validUUID,
			setup:   func(m *mockRepository) {},
			wantErr: true,
		},
		{
			name:    "error - whitespace handling",
			codeStr: "  " + validUUID + "  ",
			setup: func(m *mockRepository) {
				m.symbicodes[string(validCode)] = &models.Symbicode{
					ID: 1, Code: validCode, ProductID: 10, IsActivated: 0,
				}
			},
			wantErr:        false,
			wantFirstCheck: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := newMockRepository()
			tc.setup(repo)
			srv := service.NewService(repo, nil, nil, nil)

			sym, isFirst, err := srv.VerifySymbicode(tc.codeStr)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
			if sym == nil {
				t.Fatal("expected symbicode, got nil")
			}
			if isFirst != tc.wantFirstCheck {
				t.Fatalf("expected isFirst=%v, got %v", tc.wantFirstCheck, isFirst)
			}
		})
	}
}
