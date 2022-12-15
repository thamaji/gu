package slices

import (
	"github.com/thamaji/gu/iter"
	"github.com/thamaji/gu/tuple"
	"golang.org/x/exp/constraints"
)

// 値の数を返す。
func Len[V any](slice []V) int {
	return len(slice)
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
func Concat[V any](slice1 []V, slice2 []V) []V {
	return append(slice1, slice2...)
}

// スライスの末尾に値を追加する。
func Append[V any](slice []V, v ...V) []V {
	return append(slice, v...)
}

// 指定した位置に値を追加する。
func Insert[V any](slice []V, index int, v ...V) []V {
	return append(slice[:index], append(v, slice[index:]...)...)
}

// 指定した位置の要素を削除する。
func Remove[V any](slice []V, index int) []V {
	dst := make([]V, 0, len(slice)-1)
	for i := range slice {
		if i == index {
			continue
		}
		dst = append(dst, slice[i])
	}
	return dst
}

// 値ごとに関数を実行する。
func ForEach[V any](slice []V, f func(V)) {
	iter.ForEach(iter.FromSlice(slice), f)
}

// 他のスライスと関数で比較し、一致していたらtrueを返す。
func EqualBy[V any](slice1 []V, slice2 []V, f func(V, V) bool) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	return iter.EqualBy(iter.FromSlice(slice1), iter.FromSlice(slice2), f)
}

// 他のスライスと一致していたらtrueを返す。
func Equal[V comparable](slice1 []V, slice2 []V) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	return iter.Equal(iter.FromSlice(slice1), iter.FromSlice(slice2))
}

// 条件を満たす値の数を返す。
func CountBy[V any](slice []V, f func(V) bool) int {
	return iter.CountBy(iter.FromSlice(slice), f)
}

// 一致する値の数を返す。
func Count[V comparable](slice []V, v V) int {
	return iter.Count(iter.FromSlice(slice), v)
}

// 位置のスライスを返す。
func Indices[V any](slice []V) []int {
	return iter.ToSlice(iter.Indices(iter.FromSlice(slice)))
}

// 値を変換したスライスを返す。
func Map[V1 any, V2 any](slice []V1, f func(V1) V2) []V2 {
	return iter.ToSlice(iter.Map(iter.FromSlice(slice), f))
}

// 値を順に演算する。
func Reduce[V any](slice []V, f func(V, V) V) V {
	return iter.Reduce(iter.FromSlice(slice), f)
}

// 値の合計を演算する。
func Sum[V constraints.Ordered | constraints.Complex](slice []V) V {
	return iter.Sum(iter.FromSlice(slice))
}

// 値を変換して合計を演算する。
func SumBy[V1 any, V2 constraints.Ordered | constraints.Complex](slice []V1, f func(V1) V2) V2 {
	return iter.SumBy(iter.FromSlice(slice), f)
}

// 最大の値を返す。
func Max[V constraints.Ordered](slice []V) V {
	return iter.Max(iter.FromSlice(slice))
}

// 値を変換して最大の値を返す
func MaxBy[V1 any, V2 constraints.Ordered](slice []V1, f func(V1) V2) V2 {
	return iter.MaxBy(iter.FromSlice(slice), f)
}

// 最小の値を返す。
func Min[V constraints.Ordered](slice []V) V {
	return iter.Min(iter.FromSlice(slice))
}

// 値を変換して最小の値を返す
func MinBy[V1 any, V2 constraints.Ordered](slice []V1, f func(V1) V2) V2 {
	return iter.MinBy(iter.FromSlice(slice), f)
}

// 初期値と値を順に演算する。
func Fold[V1 any, V2 any](slice []V1, v V2, f func(V2, V1) V2) V2 {
	return iter.Fold(iter.FromSlice(slice), v, f)
}

// 条件を満たす最初の値の位置を返す。
func IndexBy[V any](slice []V, f func(V) bool) int {
	return iter.IndexBy(iter.FromSlice(slice), f)
}

