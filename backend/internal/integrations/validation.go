package integrations

import (
	"errors"
	"strings"
)

// ValidateAddress enforces front-end rule: only "Số nhà/Tên đường" may be
// typed by the user. It returns a localized error message for clients.
func ValidateAddress(addr string) error {
	s := strings.TrimSpace(addr)
	if len(s) < 10 {
		return errors.New("Vui lòng ghi rõ số nhà, tên đường")
	}

	lower := strings.ToLower(s)

	// Blacklist keywords (troll/spam words)
	blacklist := []string{"test", "ahihi", "dmm", "tao", "ko biet", "asdasd"}
	for _, b := range blacklist {
		if strings.Contains(lower, b) {
			return errors.New("Địa chỉ không hợp lệ, vui lòng nhập nghiêm túc để nhận hàng.")
		}
	}

	// Forbid typing administrative divisions (Tỉnh/Thành, Quận/Huyện, Phường/Xã)
	forbidden := []string{"tỉnh", "thành phố", "thanh pho", "quận", "huyện", "phường", "xã", "tp"}
	for _, f := range forbidden {
		if strings.Contains(lower, f) {
			return errors.New("Địa chỉ không hợp lệ, vui lòng nhập nghiêm túc để nhận hàng.")
		}
	}

	return nil
}
