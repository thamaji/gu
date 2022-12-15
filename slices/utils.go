package slices

import (
	"math/rand"

	"golang.org/x/exp/constraints"
)

// 指定した値をn個複製してスライスをつくる。
func Repeat[V any](n int, v V) []V {
	slice := make([]V, n)
	for i := 0; i < n; i++ {
		slice[i] = v
	}
	return slice
}

// 指定した値をn個複製してスライスをつくる。
func RepeatBy[V any](n int, f func() V) []V {
	slice := make([]V, n)
	for i := 0; i < n; i++ {
		slice[i] = f()
	}
	return slice
}

// 連番からスライスをつくる。
func Range[V constraints.Ordered](start V, stop V, step V) []V {
	slice := []V{}
	for i := start; i < stop; i += step {
		slice = append(slice, i)
	}
	return slice
}

// 要素をすべて削除する。
func Clear[V any](slice []V) []V {
	for i := range slice {
		slice[i] = *new(V)
	}
	return slice[:0]
}

// 要素をすべてコピーしたスライスを返す。
func Clone[V any](slice []V) []V {
	clone := make([]V, len(slice))
	copy(clone, slice)
	return clone
}

// 要素をランダムに入れ替える。
func Shuffle[V any](slice []V, r *rand.Rand) {
	var n int
	for i := 0; i < len(slice); i++ {
		n = r.Intn(i + 1)
		slice[i], slice[n] = slice[n], slice[i]
	}
}

// 要素を１つランダムに返す。
func Sample[V any](slice []V, r *rand.Rand) V {
	return slice[r.Intn(len(slice))]
}

// 逆順にしたスライスを返す。
func Reverse[T any](slice []T) []T {
	dst := make([]T, 0, len(slice))
	for i := len(slice) - 1; i >= 0; i-- {
		dst = append(dst, slice[i])
	}
	return dst
}

// すべての要素に関数の実行結果を代入する。
func FillBy[V any](slice []V, f func(int) V) {
	for i := range slice {
		slice[i] = f(i)
	}
}

// すべての要素に値を代入する。
func Fill[V any](slice []V, v V) {
	FillBy(slice, func(int) V { return v })
}

// すべての要素にゼロ値を代入する。
func FillZero[V any](slice []V) {
	Fill(slice, *new(V))
}

// 要素がn個になるまで先頭に関数の実行結果を挿入する。
func PadBy[V any](slice []V, n int, f func(int) V) []V {
	if len(slice) >= n {
		return slice
	}
	c := n - len(slice)
	t := make([]V, c)
	for i := 0; i < len(t); i++ {
		t[i] = f(i)
	}
	return append(t, slice...)
}

// 要素がn個になるまで先頭にvを挿入する。
func Pad[V any](slice []V, n int, v V) []V {
	return PadBy(slice, n, func(int) V { return v })
}

// 要素がn個になるまで先頭にゼロ値を挿入する。
func PadZero[V any](slice []V, n int) []V {
	return Pad(slice, n, *new(V))
}

// 要素がn個になるまで末尾に関数の実行結果を挿入する。
func PadRightBy[V any](slice []V, n int, f func(int) V) []V {
	if len(slice) >= n {
		return slice
	}
	c := n - len(slice)
	t := make([]V, c)
	for i := 0; i < len(t); i++ {
		t[i] = f(len(slice) + i)
	}
	return append(slice, t...)
}

// 要素がn個になるまで末尾にvを挿入する。
func PadRight[V any](slice []V, n int, v V) []V {
	return PadRightBy(slice, n, func(int) V { return v })
}

// 要素がn個になるまで末尾にゼロ値を挿入する。
func PadZeroRight[V any](slice []V, n int) []V {
	return PadRight(slice, n, *new(V))
}