// 一致する最初の値の位置を返す。
func Index[V comparable](slice []V, v V) int {
	return iter.Index(iter.FromSlice(slice), v)
}

// 条件を満たす最後の値の位置を返す。
func LastIndexBy[V any](slice []V, f func(V) bool) int {
	return iter.LastIndexBy(iter.FromSlice(slice), f)
}

// 一致する最後の値の位置を返す。
func LastIndex[V comparable](slice []V, v V) int {
	return iter.LastIndex(iter.FromSlice(slice), v)
}

// 条件を満たす値を返す。
func FindBy[V any](slice []V, f func(V) bool) (V, bool) {
	return iter.FindBy(iter.FromSlice(slice), f)
}

// 一致する値を返す。
func Find[V comparable](slice []V, v V) (V, bool) {
	return iter.Find(iter.FromSlice(slice), v)
}

// 条件を満たす値が存在したらtrueを返す。
func ExistsBy[V any](slice []V, f func(V) bool) bool {
	return iter.ExistsBy(iter.FromSlice(slice), f)
}

// 一致する値が存在したらtrueを返す。
func Exists[V comparable](slice []V, v V) bool {
	return iter.Exists(iter.FromSlice(slice), v)
}

// すべての値が条件を満たせばtrueを返す。
func ForAllBy[V any](slice []V, f func(V) bool) bool {
	return iter.ForAllBy(iter.FromSlice(slice), f)
}

// すべての値が一致したらtrueを返す。
func ForAll[V comparable](slice []V, v V) bool {
	return iter.ForAll(iter.FromSlice(slice), v)
}

// 他のスライスの値がひとつでも存在していたらtrueを返す。
func ContainsAny[V comparable](slice []V, subset []V) bool {
	return iter.ContainsAny(iter.FromSlice(slice), iter.FromSlice(subset))
}

// 他のスライスの値がすべて存在していたらtrueを返す。
func ContainsAll[V comparable](slice []V, subset []V) bool {
	return iter.ContainsAll(iter.FromSlice(slice), iter.FromSlice(subset))
}

// 先頭が他のスライスと一致していたらtrueを返す。
func StartsWith[V comparable](slice []V, subset []V) bool {
	return iter.StartsWith(iter.FromSlice(slice), iter.FromSlice(subset))
}

// 終端が他のスライスと一致していたらtrueを返す。
func EndsWith[V comparable](slice []V, subset []V) bool {
	return iter.EndsWith(iter.FromSlice(slice), iter.FromSlice(subset))
}

// ひとつめのoldをnewで置き換えたスライスを返す。
func Replace[V comparable](slice []V, old V, new V) []V {
	return iter.ToSlice(iter.Replace(iter.FromSlice(slice), old, new))
}

// すべてのoldをnewで置き換えたスライスを返す。
func ReplaceAll[V comparable](slice []V, old V, new V) []V {
	return iter.ToSlice(iter.ReplaceAll(iter.FromSlice(slice), old, new))
}

// 条件を満たす値だけのスライスを返す。
func FilterBy[V any](slice []V, f func(V) bool) []V {
	return iter.ToSlice(iter.FilterBy(iter.FromSlice(slice), f))
}

// 一致する値だけのスライスを返す。
func Filter[V comparable](slice []V, v V) []V {
	return iter.ToSlice(iter.Filter(iter.FromSlice(slice), v))
}

// 条件を満たす値を除いたスライスを返す。
func FilterNotBy[V any](slice []V, f func(V) bool) []V {
	return iter.ToSlice(iter.FilterNotBy(iter.FromSlice(slice), f))
}

// 一致する値を除いたスライスを返す。
func FilterNot[V comparable](slice []V, v V) []V {
	return iter.ToSlice(iter.FilterNot(iter.FromSlice(slice), v))
}

