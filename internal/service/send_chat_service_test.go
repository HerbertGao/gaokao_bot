package service

import (
	"testing"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSendChatTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&model.SendChat{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func setupSendChatTestService(t *testing.T) (*SendChatService, *gorm.DB) {
	db := setupSendChatTestDB(t)
	repo := repository.NewSendChatRepository(db)
	service := NewSendChatService(repo)
	return service, db
}

func TestSendChatService_Create(t *testing.T) {
	service, _ := setupSendChatTestService(t)

	chat := &model.SendChat{
		ID:     1,
		ChatID: "123456789",
	}

	err := service.Create(chat)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}
}

func TestSendChatService_GetAll(t *testing.T) {
	service, db := setupSendChatTestService(t)

	// 插入测试数据
	db.Create(&model.SendChat{ID: 1, ChatID: "111"})
	db.Create(&model.SendChat{ID: 2, ChatID: "222"})
	db.Create(&model.SendChat{ID: 3, ChatID: "333"})

	result, err := service.GetAll()
	if err != nil {
		t.Errorf("GetAll() error = %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 chats, got %d", len(result))
	}
}

func TestSendChatService_GetAll_Empty(t *testing.T) {
	service, _ := setupSendChatTestService(t)

	result, err := service.GetAll()
	if err != nil {
		t.Errorf("GetAll() error = %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 chats, got %d", len(result))
	}
}

func TestSendChatService_Delete(t *testing.T) {
	service, db := setupSendChatTestService(t)

	// 插入测试数据
	db.Create(&model.SendChat{ID: 1, ChatID: "123"})

	err := service.Delete(1)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// 验证删除
	var count int64
	db.Model(&model.SendChat{}).Where("id = ?", 1).Count(&count)
	if count != 0 {
		t.Error("Chat should be deleted")
	}
}

func TestSendChatService_Delete_NotFound(t *testing.T) {
	service, _ := setupSendChatTestService(t)

	// 删除不存在的记录不应该报错
	err := service.Delete(999)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}
}
