package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// =============================================================================
// CÁCH 1: Đảm bảo có người nhận (Synchronization)
// Nguyên tắc: Sender chờ Receiver. Phải có Receiver thì Sender mới thoát được.
// =============================================================================
func Fix1_Receiver() {
	ch := make(chan int) // Unbuffered

	// Goroutine gửi (Sender)
	go func() {
		ch <- 1 // Block chờ receiver
		// fmt.Println("Fix1: Sender done")
	}()

	// Goroutine nhận (Receiver)
	go func() {
		<-ch // Hứng lấy dữ liệu -> gỡ block cho sender
		// fmt.Println("Fix1: Receiver done")
	}()
}

// =============================================================================
// CÁCH 2: Dùng Buffered Channel (Fire & Forget)
// Nguyên tắc: Cho Sender một cái "giỏ". Bỏ vào giỏ xong đi luôn, không cần chờ ai nhận.
// =============================================================================
func Fix2_Buffered() {
	ch := make(chan int, 1) // Buffer size = 1

	go func() {
		ch <- 1 // Bỏ vào buffer -> Không bị block -> Thoát luôn
		// fmt.Println("Fix2: Buffered done (Dù không ai nhận)")
	}()
}

// =============================================================================
// CÁCH 3: Dùng Context (Timeout / Cancellation) - CHUẨN MỰC NHẤT
// Nguyên tắc: Sếp (Context) bảo nghỉ thì nghỉ. Không chờ nữa.
// =============================================================================
func Fix3_Context() {
	// Tạo context có timeout 10ms
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	ch := make(chan int)

	go func() {
		select {
		case ch <- 1: // Cố gửi...
			fmt.Println("Gửi thành công")
		case <-ctx.Done(): // ...Nhưng nếu lâu quá (hết giờ)
			// fmt.Println("Fix3: Timeout! Hủy kèo, đi về.")
			return // Thoát Goroutine
		}
	}()
}

// =============================================================================
// CÁCH 4: Select Default (Non-blocking Send)
// Nguyên tắc: Thử gửi, nếu tắc đường thì quay xe làm việc khác ngay.
// =============================================================================
func Fix4_SelectDefault() {
	ch := make(chan int)

	go func() {
		select {
		case ch <- 1:
			fmt.Println("Gửi được")
		default:
			// fmt.Println("Fix4: Kênh bận/không ai nghe -> Bỏ qua, thoát luôn.")
			return
		}
	}()
}

func main() {
	// In số Goroutine ban đầu (thường là 1 - main)
	initial := runtime.NumGoroutine()
	fmt.Printf("Initial Goroutines: %d\n", initial)

	// Chạy thử loop mỗi cách 100 lần để xem có leak không
	for range 100 {
		Fix1_Receiver()
		Fix2_Buffered()
		Fix3_Context()
		Fix4_SelectDefault()
	}

	// Đợi cho các Goroutine chạy xong và dọn dẹp
	time.Sleep(1 * time.Second)
	runtime.GC() // Ép dọn rác cho chắc

	final := runtime.NumGoroutine()
	fmt.Printf("Final Goroutines:   %d\n", final)

	if final == initial {
		fmt.Println("✅ TUYỆT VỜI! KHÔNG CÓ LEAK NÀO.")
	} else {
		fmt.Printf("❌ CẢNH BÁO: VẪN CÒN %d CON KIẾN ĐANG LƯU LẠC!\n", final-initial)
	}
}
