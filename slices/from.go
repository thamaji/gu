package slices

import (
	"github.com/thamaji/gu/iter"
	"github.com/thamaji/gu/tuple"
	"golang.org/x/exp/constraints"
)

// 範囲を指定してスライスをつくる。
func Range[V constraints.Integer | constraints.Float](start V, stop V, step V) []V {
	slice := make([]V, 0, int((stop-start)/step))
	for cursor := start; cursor < stop; cursor += step {
		slice = append(slice, cursor)
	}
	return slice
}

// 指定した値をn個複製してスライスをつくる。
func Repeat[V any](n int, v V) []V {
	slice := make([]V, n)
	for i := 0; i < n; i++ {
		slice[i] = v
	}
	return slice
}

// 指定した値をn個複製してスライスをつくる。
func RepeatBy[V any](n int, f func(int) (V, error)) ([]V, error) {
	slice := make([]V, n)
	for i := 0; i < n; i++ {
		v, err := f(i)
		if err != nil {
			return nil, err
		}
		slice[i] = v
	}
	return slice, nil
}

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

// 値と値の有無を受け取ってスライスをつくる。
func Option[V any](v V, ok bool) []V {
	if ok {
		return []V{v}
	}
	return []V{}
}
