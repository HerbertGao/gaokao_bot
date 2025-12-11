package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/api"
	"github.com/herbertgao/gaokao_bot/internal/bot"
	"github.com/herbertgao/gaokao_bot/internal/config"
	"github.com/herbertgao/gaokao_bot/internal/database"
	"github.com/herbertgao/gaokao_bot/internal/repository"
	"github.com/herbertgao/gaokao_bot/internal/service"
	"github.com/herbertgao/gaokao_bot/internal/task"
	"github.com/herbertgao/gaokao_bot/internal/util"
	"github.com/herbertgao/gaokao_bot/internal/version"
	"github.com/mymmrac/telego"
	"github.com/sirupsen/logrus"
)

func main() {
	// 解析命令行参数
	env := flag.String("env", "dev", "Environment: dev, prod")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// 显示版本信息
	if *showVersion {
		fmt.Println(version.GetFullVersionInfo())
		return
	}

	// 加载配置
	cfg, err := config.Load(*env)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	logger := initLogger(cfg.Log)

	// 初始化数据库
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		logger.Fatalf("连接数据库失败: %v", err)
	}

	// 确保数据库连接在程序退出时关闭
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatalf("获取数据库实例失败: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Errorf("关闭数据库连接失败: %v", err)
		}
	}()

	// 初始化 Snowflake
	if err := util.InitSnowflake(cfg.Snowflake.DatacenterID, cfg.Snowflake.MachineID); err != nil {
		logger.Fatalf("初始化 Snowflake 失败: %v", err)
	}

	// 初始化仓储
	examDateRepo := repository.NewExamDateRepository(db)
	userTemplateRepo := repository.NewUserTemplateRepository(db)
	sendChatRepo := repository.NewSendChatRepository(db)

	// 初始化服务
	examDateService := service.NewExamDateService(examDateRepo)
	userTemplateService := service.NewUserTemplateService(userTemplateRepo)
	sendChatService := service.NewSendChatService(sendChatRepo)

	// 初始化 Telegram Bot
	var telegramBot *telego.Bot
	if cfg.App.Env == "prod" {
		// 生产环境：使用标准 API
		telegramBot, err = telego.NewBot(cfg.Telegram.Bot.Token)
	} else {
		// 非生产环境（dev/test）：使用测试服务器 API
		logger.Infof("使用 Telegram 测试服务器 (环境: %s)", cfg.App.Env)
		telegramBot, err = telego.NewBot(cfg.Telegram.Bot.Token, telego.WithTestServerPath())
	}
	if err != nil {
		logger.Fatalf("创建 Telegram Bot 失败: %v", err)
	}

	// 初始化消息和内联查询服务
	messageService := service.NewMessageService(examDateService, userTemplateService, logger)
	inlineQueryService := service.NewInlineQueryService(examDateService, userTemplateService, logger)

	// 初始化 Bot 服务
	botService := service.NewBotService(telegramBot, messageService, inlineQueryService, logger, cfg.Telegram.MiniApp.URL)

	// 初始化高考倒计时 Bot
	gaokaoBot, err := bot.NewGaokaoBot(telegramBot, &cfg.Telegram, botService, logger)
	if err != nil {
		logger.Fatalf("初始化 Bot 失败: %v", err)
	}

	// 初始化定时任务
	var dailyTask *task.DailySendTask
	if cfg.Task.DailySend.Enabled {
		dailyTask = task.NewDailySendTask(telegramBot, examDateService, userTemplateService, sendChatService, logger)
		if err := dailyTask.Start(cfg.Task.DailySend.Cron); err != nil {
			logger.Fatalf("启动定时任务失败: %v", err)
		}
	}

	// 初始化 HTTP 服务器（Mini App API）
	// 在非生产环境下跳过 Telegram Init Data 验证（方便开发调试）
	skipValidation := cfg.App.Env != "prod"
	// 仅在 debug 日志级别下启用 GIN 访问日志
	enableGinLogger := cfg.Log.Level == "debug"
	router, rateLimiter := api.NewRouter(db, cfg.Telegram.Bot.Token, userTemplateService, skipValidation, enableGinLogger)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: router,
	}

	// 创建错误通道用于 goroutine 错误传递
	serverErr := make(chan error, 1)

	// 在 goroutine 中启动 HTTP 服务器
	go func() {
		logger.Infof("正在启动 HTTP 服务器，端口: %d...", cfg.App.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("HTTP 服务器错误: %v", err)
			serverErr <- err
		}
	}()

	// 启动 Bot
	logger.Info("正在启动高考倒计时 Bot...")
	if err := gaokaoBot.Start(); err != nil {
		logger.Fatalf("Bot 错误: %v", err)
	}

	// 等待退出信号或服务器错误
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		logger.Info("收到退出信号")
	case err := <-serverErr:
		logger.Errorf("HTTP 服务器异常退出: %v", err)
	}

	logger.Info("正在关闭...")

	// 停止 HTTP 服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Errorf("HTTP 服务器关闭错误: %v", err)
	}

	// 停止速率限制器
	rateLimiter.Stop()

	// 停止 Bot
	gaokaoBot.Stop()

	// 停止定时任务
	if dailyTask != nil {
		dailyTask.Stop()
	}

	logger.Info("高考倒计时 Bot 已停止")
}

func initLogger(cfg config.LogConfig) *logrus.Logger {
	logger := logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 设置日志格式
	if cfg.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// 设置输出
	if cfg.Output == "stdout" {
		logger.SetOutput(os.Stdout)
	}

	return logger
}
