package iter

import (
	"github.com/thamaji/gu/tuple"
	"golang.org/x/exp/constraints"
)

// 値の数を返す。
func Len[V any](iter Iter[V]) int {
	c := 0
	for {
		if _, ok := iter.Next(); ok {
			break
		}
		c++
	}
	return c
}

// 指定した位置の値を返す。
func Get[V any](iter Iter[V], index int) (V, bool) {
	i := 0
	for {
		v, ok := iter.Next()
		if !ok {
			return *new(V), false
		}
		if i == index {
			return v, true
		}
		i++
	}
}

// 指定した位置の値を返す。無い場合はvを返す。
func GetOrElse[V any](iter Iter[V], index int, v V) V {
	if v, ok := Get(iter, index); ok {
		return v
	}
	return v
}

// 先頭の値を返す。
func GetFirst[V any](iter Iter[V]) (V, bool) {
	return Get(iter, 0)
}

// 先頭の値を返す。無い場合はvを返す。
func GetFirstOrElse[V any](iter Iter[V], v V) V {
	if v, ok := GetFirst(iter); ok {
		return v
	}
	return v
}

// 終端の値を返す。
func GetLast[V any](iter Iter[V]) (V, bool) {
	v := *new(V)
	ok := false
	for {
		v1, ok1 := iter.Next()
		if !ok1 {
			break
		}
		v = v1
		ok = true
	}
	return v, ok
}

// 終端の値を返す。無い場合はvを返す。
func GetLastOrElse[V any](iter Iter[V], v V) V {
	if v, ok := GetLast(iter); ok {
		return v
	}
	return v
}

// イテレータの末尾に他のイテレータを結合する。
func Concat[V any](iter1 Iter[V], iter2 Iter[V]) Iter[V] {
	return IterFunc[V](func() (V, bool) {
		if v, ok := iter1.Next(); ok {
			return v, true
		}
		return iter2.Next()
	})
}

// イテレータの末尾に値を追加する。
func Append[V any](iter Iter[V], v ...V) Iter[V] {
	return Concat(iter, FromSlice(v))
}

// 指定した位置に値を追加する。
func Insert[V any](iter Iter[V], index int, v ...V) Iter[V] {
	i := 0
	return IterFunc[V](func() (V, bool) {
		if i >= index && i < index+len(v) {
			v1 := v[i-index]
			i++
			return v1, true
		}
		v, ok := iter.Next()
		i++
		return v, ok
	})
}

// 指定した位置の値を削除する。
func Remove[V any](iter Iter[V], index int) Iter[V] {
	i := 0
	return IterFunc[V](func() (V, bool) {
		for {
			v, ok := iter.Next()
			if !ok {
				return *new(V), false
			}
			skip := i == index
			i++
			if skip {
				continue
			}
			return v, ok
		}
	})
}

// 値ごとに関数を実行する。
func ForEach[V any](iter Iter[V], f func(V)) {
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		f(v)
	}
}

// 他のイテレータと関数で比較し、一致していたらtrueを返す。
func EqualBy[V any](iter1 Iter[V], iter2 Iter[V], f func(V, V) bool) bool {
	for {
		v1, ok1 := iter1.Next()
		v2, ok2 := iter2.Next()
		if !ok1 || !ok2 {
			return ok1 == ok2
		}
		if !f(v1, v2) {
			return false
		}
	}
}

// 他のイテレータと一致していたらtrueを返す。
func Equal[V comparable](iter1 Iter[V], iter2 Iter[V]) bool {
	return EqualBy(iter1, iter2, func(v1 V, v2 V) bool { return v1 == v2 })
}

// 条件を満たす値の数を返す。
func CountBy[V any](iter Iter[V], f func(V) bool) int {
	return Len(FilterBy(iter, f))
}

// 一致する値の数を返す。
func Count[V comparable](iter Iter[V], v V) int {
	return Len(Filter(iter, v))
}

// 位置のイテレータを返す。
func Indices[V any](iter Iter[V]) Iter[int] {
	i := 0
	return IterFunc[int](func() (int, bool) {
		if _, ok := iter.Next(); !ok {
			return 0, false
		}
		v := i
		i++
		return v, true
	})
}

