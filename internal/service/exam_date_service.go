package service

import (
	"time"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/repository"
)

// ExamDateService 考试日期服务
type ExamDateService struct {
	repo *repository.ExamDateRepository
}

// NewExamDateService 创建考试日期服务
func NewExamDateService(repo *repository.ExamDateRepository) *ExamDateService {
	return &ExamDateService{repo: repo}
}

// GetExamsInRange 获取时间范围内的考试
func (s *ExamDateService) GetExamsInRange(now time.Time) ([]model.ExamDate, error) {
	return s.repo.GetExamsInRange(now)
}

// GetExamByYear 按年份获取考试
func (s *ExamDateService) GetExamByYear(year int) ([]model.ExamDate, error) {
	return s.repo.GetExamByYear(year)
}

// GetNextExamDate 获取下一个高考日期
func (s *ExamDateService) GetNextExamDate() (*model.ExamDate, error) {
	now := time.Now()
	exams, err := s.repo.GetExamsInRange(now)
	if err != nil {
		return nil, err
	}

	if len(exams) == 0 {
		return nil, nil
	}

	// 返回第一个（最近的）考试
	return &exams[0], nil
}
