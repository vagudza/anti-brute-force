package bucket

import (
	"context"
	"sync"
	"time"

	"github.com/vagudza/anti-brute-force/internal/config"
	"go.uber.org/zap"
)

const cleanupInterval = 5 * time.Minute

type MemoryBucketStorage struct {
	logger          *zap.Logger
	buckets         map[string]*LeakyBucket
	capacity        int
	leakRate        float64
	ttl             time.Duration
	cleanupInterval time.Duration
	mu              sync.RWMutex
	stopCleanup     chan struct{}
}

func NewMemoryBucketStorage(cfg *config.LimiterConfig, logger *zap.Logger) *MemoryBucketStorage {
	storage := &MemoryBucketStorage{
		logger:          logger,
		buckets:         make(map[string]*LeakyBucket),
		capacity:        cfg.MaxAttemptsPerMinute,
		leakRate:        float64(cfg.MaxAttemptsPerMinute) / 60.0,
		ttl:             cfg.TTL,
		cleanupInterval: cleanupInterval,
		mu:              sync.RWMutex{},
		stopCleanup:     make(chan struct{}),
	}

	go storage.cleanup()

	return storage
}

func (s *MemoryBucketStorage) Allow(_ context.Context, key string) (bool, error) {
	s.mu.Lock()
	bucket, exists := s.buckets[key]
	if !exists {
		bucket = NewLeakyBucket(s.capacity, s.leakRate)
		s.buckets[key] = bucket
	}
	s.mu.Unlock()

	return bucket.Add(), nil
}

func (s *MemoryBucketStorage) Reset(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if bucket, exists := s.buckets[key]; exists {
		bucket.Reset()
	}

	return nil
}

func (s *MemoryBucketStorage) Close(_ context.Context) error {
	close(s.stopCleanup)
	return nil
}

func (s *MemoryBucketStorage) cleanup() {
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.removeStale()
		case <-s.stopCleanup:
			s.logger.Info("stopping cleanup buckets goroutine")
			return
		}
	}
}

// removeStale removes inactive buckets.
func (s *MemoryBucketStorage) removeStale() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Info("Removing buckets that not used longer than TTL")
	now := time.Now()
	for key, bucket := range s.buckets {
		if now.Sub(bucket.lastLeakTime) > s.ttl {
			s.logger.Info("Removing stale bucket", zap.String("key", key))
			delete(s.buckets, key)
		}
	}
}
