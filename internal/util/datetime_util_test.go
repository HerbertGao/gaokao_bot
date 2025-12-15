package util

import (
	"testing"
	"time"
)

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
