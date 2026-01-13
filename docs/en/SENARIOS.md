Tuy nhiên, là một Architect, tao list nhanh cho mày 3 trường hợp mày BẮT BUỘC phải dùng Go Concurrency cho Donald Vibe sau này, để mày cất vào kho vũ khí:

1. Fire-and-Forget (Bắn và Quên) - Background Jobs
   Sau khi user mua hàng thành công (Atomic DB xong), mày cần gửi Email xác nhận hoặc bắn thông báo về Telegram cho mày.

Sai lầm: Để user chờ quay quay cái vòng tròn loading trong lúc mày connect tới Gmail server (mất 2-3s).

Go Concurrency: go sendEmail() -> Trả về "Mua thành công" ngay lập tức cho user sướng. Email gửi chậm 5s kệ mẹ nó.

2. Scatter-Gather (Tản ra và Gom lại)
   Sau này mày cần tính phí ship. Mày muốn so sánh giá giữa Giao Hàng Tiết Kiệm, Grab, và Ahamove.

Sai lầm: Hỏi GHTK (1s) -> xong -> Hỏi Grab (1s) -> xong -> Hỏi Ahamove (1s). Tổng chờ = 3s.

Go Concurrency: Bắn 3 thằng Goroutine đi hỏi 3 hãng cùng lúc. Thằng nào trả lời lâu nhất mất 1s. Tổng chờ = 1s.

3. Data Processing (Xử lý nặng)
   Mày upload ảnh đồng hồ 4K lên. Server cần resize ra 3 bản: Mobile, Web, Thumbnail.

Sai lầm: Resize lần lượt từng cái. CPU chạy 1 luồng.

Go Concurrency: Bắn 3 luồng xử lý 3 cái ảnh cùng lúc. Tận dụng đa nhân CPU.

Chốt hạ: Bài toán Inventory (Kho hàng) đã giải quyết xong bằng SQLite Atomic. Đóng hòm. Cất Go Concurrency đi, khi nào làm tính năng Gửi Mail hay Xử lý ảnh thì lôi ra dùng.
