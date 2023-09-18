package slices

import (
	"math/rand"

	"github.com/thamaji/gu/must"
	"github.com/thamaji/gu/tuple"
	"golang.org/x/exp/constraints"
)

// 値の数を返す。
func Len[V any](slice []V) int {
	return len(slice)
}

// 空のときtrueを返す。
func IsEmpty[V any](slice []V) bool {
	return len(slice) == 0
}

// 指定した位置の要素を返す。
func Get[V any](slice []V, index int) (V, bool) {
	if index < len(slice) {
		return slice[index], true
	}
	return *new(V), false
}

// 指定した位置の要素を返す。無い場合はvを返す。
func GetOrElse[V any](slice []V, index int, v V) V {
	if index < len(slice) {
		return slice[index]
	}
	return v
}

// 指定した位置の要素のポインタを返す。無い場合はnilを返す。
func GetOrNil[V any](slice []V, index int) *V {
	if index < len(slice) {
		return &slice[index]
	}
	return nil
}

// 指定した位置の要素を返す。無い場合はゼロ値を返す。
func GetOrZero[V any](slice []V, index int) V {
	if index < len(slice) {
		return slice[index]
	}
	return *new(V)
}

// 指定した位置の要素を返す。無い場合は関数の実行結果を返す。
func GetOrFunc[V any](slice []V, index int, f func() (V, error)) (V, error) {
	if index < len(slice) {
		return slice[index], nil
	}
	return f()
}

// 指定した位置の要素を返す。無い場合は関数の実行結果を返す。実行中にエラーが起きた場合 panic する。
func MustGetOrFunc[V any](slice []V, index int, f func() (V, error)) V {
	return must.Must1(GetOrFunc(slice, index, f))
}

// 先頭の要素を返す。
func GetFirst[V any](slice []V) (V, bool) {
	if len(slice) > 0 {
		return slice[0], true
	}
	return *new(V), false
}

// 先頭の要素を返す。無い場合はvを返す。
func GetFirstOrElse[V any](slice []V, v V) V {
	if len(slice) > 0 {
		return slice[0]
	}
	return v
}

// 終端の要素を返す。
func GetLast[V any](slice []V) (V, bool) {
	if len(slice) > 0 {
		return slice[len(slice)-1], true
	}
	return *new(V), false
}

// 終端の要素を返す。無い場合はvを返す。
func GetLastOrElse[V any](slice []V, v V) V {
	if len(slice) > 0 {
		return slice[len(slice)-1]
	}
	return v
}

// スライスの末尾に他のスライスを結合する。
func Concat[S ~[]V, V any](slice1 S, slice2 []V) S {
	return append(slice1, slice2...)
}

// スライスの末尾に値を追加する。
func Append[S ~[]V, V any](slice S, v ...V) S {
	return append(slice, v...)
}

// 指定した位置に値を追加する。
func Insert[S ~[]V, V any](slice S, index int, v ...V) S {
	return append(slice[:index], append(v, slice[index:]...)...)
}

// 指定した位置の要素を削除する。
func Remove[S ~[]V, V any](slice S, index int) S {
	dst := make(S, 0, len(slice)-1)
	for i := range slice {
		if i == index {
			continue
		}
		dst = append(dst, slice[i])
	}
	return dst
}

// 要素をすべてコピーしたスライスを返す。
func Clone[S ~[]V, V any](slice S) S {
	clone := make(S, len(slice))
	copy(clone, slice)
	return clone
}

// 要素を１つランダムに返す。
func Sample[V any](slice []V, r *rand.Rand) V {
	return slice[r.Intn(len(slice))]
}

// 逆順にしたスライスを返す。
func Reverse[S ~[]V, V any](slice []V) S {
	dst := make(S, 0, len(slice))
	for i := len(slice) - 1; i >= 0; i-- {
		dst = append(dst, slice[i])
	}
	return dst
}

// 値ごとに関数を実行する。
func ForEach[V any](slice []V, f func(V) error) error {
	for i := range slice {
		if err := f(slice[i]); err != nil {
			return err
		}
	}
	return nil
}

// 値ごとに関数を実行する。実行中にエラーが起きた場合 panic する。
func MustForEach[V any](slice []V, f func(V) error) {
	must.Must0(ForEach(slice, f))
}