// 値を変換したイテレータを返す。
func Map[V1 any, V2 any](iter Iter[V1], f func(V1) V2) Iter[V2] {
	return IterFunc[V2](func() (V2, bool) {
		v1, ok := iter.Next()
		if !ok {
			return *new(V2), false
		}
		return f(v1), true
	})
}

// 値を順に演算する。
func Reduce[V any](iter Iter[V], f func(V, V) V) V {
	v, ok := iter.Next()
	if ok {
		for {
			v1, ok := iter.Next()
			if !ok {
				break
			}
			v = f(v, v1)
		}
	}
	return v
}

// 値の合計を演算する。
func Sum[V constraints.Ordered | constraints.Complex](iter Iter[V]) V {
	return Reduce(iter, func(sum V, v V) V {
		return sum + v
	})
}

// 値を変換して合計を演算する。
func SumBy[V1 any, V2 constraints.Ordered | constraints.Complex](iter Iter[V1], f func(V1) V2) V2 {
	return Sum(Map(iter, f))
}

// 最大の値を返す。
func Max[V constraints.Ordered](iter Iter[V]) V {
	return Reduce(iter, func(max V, v V) V {
		if max < v {
			return v
		}
		return max
	})
}

// 値を変換して最大の値を返す
func MaxBy[V1 any, V2 constraints.Ordered](iter Iter[V1], f func(V1) V2) V2 {
	return Max(Map(iter, f))
}

// 最小の値を返す。
func Min[V constraints.Ordered](iter Iter[V]) V {
	return Reduce(iter, func(min V, v V) V {
		if min > v {
			return v
		}
		return min
	})
}

// 値を変換して最小の値を返す
func MinBy[V1 any, V2 constraints.Ordered](iter Iter[V1], f func(V1) V2) V2 {
	return Min(Map(iter, f))
}

// 初期値と値を順に演算する。
func Fold[V1 any, V2 any](iter Iter[V1], v V2, f func(V2, V1) V2) V2 {
	for {
		v1, ok := iter.Next()
		if !ok {
			break
		}
		v = f(v, v1)
	}
	return v
}

// 条件を満たす最初の値の位置を返す。
func IndexBy[V any](iter Iter[V], f func(V) bool) int {
	i := 0
	for {
		v, ok := iter.Next()
		if !ok {
			return -1
		}
		if f(v) {
			return i
		}
		i++
	}
}

// 一致する最初の値の位置を返す。
func Index[V comparable](iter Iter[V], v V) int {
	return IndexBy(iter, func(v1 V) bool { return v1 == v })
}

// 条件を満たす最後の値の位置を返す。
func LastIndexBy[V any](iter Iter[V], f func(V) bool) int {
	i := 0
	index := -1
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if f(v) {
			index = i
		}
		i++
	}
	return index
}

// 一致する最後の値の位置を返す。
func LastIndex[V comparable](iter Iter[V], v V) int {
	return LastIndexBy(iter, func(v1 V) bool { return v1 == v })
}

// 条件を満たす値を返す。
func FindBy[V any](iter Iter[V], f func(V) bool) (V, bool) {
	for {
		v, ok := iter.Next()
		if !ok {
			return *new(V), false
		}
		if f(v) {
			return v, true
		}
	}
}

// 一致する値を返す。
func Find[V comparable](iter Iter[V], v V) (V, bool) {
	return FindBy(iter, func(v1 V) bool { return v1 == v })
}

// 条件を満たす値が存在したらtrueを返す。
func ExistsBy[V any](iter Iter[V], f func(V) bool) bool {
	return IndexBy(iter, f) >= 0
}

// 一致する値が存在したらtrueを返す。
func Exists[V comparable](iter Iter[V], v V) bool {
	return ExistsBy(iter, func(v1 V) bool { return v1 == v })
}

// すべての値が条件を満たせばtrueを返す。
func ForAllBy[V any](iter Iter[V], f func(V) bool) bool {
	return !ExistsBy(iter, func(v V) bool { return !f(v) })
}

