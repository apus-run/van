package sliceutils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type randomTestStruct struct {
	blah string
}

func testMap[T, U any](t *testing.T, desc string, slice []T, mapFunc func(T) U, expectedOut []U) {
	t.Run(desc, func(t *testing.T) {
		out := Map(slice, mapFunc)
		require.Equal(t, expectedOut, out)
	})
}

func TestMap(t *testing.T) {
	testMap(t,
		"simple func on string",
		[]string{"1", "2"},
		func(s string) string {
			return s + s
		},
		[]string{"11", "22"},
	)
	testMap(t,
		"func of int to string array",
		[]int{1, 2},
		func(val int) []string {
			out := make([]string, 0, val)
			for i := 0; i < val; i++ {
				out = append(out, strconv.Itoa(i))
			}
			return out
		},
		[][]string{{"0"}, {"0", "1"}},
	)
	testMap(t,
		"extract element from struct",
		[]randomTestStruct{{"1"}, {"2"}},
		func(s randomTestStruct) string {
			return s.blah
		},
		[]string{"1", "2"},
	)
}

func TestToMap(t *testing.T) {
	type Foo struct {
		ID   int
		Name string
	}
	mapper := func(f Foo) (int, string) { return f.ID, f.Name }
	assert.Equal(t, map[int]string{}, ToMap([]Foo{}, mapper))
	assert.Equal(t, map[int]string{}, ToMap(nil, mapper))
	assert.Equal(t,
		map[int]string{1: "one", 2: "two", 3: "three"},
		ToMap([]Foo{{1, "one"}, {2, "two"}, {3, "three"}}, mapper))
}
