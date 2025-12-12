package bot

import (
	"context"
	"sync"
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
	started bool              // 标记 bot 是否成功启动
	mu      sync.Mutex        // 保护并发访问
	ctx     context.Context   // 用于控制 bot 生命周期
	cancel  context.CancelFunc // 用于取消 context
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
	// 创建用于 GetMe 的临时 context
	getCtx, getCancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer getCancel()

	botUser, err := b.bot.GetMe(getCtx)
	if err != nil {
		return err
	}
	b.logger.Infof("Bot authorized on account %s", botUser.Username)

	// 创建用于 bot 生命周期的 context
	b.ctx, b.cancel = context.WithCancel(context.Background())

	// 启动长轮询（使用可取消的 context，确保 Stop() 时能优雅关闭）
	updates, err := b.bot.UpdatesViaLongPolling(b.ctx, nil)
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

	// 标记为已启动
	b.mu.Lock()
	b.started = true
	b.mu.Unlock()

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

	b.mu.Lock()
	started := b.started
	handler := b.handler
	cancel := b.cancel
	b.mu.Unlock()

	// 只有在成功启动后才尝试停止
	if started && handler != nil {
		// 取消 context，停止长轮询
		if cancel != nil {
			cancel()
		}

		if err := handler.Stop(); err != nil {
			b.logger.Errorf("Handler stop error: %v", err)
		}
		// 等待 handler 完全停止，使用 select 防止永久阻塞
		select {
		case <-b.done:
			b.logger.Info("Bot stopped")
		case <-time.After(5 * time.Second):
			b.logger.Warn("Timeout waiting for bot to stop")
		}
	} else {
		b.logger.Info("Bot was not started or already stopped")
	}
}

// Wait 等待Bot停止
func (b *GaokaoBot) Wait() {
	b.mu.Lock()
	started := b.started
	b.mu.Unlock()

	// 只有在成功启动后才等待
	if started {
		<-b.done
	}
}