// すべての値が一致したらtrueを返す。
func ForAll[V comparable](iter Iter[V], v V) bool {
	return ForAllBy(iter, func(v1 V) bool { return v1 == v })
}

// 他のイテレータの値がひとつでも存在していたらtrueを返す。
func ContainsAny[V comparable](iter Iter[V], subset Iter[V]) bool {
	slice := ToSlice(subset)
	for {
		v, ok := iter.Next()
		if !ok {
			return false
		}
		for i := range slice {
			if slice[i] == v {
				return true
			}
		}
	}
}

// 他のイテレータの値がすべて存在していたらtrueを返す。
func ContainsAll[V comparable](iter Iter[V], subset Iter[V]) bool {
	slice := ToSlice(subset)
	for {
		v, ok := iter.Next()
		if !ok {
			return false
		}
		c := 0
		for i := range slice {
			if slice[i] == v {
				c++
			}
		}
		if c == len(slice) {
			return true
		}
	}
}

// 先頭が他のイテレータと一致していたらtrueを返す。
func StartsWith[V comparable](iter Iter[V], subset Iter[V]) bool {
	for {
		v1, ok1 := iter.Next()
		v2, ok2 := subset.Next()
		if !ok2 {
			return true
		}
		if !ok1 || v1 != v2 {
			return false
		}
	}
}

// 終端が他のイテレータと一致していたらtrueを返す。
func EndsWith[V comparable](iter Iter[V], subset Iter[V]) bool {
	slice := ToSlice(subset)
	i := 0
	for {
		v, ok := iter.Next()
		if !ok {
			return i == len(slice)
		}
		if v == slice[i] {
			i++
		} else {
			i = 0
		}
	}
}

// ひとつめのoldをnewで置き換えたイテレータを返す。
func Replace[V comparable](iter Iter[V], old V, new V) Iter[V] {
	done := true
	return Map(iter, func(v V) V {
		if done && v == old {
			done = false
			return new
		}
		return v
	})
}

// すべてのoldをnewで置き換えたイテレータを返す。
func ReplaceAll[V comparable](iter Iter[V], old V, new V) Iter[V] {
	return Map(iter, func(v V) V {
		if v == old {
			return new
		}
		return v
	})
}

// 条件を満たす値だけのイテレータを返す。
func FilterBy[V any](iter Iter[V], f func(V) bool) Iter[V] {
	return IterFunc[V](func() (V, bool) {
		for {
			v, ok := iter.Next()
			if ok && !f(v) {
				continue
			}
			return v, ok
		}
	})
}

// 一致する値だけのイテレータを返す。
func Filter[V comparable](iter Iter[V], v V) Iter[V] {
	return FilterBy(iter, func(v1 V) bool { return v1 == v })
}

// 条件を満たす値を除いたイテレータを返す。
func FilterNotBy[V any](iter Iter[V], f func(V) bool) Iter[V] {
	return FilterBy(iter, func(v V) bool { return !f(v) })
}

// 一致する値を除いたイテレータを返す。
func FilterNot[V comparable](iter Iter[V], v V) Iter[V] {
	return FilterNotBy(iter, func(v1 V) bool { return v1 == v })
}

// 条件を満たす値の直前で分割したふたつのイテレータを返す。
func SplitBy[V any](iter Iter[V], f func(V) bool) (Iter[V], Iter[V]) {
	remain := []V{}
	for {
		v, ok := iter.Next()
		if !ok {
			return FromSlice(remain), Empty[V]()
		}
		if f(v) {
			return FromSlice(remain), Insert(iter, 0, v)
		}
		remain = append(remain, v)
	}
}

// 一致する値の直前で分割したふたつのイテレータを返す。
func Split[V comparable](iter Iter[V], v V) (Iter[V], Iter[V]) {
	return SplitBy(iter, func(v1 V) bool { return v1 == v })
}

// 条件を満たす値の直後で分割したふたつのイテレータを返す。
func SplitAfterBy[V any](iter Iter[V], f func(V) bool) (Iter[V], Iter[V]) {
	remain := []V{}
	for {
		v, ok := iter.Next()
		if !ok {
			return FromSlice(remain), Empty[V]()
		}
		if f(v) {
			return Insert(FromSlice(remain), 0, v), iter
		}
		remain = append(remain, v)
	}
}

