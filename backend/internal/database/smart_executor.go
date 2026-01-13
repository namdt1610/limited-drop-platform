package database

import (
	"database/sql"
)

// SmartExecutor automatically routes SELECT queries to Reader and write queries to Writer
type SmartExecutor struct {
	writer *sql.DB
	reader *sql.DB
}

// NewSmartExecutor creates a new smart executor that routes queries optimally
func NewSmartExecutor(writer *sql.DB, reader *sql.DB) *SmartExecutor {
	return &SmartExecutor{
		writer: writer,
		reader: reader,
	}
}

// Query routes SELECT queries to Reader
func (se *SmartExecutor) Query(query string, args ...interface{}) (*sql.Rows, error) {
	// Route read queries to Reader for parallelism
	return se.reader.Query(query, args...)
}

// QueryRow routes SELECT queries to Reader
func (se *SmartExecutor) QueryRow(query string, args ...interface{}) *sql.Row {
	// Route read queries to Reader for parallelism
	return se.reader.QueryRow(query, args...)
}

// Exec routes write operations to Writer
func (se *SmartExecutor) Exec(query string, args ...interface{}) (sql.Result, error) {
	// All writes go to Writer for serialization and safety
	return se.writer.Exec(query, args...)
}

// Begin creates a new transaction on Writer (for serialization)
func (se *SmartExecutor) Begin() (*sql.Tx, error) {
	// Transactions always use Writer for safety and serialization
	return se.writer.Begin()
}
