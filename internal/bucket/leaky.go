package bucket

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	capacity     int
	leakRate     float64   // speed of leak in tokens/sec
	tokens       float64   // текущее количество токенов в ведре
	lastLeakTime time.Time // время последней "утечки"
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

// Add добавляет токен в bucket и проверяет, не превышен ли лимит
func (b *LeakyBucket) Add() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()

	// Рассчитываем сколько токенов "утекло" с момента последнего обращения
	elapsed := now.Sub(b.lastLeakTime).Seconds()
	leak := elapsed * b.leakRate // количество токенов, которые "утекли" за это время

	// Обновляем состояние ведра
	b.tokens = b.tokens - leak
	if b.tokens < 0 {
		b.tokens = 0
	}

	b.lastLeakTime = now

	// Проверяем, можно ли добавить 1 новый запрос
	if b.tokens+1 <= float64(b.capacity) {
		b.tokens++
		return true
	}

	// Ведро переполнено, запрос отклоняется
	return false
}

// Reset сбрасывает состояние bucket
func (b *LeakyBucket) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.tokens = 0
	b.lastLeakTime = time.Now()
}
