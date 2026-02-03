package main

import (
	"fmt"
	"time"
)

// =============================================================================
// 1. UNBUFFERED CHANNEL: GIAO HÀNG TẬN TAY
// =============================================================================
// Sender bị chặn (block) cho đến khi Receiver nhận hàng.
// Receiver bị chặn cho đến khi Sender đưa hàng.
// -> Đồng bộ chặt chẽ (Synchronization).
func DemoUnbuffered() {
	fmt.Println("\n--- 1. UNBUFFERED (Giao tận tay) ---")
	ch := make(chan string) // Không có số size

	go func() {
		fmt.Println("[Shipper] Đang tới cửa nhà...")
		time.Sleep(1 * time.Second) // Giả lập đi đường
		fmt.Println("[Shipper] Alo! Ra lấy hàng đi!")
		ch <- "Gói hàng Iphone 15" // BLOCK tại đây cho đến khi có người nhận
		fmt.Println("[Shipper] Khách đã nhận, đi giao đơn khác.")
	}()

	fmt.Println("[Khách] Đang ngủ...")
	time.Sleep(2 * time.Second) // Shipper phải chờ khách ngủ dậy
	fmt.Println("[Khách] Đã dậy, ra mở cửa.")
	msg := <-ch // Nhận hàng
	fmt.Printf("[Khách] Đã nhận: %s\n", msg)
}

// =============================================================================
// 2. BUFFERED CHANNEL: HỘP THƯ (LOCKER)
// =============================================================================
// Sender chỉ bị chặn nếu Hộp thư ĐẦY.
// Receiver chỉ bị chặn nếu Hộp thư RỖNG.
// -> Bất đồng bộ (Asynchronous) trong giới hạn cho phép.
func DemoBuffered() {
	fmt.Println("\n--- 2. BUFFERED (Hộp thư size=2) ---")
	ch := make(chan string, 2) // Size = 2

	go func() {
		fmt.Println("[Shipper] Bỏ gói 1 vào hòm thư...")
		ch <- "Gói 1" // Không block vì còn chỗ
		fmt.Println("[Shipper] Bỏ gói 2 vào hòm thư...")
		ch <- "Gói 2" // Không block vì vẫn còn chỗ
		
		fmt.Println("[Shipper] Cố nhét gói 3...")
		ch <- "Gói 3" // BLOCK! Vì hòm thư đầy (size=2). Phải chờ khách lấy bớt.
		fmt.Println("[Shipper] Nhét xong gói 3 (sau khi khách lấy bớt).")
		close(ch)
	}()

	time.Sleep(1 * time.Second) // Giả lập khách bận
	fmt.Println("[Khách] Đi kiểm tra hòm thư...")
	for msg := range ch {
		fmt.Printf("[Khách] Lấy được: %s\n", msg)
		time.Sleep(500 * time.Millisecond) // Lấy từ từ
	}
}

// =============================================================================
// 3. DEADLOCK SCENARIOS (KỊCH BẢN CHẾT ĐỨNG)
// =============================================================================
/*
	Deadlock xảy ra khi tất cả Goroutine đều đang NGỦ (asleep) chờ đợi nhau 
	và không ai đánh thức ai.
*/

func DemoDeadlock1_Unbuffered() {
	fmt.Println("\n--- 3.1 Deadlock Unbuffered ---")
	ch := make(chan int)
	
	// Lỗi: Gửi trên Main Thread mà không có ai nhận (bên Goroutine khác)
	// Main thread tự block chính mình chờ người nhận -> Mà ai chạy nữa đâu mà nhận? -> Crash
	ch <- 1 
}

func DemoDeadlock2_BufferedFull() {
	fmt.Println("\n--- 3.2 Deadlock Buffered Full ---")
	ch := make(chan int, 1)
	
	ch <- 1 // OK
	// Lỗi: Hộp đầy rồi, cố nhét tiếp thì Block chờ người lấy.
	// Nhưng đang ở Main Thread, block rồi thì ai chạy code để lấy ra? -> Crash
	ch <- 2 
}


// =============================================================================
// 4. TRÁNH DEADLOCK BẰNG SELECT (Non-blocking & Timeout)
// =============================================================================

func DemoSelect_NonBlocking() {
	fmt.Println("\n--- 4.1. Select with Default (Non-blocking) ---")
	ch := make(chan int) // Unbuffered

	// Thử gửi
	select {
	case ch <- 1:
		fmt.Println("Gửi thành công (sẽ không vào đây vì không ai nhận)")
	default:
		fmt.Println("Kênh bận/không ai nhận -> Bỏ qua, không bị Deadlock!")
	}
}

func DemoSelect_Timeout() {
	fmt.Println("\n--- 4.2. Select with Timeout ---")
	ch := make(chan int)

	go func() {
		time.Sleep(2 * time.Second) // Giả vờ xử lý lâu
		ch <- 100
	}()

	fmt.Println("Đang chờ dữ liệu (tối đa 1s)...")
	select {
	case res := <-ch:
		fmt.Printf("Nhận được: %d\n", res)
	case <-time.After(1 * time.Second): // Timeout sau 1s
		fmt.Println("Timeout! Chờ lâu quá nên thôi.")
	}
}

func main() {
	DemoUnbuffered()
	DemoBuffered()
	
	// Cách xử lý Deadlock thông minh:
	DemoSelect_NonBlocking()
	DemoSelect_Timeout()

	// DemoDeadlock1_Unbuffered()
	// DemoDeadlock2_BufferedFull()
}
