package integrations_test

import (
	"testing"

	"ecommerce-backend/internal/integrations"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// ADDRESS VALIDATION TESTS
// =============================================================================

func TestValidateAddress_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		wantErr bool
		errMsg  string
	}{
		// Valid addresses
		{
			name:    "valid - normal address",
			addr:    "123 Nguyễn Văn Cừ, Phường 4",
			wantErr: true, // Contains "Phường" - forbidden
		},
		{
			name:    "valid - house number and street only",
			addr:    "456 Lê Lợi, Tòa nhà ABC",
			wantErr: false,
		},
		{
			name:    "valid - apartment address",
			addr:    "Căn hộ 1201, Chung cư Sunrise City",
			wantErr: false,
		},

		// Too short
		{
			name:    "error - too short",
			addr:    "123 abc",
			wantErr: true,
			errMsg:  "Vui lòng ghi rõ số nhà, tên đường",
		},
		{
			name:    "error - empty string",
			addr:    "",
			wantErr: true,
			errMsg:  "Vui lòng ghi rõ số nhà, tên đường",
		},
		{
			name:    "error - whitespace only",
			addr:    "         ",
			wantErr: true,
			errMsg:  "Vui lòng ghi rõ số nhà, tên đường",
		},

		// Blacklist words
		{
			name:    "error - contains 'test'",
			addr:    "123 Đường Test Street",
			wantErr: true,
			errMsg:  "Địa chỉ không hợp lệ, vui lòng nhập nghiêm túc để nhận hàng.",
		},
		{
			name:    "error - contains 'ahihi'",
			addr:    "ahihi đồ ngốc 123 abc",
			wantErr: true,
			errMsg:  "Địa chỉ không hợp lệ, vui lòng nhập nghiêm túc để nhận hàng.",
		},
		{
			name:    "error - contains 'asdasd'",
			addr:    "asdasd asdasd asda",
			wantErr: true,
			errMsg:  "Địa chỉ không hợp lệ, vui lòng nhập nghiêm túc để nhận hàng.",
		},

		// Forbidden administrative words
		{
			name:    "error - contains 'tỉnh'",
			addr:    "123 Đường ABC, Tỉnh Bình Dương",
			wantErr: true,
			errMsg:  "Địa chỉ không hợp lệ, vui lòng nhập nghiêm túc để nhận hàng.",
		},
		{
			name:    "error - contains 'quận'",
			addr:    "456 Đường XYZ, Quận 1",
			wantErr: true,
			errMsg:  "Địa chỉ không hợp lệ, vui lòng nhập nghiêm túc để nhận hàng.",
		},
		{
			name:    "error - contains 'phường'",
			addr:    "789 Đường ABC, Phường 5",
			wantErr: true,
			errMsg:  "Địa chỉ không hợp lệ, vui lòng nhập nghiêm túc để nhận hàng.",
		},
		{
			name:    "error - contains 'tp' (thành phố)",
			addr:    "123 ABC, TP Hồ Chí Minh",
			wantErr: true,
			errMsg:  "Địa chỉ không hợp lệ, vui lòng nhập nghiêm túc để nhận hàng.",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := integrations.ValidateAddress(tc.addr)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.errMsg != "" {
					assert.Equal(t, tc.errMsg, err.Error())
				}
				return
			}

			assert.NoError(t, err)
		})
	}
}
