package ptr

import (
	"github.com/thamaji/gu/iter"
	"github.com/thamaji/gu/tuple"
)

// 値からポインタをつくる。
func From[V any](v V) *V {
	return &v
}

// スライスからポインタをつくる。
// ふたつ目以降の値は無視される。
func FromSlice[V any](slice []V) *V {
	for _, v := range slice {
		return &v
	}
	return nil
}

// マップからポインタをつくる。
// ふたつ目以降の値は無視される。
func FromMap[K comparable, V any](m map[K]V) *tuple.T2[K, V] {
	for k, v := range m {
		return From(tuple.NewT2(k, v))
	}
	return nil
}

// マップのキーからポインタをつくる。
// ふたつ目以降の値は無視される。
func FromMapKeys[K comparable, V any](m map[K]V) *K {
	for k := range m {
		return &k
	}
	return nil
}

// マップの値からポインタをつくる。
// ふたつ目以降の値は無視される。
func FromMapValues[K comparable, V any](m map[K]V) *V {
	for _, v := range m {
		return &v
	}
	return nil
}

// イテレータからポインタをつくる。
// ふたつ目以降の値は無視される。
func FromIter[V any](iter iter.Iter[V]) *V {
	v, ok := iter.Next()
	if !ok {
		return nil
	}
	return &v
}

// 関数を実行した結果からポインタをつくる。
func FromFunc[V any](f func() *V) *V {
	return f()
}

// 値と値の有無を受け取ってポインタをつくる。
func Option[V any](v V, ok bool) *V {
	if ok {
		return &v
	}
	return nil
}
