package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/model"
)

// IsExamBeginTime 判断是否考试开始时刻
func IsExamBeginTime(exam *model.ExamDate, now time.Time) bool {
	return now.After(exam.ExamBeginDate) && now.Before(exam.ExamBeginDate.Add(time.Minute))
}

// IsExamTime 判断是否考试进行中
func IsExamTime(exam *model.ExamDate, now time.Time) bool {
	return now.After(exam.ExamBeginDate) && now.Before(exam.ExamEndDate)
}

// IsExpiredExam 判断考试是否已结束
func IsExpiredExam(exam *model.ExamDate, now time.Time) bool {
	return now.After(exam.ExamEndDate)
}

// GetCountDownString 生成倒计时字符串
func GetCountDownString(exam *model.ExamDate, template string, now time.Time) string {
	if IsExamTime(exam, now) {
		return fmt.Sprintf("%s正在进行中！", exam.ExamDesc)
	}

	if IsExpiredExam(exam, now) {
		return fmt.Sprintf("%s已经结束了。", exam.ExamDesc)
	}

	duration := exam.ExamBeginDate.Sub(now)
	countdownTime := FormatDuration(duration)

	result := template
	result = strings.ReplaceAll(result, "{exam_year}", fmt.Sprintf("%d", exam.ExamYear))
	result = strings.ReplaceAll(result, "{exam}", exam.ExamDesc)
	result = strings.ReplaceAll(result, "{exam_s}", exam.ShortDesc)
	result = strings.ReplaceAll(result, "{time}", countdownTime)

	return result
}

// GetCountDownTime 获取倒计时时间文字
func GetCountDownTime(exam *model.ExamDate, now time.Time) string {
	if IsExamTime(exam, now) || IsExpiredExam(exam, now) {
		return ""
	}

	duration := exam.ExamBeginDate.Sub(now)
	return FormatDuration(duration)
}