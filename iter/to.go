package iter

import (
	"github.com/thamaji/gu/must"
	"github.com/thamaji/gu/tuple"
)

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

// イテレータからスライスをつくる。実行中にエラーが起きた場合 panic する。
func MustToSlice[V any](iter Iter[V]) []V {
	return must.Must1(ToSlice(iter))
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

// イテレータからマップをつくる。実行中にエラーが起きた場合 panic する。
func MustToMap[K comparable, V any](iter Iter[tuple.T2[K, V]]) map[K]V {
	return must.Must1(ToMap(iter))
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

// イテレータからポインタをつくる。実行中にエラーが起きた場合 panic する。
// ふたつ目以降の値は無視される。
func MustToPtr[V any](iter Iter[V]) *V {
	return must.Must1(ToPtr(iter))
}
