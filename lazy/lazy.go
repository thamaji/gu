package lazy

import "sync"

type Value[V any] interface {
	Get() (V, error)
}

func New[V any](f func() (V, error)) Value[V] {
	return &value[V]{get: f}
}

type value[V any] struct {
	get  func() (V, error)
	once sync.Once
	v    V
	err  error
}

func (v *value[V]) Get() (V, error) {
	v.once.Do(func() {
		v.v, v.err = v.get()
	})
	return v.v, v.err
}
