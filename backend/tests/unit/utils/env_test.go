package utils_test

import (
	"testing"

	utils "ecommerce-backend/internal/utils"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantKey string
		wantVal string
		wantOk  bool
	}{
		{"Dòng chuẩn", "DB_PORT=5432", "DB_PORT", "5432", true},
		{"Có khoảng trắng", "  PORT = 8080  ", "PORT", "8080", true},
		{"Value có dấu bằng", "KEY=stark=secret", "KEY", "stark=secret", true},
		{"Dòng trống", "   ", "", "", false},
		{"Dòng comment", "# DATABASE_URL=localhost", "", "", false},
		{"Không có dấu bằng", "INVALID_LINE", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, val, ok := utils.ParseLine(tt.input)
			if ok != tt.wantOk {
				t.Fatalf("expected ok=%v, got=%v for input %q", tt.wantOk, ok, tt.input)
			}
			if key != tt.wantKey {
				t.Fatalf("expected key=%q, got=%q for input %q", tt.wantKey, key, tt.input)
			}
			if val != tt.wantVal {
				t.Fatalf("expected val=%q, got=%q for input %q", tt.wantVal, val, tt.input)
			}
		})
	}
}
