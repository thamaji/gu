package ptr

import (
	"github.com/thamaji/gu/iter"
)

// ポインタからスライスをつくる。
func ToSlice[V any](p *V) []V {
	if p == nil {
		return []V{}
	}
	return []V{*p}
}

// ポインタからイテレータをつくる。
func ToIter[V any](p *V) iter.Iter[V] {
	return iter.FromPtr(p)
}
