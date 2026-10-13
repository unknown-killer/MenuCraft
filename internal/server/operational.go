package server

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"
)

type loginWindow struct {
	count int
	until time.Time
}
type loginLimiter struct {
	mu      sync.Mutex
	entries map[string]loginWindow
}

func newLoginLimiter() *loginLimiter { return &loginLimiter{entries: map[string]loginWindow{}} }

func (l *loginLimiter) allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.entries) >= 4096 {
		for key, entry := range l.entries {
			if !now.Before(entry.until) {
				delete(l.entries, key)
			}
		}
	}
	e, exists := l.entries[key]
	if !exists && len(l.entries) >= 4096 {
		return false
	}
	if !now.Before(e.until) {
		e = loginWindow{until: now.Add(time.Minute)}
	}
	e.count++
	l.entries[key] = e
	return e.count <= 20
}

func (s *Server) limitAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Forwarded headers are deliberately not trusted. Operators behind a
		// reverse proxy should also apply per-client limits at that proxy.
		peer, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			peer = r.RemoteAddr
		}
		if !s.authLimits.allow(peer, time.Now()) {
			w.Header().Set("Retry-After", "60")
			http.Error(w, "Too many attempts. Try again in one minute. / یک دقیقه دیگر تلاش کنید.", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleReadiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := s.store.DB().PingContext(ctx); err != nil {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ready"))
}
