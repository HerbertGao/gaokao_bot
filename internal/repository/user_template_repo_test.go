package repository

import (
	"testing"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&model.UserTemplate{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func TestUserTemplateRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserTemplateRepository(db)

	template := &model.UserTemplate{
		ID:              1,
		UserID:          123,
		TemplateName:    "测试模板",
		TemplateContent: "距离{exam}还有{time}",
	}

	err := repo.Create(template)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	// 验证创建成功
	var count int64
	db.Model(&model.UserTemplate{}).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 template, got %d", count)
	}
}

func TestUserTemplateRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserTemplateRepository(db)

	// 插入测试数据
	template := &model.UserTemplate{
		ID:              1,
		UserID:          123,
		TemplateName:    "测试模板",
		TemplateContent: "距离{exam}还有{time}",
	}
	db.Create(template)

	// 测试获取
	result, err := repo.GetByID(1)
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

func TestUserTemplateRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserTemplateRepository(db)

	result, err := repo.GetByID(999)
	if err != nil {
		t.Errorf("GetByID() should not return error for not found, got %v", err)
	}

	if result != nil {
		t.Error("Expected nil for not found, got template")
	}
}

func TestUserTemplateRepository_GetByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserTemplateRepository(db)

	userID := int64(123)

	// 插入测试数据
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
	db.Create(&model.UserTemplate{
		ID:              3,
		UserID:          456, // 其他用户
		TemplateName:    "模板3",
		TemplateContent: "距离{exam}还有{time}",
	})

	result, err := repo.GetByUserID(userID)
	if err != nil {
		t.Errorf("GetByUserID() error = %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 templates, got %d", len(result))
	}
}

func TestUserTemplateRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserTemplateRepository(db)

	// 插入测试数据
	template := &model.UserTemplate{
		ID:              1,
		UserID:          123,
		TemplateName:    "旧模板",
		TemplateContent: "距离{exam}还有{time}",
	}
	db.Create(template)

	// 更新
	template.TemplateName = "新模板"
	template.TemplateContent = "{exam}倒计时{time}"

	err := repo.Update(template)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	// 验证更新
	var updated model.UserTemplate
	db.First(&updated, 1)
	if updated.TemplateName != "新模板" {
		t.Errorf("TemplateName = %s, want %s", updated.TemplateName, "新模板")
	}
}

func TestUserTemplateRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserTemplateRepository(db)

	// 插入测试数据
	db.Create(&model.UserTemplate{
		ID:              1,
		UserID:          123,
		TemplateName:    "测试模板",
		TemplateContent: "距离{exam}还有{time}",
	})

	err := repo.Delete(1)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// 验证删除
	var count int64
	db.Model(&model.UserTemplate{}).Where("id = ?", 1).Count(&count)
	if count != 0 {
		t.Error("Template should be deleted")
	}
}

func TestUserTemplateRepository_CountByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserTemplateRepository(db)

	userID := int64(123)

	// 插入测试数据
	for i := 0; i < 5; i++ {
		db.Create(&model.UserTemplate{
			ID:              int64(i + 1),
			UserID:          userID,
			TemplateName:    "模板",
			TemplateContent: "距离{exam}还有{time}",
		})
	}

	count, err := repo.CountByUserID(userID)
	if err != nil {
		t.Errorf("CountByUserID() error = %v", err)
	}

	if count != 5 {
		t.Errorf("Expected count = 5, got %d", count)
	}
}

func TestUserTemplateRepository_GetDefaultTemplate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserTemplateRepository(db)

	// 插入测试数据（默认模板的 user_id = 0）
	db.Create(&model.UserTemplate{
		ID:              1,
		UserID:          0,
		TemplateName:    "默认模板",
		TemplateContent: "距离{exam}还有{time}",
	})

	result, err := repo.GetDefaultTemplate()
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
