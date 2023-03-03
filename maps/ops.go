package maps

import (
	"math/rand"

	"github.com/thamaji/gu/must"
	"github.com/thamaji/gu/tuple"
	"golang.org/x/exp/constraints"
)

// 値の数を返す。
func Len[K comparable, V any](m map[K]V) int {
	return len(m)
}

// 指定したキーの値を返す。
func Get[K comparable, V any](m map[K]V, k K) (V, bool) {
	v, ok := m[k]
	return v, ok
}

// 指定したキーの値を返す。無い場合はvを返す。
func GetOrElse[K comparable, V any](m map[K]V, k K, v V) V {
	if v, ok := m[k]; ok {
		return v
	}
	return v
}

// 指定したキーの値を返す。無い場合は関数の実行結果を返す。
func GetOrFunc[K comparable, V any](m map[K]V, k K, f func() (V, error)) (V, error) {
	if v, ok := m[k]; ok {
		return v, nil
	}
	return f()
}

// 指定したキーの値を返す。無い場合は関数の実行結果を返す。実行中にエラーが起きた場合 panic する。
func MustGetOrFunc[K comparable, V any](m map[K]V, k K, f func() (V, error)) V {
	return must.Must1(GetOrFunc(m, k, f))
}

// 要素をすべてコピーしたマップを返す。
func Clone[K comparable, V any](m map[K]V) map[K]V {
	dst := make(map[K]V, len(m))
	for k, v := range m {
		dst[k] = v
	}
	return dst
}

// 値を１つランダムに返す。
func Sample[K comparable, V any](m map[K]V, r *rand.Rand) V {
	i, n := 0, r.Intn(len(m))
	var k K
	for k = range m {
		if i == n {
			break
		}
		i++
	}
	return m[k]
}

// 値ごとに関数を実行する。
func ForEach[K comparable, V any](m map[K]V, f func(K, V) error) error {
	for k, v := range m {
		if err := f(k, v); err != nil {
			return err
		}
	}
	return nil
}

// 値ごとに関数を実行する。実行中にエラーが起きた場合 panic する。
func MustForEach[K comparable, V any](m map[K]V, f func(K, V) error) {
	must.Must0(ForEach(m, f))
}

