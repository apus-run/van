package sql_helper

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValuesPlaceholders_ValidInputs_Success(t *testing.T) {

	v := ValuesPlaceholders(3, 2)
	assert.Equal(t, "($1,$2,$3),($4,$5,$6)", v)

	v = ValuesPlaceholders(2, 4)
	assert.Equal(t, "($1,$2),($3,$4),($5,$6),($7,$8)", v)

	v = ValuesPlaceholders(1, 1)
	assert.Equal(t, "($1)", v)

	v = ValuesPlaceholders(1, 3)
	assert.Equal(t, "($1),($2),($3)", v)
}

func TestValuesPlaceholders_InvalidInputs_Panics(t *testing.T) {

	assert.Panics(t, func() {
		ValuesPlaceholders(-3, 2)
	})
	assert.Panics(t, func() {
		ValuesPlaceholders(2, -4)
	})
	assert.Panics(t, func() {
		ValuesPlaceholders(0, 0)
	})
}

func TestNewNullBool(t *testing.T) {
	tests := []struct {
		name string
		val  bool
		want sql.NullBool
	}{
		{
			name: "nonzero",
			val:  true,
			want: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		},
		{
			name: "zero",
			val:  false,
			want: sql.NullBool{
				Bool:  false,
				Valid: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullBool(tt.val), "NewNullBool(%v)", tt.val)
		})
	}
}

func TestNewNullBytes(t *testing.T) {
	tests := []struct {
		name string
		val  []byte
		want sql.NullString
	}{
		{
			name: "nonzero",
			val:  []byte("test"),
			want: sql.NullString{
				String: "test",
				Valid:  true,
			},
		},
		{
			name: "zero",
			val:  []byte{},
			want: sql.NullString{
				String: "",
				Valid:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullBytes(tt.val), "NewNullBytes(%v)", tt.val)
		})
	}
}

func TestNewNullFloat64(t *testing.T) {
	tests := []struct {
		name string
		val  float64
		want sql.NullFloat64
	}{
		{
			name: "nonzero",
			val:  1.1,
			want: sql.NullFloat64{
				Float64: 1.1,
				Valid:   true,
			},
		},
		{
			name: "zero",
			val:  0,
			want: sql.NullFloat64{
				Float64: 0,
				Valid:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullFloat64(tt.val), "NewNullFloat64(%v)", tt.val)
		})
	}
}

func TestNewNullInt64(t *testing.T) {
	tests := []struct {
		name string
		val  int64
		want sql.NullInt64
	}{
		{
			name: "nonzero",
			val:  1,
			want: sql.NullInt64{
				Int64: 1,
				Valid: true,
			},
		},
		{
			name: "zero",
			val:  0,
			want: sql.NullInt64{
				Int64: 0,
				Valid: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullInt64(tt.val), "NewNullInt64(%v)", tt.val)
		})
	}
}

func TestNewNullString(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want sql.NullString
	}{
		{
			name: "nonzero",
			val:  "test",
			want: sql.NullString{
				String: "test",
				Valid:  true,
			},
		},
		{
			name: "zero",
			val:  "",
			want: sql.NullString{
				String: "",
				Valid:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullString(tt.val), "NewNullString(%v)", tt.val)
		})
	}
}

func TestNewNullTime(t *testing.T) {
	tests := []struct {
		name string
		val  time.Time
		want sql.NullTime
	}{
		{
			name: "nonzero",
			val:  time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
			want: sql.NullTime{
				Time:  time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
				Valid: true,
			},
		},
		{
			name: "zero",
			val:  time.Time{},
			want: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullTime(tt.val), "NewNullTime(%v)", tt.val)
		})
	}
}

func TestJsonColumn_Value(t *testing.T) {
	testCases := []struct {
		name    string
		valuer  driver.Valuer
		wantRes any
		wantErr error
	}{
		{
			name:    "user",
			valuer:  JsonColumn[User]{Valid: true, Val: User{Name: "Tom"}},
			wantRes: []byte(`{"Name":"Tom"}`),
		},
		{
			name:   "invalid",
			valuer: JsonColumn[User]{},
		},
		{
			name:   "nil",
			valuer: JsonColumn[*User]{},
		},
		{
			name:    "nil but valid",
			valuer:  JsonColumn[*User]{Valid: true},
			wantRes: []uint8("null"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := tc.valuer.Value()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, value)
		})
	}
}

func TestJsonColumn_Scan(t *testing.T) {
	testCases := []struct {
		name      string
		src       any
		wantErr   error
		wantValid bool
		wantVal   User
	}{
		{
			name:    "nil",
			wantVal: User{},
		},
		{
			name:      "string",
			src:       `{"Name":"Tom"}`,
			wantVal:   User{Name: "Tom"},
			wantValid: true,
		},
		{
			name:      "bytes",
			src:       []byte(`{"Name":"Tom"}`),
			wantVal:   User{Name: "Tom"},
			wantValid: true,
		},
		{
			name:    "int",
			src:     123,
			wantErr: errors.New("ekit：JsonColumn.Scan 不支持 src 类型 123"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			js := &JsonColumn[User]{}
			err := js.Scan(tc.src)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantValid, js.Valid)
			if !js.Valid {
				return
			}
			assert.Equal(t, tc.wantVal, js.Val)
		})
	}
}

func TestJsonColumn_ScanTypes(t *testing.T) {
	jsSlice := JsonColumn[[]string]{}
	err := jsSlice.Scan(`["a", "b", "c"]`)
	assert.Nil(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, jsSlice.Val)
	val, err := jsSlice.Value()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`["a","b","c"]`), val)

	jsMap := JsonColumn[map[string]string]{}
	err = jsMap.Scan(`{"a":"a value"}`)
	assert.Nil(t, err)
	val, err = jsMap.Value()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"a":"a value"}`), val)
}

type User struct {
	Name string
}

func ExampleJsonColumn_Value() {
	js := JsonColumn[User]{Valid: true, Val: User{Name: "Tom"}}
	value, err := js.Value()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(value.([]byte)))
	// Output:
	// {"Name":"Tom"}
}

func ExampleJsonColumn_Scan() {
	js := JsonColumn[User]{}
	err := js.Scan(`{"Name":"Tom"}`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(js.Val)
	// Output:
	// {Tom}
}