// 他のスライスと関数で比較し、一致していたらtrueを返す。
func EqualBy[V any](slice1 []V, slice2 []V, f func(V, V) (bool, error)) (bool, error) {
	if len(slice1) != len(slice2) {
		return false, nil
	}
	for i := range slice1 {
		ok, err := f(slice1[i], slice2[i])
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

// 他のスライスと関数で比較し、一致していたらtrueを返す。実行中にエラーが起きた場合 panic する。
func MustEqualBy[V any](slice1 []V, slice2 []V, f func(V, V) (bool, error)) bool {
	return must.Must1(EqualBy(slice1, slice2, f))
}

// 他のスライスと一致していたらtrueを返す。
func Equal[V comparable](slice1 []V, slice2 []V) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

// 条件を満たす値の数を返す。
func CountBy[V any](slice []V, f func(V) (bool, error)) (int, error) {
	c := 0
	for i := range slice {
		ok, err := f(slice[i])
		if err != nil {
			return 0, err
		}
		if ok {
			c++
		}
	}
	return c, nil
}

// 条件を満たす値の数を返す。実行中にエラーが起きた場合 panic する。
func MustCountBy[V any](slice []V, f func(V) (bool, error)) int {
	return must.Must1(CountBy(slice, f))
}

// 一致する値の数を返す。
func Count[V comparable](slice []V, v V) int {
	c := 0
	for i := range slice {
		if slice[i] == v {
			c++
		}
	}
	return c
}

// 位置のスライスを返す。
func Indices[V any](slice []V) []int {
	dst := make([]int, len(slice))
	for i := range slice {
		dst[i] = i
	}
	return dst
}

// 値を変換したスライスを返す。
func Map[V1 any, V2 any](slice []V1, f func(V1) (V2, error)) ([]V2, error) {
	dst := make([]V2, len(slice))
	for i := range slice {
		v2, err := f(slice[i])
		if err != nil {
			return nil, err
		}
		dst[i] = v2
	}
	return dst, nil
}

// 値を変換したスライスを返す。実行中にエラーが起きた場合 panic する。
func MustMap[V1 any, V2 any](slice []V1, f func(V1) (V2, error)) []V2 {
	return must.Must1(Map(slice, f))
}

// 値を順に演算する。
func Reduce[V any](slice []V, f func(V, V) (V, error)) (V, error) {
	var v V
	var err error
	if len(slice) > 0 {
		v = slice[0]
		for i := 1; i < len(slice); i++ {
			v, err = f(v, slice[i])
			if err != nil {
				return *new(V), err
			}
		}
	}
	return v, nil
}

// 値を順に演算する。実行中にエラーが起きた場合 panic する。
func MustReduce[V any](slice []V, f func(V, V) (V, error)) V {
	return must.Must1(Reduce(slice, f))
}

// 値の合計を演算する。
func Sum[V constraints.Ordered | constraints.Complex](slice []V) V {
	var sum V
	if len(slice) > 0 {
		sum = slice[0]
		for i := 1; i < len(slice); i++ {
			sum += slice[i]
		}
	}
	return sum
}

// 値を変換して合計を演算する。
func SumBy[V1 any, V2 constraints.Ordered | constraints.Complex](slice []V1, f func(V1) (V2, error)) (V2, error) {
	var sum V2
	var err error
	if len(slice) > 0 {
		sum, err = f(slice[0])
		if err != nil {
			return *new(V2), err
		}

		for i := 1; i < len(slice); i++ {
			v2, err := f(slice[i])
			if err != nil {
				return *new(V2), err
			}
			sum += v2
		}
	}
	return sum, nil
}

// 値を変換して合計を演算する。実行中にエラーが起きた場合 panic する。
func MustSumBy[V1 any, V2 constraints.Ordered | constraints.Complex](slice []V1, f func(V1) (V2, error)) V2 {
	return must.Must1(SumBy(slice, f))
}

// 最大の値を返す。
func Max[V constraints.Ordered](slice []V) V {
	var max V
	if len(slice) > 0 {
		max = slice[0]
		for i := 1; i < len(slice); i++ {
			if max < slice[i] {
				max = slice[i]
			}
		}
	}
	return max
}

// 値を変換して最大の値を返す。
func MaxBy[V1 any, V2 constraints.Ordered](slice []V1, f func(V1) (V2, error)) (V2, error) {
	var max V2
	var err error
	if len(slice) > 0 {
		max, err = f(slice[0])
		if err != nil {
			return *new(V2), err
		}

		for i := 1; i < len(slice); i++ {
			v2, err := f(slice[i])
			if err != nil {
				return *new(V2), err
			}
			if max < v2 {
				max = v2
			}
		}
	}
	return max, nil
}

// 値を変換して最大の値を返す。実行中にエラーが起きた場合 panic する。
func MustMaxBy[V1 any, V2 constraints.Ordered](slice []V1, f func(V1) (V2, error)) V2 {
	return must.Must1(MaxBy(slice, f))
}

// 最小の値を返す。
func Min[V constraints.Ordered](slice []V) V {
	var min V
	if len(slice) > 0 {
		min = slice[0]
		for i := 1; i < len(slice); i++ {
			if min > slice[i] {
				min = slice[i]
			}
		}
	}
	return min
}

// 値を変換して最小の値を返す。
func MinBy[V1 any, V2 constraints.Ordered](slice []V1, f func(V1) (V2, error)) (V2, error) {
	var min V2
	var err error
	if len(slice) > 0 {
		min, err = f(slice[0])
		if err != nil {
			return *new(V2), err
		}

		for i := 1; i < len(slice); i++ {
			v, err := f(slice[i])
			if err != nil {
				return *new(V2), err
			}
			if min > v {
				min = v
			}
		}
	}
	return min, nil
}

// 値を変換して最小の値を返す。実行中にエラーが起きた場合 panic する。
func MustMinBy[V1 any, V2 constraints.Ordered](slice []V1, f func(V1) (V2, error)) V2 {
	return must.Must1(MinBy(slice, f))
}

// 初期値と値を順に演算する。
func Fold[V1 any, V2 any](slice []V1, v V2, f func(V2, V1) (V2, error)) (V2, error) {
	var err error
	for i := range slice {
		v, err = f(v, slice[i])
		if err != nil {
			return *new(V2), err
		}
	}
	return v, nil
}

// 初期値と値を順に演算する。実行中にエラーが起きた場合 panic する。
func MustFold[V1 any, V2 any](slice []V1, v V2, f func(V2, V1) (V2, error)) V2 {
	return must.Must1(Fold(slice, v, f))
}

// 条件を満たす最初の値の位置を返す。
func IndexBy[V any](slice []V, f func(V) (bool, error)) (int, error) {
	for i := range slice {
		ok, err := f(slice[i])
		if err != nil {
			return 0, err
		}
		if ok {
			return i, nil
		}
	}
	return -1, nil
}

// 条件を満たす最初の値の位置を返す。実行中にエラーが起きた場合 panic する。
func MustIndexBy[V any](slice []V, f func(V) (bool, error)) int {
	return must.Must1(IndexBy(slice, f))
}

// 一致する最初の値の位置を返す。
func Index[V comparable](slice []V, v V) int {
	for i := range slice {
		if slice[i] == v {
			return i
		}
	}
	return -1
}

// 条件を満たす最後の値の位置を返す。
func LastIndexBy[V any](slice []V, f func(V) (bool, error)) (int, error) {
	for i := len(slice) - 1; i >= 0; i-- {
		ok, err := f(slice[i])
		if err != nil {
			return 0, err
		}
		if ok {
			return i, nil
		}
	}
	return -1, nil
}

// 条件を満たす最後の値の位置を返す。実行中にエラーが起きた場合 panic する。
func MustLastIndexBy[V any](slice []V, f func(V) (bool, error)) int {
	return must.Must1(LastIndexBy(slice, f))
}

// 一致する最後の値の位置を返す。
func LastIndex[V comparable](slice []V, v V) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == v {
			return i
		}
	}
	return -1
}