// 条件を満たす値の直前で分割したふたつのスライスを返す。
func SplitBy[V any](slice []V, f func(V) bool) ([]V, []V) {
	iter1, iter2 := iter.SplitBy(iter.FromSlice(slice), f)
	return iter.ToSlice(iter1), iter.ToSlice(iter2)
}

// 一致する値の直前で分割したふたつのスライスを返す。
func Split[V comparable](slice []V, v V) ([]V, []V) {
	iter1, iter2 := iter.Split(iter.FromSlice(slice), v)
	return iter.ToSlice(iter1), iter.ToSlice(iter2)
}

// 条件を満たす値の直後で分割したふたつのスライスを返す。
func SplitAfterBy[V any](slice []V, f func(V) bool) ([]V, []V) {
	iter1, iter2 := iter.SplitAfterBy(iter.FromSlice(slice), f)
	return iter.ToSlice(iter1), iter.ToSlice(iter2)
}

// 一致する値の直後で分割したふたつのスライスを返す。
func SplitAfter[V comparable](slice []V, v V) ([]V, []V) {
	iter1, iter2 := iter.SplitAfter(iter.FromSlice(slice), v)
	return iter.ToSlice(iter1), iter.ToSlice(iter2)
}

// 条件を満たすスライスと満たさないスライスを返す。
func PartitionBy[V any](slice []V, f func(V) bool) ([]V, []V) {
	iter1, iter2 := iter.PartitionBy(iter.FromSlice(slice), f)
	return iter.ToSlice(iter1), iter.ToSlice(iter2)
}

// 値の一致するスライスと一致しないスライスを返す。
func Partition[V comparable](slice []V, v V) ([]V, []V) {
	iter1, iter2 := iter.Partition(iter.FromSlice(slice), v)
	return iter.ToSlice(iter1), iter.ToSlice(iter2)
}

// 条件を満たし続ける先頭の値のスライスを返す。
func TakeWhileBy[V any](slice []V, f func(V) bool) []V {
	return iter.ToSlice(iter.TakeWhileBy(iter.FromSlice(slice), f))
}

// 一致し続ける先頭の値のスライスを返す。
func TakeWhile[V comparable](slice []V, v V) []V {
	return iter.ToSlice(iter.TakeWhile(iter.FromSlice(slice), v))
}

// 先頭n個の値のスライスを返す。
func Take[V any](slice []V, n int) []V {
	return iter.ToSlice(iter.Take(iter.FromSlice(slice), n))
}

// 条件を満たし続ける先頭の値を除いたスライスを返す。
func DropWhileBy[V any](slice []V, f func(V) bool) []V {
	return iter.ToSlice(iter.DropWhileBy(iter.FromSlice(slice), f))
}

// 一致し続ける先頭の値を除いたスライスを返す。
func DropWhile[V comparable](slice []V, v V) []V {
	return iter.ToSlice(iter.DropWhile(iter.FromSlice(slice), v))
}

// 先頭n個の値を除いたスライスを返す。
func Drop[V any](slice []V, n int) []V {
	return iter.ToSlice(iter.Drop(iter.FromSlice(slice), n))
}

// 条件を満たし続ける先頭部分と残りの部分、ふたつのスライスを返す。
func SpanBy[V any](slice []V, f func(V) bool) ([]V, []V) {
	iter1, iter2 := iter.SpanBy(iter.FromSlice(slice), f)
	return iter.ToSlice(iter1), iter.ToSlice(iter2)
}

// 一致し続ける先頭部分と残りの部分、ふたつのスライスを返す。
func Span[V comparable](slice []V, v V) ([]V, []V) {
	iter1, iter2 := iter.Span(iter.FromSlice(slice), v)
	return iter.ToSlice(iter1), iter.ToSlice(iter2)
}

// ゼロ値を除いたスライスを返す。
func Clean[V comparable](slice []V) []V {
	return iter.ToSlice(iter.Clean(iter.FromSlice(slice)))
}

