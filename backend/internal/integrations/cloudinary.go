/**
 * CLOUDINARY SERVICE
 *
 * Image upload to Cloudinary
 */

package integrations

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	// CloudinaryBaseURL can be overridden for testing
	CloudinaryBaseURL = "https://api.cloudinary.com/v1_1"
)

type CloudinaryUploadResult struct {
	PublicID  string `json:"public_id"`
	SecureURL string `json:"secure_url"`
	URL       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Format    string `json:"format"`
}

// UploadToCloudinary: Upload image to Cloudinary
func UploadToCloudinary(file io.Reader, filename string) (*CloudinaryUploadResult, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	uploadPreset := os.Getenv("CLOUDINARY_UPLOAD_PRESET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("Cloudinary not configured")
	}

	// Read file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Build form data
	formData := url.Values{}
	formData.Set("file", string(fileBytes))
	if uploadPreset != "" {
		formData.Set("upload_preset", uploadPreset)
	}

	// Generate signature if not using upload preset
	if uploadPreset == "" {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		formData.Set("timestamp", timestamp)

		// Generate signature
		params := map[string]string{
			"timestamp": timestamp,
		}
		signature := generateSignature(params, apiSecret)
		formData.Set("signature", signature)
		formData.Set("api_key", apiKey)
	}

	// Upload URL using variable base URL
	uploadURL := fmt.Sprintf("%s/%s/image/upload", CloudinaryBaseURL, cloudName)

	// Make request
	resp, err := http.PostForm(uploadURL, formData)
	if err != nil {
		return nil, fmt.Errorf("failed to upload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload failed: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result CloudinaryUploadResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// generateSignature: Generate Cloudinary signature
func generateSignature(params map[string]string, apiSecret string) string {
	// Sort keys
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build string to sign
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	signString := strings.Join(parts, "&") + apiSecret

	// Hash
	hash := sha1.Sum([]byte(signString))
	return fmt.Sprintf("%x", hash)
}

// UploadBase64ToCloudinary: Upload base64 image
func UploadBase64ToCloudinary(base64Data string) (*CloudinaryUploadResult, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	uploadPreset := os.Getenv("CLOUDINARY_UPLOAD_PRESET")

	if cloudName == "" {
		return nil, fmt.Errorf("Cloudinary not configured")
	}

	// Build form data
	formData := url.Values{}
	formData.Set("file", base64Data)
	if uploadPreset != "" {
		formData.Set("upload_preset", uploadPreset)
	}

	// Upload URL using variable base URL
	uploadURL := fmt.Sprintf("%s/%s/image/upload", CloudinaryBaseURL, cloudName)

	// Make request
	resp, err := http.PostForm(uploadURL, formData)
	if err != nil {
		return nil, fmt.Errorf("failed to upload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload failed: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result CloudinaryUploadResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}
