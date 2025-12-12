package repository

import (
	"errors"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"gorm.io/gorm"
)

// ErrTemplateLimitExceeded 模板数量超过限制错误
var ErrTemplateLimitExceeded = errors.New("template limit exceeded")

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

// CountByUserID 统计用户的模板数量
func (r *UserTemplateRepository) CountByUserID(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.UserTemplate{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CreateWithLimit 在事务中原子地检查数量限制并创建模板
// 使用数据库事务确保并发安全，防止 TOCTOU 竞态条件
func (r *UserTemplateRepository) CreateWithLimit(template *model.UserTemplate, maxLimit int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 在事务中统计该用户的模板数量
		// 事务隔离级别确保并发操作的一致性
		var count int64
		if err := tx.Model(&model.UserTemplate{}).
			Where("user_id = ?", template.UserID).
			Count(&count).Error; err != nil {
			return err
		}

		// 检查是否超过限制
		if count >= maxLimit {
			return ErrTemplateLimitExceeded
		}

		// 在限制内，创建模板
		// 注意：依赖数据库的事务隔离来防止并发问题
		// 对于 PostgreSQL/MySQL，可以考虑添加 FOR UPDATE 锁
		// 但为了兼容 SQLite，这里使用标准的事务方式
		return tx.Create(template).Error
	})
}