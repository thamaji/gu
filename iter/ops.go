package iter

import (
	"github.com/thamaji/gu/tuple"
	"golang.org/x/exp/constraints"
)

// 値の数を返す。
func Len[V any](iter Iter[V]) (int, error) {
	c := 0
	for {
		if _, ok := iter.Next(); ok {
			break
		}
		c++
	}
	if err := iter.Err(); err != nil {
		return 0, err
	}
	return c, nil
}

// 指定した位置の値を返す。
func Get[V any](iter Iter[V], index int) (V, bool, error) {
	i := 0
	for {
		v, ok := iter.Next()
		if !ok {
			return *new(V), false, iter.Err()
		}
		if i == index {
			return v, true, nil
		}
		i++
	}
}

// 指定した位置の値を返す。無い場合はvを返す。
func GetOrElse[V any](iter Iter[V], index int, v V) (V, error) {
	v, ok, err := Get(iter, index)
	if !ok {
		return *new(V), err
	}
	return v, nil
}

// 先頭の値を返す。
func GetFirst[V any](iter Iter[V]) (V, bool, error) {
	return Get(iter, 0)
}

// 先頭の値を返す。無い場合はvを返す。
func GetFirstOrElse[V any](iter Iter[V], v V) (V, error) {
	v, ok, err := GetFirst(iter)
	if !ok {
		return *new(V), err
	}
	return v, nil
}

// 終端の値を返す。
func GetLast[V any](iter Iter[V]) (V, bool, error) {
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
	return v, ok, iter.Err()
}

// 終端の値を返す。無い場合はvを返す。
func GetLastOrElse[V any](iter Iter[V], v V) (V, error) {
	v, ok, err := GetLast(iter)
	if !ok {
		return *new(V), err
	}
	return v, nil
}

// イテレータの末尾に他のイテレータを結合する。
func Concat[V any](iter1 Iter[V], iter2 ...Iter[V]) Iter[V] {
	cursor := 0
	iters := append([]Iter[V]{iter1}, iter2...)

	w := &wrappedIter[V]{}
	w.next = func() (V, bool) {
		if cursor >= len(iters) {
			return *new(V), false
		}
		v, ok := iters[cursor].Next()
		if !ok {
			w.err = iters[cursor].Err()
			cursor++
		}
		return v, ok
	}
	return w
}

// イテレータの末尾に値を追加する。
func Append[V any](iter Iter[V], v ...V) Iter[V] {
	return Concat(iter, FromSlice(v))
}

// 指定した位置に値を追加する。
func Insert[V any](iter Iter[V], index int, v ...V) Iter[V] {
	i := 0
	w := &wrappedIter[V]{}
	w.next = func() (V, bool) {
		if i >= index && i < index+len(v) {
			v1 := v[i-index]
			i++
			return v1, true
		}
		v, ok := iter.Next()
		if !ok {
			w.err = iter.Err()
		}
		i++
		return v, ok
	}
	return w
}

// 指定した位置の値を削除する。
func Remove[V any](iter Iter[V], index int) Iter[V] {
	i := 0
	w := &wrappedIter[V]{}
	w.next = func() (V, bool) {
		for {
			v, ok := iter.Next()
			if !ok {
				w.err = iter.Err()
				return *new(V), false
			}
			skip := i == index
			i++
			if skip {
				continue
			}
			return v, ok
		}
	}
	return w
}

// 値ごとに関数を実行する。
func ForEach[V any](iter Iter[V], f func(V) error) error {
	for {
		v, ok := iter.Next()
		if !ok {
			if err := iter.Err(); err != nil {
				return err
			}
			return nil
		}
		if err := f(v); err != nil {
			return err
		}
	}
}

// 他のイテレータと関数で比較し、一致していたらtrueを返す。
func EqualBy[V any](iter1 Iter[V], iter2 Iter[V], f func(V, V) (bool, error)) (bool, error) {
	for {
		v1, ok1 := iter1.Next()
		if !ok1 {
			if err := iter1.Err(); err != nil {
				return false, err
			}
		}

		v2, ok2 := iter2.Next()
		if !ok2 {
			if err := iter2.Err(); err != nil {
				return false, err
			}
		}

		if !ok1 || !ok2 {
			return ok1 == ok2, nil
		}

		ok, err := f(v1, v2)
		if err != nil {
			return false, err
		}

		if !ok {
			return false, nil
		}
	}
}

