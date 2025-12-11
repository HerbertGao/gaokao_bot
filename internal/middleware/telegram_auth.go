package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// TelegramAuthMiddleware Telegram Mini App 认证中间件
func TelegramAuthMiddleware(botToken string, skipValidation bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 initData
		initData := c.GetHeader("X-Telegram-Init-Data")

		// 开发模式：跳过验证
		if skipValidation {
			if initData == "" {
				// 没有 Init Data，使用测试用户 ID
				c.Set("user_id", int64(0))
			} else {
				// 尝试解析用户 ID（不验证签名）
				values, err := url.ParseQuery(initData)
				if err == nil {
					userStr := values.Get("user")
					var userID int64
					fmt.Sscanf(userStr, `{"id":%d`, &userID)
					if userID > 0 {
						c.Set("user_id", userID)
					} else {
						c.Set("user_id", int64(0))
					}
				} else {
					c.Set("user_id", int64(0))
				}
			}
			c.Next()
			return
		}

		// 生产模式：严格验证
		if initData == "" {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "Missing Telegram init data",
			})
			c.Abort()
			return
		}

		// 验证 initData
		userID, err := ValidateTelegramInitData(initData, botToken)
		if err != nil {
			c.JSON(401, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Invalid Telegram init data: %v", err),
			})
			c.Abort()
			return
		}

		// 将 userID 存储到上下文中
		c.Set("user_id", userID)
		c.Next()
	}
}

// ValidateTelegramInitData 验证 Telegram initData
// 返回用户 ID 和可能的错误
func ValidateTelegramInitData(initData, botToken string) (int64, error) {
	// 解析 initData
	values, err := url.ParseQuery(initData)
	if err != nil {
		return 0, fmt.Errorf("failed to parse init data: %w", err)
	}

	// 获取 hash
	hash := values.Get("hash")
	if hash == "" {
		return 0, fmt.Errorf("missing hash")
	}

	// 获取 auth_date
	authDateStr := values.Get("auth_date")
	if authDateStr == "" {
		return 0, fmt.Errorf("missing auth_date")
	}

	// 验证时间戳（不超过 24 小时）
	var authDate int64
	if _, err := fmt.Sscanf(authDateStr, "%d", &authDate); err != nil {
		return 0, fmt.Errorf("invalid auth_date: %w", err)
	}

	now := time.Now().Unix()
	if now-authDate > 86400 {
		return 0, fmt.Errorf("init data expired")
	}

	// 提取用户 ID
	userStr := values.Get("user")
	if userStr == "" {
		return 0, fmt.Errorf("missing user data")
	}

	// 从 user JSON 中提取 id（简单方式）
	// user 格式: {"id":123456789,"first_name":"Test",...}
	var userID int64
	if _, err := fmt.Sscanf(userStr, `{"id":%d`, &userID); err != nil {
		return 0, fmt.Errorf("invalid user data: %w", err)
	}

	// 构建数据字符串进行验证
	delete(values, "hash")
	var dataCheckArr []string
	for key := range values {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", key, values.Get(key)))
	}
	sort.Strings(dataCheckArr)
	dataCheckString := strings.Join(dataCheckArr, "\n")

	// 计算 secret_key
	secretKey := sha256.Sum256([]byte(botToken))

	// 计算 HMAC-SHA256
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	// 验证 hash
	if calculatedHash != hash {
		return 0, fmt.Errorf("hash mismatch")
	}

	return userID, nil
}
