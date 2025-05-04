package bucket

import (
	"math"
	"sync"
	"time"
)

type key struct {
	ip           string
	login        string
	passwordHash string
}

type state struct {
	currentLevel float64   // tokens in bucket at the last update moment
	lastLeakTime time.Time // time of the last update
}

type Service struct {
	bucketsMap map[key]state
	mu         sync.Mutex
}

func NewService() *Service {
	return &Service{
		bucketsMap: make(map[key]state),
	}
}

func (s *Service) CheckAuthAttempt(login, passwordHash, ip string) bool {
	limits := []struct {
		key   key
		limit int
	}{
		{key{login: login}, 10},                // 10/мин для логина
		{key{passwordHash: passwordHash}, 100}, // 100/мин для пароля
		{key{ip: ip}, 1000},                    // 1000/мин для IP
	}

	for _, l := range limits {
		leakRate := float64(l.limit) / 60 // Конвертируем в токены/секунду
		if !s.isAllowed(l.key, l.limit, leakRate) {
			return false
		}
	}
	return true
}

func (s *Service) isAllowed(bucketKey key, capacity int, leakRatePerSec float64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	info, exists := s.bucketsMap[bucketKey]

	if !exists {
		info = state{
			lastLeakTime: now,
			currentLevel: 1, // put 1 token for the first attempt
		}
	} else {
		elapsed := now.Sub(info.lastLeakTime).Seconds()
		leaked := elapsed * leakRatePerSec                        // count of leaked tokens for the elapsed time
		info.currentLevel = math.Max(0, info.currentLevel-leaked) // guarantee that current level is not negative
		info.lastLeakTime = now
	}

	if info.currentLevel < float64(capacity) {
		info.currentLevel++ // put 1 token in bucket (for the current attempt)
		s.bucketsMap[bucketKey] = info
		return true
	}

	return false
}
