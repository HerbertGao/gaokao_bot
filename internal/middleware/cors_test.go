package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// testAllowedOrigins 测试用的允许源列表
var testAllowedOrigins = []string{
	"https://web.telegram.org",
	"http://localhost:5173",
	"http://localhost:3000",
	"http://127.0.0.1:5173",
	"http://127.0.0.1:3000",
}

func TestIsTelegramOrigin(t *testing.T) {
	tests := []struct {
		name   string
		origin string
		want   bool
	}{
		{
			name:   "Valid - telegram.org",
			origin: "https://telegram.org",
			want:   true,
		},
		{
			name:   "Valid - web.telegram.org",
			origin: "https://web.telegram.org",
			want:   true,
		},
		{
			name:   "Valid - subdomain.telegram.org",
			origin: "https://app.mini.telegram.org",
			want:   true,
		},
		{
			name:   "Valid with port - telegram.org:443",
			origin: "https://telegram.org:443",
			want:   true,
		},
		{
			name:   "Valid with port - web.telegram.org:8443",
			origin: "https://web.telegram.org:8443",
			want:   true,
		},
		{
			name:   "Invalid - nottelegram.org",
			origin: "https://nottelegram.org",
			want:   false,
		},
		{
			name:   "Invalid - telegram.org.evil.com",
			origin: "https://telegram.org.evil.com",
			want:   false,
		},
		{
			name:   "Invalid with port - nottelegram.org:8443",
			origin: "https://nottelegram.org:8443",
			want:   false,
		},
		{
			name:   "Invalid - faketelegram.org",
			origin: "https://faketelegram.org",
			want:   false,
		},
		{
			name:   "Invalid URL",
			origin: "not-a-url",
			want:   false,
		},
		{
			name:   "Empty string",
			origin: "",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTelegramOrigin(tt.origin)
			if got != tt.want {
				t.Errorf("isTelegramOrigin(%q) = %v, want %v", tt.origin, got, tt.want)
			}
		})
	}
}

func TestIsOriginAllowed_WithTelegramInList(t *testing.T) {
	tests := []struct {
		name   string
		origin string
		want   bool
	}{
		{
			name:   "Allowed origin - telegram.org",
			origin: "https://web.telegram.org",
			want:   true,
		},
		{
			name:   "Allowed origin - localhost 5173",
			origin: "http://localhost:5173",
			want:   true,
		},
		{
			name:   "Allowed origin - localhost 3000",
			origin: "http://localhost:3000",
			want:   true,
		},
		{
			name:   "Allowed origin - 127.0.0.1:5173",
			origin: "http://127.0.0.1:5173",
			want:   true,
		},
		{
			name:   "Allowed origin - 127.0.0.1:3000",
			origin: "http://127.0.0.1:3000",
			want:   true,
		},
		{
			name:   "Telegram Mini App dynamic domain",
			origin: "https://myapp.telegram.org",
			want:   true,
		},
		{
			name:   "Telegram Mini App subdomain",
			origin: "https://app.mini.telegram.org",
			want:   true,
		},
		{
			name:   "Telegram origin with port - should be allowed",
			origin: "https://app.telegram.org:8443",
			want:   true,
		},
		{
			name:   "Telegram origin with port 443 - should be allowed",
			origin: "https://web.telegram.org:443",
			want:   true,
		},
		{
			name:   "Telegram origin with custom port - should be allowed",
			origin: "https://test.telegram.org:9000",
			want:   true,
		},
		{
			name:   "Not allowed origin",
			origin: "https://evil.com",
			want:   false,
		},
		{
			name:   "Not allowed - similar but not telegram.org",
			origin: "https://web.telegram.org.fake.com",
			want:   false,
		},
		{
			name:   "Domain spoofing - telegram.org.evil.com",
			origin: "https://telegram.org.evil.com",
			want:   false,
		},
		{
			name:   "Domain spoofing - subdomain.telegram.org.evil.com",
			origin: "https://app.telegram.org.evil.com",
			want:   false,
		},
		{
			name:   "Domain spoofing with port - telegram.org.evil.com:8443",
			origin: "https://telegram.org.evil.com:8443",
			want:   false,
		},
		{
			name:   "Empty origin",
			origin: "",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isOriginAllowed(tt.origin, testAllowedOrigins)
			if got != tt.want {
				t.Errorf("isOriginAllowed(%q) = %v, want %v", tt.origin, got, tt.want)
			}
		})
	}
}

func TestIsOriginAllowed_WithoutTelegramInList(t *testing.T) {
	// 测试当 allowedOrigins 不包含 telegram.org 时，
	// telegram.org 子域名不应被允许（安全修复验证）
	noTelegramOrigins := []string{
		"http://localhost:5173",
		"http://localhost:3000",
		"https://my-custom-domain.com",
	}

	tests := []struct {
		name   string
		origin string
		want   bool
	}{
		{
			name:   "Telegram subdomain should be rejected when telegram.org not in allowed list",
			origin: "https://myapp.telegram.org",
			want:   false,
		},
		{
			name:   "web.telegram.org should be rejected when not in allowed list",
			origin: "https://web.telegram.org",
			want:   false,
		},
		{
			name:   "Allowed origin - localhost should work",
			origin: "http://localhost:5173",
			want:   true,
		},
		{
			name:   "Allowed origin - custom domain should work",
			origin: "https://my-custom-domain.com",
			want:   true,
		},
		{
			name:   "Disallowed origin",
			origin: "https://evil.com",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isOriginAllowed(tt.origin, noTelegramOrigins)
			if got != tt.want {
				t.Errorf("isOriginAllowed(%q, %v) = %v, want %v", tt.origin, noTelegramOrigins, got, tt.want)
			}
		})
	}
}

