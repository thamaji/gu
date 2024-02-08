package iter

import (
	"errors"

	"github.com/thamaji/gu/must"
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

// 値の数を返す。実行中にエラーが起きた場合 panic する。
func MustLen[V any](iter Iter[V]) int {
	return must.Must1(Len(iter))
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

// 指定した位置の値を返す。実行中にエラーが起きた場合 panic する。
func MustGet[V any](iter Iter[V], index int) (V, bool) {
	return must.Must2(Get(iter, index))
}

// 指定した位置の値を返す。無い場合はvを返す。
func GetOrElse[V any](iter Iter[V], index int, v V) (V, error) {
	v, ok, err := Get(iter, index)
	if !ok {
		return *new(V), err
	}
	return v, nil
}

// 指定した位置の値を返す。無い場合はvを返す。実行中にエラーが起きた場合 panic する。
func MustGetOrElse[V any](iter Iter[V], index int, v V) V {
	return must.Must1(GetOrElse(iter, index, v))
}

// 指定した位置の値のポインタを返す。無い場合はnilを返す。
func GetOrNil[V any](iter Iter[V], index int) (*V, error) {
	v, ok, err := Get(iter, index)
	if !ok {
		return nil, err
	}
	return &v, nil
}

// 指定した位置の値のポインタを返す。無い場合はnilを返す。実行中にエラーが起きた場合 panic する。
func MustGetOrNil[V any](iter Iter[V], index int) *V {
	return must.Must1(GetOrNil(iter, index))
}

// 指定した位置の値を返す。無い場合はゼロ値を返す。
func GetOrZero[V any](iter Iter[V], index int) (V, error) {
	v, ok, err := Get(iter, index)
	if !ok {
		return *new(V), err
	}
	return v, nil
}

// 指定した位置の値を返す。無い場合はゼロ値を返す。実行中にエラーが起きた場合 panic する。
func MustGetOrZero[V any](iter Iter[V], index int) V {
	return must.Must1(GetOrZero(iter, index))
}

// 指定した位置の要素を返す。無い場合は関数の実行結果を返す。
func GetOrFunc[V any](iter Iter[V], index int, f func() (V, error)) (V, error) {
	v, ok, err := Get(iter, index)
	if err != nil {
		return *new(V), err
	}
	if ok {
		return v, nil
	}
	return f()
}

// 指定した位置の要素を返す。無い場合は関数の実行結果を返す。実行中にエラーが起きた場合 panic する。
func MustGetOrFunc[V any](iter Iter[V], index int, f func() (V, error)) V {
	return must.Must1(GetOrFunc(iter, index, f))
}

// 先頭の値を返す。
func GetFirst[V any](iter Iter[V]) (V, bool, error) {
	return Get(iter, 0)
}

// 先頭の値を返す。実行中にエラーが起きた場合 panic する。
func MustGetFirst[V any](iter Iter[V]) (V, bool) {
	return must.Must2(GetFirst(iter))
}

// 先頭の値を返す。無い場合はvを返す。
func GetFirstOrElse[V any](iter Iter[V], v V) (V, error) {
	v, ok, err := GetFirst(iter)
	if !ok {
		return *new(V), err
	}
	return v, nil
}

// 先頭の値を返す。無い場合はvを返す。実行中にエラーが起きた場合 panic する。
func MustGetFirstOrElse[V any](iter Iter[V], v V) V {
	return must.Must1(GetFirstOrElse(iter, v))
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

// 終端の値を返す。実行中にエラーが起きた場合 panic する。
func MustGetLast[V any](iter Iter[V]) (V, bool) {
	return must.Must2(GetLast(iter))
}

// 終端の値を返す。無い場合はvを返す。
func GetLastOrElse[V any](iter Iter[V], v V) (V, error) {
	v, ok, err := GetLast(iter)
	if !ok {
		return *new(V), err
	}
	return v, nil
}

// 終端の値を返す。無い場合はvを返す。実行中にエラーが起きた場合 panic する。
func MustGetLastOrElse[V any](iter Iter[V], v V) V {
	return must.Must1(GetLastOrElse(iter, v))
}

// イテレータの末尾に他のイテレータを結合する。
func Concat[V any](iter1 Iter[V], iter2 ...Iter[V]) Iter[V] {
	cursor := 0
	iters := append([]Iter[V]{iter1}, iter2...)

	return FromFunc(func(ctx Context) (V, bool) {
		if cursor >= len(iters) {
			return *new(V), false
		}
		v, ok := iters[cursor].Next()
		if !ok {
			ctx.SetErr(iters[cursor].Err())
			cursor++
		}
		return v, ok
	})
}

// イテレータの末尾に値を追加する。
func Append[V any](iter Iter[V], v ...V) Iter[V] {
	return Concat(iter, FromSlice(v))
}

// 指定した位置に値を追加する。
func Insert[V any](iter Iter[V], index int, v ...V) Iter[V] {
	i := 0
	return FromFunc(func(ctx Context) (V, bool) {
		if i >= index && i < index+len(v) {
			v1 := v[i-index]
			i++
			return v1, true
		}
		v, ok := iter.Next()
		if !ok {
			ctx.SetErr(iter.Err())
		}
		i++
		return v, ok
	})
}

// 指定した位置の値を削除する。
func Remove[V any](iter Iter[V], index int) Iter[V] {
	i := 0
	return FromFunc(func(ctx Context) (V, bool) {
		for {
			v, ok := iter.Next()
			if !ok {
				ctx.SetErr(iter.Err())
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

// 値ごとに関数を実行する。実行中にエラーが起きた場合 panic する。
func MustForEach[V any](iter Iter[V], f func(V) error) {
	must.Must0(ForEach(iter, f))
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

// 他のイテレータと関数で比較し、一致していたらtrueを返す。実行中にエラーが起きた場合 panic する。
func MustEqualBy[V any](iter1 Iter[V], iter2 Iter[V], f func(V, V) (bool, error)) bool {
	return must.Must1(EqualBy(iter1, iter2, f))
}

// 他のイテレータと一致していたらtrueを返す。
func Equal[V comparable](iter1 Iter[V], iter2 Iter[V]) (bool, error) {
	return EqualBy(iter1, iter2, func(v1 V, v2 V) (bool, error) { return v1 == v2, nil })
}

// 他のイテレータと一致していたらtrueを返す。実行中にエラーが起きた場合 panic する。
func MustEqual[V comparable](iter1 Iter[V], iter2 Iter[V]) bool {
	return must.Must1(Equal(iter1, iter2))
}

// 条件を満たす値の数を返す。
func CountBy[V any](iter Iter[V], f func(V) (bool, error)) (int, error) {
	return Len(FilterBy(iter, f))
}

// 条件を満たす値の数を返す。実行中にエラーが起きた場合 panic する。
func MustCountBy[V any](iter Iter[V], f func(V) (bool, error)) int {
	return must.Must1(CountBy(iter, f))
}

// 一致する値の数を返す。
func Count[V comparable](iter Iter[V], v V) (int, error) {
	return Len(Filter(iter, v))
}

// 一致する値の数を返す。実行中にエラーが起きた場合 panic する。
func MustCount[V comparable](iter Iter[V], v V) int {
	return must.Must1(Count(iter, v))
}

// 位置のイテレータを返す。
func Indices[V any](iter Iter[V]) Iter[int] {
	i := 0
	return FromFunc(func(ctx Context) (int, bool) {
		if _, ok := iter.Next(); !ok {
			ctx.SetErr(iter.Err())
			return 0, false
		}
		v := i
		i++
		return v, true
	})
}

// 値を変換したイテレータを返す。
func Map[V1 any, V2 any](iter Iter[V1], f func(V1) (V2, error)) Iter[V2] {
	return FromFunc(func(ctx Context) (V2, bool) {
		v1, ok := iter.Next()
		if !ok {
			ctx.SetErr(iter.Err())
			return *new(V2), false
		}
		v2, err := f(v1)
		if err != nil {
			ctx.SetErr(err)
			return *new(V2), false
		}
		return v2, true
	})
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

// 値を順に演算する。実行中にエラーが起きた場合 panic する。
func MustReduce[V any](iter Iter[V], f func(V, V) (V, error)) V {
	return must.Must1(Reduce(iter, f))
}

// 値の合計を演算する。
func Sum[V constraints.Ordered | constraints.Complex](iter Iter[V]) (V, error) {
	return Reduce(iter, func(sum V, v V) (V, error) {
		return sum + v, nil
	})
}

// 値の合計を演算する。実行中にエラーが起きた場合 panic する。
func MustSum[V constraints.Ordered | constraints.Complex](iter Iter[V]) V {
	return must.Must1(Sum(iter))
}

// 値を変換して合計を演算する。
func SumBy[V1 any, V2 constraints.Ordered | constraints.Complex](iter Iter[V1], f func(V1) (V2, error)) (V2, error) {
	return Sum(Map(iter, f))
}

// 値を変換して合計を演算する。実行中にエラーが起きた場合 panic する。
func MustSumBy[V1 any, V2 constraints.Ordered | constraints.Complex](iter Iter[V1], f func(V1) (V2, error)) V2 {
	return must.Must1(SumBy(iter, f))
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

// 最大の値を返す。実行中にエラーが起きた場合 panic する。
func MustMax[V constraints.Ordered](iter Iter[V]) V {
	return must.Must1(Max(iter))
}

// 値を変換して最大の値を返す。
func MaxBy[V1 any, V2 constraints.Ordered](iter Iter[V1], f func(V1) (V2, error)) (V2, error) {
	return Max(Map(iter, f))
}

// 値を変換して最大の値を返す。実行中にエラーが起きた場合 panic する。
func MustMaxBy[V1 any, V2 constraints.Ordered](iter Iter[V1], f func(V1) (V2, error)) V2 {
	return must.Must1(MaxBy(iter, f))
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

// 最小の値を返す。実行中にエラーが起きた場合 panic する。
func MustMin[V constraints.Ordered](iter Iter[V]) V {
	return must.Must1(Min(iter))
}

// 値を変換して最小の値を返す。
func MinBy[V1 any, V2 constraints.Ordered](iter Iter[V1], f func(V1) (V2, error)) (V2, error) {
	return Min(Map(iter, f))
}

// 値を変換して最小の値を返す。実行中にエラーが起きた場合 panic する。
func MustMinBy[V1 any, V2 constraints.Ordered](iter Iter[V1], f func(V1) (V2, error)) V2 {
	return must.Must1(MinBy(iter, f))
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

// 初期値と値を順に演算する。実行中にエラーが起きた場合 panic する。
func MustFold[V1 any, V2 any](iter Iter[V1], v V2, f func(V2, V1) (V2, error)) V2 {
	return must.Must1(Fold(iter, v, f))
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

// 条件を満たす最初の値の位置を返す。実行中にエラーが起きた場合 panic する。
func MustIndexBy[V any](iter Iter[V], f func(V) (bool, error)) int {
	return must.Must1(IndexBy(iter, f))
}

// 一致する最初の値の位置を返す。
func Index[V comparable](iter Iter[V], v V) (int, error) {
	return IndexBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 一致する最初の値の位置を返す。実行中にエラーが起きた場合 panic する。
func MustIndex[V comparable](iter Iter[V], v V) int {
	return must.Must1(Index(iter, v))
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

// 条件を満たす最後の値の位置を返す。実行中にエラーが起きた場合 panic する。
func MustLastIndexBy[V any](iter Iter[V], f func(V) (bool, error)) int {
	return must.Must1(LastIndexBy(iter, f))
}

// 一致する最後の値の位置を返す。
func LastIndex[V comparable](iter Iter[V], v V) (int, error) {
	return LastIndexBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 一致する最後の値の位置を返す。実行中にエラーが起きた場合 panic する。
func MustLastIndex[V comparable](iter Iter[V], v V) int {
	return must.Must1(LastIndex(iter, v))
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

// 条件を満たす値を返す。実行中にエラーが起きた場合 panic する。
func MustFindBy[V any](iter Iter[V], f func(V) (bool, error)) (V, bool) {
	return must.Must2(FindBy(iter, f))
}

// 一致する値を返す。
func Find[V comparable](iter Iter[V], v V) (V, bool, error) {
	return FindBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 一致する値を返す。実行中にエラーが起きた場合 panic する。
func MustFind[V comparable](iter Iter[V], v V) (V, bool) {
	return must.Must2(Find(iter, v))
}

// 条件を満たす値が存在したらtrueを返す。
func ExistsBy[V any](iter Iter[V], f func(V) (bool, error)) (bool, error) {
	index, err := IndexBy(iter, f)
	return index >= 0, err
}

// 条件を満たす値が存在したらtrueを返す。実行中にエラーが起きた場合 panic する。
func MustExistsBy[V any](iter Iter[V], f func(V) (bool, error)) bool {
	return must.Must1(ExistsBy(iter, f))
}

// 一致する値が存在したらtrueを返す。
func Exists[V comparable](iter Iter[V], v V) (bool, error) {
	return ExistsBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 一致する値が存在したらtrueを返す。実行中にエラーが起きた場合 panic する。
func MustExists[V comparable](iter Iter[V], v V) bool {
	return must.Must1(Exists(iter, v))
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

// すべての値が条件を満たせばtrueを返す。実行中にエラーが起きた場合 panic する。
func MustForAllBy[V any](iter Iter[V], f func(V) (bool, error)) bool {
	return must.Must1(ForAllBy(iter, f))
}

// すべての値が一致したらtrueを返す。
func ForAll[V comparable](iter Iter[V], v V) (bool, error) {
	return ForAllBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// すべての値が一致したらtrueを返す。実行中にエラーが起きた場合 panic する。
func MustForAll[V comparable](iter Iter[V], v V) bool {
	return must.Must1(ForAll(iter, v))
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

// 他のイテレータの値がひとつでも存在していたらtrueを返す。実行中にエラーが起きた場合 panic する。
func MustContainsAny[V comparable](iter Iter[V], subset Iter[V]) bool {
	return must.Must1(ContainsAny(iter, subset))
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

// 他のイテレータの値がすべて存在していたらtrueを返す。実行中にエラーが起きた場合 panic する。
func MustContainsAll[V comparable](iter Iter[V], subset Iter[V]) bool {
	return must.Must1(ContainsAll(iter, subset))
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

// 先頭が他のイテレータと一致していたらtrueを返す。実行中にエラーが起きた場合 panic する。
func MustStartsWith[V comparable](iter Iter[V], subset Iter[V]) bool {
	return must.Must1(StartsWith(iter, subset))
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

// 終端が他のイテレータと一致していたらtrueを返す。実行中にエラーが起きた場合 panic する。
func MustEndsWith[V comparable](iter Iter[V], subset Iter[V]) bool {
	return must.Must1(EndsWith(iter, subset))
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
	return FromFunc(func(ctx Context) (V, bool) {
		for {
			v, ok := iter.Next()
			if !ok {
				ctx.SetErr(iter.Err())
				return *new(V), false
			}
			ok, err := f(v)
			if err != nil {
				ctx.SetErr(err)
				return *new(V), false
			}
			if !ok {
				continue
			}
			return v, true
		}
	})
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

// 条件を満たす値の直前で分割したふたつのイテレータを返す。実行中にエラーが起きた場合 panic する。
func MustSplitBy[V any](iter Iter[V], f func(V) (bool, error)) (Iter[V], Iter[V]) {
	return must.Must2(SplitBy(iter, f))
}

// 一致する値の直前で分割したふたつのイテレータを返す。
func Split[V comparable](iter Iter[V], v V) (Iter[V], Iter[V], error) {
	return SplitBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 一致する値の直前で分割したふたつのイテレータを返す。実行中にエラーが起きた場合 panic する。
func MustSplit[V comparable](iter Iter[V], v V) (Iter[V], Iter[V]) {
	return must.Must2(Split(iter, v))
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

// 条件を満たす値の直後で分割したふたつのイテレータを返す。実行中にエラーが起きた場合 panic する。
func MustSplitAfterBy[V any](iter Iter[V], f func(V) (bool, error)) (Iter[V], Iter[V]) {
	return must.Must2(SplitAfterBy(iter, f))
}

// 一致する値の直後で分割したふたつのイテレータを返す。
func SplitAfter[V comparable](iter Iter[V], v V) (Iter[V], Iter[V], error) {
	return SplitAfterBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 一致する値の直後で分割したふたつのイテレータを返す。実行中にエラーが起きた場合 panic する。
func MustSplitAfter[V comparable](iter Iter[V], v V) (Iter[V], Iter[V]) {
	return must.Must2(SplitAfter(iter, v))
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

// 条件を満たすイテレータと満たさないイテレータを返す。実行中にエラーが起きた場合 panic する。
func MustPartitionBy[V any](iter Iter[V], f func(V) (bool, error)) (Iter[V], Iter[V]) {
	return must.Must2(PartitionBy(iter, f))
}

// 値の一致するイテレータと一致しないイテレータを返す。
func Partition[V comparable](iter Iter[V], v V) (Iter[V], Iter[V], error) {
	return PartitionBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 値の一致するイテレータと一致しないイテレータを返す。実行中にエラーが起きた場合 panic する。
func MustPartition[V comparable](iter Iter[V], v V) (Iter[V], Iter[V]) {
	return must.Must2(Partition(iter, v))
}

// 条件を満たし続ける先頭の値のイテレータを返す。
func TakeWhileBy[V any](iter Iter[V], f func(V) (bool, error)) Iter[V] {
	done := false
	return FromFunc(func(ctx Context) (V, bool) {
		if done {
			return *new(V), false
		}

		v, ok := iter.Next()
		if !ok {
			ctx.SetErr(iter.Err())
			return *new(V), false
		}

		ok, err := f(v)
		if err != nil {
			ctx.SetErr(err)
			return *new(V), false
		}
		if !ok {
			done = true
			return *new(V), false
		}
		return v, true
	})
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
	return FromFunc(func(ctx Context) (V, bool) {
		if done {
			return iter.Next()
		}

		for {
			v, ok := iter.Next()
			if !ok {
				ctx.SetErr(iter.Err())
				return *new(V), false
			}

			ok, err := f(v)
			if err != nil {
				ctx.SetErr(err)
				return *new(V), false
			}
			if ok {
				continue
			}
			done = true
			return v, true
		}
	})
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

// 条件を満たし続ける先頭部分と残りの部分、ふたつのイテレータを返す。実行中にエラーが起きた場合 panic する。
func MustSpanBy[V any](iter Iter[V], f func(V) (bool, error)) (Iter[V], Iter[V]) {
	return must.Must2(SpanBy(iter, f))
}

// 一致し続ける先頭部分と残りの部分、ふたつのイテレータを返す。
func Span[V comparable](iter Iter[V], v V) (Iter[V], Iter[V], error) {
	return SpanBy(iter, func(v1 V) (bool, error) { return v1 == v, nil })
}

// 一致し続ける先頭部分と残りの部分、ふたつのイテレータを返す。実行中にエラーが起きた場合 panic する。
func MustSpan[V comparable](iter Iter[V], v V) (Iter[V], Iter[V]) {
	return must.Must2(Span(iter, v))
}

// ゼロ値を除いたイテレータを返す。
func Clean[V comparable](iter Iter[V]) Iter[V] {
	return FilterNot(iter, *new(V))
}

// 重複を除いたイテレータを返す。
func Distinct[V comparable](iter Iter[V]) Iter[V] {
	m := map[V]struct{}{}
	return FromFunc(func(ctx Context) (V, bool) {
		for {
			v, ok := iter.Next()
			if !ok {
				ctx.SetErr(iter.Err())
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
func Collect[V1 any, V2 any](iter Iter[V1], f func(V1) (V2, bool, error)) Iter[V2] {
	return FromFunc(func(ctx Context) (V2, bool) {
		for {
			v1, ok := iter.Next()
			if !ok {
				ctx.SetErr(iter.Err())
				return *new(V2), false
			}

			v2, ok, err := f(v1)
			if err != nil {
				ctx.SetErr(err)
				return *new(V2), false
			}
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
	return FromFunc(func(ctx Context) (tuple.T2[V, int], bool) {
		v, ok := iter.Next()
		if !ok {
			ctx.SetErr(iter.Err())
			return tuple.NewT2(*new(V), 0), false
		}
		t := tuple.NewT2(v, i)
		i++
		return t, true
	})
}

// ２つのイテレータの同じ位置の値をペアにしたイテレータを返す。
func Zip2[V1 any, V2 any](iter1 Iter[V1], iter2 Iter[V2]) Iter[tuple.T2[V1, V2]] {
	return FromFunc(func(ctx Context) (tuple.T2[V1, V2], bool) {
		v1, ok1 := iter1.Next()
		v2, ok2 := iter2.Next()
		if ok1 && ok2 {
			return tuple.NewT2(v1, v2), true
		}
		errs := make([]error, 0, 2)
		if err := iter1.Err(); err != nil {
			errs = append(errs, err)
		}
		if err := iter2.Err(); err != nil {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			var b []byte
			for i, err := range errs {
				if i > 0 {
					b = append(b, '\n')
				}
				b = append(b, err.Error()...)
			}
			ctx.SetErr(errors.New(string(b)))
		}
		return tuple.NewT2(*new(V1), *new(V2)), false
	})
}

// ３つのイテレータの同じ位置の値をペアにしたイテレータを返す。
func Zip3[V1 any, V2 any, V3 any](iter1 Iter[V1], iter2 Iter[V2], iter3 Iter[V3]) Iter[tuple.T3[V1, V2, V3]] {
	return FromFunc(func(ctx Context) (tuple.T3[V1, V2, V3], bool) {
		v1, ok1 := iter1.Next()
		v2, ok2 := iter2.Next()
		v3, ok3 := iter3.Next()
		if ok1 && ok2 && ok3 {
			return tuple.NewT3(v1, v2, v3), true
		}
		errs := make([]error, 0, 3)
		if err := iter1.Err(); err != nil {
			errs = append(errs, err)
		}
		if err := iter2.Err(); err != nil {
			errs = append(errs, err)
		}
		if err := iter3.Err(); err != nil {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			var b []byte
			for i, err := range errs {
				if i > 0 {
					b = append(b, '\n')
				}
				b = append(b, err.Error()...)
			}
			ctx.SetErr(errors.New(string(b)))
		}
		return tuple.NewT3(*new(V1), *new(V2), *new(V3)), false
	})
}

// ４つのイテレータの同じ位置の値をペアにしたイテレータを返す。
func Zip4[V1 any, V2 any, V3 any, V4 any](iter1 Iter[V1], iter2 Iter[V2], iter3 Iter[V3], iter4 Iter[V4]) Iter[tuple.T4[V1, V2, V3, V4]] {
	return FromFunc(func(ctx Context) (tuple.T4[V1, V2, V3, V4], bool) {
		v1, ok1 := iter1.Next()
		v2, ok2 := iter2.Next()
		v3, ok3 := iter3.Next()
		v4, ok4 := iter4.Next()
		if ok1 && ok2 && ok3 && ok4 {
			return tuple.NewT4(v1, v2, v3, v4), true
		}
		errs := make([]error, 0, 4)
		if err := iter1.Err(); err != nil {
			errs = append(errs, err)
		}
		if err := iter2.Err(); err != nil {
			errs = append(errs, err)
		}
		if err := iter3.Err(); err != nil {
			errs = append(errs, err)
		}
		if err := iter4.Err(); err != nil {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			var b []byte
			for i, err := range errs {
				if i > 0 {
					b = append(b, '\n')
				}
				b = append(b, err.Error()...)
			}
			ctx.SetErr(errors.New(string(b)))
		}
		return tuple.NewT4(*new(V1), *new(V2), *new(V3), *new(V4)), false
	})
}

// ５つのイテレータの同じ位置の値をペアにしたイテレータを返す。
func Zip5[V1 any, V2 any, V3 any, V4 any, V5 any](iter1 Iter[V1], iter2 Iter[V2], iter3 Iter[V3], iter4 Iter[V4], iter5 Iter[V5]) Iter[tuple.T5[V1, V2, V3, V4, V5]] {
	return FromFunc(func(ctx Context) (tuple.T5[V1, V2, V3, V4, V5], bool) {
		v1, ok1 := iter1.Next()
		v2, ok2 := iter2.Next()
		v3, ok3 := iter3.Next()
		v4, ok4 := iter4.Next()
		v5, ok5 := iter5.Next()
		if ok1 && ok2 && ok3 && ok4 && ok5 {
			return tuple.NewT5(v1, v2, v3, v4, v5), true
		}
		errs := make([]error, 0, 5)
		if err := iter1.Err(); err != nil {
			errs = append(errs, err)
		}
		if err := iter2.Err(); err != nil {
			errs = append(errs, err)
		}
		if err := iter3.Err(); err != nil {
			errs = append(errs, err)
		}
		if err := iter4.Err(); err != nil {
			errs = append(errs, err)
		}
		if err := iter5.Err(); err != nil {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			var b []byte
			for i, err := range errs {
				if i > 0 {
					b = append(b, '\n')
				}
				b = append(b, err.Error()...)
			}
			ctx.SetErr(errors.New(string(b)))
		}
		return tuple.NewT5(*new(V1), *new(V2), *new(V3), *new(V4), *new(V5)), false
	})
}

// ６つのイテレータの同じ位置の値をペアにしたイテレータを返す。
func Zip6[V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](iter1 Iter[V1], iter2 Iter[V2], iter3 Iter[V3], iter4 Iter[V4], iter5 Iter[V5], iter6 Iter[V6]) Iter[tuple.T6[V1, V2, V3, V4, V5, V6]] {
	return FromFunc(func(ctx Context) (tuple.T6[V1, V2, V3, V4, V5, V6], bool) {
		v1, ok1 := iter1.Next()
		v2, ok2 := iter2.Next()
		v3, ok3 := iter3.Next()
		v4, ok4 := iter4.Next()
		v5, ok5 := iter5.Next()
		v6, ok6 := iter6.Next()
		if ok1 && ok2 && ok3 && ok4 && ok5 && ok6 {
			return tuple.NewT6(v1, v2, v3, v4, v5, v6), true
		}
		errs := make([]error, 0, 6)
		if err := iter1.Err(); err != nil {
			errs = append(errs, err)
		}
		if err := iter2.Err(); err != nil {
			errs = append(errs, err)
		}
		if err := iter3.Err(); err != nil {
			errs = append(errs, err)
		}
		if err := iter4.Err(); err != nil {
			errs = append(errs, err)
		}
		if err := iter5.Err(); err != nil {
			errs = append(errs, err)
		}
		if err := iter6.Err(); err != nil {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			var b []byte
			for i, err := range errs {
				if i > 0 {
					b = append(b, '\n')
				}
				b = append(b, err.Error()...)
			}
			ctx.SetErr(errors.New(string(b)))
		}
		return tuple.NewT6(*new(V1), *new(V2), *new(V3), *new(V4), *new(V5), *new(V6)), false
	})
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

// 値ごとに関数の返すキーでグルーピングしたマップを返す。実行中にエラーが起きた場合 panic する。
func MustGroupBy[K comparable, V any](iter Iter[V], f func(V) (K, error)) map[K][]V {
	return must.Must1(GroupBy(iter, f))
}

// 平坦化したイテレータを返す。
func Flatten[V any](iter Iter[Iter[V]]) Iter[V] {
	sub, ok := iter.Next()
	return FromFunc(func(ctx Context) (V, bool) {
		for {
			if !ok {
				ctx.SetErr(iter.Err())
				return *new(V), false
			}

			if v, ok := sub.Next(); ok {
				return v, true
			}
			if err := sub.Err(); err != nil {
				ctx.SetErr(err)
				return *new(V), false
			}

			sub, ok = iter.Next()
		}
	})
}

// 値をイテレータに変換し、それらを結合したイテレータを返す。
func FlatMap[V1 any, V2 any](iter Iter[V1], f func(V1) (Iter[V2], error)) Iter[V2] {
	return Flatten(Map(iter, f))
}

// 値のあいだにseparatorを挿入したイテレータを返す。
func Join[V any](iter Iter[V], separator V) Iter[V] {
	return Drop(FlatMap(iter, func(v V) (Iter[V], error) { return From(separator, v), nil }), 1)
}