// 一致する値の直後で分割したふたつのイテレータを返す。
func SplitAfter[V comparable](iter Iter[V], v V) (Iter[V], Iter[V]) {
	return SplitAfterBy(iter, func(v1 V) bool { return v1 == v })
}

// 条件を満たすイテレータと満たさないイテレータを返す。
func PartitionBy[V any](iter Iter[V], f func(V) bool) (Iter[V], Iter[V]) {
	slice1, slice2 := []V{}, []V{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if f(v) {
			slice1 = append(slice1, v)
		} else {
			slice2 = append(slice2, v)
		}
	}
	return FromSlice(slice1), FromSlice(slice2)
}

// 値の一致するイテレータと一致しないイテレータを返す。
func Partition[V comparable](iter Iter[V], v V) (Iter[V], Iter[V]) {
	return PartitionBy(iter, func(v1 V) bool { return v1 == v })
}

// 条件を満たし続ける先頭の値のイテレータを返す。
func TakeWhileBy[V any](iter Iter[V], f func(V) bool) Iter[V] {
	done := false
	return IterFunc[V](func() (V, bool) {
		if done {
			return *new(V), false
		}

		v, ok := iter.Next()
		if !ok || !f(v) {
			done = true
			return *new(V), false
		}
		return v, true
	})
}

// 一致し続ける先頭の値のイテレータを返す。
func TakeWhile[V comparable](iter Iter[V], v V) Iter[V] {
	return TakeWhileBy(iter, func(v1 V) bool { return v1 == v })
}

// 先頭n個の値のイテレータを返す。
func Take[V any](iter Iter[V], n int) Iter[V] {
	i := 0
	return TakeWhileBy(iter, func(v V) bool {
		ok := i < n
		i++
		return ok
	})
}

// 条件を満たし続ける先頭の値を除いたイテレータを返す。
func DropWhileBy[V any](iter Iter[V], f func(V) bool) Iter[V] {
	done := false
	return IterFunc[V](func() (V, bool) {
		if done {
			return iter.Next()
		}

		for {
			v, ok := iter.Next()
			if !ok {
				return *new(V), false
			}
			if f(v) {
				continue
			}
			done = true
			return v, true
		}
	})
}

// 一致し続ける先頭の値を除いたイテレータを返す。
func DropWhile[V comparable](iter Iter[V], v V) Iter[V] {
	return DropWhileBy(iter, func(v1 V) bool { return v1 == v })
}

// 先頭n個の値を除いたイテレータを返す。
func Drop[V any](iter Iter[V], n int) Iter[V] {
	i := 0
	return DropWhileBy(iter, func(v V) bool {
		ok := i < n
		i++
		return ok
	})
}

// 条件を満たし続ける先頭部分と残りの部分、ふたつのイテレータを返す。
func SpanBy[V any](iter Iter[V], f func(V) bool) (Iter[V], Iter[V]) {
	remain := []V{}
	for {
		v, ok := iter.Next()
		if !ok {
			return FromSlice(remain), Empty[V]()
		}
		if !f(v) {
			return FromSlice(remain), Insert(iter, 0, v)
		}
		remain = append(remain, v)
	}
}

// 一致し続ける先頭部分と残りの部分、ふたつのイテレータを返す。
func Span[V comparable](iter Iter[V], v V) (Iter[V], Iter[V]) {
	return SpanBy(iter, func(v1 V) bool { return v1 == v })
}

// ゼロ値を除いたイテレータを返す。
func Clean[V comparable](iter Iter[V]) Iter[V] {
	return FilterNot(iter, *new(V))
}

// 重複を除いたイテレータを返す。
func Distinct[V comparable](iter Iter[V]) Iter[V] {
	m := map[V]struct{}{}
	return IterFunc[V](func() (V, bool) {
		for {
			v, ok := iter.Next()
			if !ok {
				return *new(V), false
			}

			if _, ok := m[v]; ok {
				continue
			}
			m[v] = struct{}{}

			return v, true
		}
	})
}

