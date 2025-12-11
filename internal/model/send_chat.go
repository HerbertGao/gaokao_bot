package model

// SendChat 发送对话实体
type SendChat struct {
	ID     int64  `gorm:"primaryKey;autoIncrement"`
	ChatID string `gorm:"type:varchar(64);not null"`
}

// TableName 指定表名
func (SendChat) TableName() string {
	return "send_chat"
}