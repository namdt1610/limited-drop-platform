package utils

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// tests moved to backend/tests/unit/utils/generics_test.go
func TestCreate(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Mock transaction begin + create query (GORM may use transactions and RETURNING id)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "test_models"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	type TestModel struct {
		ID   uint   `gorm:"primarykey"`
		Name string `gorm:"type:varchar(100)"`
	}

	item := &TestModel{Name: "test"}
	err = Create(gormDB, item)

	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Mock transaction begin + update query (GORM may use transactions)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "test_models" SET`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	type TestModel struct {
		ID   uint   `gorm:"primarykey"`
		Name string `gorm:"type:varchar(100)"`
	}

	item := &TestModel{ID: 1, Name: "updated"}
	err = Update(gormDB, item)

	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Mock transaction begin + delete query (GORM may use transactions)
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "test_models" WHERE .*`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	type TestModel struct {
		ID uint `gorm:"primarykey"`
	}

	err = Delete[TestModel](gormDB, 1)

	assert.NoError(t, err)
}

func TestFindWithFilters(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Mock the query with filters
	mock.ExpectQuery(`SELECT \* FROM "test_models" WHERE .* ORDER BY .*`).
		WithArgs("test").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test"))

	type TestModel struct {
		ID   uint   `gorm:"primarykey"`
		Name string `gorm:"type:varchar(100)"`
	}

	results, err := FindWithFilters[TestModel](gormDB,
		WithWhere("name = ?", "test"),
		WithOrder("created_at"))

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "test", results[0].Name)
}
