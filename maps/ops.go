package maps

import (
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

// 値ごとに関数を実行する。
func ForEach[K comparable, V any](m map[K]V, f func(K, V)) {
	for k, v := range m {
		f(k, v)
	}
}

// 他のマップと関数で比較し、一致していたらtrueを返す。
func EqualBy[K comparable, V any](m1 map[K]V, m2 map[K]V, f func(V, V) bool) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k1, v1 := range m1 {
		v2, ok := m2[k1]
		if !ok || !f(v1, v2) {
			return false
		}
	}
	return true
}

// 他のマップと一致していたらtrueを返す。
func Equal[K comparable, V comparable](m1 map[K]V, m2 map[K]V) bool {
	return EqualBy(m1, m2, func(v1 V, v2 V) bool { return v1 == v2 })
}

// 条件を満たす値の数を返す。
func CountBy[K comparable, V any](m map[K]V, f func(K, V) bool) int {
	c := 0
	for k, v := range m {
		if f(k, v) {
			c++
		}
	}
	return c
}

// 一致する値の数を返す。
func Count[K comparable, V comparable](m map[K]V, v V) int {
	return CountBy(m, func(k K, v1 V) bool { return v1 == v })
}

// 値を変換したマップを返す。
func Map[K comparable, V1 any, V2 any](m map[K]V1, f func(K, V1) V2) map[K]V2 {
	m2 := make(map[K]V2, len(m))
	for k, v := range m {
		m2[k] = f(k, v)
	}
	return m2
}

// 値を順に演算する。
func Reduce[K comparable, V any](m map[K]V, f func(V, K, V) V) V {
	var ready bool
	var result V
	for k, v := range m {
		if !ready {
			result = v
			ready = true
			continue
		}
		result = f(result, k, v)
	}
	return result
}

// 値の合計を返す。
func Sum[K comparable, V constraints.Ordered | constraints.Complex](m map[K]V) V {
	return Reduce(m, func(sum V, k K, v V) V { return sum + v })
}

// 値を変換して合計を返す。
func SumBy[K comparable, V1 any, V2 constraints.Ordered | constraints.Complex](m map[K]V1, f func(K, V1) V2) V2 {
	return Sum(Map(m, f))
}

// 最大の値を返す。
func Max[K comparable, V constraints.Ordered](m map[K]V) V {
	return Reduce(m, func(max V, k K, v V) V {
		if max < v {
			return v
		}
		return max
	})
}

// 値を変換して最大の値を返す。
func MaxBy[K comparable, V1 any, V2 constraints.Ordered](m map[K]V1, f func(K, V1) V2) V2 {
	return Max(Map(m, f))
}

// 最小の値を返す。
func Min[K comparable, V constraints.Ordered](m map[K]V) V {
	return Reduce(m, func(max V, k K, v V) V {
		if max > v {
			return v
		}
		return max
	})
}

// 値を変換して最小の値を返す。
func MinBy[K comparable, V1 any, V2 constraints.Ordered](m map[K]V1, f func(K, V1) V2) V2 {
	return Min(Map(m, f))
}

// 初期値と値を順に演算する。
func Fold[K comparable, V1 any, V2 any](m map[K]V1, v V2, f func(V2, K, V1) V2) V2 {
	result := v
	for k, v := range m {
		result = f(result, k, v)
	}
	return result
}

// 条件を満たす値を返す。
func FindBy[K comparable, V any](m map[K]V, f func(K, V) bool) (V, bool) {
	for k, v := range m {
		if f(k, v) {
			return v, true
		}
	}
	return *new(V), false
}

// 一致する値を返す。
func Find[K comparable, V comparable](m map[K]V, v V) (V, bool) {
	return FindBy(m, func(k K, v1 V) bool { return v1 == v })
}

// 条件を満たす値が存在したらtrueを返す。
func ExistsBy[K comparable, V any](m map[K]V, f func(K, V) bool) bool {
	_, ok := FindBy(m, f)
	return ok
}

