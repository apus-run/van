package stringutils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ToString Change arg to string
func ToString(arg any, timeFormat ...string) string {
	switch v := arg.(type) {
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.Itoa(int(v))
	case uint8:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatInt(int64(v), 10)
	case string:
		return v
	case []byte:
		return string(v)
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		if len(timeFormat) > 0 {
			return v.Format(timeFormat[0])
		}
		return v.Format("2006-01-02 15:04:05")
	case reflect.Value:
		return ToString(v.Interface(), timeFormat...)
	case fmt.Stringer:
		return v.String()
	default:
		// Check if the type is a pointer by using reflection
		rv := reflect.ValueOf(arg)
		if rv.Kind() == reflect.Ptr && !rv.IsNil() {
			// Dereference the pointer and recursively call ToString
			return ToString(rv.Elem().Interface(), timeFormat...)
		} else if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			// handle slices
			var buf strings.Builder
			buf.WriteString("[") //nolint: revive,errcheck // no need to check error
			for i := 0; i < rv.Len(); i++ {
				if i > 0 {
					buf.WriteString(" ") //nolint: revive,errcheck // no need to check error
				}
				buf.WriteString(ToString(rv.Index(i).Interface())) //nolint: revive,errcheck // no need to check error
			}
			buf.WriteString("]") //nolint: revive,errcheck // no need to check error
			return buf.String()
		}

		// For types not explicitly handled, use fmt.Sprint to generate a string representation
		return fmt.Sprint(arg)
	}
}