// 条件を満たす値を変換したイテレータを返す。
func Collect[V1 any, V2 any](iter Iter[V1], f func(V1) (V2, bool)) Iter[V2] {
	return IterFunc[V2](func() (V2, bool) {
		for {
			v1, ok := iter.Next()
			if !ok {
				return *new(V2), false
			}

			v2, ok := f(v1)
			if !ok {
				continue
			}
			return v2, true
		}
	})
}

// 値と位置をペアにしたイテレータを返す。
func ZipWithIndex[V any](iter Iter[V]) Iter[tuple.T2[V, int]] {
	i := 0
	return IterFunc[tuple.T2[V, int]](func() (tuple.T2[V, int], bool) {
		v, ok := iter.Next()
		if !ok {
			return tuple.NewT2(*new(V), 0), false
		}
		t := tuple.NewT2(v, i)
		i++
		return t, true
	})
}

// ２つのイテレータの同じ位置の値をペアにしたイテレータを返す。
func Zip2[V1 any, V2 any](iter1 Iter[V1], iter2 Iter[V2]) Iter[tuple.T2[V1, V2]] {
	return IterFunc[tuple.T2[V1, V2]](func() (tuple.T2[V1, V2], bool) {
		v1, ok := iter1.Next()
		if !ok {
			return tuple.NewT2(*new(V1), *new(V2)), false
		}
		v2, ok := iter2.Next()
		if !ok {
			return tuple.NewT2(*new(V1), *new(V2)), false
		}
		return tuple.NewT2(v1, v2), true
	})
}

// ３つのイテレータの同じ位置の値をペアにしたイテレータを返す。
func Zip3[V1 any, V2 any, V3 any](iter1 Iter[V1], iter2 Iter[V2], iter3 Iter[V3]) Iter[tuple.T3[V1, V2, V3]] {
	return IterFunc[tuple.T3[V1, V2, V3]](func() (tuple.T3[V1, V2, V3], bool) {
		v1, ok := iter1.Next()
		if !ok {
			return tuple.NewT3(*new(V1), *new(V2), *new(V3)), false
		}
		v2, ok := iter2.Next()
		if !ok {
			return tuple.NewT3(*new(V1), *new(V2), *new(V3)), false
		}
		v3, ok := iter3.Next()
		if !ok {
			return tuple.NewT3(*new(V1), *new(V2), *new(V3)), false
		}
		return tuple.NewT3(v1, v2, v3), true
	})
}

// ４つのイテレータの同じ位置の値をペアにしたイテレータを返す。
func Zip4[V1 any, V2 any, V3 any, V4 any](iter1 Iter[V1], iter2 Iter[V2], iter3 Iter[V3], iter4 Iter[V4]) Iter[tuple.T4[V1, V2, V3, V4]] {
	return IterFunc[tuple.T4[V1, V2, V3, V4]](func() (tuple.T4[V1, V2, V3, V4], bool) {
		v1, ok := iter1.Next()
		if !ok {
			return tuple.NewT4(*new(V1), *new(V2), *new(V3), *new(V4)), false
		}
		v2, ok := iter2.Next()
		if !ok {
			return tuple.NewT4(*new(V1), *new(V2), *new(V3), *new(V4)), false
		}
		v3, ok := iter3.Next()
		if !ok {
			return tuple.NewT4(*new(V1), *new(V2), *new(V3), *new(V4)), false
		}
		v4, ok := iter4.Next()
		if !ok {
			return tuple.NewT4(*new(V1), *new(V2), *new(V3), *new(V4)), false
		}
		return tuple.NewT4(v1, v2, v3, v4), true
	})
}

