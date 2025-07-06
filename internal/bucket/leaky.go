package bucket

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	capacity     int
	leakRate     float64 // speed of leak in tokens/sec
	tokens       float64
	lastLeakTime time.Time
	mu           sync.Mutex
}

func NewLeakyBucket(capacity int, leakRate float64) *LeakyBucket {
	return &LeakyBucket{
		capacity:     capacity,
		leakRate:     leakRate,
		tokens:       0,
		lastLeakTime: time.Now(),
	}
}

func (b *LeakyBucket) Add() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()

	elapsed := now.Sub(b.lastLeakTime).Seconds()
	leak := elapsed * b.leakRate

	b.tokens -= leak
	if b.tokens < 0 {
		b.tokens = 0
	}

	b.lastLeakTime = now
	if b.tokens+1 <= float64(b.capacity) {
		b.tokens++
		return true
	}

	return false
}

func (b *LeakyBucket) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.tokens = 0
	b.lastLeakTime = time.Now()
}
