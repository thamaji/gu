package slices

import (
	"github.com/thamaji/gu/iter"
	"github.com/thamaji/gu/tuple"
)

// スライスからイテレータをつくる。
func ToIter[V any](slice []V) iter.Iter[V] {
	return iter.FromSlice(slice)
}

// スライスからポインタをつくる。
// ふたつ目以降の値は無視される。
func ToPtr[V any](slice []V) *V {
	for _, v := range slice {
		return &v
	}
	return nil
}

// キーと値のペアのスライスからマップをつくる。
func ToMap[K comparable, V any](slice []tuple.T2[K, V]) map[K]V {
	m := make(map[K]V, len(slice))
	for _, t := range slice {
		m[t.V1] = t.V2
	}
	return m
}
