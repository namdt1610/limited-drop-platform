package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	utils "ecommerce-backend/internal/utils"
)

func TestPtr(t *testing.T) {
	cases := []struct {
		name  string
		input any
		want  any
	}{
		{"string", "test", "test"},
		{"int", 42, 42},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			switch v := tc.input.(type) {
			case string:
				got := utils.Ptr(v)
				if got == nil {
					t.Fatalf("expected non-nil pointer for %s", tc.name)
				}
				if *got != tc.want.(string) {
					t.Fatalf("expected %v, got %v", tc.want, *got)
				}
			case int:
				got := utils.Ptr(v)
				if got == nil {
					t.Fatalf("expected non-nil pointer for %s", tc.name)
				}
				if *got != tc.want.(int) {
					t.Fatalf("expected %v, got %v", tc.want, *got)
				}
			default:
				t.Fatalf("unsupported type %T", v)
			}
		})
	}
}

func TestMap(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"strings to lengths", func(t *testing.T) {
			input := []string{"a", "bb", "ccc"}
			got := utils.Map(input, func(s string) int { return len(s) })
			want := []int{1, 2, 3}
			if !assert.Equal(t, want, got) {
				t.FailNow()
			}
		}},
		{"numbers to squares", func(t *testing.T) {
			input := []int{1, 2, 3, 4}
			got := utils.Map(input, func(n int) int { return n * n })
			want := []int{1, 4, 9, 16}
			if !assert.Equal(t, want, got) {
				t.FailNow()
			}
		}},
		{"empty slice", func(t *testing.T) {
			input := []string{}
			got := utils.Map(input, func(s string) string { return s + "!" })
			if !assert.Len(t, got, 0) {
				t.FailNow()
			}
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, tc.run)
	}
}

func TestFilter(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"filter evens", func(t *testing.T) {
			input := []int{1, 2, 3, 4, 5, 6}
			got := utils.Filter(input, func(n int) bool { return n%2 == 0 })
			want := []int{2, 4, 6}
			if !assert.Equal(t, want, got) {
				t.FailNow()
			}
		}},
		{"filter long strings", func(t *testing.T) {
			input := []string{"a", "bb", "ccc", "dddd"}
			got := utils.Filter(input, func(s string) bool { return len(s) > 2 })
			want := []string{"ccc", "dddd"}
			if !assert.Equal(t, want, got) {
				t.FailNow()
			}
		}},
		{"empty input", func(t *testing.T) {
			input := []int{}
			got := utils.Filter(input, func(n int) bool { return true })
			if !assert.Len(t, got, 0) {
				t.FailNow()
			}
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, tc.run)
	}
}

// Table-driven conversion of DB tests can be done later if you want.
// Leave existing DB tests as-is (kept in older file) or move them here when ready.

// Table-driven conversion of DB tests can be done later if you want.// Leave existing DB tests as-is (kept in older file) or move them here when ready.}	}		t.Run(tc.name, tc.run)	for _, tc := range cases {	}		}},			if !assert.Len(t, got, 0) { t.FailNow() }			got := utils.Filter(input, func(n int) bool { return true })			input := []int{}		{"empty input", func(t *testing.T) {		}},			if !assert.Equal(t, want, got) { t.FailNow() }			want := []string{"ccc","dddd"}			got := utils.Filter(input, func(s string) bool { return len(s) > 2 })			input := []string{"a","bb","ccc","dddd"}		{"filter long strings", func(t *testing.T) {		}},			if !assert.Equal(t, want, got) { t.FailNow() }			want := []int{2,4,6}			got := utils.Filter(input, func(n int) bool { return n%2 == 0 })			input := []int{1,2,3,4,5,6}		{"filter evens", func(t *testing.T) {	}{		run func(t *testing.T)		name string	cases := []struct {func TestFilter(t *testing.T) {}	}		t.Run(tc.name, tc.run)	for _, tc := range cases {	}		}},			if !assert.Len(t, got, 0) { t.FailNow() }			got := utils.Map(input, func(s string) string { return s + "!" })			input := []string{}		{"empty slice", func(t *testing.T) {		}},			if !assert.Equal(t, want, got) { t.FailNow() }			want := []int{1,4,9,16}			got := utils.Map(input, func(n int) int { return n*n })			input := []int{1,2,3,4}		{"numbers to squares", func(t *testing.T) {		}},			if !assert.Equal(t, want, got) { t.FailNow() }			want := []int{1,2,3}			got := utils.Map(input, func(s string) int { return len(s) })			input := []string{"a", "bb", "ccc"}		{"strings to lengths", func(t *testing.T) {	}{		run func(t *testing.T)		name string	cases := []struct {func TestMap(t *testing.T) {}	}		})			}				t.Fatalf("unsupported type %T", v)			default:				}					t.Fatalf("expected %v, got %v", tc.want, *got)				if *got != tc.want.(int) {				}					t.Fatalf("expected non-nil pointer for %s", tc.name)				if got == nil {				got := utils.Ptr(v)			case int:				}					t.Fatalf("expected %v, got %v", tc.want, *got)				if *got != tc.want.(string) {				}					t.Fatalf("expected non-nil pointer for %s", tc.name)				if got == nil {				got := utils.Ptr(v)			case string:			switch v := tc.input.(type) {		t.Run(tc.name, func(t *testing.T) {	for _, tc := range cases {	}		{"int", 42, 42},		{"string", "test", "test"},	}{		want any		input any		name string	cases := []struct{func TestPtr(t *testing.T) {)	utils "ecommerce-backend/internal/utils"	"gorm.io/gorm"	"gorm.io/driver/postgres"	"github.com/stretchr/testify/assert"	"github.com/DATA-DOG/go-sqlmock"	"testing"
