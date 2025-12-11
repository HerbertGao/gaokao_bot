package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // 每秒允许的请求数
	burst    int           // 突发请求数
	cleanup  time.Duration // 清理过期访客的间隔
	done     chan struct{} // 关闭信号
}

type visitor struct {
	lastSeen time.Time
	tokens   float64
	mu       sync.Mutex
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
		cleanup:  5 * time.Minute,
		done:     make(chan struct{}),
	}

	// 启动清理 goroutine
	go rl.cleanupVisitors()

	return rl
}

// cleanupVisitors 清理过期的访客
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			for key, v := range rl.visitors {
				if time.Since(v.lastSeen) > rl.cleanup {
					delete(rl.visitors, key)
				}
			}
			rl.mu.Unlock()
		case <-rl.done:
			return
		}
	}
}

// Stop 停止速率限制器
func (rl *RateLimiter) Stop() {
	close(rl.done)
}

// getVisitor 获取或创建访客
func (rl *RateLimiter) getVisitor(key string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[key]
	if !exists {
		v = &visitor{
			lastSeen: time.Now(),
			tokens:   float64(rl.burst),
		}
		rl.visitors[key] = v
	}

	return v
}

// allow 检查是否允许请求（令牌桶算法）
func (rl *RateLimiter) allow(key string) bool {
	v := rl.getVisitor(key)
	v.mu.Lock()
	defer v.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(v.lastSeen).Seconds()
	v.lastSeen = now

	// 添加新令牌
	v.tokens += elapsed * float64(rl.rate)
	if v.tokens > float64(rl.burst) {
		v.tokens = float64(rl.burst)
	}

	// 检查是否有足够的令牌
	if v.tokens >= 1 {
		v.tokens--
		return true
	}

	return false
}

// RateLimitMiddleware 创建速率限制中间件
// 返回中间件处理函数和 RateLimiter 实例（用于生命周期管理）
func RateLimitMiddleware(rate, burst int) (gin.HandlerFunc, *RateLimiter) {
	limiter := NewRateLimiter(rate, burst)

	handler := func(c *gin.Context) {
		// 使用 user_id 作为限制键，如果没有则使用 IP
		var key string
		userID, exists := c.Get("user_id")
		if exists && userID != nil {
			key = fmt.Sprintf("user_%v", userID)
		} else {
			key = c.ClientIP()
		}

		if !limiter.allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}

	return handler, limiter
}
