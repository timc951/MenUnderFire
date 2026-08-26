package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type visitor struct {
	count       int
	windowStart time.Time
}

// RateLimiter provides IP-based rate limiting using a fixed window.
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
	// trustProxyHeaders controls whether client-supplied forwarding headers are
	// honoured when identifying a caller. Enable it only when the service sits
	// behind a proxy that overwrites these headers; otherwise any client can
	// spoof them and mint a fresh rate-limit bucket per request.
	trustProxyHeaders bool
}

// NewRateLimiter creates a rate limiter that allows `limit` requests per `window` per IP.
// See RateLimiter.trustProxyHeaders for the meaning of trustProxyHeaders.
func NewRateLimiter(limit int, window time.Duration, trustProxyHeaders bool) *RateLimiter {
	rl := &RateLimiter{
		visitors:          make(map[string]*visitor),
		limit:             limit,
		window:            window,
		trustProxyHeaders: trustProxyHeaders,
	}
	// Periodically clean up stale entries
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(rl.window)
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			if now.Sub(v.windowStart) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// extractIP gets the client IP, consulting proxy headers only when the limiter
// is configured to trust them.
func (rl *RateLimiter) extractIP(r *http.Request) string {
	if rl.trustProxyHeaders {
		if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
			return ip
		}
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			return strings.TrimSpace(strings.Split(forwarded, ",")[0])
		}
		if ip := r.Header.Get("X-Real-Ip"); ip != "" {
			return ip
		}
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

// Middleware returns an HTTP middleware that enforces the rate limit.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := rl.extractIP(r)

		rl.mu.Lock()
		now := time.Now()
		v, exists := rl.visitors[ip]
		if !exists || now.Sub(v.windowStart) > rl.window {
			rl.visitors[ip] = &visitor{count: 1, windowStart: now}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		v.count++
		if v.count > rl.limit {
			rl.mu.Unlock()
			http.Error(w, `{"error":"Rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