// ５つのイテレータの同じ位置の値をペアにしたイテレータを返す。
func Zip5[V1 any, V2 any, V3 any, V4 any, V5 any](iter1 Iter[V1], iter2 Iter[V2], iter3 Iter[V3], iter4 Iter[V4], iter5 Iter[V5]) Iter[tuple.T5[V1, V2, V3, V4, V5]] {
	return IterFunc[tuple.T5[V1, V2, V3, V4, V5]](func() (tuple.T5[V1, V2, V3, V4, V5], bool) {
		v1, ok := iter1.Next()
		if !ok {
			return tuple.NewT5(*new(V1), *new(V2), *new(V3), *new(V4), *new(V5)), false
		}
		v2, ok := iter2.Next()
		if !ok {
			return tuple.NewT5(*new(V1), *new(V2), *new(V3), *new(V4), *new(V5)), false
		}
		v3, ok := iter3.Next()
		if !ok {
			return tuple.NewT5(*new(V1), *new(V2), *new(V3), *new(V4), *new(V5)), false
		}
		v4, ok := iter4.Next()
		if !ok {
			return tuple.NewT5(*new(V1), *new(V2), *new(V3), *new(V4), *new(V5)), false
		}
		v5, ok := iter5.Next()
		if !ok {
			return tuple.NewT5(*new(V1), *new(V2), *new(V3), *new(V4), *new(V5)), false
		}
		return tuple.NewT5(v1, v2, v3, v4, v5), true
	})
}

// ６つのイテレータの同じ位置の値をペアにしたイテレータを返す。
func Zip6[V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](iter1 Iter[V1], iter2 Iter[V2], iter3 Iter[V3], iter4 Iter[V4], iter5 Iter[V5], iter6 Iter[V6]) Iter[tuple.T6[V1, V2, V3, V4, V5, V6]] {
	return IterFunc[tuple.T6[V1, V2, V3, V4, V5, V6]](func() (tuple.T6[V1, V2, V3, V4, V5, V6], bool) {
		v1, ok := iter1.Next()
		if !ok {
			return tuple.NewT6(*new(V1), *new(V2), *new(V3), *new(V4), *new(V5), *new(V6)), false
		}
		v2, ok := iter2.Next()
		if !ok {
			return tuple.NewT6(*new(V1), *new(V2), *new(V3), *new(V4), *new(V5), *new(V6)), false
		}
		v3, ok := iter3.Next()
		if !ok {
			return tuple.NewT6(*new(V1), *new(V2), *new(V3), *new(V4), *new(V5), *new(V6)), false
		}
		v4, ok := iter4.Next()
		if !ok {
			return tuple.NewT6(*new(V1), *new(V2), *new(V3), *new(V4), *new(V5), *new(V6)), false
		}
		v5, ok := iter5.Next()
		if !ok {
			return tuple.NewT6(*new(V1), *new(V2), *new(V3), *new(V4), *new(V5), *new(V6)), false
		}
		v6, ok := iter6.Next()
		if !ok {
			return tuple.NewT6(*new(V1), *new(V2), *new(V3), *new(V4), *new(V5), *new(V6)), false
		}
		return tuple.NewT6(v1, v2, v3, v4, v5, v6), true
	})
}

// 値のペアを分離して２つのイテレータを返す。
func Unzip2[V1 any, V2 any](iter Iter[tuple.T2[V1, V2]]) (Iter[V1], Iter[V2]) {
	// TODO: イテレータを複製できれば、いちいち読み切る必要がなくなる
	slice1, slice2 := []V1{}, []V2{}
	for {
		t, ok := iter.Next()
		if !ok {
			break
		}
		slice1 = append(slice1, t.V1)
		slice2 = append(slice2, t.V2)
	}
	return FromSlice(slice1), FromSlice(slice2)
}

// 値のペアを分離して３つのイテレータを返す。
func Unzip3[V1 any, V2 any, V3 any](iter Iter[tuple.T3[V1, V2, V3]]) (Iter[V1], Iter[V2], Iter[V3]) {
	// TODO: イテレータを複製できれば、いちいち読み切る必要がなくなる
	slice1, slice2, slice3 := []V1{}, []V2{}, []V3{}
	for {
		t, ok := iter.Next()
		if !ok {
			break
		}
		slice1 = append(slice1, t.V1)
		slice2 = append(slice2, t.V2)
		slice3 = append(slice3, t.V3)
	}
	return FromSlice(slice1), FromSlice(slice2), FromSlice(slice3)
}

