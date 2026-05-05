package ratelimit

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Config struct {
	RequestsPerSecond int  `envconfig:"RATE_LIMIT_RPS" default:"10"`       // RequestsPerSecond - количество запросов в секунду
	Burst             int  `envconfig:"RATE_LIMIT_BURST" default:"20"`     // Burst - максимальное количество запросов за короткий промежуток времени
	Enabled           bool `envconfig:"RATE_LIMIT_ENABLED" default:"true"` // Enabled - включен ли rate limiting
}

type Limiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	rps      rate.Limit
	burst    int
}

func NewLimiter(rps int, burst int) *Limiter {
	return &Limiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rate.Limit(rps),
		burst:    burst,
	}
}

func (l *Limiter) getLimiter(key string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, exists := l.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(l.rps, l.burst)
		l.limiters[key] = limiter
	}

	return limiter
}

func (l *Limiter) allow(key string) bool {
	limiter := l.getLimiter(key)
	return limiter.Allow()
}

// Middleware возвращает gin middleware для rate limiting по IP
func (l *Limiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем IP клиента
		ip := c.ClientIP()

		if !l.allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "too many requests",
				"message":     "rate limit exceeded, please try again later",
				"retry_after": time.Second / time.Duration(l.rps),
			})
			return
		}

		c.Next()
	}
}

// RateLimiterMiddleware создает и возвращает middleware с конфигурацией
func RateLimiterMiddleware(cfg Config) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	limiter := NewLimiter(cfg.RequestsPerSecond, cfg.Burst)
	return limiter.Middleware()
}
