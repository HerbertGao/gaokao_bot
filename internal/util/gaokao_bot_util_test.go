package util

import (
	"strings"
	"testing"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/model"
)

// createTestExam 创建测试用的考试数据（匹配数据库真实数据）
func createTestExam() *model.ExamDate {
	// 2026年6月7日9:00 - 6月10日17:00（匹配实际数据库）
	examBeginDate := time.Date(2026, 6, 7, 9, 0, 0, 0, time.UTC)
	examEndDate := time.Date(2026, 6, 10, 17, 0, 0, 0, time.UTC)
	yearBeginDate := time.Date(2025, 6, 10, 17, 0, 0, 0, time.UTC)
	yearEndDate := time.Date(2026, 6, 10, 17, 0, 0, 0, time.UTC)

	return &model.ExamDate{
		ID:                9,
		ExamYear:          2026,
		ExamDesc:          "2026年普通高等学校招生全国统一考试",
		ShortDesc:         "2026年高考",
		ExamBeginDate:     examBeginDate,
		ExamEndDate:       examEndDate,
		ExamYearBeginDate: yearBeginDate,
		ExamYearEndDate:   yearEndDate,
		IsDelete:          false,
	}
}

func TestIsExamBeginTime(t *testing.T) {
	exam := createTestExam()
	// 考试开始时间: 2026-06-07 09:00:00

	tests := []struct {
		name     string
		now      time.Time
		expected bool
	}{
		{
			name:     "考试开始前1分钟",
			now:      exam.ExamBeginDate.Add(-1 * time.Minute),
			expected: false,
		},
		{
			name:     "考试开始时刻（0秒）",
			now:      exam.ExamBeginDate,
			expected: false, // After 不包括等于
		},
		{
			name:     "考试开始后10秒",
			now:      exam.ExamBeginDate.Add(10 * time.Second),
			expected: true,
		},
		{
			name:     "考试开始后30秒",
			now:      exam.ExamBeginDate.Add(30 * time.Second),
			expected: true,
		},
		{
			name:     "考试开始后59秒",
			now:      exam.ExamBeginDate.Add(59 * time.Second),
			expected: true,
		},
		{
			name:     "考试开始后1分钟",
			now:      exam.ExamBeginDate.Add(1 * time.Minute),
			expected: false, // Before 不包括等于
		},
		{
			name:     "考试开始后2分钟",
			now:      exam.ExamBeginDate.Add(2 * time.Minute),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsExamBeginTime(exam, tt.now)
			if result != tt.expected {
				t.Errorf("IsExamBeginTime() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsExamTime(t *testing.T) {
	exam := createTestExam()
	// 考试时间: 2026-06-07 09:00:00 ~ 2026-06-10 17:00:00

	tests := []struct {
		name     string
		now      time.Time
		expected bool
	}{
		{
			name:     "考试开始前1天",
			now:      exam.ExamBeginDate.Add(-24 * time.Hour),
			expected: false,
		},
		{
			name:     "考试开始前1分钟",
			now:      exam.ExamBeginDate.Add(-1 * time.Minute),
			expected: false,
		},
		{
			name:     "考试开始时刻",
			now:      exam.ExamBeginDate,
			expected: false, // After 不包括等于
		},
		{
			name:     "考试开始后1秒",
			now:      exam.ExamBeginDate.Add(1 * time.Second),
			expected: true,
		},
		{
			name:     "考试进行中（第一天下午）",
			now:      time.Date(2026, 6, 7, 15, 30, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "考试进行中（第二天）",
			now:      time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "考试结束前1分钟",
			now:      exam.ExamEndDate.Add(-1 * time.Minute),
			expected: true,
		},
		{
			name:     "考试结束时刻",
			now:      exam.ExamEndDate,
			expected: false, // Before 不包括等于
		},
		{
			name:     "考试结束后1分钟",
			now:      exam.ExamEndDate.Add(1 * time.Minute),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsExamTime(exam, tt.now)
			if result != tt.expected {
				t.Errorf("IsExamTime() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsExpiredExam(t *testing.T) {
	exam := createTestExam()
	// 考试结束时间: 2026-06-10 17:00:00

	tests := []struct {
		name     string
		now      time.Time
		expected bool
	}{
		{
			name:     "考试开始前",
			now:      exam.ExamBeginDate.Add(-1 * time.Hour),
			expected: false,
		},
		{
			name:     "考试进行中",
			now:      exam.ExamBeginDate.Add(1 * time.Hour),
			expected: false,
		},
		{
			name:     "考试结束前1分钟",
			now:      exam.ExamEndDate.Add(-1 * time.Minute),
			expected: false,
		},
		{
			name:     "考试结束时刻",
			now:      exam.ExamEndDate,
			expected: false, // After 不包括等于
		},
		{
			name:     "考试结束后1秒",
			now:      exam.ExamEndDate.Add(1 * time.Second),
			expected: true,
		},
		{
			name:     "考试结束后1天",
			now:      exam.ExamEndDate.Add(24 * time.Hour),
			expected: true,
		},
		{
			name:     "考试结束后1年",
			now:      exam.ExamEndDate.Add(365 * 24 * time.Hour),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsExpiredExam(exam, tt.now)
			if result != tt.expected {
				t.Errorf("IsExpiredExam() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetCountDownString(t *testing.T) {
	exam := createTestExam()
	template := "现在距离{exam}还有{time}"

	tests := []struct {
		name           string
		now            time.Time
		expectedResult string // 期望结果（可以是部分匹配）
		checkContains  bool   // 是否只检查包含关系
	}{
		{
			name:           "考试前350天",
			now:            exam.ExamBeginDate.Add(-350 * 24 * time.Hour),
			expectedResult: "350天",
			checkContains:  true,
		},
		{
			name:           "考试前100天",
			now:            exam.ExamBeginDate.Add(-100 * 24 * time.Hour),
			expectedResult: "100天",
			checkContains:  true,
		},
		{
			name:           "考试前1天",
			now:            exam.ExamBeginDate.Add(-24 * time.Hour),
			expectedResult: "1天",
			checkContains:  true,
		},
		{
			name:           "考试前5小时",
			now:            exam.ExamBeginDate.Add(-5 * time.Hour),
			expectedResult: "5小时",
			checkContains:  true,
		},
		{
			name:           "考试前30分钟",
			now:            exam.ExamBeginDate.Add(-30 * time.Minute),
			expectedResult: "30分钟",
			checkContains:  true,
		},
		{
			name:           "考试前5分钟",
			now:            exam.ExamBeginDate.Add(-5 * time.Minute),
			expectedResult: "5分钟",
			checkContains:  true,
		},
		{
			name:           "考试前59秒",
			now:            exam.ExamBeginDate.Add(-59 * time.Second),
			expectedResult: "59秒",
			checkContains:  true,
		},
		{
			name:           "考试前10秒",
			now:            exam.ExamBeginDate.Add(-10 * time.Second),
			expectedResult: "10秒",
			checkContains:  true,
		},
		{
			name:           "考试前1秒",
			now:            exam.ExamBeginDate.Add(-1 * time.Second),
			expectedResult: "1秒",
			checkContains:  true,
		},
		{
			name:           "考试前500毫秒",
			now:            exam.ExamBeginDate.Add(-500 * time.Millisecond),
			expectedResult: "1秒", // 小于1秒显示为1秒
			checkContains:  true,
		},
		{
			name:           "考试前100毫秒",
			now:            exam.ExamBeginDate.Add(-100 * time.Millisecond),
			expectedResult: "1秒", // 小于1秒显示为1秒
			checkContains:  true,
		},
		{
			name:           "考试进行中（开始后1小时）",
			now:            exam.ExamBeginDate.Add(1 * time.Hour),
			expectedResult: "2026年普通高等学校招生全国统一考试正在进行中！",
			checkContains:  false,
		},
		{
			name:           "考试进行中（第二天）",
			now:            time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC),
			expectedResult: "2026年普通高等学校招生全国统一考试正在进行中！",
			checkContains:  false,
		},
		{
			name:           "考试结束后1小时",
			now:            exam.ExamEndDate.Add(1 * time.Hour),
			expectedResult: "2026年普通高等学校招生全国统一考试已经结束了。",
			checkContains:  false,
		},
		{
			name:           "考试结束后1天",
			now:            exam.ExamEndDate.Add(24 * time.Hour),
			expectedResult: "2026年普通高等学校招生全国统一考试已经结束了。",
			checkContains:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCountDownString(exam, template, tt.now)
			if tt.checkContains {
				if !strings.Contains(result, tt.expectedResult) {
					t.Errorf("GetCountDownString() = %v, want to contain %v", result, tt.expectedResult)
				}
			} else {
				if result != tt.expectedResult {
					t.Errorf("GetCountDownString() = %v, want %v", result, tt.expectedResult)
				}
			}
		})
	}
}

func TestGetCountDownTime(t *testing.T) {
	exam := createTestExam()

	tests := []struct {
		name           string
		now            time.Time
		expectedResult string
		checkContains  bool
	}{
		{
			name:           "考试前350天",
			now:            exam.ExamBeginDate.Add(-350 * 24 * time.Hour),
			expectedResult: "350天",
			checkContains:  true,
		},
		{
			name:           "考试前1天2小时3分钟4秒",
			now:            exam.ExamBeginDate.Add(-26*time.Hour - 3*time.Minute - 4*time.Second),
			expectedResult: "1天2小时3分钟4秒",
			checkContains:  false,
		},
		{
			name:           "考试前5小时30分钟",
			now:            exam.ExamBeginDate.Add(-5*time.Hour - 30*time.Minute),
			expectedResult: "5小时30分钟",
			checkContains:  true,
		},
		{
			name:           "考试前1分钟30秒",
			now:            exam.ExamBeginDate.Add(-1*time.Minute - 30*time.Second),
			expectedResult: "1分钟30秒",
			checkContains:  false,
		},
		{
			name:           "考试前10秒",
			now:            exam.ExamBeginDate.Add(-10 * time.Second),
			expectedResult: "10秒",
			checkContains:  false,
		},
		{
			name:           "考试前1秒",
			now:            exam.ExamBeginDate.Add(-1 * time.Second),
			expectedResult: "1秒",
			checkContains:  false,
		},
		{
			name:           "考试前500毫秒",
			now:            exam.ExamBeginDate.Add(-500 * time.Millisecond),
			expectedResult: "1秒", // 小于1秒显示为1秒
			checkContains:  false,
		},
		{
			name:           "考试进行中",
			now:            exam.ExamBeginDate.Add(1 * time.Hour),
			expectedResult: "",
			checkContains:  false,
		},
		{
			name:           "考试结束后",
			now:            exam.ExamEndDate.Add(1 * time.Hour),
			expectedResult: "",
			checkContains:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCountDownTime(exam, tt.now)
			if tt.checkContains {
				if !strings.Contains(result, tt.expectedResult) {
					t.Errorf("GetCountDownTime() = %v, want to contain %v", result, tt.expectedResult)
				}
			} else {
				if result != tt.expectedResult {
					t.Errorf("GetCountDownTime() = %v, want %v", result, tt.expectedResult)
				}
			}
		})
	}
}

// TestCountDownScenarios 综合场景测试
func TestCountDownScenarios(t *testing.T) {
	exam := createTestExam()
	template := "现在距离{exam}还有{time}"

	t.Run("完整倒计时场景", func(t *testing.T) {
		// 2025年6月11日（考试前约362天）
		now := time.Date(2025, 6, 11, 9, 0, 0, 0, time.UTC)
		result := GetCountDownString(exam, template, now)
		t.Logf("考试前约362天: %s", result)

		if !strings.Contains(result, "还有") {
			t.Errorf("Expected countdown message, got: %s", result)
		}
	})

	t.Run("考试前最后一分钟", func(t *testing.T) {
		// 考试前59秒
		now := exam.ExamBeginDate.Add(-59 * time.Second)
		result := GetCountDownString(exam, template, now)
		t.Logf("考试前59秒: %s", result)

		if !strings.Contains(result, "59秒") {
			t.Errorf("Expected 59秒, got: %s", result)
		}
	})

	t.Run("考试前几毫秒", func(t *testing.T) {
		// 考试前100毫秒
		now := exam.ExamBeginDate.Add(-100 * time.Millisecond)
		result := GetCountDownString(exam, template, now)
		t.Logf("考试前100毫秒: %s", result)

		if !strings.Contains(result, "1秒") {
			t.Errorf("Expected 1秒 (for sub-second duration), got: %s", result)
		}
	})

	t.Run("考试开始瞬间", func(t *testing.T) {
		// 考试开始后10秒（在开始时刻窗口内）
		now := exam.ExamBeginDate.Add(10 * time.Second)
		result := GetCountDownString(exam, template, now)
		t.Logf("考试开始后10秒: %s", result)

		if !strings.Contains(result, "正在进行中") {
			t.Errorf("Expected '正在进行中', got: %s", result)
		}
	})

	t.Run("考试第二天", func(t *testing.T) {
		// 第二天上午
		now := time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)
		result := GetCountDownString(exam, template, now)
		t.Logf("考试第二天: %s", result)

		if !strings.Contains(result, "正在进行中") {
			t.Errorf("Expected '正在进行中', got: %s", result)
		}
	})

	t.Run("考试刚结束", func(t *testing.T) {
		// 考试结束后1分钟
		now := exam.ExamEndDate.Add(1 * time.Minute)
		result := GetCountDownString(exam, template, now)
		t.Logf("考试结束后1分钟: %s", result)

		if !strings.Contains(result, "已经结束了") {
			t.Errorf("Expected '已经结束了', got: %s", result)
		}
	})
}
