package sql_helper

import (
	"database/sql"
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
