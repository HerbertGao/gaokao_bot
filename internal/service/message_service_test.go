package service

import (
	"testing"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/util"
	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/repository"
	"github.com/mymmrac/telego"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupMessageTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&model.ExamDate{}, &model.UserTemplate{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func setupMessageTestService(t *testing.T) (*MessageService, *gorm.DB) {
	db := setupMessageTestDB(t)

	examDateRepo := repository.NewExamDateRepository(db)
	examDateService := NewExamDateService(examDateRepo)

	userTemplateRepo := repository.NewUserTemplateRepository(db)
	userTemplateService := NewUserTemplateService(userTemplateRepo)

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	messageService := NewMessageService(examDateService, userTemplateService, logger)

	return messageService, db
}

func TestMessageService_GetCountDownMessage_NoParams(t *testing.T) {
	service, db := setupMessageTestService(t)

	now := time.Now()
	futureDate := now.AddDate(1, 0, 0)

	// 插入未来的考试
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

	// 创建没有文本的消息
	msg := &telego.Message{
		Text: "",
	}

	result, err := service.GetCountDownMessage(msg)
	if err != nil {
		t.Errorf("GetCountDownMessage() error = %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestMessageService_GetCountDownMessage_WithYear(t *testing.T) {
	service, db := setupMessageTestService(t)

	year := 2026

	// 插入指定年份的考试
	db.Create(&model.ExamDate{
		ID:                1,
		ExamYear:          year,
		ExamDesc:          "2026年高考",
		ShortDesc:         "高考",
		ExamBeginDate:     time.Date(year, 6, 7, 9, 0, 0, 0, util.GetBJTLocation()),
		ExamEndDate:       time.Date(year, 6, 10, 17, 0, 0, 0, util.GetBJTLocation()),
		ExamYearBeginDate: time.Date(year-1, 6, 10, 17, 0, 0, 0, util.GetBJTLocation()),
		ExamYearEndDate:   time.Date(year, 6, 10, 17, 0, 0, 0, util.GetBJTLocation()),
		IsDelete:          false,
	})

	// 创建带年份的消息
	msg := &telego.Message{
		Text: "2026",
	}

	result, err := service.GetCountDownMessage(msg)
	if err != nil {
		t.Errorf("GetCountDownMessage() error = %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestMessageService_GetCountDownMessage_InvalidYear(t *testing.T) {
	service, _ := setupMessageTestService(t)

	// 创建带无效年份的消息
	msg := &telego.Message{
		Text: "2017", // 小于2018
	}

	result, err := service.GetCountDownMessage(msg)
	if err != nil {
		t.Errorf("GetCountDownMessage() should not error for invalid year, got %v", err)
	}

	if result != "参数暂时无法识别。" {
		t.Errorf("Expected '参数暂时无法识别。', got %s", result)
	}
}

func TestMessageService_GetCountDownMessage_NonNumericText(t *testing.T) {
	service, _ := setupMessageTestService(t)

	// 创建带非数字文本的消息
	msg := &telego.Message{
		Text: "hello",
	}

	result, err := service.GetCountDownMessage(msg)
	if err != nil {
		t.Errorf("GetCountDownMessage() should not error for non-numeric text, got %v", err)
	}

	if result != "参数暂时无法识别。" {
		t.Errorf("Expected '参数暂时无法识别。', got %s", result)
	}
}

func TestMessageService_GetCountDownMessage_YearNotFound(t *testing.T) {
	service, _ := setupMessageTestService(t)

	// 查询一个不存在的年份
	msg := &telego.Message{
		Text: "2099",
	}

	result, err := service.GetCountDownMessage(msg)
	if err != nil {
		t.Errorf("GetCountDownMessage() error = %v", err)
	}

	if result != "参数暂时无法识别。" {
		t.Errorf("Expected '参数暂时无法识别。', got %s", result)
	}
}

func TestMessageService_GetCountDownMessage_NoExamsInRange(t *testing.T) {
	service, _ := setupMessageTestService(t)

	// 不插入任何考试数据
	msg := &telego.Message{
		Text: "",
	}

	result, err := service.GetCountDownMessage(msg)
	if err != nil {
		t.Errorf("GetCountDownMessage() error = %v", err)
	}

	if result != "数据库中没有可用的信息，请联系开发者。" {
		t.Errorf("Expected '数据库中没有可用的信息，请联系开发者。', got %s", result)
	}
}

func TestMessageService_GetCountDownMessage_WithCustomTemplate(t *testing.T) {
	service, db := setupMessageTestService(t)

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

	// 插入自定义模板（UserID = 0 表示默认模板）
	db.Create(&model.UserTemplate{
		ID:              1,
		UserID:          0,
		TemplateContent: "距离{exam}倒计时：{time}",
	})

	msg := &telego.Message{
		Text: "",
	}

	result, err := service.GetCountDownMessage(msg)
	if err != nil {
		t.Errorf("GetCountDownMessage() error = %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestMessageService_GetCountDownMessage_MultipleExams(t *testing.T) {
	service, db := setupMessageTestService(t)

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

	msg := &telego.Message{
		Text: "",
	}

	result, err := service.GetCountDownMessage(msg)
	if err != nil {
		t.Errorf("GetCountDownMessage() error = %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result with multiple exams")
	}
}
