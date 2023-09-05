package retry

type Strategy interface {
	call(int, func() bool) bool
}

var defaultStrategy Strategy = Simple()

func SetDefaultStrategy(s Strategy) {
	defaultStrategy = s
}

// リトライ戦略に従って、関数がエラーを返さなくなるまで再試行する
func Run(s Strategy, f func() error) (err error) {
	if s == nil {
		s = defaultStrategy
	}
	for c := 0; s.call(c, func() bool { err = f(); return err == nil }); c++ {
	}
	return
}

// リトライ戦略に従って、関数がエラーを返さなくなるまで再試行する
func Run1[V1 any](s Strategy, f func() (V1, error)) (v1 V1, err error) {
	if s == nil {
		s = defaultStrategy
	}
	for c := 0; s.call(c, func() bool { v1, err = f(); return err == nil }); c++ {
	}
	return
}

// リトライ戦略に従って、関数がエラーを返さなくなるまで再試行する
func Run2[V1 any, V2 any](s Strategy, f func() (V1, V2, error)) (v1 V1, v2 V2, err error) {
	if s == nil {
		s = defaultStrategy
	}
	for c := 0; s.call(c, func() bool { v1, v2, err = f(); return err == nil }); c++ {
	}
	return
}

// リトライ戦略に従って、関数がエラーを返さなくなるまで再試行する
func Run3[V1 any, V2 any, V3 any](s Strategy, f func() (V1, V2, V3, error)) (v1 V1, v2 V2, v3 V3, err error) {
	if s == nil {
		s = defaultStrategy
	}
	for c := 0; s.call(c, func() bool { v1, v2, v3, err = f(); return err == nil }); c++ {
	}
	return
}

// リトライ戦略に従って、関数がエラーを返さなくなるまで再試行する
func Run4[V1 any, V2 any, V3 any, V4 any](s Strategy, f func() (V1, V2, V3, V4, error)) (v1 V1, v2 V2, v3 V3, v4 V4, err error) {
	if s == nil {
		s = defaultStrategy
	}
	for c := 0; s.call(c, func() bool { v1, v2, v3, v4, err = f(); return err == nil }); c++ {
	}
	return
}

// リトライ戦略に従って、関数がエラーを返さなくなるまで再試行する
func Run5[V1 any, V2 any, V3 any, V4 any, V5 any](s Strategy, f func() (V1, V2, V3, V4, V5, error)) (v1 V1, v2 V2, v3 V3, v4 V4, v5 V5, err error) {
	if s == nil {
		s = defaultStrategy
	}
	for c := 0; s.call(c, func() bool { v1, v2, v3, v4, v5, err = f(); return err == nil }); c++ {
	}
	return
}

// リトライ戦略に従って、関数がエラーを返さなくなるまで再試行する
func Run6[V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](s Strategy, f func() (V1, V2, V3, V4, V5, V6, error)) (v1 V1, v2 V2, v3 V3, v4 V4, v5 V5, v6 V6, err error) {
	if s == nil {
		s = defaultStrategy
	}
	for c := 0; s.call(c, func() bool { v1, v2, v3, v4, v5, v6, err = f(); return err == nil }); c++ {
	}
	return
}
