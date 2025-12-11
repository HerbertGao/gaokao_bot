package repository

import (
	"github.com/herbertgao/gaokao_bot/internal/model"
	"gorm.io/gorm"
)

// UserTemplateRepository 用户模板仓储
type UserTemplateRepository struct {
	db *gorm.DB
}

// NewUserTemplateRepository 创建用户模板仓储
func NewUserTemplateRepository(db *gorm.DB) *UserTemplateRepository {
	return &UserTemplateRepository{db: db}
}

// GetByUserID 根据用户ID获取模板列表
func (r *UserTemplateRepository) GetByUserID(userID int64) ([]model.UserTemplate, error) {
	var templates []model.UserTemplate

	err := r.db.Where("user_id = ?", userID).Find(&templates).Error

	return templates, err
}

// GetDefaultTemplate 获取默认模板
func (r *UserTemplateRepository) GetDefaultTemplate() (*model.UserTemplate, error) {
	var template model.UserTemplate

	err := r.db.Where("user_id = ?", 0).First(&template).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &template, err
}

// Create 创建模板
func (r *UserTemplateRepository) Create(template *model.UserTemplate) error {
	return r.db.Create(template).Error
}

// Update 更新模板
func (r *UserTemplateRepository) Update(template *model.UserTemplate) error {
	return r.db.Save(template).Error
}

// Delete 删除模板
func (r *UserTemplateRepository) Delete(id int64) error {
	return r.db.Delete(&model.UserTemplate{}, id).Error
}

// GetByID 根据ID获取模板
func (r *UserTemplateRepository) GetByID(id int64) (*model.UserTemplate, error) {
	var template model.UserTemplate

	err := r.db.First(&template, id).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &template, err
}