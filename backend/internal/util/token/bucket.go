package token

import (
	"sync"
	"time"
)

type Bucket struct {
	tokens       int
	capacity     int
	fillInterval time.Duration
	lastFill     time.Time
	mutex        sync.Mutex
}

func NewBucket(capacity int, fillInterval time.Duration) *Bucket {
	return &Bucket{
		tokens:       capacity,
		capacity:     capacity,
		fillInterval: fillInterval,
		lastFill:     time.Now(),
	}
}

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
