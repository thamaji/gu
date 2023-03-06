package ptr

import (
	"github.com/thamaji/gu/must"
	"github.com/thamaji/gu/tuple"
	"golang.org/x/exp/constraints"
)

// 値の数を返す。
func Len[V any](p *V) int {
	if p == nil {
		return 0
	}
	return 1
}

// 空のときtrueを返す。
func IsEmpty[V any](p *V) bool {
	return p == nil
}

// 値を返す。
func Get[V any](p *V) (V, bool) {
	if p == nil {
		return *new(V), false
	}
	return *p, true
}

// 値を返す。無い場合はvを返す。
func GetOrElse[V any](p *V, v V) V {
	if p == nil {
		return v
	}
	return *p
}

// 値を返す。無い場合は関数の実行結果を返す。
func GetOrFunc[V any](p *V, f func() (V, error)) (V, error) {
	if p == nil {
		return f()
	}
	return *p, nil
}

// 値を返す。無い場合は関数の実行結果を返す。実行中にエラーが起きた場合 panic する。
func MustGetOrFunc[V any](p *V, f func() (V, error)) V {
	return must.Must1(GetOrFunc(p, f))
}

// 値ごとに関数を実行する。
func ForEach[V any](p *V, f func(V) error) error {
	if p == nil {
		return nil
	}
	return f(*p)
}

// 値ごとに関数を実行する。実行中にエラーが起きた場合 panic する。
func MustForEach[V any](p *V, f func(V) error) {
	must.Must0(ForEach(p, f))
}

// 他のポインタと関数で比較し、一致していたらtrueを返す。
func EqualBy[V any](p1 *V, p2 *V, f func(V, V) (bool, error)) (bool, error) {
	if p1 == nil && p2 == nil {
		return true, nil
	}
	if p1 == nil || p2 == nil {
		return false, nil
	}
	return f(*p1, *p2)
}

// 他のポインタと関数で比較し、一致していたらtrueを返す。実行中にエラーが起きた場合 panic する。
func MustEqualBy[V any](p1 *V, p2 *V, f func(V, V) (bool, error)) bool {
	return must.Must1(EqualBy(p1, p2, f))
}

// 他のポインタと一致していたらtrueを返す。
func Equal[V comparable](p1 *V, p2 *V) bool {
	if p1 == nil && p2 == nil {
		return true
	}
	if p1 == nil || p2 == nil {
		return false
	}
	return *p1 == *p2
}

// 条件を満たす値の数を返す。
func CountBy[V any](p *V, f func(V) (bool, error)) (int, error) {
	if p == nil {
		return 0, nil
	}
	ok, err := f(*p)
	if err != nil {
		return 0, err
	}
	if ok {
		return 1, nil
	}
	return 0, nil
}

// 条件を満たす値の数を返す。実行中にエラーが起きた場合 panic する。
func MustCountBy[V any](p *V, f func(V) (bool, error)) int {
	return must.Must1(CountBy(p, f))
}

// 一致する値の数を返す。
func Count[V comparable](p *V, v V) int {
	if p == nil {
		return 0
	}
	if *p == v {
		return 1
	}
	return 0
}

// 位置のポインタを返す。
func Indices[V any](p *V) *int {
	if p == nil {
		return nil
	}
	i := 0
	return &i
}

// 値を変換したポインタを返す。
func Map[V1 any, V2 any](p *V1, f func(V1) (V2, error)) (*V2, error) {
	if p == nil {
		return nil, nil
	}
	v2, err := f(*p)
	if err != nil {
		return nil, err
	}
	return &v2, nil
}

// 値を変換したポインタを返す。実行中にエラーが起きた場合 panic する。
func MustMap[V1 any, V2 any](p *V1, f func(V1) (V2, error)) *V2 {
	return must.Must1(Map(p, f))
}

// 値を順に演算する。
func Reduce[V any](p *V, f func(V, V) (V, error)) (V, error) {
	var v V
	if p != nil {
		v = *p
	}
	return v, nil
}

// 値を順に演算する。実行中にエラーが起きた場合 panic する。
func MustReduce[V any](p *V, f func(V, V) (V, error)) V {
	return must.Must1(Reduce(p, f))
}

