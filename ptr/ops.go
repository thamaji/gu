package ptr

import (
	"github.com/thamaji/gu/iter"
	"github.com/thamaji/gu/tuple"
	"golang.org/x/exp/constraints"
)

// 値の数を返す。
func Len[V any](p *V) int {
	return iter.Len(iter.FromPtr(p))
}

// 値を返す。
func Get[V any](p *V) (V, bool) {
	return iter.GetFirst(iter.FromPtr(p))
}

// 値を返す。無い場合はvを返す。
func GetOrElse[V any](p *V, v V) V {
	return iter.GetFirstOrElse(iter.FromPtr(p), v)
}

// 値ごとに関数を実行する。
func ForEach[V any](p *V, f func(V)) {
	iter.ForEach(iter.FromPtr(p), f)
}

// 他のポインタと関数で比較し、一致していたらtrueを返す。
func EqualBy[V any](p1 *V, p2 *V, f func(V, V) bool) bool {
	return iter.EqualBy(iter.FromPtr(p1), iter.FromPtr(p2), f)
}

// 他のポインタと一致していたらtrueを返す。
func Equal[V comparable](p1 *V, p2 *V) bool {
	return iter.Equal(iter.FromPtr(p1), iter.FromPtr(p2))
}

// 条件を満たす値の数を返す。
func CountBy[V any](p *V, f func(V) bool) int {
	return iter.CountBy(iter.FromPtr(p), f)
}

// 一致する値の数を返す。
func Count[V comparable](p *V, v V) int {
	return iter.Count(iter.FromPtr(p), v)
}

// 位置のポインタを返す。
func Indices[V any](p *V) *int {
	return iter.ToPtr(iter.Indices(iter.FromPtr(p)))
}

// 値を変換したポインタを返す。
func Map[V1 any, V2 any](p *V1, f func(V1) V2) *V2 {
	return iter.ToPtr(iter.Map(iter.FromPtr(p), f))
}

// 値を順に演算する。
func Reduce[V any](p *V, f func(V, V) V) V {
	return iter.Reduce(iter.FromPtr(p), f)
}

// 値の合計を演算する。
func Sum[V constraints.Ordered | constraints.Complex](p *V) V {
	return iter.Sum(iter.FromPtr(p))
}

// 値を変換して合計を演算する。
func SumBy[V1 any, V2 constraints.Ordered | constraints.Complex](p *V1, f func(V1) V2) V2 {
	return iter.SumBy(iter.FromPtr(p), f)
}

// 最大の値を返す。
func Max[V constraints.Ordered](p *V) V {
	return iter.Max(iter.FromPtr(p))
}

// 値を変換して最大の値を返す
func MaxBy[V1 any, V2 constraints.Ordered](p *V1, f func(V1) V2) V2 {
	return iter.MaxBy(iter.FromPtr(p), f)
}

// 最小の値を返す。
func Min[V constraints.Ordered](p *V) V {
	return iter.Min(iter.FromPtr(p))
}

// 値を変換して最小の値を返す
func MinBy[V1 any, V2 constraints.Ordered](p *V1, f func(V1) V2) V2 {
	return iter.MinBy(iter.FromPtr(p), f)
}

// 初期値と値を順に演算する。
func Fold[V1 any, V2 any](p *V1, v V2, f func(V2, V1) V2) V2 {
	return iter.Fold(iter.FromPtr(p), v, f)
}

// 条件を満たす最初の値の位置を返す。
func IndexBy[V any](p *V, f func(V) bool) int {
	return iter.IndexBy(iter.FromPtr(p), f)
}

// 一致する最初の値の位置を返す。
func Index[V comparable](p *V, v V) int {
	return iter.Index(iter.FromPtr(p), v)
}

// 条件を満たす最後の値の位置を返す。
func LastIndexBy[V any](p *V, f func(V) bool) int {
	return iter.LastIndexBy(iter.FromPtr(p), f)
}

// 一致する最後の値の位置を返す。
func LastIndex[V comparable](p *V, v V) int {
	return iter.LastIndex(iter.FromPtr(p), v)
}

// 条件を満たす値を返す。
func FindBy[V any](p *V, f func(V) bool) (V, bool) {
	return iter.FindBy(iter.FromPtr(p), f)
}

