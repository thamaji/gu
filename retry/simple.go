package retry

import (
	"time"
)

func Simple() *simple {
	return &simple{
		MaxRetry: 10,
		Delay:    0,
	}
}

type simple struct {
	MaxRetry int
	Delay    time.Duration
}

func (simple *simple) WithMaxRetry(maxRetry int) *simple {
	simple.MaxRetry = maxRetry
	return simple
}

func (simple *simple) WithDelay(delay time.Duration) *simple {
	simple.Delay = delay
	return simple
}

func (simple simple) call(n int, f func() bool) bool {
	if f() {
		return false
	}

	if n >= simple.MaxRetry {
		return false
	}

	if simple.Delay > 0 {
		time.Sleep(simple.Delay)
	}

	return true
}