// 重複を除いたスライスを返す。
func Distinct[V comparable](slice []V) []V {
	return iter.ToSlice(iter.Distinct(iter.FromSlice(slice)))
}

// 条件を満たす値を変換したスライスを返す。
func Collect[V1 any, V2 any](slice []V1, f func(V1) (V2, bool)) []V2 {
	return iter.ToSlice(iter.Collect(iter.FromSlice(slice), f))
}

// 値と位置をペアにしたスライスを返す。
func ZipWithIndex[T any](slice []T) []tuple.T2[T, int] {
	return iter.ToSlice(iter.ZipWithIndex(iter.FromSlice(slice)))
}

// ２つのスライスの同じ位置の値をペアにしたスライスを返す。
func Zip2[V1 any, V2 any](slice1 []V1, slice2 []V2) []tuple.T2[V1, V2] {
	return iter.ToSlice(iter.Zip2(iter.FromSlice(slice1), iter.FromSlice(slice2)))
}

// ３つのスライスの同じ位置の値をペアにしたスライスを返す。
func Zip3[V1 any, V2 any, V3 any](slice1 []V1, slice2 []V2, slice3 []V3) []tuple.T3[V1, V2, V3] {
	return iter.ToSlice(iter.Zip3(iter.FromSlice(slice1), iter.FromSlice(slice2), iter.FromSlice(slice3)))
}

// ４つのスライスの同じ位置の値をペアにしたスライスを返す。
func Zip4[V1 any, V2 any, V3 any, V4 any](slice1 []V1, slice2 []V2, slice3 []V3, slice4 []V4) []tuple.T4[V1, V2, V3, V4] {
	return iter.ToSlice(iter.Zip4(iter.FromSlice(slice1), iter.FromSlice(slice2), iter.FromSlice(slice3), iter.FromSlice(slice4)))
}

// ５つのスライスの同じ位置の値をペアにしたスライスを返す。
func Zip5[V1 any, V2 any, V3 any, V4 any, V5 any](slice1 []V1, slice2 []V2, slice3 []V3, slice4 []V4, slice5 []V5) []tuple.T5[V1, V2, V3, V4, V5] {
	return iter.ToSlice(iter.Zip5(iter.FromSlice(slice1), iter.FromSlice(slice2), iter.FromSlice(slice3), iter.FromSlice(slice4), iter.FromSlice(slice5)))
}

// ６つのスライスの同じ位置の値をペアにしたスライスを返す。
func Zip6[V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](slice1 []V1, slice2 []V2, slice3 []V3, slice4 []V4, slice5 []V5, slice6 []V6) []tuple.T6[V1, V2, V3, V4, V5, V6] {
	return iter.ToSlice(iter.Zip6(iter.FromSlice(slice1), iter.FromSlice(slice2), iter.FromSlice(slice3), iter.FromSlice(slice4), iter.FromSlice(slice5), iter.FromSlice(slice6)))
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
func GroupBy[K comparable, V any](slice []V, f func(V) K) map[K][]V {
	return iter.GroupBy(iter.FromSlice(slice), f)
}

// 平坦化したスライスを返す。
func Flatten[V any](slice [][]V) []V {
	dst := make([]V, 0, len(slice))
	for i := range slice {
		dst = append(dst, slice[i]...)
	}
	return dst
}

// 値をスライスに変換し、それらを結合したスライスを返す。
func FlatMap[V1, V2 any](slice []V1, f func(V1) []V2) []V2 {
	dst := make([]V2, 0, len(slice))
	for i := range slice {
		dst = append(dst, f(slice[i])...)
	}
	return dst
}

// 値のあいだにseparatorを挿入したスライスを返す。
func Join[V any](slice []V, separator V) []V {
	return Drop(FlatMap(slice, func(v V) []V { return []V{separator, v} }), 1)
}
