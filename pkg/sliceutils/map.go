package sliceutils

// Map maps the elements of slice, using the given mapFunc
// Example usage:
//
//	Map([]string{"a", "b", "cd"}, func(s string) int {
//	  return len(s)
//	})
//
// will return []int{1, 1, 2}.
func Map[T, U any](slice []T, mapFunc func(T) U) []U {
	result := make([]U, 0, len(slice))
	for _, elem := range slice {
		result = append(result, mapFunc(elem))
	}
	return result
}

// ToMap collects elements of slice to map, both map keys and values are produced
// by mapping function f.
//
// 🚀 EXAMPLE:
//
//	type Foo struct {
//		ID   int
//		Name string
//	}
//	mapper := func(f Foo) (int, string) { return f.ID, f.Name }
//	ToMap([]Foo{}, mapper) ⏩ map[int]string{}
//	s := []Foo{{1, "one"}, {2, "two"}, {3, "three"}}
//	ToMap(s, mapper)       ⏩ map[int]string{1: "one", 2: "two", 3: "three"}
func ToMap[T, V any, K comparable](s []T, f func(T) (K, V)) map[K]V {
	m := make(map[K]V, len(s))
	for _, e := range s {
		k, v := f(e)
		m[k] = v
	}
	return m
}
