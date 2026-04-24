package middleware

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	mu      sync.Mutex
	limiter map[string]*rate.Limiter
	r       rate.Limit
	b       int
}

func NewIPRateLimiter(r rate.Limit, b int) func(http.Handler) http.Handler {
	ipLimiter := &IPRateLimiter{
		limiter: make(map[string]*rate.Limiter),
		r:       r,
		b:       b,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter := ipLimiter.getLimiter(r.RemoteAddr)

			if !limiter.Allow() {
				http.Error(w, `{"error": "rate limit exceeded"}`, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (l IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	lim, exists := l.limiter[ip]
	if !exists {
		lim = rate.NewLimiter(l.r, l.b)
		l.limiter[ip] = lim
	}
	return lim
}
