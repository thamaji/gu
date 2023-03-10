package iter

import (
	"github.com/thamaji/gu/tuple"
	"golang.org/x/exp/constraints"
)

// 範囲を指定してイテレータをつくる。
func Range[V constraints.Ordered](start V, stop V, step V) Iter[V] {
	cursor := start
	return IterFunc[V](func() (V, bool) {
		if cursor >= stop {
			return *new(V), false
		}
		v := cursor
		cursor += step
		return v, true
	})
}

// 複数の値からイテレータをつくる。
func From[V any](values ...V) Iter[V] {
	cursor := 0
	return IterFunc[V](func() (V, bool) {
		if cursor >= len(values) {
			return *new(V), false
		}
		v := values[cursor]
		cursor++
		return v, true
	})
}

// ポインタからイテレータをつくる。
func FromPtr[V any](p *V) Iter[V] {
	cursor := 0
	return IterFunc[V](func() (V, bool) {
		if p == nil || cursor >= 1 {
			return *new(V), false
		}
		v := *p
		cursor++
		return v, true
	})
}

// スライスからイテレータをつくる。
func FromSlice[V any](slice []V) Iter[V] {
	cursor := 0
	return IterFunc[V](func() (V, bool) {
		if cursor >= len(slice) {
			return *new(V), false
		}
		v := slice[cursor]
		cursor++
		return v, true
	})
}

// マップからイテレータをつくる。
func FromMap[K comparable, V any](m map[K]V) Iter[tuple.T2[K, V]] {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	cursor := 0
	return IterFunc[tuple.T2[K, V]](func() (tuple.T2[K, V], bool) {
		if cursor >= len(keys) {
			return tuple.NewT2(*new(K), *new(V)), false
		}
		k := keys[cursor]
		v := m[k]
		cursor++
		return tuple.NewT2(k, v), true
	})
}

// マップのキーからイテレータをつくる。
func FromMapKeys[K comparable, V any](m map[K]V) Iter[K] {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	cursor := 0
	return IterFunc[K](func() (K, bool) {
		if cursor >= len(keys) {
			return *new(K), false
		}
		v := keys[cursor]
		cursor++
		return v, true
	})
}

// マップの値からイテレータをつくる。
func FromMapValues[K comparable, V any](m map[K]V) Iter[V] {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	cursor := 0
	return IterFunc[V](func() (V, bool) {
		if cursor >= len(values) {
			return *new(V), false
		}
		v := values[cursor]
		cursor++
		return v, true
	})
}

// 値と値の有無を受け取ってイテレータをつくる。
func Option[V any](v V, ok bool) Iter[V] {
	if ok {
		return From(v)
	}
	return Empty[V]()
}
