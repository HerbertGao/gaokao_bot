package service

import (
	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/repository"
)

// SendChatService 发送对话服务
type SendChatService struct {
	repo *repository.SendChatRepository
}

// NewSendChatService 创建发送对话服务
func NewSendChatService(repo *repository.SendChatRepository) *SendChatService {
	return &SendChatService{repo: repo}
}

// GetAll 获取所有发送对话
func (s *SendChatService) GetAll() ([]model.SendChat, error) {
	return s.repo.GetAll()
}

// Create 创建发送对话
func (s *SendChatService) Create(chat *model.SendChat) error {
	return s.repo.Create(chat)
}

// Delete 删除发送对话
func (s *SendChatService) Delete(id int64) error {
	return s.repo.Delete(id)
}
