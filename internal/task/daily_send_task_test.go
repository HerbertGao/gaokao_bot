package task

import (
	"testing"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/sirupsen/logrus"
)

func TestNewDailySendTask(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	task := NewDailySendTask(nil, nil, nil, nil, logger)

	if task == nil {
		t.Fatal("NewDailySendTask() returned nil")
	}

	if task.cron == nil {
		t.Error("cron not initialized")
	}

	if task.logger != logger {
		t.Error("logger not set correctly")
	}
}

func TestDailySendTask_ShouldSend(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	task := NewDailySendTask(nil, nil, nil, nil, logger)

	tests := []struct {
		name        string
		examDate    time.Time
		currentTime time.Time
		want        bool
	}{
		{
			name:        "24小时内，应该发送",
			examDate:    time.Date(2025, 6, 7, 10, 0, 0, 0, time.UTC),
			currentTime: time.Date(2025, 6, 7, 1, 0, 0, 0, time.UTC), // 距离9小时
			want:        true,
		},
		{
			name:        "12小时内，应该发送",
			examDate:    time.Date(2025, 6, 7, 21, 0, 0, 0, time.UTC),
			currentTime: time.Date(2025, 6, 7, 10, 0, 0, 0, time.UTC), // 距离11小时
			want:        true,
		},
		{
			name:        "1小时内，应该发送",
			examDate:    time.Date(2025, 6, 7, 10, 30, 0, 0, time.UTC),
			currentTime: time.Date(2025, 6, 7, 10, 0, 0, 0, time.UTC), // 距离30分钟
			want:        true,
		},
		{
			name:        "超过24小时，9:00应该发送",
			examDate:    time.Date(2025, 6, 10, 10, 0, 0, 0, time.UTC),
			currentTime: time.Date(2025, 6, 7, 9, 0, 0, 0, time.UTC), // 距离3天，当前时间9:00
			want:        true,
		},
		{
			name:        "超过24小时，非9:00不发送",
			examDate:    time.Date(2025, 6, 10, 10, 0, 0, 0, time.UTC),
			currentTime: time.Date(2025, 6, 7, 10, 0, 0, 0, time.UTC), // 距离3天，当前时间10:00
			want:        false,
		},
		{
			name:        "超过24小时，9:01不发送",
			examDate:    time.Date(2025, 6, 10, 10, 0, 0, 0, time.UTC),
			currentTime: time.Date(2025, 6, 7, 9, 1, 0, 0, time.UTC), // 距离3天，当前时间9:01
			want:        false,
		},
		{
			name:        "考试已过，不发送",
			examDate:    time.Date(2025, 6, 7, 10, 0, 0, 0, time.UTC),
			currentTime: time.Date(2025, 6, 7, 11, 0, 0, 0, time.UTC), // 考试过去1小时
			want:        false,
		},
		{
			name:        "刚好24小时，应该发送",
			examDate:    time.Date(2025, 6, 8, 10, 0, 0, 0, time.UTC),
			currentTime: time.Date(2025, 6, 7, 10, 0, 0, 0, time.UTC), // 刚好24小时
			want:        true,
		},
		{
			name:        "刚好超过24小时，非9:00不发送",
			examDate:    time.Date(2025, 6, 8, 10, 0, 0, 0, time.UTC),
			currentTime: time.Date(2025, 6, 7, 9, 59, 0, 0, time.UTC), // 略超过24小时
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exam := model.ExamDate{
				ExamBeginDate: tt.examDate,
			}

			got := task.shouldSend(exam, tt.currentTime)
			if got != tt.want {
				duration := tt.examDate.Sub(tt.currentTime)
				t.Errorf("shouldSend() = %v, want %v (距离: %.2f 小时, 当前时间: %s)",
					got, tt.want, duration.Hours(), tt.currentTime.Format("15:04"))
			}
		})
	}
}

func TestDailySendTask_Stop(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	task := NewDailySendTask(nil, nil, nil, nil, logger)

	// 测试 Stop 不会 panic
	task.Stop()

	// 可以多次调用 Stop
	task.Stop()
}

func TestDailySendTask_StartInvalidCron(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	task := NewDailySendTask(nil, nil, nil, nil, logger)

	// 使用无效的 cron 表达式
	err := task.Start("invalid cron expression")
	if err == nil {
		t.Error("Start() with invalid cron should return error")
	}
}

func TestDailySendTask_StartAndStop(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	task := NewDailySendTask(nil, nil, nil, nil, logger)

	// 使用有效但不会立即触发的 cron 表达式（每年1月1日0:00）
	// 格式: 秒 分 时 日 月 周
	err := task.Start("0 0 0 1 1 *")
	if err != nil {
		t.Errorf("Start() error = %v, want nil", err)
	}

	// 确保可以正常停止
	task.Stop()
}
