package maps

// 要素をすべて削除する。
func Clear[K comparable, V any](m map[K]V) map[K]V {
	for k := range m {
		m[k] = *new(V)
		delete(m, k)
	}
	return m
}

// すべての要素に関数の実行結果を代入する。
func FillBy[K comparable, V any](m map[K]V, f func(K) (V, error)) error {
	for k := range m {
		v, err := f(k)
		if err != nil {
			return err
		}
		m[k] = v
	}
	return nil
}

// すべての要素に値を代入する。
func Fill[K comparable, V any](m map[K]V, v V) {
	for k := range m {
		m[k] = v
	}
}

// すべての要素にゼロ値を代入する。
func FillZero[K comparable, V any](m map[K]V) {
	Fill(m, *new(V))
}
