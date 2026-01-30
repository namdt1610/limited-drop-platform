package integrations_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"ecommerce-backend/internal/integrations"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendEmail_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		reqTo       []string
		reqSubject  string
		mockHandler func(w http.ResponseWriter, r *http.Request)
		mockEnv     map[string]string
		wantErr     bool
	}{
		{
			name:       "success",
			reqTo:      []string{"test@example.com"},
			reqSubject: "Test Subject",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/emails", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))

				// Check Request Body
				var req integrations.ResendEmailRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				assert.NoError(t, err)
				assert.Equal(t, "test@example.com", req.To[0])
				assert.Equal(t, "Test Subject", req.Subject)

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id":"re_123456789"}`))
			},
			mockEnv: map[string]string{
				"RESEND_API_KEY": "test-api-key",
			},
			wantErr: false,
		},
		{
			name:       "error - api error",
			reqTo:      []string{"test@example.com"},
			reqSubject: "Test Subject",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("bad request"))
			},
			mockEnv: map[string]string{
				"RESEND_API_KEY": "test-api-key",
			},
			wantErr: true,
		},
		{
			name:       "error - missing api key",
			reqTo:      []string{"test@example.com"},
			reqSubject: "Test Subject",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {},
			mockEnv:     map[string]string{}, // Missing API Key
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tc.mockHandler))
			defer server.Close()

			for k, v := range tc.mockEnv {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}
			os.Setenv("RESEND_BASE_URL", server.URL)
			defer os.Unsetenv("RESEND_BASE_URL")

			err := integrations.SendEmail(tc.reqTo, tc.reqSubject, "<p>Hello</p>")

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
