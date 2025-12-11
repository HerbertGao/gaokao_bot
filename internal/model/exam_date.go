package model

import "time"

// ExamDate 考试日期实体
type ExamDate struct {
	ID                uint      `gorm:"primaryKey;autoIncrement"`
	ExamYear          int       `gorm:"not null;index"`
	ExamDesc          string    `gorm:"type:varchar(255)"`
	ShortDesc         string    `gorm:"type:varchar(32)"`
	ExamBeginDate     time.Time `gorm:"not null"`
	ExamEndDate       time.Time `gorm:"not null"`
	ExamYearBeginDate time.Time `gorm:"not null"`
	ExamYearEndDate   time.Time `gorm:"not null"`
	IsDelete          bool      `gorm:"default:false"`
}

// TableName 指定表名
func (ExamDate) TableName() string {
	return "exam_date"
}