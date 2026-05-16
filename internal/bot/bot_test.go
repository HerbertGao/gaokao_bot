package bot

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/config"
	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/repository"
	"github.com/herbertgao/gaokao_bot/internal/service"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoapi"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 注意：由于 GaokaoBot 依赖真实的 Telegram Bot API，
// 这些测试主要覆盖不需要实际网络调用的部分

func TestNewGaokaoBot(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	cfg := &config.TelegramConfig{
		Bot: config.BotConfig{
			Username: "test_bot",
			Token:    "test_token",
		},
		MiniApp: config.MiniAppConfig{
			URL: "https://example.com",
		},
	}

	// 创建 nil bot（测试中不会调用真实 API）
	var bot *telego.Bot = nil

	// 创建 nil service（测试中不需要）
	var svc *service.BotService = nil

	gaokaoBot, err := NewGaokaoBot(bot, cfg, svc, logger)

	if err != nil {
		t.Errorf("NewGaokaoBot() error = %v, want nil", err)
	}

	if gaokaoBot == nil {
		t.Fatal("NewGaokaoBot() returned nil")
	}

	if gaokaoBot.bot != bot {
		t.Error("bot not set correctly")
	}

	if gaokaoBot.config != cfg {
		t.Error("config not set correctly")
	}

	if gaokaoBot.service != svc {
		t.Error("service not set correctly")
	}

	if gaokaoBot.logger != logger {
		t.Error("logger not set correctly")
	}

	if gaokaoBot.done == nil {
		t.Error("done channel not initialized")
	}

	// 验证初始状态
	if gaokaoBot.started {
		t.Error("bot should not be marked as started initially")
	}
}

// TestStop_BeforeStart 测试在 Start() 之前调用 Stop() 不会 panic
func TestStop_BeforeStart(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	cfg := &config.TelegramConfig{
		Bot: config.BotConfig{
			Username: "test_bot",
			Token:    "test_token",
		},
		MiniApp: config.MiniAppConfig{
			URL: "https://example.com",
		},
	}

	gaokaoBot, err := NewGaokaoBot(nil, cfg, nil, logger)
	if err != nil {
		t.Fatalf("NewGaokaoBot() error = %v", err)
	}

	// 在 Start() 之前调用 Stop() 不应该 panic
	gaokaoBot.Stop()

	// 如果没有 panic，测试通过
}

// TestWait_BeforeStart 测试在 Start() 之前调用 Wait() 不会阻塞
func TestWait_BeforeStart(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	cfg := &config.TelegramConfig{
		Bot: config.BotConfig{
			Username: "test_bot",
			Token:    "test_token",
		},
		MiniApp: config.MiniAppConfig{
			URL: "https://example.com",
		},
	}

	gaokaoBot, err := NewGaokaoBot(nil, cfg, nil, logger)
	if err != nil {
		t.Fatalf("NewGaokaoBot() error = %v", err)
	}

	// 在 Start() 之前调用 Wait() 应该立即返回
	done := make(chan struct{})
	go func() {
		gaokaoBot.Wait()
		close(done)
	}()

	// 等待一小段时间，确保 Wait() 返回
	select {
	case <-done:
		// Wait() 正确返回，测试通过
	case <-time.After(1 * time.Second):
		t.Error("Wait() blocked when bot was not started")
	}
}

// guestSpyCaller 模拟 Telegram API 调用，捕获 Guest 查询应答，避免真实网络请求
type guestSpyCaller struct {
	called chan struct{}
}

func (c *guestSpyCaller) Call(_ context.Context, _ string, _ *telegoapi.RequestData) (*telegoapi.Response, error) {
	select {
	case c.called <- struct{}{}:
	default:
	}
	return &telegoapi.Response{
		Ok:     true,
		Result: json.RawMessage(`{"inline_message_id":"m1"}`),
	}, nil
}

// TestRegisterHandlers_GuestMessageRouted 验证 Guest 消息更新被路由到 BotService.HandleGuestMessage
func TestRegisterHandlers_GuestMessageRouted(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	if err := db.AutoMigrate(&model.ExamDate{}, &model.UserTemplate{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	examDateService := service.NewExamDateService(repository.NewExamDateRepository(db))
	userTemplateService := service.NewUserTemplateService(repository.NewUserTemplateRepository(db))
	messageService := service.NewMessageService(examDateService, userTemplateService, logger)

	caller := &guestSpyCaller{called: make(chan struct{}, 1)}
	tgBot, err := telego.NewBot(
		"123456:abcdefghijklmnopqrstuvwxyz012345678",
		telego.WithAPICaller(caller),
		telego.WithDiscardLogger(),
	)
	if err != nil {
		t.Fatalf("NewBot() error = %v", err)
	}

	botService := service.NewBotService(tgBot, messageService, nil, logger, "")

	cfg := &config.TelegramConfig{
		Bot:     config.BotConfig{Username: "gaokao_bot", Token: "test_token"},
		MiniApp: config.MiniAppConfig{URL: "https://example.com"},
	}
	gaokaoBot, err := NewGaokaoBot(tgBot, cfg, botService, logger)
	if err != nil {
		t.Fatalf("NewGaokaoBot() error = %v", err)
	}

	// 预置一个 Guest 消息更新，然后构造 handler 并注册处理器
	updates := make(chan telego.Update, 1)
	updates <- telego.Update{
		GuestMessage: &telego.Message{Text: "@gaokao_bot", GuestQueryID: "q1"},
	}
	handler, err := telegohandler.NewBotHandler(tgBot, updates)
	if err != nil {
		t.Fatalf("NewBotHandler() error = %v", err)
	}
	gaokaoBot.handler = handler

	gaokaoBot.registerHandlers()

	go func() { _ = handler.Start() }()
	defer func() { _ = handler.Stop() }()

	// Guest 处理器命中后会经 BotService.HandleGuestMessage 调用 answerGuestQuery
	select {
	case <-caller.called:
		// 路由成功，测试通过
	case <-time.After(2 * time.Second):
		t.Fatal("Guest 消息未被路由到 Guest 处理器")
	}
}
