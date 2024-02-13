package iter

import "github.com/thamaji/gu/tuple"

type Iterable[V any] interface {
	Iter() Iter[V]
}

type Slice[V any] []V

func (slice Slice[V]) Iter() Iter[V] {
	return FromSlice(slice)
}

// XXX: ほんとうは Map にしたいが、func Map と被ってしまう。
type HashMap[K comparable, V any] map[K]V

func (m HashMap[K, V]) Iter() Iter[tuple.T2[K, V]] {
	return FromMap(m)
}

func (m HashMap[K, V]) Keys() Iter[K] {
	return FromMapKeys(m)
}

func (m HashMap[K, V]) Values() Iter[V] {
	return FromMapValues(m)
}
