package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/rafialfauzan/layrd/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiterStore struct {
	mu      sync.Mutex
	clients map[string]*ipLimiter
	rps     rate.Limit
	burst   int
}

func NewRateLimiterStore(rps float64, burst int) *RateLimiterStore {
	store := &RateLimiterStore{
		clients: make(map[string]*ipLimiter),
		rps:     rate.Limit(rps),
		burst:   burst,
	}

	// Periodically cleanup stale client records every 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			store.mu.Lock()
			for ip, client := range store.clients {
				if time.Since(client.lastSeen) > 10*time.Minute {
					delete(store.clients, ip)
				}
			}
			store.mu.Unlock()
		}
	}()

	return store
}

func (s *RateLimiterStore) getLimiter(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	client, exists := s.clients[ip]
	if !exists {
		limiter := rate.NewLimiter(s.rps, s.burst)
		s.clients[ip] = &ipLimiter{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	client.lastSeen = time.Now()
	return client.limiter
}

// RateLimit creates a Gin middleware limiting requests per IP address
func RateLimit(rps float64, burst int) gin.HandlerFunc {
	store := NewRateLimiterStore(rps, burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := store.getLimiter(ip)

		if !limiter.Allow() {
			response.Error(c, http.StatusTooManyRequests, "Too many requests. Please slow down.", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