// 一致する値が存在したらtrueを返す。
func Exists[K comparable, V comparable](m map[K]V, v V) bool {
	_, ok := Find(m, v)
	return ok
}

// すべての値が条件を満たせばtrueを返す。
func ForAllBy[K comparable, V any](m map[K]V, f func(K, V) bool) bool {
	return !ExistsBy(m, func(k K, v V) bool { return !f(k, v) })
}

// すべての値が一致したらtrueを返す。
func ForAll[K comparable, V comparable](m map[K]V, v V) bool {
	return ForAllBy(m, func(k K, v1 V) bool { return v1 == v })
}

// ひとつめのoldをnewで置き換えたマップを返す。
func Replace[K comparable, V comparable](m map[K]V, old V, new V) map[K]V {
	done := true
	return Map(m, func(k K, v V) V {
		if done && v == old {
			done = false
			return new
		}
		return v
	})
}

// すべてのoldをnewで置き換えたマップを返す。
func ReplaceAll[K comparable, V comparable](m map[K]V, old V, new V) map[K]V {
	return Map(m, func(k K, v V) V {
		if v == old {
			return new
		}
		return v
	})
}

// 条件を満たす値だけのマップを返す。
func FilterBy[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	m2 := map[K]V{}
	for k, v := range m {
		if f(k, v) {
			m2[k] = v
		}
	}
	return m2
}

// 一致する値だけのマップを返す。
func Filter[K comparable, V comparable](m map[K]V, v V) map[K]V {
	return FilterBy(m, func(k K, v1 V) bool { return v1 == v })
}

// 条件を満たす値を除いたマップを返す。
func FilterNotBy[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	return FilterBy(m, func(k K, v V) bool { return !f(k, v) })
}

// 一致する値を除いたマップを返す。
func FilterNot[K comparable, V comparable](m map[K]V, v V) map[K]V {
	return FilterNotBy(m, func(k K, v1 V) bool { return v1 == v })
}

// 条件を満たすマップと満たさないマップを返す。
func PartitionBy[K comparable, V any](m map[K]V, f func(K, V) bool) (map[K]V, map[K]V) {
	m1, m2 := map[K]V{}, map[K]V{}
	for k, v := range m {
		if f(k, v) {
			m1[k] = v
		} else {
			m2[k] = v
		}
	}
	return m1, m2
}

// 値の一致するイテレータと一致しないイテレータを返す。
func Partition[K comparable, V comparable](m map[K]V, v V) (map[K]V, map[K]V) {
	return PartitionBy(m, func(k K, v1 V) bool { return v1 == v })
}

// ゼロ値を除いたマップを返す。
func Clean[K comparable, V comparable](m map[K]V) map[K]V {
	return FilterNot(m, *new(V))
}

// 重複を除いたマップを返す。
func Distinct[K comparable, V comparable](m map[K]V) map[K]V {
	mv := map[V]struct{}{}
	m2 := map[K]V{}
	for k, v := range m {
		if _, ok := mv[v]; ok {
			continue
		}
		mv[v] = struct{}{}
		m2[k] = v
	}
	return m2
}

// 条件を満たす値を変換したマップを返す。
func Collect[K comparable, V1 any, V2 any](m map[K]V1, f func(K, V1) (V2, bool)) map[K]V2 {
	m2 := map[K]V2{}
	for k, v := range m {
		v, ok := f(k, v)
		if !ok {
			continue
		}
		m2[k] = v
	}
	return m2
}

// ２つのマップの同じキーの値をペアにしたマップを返す。
func Zip2[K comparable, V1 any, V2 any](m1 map[K]V1, m2 map[K]V2) map[K]tuple.T2[V1, V2] {
	m := map[K]tuple.T2[V1, V2]{}
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok {
			continue
		}
		m[k] = tuple.NewT2(v1, v2)
	}
	return m
}