// 条件を満たす値を返す。
func FindBy[V any](slice []V, f func(V) (bool, error)) (V, bool, error) {
	for i := range slice {
		ok, err := f(slice[i])
		if err != nil {
			return *new(V), false, err
		}
		if ok {
			return slice[i], true, nil
		}
	}
	return *new(V), false, nil
}

// 条件を満たす値を返す。実行中にエラーが起きた場合 panic する。
func MustFindBy[V any](slice []V, f func(V) (bool, error)) (V, bool) {
	return must.Must2(FindBy(slice, f))
}

// 一致する値を返す。
func Find[V comparable](slice []V, v V) (V, bool) {
	for i := range slice {
		if slice[i] == v {
			return slice[i], true
		}
	}
	return *new(V), false
}

// 条件を満たす値が存在したらtrueを返す。
func ExistsBy[V any](slice []V, f func(V) (bool, error)) (bool, error) {
	for i := range slice {
		ok, err := f(slice[i])
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// 条件を満たす値が存在したらtrueを返す。実行中にエラーが起きた場合 panic する。
func MustExistsBy[V any](slice []V, f func(V) (bool, error)) bool {
	return must.Must1(ExistsBy(slice, f))
}

// 一致する値が存在したらtrueを返す。
func Exists[V comparable](slice []V, v V) bool {
	for i := range slice {
		if slice[i] == v {
			return true
		}
	}
	return false
}

// すべての値が条件を満たせばtrueを返す。
func ForAllBy[V any](slice []V, f func(V) (bool, error)) (bool, error) {
	for i := range slice {
		ok, err := f(slice[i])
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

// すべての値が条件を満たせばtrueを返す。実行中にエラーが起きた場合 panic する。
func MustForAllBy[V any](slice []V, f func(V) (bool, error)) bool {
	return must.Must1(ForAllBy(slice, f))
}

// すべての値が一致したらtrueを返す。
func ForAll[V comparable](slice []V, v V) bool {
	for i := range slice {
		if slice[i] != v {
			return false
		}
	}
	return true
}

// 他のスライスの値がひとつでも存在していたらtrueを返す。
func ContainsAny[V comparable](slice []V, subset []V) bool {
	for i := range subset {
		for j := range slice {
			if subset[i] == slice[j] {
				return true
			}
		}
	}
	return false
}

// 他のスライスの値がすべて存在していたらtrueを返す。
func ContainsAll[V comparable](slice []V, subset []V) bool {
	for i := range subset {
		ok := false
		for j := range slice {
			if subset[i] == slice[j] {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

// 先頭が他のスライスと一致していたらtrueを返す。
func StartsWith[V comparable](slice []V, subset []V) bool {
	if len(subset) > len(slice) {
		return false
	}
	for i := range subset {
		if subset[i] != slice[i] {
			return false
		}
	}
	return true
}

// 終端が他のスライスと一致していたらtrueを返す。
func EndsWith[V comparable](slice []V, subset []V) bool {
	if len(subset) > len(slice) {
		return false
	}
	for i, j := len(subset)-1, len(slice)-1; i >= 0; i, j = i-1, j-1 {
		if subset[i] != slice[i] {
			return false
		}
	}
	return true
}

// ひとつめのoldをnewで置き換えたスライスを返す。
func Replace[S ~[]V, V comparable](slice S, old V, new V) S {
	done := false
	dst := make(S, len(slice))
	for i := range slice {
		if !done && slice[i] == old {
			dst[i] = new
			done = true
		} else {
			dst[i] = slice[i]
		}
	}
	return dst
}

// すべてのoldをnewで置き換えたスライスを返す。
func ReplaceAll[S ~[]V, V comparable](slice S, old V, new V) S {
	dst := make(S, len(slice))
	for i := range slice {
		if slice[i] == old {
			dst[i] = new
		} else {
			dst[i] = slice[i]
		}
	}
	return dst
}

// 条件を満たす値だけのスライスを返す。
func FilterBy[S ~[]V, V any](slice S, f func(V) (bool, error)) (S, error) {
	dst := make(S, 0, len(slice)/2)
	for i := range slice {
		ok, err := f(slice[i])
		if err != nil {
			return nil, err
		}
		if ok {
			dst = append(dst, slice[i])
		}
	}
	return dst, nil
}

// 条件を満たす値だけのスライスを返す。実行中にエラーが起きた場合 panic する。
func MustFilterBy[S ~[]V, V any](slice S, f func(V) (bool, error)) S {
	return must.Must1(FilterBy(slice, f))
}

// 一致する値だけのスライスを返す。
func Filter[S ~[]V, V comparable](slice S, v V) S {
	dst := make(S, 0, len(slice)/2)
	for i := range slice {
		if slice[i] == v {
			dst = append(dst, slice[i])
		}
	}
	return dst
}

// 条件を満たす値を除いたスライスを返す。
func FilterNotBy[S ~[]V, V any](slice S, f func(V) (bool, error)) (S, error) {
	dst := make(S, 0, len(slice)/2)
	for i := range slice {
		ok, err := f(slice[i])
		if err != nil {
			return nil, err
		}
		if !ok {
			dst = append(dst, slice[i])
		}
	}
	return dst, nil
}

// 条件を満たす値を除いたスライスを返す。実行中にエラーが起きた場合 panic する。
func MustFilterNotBy[S ~[]V, V any](slice S, f func(V) (bool, error)) S {
	return must.Must1(FilterNotBy(slice, f))
}

// 一致する値を除いたスライスを返す。
func FilterNot[S ~[]V, V comparable](slice S, v V) S {
	dst := make(S, 0, len(slice)/2)
	for i := range slice {
		if slice[i] != v {
			dst = append(dst, slice[i])
		}
	}
	return dst
}

// 条件を満たす値の直前で分割したふたつのスライスを返す。
func SplitBy[S ~[]V, V any](slice S, f func(V) (bool, error)) (S, S, error) {
	for i := range slice {
		ok, err := f(slice[i])
		if err != nil {
			return nil, nil, err
		}

		if ok {
			return slice[:i], slice[i:], nil
		}
	}
	return slice, nil, nil
}

// 条件を満たす値の直前で分割したふたつのスライスを返す。実行中にエラーが起きた場合 panic する。
func MustSplitBy[S ~[]V, V any](slice S, f func(V) (bool, error)) (S, S) {
	return must.Must2(SplitBy(slice, f))
}

// 一致する値の直前で分割したふたつのスライスを返す。
func Split[S ~[]V, V comparable](slice S, v V) (S, S) {
	for i := range slice {
		if slice[i] == v {
			return slice[:i], slice[i:]
		}
	}
	return slice, nil
}

// 条件を満たす値の直後で分割したふたつのスライスを返す。
func SplitAfterBy[S ~[]V, V any](slice S, f func(V) (bool, error)) (S, S, error) {
	for i := range slice {
		ok, err := f(slice[i])
		if err != nil {
			return nil, nil, err
		}

		if ok {
			return slice[:i+1], slice[i+1:], nil
		}
	}
	return slice, nil, nil
}

// 条件を満たす値の直後で分割したふたつのスライスを返す。実行中にエラーが起きた場合 panic する。
func MustSplitAfterBy[S ~[]V, V any](slice S, f func(V) (bool, error)) (S, S) {
	return must.Must2(SplitAfterBy(slice, f))
}

// 一致する値の直後で分割したふたつのスライスを返す。
func SplitAfter[S ~[]V, V comparable](slice S, v V) (S, S, error) {
	for i := range slice {
		if slice[i] == v {
			return slice[:i+1], slice[i+1:], nil
		}
	}
	return slice, nil, nil
}

// 一致する値の直後で分割したふたつのスライスを返す。実行中にエラーが起きた場合 panic する。
func MustSplitAfter[S ~[]V, V comparable](slice S, v V) (S, S) {
	return must.Must2(SplitAfter(slice, v))
}

// 条件を満たすスライスと満たさないスライスを返す。
func PartitionBy[S ~[]V, V any](slice S, f func(V) (bool, error)) (S, S, error) {
	dst1 := make(S, 0, len(slice)/2)
	dst2 := make(S, 0, len(slice)/2)
	for i := range slice {
		ok, err := f(slice[i])
		if err != nil {
			return nil, nil, err
		}
		if ok {
			dst1 = append(dst1, slice[i])
		} else {
			dst2 = append(dst2, slice[i])
		}
	}
	return dst1, dst2, nil
}

// 条件を満たすスライスと満たさないスライスを返す。
func MustPartitionBy[S ~[]V, V any](slice S, f func(V) (bool, error)) (S, S) {
	return must.Must2(PartitionBy(slice, f))
}

// 値の一致するスライスと一致しないスライスを返す。
func Partition[S ~[]V, V comparable](slice S, v V) (S, S) {
	dst1 := make(S, 0, len(slice)/2)
	dst2 := make(S, 0, len(slice)/2)
	for i := range slice {
		if slice[i] == v {
			dst1 = append(dst1, slice[i])
		} else {
			dst2 = append(dst2, slice[i])
		}
	}
	return dst1, dst2
}

// 条件を満たし続ける先頭の値のスライスを返す。
func TakeWhileBy[S ~[]V, V any](slice S, f func(V) (bool, error)) (S, error) {
	for i := range slice {
		ok, err := f(slice[i])
		if err != nil {
			return nil, err
		}
		if !ok {
			return slice[:i], nil
		}
	}
	return slice, nil
}

// 条件を満たし続ける先頭の値のスライスを返す。実行中にエラーが起きた場合 panic する。
func MustTakeWhileBy[S ~[]V, V any](slice S, f func(V) (bool, error)) S {
	return must.Must1(TakeWhileBy(slice, f))
}

// 一致し続ける先頭の値のスライスを返す。
func TakeWhile[S ~[]V, V comparable](slice S, v V) S {
	for i := range slice {
		if slice[i] != v {
			return slice[:i]
		}
	}
	return slice
}

// 先頭n個の値のスライスを返す。
func Take[S ~[]V, V any](slice S, n int) S {
	if n > len(slice) {
		return slice
	}
	return slice[:n]
}

// 条件を満たし続ける先頭の値を除いたスライスを返す。
func DropWhileBy[S ~[]V, V any](slice S, f func(V) (bool, error)) (S, error) {
	for i := range slice {
		ok, err := f(slice[i])
		if err != nil {
			return nil, err
		}
		if !ok {
			return slice[i:], nil
		}
	}
	return nil, nil
}

// 条件を満たし続ける先頭の値を除いたスライスを返す。実行中にエラーが起きた場合 panic する。
func MustDropWhileBy[S ~[]V, V any](slice S, f func(V) (bool, error)) S {
	return must.Must1(DropWhileBy(slice, f))
}

// 一致し続ける先頭の値を除いたスライスを返す。
func DropWhile[S ~[]V, V comparable](slice S, v V) S {
	for i := range slice {
		if slice[i] != v {
			return slice[i:]
		}
	}
	return nil
}

// 先頭n個の値を除いたスライスを返す。
func Drop[S ~[]V, V any](slice S, n int) S {
	if n > len(slice) {
		return nil
	}
	return slice[n:]
}

// 条件を満たし続ける先頭部分と残りの部分、ふたつのスライスを返す。
func SpanBy[S ~[]V, V any](slice S, f func(V) (bool, error)) (S, S, error) {
	for i := range slice {
		ok, err := f(slice[i])
		if err != nil {
			return nil, nil, err
		}
		if !ok {
			return slice[:i], slice[i:], nil
		}
	}
	return slice, nil, nil
}

// 条件を満たし続ける先頭部分と残りの部分、ふたつのスライスを返す。実行中にエラーが起きた場合 panic する。
func MustSpanBy[S ~[]V, V any](slice S, f func(V) (bool, error)) (S, S) {
	return must.Must2(SpanBy(slice, f))
}

// 一致し続ける先頭部分と残りの部分、ふたつのスライスを返す。
func Span[S ~[]V, V comparable](slice S, v V) (S, S) {
	for i := range slice {
		if slice[i] != v {
			return slice[:i], slice[i:]
		}
	}
	return slice, nil
}

// ゼロ値を除いたスライスを返す。
func Clean[S ~[]V, V comparable](slice S) S {
	zero := *new(V)
	dst := make(S, 0, len(slice)/2)
	for i := range slice {
		if slice[i] == zero {
			continue
		}
		dst = append(dst, slice[i])
	}
	return dst
}

// 重複を除いたスライスを返す。
func Distinct[S ~[]V, V comparable](slice S) S {
	dst := make(S, 0, len(slice)/2)
	for i := range slice {
		skip := false
		for j := range dst {
			if slice[i] == dst[j] {
				skip = true
			}
		}
		if skip {
			continue
		}
		dst = append(dst, slice[i])
	}
	return dst
}

// 条件を満たす値を変換したスライスを返す。
func Collect[V1 any, V2 any](slice []V1, f func(V1) (V2, bool, error)) ([]V2, error) {
	dst := make([]V2, 0, len(slice)/2)
	for i := range slice {
		v2, ok, err := f(slice[i])
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		dst = append(dst, v2)
	}
	return dst, nil
}

// 条件を満たす値を変換したスライスを返す。実行中にエラーが起きた場合 panic する。
func MustCollect[V1 any, V2 any](slice []V1, f func(V1) (V2, bool, error)) []V2 {
	return must.Must1(Collect(slice, f))
}

// 値と位置をペアにしたスライスを返す。
func ZipWithIndex[V any](slice []V) []tuple.T2[V, int] {
	dst := make([]tuple.T2[V, int], len(slice))
	for i := range slice {
		dst = append(dst, tuple.NewT2(slice[i], i))
	}
	return dst
}

// ２つのスライスの同じ位置の値をペアにしたスライスを返す。
func Zip2[V1 any, V2 any](slice1 []V1, slice2 []V2) []tuple.T2[V1, V2] {
	n := len(slice1)
	if n > len(slice2) {
		n = len(slice2)
	}
	dst := make([]tuple.T2[V1, V2], n)
	for i := 0; i < n; i++ {
		dst[i] = tuple.NewT2(slice1[i], slice2[i])
	}
	return dst
}

// ３つのスライスの同じ位置の値をペアにしたスライスを返す。
func Zip3[V1 any, V2 any, V3 any](slice1 []V1, slice2 []V2, slice3 []V3) []tuple.T3[V1, V2, V3] {
	n := len(slice1)
	if n > len(slice2) {
		n = len(slice2)
	}
	if n > len(slice3) {
		n = len(slice3)
	}
	dst := make([]tuple.T3[V1, V2, V3], n)
	for i := 0; i < n; i++ {
		dst[i] = tuple.NewT3(slice1[i], slice2[i], slice3[i])
	}
	return dst
}

// ４つのスライスの同じ位置の値をペアにしたスライスを返す。
func Zip4[V1 any, V2 any, V3 any, V4 any](slice1 []V1, slice2 []V2, slice3 []V3, slice4 []V4) []tuple.T4[V1, V2, V3, V4] {
	n := len(slice1)
	if n > len(slice2) {
		n = len(slice2)
	}
	if n > len(slice3) {
		n = len(slice3)
	}
	if n > len(slice4) {
		n = len(slice4)
	}
	dst := make([]tuple.T4[V1, V2, V3, V4], n)
	for i := 0; i < n; i++ {
		dst[i] = tuple.NewT4(slice1[i], slice2[i], slice3[i], slice4[i])
	}
	return dst
}

// ５つのスライスの同じ位置の値をペアにしたスライスを返す。
func Zip5[V1 any, V2 any, V3 any, V4 any, V5 any](slice1 []V1, slice2 []V2, slice3 []V3, slice4 []V4, slice5 []V5) []tuple.T5[V1, V2, V3, V4, V5] {
	n := len(slice1)
	if n > len(slice2) {
		n = len(slice2)
	}
	if n > len(slice3) {
		n = len(slice3)
	}
	if n > len(slice4) {
		n = len(slice4)
	}
	if n > len(slice5) {
		n = len(slice5)
	}
	dst := make([]tuple.T5[V1, V2, V3, V4, V5], n)
	for i := 0; i < n; i++ {
		dst[i] = tuple.NewT5(slice1[i], slice2[i], slice3[i], slice4[i], slice5[i])
	}
	return dst
}

// ６つのスライスの同じ位置の値をペアにしたスライスを返す。
func Zip6[V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](slice1 []V1, slice2 []V2, slice3 []V3, slice4 []V4, slice5 []V5, slice6 []V6) []tuple.T6[V1, V2, V3, V4, V5, V6] {
	n := len(slice1)
	if n > len(slice2) {
		n = len(slice2)
	}
	if n > len(slice3) {
		n = len(slice3)
	}
	if n > len(slice4) {
		n = len(slice4)
	}
	if n > len(slice5) {
		n = len(slice5)
	}
	if n > len(slice6) {
		n = len(slice6)
	}
	dst := make([]tuple.T6[V1, V2, V3, V4, V5, V6], n)
	for i := 0; i < n; i++ {
		dst[i] = tuple.NewT6(slice1[i], slice2[i], slice3[i], slice4[i], slice5[i], slice6[i])
	}
	return dst
}

// 値のペアを分離して２つのスライスを返す。
func Unzip2[V1 any, V2 any](slice []tuple.T2[V1, V2]) ([]V1, []V2) {
	dst1 := make([]V1, 0, len(slice))
	dst2 := make([]V2, 0, len(slice))
	for i := range slice {
		dst1 = append(dst1, slice[i].V1)
		dst2 = append(dst2, slice[i].V2)
	}
	return dst1, dst2
}

// 値のペアを分離して３つのスライスを返す。
func Unzip3[V1 any, V2 any, V3 any](slice []tuple.T3[V1, V2, V3]) ([]V1, []V2, []V3) {
	dst1 := make([]V1, 0, len(slice))
	dst2 := make([]V2, 0, len(slice))
	dst3 := make([]V3, 0, len(slice))
	for i := range slice {
		dst1 = append(dst1, slice[i].V1)
		dst2 = append(dst2, slice[i].V2)
		dst3 = append(dst3, slice[i].V3)
	}
	return dst1, dst2, dst3
}

// 値のペアを分離して４つのスライスを返す。
func Unzip4[V1 any, V2 any, V3 any, V4 any](slice []tuple.T4[V1, V2, V3, V4]) ([]V1, []V2, []V3, []V4) {
	dst1 := make([]V1, 0, len(slice))
	dst2 := make([]V2, 0, len(slice))
	dst3 := make([]V3, 0, len(slice))
	dst4 := make([]V4, 0, len(slice))
	for i := range slice {
		dst1 = append(dst1, slice[i].V1)
		dst2 = append(dst2, slice[i].V2)
		dst3 = append(dst3, slice[i].V3)
		dst4 = append(dst4, slice[i].V4)
	}
	return dst1, dst2, dst3, dst4
}

// 値のペアを分離して５つのスライスを返す。
func Unzip5[V1 any, V2 any, V3 any, V4 any, V5 any](slice []tuple.T5[V1, V2, V3, V4, V5]) ([]V1, []V2, []V3, []V4, []V5) {
	dst1 := make([]V1, 0, len(slice))
	dst2 := make([]V2, 0, len(slice))
	dst3 := make([]V3, 0, len(slice))
	dst4 := make([]V4, 0, len(slice))
	dst5 := make([]V5, 0, len(slice))
	for i := range slice {
		dst1 = append(dst1, slice[i].V1)
		dst2 = append(dst2, slice[i].V2)
		dst3 = append(dst3, slice[i].V3)
		dst4 = append(dst4, slice[i].V4)
		dst5 = append(dst5, slice[i].V5)
	}
	return dst1, dst2, dst3, dst4, dst5
}

// 値のペアを分離して６つのスライスを返す。
func Unzip6[V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](slice []tuple.T6[V1, V2, V3, V4, V5, V6]) ([]V1, []V2, []V3, []V4, []V5, []V6) {
	dst1 := make([]V1, 0, len(slice))
	dst2 := make([]V2, 0, len(slice))
	dst3 := make([]V3, 0, len(slice))
	dst4 := make([]V4, 0, len(slice))
	dst5 := make([]V5, 0, len(slice))
	dst6 := make([]V6, 0, len(slice))
	for i := range slice {
		dst1 = append(dst1, slice[i].V1)
		dst2 = append(dst2, slice[i].V2)
		dst3 = append(dst3, slice[i].V3)
		dst4 = append(dst4, slice[i].V4)
		dst5 = append(dst5, slice[i].V5)
		dst6 = append(dst6, slice[i].V6)
	}
	return dst1, dst2, dst3, dst4, dst5, dst6
}

// 値ごとに関数の返すキーでグルーピングしたマップを返す。
func GroupBy[S ~[]V, K comparable, V any](slice S, f func(V) (K, error)) (map[K]S, error) {
	dst := map[K]S{}
	for i := range slice {
		k, err := f(slice[i])
		if err != nil {
			return nil, err
		}
		dst[k] = append(dst[k], slice[i])
	}
	return dst, nil
}

// 値ごとに関数の返すキーでグルーピングしたマップを返す。実行中にエラーが起きた場合 panic する。
func MustGroupBy[S ~[]V, K comparable, V any](slice S, f func(V) (K, error)) map[K]S {
	return must.Must1(GroupBy(slice, f))
}

// 平坦化したスライスを返す。
func Flatten[S ~[]V, V any](slice []S) S {
	dst := make(S, 0, len(slice))
	for i := range slice {
		dst = append(dst, slice[i]...)
	}
	return dst
}

// 値をスライスに変換し、それらを結合したスライスを返す。
func FlatMap[V1, V2 any](slice []V1, f func(V1) ([]V2, error)) ([]V2, error) {
	dst := make([]V2, 0, len(slice))
	for i := range slice {
		slice, err := f(slice[i])
		if err != nil {
			return nil, err
		}
		dst = append(dst, slice...)
	}
	return dst, nil
}

// 値をスライスに変換し、それらを結合したスライスを返す。
func MustFlatMap[V1, V2 any](slice []V1, f func(V1) ([]V2, error)) []V2 {
	return must.Must1(FlatMap(slice, f))
}

// 値のあいだにseparatorを挿入したスライスを返す。
func Join[S ~[]V, V any](slice S, separator V) S {
	dst := make(S, 0, len(slice)*2)
	for i := range slice {
		dst = append(dst, separator, slice[i])
	}
	return dst[1:]
}

// n個ごとのスライスを返す。
func Grouped[S ~[]V, V any](slice S, n int) []S {
	if n == 0 {
		return nil
	}
	dst := make([]S, 0, len(slice)/n+1)
	for i := range slice {
		if i%n == 0 {
			dst = append(dst, make(S, 0, n))
		}
		dst[len(dst)-1] = append(dst[len(dst)-1], slice[i])
	}
	return dst
}

// stepずつズラしたn個ごとのスライスを返す。
func Sliding[S ~[]V, V any](slice S, n int, step int) []S {
	dst := make([]S, 0, len(slice)/step)
	l, r := 0, n
	for ; r < len(slice); l, r = l+step, r+step {
		dst = append(dst, slice[l:r])
	}
	if l < len(slice) {
		dst = append(dst, slice[l:])
	}
	return dst
}

// 要素がn個になるまで先頭に関数の実行結果を挿入する。
func PadBy[S ~[]V, V any](slice S, n int, f func(int) (V, error)) (S, error) {
	if len(slice) >= n {
		return slice, nil
	}
	left := make(S, n-len(slice))
	for i := 0; i < len(left); i++ {
		v, err := f(i)
		if err != nil {
			return nil, err
		}
		left[i] = v
	}
	return append(left, slice...), nil
}

// 要素がn個になるまで先頭に関数の実行結果を挿入する。
func MustPadBy[S ~[]V, V any](slice S, n int, f func(int) (V, error)) S {
	return must.Must1(PadBy(slice, n, f))
}

// 要素がn個になるまで先頭にvを挿入する。
func Pad[S ~[]V, V any](slice S, n int, v V) S {
	if len(slice) >= n {
		return slice
	}
	left := make(S, n-len(slice))
	for i := 0; i < len(left); i++ {
		left[i] = v
	}
	return append(left, slice...)
}

// 要素がn個になるまで先頭にゼロ値を挿入する。
func PadZero[S ~[]V, V any](slice S, n int) S {
	return Pad(slice, n, *new(V))
}

// 要素がn個になるまで末尾に関数の実行結果を挿入する。
func PadRightBy[S ~[]V, V any](slice S, n int, f func(int) (V, error)) (S, error) {
	if len(slice) >= n {
		return slice, nil
	}
	right := make(S, n-len(slice))
	for i := 0; i < len(right); i++ {
		v, err := f(len(slice) + i)
		if err != nil {
			return nil, err
		}
		right[i] = v
	}
	return append(slice, right...), nil
}

// 要素がn個になるまで末尾に関数の実行結果を挿入する。
func MustPadRightBy[S ~[]V, V any](slice S, n int, f func(int) (V, error)) S {
	return must.Must1(PadRightBy(slice, n, f))
}

// 要素がn個になるまで末尾にvを挿入する。
func PadRight[S ~[]V, V any](slice S, n int, v V) S {
	if len(slice) >= n {
		return slice
	}
	right := make(S, n-len(slice))
	for i := 0; i < len(right); i++ {
		right[i] = v
	}
	return append(slice, right...)
}

// 要素がn個になるまで末尾にゼロ値を挿入する。
func PadZeroRight[S ~[]V, V any](slice S, n int) S {
	return PadRight(slice, n, *new(V))
}