// 一致する値を返す。
func Find[V comparable](p *V, v V) (V, bool) {
	return iter.Find(iter.FromPtr(p), v)
}

// 条件を満たす値が存在したらtrueを返す。
func ExistsBy[V any](p *V, f func(V) bool) bool {
	return iter.ExistsBy(iter.FromPtr(p), f)
}

// 一致する値が存在したらtrueを返す。
func Exists[V comparable](p *V, v V) bool {
	return iter.Exists(iter.FromPtr(p), v)
}

// すべての値が条件を満たせばtrueを返す。
func ForAllBy[V any](p *V, f func(V) bool) bool {
	return iter.ForAllBy(iter.FromPtr(p), f)
}

// すべての値が一致したらtrueを返す。
func ForAll[V comparable](p *V, v V) bool {
	return iter.ForAll(iter.FromPtr(p), v)
}

// 他のポインタの値がひとつでも存在していたらtrueを返す。
func ContainsAny[V comparable](p *V, subset *V) bool {
	return iter.ContainsAny(iter.FromPtr(p), iter.FromPtr(subset))
}

// 他のポインタの値がすべて存在していたらtrueを返す。
func ContainsAll[V comparable](p *V, subset *V) bool {
	return iter.ContainsAll(iter.FromPtr(p), iter.FromPtr(subset))
}

// 先頭が他のポインタと一致していたらtrueを返す。
func StartsWith[V comparable](p *V, subset *V) bool {
	return iter.StartsWith(iter.FromPtr(p), iter.FromPtr(subset))
}

// 終端が他のポインタと一致していたらtrueを返す。
func EndsWith[V comparable](p *V, subset *V) bool {
	return iter.EndsWith(iter.FromPtr(p), iter.FromPtr(subset))
}

// ひとつめのoldをnewで置き換えたポインタを返す。
func Replace[V comparable](p *V, old V, new V) *V {
	return iter.ToPtr(iter.Replace(iter.FromPtr(p), old, new))
}

// すべてのoldをnewで置き換えたポインタを返す。
func ReplaceAll[V comparable](p *V, old V, new V) *V {
	return iter.ToPtr(iter.ReplaceAll(iter.FromPtr(p), old, new))
}

// 条件を満たす値だけのポインタを返す。
func FilterBy[V any](p *V, f func(V) bool) *V {
	return iter.ToPtr(iter.FilterBy(iter.FromPtr(p), f))
}

// 一致する値だけのポインタを返す。
func Filter[V comparable](p *V, v V) *V {
	return iter.ToPtr(iter.Filter(iter.FromPtr(p), v))
}

// 条件を満たす値を除いたポインタを返す。
func FilterNotBy[V any](p *V, f func(V) bool) *V {
	return iter.ToPtr(iter.FilterNotBy(iter.FromPtr(p), f))
}

// 一致する値を除いたポインタを返す。
func FilterNot[V comparable](p *V, v V) *V {
	return iter.ToPtr(iter.FilterNot(iter.FromPtr(p), v))
}

// 条件を満たす値の直前で分割したふたつのポインタを返す。
func SplitBy[V any](p *V, f func(V) bool) (*V, *V) {
	iter1, iter2 := iter.SplitBy(iter.FromPtr(p), f)
	return iter.ToPtr(iter1), iter.ToPtr(iter2)
}

// 一致する値の直前で分割したふたつのポインタを返す。
func Split[V comparable](p *V, v V) (*V, *V) {
	iter1, iter2 := iter.Split(iter.FromPtr(p), v)
	return iter.ToPtr(iter1), iter.ToPtr(iter2)
}

// 条件を満たす値の直後で分割したふたつのポインタを返す。
func SplitAfterBy[V any](p *V, f func(V) bool) (*V, *V) {
	iter1, iter2 := iter.SplitAfterBy(iter.FromPtr(p), f)
	return iter.ToPtr(iter1), iter.ToPtr(iter2)
}

// 一致する値の直後で分割したふたつのポインタを返す。
func SplitAfter[V comparable](p *V, v V) (*V, *V) {
	iter1, iter2 := iter.SplitAfter(iter.FromPtr(p), v)
	return iter.ToPtr(iter1), iter.ToPtr(iter2)
}

