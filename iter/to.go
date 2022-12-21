package iter

import "github.com/thamaji/gu/tuple"

// イテレータからスライスをつくる。
func ToSlice[V any](iter Iter[V]) ([]V, error) {
	slice := []V{}
	for {
		v, ok := iter.Next()
		if !ok {
			if err := iter.Err(); err != nil {
				return nil, err
			}
			break
		}
		slice = append(slice, v)
	}
	return slice, nil
}

// イテレータからマップをつくる。
func ToMap[K comparable, V any](iter Iter[tuple.T2[K, V]]) (map[K]V, error) {
	m := map[K]V{}
	for {
		v, ok := iter.Next()
		if !ok {
			if err := iter.Err(); err != nil {
				return nil, err
			}
			break
		}
		m[v.V1] = v.V2
	}
	return m, nil
}

// イテレータからポインタをつくる。
// ふたつ目以降の値は無視される。
func ToPtr[V any](iter Iter[V]) (*V, error) {
	for {
		v, ok := iter.Next()
		if !ok {
			return nil, iter.Err()
		}
		return &v, nil
	}
}
