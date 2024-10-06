package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter implements a custom rate limiter based on the token bucket algorithm.
// `r` is the rate (tokens per second), and `b` is the burst size (maximum number of tokens).
func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	limiter := NewTokenBucketLimiter(float64(r), b)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.String(http.StatusTooManyRequests, "Rate limit exceeded")
			c.Abort()
			return
		}
		c.Next()
	}
}

// TokenBucketLimiter represents a token bucket rate limiter.
type TokenBucketLimiter struct {
	rate         float64       // tokens per second
	burst        int           // maximum burst size
	tokens       float64       // available tokens
	lastRefill   time.Time     // last time tokens were refilled
	mutex        sync.Mutex    // ensures thread-safe access to the limiter
	refillPeriod time.Duration // duration between token refills
}

// NewTokenBucketLimiter initializes a new rate limiter with the specified rate and burst.
// Set the minimum rate to 1
func NewTokenBucketLimiter(rate float64, burst int) *TokenBucketLimiter {
	if rate < 1 {
		rate = 1
	}

	return &TokenBucketLimiter{
		rate:         rate,
		burst:        burst,
		tokens:       float64(burst),
		lastRefill:   time.Now(),
		refillPeriod: time.Second / time.Duration(rate),
	}
}

// Allow checks if a request can be processed or should be rate-limited.
func (rl *TokenBucketLimiter) Allow() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()

	// Refill tokens based on the elapsed time
	// Cap tokens at burst size
	rl.tokens += elapsed * rl.rate
	if rl.tokens > float64(rl.burst) {
		rl.tokens = float64(rl.burst)
	}
	rl.lastRefill = now

	// Check if there's at least 1 token available
	// Consume one token
	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}

	return false
}
