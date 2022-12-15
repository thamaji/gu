package maps

import (
	"github.com/thamaji/gu/iter"
	"github.com/thamaji/gu/tuple"
)

// マップからキーのスライスをつくる。
func Keys[K comparable, V any](m map[K]V) []K {
	dst := make([]K, 0, len(m))
	for k := range m {
		dst = append(dst, k)
	}
	return dst
}

// マップから値のスライスをつくる。
func Values[K comparable, V any](m map[K]V) []V {
	dst := make([]V, 0, len(m))
	for _, v := range m {
		dst = append(dst, v)
	}
	return dst
}

// マップからスライスをつくる。
func ToSlice[K comparable, V any](m map[K]V) []tuple.T2[K, V] {
	dst := make([]tuple.T2[K, V], 0, len(m))
	for k, v := range m {
		dst = append(dst, tuple.NewT2(k, v))
	}
	return dst
}

// マップからイテレータをつくる。
func ToIter[K comparable, V any](m map[K]V) iter.Iter[tuple.T2[K, V]] {
	return iter.FromMap(m)
}
