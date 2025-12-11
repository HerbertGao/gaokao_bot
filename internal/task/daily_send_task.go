package task

import (
	"context"
	"strconv"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/service"
	"github.com/herbertgao/gaokao_bot/internal/util"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// DailySendTask 每日发送任务
type DailySendTask struct {
	cron                *cron.Cron
	bot                 *telego.Bot
	examDateService     *service.ExamDateService
	userTemplateService *service.UserTemplateService
	sendChatService     *service.SendChatService
	logger              *logrus.Logger
}

// NewDailySendTask 创建每日发送任务
func NewDailySendTask(
	bot *telego.Bot,
	examDateService *service.ExamDateService,
	userTemplateService *service.UserTemplateService,
	sendChatService *service.SendChatService,
	logger *logrus.Logger,
) *DailySendTask {
	return &DailySendTask{
		cron:                cron.New(cron.WithSeconds()),
		bot:                 bot,
		examDateService:     examDateService,
		userTemplateService: userTemplateService,
		sendChatService:     sendChatService,
		logger:              logger,
	}
}

// Start 启动定时任务
func (t *DailySendTask) Start(cronExpr string) error {
	_, err := t.cron.AddFunc(cronExpr, t.execute)
	if err != nil {
		return err
	}

	t.cron.Start()
	t.logger.Info("Daily send task started")
	return nil
}

// Stop 停止定时任务
func (t *DailySendTask) Stop() {
	t.cron.Stop()
	t.logger.Info("Daily send task stopped")
}

// execute 执行任务
func (t *DailySendTask) execute() {
	now := time.Now()

	// 获取符合条件的考试
	exams, err := t.examDateService.GetExamsInRange(now)
	if err != nil {
		t.logger.Errorf("Failed to get exams: %v", err)
		return
	}

	if len(exams) == 0 {
		return
	}

	for _, exam := range exams {
		if !t.shouldSend(exam, now) {
			continue
		}

		// 获取默认模板
		template, err := t.userTemplateService.GetDefaultTemplate()
		if err != nil {
			t.logger.Errorf("Failed to get default template: %v", err)
			continue
		}

		templateContent := "现在距离{exam}还有{time}"
		if template != nil {
			templateContent = template.TemplateContent
		}

		// 生成倒计时消息
		message := util.GetCountDownString(&exam, templateContent, now)

		// 获取发送目标
		chats, err := t.sendChatService.GetAll()
		if err != nil {
			t.logger.Errorf("Failed to get chat list: %v", err)
			continue
		}

		// 发送消息
		for _, chat := range chats {
			chatID, err := strconv.ParseInt(chat.ChatID, 10, 64)
			if err != nil {
				t.logger.Errorf("Invalid chat ID %s: %v", chat.ChatID, err)
				continue
			}

			sentMsg, err := t.bot.SendMessage(context.Background(), telegoutil.Message(
				telegoutil.ID(chatID),
				message,
			))

			if err != nil {
				t.logger.Errorf("Failed to send to chat %s: %v", chat.ChatID, err)
			} else if t.logger.Level >= logrus.DebugLevel {
				// Debug 模式下打印发送的消息
				t.logger.Debugf("[Telegram] -> Sent daily task message to Chat %d (MsgID: %d)",
					chatID,
					sentMsg.MessageID)
			}
		}
	}
}

// shouldSend 判断是否应该发送
func (t *DailySendTask) shouldSend(exam model.ExamDate, now time.Time) bool {
	duration := exam.ExamBeginDate.Sub(now)
	hours := duration.Hours()

	// 距离考试 <= 24 小时，每小时发送
	if hours <= 24 && hours > 0 {
		return true
	}

	// 距离考试 > 1 天，仅在9:00发送
	if hours > 24 && now.Hour() == 9 && now.Minute() == 0 {
		return true
	}

	return false
}
