package limiter

import (
	"sync"
	"time"
)

type localBucket struct {
	tokens     float64
	lastRefill time.Time
}

type LocalLimiter struct {
	mu         sync.Mutex
	capacity   float64
	refillRate float64
	buckets    map[string]*localBucket
}

func NewLocalLimiter(capacity, refillRate float64) *LocalLimiter {
	return &LocalLimiter{
		capacity:   capacity,
		refillRate: refillRate,
		buckets:    make(map[string]*localBucket),
	}
}

func (l *LocalLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	b, ok := l.buckets[key]
	if !ok {
		l.buckets[key] = &localBucket{
			tokens:     l.capacity - 1,
			lastRefill: now,
		}
		return true
	}

	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens = min(l.capacity, b.tokens+elapsed*l.refillRate)
	b.lastRefill = now

	if b.tokens < 1 {
		return false
	}

	b.tokens--
	return true
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
