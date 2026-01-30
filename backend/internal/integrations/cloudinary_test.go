package integrations

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSignature(t *testing.T) {
	params := map[string]string{
		"timestamp": "1234567890",
		"folder":    "products",
	}
	secret := "my_secret_key"
	
	// Expected signature based on Cloudinary logic: sha1("folder=products&timestamp=1234567890my_secret_key")
	// SHA1("folder=products&timestamp=1234567890my_secret_key")
	// e.g. using online calculator to verify or rely on deterministic output
	
	sig := generateSignature(params, secret)
	assert.NotEmpty(t, sig)
	assert.Len(t, sig, 40) // SHA1 hex length
}

func TestUploadToCloudinary(t *testing.T) {
	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-cloud/image/upload", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		
		// Check form values
		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "test-content", r.Form.Get("file"))
		
		w.WriteHeader(http.StatusOK)
		resp := CloudinaryUploadResult{
			PublicID:  "test_id",
			SecureURL: "https://res.cloudinary.com/demo/image/upload/v1/test_id.jpg",
			Format:   "jpg",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Override Base URL and Env
	originalBaseURL := CloudinaryBaseURL
	CloudinaryBaseURL = server.URL
	defer func() { CloudinaryBaseURL = originalBaseURL }()

	t.Setenv("CLOUDINARY_CLOUD_NAME", "test-cloud")
	t.Setenv("CLOUDINARY_API_KEY", "test-key")
	t.Setenv("CLOUDINARY_API_SECRET", "test-secret")
	
	// Test Upload
	reader := strings.NewReader("test-content")
	result, err := UploadToCloudinary(reader, "test.jpg")
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test_id", result.PublicID)
}

func TestUploadBase64ToCloudinary(t *testing.T) {
	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-cloud/image/upload", r.URL.Path)
		
		r.ParseForm()
		assert.Equal(t, "data:image/png;base64,xxxx", r.Form.Get("file"))
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CloudinaryUploadResult{PublicID: "b64_id"})
	}))
	defer server.Close()

	// Override
	originalBaseURL := CloudinaryBaseURL
	CloudinaryBaseURL = server.URL
	defer func() { CloudinaryBaseURL = originalBaseURL }()

	t.Setenv("CLOUDINARY_CLOUD_NAME", "test-cloud")
	
	// Test
	result, err := UploadBase64ToCloudinary("data:image/png;base64,xxxx")
	assert.NoError(t, err)
	assert.Equal(t, "b64_id", result.PublicID)
}

func TestUploadToCloudinary_Error(t *testing.T) {
	// Mock Failure
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
	}))
	defer server.Close()

	originalBaseURL := CloudinaryBaseURL
	CloudinaryBaseURL = server.URL
	defer func() { CloudinaryBaseURL = originalBaseURL }()

	t.Setenv("CLOUDINARY_CLOUD_NAME", "test-cloud")
	t.Setenv("CLOUDINARY_API_KEY", "test-key")
	t.Setenv("CLOUDINARY_API_SECRET", "test-secret")

	reader := strings.NewReader("test")
	_, err := UploadToCloudinary(reader, "test.jpg")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "upload failed")
}
