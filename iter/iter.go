package iter

type Iter[V any] interface {
	// 次の値を返す。
	// 第２戻り値は値があるときは true、なければ false を返す。
	Next() (V, bool)
	Err() error
}

type IterFunc[V any] func() (V, bool)

func (f IterFunc[V]) Next() (V, bool) {
	return f()
}

func (f IterFunc[V]) Err() error {
	return nil
}

func Empty[V any]() Iter[V] {
	return IterFunc[V](func() (V, bool) {
		return *new(V), false
	})
}
