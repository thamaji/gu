package slices

import (
	"github.com/thamaji/gu/iter"
	"github.com/thamaji/gu/tuple"
)

// 値からスライスをつくる。
func From[T any](values ...T) []T {
	return append(make([]T, 0, len(values)), values...)
}

// ポインタからスライスをつくる。
func FromPtr[T any](v *T) []T {
	if v == nil {
		return []T{}
	}
	return []T{*v}
}

// マップからスライスをつくる。
func FromMap[K comparable, V any](m map[K]V) []tuple.T2[K, V] {
	entries := make([]tuple.T2[K, V], 0, len(m))
	for key, value := range m {
		entries = append(entries, tuple.NewT2(key, value))
	}
	return entries
}

// マップのキーをスライスをつくる。
func FromMapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// マップの値をスライスをつくる。
func FromMapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

// イテレータからスライスをつくる。
func FromIter[V any](iter iter.Iter[V]) []V {
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

// 関数をn回実行した結果からスライスをつくる。
func FromFunc[V any](n int, f func(int) V) []V {
	slice := make([]V, 0, n)
	for i := 0; i < n; i++ {
		slice = append(slice, f(i))
	}
	return slice
}

// 値と値の有無を受け取ってスライスをつくる。
func Option[V any](v V, ok bool) []V {
	if ok {
		return []V{v}
	}
	return []V{}
}