// 他のイテレータと一致していたらtrueを返す。
func Equal[V comparable](iter1 Iter[V], iter2 Iter[V]) (bool, error) {
	return EqualBy(iter1, iter2, func(v1 V, v2 V) (bool, error) { return v1 == v2, nil })
}

// 条件を満たす値の数を返す。
func CountBy[V any](iter Iter[V], f func(V) (bool, error)) (int, error) {
	return Len(FilterBy(iter, f))
}

// 一致する値の数を返す。
func Count[V comparable](iter Iter[V], v V) (int, error) {
	return Len(Filter(iter, v))
}

// 位置のイテレータを返す。
func Indices[V any](iter Iter[V]) Iter[int] {
	i := 0
	w := &wrappedIter[int]{}
	w.next = func() (int, bool) {
		if _, ok := iter.Next(); !ok {
			w.err = iter.Err()
			return 0, false
		}
		v := i
		i++
		return v, true
	}
	return w
}

// 値を変換したイテレータを返す。
func Map[V1 any, V2 any](iter Iter[V1], f func(V1) (V2, error)) Iter[V2] {
	w := &wrappedIter[V2]{}
	w.next = func() (V2, bool) {
		v1, ok := iter.Next()
		if !ok {
			w.err = iter.Err()
			return *new(V2), false
		}
		v2, err := f(v1)
		if err != nil {
			w.err = err
			return *new(V2), false
		}
		return v2, true
	}
	return w
}

// 値を順に演算する。
func Reduce[V any](iter Iter[V], f func(V, V) (V, error)) (V, error) {
	var err error
	v, ok := iter.Next()
	if ok {
		for {
			v1, ok := iter.Next()
			if !ok {
				break
			}
			v, err = f(v, v1)
			if err != nil {
				break
			}
		}
	}
	if err1 := iter.Err(); err == nil {
		err = err1
	}
	return v, err
}

// 値の合計を演算する。
func Sum[V constraints.Ordered | constraints.Complex](iter Iter[V]) (V, error) {
	return Reduce(iter, func(sum V, v V) (V, error) {
		return sum + v, nil
	})
}

// 値を変換して合計を演算する。
func SumBy[V1 any, V2 constraints.Ordered | constraints.Complex](iter Iter[V1], f func(V1) (V2, error)) (V2, error) {
	return Sum(Map(iter, f))
}

// 最大の値を返す。
func Max[V constraints.Ordered](iter Iter[V]) (V, error) {
	return Reduce(iter, func(max V, v V) (V, error) {
		if max < v {
			return v, nil
		}
		return max, nil
	})
}

// 値を変換して最大の値を返す
func MaxBy[V1 any, V2 constraints.Ordered](iter Iter[V1], f func(V1) (V2, error)) (V2, error) {
	return Max(Map(iter, f))
}

// 最小の値を返す。
func Min[V constraints.Ordered](iter Iter[V]) (V, error) {
	return Reduce(iter, func(min V, v V) (V, error) {
		if min > v {
			return v, nil
		}
		return min, nil
	})
}

// 値を変換して最小の値を返す
func MinBy[V1 any, V2 constraints.Ordered](iter Iter[V1], f func(V1) (V2, error)) (V2, error) {
	return Min(Map(iter, f))
}

// 初期値と値を順に演算する。
func Fold[V1 any, V2 any](iter Iter[V1], v V2, f func(V2, V1) (V2, error)) (V2, error) {
	var err error
	for {
		v1, ok := iter.Next()
		if !ok {
			break
		}
		v, err = f(v, v1)
		if err != nil {
			break
		}
	}
	if err1 := iter.Err(); err == nil {
		err = err1
	}
	return v, err
}

