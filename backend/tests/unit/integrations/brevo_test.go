package integrations_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"ecommerce-backend/internal/integrations"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendEmailBrevo_TableDriven(t *testing.T) {
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
			reqSubject: "Subject",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/smtp/email", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "test-api-key", r.Header.Get("api-key"))
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{"messageId":12345}`))
			},
			mockEnv: map[string]string{
				"BREVO_API_KEY": "test-api-key",
			},
			wantErr: false,
		},
		{
			name:       "error - api error",
			reqTo:      []string{"test@example.com"},
			reqSubject: "Subject",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			mockEnv: map[string]string{
				"BREVO_API_KEY": "test-api-key",
			},
			wantErr: true,
		},
		{
			name:       "error - missing api key",
			reqTo:      []string{"test@example.com"},
			reqSubject: "Subject",
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
			os.Setenv("BREVO_BASE_URL", server.URL)
			defer os.Unsetenv("BREVO_BASE_URL")

			err := integrations.SendEmailBrevo(tc.reqTo, tc.reqSubject, "<h1>Hello</h1>")

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
