package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter kiểm soát số request theo IP
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

// visitor lưu thời điểm request cuối cùng & số request trong khoảng thời gian
type visitor struct {
	lastSeen time.Time
	requests int
}

// NewRateLimiter tạo Rate Limiter mới
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}
}

// Middleware áp dụng rate limiting
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		rl.mu.Lock()
		v, exists := rl.visitors[ip]

		if !exists || time.Since(v.lastSeen) > rl.window {
			v = &visitor{lastSeen: time.Now(), requests: 1}
			rl.visitors[ip] = v
		} else {
			v.requests++
			if v.requests > rl.limit {
				rl.mu.Unlock()
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
		}

		v.lastSeen = time.Now()
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
