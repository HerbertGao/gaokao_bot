package repository

import (
	"github.com/herbertgao/gaokao_bot/internal/model"
	"gorm.io/gorm"
)

// SendChatRepository 发送对话仓储
type SendChatRepository struct {
	db *gorm.DB
}

// NewSendChatRepository 创建发送对话仓储
func NewSendChatRepository(db *gorm.DB) *SendChatRepository {
	return &SendChatRepository{db: db}
}

// GetAll 获取所有发送对话
func (r *SendChatRepository) GetAll() ([]model.SendChat, error) {
	var chats []model.SendChat

	err := r.db.Find(&chats).Error

	return chats, err
}

// Create 创建发送对话
func (r *SendChatRepository) Create(chat *model.SendChat) error {
	return r.db.Create(chat).Error
}

// Delete 删除发送对话
func (r *SendChatRepository) Delete(id int64) error {
	return r.db.Delete(&model.SendChat{}, id).Error
}