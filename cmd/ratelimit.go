package cmd

import (
	"math"
	"math/rand"
	"time"
)

// RateLimiter implements exponential backoff for GitHub API rate limiting
type RateLimiter struct {
	maxRetries  int
	baseDelay   time.Duration
	maxDelay    time.Duration
	jitter      bool
	currentRate int // requests per window
	windowStart time.Time
	windowSize  time.Duration
}

// NewRateLimiter creates a rate limiter with default settings
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		maxRetries:  5,
		baseDelay:   1 * time.Second,
		maxDelay:    60 * time.Second,
		jitter:      true,
		windowSize:  time.Minute,
	}
}

// Wait blocks until the rate limit allows a new request
func (rl *RateLimiter) Wait() {
	now := time.Now()
	
	// Reset window if expired
	if now.Sub(rl.windowStart) > rl.windowSize {
		rl.currentRate = 0
		rl.windowStart = now
	}
	
	// If under limit, allow immediately
	if rl.currentRate < 100 { // GitHub free tier: 5000/hr, ~83/min
		rl.currentRate++
		return
	}
	
	// Calculate backoff delay
	rl.applyBackoff()
}

// applyBackoff implements exponential backoff with jitter
func (rl *RateLimiter) applyBackoff() {
	// Check for retry-after header simulation (GitHub returns 403 with Retry-After)
	// In real usage, this would parse the Retry-After header
	// For now, use exponential backoff
	
	for i := 0; i < rl.maxRetries; i++ {
		// Exponential delay: baseDelay * 2^i
		expDelay := time.Duration(float64(rl.baseDelay) * math.Pow(2, float64(i)))
		
		// Cap at maxDelay
		if expDelay > rl.maxDelay {
			expDelay = rl.maxDelay
		}
		
		// Add jitter (±25%)
		if rl.jitter {
			jitter := time.Duration(float64(expDelay) * 0.25 * (rand.Float64()*2 - 1))
			expDelay += jitter
		}
		
		time.Sleep(expDelay)
		
		// Reset rate counter after delay
		rl.currentRate = 0
		rl.windowStart = time.Now()
		
		// Check if we can proceed
		if rl.currentRate < 100 {
			rl.currentRate++
			return
		}
	}
	
	// If all retries failed, use maxDelay
	time.Sleep(rl.maxDelay)
}

// Reset resets the rate limiter state
func (rl *RateLimiter) Reset() {
	rl.currentRate = 0
	rl.windowStart = time.Now()
}

// WithJitter enables or disables jitter
func (rl *RateLimiter) WithJitter(enabled bool) *RateLimiter {
	rl.jitter = enabled
	return rl
}

// WithMaxRetries sets the maximum number of retries
func (rl *RateLimiter) WithMaxRetries(retries int) *RateLimiter {
	rl.maxRetries = retries
	return rl
}
