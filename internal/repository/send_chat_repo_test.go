package repository

import (
	"fmt"
	"testing"

	"github.com/herbertgao/gaokao_bot/internal/model"
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

func TestSendChatRepository_Create(t *testing.T) {
	db := setupSendChatTestDB(t)
	repo := NewSendChatRepository(db)

	chat := &model.SendChat{
		ID:     1,
		ChatID: "123456789",
	}

	err := repo.Create(chat)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	// 验证创建成功
	var count int64
	db.Model(&model.SendChat{}).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 chat, got %d", count)
	}
}

func TestSendChatRepository_GetAll(t *testing.T) {
	db := setupSendChatTestDB(t)
	repo := NewSendChatRepository(db)

	// 插入测试数据
	db.Create(&model.SendChat{ID: 1, ChatID: "111"})
	db.Create(&model.SendChat{ID: 2, ChatID: "222"})
	db.Create(&model.SendChat{ID: 3, ChatID: "333"})

	result, err := repo.GetAll()
	if err != nil {
		t.Errorf("GetAll() error = %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 chats, got %d", len(result))
	}
}

func TestSendChatRepository_GetAll_Empty(t *testing.T) {
	db := setupSendChatTestDB(t)
	repo := NewSendChatRepository(db)

	result, err := repo.GetAll()
	if err != nil {
		t.Errorf("GetAll() error = %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 chats, got %d", len(result))
	}
}

func TestSendChatRepository_Delete(t *testing.T) {
	db := setupSendChatTestDB(t)
	repo := NewSendChatRepository(db)

	// 插入测试数据
	db.Create(&model.SendChat{ID: 1, ChatID: "123"})

	err := repo.Delete(1)
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

func TestSendChatRepository_Delete_NotFound(t *testing.T) {
	db := setupSendChatTestDB(t)
	repo := NewSendChatRepository(db)

	// 删除不存在的记录不应该报错
	err := repo.Delete(999)
	if err != nil {
		t.Errorf("Delete() should not error for not found, got %v", err)
	}
}

func TestSendChatRepository_CreateMultiple(t *testing.T) {
	db := setupSendChatTestDB(t)
	repo := NewSendChatRepository(db)

	// 创建多个chat
	for i := 1; i <= 5; i++ {
		chat := &model.SendChat{
			ID:     int64(i),
			ChatID: fmt.Sprintf("%d", 100+i),
		}
		err := repo.Create(chat)
		if err != nil {
			t.Errorf("Create() error = %v for chat %d", err, i)
		}
	}

	// 验证总数
	result, _ := repo.GetAll()
	if len(result) != 5 {
		t.Errorf("Expected 5 chats, got %d", len(result))
	}
}
