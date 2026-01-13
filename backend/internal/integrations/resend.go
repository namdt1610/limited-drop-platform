/**
 * RESEND SERVICE
 *
 * Email service using Resend
 */

package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/utils/base32"
)

type ResendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
	Text    string   `json:"text,omitempty"`
}

type ResendEmailResponse struct {
	ID string `json:"id"`
}

// SendEmail: Send email via Resend
func SendEmail(to []string, subject, htmlContent string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	fromEmail := os.Getenv("RESEND_FROM_EMAIL")

	if apiKey == "" {
		return fmt.Errorf("resend api key not configured")
	}

	if fromEmail == "" {
		fromEmail = "noreply@yourdomain.com"
	}

	req := ResendEmailRequest{
		From:    fromEmail,
		To:      to,
		Subject: subject,
		HTML:    htmlContent,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend error: %s", string(body))
	}

	return nil
}

// SendWelcomeEmail: Send welcome email to new user
func SendWelcomeEmail(email, name string) error {
	html := fmt.Sprintf(`
		<h1>Chào mừng %s!</h1>
		<p>Cảm ơn bạn đã đăng ký tài khoản tại Donald Watch.</p>
		<p>Chúc bạn mua sắm vui vẻ!</p>
	`, name)

	return SendEmailBrevo([]string{email}, "Chào mừng đến với Donald Watch", html)
}

// SendOrderConfirmationEmail: Send order confirmation email
func SendOrderConfirmationEmail(email, orderNumber string, totalAmount float64) error {
	html := fmt.Sprintf(`
		<h1>Đơn hàng của bạn đã được xác nhận</h1>
		<p>Mã đơn hàng: <strong>%s</strong></p>
		<p>Tổng tiền: <strong>%.0f VND</strong></p>
		<p>Cảm ơn bạn đã mua sắm tại Donald Watch!</p>
	`, orderNumber, totalAmount)

	return SendEmailBrevo([]string{email}, fmt.Sprintf("Xác nhận đơn hàng #%s", orderNumber), html)
}

// SendOrderDetailsEmail: Send full order details (guest lookup)
func SendOrderDetailsEmail(email string, order *models.Order) error {
	if email == "" || order == nil {
		return fmt.Errorf("missing email or order")
	}

	var itemsBuilder strings.Builder
	var orderItems []struct {
		ProductName string  `json:"product_name"`
		Quantity    int     `json:"quantity"`
		Price       float64 `json:"price"`
	}
	if err := json.Unmarshal(order.Items, &orderItems); err == nil {
		for _, item := range orderItems {
			itemsBuilder.WriteString(fmt.Sprintf(
				"<li>%s x %d - %.0f VND</li>",
				item.ProductName, item.Quantity, item.Price*float64(item.Quantity),
			))
		}
	}

	address := string(order.ShippingAddress)
	if address == "" {
		address = "Không có địa chỉ"
	}

	html := fmt.Sprintf(`
		<h2>Thông tin đơn hàng #%s</h2>
		<p>Cảm ơn bạn đã đặt hàng tại Donald Watch.</p>
		<p><strong>Tổng tiền:</strong> %d VND</p>
		<p><strong>Địa chỉ giao:</strong> %s</p>
		<p><strong>Trạng thái:</strong> %d</p>
		<p><strong>Sản phẩm:</strong></p>
		<ul>%s</ul>
		<p>Bạn có thể tra cứu đơn bằng mã đơn và email/số điện thoại tại trang: https://donaldwatch.vn/orders</p>
	`, base32.GenerateOrderNumber(order.ID), order.TotalAmount, address, order.Status, itemsBuilder.String())

	return SendEmailBrevo([]string{email}, fmt.Sprintf("Chi tiết đơn hàng #%s", base32.GenerateOrderNumber(order.ID)), html)
}

