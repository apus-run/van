package timeutils

import (
	"time"

	"github.com/apus-run/van/pkg/utils"
)

// MustParse parses the given value into a `time.Time` according to the layout, or panics if there is a parse error.
func MustParse(layout string, value string) time.Time {
	ts, err := time.Parse(layout, value)
	utils.CrashOnError(err)
	return ts
}
