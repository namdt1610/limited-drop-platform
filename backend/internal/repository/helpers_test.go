package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestMarshalJSON(t *testing.T) {
	input := map[string]interface{}{"key": "value"}
	bytes, err := marshalJSON(input)
	assert.NoError(t, err)
	assert.Contains(t, string(bytes), `"key":"value"`)
}

func TestJsonToBytes_BytesToJson(t *testing.T) {
	// Test casting helpers
	raw := []byte(`{"key":"value"}`)
	
	// bytesToJSON: []byte -> datatypes.JSON
	dt := bytesToJSON(raw)
	assert.Equal(t, datatypes.JSON(raw), dt)

	// jsonToBytes: datatypes.JSON -> []byte
	b := jsonToBytes(dt)
	assert.Equal(t, raw, b)
}

func TestNullTimeHelpers(t *testing.T) {
	now := time.Now()
	
	// Test nullTimeToPtr
	ntValid := sql.NullTime{Time: now, Valid: true}
	ptr := nullTimeToPtr(ntValid)
	assert.NotNil(t, ptr)
	assert.Equal(t, now, *ptr)

	ntInvalid := sql.NullTime{Valid: false}
	ptrNil := nullTimeToPtr(ntInvalid)
	assert.Nil(t, ptrNil)

	// Test ptrToNullTime
	nt, valid := ptrToNullTime(&now).Value() // sql.NullTime implements Value()
	assert.Nil(t, valid) // error is nil
	assert.Equal(t, now, nt.(time.Time))

	ntNilVal, _ := ptrToNullTime(nil).Value()
	assert.Nil(t, ntNilVal)
}
