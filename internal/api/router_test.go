package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herbertgao/gaokao_bot/internal/repository"
	"github.com/herbertgao/gaokao_bot/internal/service"
	"github.com/herbertgao/gaokao_bot/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&model.UserTemplate{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func TestNewRouter(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserTemplateRepository(db)
	templateService := service.NewUserTemplateService(repo)
	botToken := "test_token"

	router, rateLimiter := NewRouter(db, botToken, templateService, true, false)
	defer rateLimiter.Stop()

	if router == nil {
		t.Error("Expected router, got nil")
	}

	if rateLimiter == nil {
		t.Error("Expected rateLimiter, got nil")
	}
}

func TestHealthCheck(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserTemplateRepository(db)
	templateService := service.NewUserTemplateService(repo)
	botToken := "test_token"

	router, rateLimiter := NewRouter(db, botToken, templateService, true, false)
	defer rateLimiter.Stop()

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Health check status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestHealthCheck_ResponseFormat(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserTemplateRepository(db)
	templateService := service.NewUserTemplateService(repo)
	botToken := "test_token"

	router, rateLimiter := NewRouter(db, botToken, templateService, true, false)
	defer rateLimiter.Stop()

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	// 验证响应包含必要的字段
	body := w.Body.String()
	if body == "" {
		t.Error("Health check response should not be empty")
	}
}

func TestRouter_WithLogger(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserTemplateRepository(db)
	templateService := service.NewUserTemplateService(repo)
	botToken := "test_token"

	// 测试启用日志
	router, rateLimiter := NewRouter(db, botToken, templateService, true, true)
	defer rateLimiter.Stop()

	if router == nil {
		t.Error("Expected router with logger, got nil")
	}
}

func TestRouter_WithoutLogger(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserTemplateRepository(db)
	templateService := service.NewUserTemplateService(repo)
	botToken := "test_token"

	// 测试禁用日志
	router, rateLimiter := NewRouter(db, botToken, templateService, true, false)
	defer rateLimiter.Stop()

	if router == nil {
		t.Error("Expected router without logger, got nil")
	}
}