// 他のマップと関数で比較し、一致していたらtrueを返す。
func EqualBy[K comparable, V any](m1 map[K]V, m2 map[K]V, f func(V, V) (bool, error)) (bool, error) {
	if len(m1) != len(m2) {
		return false, nil
	}
	for k1, v1 := range m1 {
		v2, ok := m2[k1]
		if !ok {
			return false, nil
		}

		ok, err := f(v1, v2)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

// 他のマップと関数で比較し、一致していたらtrueを返す。実行中にエラーが起きた場合 panic する。
func MustEqualBy[K comparable, V any](m1 map[K]V, m2 map[K]V, f func(V, V) (bool, error)) bool {
	return must.Must1(EqualBy(m1, m2, f))
}

// 他のマップと一致していたらtrueを返す。
func Equal[K comparable, V comparable](m1 map[K]V, m2 map[K]V) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k1, v1 := range m1 {
		v2, ok := m2[k1]
		if !ok {
			return false
		}
		if v1 != v2 {
			return false
		}
	}
	return true
}

// 条件を満たす値の数を返す。
func CountBy[K comparable, V any](m map[K]V, f func(K, V) (bool, error)) (int, error) {
	c := 0
	for k, v := range m {
		ok, err := f(k, v)
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
func MustCountBy[K comparable, V any](m map[K]V, f func(K, V) (bool, error)) int {
	return must.Must1(CountBy(m, f))
}

// 一致する値の数を返す。
func Count[K comparable, V comparable](m map[K]V, v V) int {
	c := 0
	for _, v1 := range m {
		if v1 == v {
			c++
		}
	}
	return c
}

// 値を変換したマップを返す。
func Map[K comparable, V1 any, V2 any](m map[K]V1, f func(K, V1) (V2, error)) (map[K]V2, error) {
	m2 := make(map[K]V2, len(m))
	for k, v := range m {
		v2, err := f(k, v)
		if err != nil {
			return nil, err
		}
		m2[k] = v2
	}
	return m2, nil
}

// 値を変換したマップを返す。実行中にエラーが起きた場合 panic する。
func MustMap[K comparable, V1 any, V2 any](m map[K]V1, f func(K, V1) (V2, error)) map[K]V2 {
	return must.Must1(Map(m, f))
}

// 値を順に演算する。
func Reduce[K comparable, V any](m map[K]V, f func(V, K, V) (V, error)) (V, error) {
	head := true
	var result V
	var err error
	for k, v := range m {
		if !head {
			result = v
			head = false
			continue
		}
		result, err = f(result, k, v)
		if err != nil {
			return *new(V), err
		}
	}
	return result, nil
}

// 値を順に演算する。実行中にエラーが起きた場合 panic する。
func MustReduce[K comparable, V any](m map[K]V, f func(V, K, V) (V, error)) V {
	return must.Must1(Reduce(m, f))
}

// 値の合計を返す。
func Sum[K comparable, V constraints.Ordered | constraints.Complex](m map[K]V) V {
	head := true
	var sum V
	for _, v := range m {
		if !head {
			sum = v
			head = false
			continue
		}
		sum += v
	}
	return sum
}

// 値を変換して合計を返す。
func SumBy[K comparable, V1 any, V2 constraints.Ordered | constraints.Complex](m map[K]V1, f func(K, V1) (V2, error)) (V2, error) {
	head := true
	var sum V2
	var err error
	for k, v1 := range m {
		if !head {
			sum, err = f(k, v1)
			if err != nil {
				return *new(V2), err
			}
			head = false
			continue
		}

		v2, err := f(k, v1)
		if err != nil {
			return *new(V2), err
		}
		sum += v2
	}
	return sum, nil
}

// 値を変換して合計を返す。実行中にエラーが起きた場合 panic する。
func MustSumBy[K comparable, V1 any, V2 constraints.Ordered | constraints.Complex](m map[K]V1, f func(K, V1) (V2, error)) V2 {
	return must.Must1(SumBy(m, f))
}

// 最大の値を返す。
func Max[K comparable, V constraints.Ordered](m map[K]V) V {
	head := true
	var max V
	for _, v := range m {
		if !head {
			max = v
			head = false
			continue
		}
		if max < v {
			max = v
		}
	}
	return max
}

// 値を変換して最大の値を返す。
func MaxBy[K comparable, V1 any, V2 constraints.Ordered](m map[K]V1, f func(K, V1) (V2, error)) (V2, error) {
	head := true
	var max V2
	var err error
	for k, v1 := range m {
		if !head {
			max, err = f(k, v1)
			if err != nil {
				return *new(V2), err
			}
			head = false
			continue
		}
		v2, err := f(k, v1)
		if err != nil {
			return *new(V2), err
		}
		if max < v2 {
			max = v2
		}
	}
	return max, nil
}

// 値を変換して最大の値を返す。実行中にエラーが起きた場合 panic する。
func MustMaxBy[K comparable, V1 any, V2 constraints.Ordered](m map[K]V1, f func(K, V1) (V2, error)) V2 {
	return must.Must1(MaxBy(m, f))
}

// 最小の値を返す。
func Min[K comparable, V constraints.Ordered](m map[K]V) V {
	head := true
	var min V
	for _, v := range m {
		if !head {
			min = v
			head = false
			continue
		}
		if min > v {
			min = v
		}
	}
	return min
}

// 値を変換して最小の値を返す。
func MinBy[K comparable, V1 any, V2 constraints.Ordered](m map[K]V1, f func(K, V1) (V2, error)) (V2, error) {
	head := true
	var min V2
	var err error
	for k, v1 := range m {
		if !head {
			min, err = f(k, v1)
			if err != nil {
				return *new(V2), err
			}
			head = false
			continue
		}
		v2, err := f(k, v1)
		if err != nil {
			return *new(V2), err
		}
		if min > v2 {
			min = v2
		}
	}
	return min, nil
}

// 値を変換して最小の値を返す。実行中にエラーが起きた場合 panic する。
func MustMinBy[K comparable, V1 any, V2 constraints.Ordered](m map[K]V1, f func(K, V1) (V2, error)) V2 {
	return must.Must1(MinBy(m, f))
}

// 初期値と値を順に演算する。
func Fold[K comparable, V1 any, V2 any](m map[K]V1, v V2, f func(V2, K, V1) (V2, error)) (V2, error) {
	var err error
	for k, v1 := range m {
		v, err = f(v, k, v1)
		if err != nil {
			return *new(V2), err
		}
	}
	return v, nil
}

// 初期値と値を順に演算する。実行中にエラーが起きた場合 panic する。
func MustFold[K comparable, V1 any, V2 any](m map[K]V1, v V2, f func(V2, K, V1) (V2, error)) V2 {
	return must.Must1(Fold(m, v, f))
}

// 条件を満たす値を返す。
func FindBy[K comparable, V any](m map[K]V, f func(K, V) (bool, error)) (V, bool, error) {
	for k, v := range m {
		ok, err := f(k, v)
		if err != nil {
			return *new(V), false, err
		}
		if ok {
			return v, true, nil
		}
	}
	return *new(V), false, nil
}

// 条件を満たす値を返す。実行中にエラーが起きた場合 panic する。
func MustFindBy[K comparable, V any](m map[K]V, f func(K, V) (bool, error)) (V, bool) {
	return must.Must2(FindBy(m, f))
}

// 一致する値を返す。
func Find[K comparable, V comparable](m map[K]V, v V) (V, bool) {
	for _, v1 := range m {
		if v1 == v {
			return v1, true
		}
	}
	return *new(V), false
}

// 条件を満たす値が存在したらtrueを返す。
func ExistsBy[K comparable, V any](m map[K]V, f func(K, V) (bool, error)) (bool, error) {
	for k, v := range m {
		ok, err := f(k, v)
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
func MustExistsBy[K comparable, V any](m map[K]V, f func(K, V) (bool, error)) bool {
	return must.Must1(ExistsBy(m, f))
}

// 一致する値が存在したらtrueを返す。
func Exists[K comparable, V comparable](m map[K]V, v V) bool {
	for _, v1 := range m {
		if v1 == v {
			return true
		}
	}
	return false
}

// すべての値が条件を満たせばtrueを返す。
func ForAllBy[K comparable, V any](m map[K]V, f func(K, V) (bool, error)) (bool, error) {
	for k, v := range m {
		ok, err := f(k, v)
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
func MustForAllBy[K comparable, V any](m map[K]V, f func(K, V) (bool, error)) bool {
	return must.Must1(ForAllBy(m, f))
}

// すべての値が一致したらtrueを返す。
func ForAll[K comparable, V comparable](m map[K]V, v V) bool {
	for _, v1 := range m {
		if v1 != v {
			return false
		}
	}
	return true
}

// ひとつめのoldをnewで置き換えたマップを返す。
func Replace[K comparable, V comparable](m map[K]V, old V, new V) map[K]V {
	done := false
	dst := make(map[K]V, len(m))
	for k, v := range m {
		if !done && v == old {
			dst[k] = new
			done = true
		} else {
			dst[k] = v
		}
	}
	return dst
}

// すべてのoldをnewで置き換えたマップを返す。
func ReplaceAll[K comparable, V comparable](m map[K]V, old V, new V) map[K]V {
	dst := make(map[K]V, len(m))
	for k, v := range m {
		if v == old {
			dst[k] = new
		} else {
			dst[k] = v
		}
	}
	return dst
}

// 条件を満たす値だけのマップを返す。
func FilterBy[K comparable, V any](m map[K]V, f func(K, V) (bool, error)) (map[K]V, error) {
	dst := make(map[K]V, len(m)/2)
	for k, v := range m {
		ok, err := f(k, v)
		if err != nil {
			return nil, err
		}
		if ok {
			dst[k] = v
		}
	}
	return dst, nil
}

// 条件を満たす値だけのマップを返す。実行中にエラーが起きた場合 panic する。
func MustFilterBy[K comparable, V any](m map[K]V, f func(K, V) (bool, error)) map[K]V {
	return must.Must1(FilterBy(m, f))
}

// 一致する値だけのマップを返す。
func Filter[K comparable, V comparable](m map[K]V, v V) map[K]V {
	dst := make(map[K]V, len(m)/2)
	for k, v1 := range m {
		if v1 == v {
			dst[k] = v1
		}
	}
	return dst
}

// 条件を満たす値を除いたマップを返す。
func FilterNotBy[K comparable, V any](m map[K]V, f func(K, V) (bool, error)) (map[K]V, error) {
	dst := make(map[K]V, len(m)/2)
	for k, v := range m {
		ok, err := f(k, v)
		if err != nil {
			return nil, err
		}
		if !ok {
			dst[k] = v
		}
	}
	return dst, nil
}

// 条件を満たす値を除いたマップを返す。実行中にエラーが起きた場合 panic する。
func MustFilterNotBy[K comparable, V any](m map[K]V, f func(K, V) (bool, error)) map[K]V {
	return must.Must1(FilterNotBy(m, f))
}

// 一致する値を除いたマップを返す。
func FilterNot[K comparable, V comparable](m map[K]V, v V) map[K]V {
	dst := make(map[K]V, len(m)/2)
	for k, v1 := range m {
		if v1 != v {
			dst[k] = v1
		}
	}
	return dst
}

// 条件を満たすマップと満たさないマップを返す。
func PartitionBy[K comparable, V any](m map[K]V, f func(K, V) (bool, error)) (map[K]V, map[K]V, error) {
	dst1 := make(map[K]V, len(m)/2)
	dst2 := make(map[K]V, len(m)/2)
	for k, v := range m {
		ok, err := f(k, v)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			dst1[k] = v
		} else {
			dst2[k] = v
		}
	}
	return dst1, dst2, nil
}

// 条件を満たすマップと満たさないマップを返す。実行中にエラーが起きた場合 panic する。
func MustPartitionBy[K comparable, V any](m map[K]V, f func(K, V) (bool, error)) (map[K]V, map[K]V) {
	return must.Must2(PartitionBy(m, f))
}

// 値の一致するイテレータと一致しないイテレータを返す。
func Partition[K comparable, V comparable](m map[K]V, v V) (map[K]V, map[K]V) {
	dst1 := make(map[K]V, len(m)/2)
	dst2 := make(map[K]V, len(m)/2)
	for k, v1 := range m {
		if v1 == v {
			dst1[k] = v1
		} else {
			dst2[k] = v1
		}
	}
	return dst1, dst2
}

// ゼロ値を除いたマップを返す。
func Clean[K comparable, V comparable](m map[K]V) map[K]V {
	zero := *new(V)
	dst := make(map[K]V, len(m)/2)
	for k, v := range m {
		if v == zero {
			continue
		}
		dst[k] = v
	}
	return dst
}

// 重複を除いたマップを返す。
func Distinct[K comparable, V comparable](m map[K]V) map[K]V {
	dst := make(map[K]V, len(m)/2)
	for k, v := range m {
		skip := false
		for _, v2 := range dst {
			if v == v2 {
				skip = true
			}
		}
		if skip {
			continue
		}
		dst[k] = v
	}
	return dst
}

// 条件を満たす値を変換したマップを返す。
func Collect[K comparable, V1 any, V2 any](m map[K]V1, f func(K, V1) (V2, bool, error)) (map[K]V2, error) {
	dst := make(map[K]V2, len(m)/2)
	for k, v1 := range m {
		v2, ok, err := f(k, v1)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		dst[k] = v2
	}
	return dst, nil
}

// 条件を満たす値を変換したマップを返す。実行中にエラーが起きた場合 panic する。
func MustCollect[K comparable, V1 any, V2 any](m map[K]V1, f func(K, V1) (V2, bool, error)) map[K]V2 {
	return must.Must1(Collect(m, f))
}

// ２つのマップの同じキーの値をペアにしたマップを返す。
func Zip2[K comparable, V1 any, V2 any](m1 map[K]V1, m2 map[K]V2) map[K]tuple.T2[V1, V2] {
	n := len(m1)
	if n > len(m2) {
		n = len(m2)
	}
	dst := make(map[K]tuple.T2[V1, V2], n)
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok {
			continue
		}
		dst[k] = tuple.NewT2(v1, v2)
	}
	return dst
}

// ３つのマップの同じキーの値をペアにしたマップを返す。
func Zip3[K comparable, V1 any, V2 any, V3 any](m1 map[K]V1, m2 map[K]V2, m3 map[K]V3) map[K]tuple.T3[V1, V2, V3] {
	n := len(m1)
	if n > len(m2) {
		n = len(m2)
	}
	if n > len(m3) {
		n = len(m3)
	}
	dst := make(map[K]tuple.T3[V1, V2, V3], n)
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok {
			continue
		}
		v3, ok := m3[k]
		if !ok {
			continue
		}
		dst[k] = tuple.NewT3(v1, v2, v3)
	}
	return dst
}

// ４つのマップの同じキーの値をペアにしたマップを返す。
func Zip4[K comparable, V1 any, V2 any, V3 any, V4 any](m1 map[K]V1, m2 map[K]V2, m3 map[K]V3, m4 map[K]V4) map[K]tuple.T4[V1, V2, V3, V4] {
	n := len(m1)
	if n > len(m2) {
		n = len(m2)
	}
	if n > len(m3) {
		n = len(m3)
	}
	if n > len(m4) {
		n = len(m4)
	}
	dst := make(map[K]tuple.T4[V1, V2, V3, V4], n)
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok {
			continue
		}
		v3, ok := m3[k]
		if !ok {
			continue
		}
		v4, ok := m4[k]
		if !ok {
			continue
		}
		dst[k] = tuple.NewT4(v1, v2, v3, v4)
	}
	return dst
}

// ５つのマップの同じキーの値をペアにしたマップを返す。
func Zip5[K comparable, V1 any, V2 any, V3 any, V4 any, V5 any](m1 map[K]V1, m2 map[K]V2, m3 map[K]V3, m4 map[K]V4, m5 map[K]V5) map[K]tuple.T5[V1, V2, V3, V4, V5] {
	n := len(m1)
	if n > len(m2) {
		n = len(m2)
	}
	if n > len(m3) {
		n = len(m3)
	}
	if n > len(m4) {
		n = len(m4)
	}
	if n > len(m5) {
		n = len(m5)
	}
	dst := make(map[K]tuple.T5[V1, V2, V3, V4, V5], n)
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok {
			continue
		}
		v3, ok := m3[k]
		if !ok {
			continue
		}
		v4, ok := m4[k]
		if !ok {
			continue
		}
		v5, ok := m5[k]
		if !ok {
			continue
		}
		dst[k] = tuple.NewT5(v1, v2, v3, v4, v5)
	}
	return dst
}

// ６つのマップの同じキーの値をペアにしたマップを返す。
func Zip6[K comparable, V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](m1 map[K]V1, m2 map[K]V2, m3 map[K]V3, m4 map[K]V4, m5 map[K]V5, m6 map[K]V6) map[K]tuple.T6[V1, V2, V3, V4, V5, V6] {
	n := len(m1)
	if n > len(m2) {
		n = len(m2)
	}
	if n > len(m3) {
		n = len(m3)
	}
	if n > len(m4) {
		n = len(m4)
	}
	if n > len(m5) {
		n = len(m5)
	}
	dst := make(map[K]tuple.T6[V1, V2, V3, V4, V5, V6], n)
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok {
			continue
		}
		v3, ok := m3[k]
		if !ok {
			continue
		}
		v4, ok := m4[k]
		if !ok {
			continue
		}
		v5, ok := m5[k]
		if !ok {
			continue
		}
		v6, ok := m6[k]
		if !ok {
			continue
		}
		dst[k] = tuple.NewT6(v1, v2, v3, v4, v5, v6)
	}
	return dst
}

// 値のペアを分離して２つのマップを返す。
func Unzip2[K comparable, V1 any, V2 any](m map[K]tuple.T2[V1, V2]) (map[K]V1, map[K]V2) {
	m1, m2 := make(map[K]V1, len(m)), make(map[K]V2, len(m))
	for k, v := range m {
		m1[k] = v.V1
		m2[k] = v.V2
	}
	return m1, m2
}

// 値のペアを分離して３つのマップを返す。
func Unzip3[K comparable, V1 any, V2 any, V3 any](m map[K]tuple.T3[V1, V2, V3]) (map[K]V1, map[K]V2, map[K]V3) {
	m1, m2, m3 := make(map[K]V1, len(m)), make(map[K]V2, len(m)), make(map[K]V3, len(m))
	for k, v := range m {
		m1[k] = v.V1
		m2[k] = v.V2
		m3[k] = v.V3
	}
	return m1, m2, m3
}

// 値のペアを分離して４つのマップを返す。
func Unzip4[K comparable, V1 any, V2 any, V3 any, V4 any](m map[K]tuple.T4[V1, V2, V3, V4]) (map[K]V1, map[K]V2, map[K]V3, map[K]V4) {
	m1, m2, m3, m4 := make(map[K]V1, len(m)), make(map[K]V2, len(m)), make(map[K]V3, len(m)), make(map[K]V4, len(m))
	for k, v := range m {
		m1[k] = v.V1
		m2[k] = v.V2
		m3[k] = v.V3
		m4[k] = v.V4
	}
	return m1, m2, m3, m4
}

// 値のペアを分離して５つのマップを返す。
func Unzip5[K comparable, V1 any, V2 any, V3 any, V4 any, V5 any](m map[K]tuple.T5[V1, V2, V3, V4, V5]) (map[K]V1, map[K]V2, map[K]V3, map[K]V4, map[K]V5) {
	m1, m2, m3, m4, m5 := make(map[K]V1, len(m)), make(map[K]V2, len(m)), make(map[K]V3, len(m)), make(map[K]V4, len(m)), make(map[K]V5, len(m))
	for k, v := range m {
		m1[k] = v.V1
		m2[k] = v.V2
		m3[k] = v.V3
		m4[k] = v.V4
		m5[k] = v.V5
	}
	return m1, m2, m3, m4, m5
}

// 値のペアを分離して５つのマップを返す。
func Unzip6[K comparable, V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](m map[K]tuple.T6[V1, V2, V3, V4, V5, V6]) (map[K]V1, map[K]V2, map[K]V3, map[K]V4, map[K]V5, map[K]V6) {
	m1, m2, m3, m4, m5, m6 := make(map[K]V1, len(m)), make(map[K]V2, len(m)), make(map[K]V3, len(m)), make(map[K]V4, len(m)), make(map[K]V5, len(m)), make(map[K]V6, len(m))
	for k, v := range m {
		m1[k] = v.V1
		m2[k] = v.V2
		m3[k] = v.V3
		m4[k] = v.V4
		m5[k] = v.V5
		m6[k] = v.V6
	}
	return m1, m2, m3, m4, m5, m6
}

// 値ごとに関数の返すキーでグルーピングしたマップを返す。
func GroupBy[K1 comparable, K2 comparable, V any](m map[K1]V, f func(K1, V) (K2, error)) (map[K2][]V, error) {
	dst := map[K2][]V{}
	for k1, v := range m {
		k2, err := f(k1, v)
		if err != nil {
			return nil, err
		}
		dst[k2] = append(dst[k2], v)
	}
	return dst, nil
}

// 値ごとに関数の返すキーでグルーピングしたマップを返す。実行中にエラーが起きた場合 panic する。
func MustGroupBy[K1 comparable, K2 comparable, V any](m map[K1]V, f func(K1, V) (K2, error)) map[K2][]V {
	return must.Must1(GroupBy(m, f))
}
