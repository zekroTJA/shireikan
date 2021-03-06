package ratelimit

import "time"

// Limiter handles and calculates the
// rate limit tokens using a simple
// token bucket system.
type Limiter struct {
	burst       int
	restoration time.Duration

	tokens         int
	lastActivation time.Time
}

// NewLimiter initializes a new limiter with the
// given burst and restoration values.
func NewLimiter(burst int, restoration time.Duration) *Limiter {
	return &Limiter{
		burst:          burst,
		restoration:    restoration,
		tokens:         burst,
		lastActivation: time.Time{},
	}
}

// Take returns true when a token was available
// to be taken. Otherwise, false is returned as
// well as a duration until a next token will be
// available.
func (l *Limiter) Take() (ok bool, next time.Duration) {
	tokens := l.getVirtualTokens()
	if tokens == 0 {
		next = l.restoration - time.Since(l.lastActivation)
		return
	}

	l.tokens = tokens - 1
	l.lastActivation = time.Now()
	ok = true

	return
}

func (l *Limiter) getVirtualTokens() int {
	tokens := int(time.Since(l.lastActivation)/l.restoration) + l.tokens
	if tokens > l.burst {
		return l.burst
	}
	return tokens
}
