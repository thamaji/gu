package ptr

// エラーがあるときpanicにする。
func Must[V1 any](v1 V1, err error) V1 {
	if err != nil {
		panic(err)
	}
	return v1
}

// エラーがあるときpanicにする。
func Must2[V1 any, V2 any](v1 V1, v2 V2, err error) (V1, V2) {
	if err != nil {
		panic(err)
	}
	return v1, v2
}

// エラーがあるときpanicにする。
func Must3[V1 any, V2 any, V3 any](v1 V1, v2 V2, v3 V3, err error) (V1, V2, V3) {
	if err != nil {
		panic(err)
	}
	return v1, v2, v3
}

// エラーがあるときpanicにする。
func Must4[V1 any, V2 any, V3 any, V4 any](v1 V1, v2 V2, v3 V3, v4 V4, err error) (V1, V2, V3, V4) {
	if err != nil {
		panic(err)
	}
	return v1, v2, v3, v4
}

// エラーがあるときpanicにする。
func Must5[V1 any, V2 any, V3 any, V4 any, V5 any](v1 V1, v2 V2, v3 V3, v4 V4, v5 V5, err error) (V1, V2, V3, V4, V5) {
	if err != nil {
		panic(err)
	}
	return v1, v2, v3, v4, v5
}

// エラーがあるときpanicにする。
func Must6[V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](v1 V1, v2 V2, v3 V3, v4 V4, v5 V5, v6 V6, err error) (V1, V2, V3, V4, V5, V6) {
	if err != nil {
		panic(err)
	}
	return v1, v2, v3, v4, v5, v6
}
