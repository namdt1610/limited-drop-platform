/**
 * Go Generics Utilities
 * Type-safe, reusable functions for service layer
 */

package utils

import (
	"gorm.io/gorm"
)

// PaginatedResult - Generic pagination response
type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

// Paginate - Generic pagination helper
// Usage: Paginate[Product](db, page, limit, &products)
func Paginate[T any](db *gorm.DB, page, limit int, dest *[]T) (*PaginatedResult[T], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var total int64
	if err := db.Model(new(T)).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * limit
	if err := db.Offset(offset).Limit(limit).Find(dest).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &PaginatedResult[T]{
		Items:      *dest,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// FindByID - Generic find by ID
func FindByID[T any](db *gorm.DB, id uint) (*T, error) {
	var result T
	if err := db.First(&result, id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// FindAll - Generic find all with optional preloads
func FindAll[T any](db *gorm.DB, preloads ...string) ([]T, error) {
	var results []T
	query := db
	for _, p := range preloads {
		query = query.Preload(p)
	}
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// Create - Generic create
func Create[T any](db *gorm.DB, item *T) error {
	return db.Create(item).Error
}

// Update - Generic update
func Update[T any](db *gorm.DB, item *T) error {
	return db.Save(item).Error
}

// Delete - Generic soft delete by ID
func Delete[T any](db *gorm.DB, id uint) error {
	return db.Delete(new(T), id).Error
}

// Filter - Generic filter helper
type FilterOption func(*gorm.DB) *gorm.DB

func WithWhere(condition string, args ...interface{}) FilterOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(condition, args...)
	}
}

func WithOrder(order string) FilterOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(order)
	}
}

func WithPreload(preload string) FilterOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(preload)
	}
}

// FindWithFilters - Generic find with filters
func FindWithFilters[T any](db *gorm.DB, opts ...FilterOption) ([]T, error) {
	var results []T
	query := db.Model(new(T))
	for _, opt := range opts {
		query = opt(query)
	}
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// Map - Transform slice of T to slice of R
func Map[T, R any](items []T, fn func(T) R) []R {
	result := make([]R, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}

// Filter - Filter slice by predicate
func Filter[T any](items []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Ptr - Return pointer to value (useful for optional fields)
func Ptr[T any](v T) *T {
	return &v
}
