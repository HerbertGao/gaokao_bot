package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/herbertgao/gaokao_bot/internal/util"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	visitors   map[string]*visitor
	mu         sync.RWMutex
	rate       int           // 每秒允许的请求数
	burst      int           // 突发请求数
	cleanup    time.Duration // 清理过期访客的间隔
	maxVisitors int          // 最大访客数量限制，防止内存泄漏
	done       chan struct{} // 关闭信号
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
		// 缩短清理间隔为 2 分钟，更积极地释放内存
		// 平衡内存使用和清理频率：
		// - 太频繁：浪费 CPU
		// - 太少：占用内存
		cleanup: 2 * time.Minute,
		// 最大访客数量限制为 10000
		// 防止恶意请求导致内存无限增长
		maxVisitors: 10000,
		done:        make(chan struct{}),
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
			// 收集需要删除的 visitor 及其 key
			// 保持 visitor 锁定直到删除完成，防止竞态条件
			type expiredVisitor struct {
				key string
				v   *visitor
			}
			var toDelete []expiredVisitor

			rl.mu.Lock()
			for key, v := range rl.visitors {
				v.mu.Lock()
				if time.Since(v.lastSeen) > rl.cleanup {
					toDelete = append(toDelete, expiredVisitor{key: key, v: v})
				} else {
					// 如果不删除，立即解锁
					v.mu.Unlock()
				}
			}

			// 删除过期的访客（仍持有 visitor 锁，防止并发更新）
			for _, item := range toDelete {
				delete(rl.visitors, item.key)
				item.v.mu.Unlock() // 删除后才解锁
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
		// 如果访客数量达到上限，立即触发清理
		if len(rl.visitors) >= rl.maxVisitors {
			rl.cleanupExpiredVisitorsUnsafe()
		}

		// 如果清理后仍然超过上限，拒绝创建新访客（使用现有最严格的限制）
		if len(rl.visitors) >= rl.maxVisitors {
			// 返回一个临时的、tokens 为 0 的访客，这会导致请求被限流
			return &visitor{
				lastSeen: util.NowBJT(),
				tokens:   0,
			}
		}

		v = &visitor{
			lastSeen: util.NowBJT(),
			tokens:   float64(rl.burst),
		}
		rl.visitors[key] = v
	}

	return v
}

// cleanupExpiredVisitorsUnsafe 立即清理过期访客（不加锁，由调用者负责加锁）
func (rl *RateLimiter) cleanupExpiredVisitorsUnsafe() {
	var keysToDelete []string

	for key, v := range rl.visitors {
		v.mu.Lock()
		if time.Since(v.lastSeen) > rl.cleanup {
			keysToDelete = append(keysToDelete, key)
		}
		v.mu.Unlock()
	}

	for _, key := range keysToDelete {
		delete(rl.visitors, key)
	}
}

// allow 检查是否允许请求（令牌桶算法）
func (rl *RateLimiter) allow(key string) bool {
	v := rl.getVisitor(key)
	v.mu.Lock()
	defer v.mu.Unlock()

	now := util.NowBJT()
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
				"error":   "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}

	return handler, limiter
}