// 値の合計を演算する。
func Sum[V constraints.Ordered | constraints.Complex](p *V) V {
	var v V
	if p != nil {
		v = *p
	}
	return v
}

// 値を変換して合計を演算する。
func SumBy[V1 any, V2 constraints.Ordered | constraints.Complex](p *V1, f func(V1) (V2, error)) (V2, error) {
	var v V2
	var err error
	if p != nil {
		v, err = f(*p)
	}
	return v, err
}

// 値を変換して合計を演算する。実行中にエラーが起きた場合 panic する。
func MustSumBy[V1 any, V2 constraints.Ordered | constraints.Complex](p *V1, f func(V1) (V2, error)) V2 {
	return must.Must1(SumBy(p, f))
}

// 最大の値を返す。
func Max[V constraints.Ordered](p *V) V {
	var v V
	if p != nil {
		v = *p
	}
	return v
}

// 値を変換して最大の値を返す。
func MaxBy[V1 any, V2 constraints.Ordered](p *V1, f func(V1) (V2, error)) (V2, error) {
	var v V2
	var err error
	if p != nil {
		v, err = f(*p)
	}
	return v, err
}

// 値を変換して最大の値を返す。実行中にエラーが起きた場合 panic する。
func MustMaxBy[V1 any, V2 constraints.Ordered](p *V1, f func(V1) (V2, error)) V2 {
	return must.Must1(MaxBy(p, f))
}

// 最小の値を返す。
func Min[V constraints.Ordered](p *V) V {
	var v V
	if p != nil {
		v = *p
	}
	return v
}

// 値を変換して最小の値を返す。
func MinBy[V1 any, V2 constraints.Ordered](p *V1, f func(V1) (V2, error)) (V2, error) {
	var v V2
	var err error
	if p != nil {
		v, err = f(*p)
	}
	return v, err
}

// 値を変換して最小の値を返す。実行中にエラーが起きた場合 panic する。
func MustMinBy[V1 any, V2 constraints.Ordered](p *V1, f func(V1) (V2, error)) V2 {
	return must.Must1(MinBy(p, f))
}

// 初期値と値を順に演算する。
func Fold[V1 any, V2 any](p *V1, v V2, f func(V2, V1) (V2, error)) (V2, error) {
	var err error
	if p != nil {
		v, err = f(v, *p)
	}
	return v, err
}

// 初期値と値を順に演算する。実行中にエラーが起きた場合 panic する。
func MustFold[V1 any, V2 any](p *V1, v V2, f func(V2, V1) (V2, error)) V2 {
	return must.Must1(Fold(p, v, f))
}

// 条件を満たす最初の値の位置を返す。
func IndexBy[V any](p *V, f func(V) (bool, error)) (int, error) {
	if p != nil {
		ok, err := f(*p)
		if err != nil {
			return 0, err
		}
		if ok {
			return 0, nil
		}
	}
	return -1, nil
}

// 条件を満たす最初の値の位置を返す。実行中にエラーが起きた場合 panic する。
func MustIndexBy[V any](p *V, f func(V) (bool, error)) int {
	return must.Must1(IndexBy(p, f))
}

// 一致する最初の値の位置を返す。
func Index[V comparable](p *V, v V) int {
	if p != nil {
		if *p == v {
			return 0
		}
	}
	return -1
}

// 条件を満たす最後の値の位置を返す。
func LastIndexBy[V any](p *V, f func(V) (bool, error)) (int, error) {
	if p != nil {
		ok, err := f(*p)
		if err != nil {
			return 0, err
		}
		if ok {
			return 0, nil
		}
	}
	return -1, nil
}

// 条件を満たす最後の値の位置を返す。実行中にエラーが起きた場合 panic する。
func MustLastIndexBy[V any](p *V, f func(V) (bool, error)) int {
	return must.Must1(LastIndexBy(p, f))
}

// 一致する最後の値の位置を返す。
func LastIndex[V comparable](p *V, v V) int {
	if p != nil {
		if *p == v {
			return 0
		}
	}
	return -1
}

// 条件を満たす値を返す。
func FindBy[V any](p *V, f func(V) (bool, error)) (V, bool, error) {
	if p != nil {
		ok, err := f(*p)
		if err != nil {
			return *new(V), false, err
		}
		if ok {
			return *p, true, nil
		}
	}
	return *new(V), false, nil
}

