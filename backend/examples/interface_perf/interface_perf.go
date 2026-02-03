package main

import (
	"fmt"
	"time"
)

// =============================================================================
// 1. interface{} vs any
// =============================================================================
// `any` chỉ là alias của `interface{}` (Go 1.18+).
// Không có sự khác biệt nào về lúc biên dịch hay lúc chạy.
// `any` ngắn hơn, dễ đọc hơn.
type EmptyInterface interface{}
type AnyAlias = any

// =============================================================================
// BENCHMARK HELPER
// =============================================================================
func measure(name string, loops int, f func()) {
	start := time.Now()
	for range loops {
		f()
	}
	duration := time.Since(start)
	fmt.Printf("[%s] %v (%.2f ns/op)\n", name, duration, float64(duration.Nanoseconds())/float64(loops))
}

func main() {
	var val any = 42 // Một biến interface giữ giá trị int
	loops := 100_000_000

	fmt.Println("--- SO SÁNH HIỆU NĂNG ---")

	// 1. TYPE ASSERTION (Direct check)
	// Kiểm tra trực tiếp: "Mày có phải là int không?"
	measure("Type Assertion (OK)", loops, func() {
		if i, ok := val.(int); ok {
			_ = i
		}
	})

	// 2. TYPE ASSERTION (Fail check)
	// Kiểm tra sai kiểu: "Mày có phải là string không?"
	measure("Type Assertion (Fail)", loops, func() {
		if _, ok := val.(string); ok {
			// Không bao giờ vào đây
		}
	})

	// 3. TYPE SWITCH (Single case)
	// Switch với 1 case: Tương đương assertion nhưng overhead context switch chút xíu
	measure("Type Switch (1 case)", loops, func() {
		switch v := val.(type) {
		case int:
			_ = v
		}
	})

	// 4. TYPE SWITCH (Multiple cases)
	// Switch phải so sánh tuần tự hoặc nhảy bảng (tùy compiler optimization)
	measure("Type Switch (Many cases)", loops, func() {
		switch v := val.(type) {
		case string:
			_ = v
		case bool:
			_ = v
		case float64:
			_ = v
		case int: // Case đúng nằm ở cuối
			_ = v
		}
	})

	fmt.Println("\n--- KẾT LUẬN ---")
	fmt.Println("1. interface{} và any là MỘT. 100% giống nhau.")
	fmt.Println("2. Type Assertion nhanh nhất nếu bạn biết chính xác mình cần gì.")
	fmt.Println("3. Type Switch chậm hơn xíu do cơ chế nhảy case, nhưng code sạch hơn nếu handle nhiều type.")
}
