package integrations

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	sheets "google.golang.org/api/sheets/v4"
)

// SubmitOrderToGoogleSheet appends a row to a Google Spreadsheet using a
// service account. If not configured, it is a silent no-op to avoid
// impacting order processing.
func SubmitOrderToGoogleSheet(name, phone, email, address, dropName string, amount float64, paidAt time.Time) error {
	sheetID := os.Getenv("GSSHEET_SPREADSHEET_ID")
	if sheetID == "" {
		// Not configured; noop
		return nil
	}

	credPath := os.Getenv("GDRIVE_SERVICE_ACCOUNT")
	if credPath == "" {
		credPath = "./gdrive-service-account.json"
	}

	b, err := os.ReadFile(credPath)
	if err != nil {
		return fmt.Errorf("read service account: %w", err)
	}

	conf, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		return fmt.Errorf("jwt config: %w", err)
	}

	client := conf.Client(context.Background())
	srv, err := sheets.New(client)
	if err != nil {
		return fmt.Errorf("sheets service: %w", err)
	}

	sheetName := os.Getenv("GSSHEET_SHEET_NAME")
	if sheetName == "" {
		sheetName = "Sheet1"
	}
	rangeStr := fmt.Sprintf("%s!A:Z", sheetName)

	vr := &sheets.ValueRange{Values: [][]interface{}{{
		paidAt.Format(time.RFC3339),
		name,
		phone,
		email,
		address,
		dropName,
		fmt.Sprintf("%.0f", amount),
	}}}

	_, err = srv.Spreadsheets.Values.Append(sheetID, rangeStr, vr).
		ValueInputOption("USER_ENTERED").
		InsertDataOption("INSERT_ROWS").
		Context(context.Background()).
		Do()
	if err != nil {
		return fmt.Errorf("append values: %w", err)
	}

	return nil
}
