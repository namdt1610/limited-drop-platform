package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnv(t *testing.T) {
	// Create temporary .env file
	content := []byte("TEST_KEY=test_value\n#Comment\nINVALID_LINE\n   \nSPACED_KEY = spaced_value  ")
	err := os.WriteFile(".env", content, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(".env")

	// Clear env vars to be sure
	os.Unsetenv("TEST_KEY")
	os.Unsetenv("SPACED_KEY")

	// Call private loadEnv (since we are in package utils)
	loadEnv()

	// Verify
	assert.Equal(t, "test_value", os.Getenv("TEST_KEY"))
	assert.Equal(t, "spaced_value", os.Getenv("SPACED_KEY"))
}

func TestLoadEnv_NoFile(t *testing.T) {
	// Ensure no .env file exists
	os.Remove(".env")

	// Should not panic or error
	loadEnv()
}
