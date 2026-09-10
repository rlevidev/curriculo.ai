package main

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type Visitor struct {
	lastSeen time.Time
	tokens   int
}

type Limiter struct {
	mu       sync.Mutex
	visitors map[string]*Visitor
}

// NewLimiter creates the limiter and starts the periodic cleanup of expired visitors.
func NewLimiter() *Limiter {
	l := &Limiter{visitors: make(map[string]*Visitor)}
	go l.cleanupLoop()
	return l
}

// Allow enforces a per-IP token bucket: 5 tokens, full refill after 1h idle; consumes one token if available.
func (l *Limiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	v, ok := l.visitors[ip]
	if !ok {
		l.visitors[ip] = &Visitor{lastSeen: time.Now(), tokens: 4}
		return true
	}
	if time.Since(v.lastSeen) > time.Hour {
		v.tokens = 5
	}
	v.lastSeen = time.Now()
	if v.tokens <= 0 {
		return false
	}
	v.tokens--
	return true
}

// Reset clears all visitors (test helper).
func (l *Limiter) Reset() {
	l.mu.Lock()
	l.visitors = make(map[string]*Visitor)
	l.mu.Unlock()
}

// SetLastSeenForTest sets lastSeen for an IP to simulate the refill window in tests.
func (l *Limiter) SetLastSeenForTest(ip string, t time.Time) {
	l.mu.Lock()
	if v, ok := l.visitors[ip]; ok {
		v.lastSeen = t
	}
	l.mu.Unlock()
}

// cleanupLoop evicts visitors idle for more than 1h every 10 minutes.
func (l *Limiter) cleanupLoop() {
	for {
		time.Sleep(10 * time.Minute)
		l.mu.Lock()
		for ip, v := range l.visitors {
			if time.Since(v.lastSeen) > time.Hour {
				delete(l.visitors, ip)
			}
		}
		l.mu.Unlock()
	}
}

// getVisitorIP extracts the IP from RemoteAddr (host:port); falls back to the raw value if parsing fails.
func getVisitorIP(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return ip
}

var defaultLimiter = NewLimiter()

// rateLimiter is the middleware that returns 429 when Allow denies the IP.
func rateLimiter(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getVisitorIP(r.RemoteAddr)
		if !defaultLimiter.Allow(ip) {
			loggerFor(r).Warn("rate limited", "op", "rate_limit", "status", http.StatusTooManyRequests)
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}
