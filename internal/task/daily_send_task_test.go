package task

import (
	"strings"
	"testing"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/util"
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

	// 使用北京时区（与生产代码保持一致）
	bjtZone := util.GetBJTLocation()

	tests := []struct {
		name        string
		examDate    time.Time
		currentTime time.Time
		want        bool
	}{
		{
			name:        "24小时内，应该发送",
			examDate:    time.Date(2025, 6, 7, 10, 0, 0, 0, bjtZone),
			currentTime: time.Date(2025, 6, 7, 1, 0, 0, 0, bjtZone), // 距离9小时
			want:        true,
		},
		{
			name:        "12小时内，应该发送",
			examDate:    time.Date(2025, 6, 7, 21, 0, 0, 0, bjtZone),
			currentTime: time.Date(2025, 6, 7, 10, 0, 0, 0, bjtZone), // 距离11小时
			want:        true,
		},
		{
			name:        "1小时内，应该发送",
			examDate:    time.Date(2025, 6, 7, 10, 30, 0, 0, bjtZone),
			currentTime: time.Date(2025, 6, 7, 10, 0, 0, 0, bjtZone), // 距离30分钟
			want:        true,
		},
		{
			name:        "超过24小时，9:00应该发送",
			examDate:    time.Date(2025, 6, 10, 10, 0, 0, 0, bjtZone),
			currentTime: time.Date(2025, 6, 7, 9, 0, 0, 0, bjtZone), // 距离3天，当前时间9:00
			want:        true,
		},
		{
			name:        "超过24小时，非9:00不发送",
			examDate:    time.Date(2025, 6, 10, 10, 0, 0, 0, bjtZone),
			currentTime: time.Date(2025, 6, 7, 10, 0, 0, 0, bjtZone), // 距离3天，当前时间10:00
			want:        false,
		},
		{
			name:        "超过24小时，9:01应该发送（时间窗口容错）",
			examDate:    time.Date(2025, 6, 10, 10, 0, 0, 0, bjtZone),
			currentTime: time.Date(2025, 6, 7, 9, 1, 0, 0, bjtZone), // 距离3天，当前时间9:01
			want:        true,
		},
		{
			name:        "超过24小时，9:02不发送（超出时间窗口）",
			examDate:    time.Date(2025, 6, 10, 10, 0, 0, 0, bjtZone),
			currentTime: time.Date(2025, 6, 7, 9, 2, 0, 0, bjtZone), // 距离3天，当前时间9:02
			want:        false,
		},
		{
			name:        "考试已过，不发送",
			examDate:    time.Date(2025, 6, 7, 10, 0, 0, 0, bjtZone),
			currentTime: time.Date(2025, 6, 7, 11, 0, 0, 0, bjtZone), // 考试过去1小时
			want:        false,
		},
		{
			name:        "刚好24小时，应该发送",
			examDate:    time.Date(2025, 6, 8, 10, 0, 0, 0, bjtZone),
			currentTime: time.Date(2025, 6, 7, 10, 0, 0, 0, bjtZone), // 刚好24小时
			want:        true,
		},
		{
			name:        "刚好超过24小时，非9:00不发送",
			examDate:    time.Date(2025, 6, 8, 10, 0, 0, 0, bjtZone),
			currentTime: time.Date(2025, 6, 7, 9, 59, 0, 0, bjtZone), // 略超过24小时
			want:        false,
		},
		{
			name:        "超过24小时，9:00:35应该发送（Bug修复：时间标准化不影响判断）",
			examDate:    time.Date(2025, 6, 10, 10, 0, 0, 0, bjtZone),
			currentTime: time.Date(2025, 6, 7, 9, 0, 35, 0, bjtZone), // 9:00:35，秒>=30会标准化为9:01
			want:        true, // 应该发送，因为 shouldSend 使用原始时间判断
		},
		{
			name:        "超过24小时，9:00:59应该发送（Bug修复：时间标准化不影响判断）",
			examDate:    time.Date(2025, 6, 10, 10, 0, 0, 0, bjtZone),
			currentTime: time.Date(2025, 6, 7, 9, 0, 59, 0, bjtZone), // 9:00:59，秒>=30会标准化为9:01
			want:        true, // 应该发送，因为 shouldSend 使用原始时间判断
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

// TestTimeNormalization 测试时间标准化逻辑
// 定时任务一般设置在整点，可能有±1分钟以内的误差
// 将时间标准化为最近的整分钟（基于北京时间 UTC+8）
// 四舍五入：秒 >= 30 进位到下一分钟，< 30 保持当前分钟
func TestTimeNormalization(t *testing.T) {
	bjtZone := time.FixedZone("BJT", 8*3600) // UTC+8

	tests := []struct {
		name           string
		originalTime   time.Time
		normalizedTime time.Time
	}{
		{
			name:           "9:00:01 触发应标准化为 9:00:00（< 30秒，保持当前分钟）",
			originalTime:   time.Date(2025, 6, 7, 9, 0, 1, 0, bjtZone),
			normalizedTime: time.Date(2025, 6, 7, 9, 0, 0, 0, bjtZone),
		},
		{
			name:           "9:00:29 触发应标准化为 9:00:00（< 30秒，保持当前分钟）",
			originalTime:   time.Date(2025, 6, 7, 9, 0, 29, 0, bjtZone),
			normalizedTime: time.Date(2025, 6, 7, 9, 0, 0, 0, bjtZone),
		},
		{
			name:           "9:00:30 触发应标准化为 9:01:00（= 30秒，进位到下一分钟）",
			originalTime:   time.Date(2025, 6, 7, 9, 0, 30, 0, bjtZone),
			normalizedTime: time.Date(2025, 6, 7, 9, 1, 0, 0, bjtZone),
		},
		{
			name:           "9:00:59 触发应标准化为 9:01:00（>= 30秒，进位到下一分钟）",
			originalTime:   time.Date(2025, 6, 7, 9, 0, 59, 0, bjtZone),
			normalizedTime: time.Date(2025, 6, 7, 9, 1, 0, 0, bjtZone),
		},
		{
			name:           "8:59:29 触发应标准化为 8:59:00（< 30秒，保持当前分钟）",
			originalTime:   time.Date(2025, 6, 7, 8, 59, 29, 0, bjtZone),
			normalizedTime: time.Date(2025, 6, 7, 8, 59, 0, 0, bjtZone),
		},
		{
			name:           "8:59:30 触发应标准化为 9:00:00（>= 30秒，进位到下一分钟）",
			originalTime:   time.Date(2025, 6, 7, 8, 59, 30, 0, bjtZone),
			normalizedTime: time.Date(2025, 6, 7, 9, 0, 0, 0, bjtZone),
		},
		{
			name:           "9:00:00 触发应保持 9:00:00",
			originalTime:   time.Date(2025, 6, 7, 9, 0, 0, 0, bjtZone),
			normalizedTime: time.Date(2025, 6, 7, 9, 0, 0, 0, bjtZone),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用 util.NormalizeToMinute
			normalized := util.NormalizeToMinute(tt.originalTime)

			if !normalized.Equal(tt.normalizedTime) {
				t.Errorf("标准化时间 = %v, 期望 %v", normalized.Format("2006-01-02 15:04:05"), tt.normalizedTime.Format("2006-01-02 15:04:05"))
			}
		})
	}
}

// TestCountdownWithTimeNormalization 测试标准化时间后的倒计时显示（基于北京时间 UTC+8）
func TestCountdownWithTimeNormalization(t *testing.T) {
	bjtZone := time.FixedZone("BJT", 8*3600) // UTC+8

	tests := []struct {
		name            string
		triggerTime     time.Time // 实际触发时间
		examDate        time.Time // 考试时间
		wantContains    string    // 期望倒计时包含的内容
		wantNotContains string    // 不应该包含的内容（如 "59分钟59秒"）
	}{
		{
			name:            "9:00:01 触发，距离考试4天，应显示4天而不是3天23小时59分59秒",
			triggerTime:     time.Date(2025, 6, 3, 9, 0, 1, 0, bjtZone),
			examDate:        time.Date(2025, 6, 7, 9, 0, 0, 0, bjtZone),
			wantContains:    "4天",
			wantNotContains: "59",
		},
		{
			name:            "8:59:59 触发，标准化为 9:00:00，距离考试刚好5天",
			triggerTime:     time.Date(2025, 6, 2, 8, 59, 59, 0, bjtZone),
			examDate:        time.Date(2025, 6, 7, 9, 0, 0, 0, bjtZone),
			wantContains:    "5天",
			wantNotContains: "小时",
		},
		{
			name:            "9:00:00 触发，距离考试刚好4天，应显示4天",
			triggerTime:     time.Date(2025, 6, 3, 9, 0, 0, 0, bjtZone),
			examDate:        time.Date(2025, 6, 7, 9, 0, 0, 0, bjtZone),
			wantContains:    "4天",
			wantNotContains: "小时",
		},
		{
			name:            "10:30:45 触发，标准化为 10:31:00，距离考试3天22小时29分钟",
			triggerTime:     time.Date(2025, 6, 3, 10, 30, 45, 123456789, bjtZone),
			examDate:        time.Date(2025, 6, 7, 9, 0, 0, 0, bjtZone),
			wantContains:    "3天22小时29分钟",
			wantNotContains: "秒",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用 util.NormalizeToMinute（与生产代码一致）
			normalized := util.NormalizeToMinute(tt.triggerTime)

			// 使用 util.FormatDuration 计算倒计时（与生产代码一致）
			duration := tt.examDate.Sub(normalized)
			result := util.FormatDuration(duration)

			t.Logf("原始触发时间: %v", tt.triggerTime.Format("2006-01-02 15:04:05 MST"))
			t.Logf("标准化时间: %v", normalized.Format("2006-01-02 15:04:05 MST"))
			t.Logf("倒计时结果: %s", result)

			if !strings.Contains(result, tt.wantContains) {
				t.Errorf("倒计时 = %v, 应该包含 %v", result, tt.wantContains)
			}

			if tt.wantNotContains != "" && strings.Contains(result, tt.wantNotContains) {
				t.Errorf("倒计时 = %v, 不应该包含 %v", result, tt.wantNotContains)
			}
		})
	}
}
