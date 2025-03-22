package safe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicErr(t *testing.T) {
	// defer func() {
	// 	panicInfo := recover()
	// 	if panicInfo != nil {
	// 		err := NewPanicErr(panicInfo, debug.Stack())
	// 		panic(err)
	// 	}
	// }()

	err := NewPanicErr("info", []byte("stack"))
	assert.Equal(t, "panic error: info, \nstack: stack", err.Error())

}
