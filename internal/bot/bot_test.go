package bot

import (
	"testing"

	"github.com/herbertgao/gaokao_bot/internal/config"
	"github.com/herbertgao/gaokao_bot/internal/service"
	"github.com/mymmrac/telego"
	"github.com/sirupsen/logrus"
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
}
