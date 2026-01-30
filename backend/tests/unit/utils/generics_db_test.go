package utils_test

import (
	"testing"

	"ecommerce-backend/internal/utils"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestModel for generic DB tests
type TestModel struct {
	ID   uint `gorm:"primaryKey"`
	Name string
	Age  int
}

func setupTestDB(t *testing.T) *gorm.DB {
	// Use t.Name() to ensure unique DB per test case
	dsn := "file:" + t.Name() + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&TestModel{})
	return db
}

func TestGenerics_Create_Find_Update_Delete(t *testing.T) {
	db := setupTestDB(t)

	// Test Create
	item := &TestModel{Name: "Alice", Age: 30}
	err := utils.Create(db, item)
	assert.NoError(t, err)
	assert.NotZero(t, item.ID)

	// Test FindByID
	found, err := utils.FindByID[TestModel](db, item.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", found.Name)

	// Test Update
	item.Age = 31
	err = utils.Update(db, item)
	assert.NoError(t, err)

	updated, _ := utils.FindByID[TestModel](db, item.ID)
	assert.Equal(t, 31, updated.Age)

	// Test FindAll
	utils.Create(db, &TestModel{Name: "Bob", Age: 25})
	all, err := utils.FindAll[TestModel](db)
	assert.NoError(t, err)
	assert.Len(t, all, 2)

	// Test Delete
	err = utils.Delete[TestModel](db, item.ID)
	assert.NoError(t, err)

	allAfter, _ := utils.FindAll[TestModel](db)
	assert.Len(t, allAfter, 1)
	assert.Equal(t, "Bob", allAfter[0].Name)
}

func TestGenerics_Pagination_Filters(t *testing.T) {
	db := setupTestDB(t)
	utils.Create(db, &TestModel{Name: "A", Age: 10})
	utils.Create(db, &TestModel{Name: "B", Age: 20})
	utils.Create(db, &TestModel{Name: "C", Age: 30})
	utils.Create(db, &TestModel{Name: "D", Age: 40})
	utils.Create(db, &TestModel{Name: "E", Age: 50})

	// Test Paginate
	var results []TestModel
	paginated, err := utils.Paginate[TestModel](db, 2, 2, &results) // Page 2, Size 2 (C, D)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), paginated.Total)
	assert.Len(t, paginated.Items, 2)
	assert.Equal(t, "C", paginated.Items[0].Name)
	assert.Equal(t, "D", paginated.Items[1].Name)

	// Test FindWithFilters (Where)
	oldies, err := utils.FindWithFilters[TestModel](db, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("age > ?", 30)
	})
	assert.NoError(t, err)
	assert.Len(t, oldies, 2) // D, E
	assert.Equal(t, "D", oldies[0].Name)
}

func TestGenerics_Scopes(t *testing.T) {
	db := setupTestDB(t) 
	// Just verify scopes return correct gorm instance wrapper
	tx := db.Model(&TestModel{})
	
	scope := utils.WithWhere("age > ?", 10)
	tx2 := scope(tx)
	assert.NotNil(t, tx2)
	
	scopeOrder := utils.WithOrder("age desc")
	tx3 := scopeOrder(tx)
	assert.NotNil(t, tx3)
}

func TestGenerics_WithPreload(t *testing.T) {
	db := setupTestDB(t)
	tx := db.Model(&TestModel{})

	scope := utils.WithPreload("SomeRelation")
	tx2 := scope(tx)
	assert.NotNil(t, tx2)
}