func TestIsOriginAllowed_SubstringMatchAttack(t *testing.T) {
	// 测试子串匹配漏洞修复：
	// 包含 "telegram.org" 子串但不是真正的 Telegram 域名的恶意域名
	// 不应触发通配符匹配
	fakeOrigins := []string{
		"https://nottelegram.org", // 包含 "telegram.org" 子串
		"https://faketelegram.org.example.com", // 包含 "telegram.org" 子串
		"https://my-telegram.org-service.com", // 包含 "telegram.org" 子串
	}

	tests := []struct {
		name   string
		origin string
		want   bool
	}{
		{
			name:   "Real telegram.org subdomain should be rejected (wildcard not triggered)",
			origin: "https://app.telegram.org",
			want:   false,
		},
		{
			name:   "Real web.telegram.org should be rejected (wildcard not triggered)",
			origin: "https://web.telegram.org",
			want:   false,
		},
		{
			name:   "Exact match - nottelegram.org should work",
			origin: "https://nottelegram.org",
			want:   true,
		},
		{
			name:   "Exact match - faketelegram.org.example.com should work",
			origin: "https://faketelegram.org.example.com",
			want:   true,
		},
		{
			name:   "Non-matching origin should be rejected",
			origin: "https://evil.com",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isOriginAllowed(tt.origin, fakeOrigins)
			if got != tt.want {
				t.Errorf("isOriginAllowed(%q, %v) = %v, want %v", tt.origin, fakeOrigins, got, tt.want)
			}
		})
	}
}

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		origin             string
		method             string
		expectAllowOrigin  bool
		expectStatus       int
		expectAbort        bool
	}{
		{
			name:               "Allowed origin with GET",
			origin:             "https://web.telegram.org",
			method:             http.MethodGet,
			expectAllowOrigin:  true,
			expectStatus:       http.StatusOK,
			expectAbort:        false,
		},
		{
			name:               "Allowed origin with POST",
			origin:             "http://localhost:5173",
			method:             http.MethodPost,
			expectAllowOrigin:  true,
			expectStatus:       http.StatusOK,
			expectAbort:        false,
		},
		{
			name:               "OPTIONS preflight request",
			origin:             "https://web.telegram.org",
			method:             http.MethodOptions,
			expectAllowOrigin:  true,
			expectStatus:       http.StatusNoContent,
			expectAbort:        true,
		},
		{
			name:               "Not allowed origin",
			origin:             "https://evil.com",
			method:             http.MethodGet,
			expectAllowOrigin:  false,
			expectStatus:       http.StatusOK,
			expectAbort:        false,
		},
		{
			name:               "No origin header",
			origin:             "",
			method:             http.MethodGet,
			expectAllowOrigin:  false,
			expectStatus:       http.StatusOK,
			expectAbort:        false,
		},
		{
			name:               "Telegram subdomain with OPTIONS",
			origin:             "https://app.telegram.org",
			method:             http.MethodOptions,
			expectAllowOrigin:  true,
			expectStatus:       http.StatusNoContent,
			expectAbort:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试路由
			router := gin.New()
			router.Use(CORSMiddleware(testAllowedOrigins))
			router.GET("/test", func(c *gin.Context) {
				c.String(http.StatusOK, "OK")
			})
			router.POST("/test", func(c *gin.Context) {
				c.String(http.StatusOK, "OK")
			})
			router.OPTIONS("/test", func(c *gin.Context) {
				// OPTIONS 应该被 CORS 中间件拦截，不会到这里
				c.String(http.StatusOK, "OK")
			})

			// 创建请求
			req, _ := http.NewRequest(tt.method, "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			// 执行请求
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证状态码
			if w.Code != tt.expectStatus {
				t.Errorf("Status code = %d, want %d", w.Code, tt.expectStatus)
			}

			// 验证 CORS 头
			allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
			if tt.expectAllowOrigin {
				if allowOrigin != tt.origin {
					t.Errorf("Access-Control-Allow-Origin = %q, want %q", allowOrigin, tt.origin)
				}

				// 验证其他 CORS 头
				if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
					t.Error("Access-Control-Allow-Credentials should be 'true'")
				}

				allowHeaders := w.Header().Get("Access-Control-Allow-Headers")
				if allowHeaders == "" {
					t.Error("Access-Control-Allow-Headers should not be empty")
				}

				allowMethods := w.Header().Get("Access-Control-Allow-Methods")
				if allowMethods == "" {
					t.Error("Access-Control-Allow-Methods should not be empty")
				}
			} else {
				if allowOrigin != "" {
					t.Errorf("Access-Control-Allow-Origin should be empty for disallowed origin, got %q", allowOrigin)
				}
			}
		})
	}
}
