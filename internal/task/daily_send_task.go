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

const (
	// DefaultContextTimeout 默认上下文超时时间
	DefaultContextTimeout = 10 * time.Second
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
	t.logger.Info("每日发送任务已启动")
	return nil
}

// Stop 停止定时任务
func (t *DailySendTask) Stop() {
	t.cron.Stop()
	t.logger.Info("每日发送任务已停止")
}

// execute 执行任务
func (t *DailySendTask) execute() {
	// 定时任务一般设置在整点，可能有±1分钟以内的误差
	// 将时间标准化为最近的整分钟，防止出现"还剩3天23小时59分钟59秒"或"还剩4天1秒"的情况
	// 四舍五入：秒 >= 30 则进位到下一分钟，< 30 则保持当前分钟
	now := time.Now()
	minute := now.Minute()
	hour := now.Hour()

	if now.Second() >= 30 {
		// 进位到下一分钟
		now = now.Add(time.Minute)
		minute = now.Minute()
		hour = now.Hour()
	}

	now = time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())

	// 获取符合条件的考试
	exams, err := t.examDateService.GetExamsInRange(now)
	if err != nil {
		t.logger.Errorf("获取考试列表失败: %v", err)
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
			t.logger.Errorf("获取默认模板失败: %v", err)
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
			t.logger.Errorf("获取聊天列表失败: %v", err)
			continue
		}

		// 发送消息
		for _, chat := range chats {
			chatID, err := strconv.ParseInt(chat.ChatID, 10, 64)
			if err != nil {
				t.logger.Errorf("无效的聊天ID %s: %v", chat.ChatID, err)
				continue
			}

			// 使用带超时的 context 防止 API 调用挂起
			ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
			sentMsg, err := t.bot.SendMessage(ctx, telegoutil.Message(
				telegoutil.ID(chatID),
				message,
			))
			cancel()

			if err != nil {
				t.logger.Errorf("发送消息到聊天 %s 失败: %v", chat.ChatID, err)
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
