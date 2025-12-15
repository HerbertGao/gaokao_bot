package service

import (
	"fmt"
	"strconv"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/util"
	"github.com/herbertgao/gaokao_bot/pkg/constant"
	"github.com/mymmrac/telego"
	"github.com/sirupsen/logrus"
)

// InlineQueryService 内联查询服务
type InlineQueryService struct {
	examDateService     *ExamDateService
	userTemplateService *UserTemplateService
	logger              *logrus.Logger
}

// NewInlineQueryService 创建内联查询服务
func NewInlineQueryService(
	examDateService *ExamDateService,
	userTemplateService *UserTemplateService,
	logger *logrus.Logger,
) *InlineQueryService {
	return &InlineQueryService{
		examDateService:     examDateService,
		userTemplateService: userTemplateService,
		logger:              logger,
	}
}

// GetInlineQueryResults 获取内联查询结果
func (s *InlineQueryService) GetInlineQueryResults(query *telego.InlineQuery) []telego.InlineQueryResult {
	now := util.NowBJT()

	var examList []model.ExamDate
	var err error

	// 解析年份
	if query.Query != "" {
		if y, err := strconv.Atoi(query.Query); err == nil && y >= constant.MinExamYear && y <= constant.MaxExamYear {
			examList, err = s.examDateService.GetExamByYear(y)
			if err != nil {
				s.logger.Errorf("按年份 %d 查询考试失败: %v", y, err)
				return []telego.InlineQueryResult{}
			}
		}
	} else {
		// 没有参数时，获取当前时间范围内的所有考试
		examList, err = s.examDateService.GetExamsInRange(now)
		if err != nil {
			s.logger.Errorf("查询时间范围内的考试失败: %v", err)
			return []telego.InlineQueryResult{}
		}
	}

	// 如果没有找到任何考试
	if len(examList) == 0 {
		return []telego.InlineQueryResult{}
	}

	// 获取默认模板
	defaultTemplate, err := s.userTemplateService.GetDefaultTemplate()
	if err != nil {
		s.logger.Errorf("获取默认模板失败: %v", err)
		return []telego.InlineQueryResult{}
	}

	// 获取用户自定义模板
	var userTemplates []model.UserTemplate
	if query.From.ID != 0 {
		userTemplates, _ = s.userTemplateService.GetByUserID(query.From.ID)
	}

	results := []telego.InlineQueryResult{}

	// 为每个考试生成inline结果
	for idx, exam := range examList {
		examDesc := exam.ShortDesc

		// 默认模板结果
		if defaultTemplate != nil {
			defaultTitle := fmt.Sprintf("查看%s倒计时", examDesc)
			defaultMessage := util.GetCountDownString(&exam, defaultTemplate.TemplateContent, now)
			result := &telego.InlineQueryResultArticle{
				Type:  telego.ResultTypeArticle,
				ID:    fmt.Sprintf("default_%d", idx),
				Title: defaultTitle,
				InputMessageContent: &telego.InputTextMessageContent{
					MessageText: defaultMessage,
				},
			}
			results = append(results, result)
		}

		// 用户自定义模板结果
		for tidx, template := range userTemplates {
			title := fmt.Sprintf("查看%s倒计时", examDesc)
			if template.TemplateName != "" {
				title = fmt.Sprintf("%s (%s)", title, template.TemplateName)
			}
			message := util.GetCountDownString(&exam, template.TemplateContent, now)
			result := &telego.InlineQueryResultArticle{
				Type:  telego.ResultTypeArticle,
				ID:    fmt.Sprintf("user_%d_%d", idx, tidx),
				Title: title,
				InputMessageContent: &telego.InputTextMessageContent{
					MessageText: message,
				},
			}
			results = append(results, result)
		}
	}

	return results
}