// ３つのマップの同じキーの値をペアにしたマップを返す。
func Zip3[K comparable, V1 any, V2 any, V3 any](m1 map[K]V1, m2 map[K]V2, m3 map[K]V3) map[K]tuple.T3[V1, V2, V3] {
	m := map[K]tuple.T3[V1, V2, V3]{}
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok {
			continue
		}
		v3, ok := m3[k]
		if !ok {
			continue
		}
		m[k] = tuple.NewT3(v1, v2, v3)
	}
	return m
}

// ４つのマップの同じキーの値をペアにしたマップを返す。
func Zip4[K comparable, V1 any, V2 any, V3 any, V4 any](m1 map[K]V1, m2 map[K]V2, m3 map[K]V3, m4 map[K]V4) map[K]tuple.T4[V1, V2, V3, V4] {
	m := map[K]tuple.T4[V1, V2, V3, V4]{}
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
		m[k] = tuple.NewT4(v1, v2, v3, v4)
	}
	return m
}

// ５つのマップの同じキーの値をペアにしたマップを返す。
func Zip5[K comparable, V1 any, V2 any, V3 any, V4 any, V5 any](m1 map[K]V1, m2 map[K]V2, m3 map[K]V3, m4 map[K]V4, m5 map[K]V5) map[K]tuple.T5[V1, V2, V3, V4, V5] {
	m := map[K]tuple.T5[V1, V2, V3, V4, V5]{}
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
		m[k] = tuple.NewT5(v1, v2, v3, v4, v5)
	}
	return m
}

// ６つのマップの同じキーの値をペアにしたマップを返す。
func Zip6[K comparable, V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](m1 map[K]V1, m2 map[K]V2, m3 map[K]V3, m4 map[K]V4, m5 map[K]V5, m6 map[K]V6) map[K]tuple.T6[V1, V2, V3, V4, V5, V6] {
	m := map[K]tuple.T6[V1, V2, V3, V4, V5, V6]{}
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
		m[k] = tuple.NewT6(v1, v2, v3, v4, v5, v6)
	}
	return m
}

// 値のペアを分離して２つのイテレータを返す。
func Unzip2[K comparable, V1 any, V2 any](m map[K]tuple.T2[V1, V2]) (map[K]V1, map[K]V2) {
	m1, m2 := make(map[K]V1, len(m)), make(map[K]V2, len(m))
	for k, v := range m {
		m1[k] = v.V1
		m2[k] = v.V2
	}
	return m1, m2
}

// 値のペアを分離して３つのイテレータを返す。
func Unzip3[K comparable, V1 any, V2 any, V3 any](m map[K]tuple.T3[V1, V2, V3]) (map[K]V1, map[K]V2, map[K]V3) {
	m1, m2, m3 := make(map[K]V1, len(m)), make(map[K]V2, len(m)), make(map[K]V3, len(m))
	for k, v := range m {
		m1[k] = v.V1
		m2[k] = v.V2
		m3[k] = v.V3
	}
	return m1, m2, m3
}

// 値のペアを分離して４つのイテレータを返す。
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

// 値のペアを分離して５つのイテレータを返す。
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

// 値のペアを分離して５つのイテレータを返す。
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
func GroupBy[K1 comparable, K2 comparable, V any](m map[K1]V, f func(K1, V) K2) map[K2][]V {
	m2 := map[K2][]V{}
	for k1, v := range m {
		k2 := f(k1, v)
		if _, ok := m2[k2]; ok {
			m2[k2] = append(m2[k2], v)
		} else {
			m2[k2] = []V{v}
		}
	}
	return m2
}

// 平坦化したマップを返す。
func Flatten[K1 comparable, K2 comparable, V any](m map[K1]map[K2]V) map[K2]V {
	m2 := map[K2]V{}
	for _, mm := range m {
		for k, v := range mm {
			m2[k] = v
		}
	}
	return m2
}

// 値をマップに変換し、それらを結合したマップを返す。
func FlatMap[K1 comparable, K2 comparable, V1 any, V2 any](m map[K1]V1, f func(K1, V1) map[K2]V2) map[K2]V2 {
	return Flatten(Map(m, f))
}
