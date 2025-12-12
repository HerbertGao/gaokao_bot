package service

import (
	"testing"

	"github.com/mymmrac/telego"
	"github.com/sirupsen/logrus"
)

// 注意：由于 BotService 依赖真实的 Telegram Bot API，
// 这些测试主要覆盖不需要实际网络调用的逻辑部分

func TestNewBotService(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// 创建 nil bot（在测试中我们不会真正调用 API）
	messageService := &MessageService{}
	inlineQueryService := &InlineQueryService{}
	miniAppURL := "https://example.com"

	service := NewBotService(nil, messageService, inlineQueryService, logger, miniAppURL)

	if service == nil {
		t.Fatal("NewBotService() returned nil")
	}

	if service.miniAppURL != miniAppURL {
		t.Errorf("miniAppURL = %s, want %s", service.miniAppURL, miniAppURL)
	}

	if service.logger != logger {
		t.Error("logger not set correctly")
	}
}

func TestHandleMessage_NilMessage(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewBotService(nil, nil, nil, logger, "")

	// 测试 nil 消息不应该导致 panic
	service.HandleMessage(nil, nil)
	// 如果没有 panic，测试通过
}

func TestHandleMessage_EmptyText(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewBotService(nil, nil, nil, logger, "")

	msg := &telego.Message{
		Text: "",
		Chat: telego.Chat{ID: 123},
	}

	// 测试空文本消息不应该处理
	service.HandleMessage(nil, msg)
	// 如果没有 panic，测试通过
}

func TestHandleMessage_NonCommandText(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewBotService(nil, nil, nil, logger, "")

	msg := &telego.Message{
		Text: "Hello, this is not a command",
		Chat: telego.Chat{ID: 123},
	}

	// 测试非命令消息不应该被处理
	service.HandleMessage(nil, msg)
	// 如果没有 panic，测试通过
}

func TestHandleInlineQuery_NilQuery(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewBotService(nil, nil, nil, logger, "")

	// 测试 nil 查询不应该导致 panic
	service.HandleInlineQuery(nil, nil)
	// 如果没有 panic，测试通过
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		maxLen  int
		want    string
	}{
		{
			name:   "短字符串不截断",
			input:  "Hello",
			maxLen: 10,
			want:   "Hello",
		},
		{
			name:   "相等长度不截断",
			input:  "Hello",
			maxLen: 5,
			want:   "Hello",
		},
		{
			name:   "超长字符串截断",
			input:  "Hello, World!",
			maxLen: 5,
			want:   "Hello...",
		},
		{
			name:   "空字符串",
			input:  "",
			maxLen: 10,
			want:   "",
		},
		{
			name:   "单字符截断",
			input:  "ABC",
			maxLen: 1,
			want:   "A...",
		},
		{
			name:   "长英文文本截断",
			input:  "This is a very long string that needs to be truncated",
			maxLen: 20,
			want:   "This is a very long ...",
		},
		// 多字节 UTF-8 字符（中文）测试
		{
			name:   "中文字符串不截断",
			input:  "你好世界",
			maxLen: 10,
			want:   "你好世界",
		},
		{
			name:   "中文字符串相等长度不截断",
			input:  "你好世界",
			maxLen: 4,
			want:   "你好世界",
		},
		{
			name:   "中文字符串截断",
			input:  "你好世界，欢迎使用高考倒计时机器人",
			maxLen: 5,
			want:   "你好世界，...",
		},
		{
			name:   "中英文混合截断",
			input:  "Hello你好World世界",
			maxLen: 8,
			want:   "Hello你好W...",
		},
		{
			name:   "单个中文字符截断",
			input:  "你好世界",
			maxLen: 1,
			want:   "你...",
		},
		// Emoji 测试（也是多字节 UTF-8）
		{
			name:   "包含 Emoji 的字符串截断",
			input:  "加油💪努力学习📚考上理想大学🎓",
			maxLen: 6,
			want:   "加油💪努力学...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateString(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestTruncateString_EdgeCases(t *testing.T) {
	// 测试 maxLen = 0
	result := truncateString("test", 0)
	if result != "..." {
		t.Errorf("truncateString with maxLen=0, got %q, want %q", result, "...")
	}

	// 注意：负数 maxLen 会导致 panic，这是函数的预期行为
	// 在实际使用中应该避免传入负数
}
