package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/repository"
	"github.com/herbertgao/gaokao_bot/internal/service"
	"github.com/herbertgao/gaokao_bot/internal/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 创建测试用的内存数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// 自动迁移表结构
	if err := db.AutoMigrate(&model.UserTemplate{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func setupTestHandler(t *testing.T) (*TemplateHandler, *gorm.DB) {
	// 初始化 Snowflake（如果未初始化）
	_ = util.InitSnowflake(0, 1)

	db := setupTestDB(t)
	repo := repository.NewUserTemplateRepository(db)
	templateService := service.NewUserTemplateService(repo)
	handler := NewTemplateHandler(templateService)
	return handler, db
}

func setupTestRouter(handler *TemplateHandler, userID int64) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 设置用户ID中间件
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})

	router.GET("/templates", handler.GetTemplates)
	router.POST("/templates", handler.CreateTemplate)
	router.PUT("/templates/:id", handler.UpdateTemplate)
	router.DELETE("/templates/:id", handler.DeleteTemplate)

	return router
}

func TestGetTemplates(t *testing.T) {
	handler, db := setupTestHandler(t)
	userID := int64(123)

	// 添加测试数据
	db.Create(&model.UserTemplate{
		ID:              1,
		UserID:          userID,
		TemplateName:    "测试模板1",
		TemplateContent: "距离{exam}还有{time}",
	})
	db.Create(&model.UserTemplate{
		ID:              2,
		UserID:          userID,
		TemplateName:    "测试模板2",
		TemplateContent: "{exam}倒计时{time}",
	})

	router := setupTestRouter(handler, userID)

	req, _ := http.NewRequest(http.MethodGet, "/templates", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if success, ok := response["success"].(bool); !ok || !success {
		t.Error("Expected success to be true")
	}
}

func TestCreateTemplate_Success(t *testing.T) {
	handler, _ := setupTestHandler(t)
	userID := int64(123)

	router := setupTestRouter(handler, userID)

	reqBody := CreateTemplateRequest{
		TemplateName:    "新模板",
		TemplateContent: "距离{exam}还有{time}",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/templates", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestCreateTemplate_ExceedLimit(t *testing.T) {
	handler, db := setupTestHandler(t)
	userID := int64(123)

	// 创建最大数量的模板
	for i := int64(0); i < MaxTemplatesPerUser; i++ {
		db.Create(&model.UserTemplate{
			ID:              1000 + i,
			UserID:          userID,
			TemplateName:    fmt.Sprintf("模板%d", i),
			TemplateContent: "距离{exam}还有{time}",
		})
	}

	router := setupTestRouter(handler, userID)

	reqBody := CreateTemplateRequest{
		TemplateName:    "新模板",
		TemplateContent: "距离{exam}还有{time}",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/templates", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateTemplate_InvalidContent(t *testing.T) {
	handler, _ := setupTestHandler(t)
	userID := int64(123)

	router := setupTestRouter(handler, userID)

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "Empty content",
			content: "",
		},
		{
			name:    "Missing {exam}",
			content: "倒计时{time}",
		},
		{
			name:    "Missing {time}",
			content: "{exam}倒计时",
		},
		{
			name:    "Too long content",
			content: "{exam}{time}" + string(make([]byte, 150)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := CreateTemplateRequest{
				TemplateName:    "测试",
				TemplateContent: tt.content,
			}

			body, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest(http.MethodPost, "/templates", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Status code = %d, want %d for %s", w.Code, http.StatusBadRequest, tt.name)
			}
		})
	}
}

func TestUpdateTemplate_Success(t *testing.T) {
	handler, db := setupTestHandler(t)
	userID := int64(123)
	templateID := int64(1)

	// 添加现有模板
	db.Create(&model.UserTemplate{
		ID:              templateID,
		UserID:          userID,
		TemplateName:    "旧模板",
		TemplateContent: "距离{exam}还有{time}",
	})

	router := setupTestRouter(handler, userID)

	reqBody := UpdateTemplateRequest{
		TemplateName:    "新模板",
		TemplateContent: "{exam}倒计时{time}",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/templates/%d", templateID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestUpdateTemplate_NotFound(t *testing.T) {
	handler, _ := setupTestHandler(t)
	userID := int64(123)

	router := setupTestRouter(handler, userID)

	reqBody := UpdateTemplateRequest{
		TemplateName:    "新模板",
		TemplateContent: "{exam}倒计时{time}",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/templates/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestUpdateTemplate_Forbidden(t *testing.T) {
	handler, db := setupTestHandler(t)
	userID := int64(123)
	otherUserID := int64(456)
	templateID := int64(1)

	// 添加属于其他用户的模板
	db.Create(&model.UserTemplate{
		ID:              templateID,
		UserID:          otherUserID,
		TemplateName:    "其他用户模板",
		TemplateContent: "距离{exam}还有{time}",
	})

	router := setupTestRouter(handler, userID)

	reqBody := UpdateTemplateRequest{
		TemplateName:    "新模板",
		TemplateContent: "{exam}倒计时{time}",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/templates/%d", templateID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestDeleteTemplate_Success(t *testing.T) {
	handler, db := setupTestHandler(t)
	userID := int64(123)
	templateID := int64(1)

	// 添加模板
	db.Create(&model.UserTemplate{
		ID:              templateID,
		UserID:          userID,
		TemplateName:    "测试模板",
		TemplateContent: "距离{exam}还有{time}",
	})

	router := setupTestRouter(handler, userID)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/templates/%d", templateID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	// 验证模板已删除
	var count int64
	db.Model(&model.UserTemplate{}).Where("id = ?", templateID).Count(&count)
	if count != 0 {
		t.Error("Template should be deleted")
	}
}

func TestDeleteTemplate_NotFound(t *testing.T) {
	handler, _ := setupTestHandler(t)
	userID := int64(123)

	router := setupTestRouter(handler, userID)

	req, _ := http.NewRequest(http.MethodDelete, "/templates/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestValidateTemplateContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "Valid content",
			content: "距离{exam}还有{time}",
			wantErr: false,
		},
		{
			name:    "Empty content",
			content: "",
			wantErr: true,
		},
		{
			name:    "Missing {exam}",
			content: "倒计时{time}",
			wantErr: true,
		},
		{
			name:    "Missing {time}",
			content: "{exam}倒计时",
			wantErr: true,
		},
		{
			name:    "Too long content (Chinese chars)",
			content: "{exam}{time}" + string([]rune("这是一个非常长的模板内容用于测试字符数限制功能这个模板包含了很多中文字符每个中文字符在UTF8编码中占用三个字节所以我们需要确保使用正确的字符计数方法而不是字节计数方法这样才能正确验证模板内容的长度限制功能是否正常工作现在继续添加更多的中文字符用来确保真的超过一百四十个字符的限制")),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTemplateContent(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTemplateContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTemplateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid name",
			input:   "我的模板",
			wantErr: false,
		},
		{
			name:    "Empty name (allowed)",
			input:   "",
			wantErr: false,
		},
		{
			name:    "Too long name",
			input:   "这是一个超级超级长的模板名称用于测试验证字符长度限制",
			wantErr: true,
		},
		{
			name:    "Exactly at limit",
			input:   "12345678901234567890", // 20 chars
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTemplateName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTemplateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
