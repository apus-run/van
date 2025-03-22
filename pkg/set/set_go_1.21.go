//go:build go1.21

package set

// Clear empties the set.
// It is preferable to replace the set with a newly constructed set,
// but not all callers can do that (when there are other references to the map).
func (s Set[T]) Clear() Set[T] {
	clear(s)
	return s
}