// 条件を満たすポインタと満たさないポインタを返す。
func PartitionBy[V any](p *V, f func(V) bool) (*V, *V) {
	iter1, iter2 := iter.PartitionBy(iter.FromPtr(p), f)
	return iter.ToPtr(iter1), iter.ToPtr(iter2)
}

// 値の一致するポインタと一致しないポインタを返す。
func Partition[V comparable](p *V, v V) (*V, *V) {
	iter1, iter2 := iter.Partition(iter.FromPtr(p), v)
	return iter.ToPtr(iter1), iter.ToPtr(iter2)
}

// 条件を満たし続ける先頭の値のポインタを返す。
func TakeWhileBy[V any](p *V, f func(V) bool) *V {
	return iter.ToPtr(iter.TakeWhileBy(iter.FromPtr(p), f))
}

// 一致し続ける先頭の値のポインタを返す。
func TakeWhile[V comparable](p *V, v V) *V {
	return iter.ToPtr(iter.TakeWhile(iter.FromPtr(p), v))
}

// 先頭n個の値のポインタを返す。
func Take[V any](p *V, n int) *V {
	return iter.ToPtr(iter.Take(iter.FromPtr(p), n))
}

// 条件を満たし続ける先頭の値を除いたポインタを返す。
func DropWhileBy[V any](p *V, f func(V) bool) *V {
	return iter.ToPtr(iter.DropWhileBy(iter.FromPtr(p), f))
}

// 一致し続ける先頭の値を除いたポインタを返す。
func DropWhile[V comparable](p *V, v V) *V {
	return iter.ToPtr(iter.DropWhile(iter.FromPtr(p), v))
}

// 先頭n個の値を除いたポインタを返す。
func Drop[V any](p *V, n int) *V {
	return iter.ToPtr(iter.Drop(iter.FromPtr(p), n))
}

// 条件を満たし続ける先頭部分と残りの部分、ふたつのポインタを返す。
func SpanBy[V any](p *V, f func(V) bool) (*V, *V) {
	iter1, iter2 := iter.SpanBy(iter.FromPtr(p), f)
	return iter.ToPtr(iter1), iter.ToPtr(iter2)
}

// 一致し続ける先頭部分と残りの部分、ふたつのポインタを返す。
func Span[V comparable](p *V, v V) (*V, *V) {
	iter1, iter2 := iter.Span(iter.FromPtr(p), v)
	return iter.ToPtr(iter1), iter.ToPtr(iter2)
}

// ゼロ値を除いたポインタを返す。
func Clean[V comparable](p *V) *V {
	return iter.ToPtr(iter.Clean(iter.FromPtr(p)))
}

// 重複を除いたポインタを返す。
func Distinct[V comparable](p *V) *V {
	return p
}

// 条件を満たす値を変換したポインタを返す。
func Collect[V1 any, V2 any](p *V1, f func(V1) (V2, bool)) *V2 {
	return iter.ToPtr(iter.Collect(iter.FromPtr(p), f))
}

// 値と位置をペアにしたポインタを返す。
func ZipWithIndex[V any](p *V) *tuple.T2[V, int] {
	return iter.ToPtr(iter.ZipWithIndex(iter.FromPtr(p)))
}

// ２つのポインタの同じ位置の値をペアにしたポインタを返す。
func Zip2[V1 any, V2 any](p1 *V1, p2 *V2) *tuple.T2[V1, V2] {
	return iter.ToPtr(iter.Zip2(iter.FromPtr(p1), iter.FromPtr(p2)))
}

// ３つのポインタの同じ位置の値をペアにしたポインタを返す。
func Zip3[V1 any, V2 any, V3 any](p1 *V1, p2 *V2, p3 *V3) *tuple.T3[V1, V2, V3] {
	return iter.ToPtr(iter.Zip3(iter.FromPtr(p1), iter.FromPtr(p2), iter.FromPtr(p3)))
}

// ４つのポインタの同じ位置の値をペアにしたポインタを返す。
func Zip4[V1 any, V2 any, V3 any, V4 any](p1 *V1, p2 *V2, p3 *V3, p4 *V4) *tuple.T4[V1, V2, V3, V4] {
	return iter.ToPtr(iter.Zip4(iter.FromPtr(p1), iter.FromPtr(p2), iter.FromPtr(p3), iter.FromPtr(p4)))
}

