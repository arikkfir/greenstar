package token

import (
	"sync"
	"time"
)

// Bucket represents a token bucket that regulates the flow of tokens at a specified rate.
// It has a certain capacity of tokens that gets refilled in specified time intervals.
type Bucket struct {
	tokens       int
	capacity     int
	fillInterval time.Duration
	lastFill     time.Time
	mutex        sync.Mutex
}

// NewBucket creates a new Bucket with a given capacity and fill interval.
// It initializes the token bucket with the max capacity and sets the last fill time to now.
func NewBucket(capacity int, fillInterval time.Duration) *Bucket {
	return &Bucket{
		tokens:       capacity,
		capacity:     capacity,
		fillInterval: fillInterval,
		lastFill:     time.Now(),
	}
}

// AcquireToken attempts to acquire a token from the bucket.
// If no tokens are available, it waits until a token becomes available.
func (tb *Bucket) AcquireToken() {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	for {
		now := time.Now()

		// Determine number of tokens we "deserve" to add as the amount of time passed since last fill, divided by the
		// fill interval; e.g. 5min passed since last fill, and fill interval is 1min - we deserve 5 tokens.
		timeSinceLastFill := now.Sub(tb.lastFill)
		numberOfTokensWeDeserveSinceLastFill := int(timeSinceLastFill / tb.fillInterval)

		// If indeed we deserve more tokens, add them to the bucket
		if numberOfTokensWeDeserveSinceLastFill > 0 {
			tb.tokens += numberOfTokensWeDeserveSinceLastFill
			if tb.tokens > tb.capacity {
				tb.tokens = tb.capacity
			}
			tb.lastFill = now
		}

		// If we have a token available, grab one and return
		if tb.tokens > 0 {
			tb.tokens--
			return
		}

		// Otherwise, no tokens - sleep a bit and try again
		time.Sleep(1 * time.Second)
	}
}

// fillTokensIfNecessary refills the token bucket if enough time has elapsed since the last fill.
// It calculates the number of tokens to add based on the time passed and adds them up to the bucket's capacity.
func (tb *Bucket) fillTokensIfNecessary() {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	tokensToAdd := int(now.Sub(tb.lastFill) / tb.fillInterval)
	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastFill = now
	}
}
