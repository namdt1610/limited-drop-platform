package integrations

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// SubmitOrderToGoogleForm sends a completed limited-drop order to a Google Form
// for lightweight order tracking. If misconfigured, it fails silently so that
// core payment/order flow is not affected.
func SubmitOrderToGoogleForm(name, phone, email, address, dropName string, amount float64, paidAt time.Time) {
	formURL := os.Getenv("GFORM_ORDER_URL")
	if formURL == "" {
		return
	}

	fieldName := os.Getenv("GFORM_FIELD_NAME")
	fieldPhone := os.Getenv("GFORM_FIELD_PHONE")
	fieldEmail := os.Getenv("GFORM_FIELD_EMAIL")
	fieldAddress := os.Getenv("GFORM_FIELD_ADDRESS")
	fieldDrop := os.Getenv("GFORM_FIELD_DROP")
	fieldAmount := os.Getenv("GFORM_FIELD_AMOUNT")
	fieldTxTime := os.Getenv("GFORM_FIELD_TXTIME")

	values := url.Values{}
	if fieldName != "" {
		values.Set(fieldName, name)
	}
	if fieldPhone != "" {
		values.Set(fieldPhone, phone)
	}
	if fieldEmail != "" {
		values.Set(fieldEmail, email)
	}
	if fieldAddress != "" {
		values.Set(fieldAddress, address)
	}
	if fieldDrop != "" {
		values.Set(fieldDrop, dropName)
	}
	if fieldAmount != "" {
		values.Set(fieldAmount, fmt.Sprintf("%.0f", amount))
	}
	if fieldTxTime != "" {
		values.Set(fieldTxTime, paidAt.Format(time.RFC3339))
	}

	// Fire-and-forget HTTP POST; ignore response and errors
	resp, err := http.PostForm(formURL, values)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
