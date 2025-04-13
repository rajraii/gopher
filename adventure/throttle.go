package adventure

import (
	"sync"
	"time"
)

// Throttler provides a way to limit the rate of function calls
type Throttler struct {
	mu       sync.Mutex
	lastCall time.Time
	interval time.Duration
}

// NewThrottler creates a new throttler with the specified interval
func NewThrottler(interval time.Duration) *Throttler {
	return &Throttler{
		interval: interval,
	}
}

// Throttle executes the given function only if the specified interval has passed
// since the last call. Returns true if the function was executed, false if it was throttled.
func (t *Throttler) Throttle(f func()) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if now.Sub(t.lastCall) < t.interval {
		return false
	}

	t.lastCall = now
	f()
	return true
}

// Reset resets the throttler's last call time
func (t *Throttler) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastCall = time.Time{}
}

// SetInterval updates the throttling interval
func (t *Throttler) SetInterval(interval time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.interval = interval
}