// 値のペアを分離して４つのイテレータを返す。
func Unzip4[V1 any, V2 any, V3 any, V4 any](iter Iter[tuple.T4[V1, V2, V3, V4]]) (Iter[V1], Iter[V2], Iter[V3], Iter[V4]) {
	// TODO: イテレータを複製できれば、いちいち読み切る必要がなくなる
	slice1, slice2, slice3, slice4 := []V1{}, []V2{}, []V3{}, []V4{}
	for {
		t, ok := iter.Next()
		if !ok {
			break
		}
		slice1 = append(slice1, t.V1)
		slice2 = append(slice2, t.V2)
		slice3 = append(slice3, t.V3)
		slice4 = append(slice4, t.V4)
	}
	return FromSlice(slice1), FromSlice(slice2), FromSlice(slice3), FromSlice(slice4)
}

// 値のペアを分離して５つのイテレータを返す。
func Unzip5[V1 any, V2 any, V3 any, V4 any, V5 any](iter Iter[tuple.T5[V1, V2, V3, V4, V5]]) (Iter[V1], Iter[V2], Iter[V3], Iter[V4], Iter[V5]) {
	// TODO: イテレータを複製できれば、いちいち読み切る必要がなくなる
	slice1, slice2, slice3, slice4, slice5 := []V1{}, []V2{}, []V3{}, []V4{}, []V5{}
	for {
		t, ok := iter.Next()
		if !ok {
			break
		}
		slice1 = append(slice1, t.V1)
		slice2 = append(slice2, t.V2)
		slice3 = append(slice3, t.V3)
		slice4 = append(slice4, t.V4)
		slice5 = append(slice5, t.V5)
	}
	return FromSlice(slice1), FromSlice(slice2), FromSlice(slice3), FromSlice(slice4), FromSlice(slice5)
}

// 値のペアを分離して６つのイテレータを返す。
func Unzip6[V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](iter Iter[tuple.T6[V1, V2, V3, V4, V5, V6]]) (Iter[V1], Iter[V2], Iter[V3], Iter[V4], Iter[V5], Iter[V6]) {
	// TODO: イテレータを複製できれば、いちいち読み切る必要がなくなる
	slice1, slice2, slice3, slice4, slice5, slice6 := []V1{}, []V2{}, []V3{}, []V4{}, []V5{}, []V6{}
	for {
		t, ok := iter.Next()
		if !ok {
			break
		}
		slice1 = append(slice1, t.V1)
		slice2 = append(slice2, t.V2)
		slice3 = append(slice3, t.V3)
		slice4 = append(slice4, t.V4)
		slice5 = append(slice5, t.V5)
		slice6 = append(slice6, t.V6)
	}
	return FromSlice(slice1), FromSlice(slice2), FromSlice(slice3), FromSlice(slice4), FromSlice(slice5), FromSlice(slice6)
}

// 値ごとに関数の返すキーでグルーピングしたマップを返す。
func GroupBy[K comparable, V any](iter Iter[V], f func(V) K) map[K][]V {
	m := map[K][]V{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		k := f(v)
		if _, ok := m[k]; ok {
			m[k] = append(m[k], v)
		} else {
			m[k] = []V{v}
		}
	}
	return m
}

// 平坦化したイテレータを返す。
func Flatten[V any](iter Iter[Iter[V]]) Iter[V] {
	sub, ok := iter.Next()
	return IterFunc[V](func() (V, bool) {
		for {
			if !ok {
				return *new(V), false
			}

			if v, ok := sub.Next(); ok {
				return v, true
			}

			sub, ok = iter.Next()
		}
	})
}

// 値をイテレータに変換し、それらを結合したイテレータを返す。
func FlatMap[V1 any, V2 any](iter Iter[V1], f func(V1) Iter[V2]) Iter[V2] {
	return Flatten(Map(iter, f))
}

// 値のあいだにseparatorを挿入したイテレータを返す。
func Join[V any](iter Iter[V], separator V) Iter[V] {
	return Drop(FlatMap(iter, func(v V) Iter[V] { return From(separator, v) }), 1)
}
