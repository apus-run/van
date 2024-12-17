package sql_helper

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ValuesPlaceholders returns a set of SQL placeholder numbers grouped for use in an INSERT
// statement. For example, ValuesPlaceholders(2,3) returns ($1, $2), ($3, $4), ($5, $6)
// It panics if either param is <= 0.
func ValuesPlaceholders(valuesPerRow, numRows int) string {
	if valuesPerRow <= 0 || numRows <= 0 {
		panic("Cannot make ValuesPlaceholder with 0 rows or 0 values per row")
	}
	values := strings.Builder{}
	// There are at most 5 bytes per value that need to be written
	values.Grow(5 * valuesPerRow * numRows)
	// All WriteString calls below return nil errors, as specified in the documentation of
	// strings.Builder, so it is safe to ignore them.
	for argIdx := 1; argIdx <= valuesPerRow*numRows; argIdx += valuesPerRow {
		if argIdx != 1 {
			_, _ = values.WriteString(",")
		}
		_, _ = values.WriteString("(")
		for i := 0; i < valuesPerRow; i++ {
			if i != 0 {
				_, _ = values.WriteString(",")
			}
			_, _ = values.WriteString("$")
			_, _ = values.WriteString(strconv.Itoa(argIdx + i))
		}
		_, _ = values.WriteString(")")
	}
	return values.String()
}

//这一个系列的方法，会在数据为零值时，将valid设置为false；否则设置为true；

func NewNullString(val string) sql.NullString {
	return sql.NullString{String: val, Valid: val != ""}
}

func NewNullInt64(val int64) sql.NullInt64 {
	return sql.NullInt64{Int64: val, Valid: val != 0}
}

func NewNullFloat64(val float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: val, Valid: val != 0}
}

func NewNullBool(val bool) sql.NullBool {
	return sql.NullBool{Bool: val, Valid: val}
}

func NewNullTime(val time.Time) sql.NullTime {
	return sql.NullTime{Time: val, Valid: !val.IsZero()}
}

func NewNullBytes(val []byte) sql.NullString {
	return sql.NullString{String: string(val), Valid: len(val) > 0}
}

// JsonColumn 代表存储字段的 json 类型
// 主要用于没有提供默认 json 类型的数据库
// T 可以是结构体，也可以是切片或者 map
// 理论上来说一切可以被 json 库所处理的类型都能被用作 T
// 不建议使用指针作为 T 的类型
// 如果 T 是指针，那么在 Val 为 nil 的情况下，一定要把 Valid 设置为 false
type JsonColumn[T any] struct {
	Val   T
	Valid bool
}

// Value 返回一个 json 串。类型是 []byte
func (j JsonColumn[T]) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	res, err := json.Marshal(j.Val)
	return res, err
}

// Scan 将 src 转化为对象
// src 的类型必须是 []byte, string 或者 nil
// 如果是 nil，我们不会做任何处理
func (j *JsonColumn[T]) Scan(src any) error {
	var bs []byte
	switch val := src.(type) {
	case nil:
		return nil
	case []byte:
		bs = val
	case string:
		bs = []byte(val)
	default:
		return fmt.Errorf("ekit：JsonColumn.Scan 不支持 src 类型 %v", src)
	}

	if err := json.Unmarshal(bs, &j.Val); err != nil {
		return err
	}
	j.Valid = true
	return nil
}
