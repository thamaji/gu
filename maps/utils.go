package maps

import (
	"math/rand"
)

// 要素をすべて削除する。
func Clear[K comparable, V any](m map[K]V) map[K]V {
	for k := range m {
		m[k] = *new(V)
		delete(m, k)
	}
	return m
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

// すべての要素に関数の実行結果を代入する。
func FillBy[K comparable, V any](m map[K]V, f func(K) V) {
	for k := range m {
		m[k] = f(k)
	}
}

// すべての要素に値を代入する。
func Fill[K comparable, V any](m map[K]V, v V) {
	FillBy(m, func(k K) V { return *new(V) })
}

// すべての要素にゼロ値を代入する。
func FillZero[K comparable, V any](m map[K]V) {
	Fill(m, *new(V))
}
