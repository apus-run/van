package stringutils

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// MyStringer for testing fmt.Stringer support.
type MyStringer struct {
	value string
}

// CustomType not handled by switch cases.
type CustomType struct {
	Field1 string
	Field2 int
}

func (ms MyStringer) String() string {
	return ms.value
}

func Test_ToString(t *testing.T) {
	t.Parallel()

	customInstance := CustomType{Field1: "example", Field2: 42}
	timeSample := time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC)
	stringerSample := MyStringer{value: "Stringer Value"}

	tests := []struct {
		name       string
		input      any
		timeFormat []string
		expected   string
	}{
		// Primitive types and string
		{name: "int", input: int(42), expected: "42"},
		{name: "int8", input: int8(42), expected: "42"},
		{name: "int16", input: int16(42), expected: "42"},
		{name: "int32", input: int32(42), expected: "42"},
		{name: "int64", input: int64(42), expected: "42"},
		{name: "uint", input: uint(100), expected: "100"},
		{name: "uint8", input: uint8(100), expected: "100"},
		{name: "uint16", input: uint16(100), expected: "100"},
		{name: "uint32", input: uint32(100), expected: "100"},
		{name: "uint64", input: uint64(100), expected: "100"},
		{name: "string", input: "test string", expected: "test string"},
		{name: "[]byte", input: []byte("Hello, World!"), expected: "Hello, World!"},
		{name: "bool", input: true, expected: "true"},
		{name: "float32", input: float32(3.14), expected: "3.14"},
		{name: "float64", input: float64(3.14159), expected: "3.14159"},

		// time.Time
		{name: "time.Time default format", input: timeSample, expected: "2000-01-01 12:34:56"},
		{name: "time.Time custom format", input: timeSample, timeFormat: []string{"Jan 02, 2006"}, expected: "Jan 01, 2000"},

		// reflect.Value
		{name: "reflect.Value", input: reflect.ValueOf(42), expected: "42"},

		// fmt.Stringer
		{name: "fmt.Stringer", input: stringerSample, expected: "Stringer Value"},

		// Composite types (arrays, slices)
		{name: "[]string", input: []string{"Hello", "World"}, expected: "[Hello World]"},
		{name: "[]int", input: []int{42, 21}, expected: "[42 21]"},
		{name: "[][]int", input: [][]int{{42, 21}, {42, 21}}, expected: "[[42 21] [42 21]]"},
		{name: "[]any", input: []any{[]int{42, 21}, 42, "Hello", true, []string{"Hello", "World"}}, expected: "[[42 21] 42 Hello true [Hello World]]"},

		// Custom unhandled type
		{name: "CustomType", input: customInstance, expected: "{example 42}"},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var res string
			if len(tc.timeFormat) > 0 {
				res = ToString(tc.input, tc.timeFormat...)
			} else {
				res = ToString(tc.input)
			}
			require.Equal(t, tc.expected, res)
		})
	}

	// Testing pointer to int
	intPtr := 42
	testsPtr := []struct {
		input    any
		expected string
	}{
		{&intPtr, "42"},
	}
	for _, tc := range testsPtr {
		tc := tc
		t.Run("pointer to "+reflect.TypeOf(tc.input).Elem().String(), func(t *testing.T) {
			t.Parallel()
			res := ToString(tc.input)
			require.Equal(t, tc.expected, res)
		})
	}
}
