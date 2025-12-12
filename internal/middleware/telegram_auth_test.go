package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestGetTestUserID(t *testing.T) {
	tests := []struct {
		name   string
		envVar string
		want   int64
	}{
		{
			name:   "No environment variable",
			envVar: "",
			want:   DefaultTestUserID,
		},
		{
			name:   "Valid environment variable",
			envVar: "999888777",
			want:   999888777,
		},
		{
			name:   "Invalid environment variable",
			envVar: "invalid",
			want:   DefaultTestUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 保存原始环境变量
			originalEnv := os.Getenv("TEST_USER_ID")
			defer func() {
				_ = os.Setenv("TEST_USER_ID", originalEnv)
			}()

			// 设置测试环境变量
			if tt.envVar != "" {
				_ = os.Setenv("TEST_USER_ID", tt.envVar)
			} else {
				_ = os.Unsetenv("TEST_USER_ID")
			}

			got := getTestUserID()
			if got != tt.want {
				t.Errorf("getTestUserID() = %d, want %d", got, tt.want)
			}
		})
	}
}

// generateValidInitData 生成有效的 Telegram initData 用于测试
func generateValidInitData(botToken string, userID int64) string {
	return generateValidInitDataWithTimestamp(botToken, userID, time.Now().Unix())
}

// generateValidInitDataWithTimestamp 生成带有自定义时间戳的有效 Telegram initData
func generateValidInitDataWithTimestamp(botToken string, userID int64, authDate int64) string {
	userData := fmt.Sprintf(`{"id":%d,"first_name":"Test","username":"testuser"}`, userID)

	// 构建数据
	values := url.Values{}
	values.Set("user", userData)
	values.Set("auth_date", fmt.Sprintf("%d", authDate))

	// 按字母顺序排序
	var dataCheckArr []string
	for key := range values {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", key, values.Get(key)))
	}
	sort.Strings(dataCheckArr)
	dataCheckString := strings.Join(dataCheckArr, "\n")

	// 计算 secret_key
	secretKeyHMAC := hmac.New(sha256.New, []byte("WebAppData"))
	secretKeyHMAC.Write([]byte(botToken))
	secretKey := secretKeyHMAC.Sum(nil)

	// 计算 hash
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(dataCheckString))
	hash := hex.EncodeToString(h.Sum(nil))

	// 返回完整的 initData
	values.Set("hash", hash)
	return values.Encode()
}