// ５つのポインタの同じ位置の値をペアにしたポインタを返す。
func Zip5[V1 any, V2 any, V3 any, V4 any, V5 any](p1 *V1, p2 *V2, p3 *V3, p4 *V4, p5 *V5) *tuple.T5[V1, V2, V3, V4, V5] {
	return iter.ToPtr(iter.Zip5(iter.FromPtr(p1), iter.FromPtr(p2), iter.FromPtr(p3), iter.FromPtr(p4), iter.FromPtr(p5)))
}

// ６つのポインタの同じ位置の値をペアにしたポインタを返す。
func Zip6[V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](p1 *V1, p2 *V2, p3 *V3, p4 *V4, p5 *V5, p6 *V6) *tuple.T6[V1, V2, V3, V4, V5, V6] {
	return iter.ToPtr(iter.Zip6(iter.FromPtr(p1), iter.FromPtr(p2), iter.FromPtr(p3), iter.FromPtr(p4), iter.FromPtr(p5), iter.FromPtr(p6)))
}

// 値のペアを分離して２つのポインタを返す。
func Unzip2[V1 any, V2 any](p *tuple.T2[V1, V2]) (*V1, *V2) {
	iter1, iter2 := iter.Unzip2(iter.FromPtr(p))
	return iter.ToPtr(iter1), iter.ToPtr(iter2)
}

// 値のペアを分離して３つのポインタを返す。
func Unzip3[V1 any, V2 any, V3 any](p *tuple.T3[V1, V2, V3]) (*V1, *V2, *V3) {
	iter1, iter2, iter3 := iter.Unzip3(iter.FromPtr(p))
	return iter.ToPtr(iter1), iter.ToPtr(iter2), iter.ToPtr(iter3)
}

// 値のペアを分離して４つのポインタを返す。
func Unzip4[V1 any, V2 any, V3 any, V4 any](p *tuple.T4[V1, V2, V3, V4]) (*V1, *V2, *V3, *V4) {
	iter1, iter2, iter3, iter4 := iter.Unzip4(iter.FromPtr(p))
	return iter.ToPtr(iter1), iter.ToPtr(iter2), iter.ToPtr(iter3), iter.ToPtr(iter4)
}

// 値のペアを分離して５つのポインタを返す。
func Unzip5[V1 any, V2 any, V3 any, V4 any, V5 any](p *tuple.T5[V1, V2, V3, V4, V5]) (*V1, *V2, *V3, *V4, *V5) {
	iter1, iter2, iter3, iter4, iter5 := iter.Unzip5(iter.FromPtr(p))
	return iter.ToPtr(iter1), iter.ToPtr(iter2), iter.ToPtr(iter3), iter.ToPtr(iter4), iter.ToPtr(iter5)
}

// 値のペアを分離して６つのポインタを返す。
func Unzip6[V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](p *tuple.T6[V1, V2, V3, V4, V5, V6]) (*V1, *V2, *V3, *V4, *V5, *V6) {
	iter1, iter2, iter3, iter4, iter5, iter6 := iter.Unzip6(iter.FromPtr(p))
	return iter.ToPtr(iter1), iter.ToPtr(iter2), iter.ToPtr(iter3), iter.ToPtr(iter4), iter.ToPtr(iter5), iter.ToPtr(iter6)
}

// 値ごとに関数の返すキーでグルーピングしたマップを返す。
func GroupBy[K comparable, V any](p *V, f func(V) K) map[K][]V {
	return iter.GroupBy(iter.FromPtr(p), f)
}

// 平坦化したポインタを返す。
func Flatten[V any](p **V) *V {
	if p == nil {
		return nil
	}
	return *p
}

// 値をイテレータに変換し、それらを結合したイテレータを返す。
func FlatMap[V1 any, V2 any](p *V1, f func(V1) *V2) *V2 {
	return Flatten(Map(p, f))
}

// 値のあいだにseparatorを挿入したイテレータを返す。
func Join[V any](p *V, separator V) *V {
	return p
}