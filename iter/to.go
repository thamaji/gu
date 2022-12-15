package iter

import "github.com/thamaji/gu/tuple"

// イテレータからスライスをつくる。
func ToSlice[V any](iter Iter[V]) []V {
	slice := []V{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		slice = append(slice, v)
	}
	return slice
}

// イテレータからマップをつくる。
func ToMap[K comparable, V any](iter Iter[tuple.T2[K, V]]) map[K]V {
	m := map[K]V{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		m[v.V1] = v.V2
	}
	return m
}

// イテレータからポインタをつくる。
// ふたつ目以降の値は無視される。
func ToPtr[V any](iter Iter[V]) *V {
	for {
		v, ok := iter.Next()
		if !ok {
			return nil
		}
		return &v
	}
}
