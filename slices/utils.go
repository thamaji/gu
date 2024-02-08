package slices

import "math/rand"

// 要素をすべて削除する。
func Clear[S ~[]V, V any](slice S) S {
	for i := range slice {
		slice[i] = *new(V)
	}
	return slice[:0]
}

// 要素をランダムに入れ替える。
func Shuffle[V any](slice []V, r *rand.Rand) {
	var n int
	for i := 0; i < len(slice); i++ {
		n = r.Intn(i + 1)
		slice[i], slice[n] = slice[n], slice[i]
	}
}

// すべての要素に関数の実行結果を代入する。
func FillBy[V any](slice []V, f func(int) (V, error)) error {
	for i := range slice {
		v, err := f(i)
		if err != nil {
			return err
		}
		slice[i] = v
	}
	return nil
}

// すべての要素に値を代入する。
func Fill[V any](slice []V, v V) {
	for i := range slice {
		slice[i] = v
	}
}

// すべての要素にゼロ値を代入する。
func FillZero[V any](slice []V) {
	Fill(slice, *new(V))
}
