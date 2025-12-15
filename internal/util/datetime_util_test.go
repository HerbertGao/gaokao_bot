package util

import (
	"testing"
	"time"
)

func TestNormalizeToMinute(t *testing.T) {
	bjtZone := time.FixedZone("BJT", 8*3600) // UTC+8

	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "秒数<30，保持当前分钟",
			input:    time.Date(2025, 6, 7, 9, 15, 29, 0, bjtZone),
			expected: time.Date(2025, 6, 7, 9, 15, 0, 0, bjtZone),
		},
		{
			name:     "秒数=30，进位到下一分钟",
			input:    time.Date(2025, 6, 7, 9, 15, 30, 0, bjtZone),
			expected: time.Date(2025, 6, 7, 9, 16, 0, 0, bjtZone),
		},
		{
			name:     "秒数>30，进位到下一分钟",
			input:    time.Date(2025, 6, 7, 9, 15, 45, 0, bjtZone),
			expected: time.Date(2025, 6, 7, 9, 16, 0, 0, bjtZone),
		},
		{
			name:     "59秒进位到下一分钟",
			input:    time.Date(2025, 6, 7, 9, 15, 59, 0, bjtZone),
			expected: time.Date(2025, 6, 7, 9, 16, 0, 0, bjtZone),
		},
		{
			name:     "8:59:30进位到9:00",
			input:    time.Date(2025, 6, 7, 8, 59, 30, 0, bjtZone),
			expected: time.Date(2025, 6, 7, 9, 0, 0, 0, bjtZone),
		},
		{
			name:     "整分钟保持不变",
			input:    time.Date(2025, 6, 7, 9, 0, 0, 0, bjtZone),
			expected: time.Date(2025, 6, 7, 9, 0, 0, 0, bjtZone),
		},
		{
			name:     "带纳秒的时间也会被截断",
			input:    time.Date(2025, 6, 7, 9, 15, 10, 500000000, bjtZone),
			expected: time.Date(2025, 6, 7, 9, 15, 0, 0, bjtZone),
		},
		{
			name:     "23:59:59跨日边界，进位到次日00:00",
			input:    time.Date(2025, 6, 6, 23, 59, 59, 0, bjtZone),
			expected: time.Date(2025, 6, 7, 0, 0, 0, 0, bjtZone),
		},
		{
			name:     "23:59:30跨日边界，进位到次日00:00",
			input:    time.Date(2025, 6, 6, 23, 59, 30, 0, bjtZone),
			expected: time.Date(2025, 6, 7, 0, 0, 0, 0, bjtZone),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeToMinute(tt.input)
			if !result.Equal(tt.expected) {
				t.Errorf("NormalizeToMinute() = %v, want %v",
					result.Format("2006-01-02 15:04:05.000"),
					tt.expected.Format("2006-01-02 15:04:05.000"))
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "负数时长",
			duration: -5 * time.Second,
			expected: "0秒",
		},
		{
			name:     "零时长",
			duration: 0,
			expected: "0秒",
		},
		{
			name:     "小于1秒",
			duration: 500 * time.Millisecond,
			expected: "1秒",
		},
		{
			name:     "只有秒",
			duration: 45 * time.Second,
			expected: "45秒",
		},
		{
			name:     "只有分",
			duration: 3 * time.Minute,
			expected: "3分钟",
		},
		{
			name:     "分秒",
			duration: 3*time.Minute + 25*time.Second,
			expected: "3分钟25秒",
		},
		{
			name:     "小时分秒",
			duration: 2*time.Hour + 15*time.Minute + 30*time.Second,
			expected: "2小时15分钟30秒",
		},
		{
			name:     "天小时分秒",
			duration: 350*24*time.Hour + 23*time.Hour + 59*time.Minute + 59*time.Second,
			expected: "350天23小时59分钟59秒",
		},
		{
			name:     "只有天和分",
			duration: 18*24*time.Hour + 3*time.Minute,
			expected: "18天3分钟",
		},
		{
			name:     "只有小时",
			duration: 5 * time.Hour,
			expected: "5小时",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("FormatDuration() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFormatDurationWithMillis(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "负数时长",
			duration: -5 * time.Second,
			expected: "0秒",
		},
		{
			name:     "零时长",
			duration: 0,
			expected: "0秒",
		},
		{
			name:     "只有毫秒",
			duration: 500 * time.Millisecond,
			expected: "500毫秒",
		},
		{
			name:     "秒和毫秒",
			duration: 45*time.Second + 123*time.Millisecond,
			expected: "45秒123毫秒",
		},
		{
			name:     "分秒毫秒",
			duration: 3*time.Minute + 25*time.Second + 678*time.Millisecond,
			expected: "3分钟25秒678毫秒",
		},
		{
			name:     "小时分秒毫秒",
			duration: 2*time.Hour + 15*time.Minute + 30*time.Second + 999*time.Millisecond,
			expected: "2小时15分钟30秒999毫秒",
		},
		{
			name:     "天小时分秒毫秒",
			duration: 350*24*time.Hour + 23*time.Hour + 59*time.Minute + 59*time.Second + 500*time.Millisecond,
			expected: "350天23小时59分钟59秒500毫秒",
		},
		{
			name:     "只有天和分和毫秒",
			duration: 18*24*time.Hour + 3*time.Minute + 200*time.Millisecond,
			expected: "18天3分钟200毫秒",
		},
		{
			name:     "只有小时（无毫秒）",
			duration: 5 * time.Hour,
			expected: "5小时",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDurationWithMillis(tt.duration)
			if result != tt.expected {
				t.Errorf("FormatDurationWithMillis() = %v, want %v", result, tt.expected)
			}
		})
	}
}
