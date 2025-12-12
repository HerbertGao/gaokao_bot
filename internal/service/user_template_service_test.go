package service

import (
	"testing"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/repository"
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

func setupTestService(t *testing.T) (*UserTemplateService, *gorm.DB) {
	db := setupTestDB(t)
	repo := repository.NewUserTemplateRepository(db)
	service := NewUserTemplateService(repo)
	return service, db
}

func TestUserTemplateService_Create(t *testing.T) {
	service, _ := setupTestService(t)

	template := &model.UserTemplate{
		ID:              1,
		UserID:          123,
		TemplateName:    "测试模板",
		TemplateContent: "距离{exam}还有{time}",
	}

	err := service.Create(template)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}
}

func TestUserTemplateService_GetByID(t *testing.T) {
	service, db := setupTestService(t)

	db.Create(&model.UserTemplate{
		ID:              1,
		UserID:          123,
		TemplateName:    "测试模板",
		TemplateContent: "距离{exam}还有{time}",
	})

	result, err := service.GetByID(1)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}

	if result == nil {
		t.Error("Expected template, got nil")
		return
	}

	if result.TemplateName != "测试模板" {
		t.Errorf("TemplateName = %s, want %s", result.TemplateName, "测试模板")
	}
}

func TestUserTemplateService_GetByUserID(t *testing.T) {
	service, db := setupTestService(t)

	userID := int64(123)

	db.Create(&model.UserTemplate{
		ID:              1,
		UserID:          userID,
		TemplateName:    "模板1",
		TemplateContent: "距离{exam}还有{time}",
	})
	db.Create(&model.UserTemplate{
		ID:              2,
		UserID:          userID,
		TemplateName:    "模板2",
		TemplateContent: "{exam}倒计时{time}",
	})

	result, err := service.GetByUserID(userID)
	if err != nil {
		t.Errorf("GetByUserID() error = %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 templates, got %d", len(result))
	}
}

func TestUserTemplateService_Update(t *testing.T) {
	service, db := setupTestService(t)

	template := &model.UserTemplate{
		ID:              1,
		UserID:          123,
		TemplateName:    "旧模板",
		TemplateContent: "距离{exam}还有{time}",
	}
	db.Create(template)

	template.TemplateName = "新模板"
	err := service.Update(template)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	updated, _ := service.GetByID(1)
	if updated.TemplateName != "新模板" {
		t.Errorf("TemplateName = %s, want %s", updated.TemplateName, "新模板")
	}
}

func TestUserTemplateService_Delete(t *testing.T) {
	service, db := setupTestService(t)

	db.Create(&model.UserTemplate{
		ID:              1,
		UserID:          123,
		TemplateName:    "测试模板",
		TemplateContent: "距离{exam}还有{time}",
	})

	err := service.Delete(1)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	result, _ := service.GetByID(1)
	if result != nil {
		t.Error("Template should be deleted")
	}
}

func TestUserTemplateService_CountByUserID(t *testing.T) {
	service, db := setupTestService(t)

	userID := int64(123)

	for i := 0; i < 3; i++ {
		db.Create(&model.UserTemplate{
			ID:              int64(i + 1),
			UserID:          userID,
			TemplateName:    "模板",
			TemplateContent: "距离{exam}还有{time}",
		})
	}

	count, err := service.CountByUserID(userID)
	if err != nil {
		t.Errorf("CountByUserID() error = %v", err)
	}

	if count != 3 {
		t.Errorf("Expected count = 3, got %d", count)
	}
}

func TestUserTemplateService_GetDefaultTemplate(t *testing.T) {
	service, db := setupTestService(t)

	db.Create(&model.UserTemplate{
		ID:              1,
		UserID:          0,
		TemplateName:    "默认模板",
		TemplateContent: "距离{exam}还有{time}",
	})

	result, err := service.GetDefaultTemplate()
	if err != nil {
		t.Errorf("GetDefaultTemplate() error = %v", err)
	}

	if result == nil {
		t.Error("Expected default template, got nil")
		return
	}

	if result.UserID != 0 {
		t.Errorf("Expected UserID = 0, got %d", result.UserID)
	}
}