// 条件を満たす値を返す。実行中にエラーが起きた場合 panic する。
func MustFindBy[V any](p *V, f func(V) (bool, error)) (V, bool) {
	return must.Must2(FindBy(p, f))
}

// 一致する値を返す。
func Find[V comparable](p *V, v V) (V, bool) {
	if p != nil {
		if *p == v {
			return *p, true
		}
	}
	return *new(V), false
}

// 条件を満たす値が存在したらtrueを返す。
func ExistsBy[V any](p *V, f func(V) (bool, error)) (bool, error) {
	if p != nil {
		ok, err := f(*p)
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
func MustExistsBy[V any](p *V, f func(V) (bool, error)) bool {
	return must.Must1(ExistsBy(p, f))
}

// 一致する値が存在したらtrueを返す。
func Exists[V comparable](p *V, v V) bool {
	if p != nil {
		if *p == v {
			return true
		}
	}
	return false
}

// すべての値が条件を満たせばtrueを返す。
func ForAllBy[V any](p *V, f func(V) (bool, error)) (bool, error) {
	if p != nil {
		ok, err := f(*p)
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
func MustForAllBy[V any](p *V, f func(V) (bool, error)) bool {
	return must.Must1(ForAllBy(p, f))
}

// すべての値が一致したらtrueを返す。
func ForAll[V comparable](p *V, v V) bool {
	if p != nil {
		if *p != v {
			return false
		}
	}
	return true
}

// 他のポインタの値がひとつでも存在していたらtrueを返す。
func ContainsAny[V comparable](p *V, subset *V) bool {
	if p == nil || subset == nil {
		return false
	}
	return *p == *subset
}

// 他のポインタの値がすべて存在していたらtrueを返す。
func ContainsAll[V comparable](p *V, subset *V) bool {
	if p == nil || subset == nil {
		return false
	}
	return *p == *subset
}

// 先頭が他のポインタと一致していたらtrueを返す。
func StartsWith[V comparable](p *V, subset *V) bool {
	if p == nil || subset == nil {
		return false
	}
	return *p == *subset
}

// 終端が他のポインタと一致していたらtrueを返す。
func EndsWith[V comparable](p *V, subset *V) bool {
	if p == nil || subset == nil {
		return false
	}
	return *p == *subset
}

// ひとつめのoldをnewで置き換えたポインタを返す。
func Replace[V comparable](p *V, old V, new V) *V {
	if p == nil {
		return nil
	}
	if *p == old {
		return &new
	}
	return p
}

// すべてのoldをnewで置き換えたポインタを返す。
func ReplaceAll[V comparable](p *V, old V, new V) *V {
	if p == nil {
		return nil
	}
	if *p == old {
		return &new
	}
	return p
}

// 条件を満たす値だけのポインタを返す。
func FilterBy[V any](p *V, f func(V) (bool, error)) (*V, error) {
	if p == nil {
		return nil, nil
	}
	ok, err := f(*p)
	if err != nil {
		return nil, err
	}
	if ok {
		return p, nil
	}
	return nil, nil
}

// 条件を満たす値だけのポインタを返す。実行中にエラーが起きた場合 panic する。
func MustFilterBy[V any](p *V, f func(V) (bool, error)) *V {
	return must.Must1(FilterBy(p, f))
}

// 一致する値だけのポインタを返す。
func Filter[V comparable](p *V, v V) *V {
	if p == nil {
		return nil
	}
	if *p == v {
		return p
	}
	return nil
}

// 条件を満たす値を除いたポインタを返す。
func FilterNotBy[V any](p *V, f func(V) (bool, error)) (*V, error) {
	if p == nil {
		return nil, nil
	}
	ok, err := f(*p)
	if err != nil {
		return nil, err
	}
	if !ok {
		return p, nil
	}
	return nil, nil
}

// 条件を満たす値を除いたポインタを返す。実行中にエラーが起きた場合 panic する。
func MustFilterNotBy[V any](p *V, f func(V) (bool, error)) *V {
	return must.Must1(FilterNotBy(p, f))
}

// 一致する値を除いたポインタを返す。
func FilterNot[V comparable](p *V, v V) *V {
	if p == nil {
		return nil
	}
	if *p != v {
		return p
	}
	return nil
}

// 条件を満たす値の直前で分割したふたつのポインタを返す。
func SplitBy[V any](p *V, f func(V) (bool, error)) (*V, *V, error) {
	if p == nil {
		return nil, nil, nil
	}
	ok, err := f(*p)
	if err != nil {
		return nil, nil, err
	}
	if ok {
		return nil, p, nil
	}
	return nil, nil, nil
}

// 条件を満たす値の直前で分割したふたつのポインタを返す。実行中にエラーが起きた場合 panic する。
func MustSplitBy[V any](p *V, f func(V) (bool, error)) (*V, *V) {
	return must.Must2(SplitBy(p, f))
}

// 一致する値の直前で分割したふたつのポインタを返す。
func Split[V comparable](p *V, v V) (*V, *V) {
	if p == nil {
		return nil, nil
	}
	if *p == v {
		return nil, p
	}
	return nil, nil
}

// 条件を満たす値の直後で分割したふたつのポインタを返す。
func SplitAfterBy[V any](p *V, f func(V) (bool, error)) (*V, *V, error) {
	if p == nil {
		return nil, nil, nil
	}
	ok, err := f(*p)
	if err != nil {
		return nil, nil, err
	}
	if ok {
		return p, nil, nil
	}
	return nil, nil, nil
}

// 条件を満たす値の直後で分割したふたつのポインタを返す。実行中にエラーが起きた場合 panic する。
func MustSplitAfterBy[V any](p *V, f func(V) (bool, error)) (*V, *V) {
	return must.Must2(SplitAfterBy(p, f))
}

// 一致する値の直後で分割したふたつのポインタを返す。
func SplitAfter[V comparable](p *V, v V) (*V, *V) {
	if p == nil {
		return nil, nil
	}
	if *p == v {
		return p, nil
	}
	return nil, nil
}

// 条件を満たすポインタと満たさないポインタを返す。
func PartitionBy[V any](p *V, f func(V) (bool, error)) (*V, *V, error) {
	if p == nil {
		return nil, nil, nil
	}
	ok, err := f(*p)
	if err != nil {
		return nil, nil, err
	}
	if ok {
		return p, nil, nil
	}
	return nil, p, nil
}

// 条件を満たすポインタと満たさないポインタを返す。実行中にエラーが起きた場合 panic する。
func MustPartitionBy[V any](p *V, f func(V) (bool, error)) (*V, *V) {
	return must.Must2(PartitionBy(p, f))
}

// 値の一致するポインタと一致しないポインタを返す。
func Partition[V comparable](p *V, v V) (*V, *V) {
	if p == nil {
		return nil, nil
	}
	if *p == v {
		return p, nil
	}
	return nil, p
}

// 条件を満たし続ける先頭の値のポインタを返す。
func TakeWhileBy[V any](p *V, f func(V) (bool, error)) (*V, error) {
	if p == nil {
		return nil, nil
	}
	ok, err := f(*p)
	if err != nil {
		return nil, err
	}
	if ok {
		return p, nil
	}
	return nil, nil
}

// 条件を満たし続ける先頭の値のポインタを返す。実行中にエラーが起きた場合 panic する。
func MustTakeWhileBy[V any](p *V, f func(V) (bool, error)) *V {
	return must.Must1(TakeWhileBy(p, f))
}

// 一致し続ける先頭の値のポインタを返す。
func TakeWhile[V comparable](p *V, v V) *V {
	if p == nil {
		return nil
	}
	if *p == v {
		return p
	}
	return nil
}

// 先頭n個の値のポインタを返す。
func Take[V any](p *V, n int) *V {
	if p == nil {
		return nil
	}
	if n <= 0 {
		return nil
	}
	return p
}

// 条件を満たし続ける先頭の値を除いたポインタを返す。
func DropWhileBy[V any](p *V, f func(V) (bool, error)) (*V, error) {
	if p == nil {
		return nil, nil
	}
	ok, err := f(*p)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, nil
	}
	return p, nil
}

// 条件を満たし続ける先頭の値を除いたポインタを返す。実行中にエラーが起きた場合 panic する。
func MustDropWhileBy[V any](p *V, f func(V) (bool, error)) *V {
	return must.Must1(DropWhileBy(p, f))
}

// 一致し続ける先頭の値を除いたポインタを返す。
func DropWhile[V comparable](p *V, v V) *V {
	if p == nil {
		return nil
	}
	if *p == v {
		return nil
	}
	return p
}

// 先頭n個の値を除いたポインタを返す。
func Drop[V any](p *V, n int) *V {
	if p == nil {
		return nil
	}
	if n <= 0 {
		return p
	}
	return nil
}

// 条件を満たし続ける先頭部分と残りの部分、ふたつのポインタを返す。
func SpanBy[V any](p *V, f func(V) (bool, error)) (*V, *V, error) {
	if p == nil {
		return nil, nil, nil
	}
	ok, err := f(*p)
	if err != nil {
		return nil, nil, err
	}
	if ok {
		return p, nil, nil
	}
	return nil, nil, nil
}

// 条件を満たし続ける先頭部分と残りの部分、ふたつのポインタを返す。実行中にエラーが起きた場合 panic する。
func MustSpanBy[V any](p *V, f func(V) (bool, error)) (*V, *V) {
	return must.Must2(SpanBy(p, f))
}

// 一致し続ける先頭部分と残りの部分、ふたつのポインタを返す。
func Span[V comparable](p *V, v V) (*V, *V) {
	if p == nil {
		return nil, nil
	}
	if *p == v {
		return p, nil
	}
	return nil, nil
}

// ゼロ値を除いたポインタを返す。
func Clean[V comparable](p *V) *V {
	if p == nil {
		return nil
	}
	if *p == *new(V) {
		return nil
	}
	return p
}

// 重複を除いたポインタを返す。
func Distinct[V comparable](p *V) *V {
	return p
}

// 条件を満たす値を変換したポインタを返す。
func Collect[V1 any, V2 any](p *V1, f func(V1) (V2, bool, error)) (*V2, error) {
	if p == nil {
		return nil, nil
	}
	v2, ok, err := f(*p)
	if err != nil {
		return nil, err
	}
	if ok {
		return &v2, nil
	}
	return nil, nil
}

// 条件を満たす値を変換したポインタを返す。実行中にエラーが起きた場合 panic する。
func MustCollect[V1 any, V2 any](p *V1, f func(V1) (V2, bool, error)) *V2 {
	return must.Must1(Collect(p, f))
}

// 値と位置をペアにしたポインタを返す。
func ZipWithIndex[V any](p *V) *tuple.T2[V, int] {
	if p == nil {
		return nil
	}
	return &tuple.T2[V, int]{V1: *p, V2: 0}
}

// ２つのポインタの同じ位置の値をペアにしたポインタを返す。
func Zip2[V1 any, V2 any](p1 *V1, p2 *V2) *tuple.T2[V1, V2] {
	if p1 == nil || p2 == nil {
		return nil
	}
	return &tuple.T2[V1, V2]{V1: *p1, V2: *p2}
}

// ３つのポインタの同じ位置の値をペアにしたポインタを返す。
func Zip3[V1 any, V2 any, V3 any](p1 *V1, p2 *V2, p3 *V3) *tuple.T3[V1, V2, V3] {
	if p1 == nil || p2 == nil || p3 == nil {
		return nil
	}
	return &tuple.T3[V1, V2, V3]{V1: *p1, V2: *p2, V3: *p3}
}

// ４つのポインタの同じ位置の値をペアにしたポインタを返す。
func Zip4[V1 any, V2 any, V3 any, V4 any](p1 *V1, p2 *V2, p3 *V3, p4 *V4) *tuple.T4[V1, V2, V3, V4] {
	if p1 == nil || p2 == nil || p3 == nil || p4 == nil {
		return nil
	}
	return &tuple.T4[V1, V2, V3, V4]{V1: *p1, V2: *p2, V3: *p3, V4: *p4}
}

// ５つのポインタの同じ位置の値をペアにしたポインタを返す。
func Zip5[V1 any, V2 any, V3 any, V4 any, V5 any](p1 *V1, p2 *V2, p3 *V3, p4 *V4, p5 *V5) *tuple.T5[V1, V2, V3, V4, V5] {
	if p1 == nil || p2 == nil || p3 == nil || p4 == nil || p5 == nil {
		return nil
	}
	return &tuple.T5[V1, V2, V3, V4, V5]{V1: *p1, V2: *p2, V3: *p3, V4: *p4, V5: *p5}
}

// ６つのポインタの同じ位置の値をペアにしたポインタを返す。
func Zip6[V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](p1 *V1, p2 *V2, p3 *V3, p4 *V4, p5 *V5, p6 *V6) *tuple.T6[V1, V2, V3, V4, V5, V6] {
	if p1 == nil || p2 == nil || p3 == nil || p4 == nil || p5 == nil || p6 == nil {
		return nil
	}
	return &tuple.T6[V1, V2, V3, V4, V5, V6]{V1: *p1, V2: *p2, V3: *p3, V4: *p4, V5: *p5, V6: *p6}
}

// 値のペアを分離して２つのポインタを返す。
func Unzip2[V1 any, V2 any](p *tuple.T2[V1, V2]) (*V1, *V2) {
	if p == nil {
		return nil, nil
	}
	v1, v2 := p.V1, p.V2
	return &v1, &v2
}

// 値のペアを分離して３つのポインタを返す。
func Unzip3[V1 any, V2 any, V3 any](p *tuple.T3[V1, V2, V3]) (*V1, *V2, *V3) {
	if p == nil {
		return nil, nil, nil
	}
	v1, v2, v3 := p.V1, p.V2, p.V3
	return &v1, &v2, &v3
}

// 値のペアを分離して４つのポインタを返す。
func Unzip4[V1 any, V2 any, V3 any, V4 any](p *tuple.T4[V1, V2, V3, V4]) (*V1, *V2, *V3, *V4) {
	if p == nil {
		return nil, nil, nil, nil
	}
	v1, v2, v3, v4 := p.V1, p.V2, p.V3, p.V4
	return &v1, &v2, &v3, &v4
}

// 値のペアを分離して５つのポインタを返す。
func Unzip5[V1 any, V2 any, V3 any, V4 any, V5 any](p *tuple.T5[V1, V2, V3, V4, V5]) (*V1, *V2, *V3, *V4, *V5) {
	if p == nil {
		return nil, nil, nil, nil, nil
	}
	v1, v2, v3, v4, v5 := p.V1, p.V2, p.V3, p.V4, p.V5
	return &v1, &v2, &v3, &v4, &v5
}

// 値のペアを分離して６つのポインタを返す。
func Unzip6[V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](p *tuple.T6[V1, V2, V3, V4, V5, V6]) (*V1, *V2, *V3, *V4, *V5, *V6) {
	if p == nil {
		return nil, nil, nil, nil, nil, nil
	}
	v1, v2, v3, v4, v5, v6 := p.V1, p.V2, p.V3, p.V4, p.V5, p.V6
	return &v1, &v2, &v3, &v4, &v5, &v6
}

// 値ごとに関数の返すキーでグルーピングしたマップを返す。
func GroupBy[K comparable, V any](p *V, f func(V) (K, error)) (map[K][]V, error) {
	if p == nil {
		return nil, nil
	}
	k, err := f(*p)
	if err != nil {
		return nil, err
	}
	return map[K][]V{k: {*p}}, nil
}

// 値ごとに関数の返すキーでグルーピングしたマップを返す。実行中にエラーが起きた場合 panic する。
func MustGroupBy[K comparable, V any](p *V, f func(V) (K, error)) map[K][]V {
	return must.Must1(GroupBy(p, f))
}

// 平坦化したポインタを返す。
func Flatten[V any](p **V) *V {
	if p == nil {
		return nil
	}
	return *p
}

// 値をイテレータに変換し、それらを結合したイテレータを返す。
func FlatMap[V1 any, V2 any](p *V1, f func(V1) (*V2, error)) (*V2, error) {
	if p == nil {
		return nil, nil
	}
	v2, err := f(*p)
	if err != nil {
		return nil, err
	}
	return v2, nil
}

// 値をイテレータに変換し、それらを結合したイテレータを返す。実行中にエラーが起きた場合 panic する。
func MustFlatMap[V1 any, V2 any](p *V1, f func(V1) (*V2, error)) *V2 {
	return must.Must1(FlatMap(p, f))
}

// 値のあいだにseparatorを挿入したイテレータを返す。
func Join[V any](p *V, separator V) *V {
	return p
}
