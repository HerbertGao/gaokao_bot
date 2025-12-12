package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestNewRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(10, 20)

	if limiter == nil {
		t.Fatal("NewRateLimiter returned nil")
	}

	if limiter.rate != 10 {
		t.Errorf("rate = %d, want 10", limiter.rate)
	}

	if limiter.burst != 20 {
		t.Errorf("burst = %d, want 20", limiter.burst)
	}

	if limiter.visitors == nil {
		t.Error("visitors map should be initialized")
	}

	if limiter.done == nil {
		t.Error("done channel should be initialized")
	}

	// 清理
	limiter.Stop()
}

func TestRateLimiter_Allow(t *testing.T) {
	limiter := NewRateLimiter(10, 5) // 每秒10个请求，突发5个
	defer limiter.Stop()

	key := "test_user"

	// 测试突发请求（应该允许前5个）
	for i := 0; i < 5; i++ {
		if !limiter.allow(key) {
			t.Errorf("Request %d should be allowed (within burst limit)", i+1)
		}
	}

	// 第6个请求应该被拒绝（超过突发限制）
	if limiter.allow(key) {
		t.Error("Request 6 should be denied (exceeds burst limit)")
	}

	// 等待令牌恢复（100ms = 1个令牌）
	time.Sleep(150 * time.Millisecond)

	// 现在应该允许1个请求
	if !limiter.allow(key) {
		t.Error("Request should be allowed after token refill")
	}
}

func TestRateLimiter_MultipleKeys(t *testing.T) {
	limiter := NewRateLimiter(10, 3)
	defer limiter.Stop()

	user1 := "user_1"
	user2 := "user_2"

	// 用户1消耗所有令牌
	for i := 0; i < 3; i++ {
		if !limiter.allow(user1) {
			t.Errorf("User1 request %d should be allowed", i+1)
		}
	}

	// 用户1的第4个请求被拒绝
	if limiter.allow(user1) {
		t.Error("User1 request 4 should be denied")
	}

	// 用户2应该还能发送请求（独立的令牌桶）
	for i := 0; i < 3; i++ {
		if !limiter.allow(user2) {
			t.Errorf("User2 request %d should be allowed", i+1)
		}
	}
}

func TestRateLimiter_Stop(t *testing.T) {
	limiter := NewRateLimiter(10, 20)

	// 停止 limiter
	limiter.Stop()

	// 验证 done channel 已关闭
	select {
	case <-limiter.done:
		// 正确：channel 已关闭
	default:
		t.Error("done channel should be closed after Stop()")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         interface{}
		requestCount   int
		rate           int
		burst          int
		expectDenied   int // 期望被拒绝的请求数
	}{
		{
			name:           "All requests allowed within burst",
			userID:         int64(123),
			requestCount:   3,
			rate:           10,
			burst:          5,
			expectDenied:   0,
		},
		{
			name:           "Some requests denied exceeding burst",
			userID:         int64(456),
			requestCount:   8,
			rate:           10,
			burst:          5,
			expectDenied:   3,
		},
		{
			name:           "No user ID - use IP",
			userID:         nil,
			requestCount:   3,
			rate:           10,
			burst:          2,
			expectDenied:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建路由和中间件
			router := gin.New()
			handler, limiter := RateLimitMiddleware(tt.rate, tt.burst)
			defer limiter.Stop()

			// 设置中间件
			router.Use(func(c *gin.Context) {
				if tt.userID != nil {
					c.Set("user_id", tt.userID)
				}
				c.Next()
			})
			router.Use(handler)
			router.GET("/test", func(c *gin.Context) {
				c.String(http.StatusOK, "OK")
			})

			// 发送多个请求
			deniedCount := 0
			for i := 0; i < tt.requestCount; i++ {
				req, _ := http.NewRequest(http.MethodGet, "/test", nil)
				req.RemoteAddr = "127.0.0.1:12345" // 固定 IP 以便测试
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				if w.Code == http.StatusTooManyRequests {
					deniedCount++
				}
			}

			if deniedCount != tt.expectDenied {
				t.Errorf("Denied requests = %d, want %d", deniedCount, tt.expectDenied)
			}
		})
	}
}

func TestRateLimitMiddleware_ErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler, limiter := RateLimitMiddleware(10, 1) // burst=1，第二个请求就会被拒绝
	defer limiter.Stop()

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(789))
		c.Next()
	})
	router.Use(handler)
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// 第一个请求应该成功
	req1, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("First request status = %d, want %d", w1.Code, http.StatusOK)
	}

	// 第二个请求应该被限流
	req2, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("Second request status = %d, want %d", w2.Code, http.StatusTooManyRequests)
	}

	// 验证错误消息（应该是中文）
	body := w2.Body.String()
	if body == "" {
		t.Error("Rate limit error response should not be empty")
	}
}

func TestRateLimiter_CleanupVisitors(t *testing.T) {
	// 创建一个清理间隔很短的 limiter
	limiter := &RateLimiter{
		visitors:    make(map[string]*visitor),
		rate:        10,
		burst:       20,
		cleanup:     50 * time.Millisecond, // 短的清理间隔
		maxVisitors: 10000,                 // 设置最大访客数量
		done:        make(chan struct{}),
	}

	// 启动清理 goroutine
	go limiter.cleanupVisitors()
	defer limiter.Stop()

	// 添加一个访客
	key := "test_cleanup"
	limiter.allow(key)

	// 验证访客存在
	limiter.mu.RLock()
	_, exists := limiter.visitors[key]
	limiter.mu.RUnlock()

	if !exists {
		t.Error("Visitor should exist after allow()")
	}

	// 等待清理（设置 lastSeen 为过去的时间以触发清理）
	limiter.mu.Lock()
	if v, ok := limiter.visitors[key]; ok {
		v.lastSeen = time.Now().Add(-10 * time.Minute)
	}
	limiter.mu.Unlock()

	// 等待清理运行
	time.Sleep(100 * time.Millisecond)

	// 验证访客被清理
	limiter.mu.RLock()
	_, exists = limiter.visitors[key]
	limiter.mu.RUnlock()

	if exists {
		t.Error("Visitor should be cleaned up after expiration")
	}
}
