package service

import (
	"strconv"
	"strings"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/util"
	"github.com/mymmrac/telego"
	"github.com/sirupsen/logrus"
)

// MessageService 消息处理服务
type MessageService struct {
	examDateService     *ExamDateService
	userTemplateService *UserTemplateService
	logger              *logrus.Logger
}

// NewMessageService 创建消息处理服务
func NewMessageService(
	examDateService *ExamDateService,
	userTemplateService *UserTemplateService,
	logger *logrus.Logger,
) *MessageService {
	return &MessageService{
		examDateService:     examDateService,
		userTemplateService: userTemplateService,
		logger:              logger,
	}
}

// GetCountDownMessage 获取倒计时消息
func (s *MessageService) GetCountDownMessage(msg *telego.Message) (string, error) {
	now := time.Now()
	text := util.GetTextByMessage(msg)

	var examList []model.ExamDate
	var err error

	// 如果有参数且是数字，按年份查询
	if text != "" {
		if y, err := strconv.Atoi(text); err == nil && y >= 2018 && y <= 2100 {
			examList, err = s.examDateService.GetExamByYear(y)
			if err != nil {
				s.logger.Errorf("Failed to get exam by year %d: %v", y, err)
				return "查询考试信息失败", err
			}
			if len(examList) == 0 {
				return "参数暂时无法识别。", nil
			}
		} else {
			return "参数暂时无法识别。", nil
		}
	} else {
		// 没有参数时，获取当前时间范围内的所有考试
		examList, err = s.examDateService.GetExamsInRange(now)
		if err != nil {
			s.logger.Errorf("Failed to get exams in range: %v", err)
			return "查询考试信息失败", err
		}
	}

	// 如果没有找到任何考试
	if len(examList) == 0 {
		return "数据库中没有可用的信息，请联系开发者。", nil
	}

	// 获取默认模板
	template, err := s.userTemplateService.GetDefaultTemplate()
	if err != nil {
		s.logger.Errorf("Failed to get default template: %v", err)
		return "获取模板失败", err
	}

	templateContent := "现在距离{exam}还有{time}"
	if template != nil {
		templateContent = template.TemplateContent
	}

	// 生成倒计时消息（循环处理所有考试）
	var sb strings.Builder
	for _, exam := range examList {
		message := util.GetCountDownString(&exam, templateContent, now)
		sb.WriteString(message)
	}

	return sb.String(), nil
}
