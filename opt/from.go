package opt

// 値と値の有無を受け取って Option をつくる。
func From[V any](v V, ok bool) Option[V] {
	if !ok {
		return None[V]()
	}
	return Some(v)
}

// ポインタから Option をつくる。
func FromPtr[V any](v *V) Option[V] {
	if v == nil {
		return None[V]()
	}
	return Some(*v)
}
