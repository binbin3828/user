package ratelimit

import (
	"sync"
	"time"
)

type memoryLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int
	window   time.Duration
}

type visitor struct {
	count    int
	lastSeen time.Time
}

func NewMemoryLimiter(rate int, window time.Duration) Limiter {
	ml := &memoryLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}
	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for range ticker.C {
			ml.mu.Lock()
			for ip, v := range ml.visitors {
				if time.Since(v.lastSeen) > window {
					delete(ml.visitors, ip)
				}
			}
			ml.mu.Unlock()
		}
	}()
	return ml
}

func (ml *memoryLimiter) Allow(key string) bool {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	v, ok := ml.visitors[key]
	if !ok {
		ml.visitors[key] = &visitor{count: 1, lastSeen: time.Now()}
		return true
	}

	if time.Since(v.lastSeen) > ml.window {
		v.count = 1
		v.lastSeen = time.Now()
		return true
	}

	v.lastSeen = time.Now()
	v.count++
	return v.count <= ml.rate
}
