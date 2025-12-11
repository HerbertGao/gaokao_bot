package service

import (
	"context"
	"strings"
	"time"

	"github.com/herbertgao/gaokao_bot/pkg/constant"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

const (
	// DefaultContextTimeout 默认上下文超时时间
	DefaultContextTimeout = 10 * time.Second
)

// BotService Bot业务服务
type BotService struct {
	bot                *telego.Bot
	messageService     *MessageService
	inlineQueryService *InlineQueryService
	logger             *logrus.Logger
	miniAppURL         string
}

// NewBotService 创建Bot业务服务
func NewBotService(
	bot *telego.Bot,
	messageService *MessageService,
	inlineQueryService *InlineQueryService,
	logger *logrus.Logger,
	miniAppURL string,
) *BotService {
	return &BotService{
		bot:                bot,
		messageService:     messageService,
		inlineQueryService: inlineQueryService,
		logger:             logger,
		miniAppURL:         miniAppURL,
	}
}

// HandleMessage 处理消息
func (s *BotService) HandleMessage(bot *telego.Bot, msg *telego.Message) {
	if msg == nil {
		return
	}

	// 检查消息文本是否为空
	if msg.Text == "" {
		return
	}

	// 处理命令（检查是否以 / 开头）
	if strings.HasPrefix(msg.Text, "/") {
		s.handleCommand(msg)
		return
	}
}

// HandleInlineQuery 处理内联查询
func (s *BotService) HandleInlineQuery(bot *telego.Bot, query *telego.InlineQuery) {
	if query == nil {
		return
	}

	results := s.inlineQueryService.GetInlineQueryResults(query)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	err := s.bot.AnswerInlineQuery(ctx, &telego.AnswerInlineQueryParams{
		InlineQueryID: query.ID,
		Results:       results,
	})

	if err != nil {
		s.logger.Errorf("回复内联查询失败: %v", err)
	} else if s.logger.Level >= logrus.DebugLevel {
		// Debug 模式下打印内联查询回复
		s.logger.Debugf("[Telegram] -> Answered inline query (ID: %s) with %d results",
			query.ID,
			len(results))
	}
}

// handleCommand 处理命令
func (s *BotService) handleCommand(msg *telego.Message) {
	// 提取命令
	parts := strings.Fields(msg.Text)
	if len(parts) == 0 {
		return
	}

	cmd := strings.TrimPrefix(parts[0], "/")
	// 移除 @botname 部分
	if atIndex := strings.Index(cmd, "@"); atIndex != -1 {
		cmd = cmd[:atIndex]
	}

	var response string
	var err error

	switch cmd {
	case constant.CountdownCommand:
		response, err = s.messageService.GetCountDownMessage(msg)
	case constant.DebugCommand:
		s.handleDebugCommand(msg)
		return
	case constant.TemplateCommand:
		s.handleTemplateCommand(msg)
		return
	default:
		// 未知命令，忽略
		return
	}

	if err != nil {
		s.logger.Errorf("命令执行错误: %v", err)
		response = "处理命令时出错，请稍后重试"
	}

	// 发送回复
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	sentMsg, err := s.bot.SendMessage(ctx, telegoutil.Message(
		telegoutil.ID(msg.Chat.ID),
		response,
	))

	if err != nil {
		s.logger.Errorf("发送消息失败: %v", err)
	} else if s.logger.Level >= logrus.DebugLevel {
		// Debug 模式下打印发送的消息
		s.logger.Debugf("[Telegram] -> Sent message to Chat %d (MsgID: %d): %s",
			msg.Chat.ID,
			sentMsg.MessageID,
			truncateString(response, 100))
	}
}

// handleDebugCommand 处理 debug 命令
func (s *BotService) handleDebugCommand(msg *telego.Message) {
	// 构建带 debug=1 参数的 URL
	debugURL := s.miniAppURL + "?debug=1"

	// 创建 Web App 按钮
	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{
					Text:   "打开调试模式小程序",
					WebApp: &telego.WebAppInfo{URL: debugURL},
				},
			},
		},
	}

	// 发送带按钮的消息
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	sentMsg, err := s.bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      telegoutil.ID(msg.Chat.ID),
		Text:        "点击下方按钮打开调试模式的小程序",
		ReplyMarkup: keyboard,
	})

	if err != nil {
		s.logger.Errorf("发送调试消息失败: %v", err)
	} else if s.logger.Level >= logrus.DebugLevel {
		// Debug 模式下打印发送的消息
		s.logger.Debugf("[Telegram] -> Sent /debug response to Chat %d (MsgID: %d) with WebApp button",
			msg.Chat.ID,
			sentMsg.MessageID)
	}
}

// handleTemplateCommand 处理 template 命令
func (s *BotService) handleTemplateCommand(msg *telego.Message) {
	// 创建 Web App 按钮（使用小程序主页URL）
	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{
					Text:   "打开模板配置",
					WebApp: &telego.WebAppInfo{URL: s.miniAppURL},
				},
			},
		},
	}

	// 发送带按钮的消息
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	sentMsg, err := s.bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      telegoutil.ID(msg.Chat.ID),
		Text:        "点击下方按钮打开小程序配置你的自定义模板",
		ReplyMarkup: keyboard,
	})

	if err != nil {
		s.logger.Errorf("发送模板消息失败: %v", err)
	} else if s.logger.Level >= logrus.DebugLevel {
		// Debug 模式下打印发送的消息
		s.logger.Debugf("[Telegram] -> Sent /template response to Chat %d (MsgID: %d) with WebApp button",
			msg.Chat.ID,
			sentMsg.MessageID)
	}
}

// truncateString 截断字符串到指定长度
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