// SendSymbioteReceipt sends a high-touch "ACCESS GRANTED" receipt email
func SendSymbioteReceipt(email, maskedPhone, status, elapsed string) error {
	if email == "" {
		return fmt.Errorf("email is required for receipt")
	}

	if maskedPhone == "" {
		maskedPhone = "???"
	}
	if elapsed == "" {
		elapsed = "0.042s"
	}

	var subject string
	var bodyHTML string

	switch status {
	case "WINNER":
		status = "SECURED"
		subject = "[ACCESS GRANTED] PROJECT SYMBIOTE"
		bodyHTML = fmt.Sprintf(
			`<pre style="font-family: monospace; line-height: 1.5;">
User ID: %s
Status: %s
Time: %s

Chúc mừng. Bạn đã đánh bại hàng trăm kẻ khác.
Hệ thống đã khoá slot lại và ghi nhận quyền sở hữu của bạn.
Những kẻ chậm tay hơn chỉ còn quyền than vãn.

Signed,
DAEMON / System Architect.
</pre>`, maskedPhone, status, elapsed)
	case "LOSER":
		status = "REJECTED"
		subject = "[ACCESS DENIED] PROJECT SYMBIOTE"
		bodyHTML = fmt.Sprintf(
			`<pre style="font-family: monospace; line-height: 1.5;">
User ID: %s
Status: %s
Time: %s

Thanh toán của bạn đã đến SAU người khác.
Slot đã bị cướp trước mặt bạn, hệ thống sẽ hoàn tiền (nếu đã trừ).
Lần sau hãy quyết đoán hơn, Colosseum không chờ kẻ do dự.

Signed,
DAEMON / System Architect.
</pre>`, maskedPhone, status, elapsed)
	default:
		if status == "" {
			status = "INFO"
		}
		subject = "[NOTICE] PROJECT SYMBIOTE"
		bodyHTML = fmt.Sprintf(
			`<pre style="font-family: monospace; line-height: 1.5;">
User ID: %s
Status: %s
Time: %s

Thông báo trạng thái từ hệ thống.

Signed,
DAEMON / System Architect.
</pre>`, maskedPhone, status, elapsed)
	}

	return SendEmailBrevo([]string{email}, subject, bodyHTML)
}

// getAdminRecipients returns recipients from env or default list
func getAdminRecipients() []string {
	// Allow override via env (comma-separated)
	env := os.Getenv("ADMIN_ORDER_EMAILS")
	if env != "" {
		parts := strings.Split(env, ",")
		var recipients []string
		for _, p := range parts {
			if email := strings.TrimSpace(p); email != "" {
				recipients = append(recipients, email)
			}
		}
		if len(recipients) > 0 {
			return recipients
		}
	}

	// Default admin emails
	return []string{"n3r0.corp@gmail.com", "nam.dt161@gmail.com"}
}

// SendOrderCreatedAdminEmail notifies admins when a new order is placed
func SendOrderCreatedAdminEmail(order *models.Order) error {
	if order == nil {
		return fmt.Errorf("order is nil")
	}

	recipients := getAdminRecipients()
	if len(recipients) == 0 {
		return fmt.Errorf("no admin recipients configured")
	}

	// Parse shipping address snapshot
	address := string(order.ShippingAddress)
	var addressMap map[string]any
	if err := json.Unmarshal(order.ShippingAddress, &addressMap); err == nil {
		var parts []string
		for _, key := range []string{"name", "email", "phone", "address"} {
			if val, ok := addressMap[key]; ok && val != nil && strings.TrimSpace(fmt.Sprint(val)) != "" {
				parts = append(parts, fmt.Sprintf("<strong>%s:</strong> %s", strings.Title(key), fmt.Sprint(val)))
			}
		}
		if len(parts) > 0 {
			address = strings.Join(parts, "<br/>")
		}
	}

	var itemsBuilder strings.Builder
	var orderItems []struct {
		ProductName string  `json:"product_name"`
		Quantity    int     `json:"quantity"`
		Price       float64 `json:"price"`
	}
	if err := json.Unmarshal(order.Items, &orderItems); err == nil {
		for _, item := range orderItems {
			itemsBuilder.WriteString(fmt.Sprintf(
				"<li>%s x %d - %.0f VND</li>",
				item.ProductName, item.Quantity, item.Price*float64(item.Quantity),
			))
		}
	}

	html := fmt.Sprintf(`
		<h2>Đơn hàng mới được tạo</h2>
		<p><strong>Mã đơn:</strong> %s</p>
		<p><strong>Tổng tiền:</strong> %d VND</p>
		<p><strong>Trạng thái:</strong> %d</p>
		<p><strong>Thông tin giao:</strong><br/>%s</p>
		<p><strong>Sản phẩm:</strong></p>
		<ul>%s</ul>
	`, base32.GenerateOrderNumber(order.ID), order.TotalAmount, order.Status, address, itemsBuilder.String())

	subject := fmt.Sprintf("[DW] Đơn hàng mới #%s", base32.GenerateOrderNumber(order.ID))
	return SendEmailBrevo(recipients, subject, html)
}

// SendPasswordResetEmail: Send password reset email
func SendPasswordResetEmail(email, resetToken string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("FRONTEND_URL"), resetToken)

	html := fmt.Sprintf(`
		<h1>Đặt lại mật khẩu</h1>
		<p>Bạn đã yêu cầu đặt lại mật khẩu.</p>
		<p><a href="%s">Click vào đây để đặt lại mật khẩu</a></p>
		<p>Link này sẽ hết hạn sau 1 giờ.</p>
	`, resetURL)

	return SendEmailBrevo([]string{email}, "Đặt lại mật khẩu", html)
}
