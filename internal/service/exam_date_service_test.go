package service

import (
	"testing"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/repository"
	"github.com/herbertgao/gaokao_bot/internal/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupExamDateTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&model.ExamDate{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func setupExamDateTestService(t *testing.T) (*ExamDateService, *gorm.DB) {
	db := setupExamDateTestDB(t)
	repo := repository.NewExamDateRepository(db)
	service := NewExamDateService(repo)
	return service, db
}

func TestExamDateService_GetExamsInRange(t *testing.T) {
	service, db := setupExamDateTestService(t)

	// 插入测试数据
	now := time.Now()
	futureDate := now.AddDate(1, 0, 0) // 1年后

	db.Create(&model.ExamDate{
		ID:                1,
		ExamYear:          futureDate.Year(),
		ExamDesc:          "未来高考",
		ShortDesc:         "高考",
		ExamBeginDate:     futureDate,
		ExamEndDate:       futureDate.AddDate(0, 0, 3),
		ExamYearBeginDate: now,
		ExamYearEndDate:   futureDate.AddDate(0, 0, 3),
		IsDelete:          false,
	})

	result, err := service.GetExamsInRange(now)
	if err != nil {
		t.Errorf("GetExamsInRange() error = %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected at least one exam")
	}
}

func TestExamDateService_GetExamByYear(t *testing.T) {
	service, db := setupExamDateTestService(t)

	year := 2026
	bjtZone := util.GetBJTLocation()

	db.Create(&model.ExamDate{
		ID:                1,
		ExamYear:          year,
		ExamDesc:          "2026年高考",
		ShortDesc:         "2026高考",
		ExamBeginDate:     time.Date(year, 6, 7, 9, 0, 0, 0, bjtZone),
		ExamEndDate:       time.Date(year, 6, 10, 17, 0, 0, 0, bjtZone),
		ExamYearBeginDate: time.Date(year-1, 6, 10, 17, 0, 0, 0, bjtZone),
		ExamYearEndDate:   time.Date(year, 6, 10, 17, 0, 0, 0, bjtZone),
		IsDelete:          false,
	})

	result, err := service.GetExamByYear(year)
	if err != nil {
		t.Errorf("GetExamByYear() error = %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 exam, got %d", len(result))
	}

	if len(result) > 0 && result[0].ExamYear != year {
		t.Errorf("ExamYear = %d, want %d", result[0].ExamYear, year)
	}
}

func TestExamDateService_GetNextExamDate(t *testing.T) {
	service, db := setupExamDateTestService(t)

	now := time.Now()
	futureDate := now.AddDate(1, 0, 0)

	db.Create(&model.ExamDate{
		ID:                1,
		ExamYear:          futureDate.Year(),
		ExamDesc:          "下一个高考",
		ShortDesc:         "高考",
		ExamBeginDate:     futureDate,
		ExamEndDate:       futureDate.AddDate(0, 0, 3),
		ExamYearBeginDate: now,
		ExamYearEndDate:   futureDate.AddDate(0, 0, 3),
		IsDelete:          false,
	})

	result, err := service.GetNextExamDate()
	if err != nil {
		t.Errorf("GetNextExamDate() error = %v", err)
	}

	if result == nil {
		t.Error("Expected exam, got nil")
		return
	}

	if result.ExamYear != futureDate.Year() {
		t.Errorf("ExamYear = %d, want %d", result.ExamYear, futureDate.Year())
	}
}

func TestExamDateService_GetNextExamDate_NoExams(t *testing.T) {
	service, _ := setupExamDateTestService(t)

	result, err := service.GetNextExamDate()
	if err != nil {
		t.Errorf("GetNextExamDate() error = %v", err)
	}

	if result != nil {
		t.Error("Expected nil when no exams, got exam")
	}
}
