package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Flags
	dbPath := flag.String("db", "database.db", "Path to SQLite database file")
	outputDir := flag.String("output", "./backups", "Output directory for backups")
	compress := flag.Bool("compress", true, "Compress backup with gzip")
	flag.Parse()

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating backup directory: %v\n", err)
		os.Exit(1)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupFile := filepath.Join(*outputDir, fmt.Sprintf("database_%s.db", timestamp))
	if *compress {
		backupFile += ".gz"
	}

	// Open source database
	source, err := os.Open(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer source.Close()

	// Create backup file
	destination, err := os.Create(backupFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating backup file: %v\n", err)
		os.Exit(1)
	}
	defer destination.Close()

	// Copy with optional compression
	var writer io.Writer = destination
	if *compress {
		gzWriter := gzip.NewWriter(destination)
		defer gzWriter.Close()
		writer = gzWriter
	}

	bytesWritten, err := io.Copy(writer, source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error copying database: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Backup created: %s\n", backupFile)
	fmt.Printf("   Size: %.2f MB\n", float64(bytesWritten)/1024/1024)

	// Cleanup old backups (keep last 30 days)
	if err := cleanupOldBackups(*outputDir, 30); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to cleanup old backups: %v\n", err)
	}
}

// cleanupOldBackups removes backups older than days
func cleanupOldBackups(dir string, days int) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	cutoffTime := time.Now().AddDate(0, 0, -days)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Skip non-backup files
		if !filepath.HasPrefix(entry.Name(), "database_") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoffTime) {
			path := filepath.Join(dir, entry.Name())
			if err := os.Remove(path); err == nil {
				fmt.Printf("🗑️  Removed old backup: %s\n", entry.Name())
			}
		}
	}

	return nil
}
