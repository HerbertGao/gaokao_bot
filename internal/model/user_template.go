package model

import "time"

// UserTemplate 用户模板实体
type UserTemplate struct {
	ID              int64     `gorm:"primaryKey" json:"id,string"`
	UserID          int64     `gorm:"not null;index" json:"user_id,string"`
	TemplateName    string    `gorm:"type:varchar(40)" json:"template_name"`
	TemplateContent string    `gorm:"type:varchar(160)" json:"template_content"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (UserTemplate) TableName() string {
	return "user_template"
}