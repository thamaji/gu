package try

func Try(f func()) (err error) {
	defer func() {
		if v, ok := recover().(error); ok {
			err = v
		}
	}()
	f()
	return
}

func Try1[V1 any](f func() V1) (v1 V1, err error) {
	defer func() {
		if v, ok := recover().(error); ok {
			err = v
		}
	}()
	v1 = f()
	return
}

func Try2[V1, V2 any](f func() (V1, V2)) (v1 V1, v2 V2, err error) {
	defer func() {
		if v, ok := recover().(error); ok {
			err = v
		}
	}()
	v1, v2 = f()
	return
}

func Try3[V1, V2, V3 any](f func() (V1, V2, V3)) (v1 V1, v2 V2, v3 V3, err error) {
	defer func() {
		if v, ok := recover().(error); ok {
			err = v
		}
	}()
	v1, v2, v3 = f()
	return
}

func Try4[V1, V2, V3, V4 any](f func() (V1, V2, V3, V4)) (v1 V1, v2 V2, v3 V3, v4 V4, err error) {
	defer func() {
		if v, ok := recover().(error); ok {
			err = v
		}
	}()
	v1, v2, v3, v4 = f()
	return
}

func Try5[V1, V2, V3, V4, V5 any](f func() (V1, V2, V3, V4, V5)) (v1 V1, v2 V2, v3 V3, v4 V4, v5 V5, err error) {
	defer func() {
		if v, ok := recover().(error); ok {
			err = v
		}
	}()
	v1, v2, v3, v4, v5 = f()
	return
}

func Try6[V1, V2, V3, V4, V5, V6 any](f func() (V1, V2, V3, V4, V5, V6)) (v1 V1, v2 V2, v3 V3, v4 V4, v5 V5, v6 V6, err error) {
	defer func() {
		if v, ok := recover().(error); ok {
			err = v
		}
	}()
	v1, v2, v3, v4, v5, v6 = f()
	return
}
