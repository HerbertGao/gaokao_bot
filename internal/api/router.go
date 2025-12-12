package api

import (
	"github.com/gin-gonic/gin"
	"github.com/herbertgao/gaokao_bot/internal/handler"
	"github.com/herbertgao/gaokao_bot/internal/middleware"
	"github.com/herbertgao/gaokao_bot/internal/service"
	"gorm.io/gorm"
)

// NewRouter 创建路由器
// 返回路由器和 RateLimiter 实例（用于生命周期管理）
func NewRouter(
	db *gorm.DB,
	botToken string,
	templateService *service.UserTemplateService,
	skipValidation bool,
	enableLogger bool,
) (*gin.Engine, *middleware.RateLimiter) {
	// 根据是否启用日志来创建路由器
	var router *gin.Engine
	if enableLogger {
		// Debug 模式：使用 Default (包含 Logger 和 Recovery)
		router = gin.Default()
	} else {
		// 非 Debug 模式：只使用 Recovery，不使用 Logger
		router = gin.New()
		router.Use(gin.Recovery())
	}

	// 添加 CORS 中间件
	router.Use(middleware.CORSMiddleware())

	// 创建处理器
	templateHandler := handler.NewTemplateHandler(templateService)

	// 创建速率限制中间件
	rateLimitHandler, rateLimiter := middleware.RateLimitMiddleware(10, 20) // 每秒10个请求，突发20个

	// API 路由组
	api := router.Group("/api")
	{
		// 模板相关 API（需要认证和速率限制）
		templates := api.Group("/templates")
		templates.Use(middleware.TelegramAuthMiddleware(botToken, skipValidation))
		templates.Use(rateLimitHandler)
		{
			templates.GET("", templateHandler.GetTemplates)
			templates.POST("", templateHandler.CreateTemplate)
			templates.PUT("/:id", templateHandler.UpdateTemplate)
			templates.DELETE("/:id", templateHandler.DeleteTemplate)
		}
	}

	// 健康检查（包含数据库连接检查）
	router.GET("/health", func(c *gin.Context) {
		// 检查数据库连接
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(503, gin.H{
				"status": "unhealthy",
				"error":  "无法获取数据库连接",
			})
			return
		}

		// Ping 数据库
		if err := sqlDB.Ping(); err != nil {
			c.JSON(503, gin.H{
				"status": "unhealthy",
				"error":  "数据库连接失败",
			})
			return
		}

		c.JSON(200, gin.H{
			"status":   "ok",
			"database": "connected",
		})
	})

	return router, rateLimiter
}
