package repository

import (
	"time"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"gorm.io/gorm"
)

// ExamDateRepository 考试日期仓储
type ExamDateRepository struct {
	db *gorm.DB
}

// NewExamDateRepository 创建考试日期仓储
func NewExamDateRepository(db *gorm.DB) *ExamDateRepository {
	return &ExamDateRepository{db: db}
}

// GetExamsInRange 获取时间范围内的考试
func (r *ExamDateRepository) GetExamsInRange(now time.Time) ([]model.ExamDate, error) {
	var exams []model.ExamDate

	err := r.db.Where("exam_year_begin_date <= ? AND exam_year_end_date >= ? AND is_delete = ?",
		now, now, false).
		Find(&exams).Error

	return exams, err
}

// GetExamByYear 按年份获取考试
func (r *ExamDateRepository) GetExamByYear(year int) ([]model.ExamDate, error) {
	var exams []model.ExamDate

	err := r.db.Where("exam_year = ? AND is_delete = ?", year, false).
		Find(&exams).Error

	return exams, err
}