package bot

import (
	"context"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/config"
	"github.com/herbertgao/gaokao_bot/internal/service"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/sirupsen/logrus"
)

const (
	// DefaultContextTimeout 默认上下文超时时间
	DefaultContextTimeout = 10 * time.Second
)

// GaokaoBot 高考Bot
type GaokaoBot struct {
	bot     *telego.Bot
	config  *config.TelegramConfig
	service *service.BotService
	logger  *logrus.Logger
	handler *telegohandler.BotHandler
	updates <-chan telego.Update
	done    chan struct{}
}

// NewGaokaoBot 创建高考Bot
func NewGaokaoBot(
	bot *telego.Bot,
	cfg *config.TelegramConfig,
	svc *service.BotService,
	logger *logrus.Logger,
) (*GaokaoBot, error) {
	return &GaokaoBot{
		bot:     bot,
		config:  cfg,
		service: svc,
		logger:  logger,
		done:    make(chan struct{}),
	}, nil
}

// Start 启动Bot
func (b *GaokaoBot) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	botUser, err := b.bot.GetMe(ctx)
	if err != nil {
		return err
	}
	b.logger.Infof("Bot authorized on account %s", botUser.Username)

	// 启动长轮询
	updates, err := b.bot.UpdatesViaLongPolling(context.TODO(), nil)
	if err != nil {
		return err
	}
	b.updates = updates

	// 创建 handler
	handler, err := telegohandler.NewBotHandler(b.bot, updates)
	if err != nil {
		return err
	}
	b.handler = handler

	// 注册消息处理器
	b.handler.Handle(func(ctx *telegohandler.Context, update telego.Update) error {
		if update.Message != nil {
			// Debug 模式下打印接收到的消息
			if b.logger.Level >= logrus.DebugLevel {
				var username string
				var userID int64
				if update.Message.From != nil {
					username = update.Message.From.Username
					userID = update.Message.From.ID
				} else {
					username = "unknown"
					userID = 0
				}
				b.logger.Debugf("[Telegram] <- Received message from @%s (ID: %d, Chat: %d): %s",
					username,
					userID,
					update.Message.Chat.ID,
					update.Message.Text)
			}
			b.service.HandleMessage(ctx.Bot(), update.Message)
		}
		return nil
	}, telegohandler.AnyMessage())

	// 注册内联查询处理器
	b.handler.Handle(func(ctx *telegohandler.Context, update telego.Update) error {
		if update.InlineQuery != nil {
			// Debug 模式下打印接收到的内联查询
			if b.logger.Level >= logrus.DebugLevel {
				var username string
				var userID int64
				if update.InlineQuery.From.ID != 0 {
					username = update.InlineQuery.From.Username
					userID = update.InlineQuery.From.ID
				} else {
					username = "unknown"
					userID = 0
				}
				b.logger.Debugf("[Telegram] <- Received inline query from @%s (ID: %d): %s",
					username,
					userID,
					update.InlineQuery.Query)
			}
			b.service.HandleInlineQuery(ctx.Bot(), update.InlineQuery)
		}
		return nil
	}, telegohandler.AnyInlineQuery())

	// 开始处理更新
	go func() {
		if err := b.handler.Start(); err != nil {
			b.logger.Errorf("Handler start error: %v", err)
		}
		close(b.done)
	}()

	b.logger.Info("Bot started successfully")

	return nil
}

// Stop 停止Bot
func (b *GaokaoBot) Stop() {
	b.logger.Info("Stopping bot...")
	if err := b.handler.Stop(); err != nil {
		b.logger.Errorf("Handler stop error: %v", err)
	}
	<-b.done // 等待 handler 完全停止
	b.logger.Info("Bot stopped")
}

// Wait 等待Bot停止
func (b *GaokaoBot) Wait() {
	<-b.done
}
