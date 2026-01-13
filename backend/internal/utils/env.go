package utils

import (
	"log"
	"os"
	"strings"
)

func loadEnv() {
	file, err := os.ReadFile(".env")
	if err != nil {
		log.Println(".env path not found")
		return
	}
	data := string(file)
	content := strings.Split(data, "\n")
	for _, line := range content {
		if k, v, ok := ParseLine(line); ok {
			os.Setenv(k, v)
		}
	}

}

// ParseLine parses a single env file line and returns key, value and ok
// Rules inferred from tests in env_test.go:
// - trim whitespace
// - ignore empty lines or comments starting with '#'
// - split on the first '=' and trim key/value
// - return ok=false for invalid lines
func ParseLine(line string) (string, string, bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])
	if key == "" {
		return "", "", false
	}
	return key, val, true
}
