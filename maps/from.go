package maps

import (
	"github.com/thamaji/gu/iter"
	"github.com/thamaji/gu/tuple"
)

// キーと値のペアのスライスからマップをつくる。
func From[K comparable, V any](tuples ...tuple.T2[K, V]) map[K]V {
	m := make(map[K]V, len(tuples))
	for _, t := range tuples {
		m[t.V1] = t.V2
	}
	return m
}

// キーと値のペアのイテレータからマップをつくる。
func FromIter[K comparable, V any](iter iter.Iter[tuple.T2[K, V]]) map[K]V {
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
