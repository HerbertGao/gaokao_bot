package repository

import (
	"testing"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/model"
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

func TestExamDateRepository_GetExamsInRange(t *testing.T) {
	db := setupExamDateTestDB(t)
	repo := NewExamDateRepository(db)

	now := time.Now()
	futureDate := now.AddDate(1, 0, 0)

	// 插入在范围内的考试
	db.Create(&model.ExamDate{
		ID:                1,
		ExamYear:          futureDate.Year(),
		ExamDesc:          "范围内考试",
		ShortDesc:         "考试",
		ExamBeginDate:     futureDate,
		ExamEndDate:       futureDate.AddDate(0, 0, 3),
		ExamYearBeginDate: now.AddDate(-1, 0, 0),
		ExamYearEndDate:   futureDate.AddDate(0, 0, 3),
		IsDelete:          false,
	})

	// 插入不在范围内的考试（过去的）
	pastDate := now.AddDate(-2, 0, 0)
	db.Create(&model.ExamDate{
		ID:                2,
		ExamYear:          pastDate.Year(),
		ExamDesc:          "过去考试",
		ShortDesc:         "考试",
		ExamBeginDate:     pastDate,
		ExamEndDate:       pastDate.AddDate(0, 0, 3),
		ExamYearBeginDate: pastDate.AddDate(-1, 0, 0),
		ExamYearEndDate:   pastDate.AddDate(0, 0, 3),
		IsDelete:          false,
	})

	// 插入已删除的考试
	db.Create(&model.ExamDate{
		ID:                3,
		ExamYear:          futureDate.Year(),
		ExamDesc:          "已删除考试",
		ShortDesc:         "考试",
		ExamBeginDate:     futureDate,
		ExamEndDate:       futureDate.AddDate(0, 0, 3),
		ExamYearBeginDate: now.AddDate(-1, 0, 0),
		ExamYearEndDate:   futureDate.AddDate(0, 0, 3),
		IsDelete:          true,
	})

	result, err := repo.GetExamsInRange(now)
	if err != nil {
		t.Errorf("GetExamsInRange() error = %v", err)
	}

	// 应该只返回未删除且在范围内的考试
	if len(result) != 1 {
		t.Errorf("Expected 1 exam, got %d", len(result))
	}

	if len(result) > 0 && result[0].ID != 1 {
		t.Errorf("Expected exam ID 1, got %d", result[0].ID)
	}
}

func TestExamDateRepository_GetExamByYear(t *testing.T) {
	db := setupExamDateTestDB(t)
	repo := NewExamDateRepository(db)

	year := 2026

	// 插入指定年份的考试
	db.Create(&model.ExamDate{
		ID:                1,
		ExamYear:          year,
		ExamDesc:          "2026年考试",
		ShortDesc:         "考试",
		ExamBeginDate:     time.Date(year, 6, 7, 9, 0, 0, 0, time.UTC),
		ExamEndDate:       time.Date(year, 6, 10, 17, 0, 0, 0, time.UTC),
		ExamYearBeginDate: time.Date(year-1, 6, 10, 17, 0, 0, 0, time.UTC),
		ExamYearEndDate:   time.Date(year, 6, 10, 17, 0, 0, 0, time.UTC),
		IsDelete:          false,
	})

	// 插入其他年份的考试
	db.Create(&model.ExamDate{
		ID:                2,
		ExamYear:          year + 1,
		ExamDesc:          "2027年考试",
		ShortDesc:         "考试",
		ExamBeginDate:     time.Date(year+1, 6, 7, 9, 0, 0, 0, time.UTC),
		ExamEndDate:       time.Date(year+1, 6, 10, 17, 0, 0, 0, time.UTC),
		ExamYearBeginDate: time.Date(year, 6, 10, 17, 0, 0, 0, time.UTC),
		ExamYearEndDate:   time.Date(year+1, 6, 10, 17, 0, 0, 0, time.UTC),
		IsDelete:          false,
	})

	// 插入已删除的指定年份考试
	db.Create(&model.ExamDate{
		ID:                3,
		ExamYear:          year,
		ExamDesc:          "已删除的2026年考试",
		ShortDesc:         "考试",
		ExamBeginDate:     time.Date(year, 6, 7, 9, 0, 0, 0, time.UTC),
		ExamEndDate:       time.Date(year, 6, 10, 17, 0, 0, 0, time.UTC),
		ExamYearBeginDate: time.Date(year-1, 6, 10, 17, 0, 0, 0, time.UTC),
		ExamYearEndDate:   time.Date(year, 6, 10, 17, 0, 0, 0, time.UTC),
		IsDelete:          true,
	})

	result, err := repo.GetExamByYear(year)
	if err != nil {
		t.Errorf("GetExamByYear() error = %v", err)
	}

	// 应该只返回未删除且年份匹配的考试
	if len(result) != 1 {
		t.Errorf("Expected 1 exam, got %d", len(result))
	}

	if len(result) > 0 {
		if result[0].ExamYear != year {
			t.Errorf("ExamYear = %d, want %d", result[0].ExamYear, year)
		}
		if result[0].ID != 1 {
			t.Errorf("Expected exam ID 1, got %d", result[0].ID)
		}
	}
}

func TestExamDateRepository_GetExamByYear_NotFound(t *testing.T) {
	db := setupExamDateTestDB(t)
	repo := NewExamDateRepository(db)

	result, err := repo.GetExamByYear(2099)
	if err != nil {
		t.Errorf("GetExamByYear() error = %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 exams, got %d", len(result))
	}
}
