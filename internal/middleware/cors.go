package middleware

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// CORSMaxAge CORS 预检请求缓存时间（秒）- 24小时
	CORSMaxAge = "86400"
)

// isTelegramOrigin 检查一个 origin 是否是真正的 Telegram 域名
// 使用 URL 解析确保精确匹配，防止子串攻击
func isTelegramOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	host := u.Hostname() // 获取主机名（不包含端口）

	// 必须是 telegram.org 或其子域名
	return host == "telegram.org" || strings.HasSuffix(host, ".telegram.org")
}

// isOriginAllowed 检查 origin 是否被允许
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	// 检查是否有 telegram.org 域名在允许列表中
	hasTelegramOrigin := false
	for _, allowed := range allowedOrigins {
		if allowed == origin {
			return true
		}
		// 检查是否显式允许了真正的 telegram.org 域名
		// 使用 URL 解析防止 "nottelegram.org" 等欺骗性域名触发通配符
		if isTelegramOrigin(allowed) {
			hasTelegramOrigin = true
		}
	}

	// 仅当 allowedOrigins 中包含 telegram.org 相关域名时，才启用通配符匹配
	// 支持 Telegram Mini App 的动态子域名
	// 注意：这允许任何 *.telegram.org 子域名（例如 myapp.telegram.org）
	// 这是 Telegram Mini Apps 的预期行为，因为它们可以使用各种子域名。
	// strings.HasSuffix() 确保域名以 ".telegram.org" 结尾，
	// 防止类似 "telegram.org.evil.com" 的欺骗攻击。
	//
	// 安全权衡：信任所有 telegram.org 子域名是可接受的，因为：
	// 1. telegram.org 域名由 Telegram 控制
	// 2. Telegram Mini Apps 需要动态子域名支持
	// 3. HasSuffix 检查防止域名欺骗
	// 4. 仅在显式配置 telegram.org 时启用（非全局通配符）
	if hasTelegramOrigin && strings.HasSuffix(origin, ".telegram.org") {
		return true
	}

	return false
}

// CORSMiddleware CORS 中间件
// allowedOrigins: 允许的源列表，从配置文件读取
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 如果没有 Origin 头（同源请求或某些非浏览器客户端），跳过 CORS 头设置
		if origin == "" {
			c.Next()
			return
		}

		// 验证 origin 是否在白名单中
		if isOriginAllowed(origin, allowedOrigins) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Telegram-Init-Data")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Max-Age", CORSMaxAge)
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