// 条件を満たす最初の値の位置を返す。
func IndexBy[V any](iter Iter[V], f func(V) (bool, error)) (int, error) {
	i := 0
	for {
		v, ok := iter.Next()
		if !ok {
			return -1, iter.Err()
		}
		ok, err := f(v)
		if err != nil {
			return -1, err
		}
		if ok {
			return i, nil
		}
		i++
	}
}

// 一致する最初の値の位置を返す。
func Index[V comparable](iter Iter[V], v V) (int, error) {
	return IndexBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 条件を満たす最後の値の位置を返す。
func LastIndexBy[V any](iter Iter[V], f func(V) (bool, error)) (int, error) {
	var err error
	i := 0
	index := -1
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		ok, err = f(v)
		if err != nil {
			break
		}
		if ok {
			index = i
		}
		i++
	}
	if err1 := iter.Err(); err == nil {
		err = err1
	}
	return index, err
}

// 一致する最後の値の位置を返す。
func LastIndex[V comparable](iter Iter[V], v V) (int, error) {
	return LastIndexBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 条件を満たす値を返す。
func FindBy[V any](iter Iter[V], f func(V) (bool, error)) (V, bool, error) {
	for {
		v, ok := iter.Next()
		if !ok {
			return *new(V), false, iter.Err()
		}
		ok, err := f(v)
		if err != nil {
			return *new(V), false, err
		}
		if ok {
			return v, true, nil
		}
	}
}

// 一致する値を返す。
func Find[V comparable](iter Iter[V], v V) (V, bool, error) {
	return FindBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 条件を満たす値が存在したらtrueを返す。
func ExistsBy[V any](iter Iter[V], f func(V) (bool, error)) (bool, error) {
	index, err := IndexBy(iter, f)
	return index >= 0, err
}

// 一致する値が存在したらtrueを返す。
func Exists[V comparable](iter Iter[V], v V) (bool, error) {
	return ExistsBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// すべての値が条件を満たせばtrueを返す。
func ForAllBy[V any](iter Iter[V], f func(V) (bool, error)) (bool, error) {
	ok, err := ExistsBy(iter, func(v V) (bool, error) {
		ok, err := f(v)
		return !ok, err
	})
	if err != nil {
		return false, err
	}
	return !ok, nil
}

// すべての値が一致したらtrueを返す。
func ForAll[V comparable](iter Iter[V], v V) (bool, error) {
	return ForAllBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 他のイテレータの値がひとつでも存在していたらtrueを返す。
func ContainsAny[V comparable](iter Iter[V], subset Iter[V]) (bool, error) {
	slice, err := ToSlice(subset)
	if err != nil {
		return false, err
	}

	for {
		v, ok := iter.Next()
		if !ok {
			return false, iter.Err()
		}
		for i := range slice {
			if slice[i] == v {
				return true, nil
			}
		}
	}
}

// 他のイテレータの値がすべて存在していたらtrueを返す。
func ContainsAll[V comparable](iter Iter[V], subset Iter[V]) (bool, error) {
	slice, err := ToSlice(subset)
	if err != nil {
		return false, err
	}
	for {
		v, ok := iter.Next()
		if !ok {
			return false, iter.Err()
		}
		c := 0
		for i := range slice {
			if slice[i] == v {
				c++
			}
		}
		if c == len(slice) {
			return true, nil
		}
	}
}

// 先頭が他のイテレータと一致していたらtrueを返す。
func StartsWith[V comparable](iter Iter[V], subset Iter[V]) (bool, error) {
	for {
		v1, ok1 := iter.Next()
		v2, ok2 := subset.Next()
		if !ok2 {
			if err := subset.Err(); err != nil {
				return false, err
			}
			return true, nil
		}
		if !ok1 {
			return false, iter.Err()
		}
		if v1 != v2 {
			return false, nil
		}
	}
}

// 終端が他のイテレータと一致していたらtrueを返す。
func EndsWith[V comparable](iter Iter[V], subset Iter[V]) (bool, error) {
	slice, err := ToSlice(subset)
	if err != nil {
		return false, err
	}
	i := 0
	for {
		v, ok := iter.Next()
		if !ok {
			if err := subset.Err(); err != nil {
				return false, err
			}
			return i == len(slice), nil
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
	return Map(iter, func(v V) (V, error) {
		if done && v == old {
			done = false
			return new, nil
		}
		return v, nil
	})
}

// すべてのoldをnewで置き換えたイテレータを返す。
func ReplaceAll[V comparable](iter Iter[V], old V, new V) Iter[V] {
	return Map(iter, func(v V) (V, error) {
		if v == old {
			return new, nil
		}
		return v, nil
	})
}

// 条件を満たす値だけのイテレータを返す。
func FilterBy[V any](iter Iter[V], f func(V) (bool, error)) Iter[V] {
	w := &wrappedIter[V]{}
	w.next = func() (V, bool) {
		for {
			v, ok := iter.Next()
			if !ok {
				w.err = iter.Err()
				return *new(V), false
			}
			ok, err := f(v)
			if err != nil {
				w.err = err
				return *new(V), false
			}
			if !ok {
				continue
			}
			return v, true
		}
	}
	return w
}

// 一致する値だけのイテレータを返す。
func Filter[V comparable](iter Iter[V], v V) Iter[V] {
	return FilterBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 条件を満たす値を除いたイテレータを返す。
func FilterNotBy[V any](iter Iter[V], f func(V) (bool, error)) Iter[V] {
	return FilterBy(iter, func(v V) (bool, error) {
		ok, err := f(v)
		return !ok, err
	})
}

// 一致する値を除いたイテレータを返す。
func FilterNot[V comparable](iter Iter[V], v V) Iter[V] {
	return FilterNotBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 条件を満たす値の直前で分割したふたつのイテレータを返す。
func SplitBy[V any](iter Iter[V], f func(V) (bool, error)) (Iter[V], Iter[V], error) {
	remain := []V{}
	for {
		v, ok := iter.Next()
		if !ok {
			if err := iter.Err(); err != nil {
				return nil, nil, err
			}
			return FromSlice(remain), Empty[V](), nil
		}

		ok, err := f(v)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			return FromSlice(remain), Insert(iter, 0, v), nil
		}
		remain = append(remain, v)
	}
}

// 一致する値の直前で分割したふたつのイテレータを返す。
func Split[V comparable](iter Iter[V], v V) (Iter[V], Iter[V], error) {
	return SplitBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 条件を満たす値の直後で分割したふたつのイテレータを返す。
func SplitAfterBy[V any](iter Iter[V], f func(V) (bool, error)) (Iter[V], Iter[V], error) {
	remain := []V{}
	for {
		v, ok := iter.Next()
		if !ok {
			if err := iter.Err(); err != nil {
				return nil, nil, err
			}
			return FromSlice(remain), Empty[V](), nil
		}

		ok, err := f(v)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			return Insert(FromSlice(remain), 0, v), iter, nil
		}
		remain = append(remain, v)
	}
}

// 一致する値の直後で分割したふたつのイテレータを返す。
func SplitAfter[V comparable](iter Iter[V], v V) (Iter[V], Iter[V], error) {
	return SplitAfterBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 条件を満たすイテレータと満たさないイテレータを返す。
func PartitionBy[V any](iter Iter[V], f func(V) (bool, error)) (Iter[V], Iter[V], error) {
	slice1, slice2 := []V{}, []V{}
	for {
		v, ok := iter.Next()
		if !ok {
			if err := iter.Err(); err != nil {
				return nil, nil, err
			}
			break
		}

		ok, err := f(v)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			slice1 = append(slice1, v)
		} else {
			slice2 = append(slice2, v)
		}
	}
	return FromSlice(slice1), FromSlice(slice2), nil
}

// 値の一致するイテレータと一致しないイテレータを返す。
func Partition[V comparable](iter Iter[V], v V) (Iter[V], Iter[V], error) {
	return PartitionBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 条件を満たし続ける先頭の値のイテレータを返す。
func TakeWhileBy[V any](iter Iter[V], f func(V) (bool, error)) Iter[V] {
	done := false
	w := &wrappedIter[V]{}
	w.next = func() (V, bool) {
		if done {
			return *new(V), false
		}

		v, ok := iter.Next()
		if !ok {
			w.err = iter.Err()
			return *new(V), false
		}

		ok, err := f(v)
		if err != nil {
			w.err = err
			return *new(V), false
		}
		if !ok {
			done = true
			return *new(V), false
		}
		return v, true
	}
	return w
}

// 一致し続ける先頭の値のイテレータを返す。
func TakeWhile[V comparable](iter Iter[V], v V) Iter[V] {
	return TakeWhileBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 先頭n個の値のイテレータを返す。
func Take[V any](iter Iter[V], n int) Iter[V] {
	i := 0
	return TakeWhileBy(iter, func(v V) (bool, error) {
		ok := i < n
		i++
		return ok, nil
	})
}

// 条件を満たし続ける先頭の値を除いたイテレータを返す。
func DropWhileBy[V any](iter Iter[V], f func(V) (bool, error)) Iter[V] {
	done := false
	w := &wrappedIter[V]{}
	w.next = func() (V, bool) {
		if done {
			return iter.Next()
		}

		for {
			v, ok := iter.Next()
			if !ok {
				w.err = w.Err()
				return *new(V), false
			}

			ok, err := f(v)
			if err != nil {
				w.err = w.Err()
				return *new(V), false
			}
			if ok {
				continue
			}
			done = true
			return v, true
		}
	}
	return w
}

// 一致し続ける先頭の値を除いたイテレータを返す。
func DropWhile[V comparable](iter Iter[V], v V) Iter[V] {
	return DropWhileBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 先頭n個の値を除いたイテレータを返す。
func Drop[V any](iter Iter[V], n int) Iter[V] {
	i := 0
	return DropWhileBy(iter, func(v V) (bool, error) {
		ok := i < n
		i++
		return ok, nil
	})
}

// 条件を満たし続ける先頭部分と残りの部分、ふたつのイテレータを返す。
func SpanBy[V any](iter Iter[V], f func(V) (bool, error)) (Iter[V], Iter[V], error) {
	remain := []V{}
	for {
		v, ok := iter.Next()
		if !ok {
			if err := iter.Err(); err != nil {
				return nil, nil, err
			}
			return FromSlice(remain), Empty[V](), nil
		}

		ok, err := f(v)
		if err != nil {
			return nil, nil, err
		}
		if !ok {
			return FromSlice(remain), Insert(iter, 0, v), nil
		}
		remain = append(remain, v)
	}
}

// 一致し続ける先頭部分と残りの部分、ふたつのイテレータを返す。
func Span[V comparable](iter Iter[V], v V) (Iter[V], Iter[V], error) {
	return SpanBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// ゼロ値を除いたイテレータを返す。
func Clean[V comparable](iter Iter[V]) Iter[V] {
	return FilterNot(iter, *new(V))
}

// 重複を除いたイテレータを返す。
func Distinct[V comparable](iter Iter[V]) Iter[V] {
	m := map[V]struct{}{}
	w := &wrappedIter[V]{}
	w.next = func() (V, bool) {
		for {
			v, ok := iter.Next()
			if !ok {
				w.err = iter.Err()
				return *new(V), false
			}

			if _, ok := m[v]; ok {
				continue
			}
			m[v] = struct{}{}

			return v, true
		}
	}
	return w
}

// 条件を満たす値を変換したイテレータを返す。
func Collect[V1 any, V2 any](iter Iter[V1], f func(V1) (V2, bool, error)) Iter[V2] {
	w := &wrappedIter[V2]{}
	w.next = func() (V2, bool) {
		for {
			v1, ok := iter.Next()
			if !ok {
				w.err = iter.Err()
				return *new(V2), false
			}

			v2, ok, err := f(v1)
			if err != nil {
				w.err = iter.Err()
				return *new(V2), false
			}
			if !ok {
				continue
			}
			return v2, true
		}
	}
	return w
}

// 値と位置をペアにしたイテレータを返す。
func ZipWithIndex[V any](iter Iter[V]) Iter[tuple.T2[V, int]] {
	i := 0
	w := &wrappedIter[tuple.T2[V, int]]{}
	w.next = func() (tuple.T2[V, int], bool) {
		v, ok := iter.Next()
		if !ok {
			w.err = iter.Err()
			return tuple.NewT2(*new(V), 0), false
		}
		t := tuple.NewT2(v, i)
		i++
		return t, true
	}
	return w
}

// ２つのイテレータの同じ位置の値をペアにしたイテレータを返す。
func Zip2[V1 any, V2 any](iter1 Iter[V1], iter2 Iter[V2]) Iter[tuple.T2[V1, V2]] {
	w := &wrappedIter[tuple.T2[V1, V2]]{}
	w.next = func() (tuple.T2[V1, V2], bool) {
		v1, ok1 := iter1.Next()
		v2, ok2 := iter2.Next()
		if ok1 && ok2 {
			return tuple.NewT2(v1, v2), true
		}

		err := iter1.Err()
		if err1 := iter2.Err(); err == nil {
			err = err1
		}
		if err != nil {
			w.err = err
		}
		return tuple.NewT2(*new(V1), *new(V2)), false
	}
	return w
}

// ３つのイテレータの同じ位置の値をペアにしたイテレータを返す。
func Zip3[V1 any, V2 any, V3 any](iter1 Iter[V1], iter2 Iter[V2], iter3 Iter[V3]) Iter[tuple.T3[V1, V2, V3]] {
	w := &wrappedIter[tuple.T3[V1, V2, V3]]{}
	w.next = func() (tuple.T3[V1, V2, V3], bool) {
		v1, ok1 := iter1.Next()
		v2, ok2 := iter2.Next()
		v3, ok3 := iter3.Next()
		if ok1 && ok2 && ok3 {
			return tuple.NewT3(v1, v2, v3), true
		}

		err := iter1.Err()
		if err1 := iter2.Err(); err == nil {
			err = err1
		}
		if err1 := iter3.Err(); err == nil {
			err = err1
		}
		if err != nil {
			w.err = err
		}
		return tuple.NewT3(*new(V1), *new(V2), *new(V3)), false
	}
	return w
}

// ４つのイテレータの同じ位置の値をペアにしたイテレータを返す。
func Zip4[V1 any, V2 any, V3 any, V4 any](iter1 Iter[V1], iter2 Iter[V2], iter3 Iter[V3], iter4 Iter[V4]) Iter[tuple.T4[V1, V2, V3, V4]] {
	w := &wrappedIter[tuple.T4[V1, V2, V3, V4]]{}
	w.next = func() (tuple.T4[V1, V2, V3, V4], bool) {
		v1, ok1 := iter1.Next()
		v2, ok2 := iter2.Next()
		v3, ok3 := iter3.Next()
		v4, ok4 := iter4.Next()
		if ok1 && ok2 && ok3 && ok4 {
			return tuple.NewT4(v1, v2, v3, v4), true
		}

		err := iter1.Err()
		if err1 := iter2.Err(); err == nil {
			err = err1
		}
		if err1 := iter3.Err(); err == nil {
			err = err1
		}
		if err1 := iter4.Err(); err == nil {
			err = err1
		}
		if err != nil {
			w.err = err
		}
		return tuple.NewT4(*new(V1), *new(V2), *new(V3), *new(V4)), false
	}
	return w
}

// ５つのイテレータの同じ位置の値をペアにしたイテレータを返す。
func Zip5[V1 any, V2 any, V3 any, V4 any, V5 any](iter1 Iter[V1], iter2 Iter[V2], iter3 Iter[V3], iter4 Iter[V4], iter5 Iter[V5]) Iter[tuple.T5[V1, V2, V3, V4, V5]] {
	w := &wrappedIter[tuple.T5[V1, V2, V3, V4, V5]]{}
	w.next = func() (tuple.T5[V1, V2, V3, V4, V5], bool) {
		v1, ok1 := iter1.Next()
		v2, ok2 := iter2.Next()
		v3, ok3 := iter3.Next()
		v4, ok4 := iter4.Next()
		v5, ok5 := iter5.Next()
		if ok1 && ok2 && ok3 && ok4 && ok5 {
			return tuple.NewT5(v1, v2, v3, v4, v5), true
		}

		err := iter1.Err()
		if err1 := iter2.Err(); err == nil {
			err = err1
		}
		if err1 := iter3.Err(); err == nil {
			err = err1
		}
		if err1 := iter4.Err(); err == nil {
			err = err1
		}
		if err1 := iter5.Err(); err == nil {
			err = err1
		}
		if err != nil {
			w.err = err
		}
		return tuple.NewT5(*new(V1), *new(V2), *new(V3), *new(V4), *new(V5)), false
	}
	return w
}

// ６つのイテレータの同じ位置の値をペアにしたイテレータを返す。
func Zip6[V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](iter1 Iter[V1], iter2 Iter[V2], iter3 Iter[V3], iter4 Iter[V4], iter5 Iter[V5], iter6 Iter[V6]) Iter[tuple.T6[V1, V2, V3, V4, V5, V6]] {
	w := &wrappedIter[tuple.T6[V1, V2, V3, V4, V5, V6]]{}
	w.next = func() (tuple.T6[V1, V2, V3, V4, V5, V6], bool) {
		v1, ok1 := iter1.Next()
		v2, ok2 := iter2.Next()
		v3, ok3 := iter3.Next()
		v4, ok4 := iter4.Next()
		v5, ok5 := iter5.Next()
		v6, ok6 := iter6.Next()
		if ok1 && ok2 && ok3 && ok4 && ok5 && ok6 {
			return tuple.NewT6(v1, v2, v3, v4, v5, v6), true
		}

		err := iter1.Err()
		if err1 := iter2.Err(); err == nil {
			err = err1
		}
		if err1 := iter3.Err(); err == nil {
			err = err1
		}
		if err1 := iter4.Err(); err == nil {
			err = err1
		}
		if err1 := iter5.Err(); err == nil {
			err = err1
		}
		if err1 := iter6.Err(); err == nil {
			err = err1
		}
		if err != nil {
			w.err = err
		}
		return tuple.NewT6(*new(V1), *new(V2), *new(V3), *new(V4), *new(V5), *new(V6)), false
	}
	return w
}

// 値のペアを分離して２つのイテレータを返す。
func Unzip2[V1 any, V2 any](iter Iter[tuple.T2[V1, V2]]) (Iter[V1], Iter[V2], error) {
	// TODO: イテレータを複製できれば、いちいち読み切る必要がなくなる
	slice1, slice2 := []V1{}, []V2{}
	for {
		t, ok := iter.Next()
		if !ok {
			if err := iter.Err(); err != nil {
				return nil, nil, err
			}
			break
		}
		slice1 = append(slice1, t.V1)
		slice2 = append(slice2, t.V2)
	}
	return FromSlice(slice1), FromSlice(slice2), nil
}

// 値のペアを分離して３つのイテレータを返す。
func Unzip3[V1 any, V2 any, V3 any](iter Iter[tuple.T3[V1, V2, V3]]) (Iter[V1], Iter[V2], Iter[V3], error) {
	// TODO: イテレータを複製できれば、いちいち読み切る必要がなくなる
	slice1, slice2, slice3 := []V1{}, []V2{}, []V3{}
	for {
		t, ok := iter.Next()
		if !ok {
			if err := iter.Err(); err != nil {
				return nil, nil, nil, err
			}
			break
		}
		slice1 = append(slice1, t.V1)
		slice2 = append(slice2, t.V2)
		slice3 = append(slice3, t.V3)
	}
	return FromSlice(slice1), FromSlice(slice2), FromSlice(slice3), nil
}

// 値のペアを分離して４つのイテレータを返す。
func Unzip4[V1 any, V2 any, V3 any, V4 any](iter Iter[tuple.T4[V1, V2, V3, V4]]) (Iter[V1], Iter[V2], Iter[V3], Iter[V4], error) {
	// TODO: イテレータを複製できれば、いちいち読み切る必要がなくなる
	slice1, slice2, slice3, slice4 := []V1{}, []V2{}, []V3{}, []V4{}
	for {
		t, ok := iter.Next()
		if !ok {
			if err := iter.Err(); err != nil {
				return nil, nil, nil, nil, err
			}
			break
		}
		slice1 = append(slice1, t.V1)
		slice2 = append(slice2, t.V2)
		slice3 = append(slice3, t.V3)
		slice4 = append(slice4, t.V4)
	}
	return FromSlice(slice1), FromSlice(slice2), FromSlice(slice3), FromSlice(slice4), nil
}

// 値のペアを分離して５つのイテレータを返す。
func Unzip5[V1 any, V2 any, V3 any, V4 any, V5 any](iter Iter[tuple.T5[V1, V2, V3, V4, V5]]) (Iter[V1], Iter[V2], Iter[V3], Iter[V4], Iter[V5], error) {
	// TODO: イテレータを複製できれば、いちいち読み切る必要がなくなる
	slice1, slice2, slice3, slice4, slice5 := []V1{}, []V2{}, []V3{}, []V4{}, []V5{}
	for {
		t, ok := iter.Next()
		if !ok {
			if err := iter.Err(); err != nil {
				return nil, nil, nil, nil, nil, err
			}
			break
		}
		slice1 = append(slice1, t.V1)
		slice2 = append(slice2, t.V2)
		slice3 = append(slice3, t.V3)
		slice4 = append(slice4, t.V4)
		slice5 = append(slice5, t.V5)
	}
	return FromSlice(slice1), FromSlice(slice2), FromSlice(slice3), FromSlice(slice4), FromSlice(slice5), nil
}

// 値のペアを分離して６つのイテレータを返す。
func Unzip6[V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](iter Iter[tuple.T6[V1, V2, V3, V4, V5, V6]]) (Iter[V1], Iter[V2], Iter[V3], Iter[V4], Iter[V5], Iter[V6], error) {
	// TODO: イテレータを複製できれば、いちいち読み切る必要がなくなる
	slice1, slice2, slice3, slice4, slice5, slice6 := []V1{}, []V2{}, []V3{}, []V4{}, []V5{}, []V6{}
	for {
		t, ok := iter.Next()
		if !ok {
			if err := iter.Err(); err != nil {
				return nil, nil, nil, nil, nil, nil, err
			}
			break
		}
		slice1 = append(slice1, t.V1)
		slice2 = append(slice2, t.V2)
		slice3 = append(slice3, t.V3)
		slice4 = append(slice4, t.V4)
		slice5 = append(slice5, t.V5)
		slice6 = append(slice6, t.V6)
	}
	return FromSlice(slice1), FromSlice(slice2), FromSlice(slice3), FromSlice(slice4), FromSlice(slice5), FromSlice(slice6), nil
}

// 値ごとに関数の返すキーでグルーピングしたマップを返す。
func GroupBy[K comparable, V any](iter Iter[V], f func(V) (K, error)) (map[K][]V, error) {
	m := map[K][]V{}
	for {
		v, ok := iter.Next()
		if !ok {
			if err := iter.Err(); err != nil {
				return nil, err
			}
			break
		}

		k, err := f(v)
		if err != nil {
			return nil, err
		}
		if _, ok := m[k]; ok {
			m[k] = append(m[k], v)
		} else {
			m[k] = []V{v}
		}
	}
	return m, nil
}

// 平坦化したイテレータを返す。
func Flatten[V any](iter Iter[Iter[V]]) Iter[V] {
	sub, ok := iter.Next()
	w := &wrappedIter[V]{}
	w.next = func() (V, bool) {
		for {
			if !ok {
				w.err = iter.Err()
				return *new(V), false
			}

			if v, ok := sub.Next(); ok {
				return v, true
			}
			if err := sub.Err(); err != nil {
				w.err = err
				return *new(V), false
			}

			sub, ok = iter.Next()
		}
	}
	return w
}

// 値をイテレータに変換し、それらを結合したイテレータを返す。
func FlatMap[V1 any, V2 any](iter Iter[V1], f func(V1) (Iter[V2], error)) Iter[V2] {
	return Flatten(Map(iter, f))
}

// 値のあいだにseparatorを挿入したイテレータを返す。
func Join[V any](iter Iter[V], separator V) Iter[V] {
	return Drop(FlatMap(iter, func(v V) (Iter[V], error) { return From(separator, v), nil }), 1)
}
