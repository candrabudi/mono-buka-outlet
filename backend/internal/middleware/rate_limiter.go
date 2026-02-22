package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	attempts int
	resetAt  time.Time
	lockedAt *time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*rateLimitEntry
	max     int
	window  time.Duration
	lockout time.Duration
}

func NewRateLimiter(maxAttempts int, window, lockout time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*rateLimitEntry),
		max:     maxAttempts,
		window:  window,
		lockout: lockout,
	}
	// Cleanup goroutine
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			rl.cleanup()
		}
	}()
	return rl
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	for key, entry := range rl.entries {
		if now.After(entry.resetAt) && (entry.lockedAt == nil || now.After(entry.lockedAt.Add(rl.lockout))) {
			delete(rl.entries, key)
		}
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP() + ":" + c.Request.URL.Path

		rl.mu.Lock()
		entry, exists := rl.entries[key]
		now := time.Now()

		if !exists || now.After(entry.resetAt) {
			entry = &rateLimitEntry{
				attempts: 0,
				resetAt:  now.Add(rl.window),
			}
			rl.entries[key] = entry
		}

		// Check if locked out
		if entry.lockedAt != nil {
			remaining := entry.lockedAt.Add(rl.lockout).Sub(now)
			if remaining > 0 {
				rl.mu.Unlock()
				c.JSON(http.StatusTooManyRequests, gin.H{
					"success":        false,
					"error":          "Terlalu banyak percobaan. Silakan coba lagi nanti.",
					"retry_after_ms": remaining.Milliseconds(),
				})
				c.Abort()
				return
			}
			// Lockout expired, reset
			entry.lockedAt = nil
			entry.attempts = 0
			entry.resetAt = now.Add(rl.window)
		}

		entry.attempts++
		if entry.attempts > rl.max {
			lockedNow := now
			entry.lockedAt = &lockedNow
			rl.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success":        false,
				"error":          "Terlalu banyak percobaan. Silakan coba lagi nanti.",
				"retry_after_ms": rl.lockout.Milliseconds(),
			})
			c.Abort()
			return
		}
		rl.mu.Unlock()

		c.Next()
	}
}
