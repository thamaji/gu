package iter

import (
	"reflect"

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
	iter := reflect.ValueOf(m).MapRange()
	next := iter.Next()
	return IterFunc[tuple.T2[K, V]](func() (tuple.T2[K, V], bool) {
		if !next {
			return tuple.NewT2(*new(K), *new(V)), false
		}
		k := iter.Key().Interface().(K)
		v := iter.Value().Interface().(V)
		next = iter.Next()
		return tuple.NewT2(k, v), true
	})
}

// マップのキーからイテレータをつくる。
func FromMapKeys[K comparable, V any](m map[K]V) Iter[K] {
	iter := reflect.ValueOf(m).MapRange()
	next := iter.Next()
	return IterFunc[K](func() (K, bool) {
		if !next {
			return *new(K), false
		}
		k := iter.Key().Interface().(K)
		next = iter.Next()
		return k, true
	})
}

// マップの値からイテレータをつくる。
func FromMapValues[K comparable, V any](m map[K]V) Iter[V] {
	iter := reflect.ValueOf(m).MapRange()
	next := iter.Next()
	return IterFunc[V](func() (V, bool) {
		if !next {
			return *new(V), false
		}
		v := iter.Value().Interface().(V)
		next = iter.Next()
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

// 関数からイテレータをつくる。
func FromFunc[V any](f func(Context) (V, bool)) Iter[V] {
	return &customIter[V]{
		next: f,
	}
}

type Context interface {
	SetErr(error)
	Err() error
}

type customIter[V any] struct {
	next func(Context) (V, bool)
	err  error
}

func (iter *customIter[V]) Next() (V, bool) {
	if iter.err != nil {
		return *new(V), false
	}
	return iter.next(iter)
}

func (iter *customIter[V]) SetErr(err error) {
	iter.err = err
}

func (iter *customIter[V]) Err() error {
	return iter.err
}
