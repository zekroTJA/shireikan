package ratelimit

import (
	"testing"
	"time"
)

func TestTake(t *testing.T) {
	const burst = 3
	const restoration = time.Second

	l := NewLimiter(burst, restoration)

	for i := 0; i < burst; i++ {
		if ok, _ := l.Take(); !ok {
			t.Errorf("token %d not retrieved", i)
		}
	}

	if ok, _ := l.Take(); ok {
		t.Error("token retrieved even tho no tokens are available")
	}

	time.Sleep(2 * restoration)

	for i := 0; i < burst-1; i++ {
		if ok, _ := l.Take(); !ok {
			t.Errorf("token %d not retrieved", i)
		}
	}

	ok, next := l.Take()
	if ok {
		t.Error("token retrieved even tho no tokens are available")
	}
	if next > restoration || next < restoration-100*time.Microsecond {
		t.Errorf("returned next value is not in error margin: %s", next)
	}
}
