package opt

import (
	"github.com/thamaji/gu/iter"
	"github.com/thamaji/gu/must"
)

func Some[V any](v V) Option[V] {
	return Option[V]{v}
}

func None[V any]() Option[V] {
	return Option[V]{}
}

type Option[V any] []V

// 値の数を返す。
func (opt Option[V]) Len() int {
	return len(opt)
}

// 値があるときtrueを返す。
func (opt Option[V]) IsDefined() bool {
	return len(opt) != 0
}

// 空のときtrueを返す。
func (opt Option[V]) IsEmpty() bool {
	return len(opt) == 0
}

// 値を返す。
func (opt Option[V]) Get() (V, bool) {
	if len(opt) == 0 {
		return *new(V), false
	}
	return opt[0], true
}

// 値を返す。無い場合はvを返す。
func (opt Option[V]) GetOrElse(v V) V {
	if len(opt) == 0 {
		return v
	}
	return opt[0]
}

// 値のポインタを返す。無い場合はnilを返す。
func (opt Option[V]) GetOrNil() *V {
	if len(opt) == 0 {
		return nil
	}
	return &opt[0]
}

// 値を返す。無い場合はゼロ値を返す。
func (opt Option[V]) GetOrZero() V {
	if len(opt) == 0 {
		return *new(V)
	}
	return opt[0]
}

// 値を返す。無い場合は関数の実行結果を返す。
func (opt Option[V]) GetOrFunc(f func() (V, error)) (V, error) {
	if len(opt) == 0 {
		return f()
	}
	return opt[0], nil
}

// 値を返す。無い場合は関数の実行結果を返す。実行中にエラーが起きた場合 panic する。
func (opt Option[V]) MustGetOrFunc(f func() (V, error)) V {
	return must.Must1(opt.GetOrFunc(f))
}

// スライスを返す。
func (opt Option[V]) Slice() []V {
	return opt
}

// イテレータを返す。
func (opt Option[V]) Iter() iter.Iter[V] {
	return iter.FromSlice(opt)
}
