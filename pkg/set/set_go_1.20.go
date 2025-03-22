//go:build !go1.21

package set

// Clear empties the set.
// It is preferable to replace the set with a newly constructed set,
// but not all callers can do that (when there are other references to the map).
// In some cases the set *won't* be fully cleared, e.g. a Set[float32] containing NaN
// can't be cleared because NaN can't be removed.
// For sets containing items of a type that is reflexive for ==,
// this is optimized to a single call to runtime.mapclear().
func (s Set[T]) Clear() Set[T] {
	for key := range s {
		delete(s, key)
	}
	return s
}
