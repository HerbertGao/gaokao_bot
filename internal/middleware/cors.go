package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// isOriginAllowed 检查 origin 是否被允许
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == origin {
			return true
		}
		// 支持 Telegram Mini App 的动态域名
		if strings.HasSuffix(origin, ".telegram.org") {
			return true
		}
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
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
