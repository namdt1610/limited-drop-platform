# VPS Configuration - Hệ thống Drop 9,000 RPS

## Cấu hình yêu cầu

- **2-4 cores CPU** (xử lý 4,000-10,000 RPS)
- **1-2GB RAM** (100 reader connections)
- **SSD NVMe** (WAL mode)
- **Ubuntu 24.04 LTS**

Chi phí: **72,000 - 240,000 VND/tháng (CloudFly)**

---

## Giá CloudFly (Khuyên dùng)

| Tier          | Cấu hình       | Giá/tháng   | RPS   | Dùng cho        |
| ------------- | -------------- | ----------- | ----- | --------------- |
| **Tier 1**    | 1CPU/1GB/20GB  | 72,000 VND  | 4-5k  | Test            |
| **Tier 2** ⭐ | 2CPU/2GB/60GB  | 176,000 VND | 9-10k | **Drop đầu**    |
| **Tier 2+**   | 2CPU/4GB/80GB  | 240,000 VND | 9-10k | **Safe margin** |
| Tier 3        | 4CPU/8GB/120GB | 416,000 VND | 15k+  | Recurring       |

**Khuyên**: CloudFly Standard 2CPU/2GB (176k/tháng) = tốt nhất cho Việt Nam

---

## Deploy CloudFly - 3 bước

### Bước 1: Tạo server

1. Vào: https://my.cloudfly.vn/cloud/server/deploy
2. Chọn: **Standard**
3. Cấu hình: **2CPU / 2GB / 60GB**
4. Region: **Việt Nam 01 (Đà Nẵng)**
5. OS: **Ubuntu 24.04 LTS**

### Bước 2: Setup

```bash
apt update && apt upgrade -y

# Install Go 1.25
wget https://go.dev/dl/go1.25.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install Bun
curl -fsSL https://bun.sh/install | bash

# Install tools
apt install -y nginx build-essential git

# Clone & build
git clone <repo> /opt/ecommerce
cd /opt/ecommerce/backend
make build-prod
./bin/server
```

### Bước 3: Test

```bash
k6 run k6-runner.js -e K6_MODE=kpi-10k -e RATE=9000 -e DURATION=60s
```

Kỳ vọng: 540k+ requests, < 0.1% errors, 50-100ms latency

---

## Performance

| Tier   | CPU | RAM | RPS   | P99 Latency |
| ------ | --- | --- | ----- | ----------- |
| Tier 1 | 2c  | 1GB | 4-5k  | 150ms       |
| Tier 2 | 2c  | 2GB | 9-10k | 50-100ms    |
| Tier 3 | 4c  | 8GB | 15k+  | < 50ms      |

---

## Liên hệ CloudFly

- **Hotline**: 0904.558.448
- **Website**: https://cloudfly.vn/
- **Deploy**: https://my.cloudfly.vn/cloud/server/deploy
- **Trial 3 ngày**: https://my.cloudfly.vn/cloud/server/trial

---

## Lưu ý

- **Region**: Việt Nam 01 (Đà Nẵng) - latency thấp
- **Tier 2 đủ dùng**: 176k/tháng xử lý 9-10k RPS
- **Upgrade dễ**: Nếu traffic > 5k, upgrade trong 5 phút
- **Thực tế**: Launch day thường 2-5k RPS, không phải 9k
- **Backup**: Dùng systemd script (tự động hàng ngày)
