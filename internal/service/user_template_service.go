package service

import (
	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/repository"
)

// UserTemplateService 用户模板服务
type UserTemplateService struct {
	repo *repository.UserTemplateRepository
}

// NewUserTemplateService 创建用户模板服务
func NewUserTemplateService(repo *repository.UserTemplateRepository) *UserTemplateService {
	return &UserTemplateService{repo: repo}
}

// GetByUserID 根据用户ID获取模板列表
func (s *UserTemplateService) GetByUserID(userID int64) ([]model.UserTemplate, error) {
	return s.repo.GetByUserID(userID)
}

// GetDefaultTemplate 获取默认模板
func (s *UserTemplateService) GetDefaultTemplate() (*model.UserTemplate, error) {
	return s.repo.GetDefaultTemplate()
}

// Create 创建模板
func (s *UserTemplateService) Create(template *model.UserTemplate) error {
	return s.repo.Create(template)
}

// Update 更新模板
func (s *UserTemplateService) Update(template *model.UserTemplate) error {
	return s.repo.Update(template)
}

// Delete 删除模板
func (s *UserTemplateService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// GetByID 根据ID获取模板
func (s *UserTemplateService) GetByID(id int64) (*model.UserTemplate, error) {
	return s.repo.GetByID(id)
}

// CountByUserID 统计用户的模板数量
func (s *UserTemplateService) CountByUserID(userID int64) (int64, error) {
	return s.repo.CountByUserID(userID)
}
