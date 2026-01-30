/**
 * BREVO SERVICE
 *
 * Email service using Brevo (Sendinblue)
 * Free: 300 emails/day
 * Paid: $20/month for unlimited
 */

package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type BrevoEmailRequest struct {
	Sender      BrevoSender      `json:"sender"`
	To          []BrevoRecipient `json:"to"`
	Subject     string           `json:"subject"`
	HtmlContent string           `json:"htmlContent"`
}

type BrevoSender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type BrevoRecipient struct {
	Email string `json:"email"`
}

type BrevoEmailResponse struct {
	MessageID int64 `json:"messageId"`
}

// SendEmailBrevo: Send email via Brevo
func SendEmailBrevo(to []string, subject, htmlContent string) error {
	apiKey := os.Getenv("BREVO_API_KEY")

	if apiKey == "" {
		return fmt.Errorf("brevo api key not configured")
	}

	// Build recipients
	recipients := make([]BrevoRecipient, len(to))
	for i, email := range to {
		recipients[i] = BrevoRecipient{Email: email}
	}

	req := BrevoEmailRequest{
		Sender: BrevoSender{
			Name:  "Donald Watch",
			Email: "noreply@donaldwatch.xyz",
		},
		To:          recipients,
		Subject:     subject,
		HtmlContent: htmlContent,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	baseURL := os.Getenv("BREVO_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.brevo.com/v3"
	}

	httpReq, err := http.NewRequest("POST", baseURL+"/smtp/email", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("api-key", apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("brevo returned status %d", resp.StatusCode)
	}

	var respBody BrevoEmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if respBody.MessageID == 0 {
		return fmt.Errorf("invalid message id from brevo")
	}

	return nil
}
