package retry

import (
	"math/rand"
	"time"
)

func ExponentialBackoff() *exponentialBackoff {
	return &exponentialBackoff{
		MaxRetry:   10,
		Delay:      time.Second,
		MaxBackoff: 64 * time.Second,
	}
}

type exponentialBackoff struct {
	MaxRetry   int
	Delay      time.Duration
	MaxBackoff time.Duration
}

func (backoff *exponentialBackoff) WithMaxRetry(maxRetry int) *exponentialBackoff {
	backoff.MaxRetry = maxRetry
	return backoff
}

func (backoff *exponentialBackoff) WithDelay(delay time.Duration) *exponentialBackoff {
	backoff.Delay = delay
	return backoff
}

func (backoff *exponentialBackoff) WithMaxBackoff(maxBackoff time.Duration) *exponentialBackoff {
	backoff.Delay = maxBackoff
	return backoff
}

func (backoff exponentialBackoff) call(n int, f func() bool) bool {
	if f() {
		return false
	}

	if n >= backoff.MaxRetry {
		return false
	}

	var delay time.Duration
	if n == 0 {
		delay = 1 * backoff.Delay
	} else {
		delay = (1 << n) * backoff.Delay
	}
	if delay > backoff.MaxBackoff {
		delay = backoff.MaxBackoff
	}

	time.Sleep(delay + time.Duration(rand.Int63n(backoff.Delay.Nanoseconds()))*time.Nanosecond)

	return true
}
