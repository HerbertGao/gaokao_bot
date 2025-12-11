package service

import (
	"testing"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/repository"
	"github.com/mymmrac/telego"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupInlineQueryTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&model.ExamDate{}, &model.UserTemplate{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func setupInlineQueryTestService(t *testing.T) (*InlineQueryService, *gorm.DB) {
	db := setupInlineQueryTestDB(t)

	examDateRepo := repository.NewExamDateRepository(db)
	examDateService := NewExamDateService(examDateRepo)

	userTemplateRepo := repository.NewUserTemplateRepository(db)
	userTemplateService := NewUserTemplateService(userTemplateRepo)

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	inlineQueryService := NewInlineQueryService(examDateService, userTemplateService, logger)

	return inlineQueryService, db
}

func TestInlineQueryService_GetInlineQueryResults_NoQuery(t *testing.T) {
	service, db := setupInlineQueryTestService(t)

	now := time.Now()
	futureDate := now.AddDate(1, 0, 0)

	// 插入考试数据
	db.Create(&model.ExamDate{
		ID:                1,
		ExamYear:          futureDate.Year(),
		ExamDesc:          "高考",
		ShortDesc:         "高考",
		ExamBeginDate:     futureDate,
		ExamEndDate:       futureDate.AddDate(0, 0, 3),
		ExamYearBeginDate: now,
		ExamYearEndDate:   futureDate.AddDate(0, 0, 3),
		IsDelete:          false,
	})

	// 插入默认模板
	db.Create(&model.UserTemplate{
		ID:              1,
		UserID:          0,
		TemplateContent: "距离{exam}还有{time}",
	})

	query := &telego.InlineQuery{
		ID:    "test",
		Query: "",
		From:  telego.User{ID: 123},
	}

	results := service.GetInlineQueryResults(query)

	if len(results) == 0 {
		t.Error("Expected at least one result")
	}
}

func TestInlineQueryService_GetInlineQueryResults_WithYear(t *testing.T) {
	service, db := setupInlineQueryTestService(t)

	year := 2026

	// 插入指定年份的考试
	db.Create(&model.ExamDate{
		ID:                1,
		ExamYear:          year,
		ExamDesc:          "2026年高考",
		ShortDesc:         "高考",
		ExamBeginDate:     time.Date(year, 6, 7, 9, 0, 0, 0, time.UTC),
		ExamEndDate:       time.Date(year, 6, 10, 17, 0, 0, 0, time.UTC),
		ExamYearBeginDate: time.Date(year-1, 6, 10, 17, 0, 0, 0, time.UTC),
		ExamYearEndDate:   time.Date(year, 6, 10, 17, 0, 0, 0, time.UTC),
		IsDelete:          false,
	})

	// 插入默认模板
	db.Create(&model.UserTemplate{
		ID:              1,
		UserID:          0,
		TemplateContent: "距离{exam}还有{time}",
	})

	query := &telego.InlineQuery{
		ID:    "test",
		Query: "2026",
		From:  telego.User{ID: 123},
	}

	results := service.GetInlineQueryResults(query)

	if len(results) == 0 {
		t.Error("Expected at least one result")
	}
}

func TestInlineQueryService_GetInlineQueryResults_InvalidYear(t *testing.T) {
	service, _ := setupInlineQueryTestService(t)

	query := &telego.InlineQuery{
		ID:    "test",
		Query: "2017", // 小于2018
		From:  telego.User{ID: 123},
	}

	results := service.GetInlineQueryResults(query)

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestInlineQueryService_GetInlineQueryResults_NonNumericQuery(t *testing.T) {
	service, _ := setupInlineQueryTestService(t)

	query := &telego.InlineQuery{
		ID:    "test",
		Query: "hello",
		From:  telego.User{ID: 123},
	}

	results := service.GetInlineQueryResults(query)

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestInlineQueryService_GetInlineQueryResults_NoExams(t *testing.T) {
	service, _ := setupInlineQueryTestService(t)

	query := &telego.InlineQuery{
		ID:    "test",
		Query: "",
		From:  telego.User{ID: 123},
	}

	results := service.GetInlineQueryResults(query)

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestInlineQueryService_GetInlineQueryResults_WithUserTemplates(t *testing.T) {
	service, db := setupInlineQueryTestService(t)

	now := time.Now()
	futureDate := now.AddDate(1, 0, 0)
	userID := int64(123)

	// 插入考试数据
	db.Create(&model.ExamDate{
		ID:                1,
		ExamYear:          futureDate.Year(),
		ExamDesc:          "高考",
		ShortDesc:         "高考",
		ExamBeginDate:     futureDate,
		ExamEndDate:       futureDate.AddDate(0, 0, 3),
		ExamYearBeginDate: now,
		ExamYearEndDate:   futureDate.AddDate(0, 0, 3),
		IsDelete:          false,
	})

	// 插入默认模板
	db.Create(&model.UserTemplate{
		ID:              1,
		UserID:          0,
		TemplateContent: "距离{exam}还有{time}",
	})

	// 插入用户自定义模板
	db.Create(&model.UserTemplate{
		ID:              2,
		UserID:          userID,
		TemplateName:    "我的模板",
		TemplateContent: "倒计时：{exam} - {time}",
	})

	query := &telego.InlineQuery{
		ID:    "test",
		Query: "",
		From:  telego.User{ID: userID},
	}

	results := service.GetInlineQueryResults(query)

	// 应该有2个结果：默认模板 + 用户模板
	if len(results) != 2 {
		t.Errorf("Expected 2 results (default + user template), got %d", len(results))
	}
}

func TestInlineQueryService_GetInlineQueryResults_MultipleExams(t *testing.T) {
	service, db := setupInlineQueryTestService(t)

	now := time.Now()
	futureDate1 := now.AddDate(1, 0, 0)
	futureDate2 := now.AddDate(2, 0, 0)

	// 插入多个考试
	db.Create(&model.ExamDate{
		ID:                1,
		ExamYear:          futureDate1.Year(),
		ExamDesc:          "高考1",
		ShortDesc:         "高考1",
		ExamBeginDate:     futureDate1,
		ExamEndDate:       futureDate1.AddDate(0, 0, 3),
		ExamYearBeginDate: now,
		ExamYearEndDate:   futureDate2.AddDate(0, 0, 3),
		IsDelete:          false,
	})

	db.Create(&model.ExamDate{
		ID:                2,
		ExamYear:          futureDate2.Year(),
		ExamDesc:          "高考2",
		ShortDesc:         "高考2",
		ExamBeginDate:     futureDate2,
		ExamEndDate:       futureDate2.AddDate(0, 0, 3),
		ExamYearBeginDate: now,
		ExamYearEndDate:   futureDate2.AddDate(0, 0, 3),
		IsDelete:          false,
	})

	// 插入默认模板
	db.Create(&model.UserTemplate{
		ID:              1,
		UserID:          0,
		TemplateContent: "距离{exam}还有{time}",
	})

	query := &telego.InlineQuery{
		ID:    "test",
		Query: "",
		From:  telego.User{ID: 123},
	}

	results := service.GetInlineQueryResults(query)

	// 应该有2个结果（每个考试一个）
	if len(results) != 2 {
		t.Errorf("Expected 2 results (one per exam), got %d", len(results))
	}
}

func TestInlineQueryService_GetInlineQueryResults_NoDefaultTemplate(t *testing.T) {
	service, db := setupInlineQueryTestService(t)

	now := time.Now()
	futureDate := now.AddDate(1, 0, 0)

	// 插入考试数据
	db.Create(&model.ExamDate{
		ID:                1,
		ExamYear:          futureDate.Year(),
		ExamDesc:          "高考",
		ShortDesc:         "高考",
		ExamBeginDate:     futureDate,
		ExamEndDate:       futureDate.AddDate(0, 0, 3),
		ExamYearBeginDate: now,
		ExamYearEndDate:   futureDate.AddDate(0, 0, 3),
		IsDelete:          false,
	})

	// 不插入默认模板

	query := &telego.InlineQuery{
		ID:    "test",
		Query: "",
		From:  telego.User{ID: 123},
	}

	results := service.GetInlineQueryResults(query)

	// 没有默认模板，应该没有结果
	if len(results) != 0 {
		t.Errorf("Expected 0 results without default template, got %d", len(results))
	}
}

func TestInlineQueryService_GetInlineQueryResults_WithMultipleUserTemplates(t *testing.T) {
	service, db := setupInlineQueryTestService(t)

	now := time.Now()
	futureDate := now.AddDate(1, 0, 0)
	userID := int64(456)

	// 插入考试数据
	db.Create(&model.ExamDate{
		ID:                1,
		ExamYear:          futureDate.Year(),
		ExamDesc:          "高考",
		ShortDesc:         "高考",
		ExamBeginDate:     futureDate,
		ExamEndDate:       futureDate.AddDate(0, 0, 3),
		ExamYearBeginDate: now,
		ExamYearEndDate:   futureDate.AddDate(0, 0, 3),
		IsDelete:          false,
	})

	// 插入默认模板
	db.Create(&model.UserTemplate{
		ID:              1,
		UserID:          0,
		TemplateContent: "距离{exam}还有{time}",
	})

	// 插入多个用户自定义模板
	db.Create(&model.UserTemplate{
		ID:              2,
		UserID:          userID,
		TemplateName:    "模板1",
		TemplateContent: "倒计时1：{exam} - {time}",
	})

	db.Create(&model.UserTemplate{
		ID:              3,
		UserID:          userID,
		TemplateName:    "模板2",
		TemplateContent: "倒计时2：{exam} - {time}",
	})

	query := &telego.InlineQuery{
		ID:    "test",
		Query: "",
		From:  telego.User{ID: userID},
	}

	results := service.GetInlineQueryResults(query)

	// 应该有3个结果：默认模板 + 2个用户模板
	if len(results) != 3 {
		t.Errorf("Expected 3 results (default + 2 user templates), got %d", len(results))
	}
}
