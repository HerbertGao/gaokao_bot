package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// InitDataMaxAge Telegram initData 最大有效期（秒）
	InitDataMaxAge = 86400 // 24小时
	// DefaultTestUserID 默认测试用户 ID（用于开发环境）
	DefaultTestUserID = 123456789
)

// getTestUserID 从环境变量获取测试用户 ID，如果未设置则使用默认值
func getTestUserID() int64 {
	if envUserID := os.Getenv("TEST_USER_ID"); envUserID != "" {
		if userID, err := strconv.ParseInt(envUserID, 10, 64); err == nil {
			return userID
		}
	}
	return DefaultTestUserID
}

// TelegramAuthMiddleware Telegram Mini App 认证中间件
func TelegramAuthMiddleware(botToken string, skipValidation bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 initData
		initData := c.GetHeader("X-Telegram-Init-Data")

		// 开发模式：跳过验证
		if skipValidation {
			if initData == "" {
				// 没有 Init Data，使用测试用户 ID（从环境变量或默认值）
				c.Set("user_id", getTestUserID())
			} else {
				// 尝试解析用户 ID（不验证签名）
				values, err := url.ParseQuery(initData)
				if err == nil {
					userStr := values.Get("user")
					var userData struct {
						ID int64 `json:"id"`
					}
					if err := json.Unmarshal([]byte(userStr), &userData); err == nil && userData.ID > 0 {
						c.Set("user_id", userData.ID)
					} else {
						c.Set("user_id", getTestUserID())
					}
				} else {
					c.Set("user_id", getTestUserID())
				}
			}
			c.Next()
			return
		}

		// 生产模式：严格验证
		if initData == "" {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "缺少 Telegram 初始化数据",
			})
			c.Abort()
			return
		}

		// 验证 initData
		userID, err := ValidateTelegramInitData(initData, botToken)
		if err != nil {
			// 不暴露具体的验证失败原因，防止信息泄露
			c.JSON(401, gin.H{
				"success": false,
				"error":   "身份验证失败",
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
		return 0, fmt.Errorf("解析初始化数据失败: %w", err)
	}

	// 获取 hash
	hash := values.Get("hash")
	if hash == "" {
		return 0, fmt.Errorf("缺少签名哈希")
	}

	// 获取 auth_date
	authDateStr := values.Get("auth_date")
	if authDateStr == "" {
		return 0, fmt.Errorf("缺少认证时间")
	}

	// 验证时间戳（不超过 24 小时）
	var authDate int64
	if _, err := fmt.Sscanf(authDateStr, "%d", &authDate); err != nil {
		return 0, fmt.Errorf("认证时间格式无效: %w", err)
	}

	now := time.Now().Unix()
	if now-authDate > InitDataMaxAge {
		return 0, fmt.Errorf("初始化数据已过期")
	}

	// 提取用户 ID
	userStr := values.Get("user")
	if userStr == "" {
		return 0, fmt.Errorf("缺少用户数据")
	}

	// 从 user JSON 中提取 id
	// user 格式: {"id":123456789,"first_name":"Test",...}
	var userData struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal([]byte(userStr), &userData); err != nil {
		return 0, fmt.Errorf("用户数据格式无效: %w", err)
	}
	if userData.ID == 0 {
		return 0, fmt.Errorf("用户ID无效")
	}

	// 构建数据字符串进行验证
	delete(values, "hash")
	var dataCheckArr []string
	for key := range values {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", key, values.Get(key)))
	}
	sort.Strings(dataCheckArr)
	dataCheckString := strings.Join(dataCheckArr, "\n")

	// 计算 secret_key（根据 Telegram 规范: HMAC_SHA256("WebAppData", botToken)）
	secretKeyHMAC := hmac.New(sha256.New, []byte("WebAppData"))
	secretKeyHMAC.Write([]byte(botToken))
	secretKey := secretKeyHMAC.Sum(nil)

	// 计算 HMAC-SHA256
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	// 验证 hash（使用常量时间比较防止定时攻击）
	if !hmac.Equal([]byte(calculatedHash), []byte(hash)) {
		return 0, fmt.Errorf("签名验证失败")
	}

	return userData.ID, nil
}