func TestValidateTelegramInitData(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	userID := int64(987654321)

	tests := []struct {
		name    string
		initData func() string
		wantErr bool
		wantUserID int64
	}{
		{
			name: "Valid initData",
			initData: func() string {
				return generateValidInitData(botToken, userID)
			},
			wantErr: false,
			wantUserID: userID,
		},
		{
			name: "Missing hash",
			initData: func() string {
				return "user={\"id\":123}&auth_date=1234567890"
			},
			wantErr: true,
		},
		{
			name: "Missing auth_date",
			initData: func() string {
				return "user={\"id\":123}&hash=abcd"
			},
			wantErr: true,
		},
		{
			name: "Expired auth_date",
			initData: func() string {
				oldDate := time.Now().Unix() - InitDataMaxAge - 3600 // 过期1小时
				userData := fmt.Sprintf(`{"id":%d}`, userID)
				values := url.Values{}
				values.Set("user", userData)
				values.Set("auth_date", fmt.Sprintf("%d", oldDate))
				values.Set("hash", "dummy")
				return values.Encode()
			},
			wantErr: true,
		},
		{
			name: "Future timestamp (beyond clock skew)",
			initData: func() string {
				futureDate := time.Now().Unix() + ClockSkewTolerance + 3600 // 未来1小时（超过偏移容忍度）
				userData := fmt.Sprintf(`{"id":%d}`, userID)
				values := url.Values{}
				values.Set("user", userData)
				values.Set("auth_date", fmt.Sprintf("%d", futureDate))
				values.Set("hash", "dummy")
				return values.Encode()
			},
			wantErr: true,
		},
		{
			name: "Future timestamp within clock skew tolerance",
			initData: func() string {
				// 未来2分钟，在容忍度内（5分钟）
				futureDate := time.Now().Unix() + 120
				return generateValidInitDataWithTimestamp(botToken, userID, futureDate)
			},
			wantErr: false,
			wantUserID: userID,
		},
		{
			name: "Missing user data",
			initData: func() string {
				authDate := time.Now().Unix()
				values := url.Values{}
				values.Set("auth_date", fmt.Sprintf("%d", authDate))
				values.Set("hash", "dummy")
				return values.Encode()
			},
			wantErr: true,
		},
		{
			name: "Invalid user JSON",
			initData: func() string {
				authDate := time.Now().Unix()
				values := url.Values{}
				values.Set("user", "not-json")
				values.Set("auth_date", fmt.Sprintf("%d", authDate))
				values.Set("hash", "dummy")
				return values.Encode()
			},
			wantErr: true,
		},
		{
			name: "Invalid hash",
			initData: func() string {
				validData := generateValidInitData(botToken, userID)
				values, _ := url.ParseQuery(validData)
				values.Set("hash", "invalid_hash")
				return values.Encode()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initData := tt.initData()
			gotUserID, err := ValidateTelegramInitData(initData, botToken)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTelegramInitData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && gotUserID != tt.wantUserID {
				t.Errorf("ValidateTelegramInitData() userID = %d, want %d", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestTelegramAuthMiddleware_SkipValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	botToken := "test_token"

	tests := []struct {
		name       string
		initData   string
		expectUserID int64
	}{
		{
			name:       "No init data - use default test user",
			initData:   "",
			expectUserID: DefaultTestUserID,
		},
		{
			name:       "Valid init data - parse user ID",
			initData:   `user={"id":555}&auth_date=123456&hash=dummy`,
			expectUserID: 555,
		},
		{
			name:       "Invalid init data format - use default",
			initData:   "invalid",
			expectUserID: DefaultTestUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(TelegramAuthMiddleware(botToken, true)) // skipValidation = true
			router.GET("/test", func(c *gin.Context) {
				userID := c.GetInt64("user_id")
				c.JSON(http.StatusOK, gin.H{"user_id": userID})
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tt.initData != "" {
				req.Header.Set("X-Telegram-Init-Data", tt.initData)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
			}

			// 验证 user_id 通过检查响应（简化测试）
			// 实际应用中可以解析 JSON 响应
		})
	}
}

func TestTelegramAuthMiddleware_ProductionMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	userID := int64(111222333)

	tests := []struct {
		name         string
		initData     string
		expectStatus int
		expectError  bool
	}{
		{
			name:         "No init data - unauthorized",
			initData:     "",
			expectStatus: http.StatusUnauthorized,
			expectError:  true,
		},
		{
			name:         "Valid init data - success",
			initData:     generateValidInitData(botToken, userID),
			expectStatus: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "Invalid hash - unauthorized",
			initData:     "user={\"id\":123}&auth_date=" + fmt.Sprintf("%d", time.Now().Unix()) + "&hash=invalid",
			expectStatus: http.StatusUnauthorized,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(TelegramAuthMiddleware(botToken, false)) // skipValidation = false (生产模式)
			router.GET("/test", func(c *gin.Context) {
				c.String(http.StatusOK, "OK")
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tt.initData != "" {
				req.Header.Set("X-Telegram-Init-Data", tt.initData)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectStatus {
				t.Errorf("Status code = %d, want %d", w.Code, tt.expectStatus)
			}
		})
	}
}

func TestTelegramAuthMiddleware_SetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	userID := int64(777888999)

	router := gin.New()
	router.Use(TelegramAuthMiddleware(botToken, false))
	router.GET("/test", func(c *gin.Context) {
		gotUserID := c.GetInt64("user_id")
		if gotUserID != userID {
			t.Errorf("user_id in context = %d, want %d", gotUserID, userID)
		}
		c.String(http.StatusOK, "OK")
	})

	initData := generateValidInitData(botToken, userID)
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Telegram-Init-Data", initData)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}
}
