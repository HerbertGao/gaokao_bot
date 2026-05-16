package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/util"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoapi"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

// mockGuestCaller 模拟 Telegram API 调用，用于 Guest 查询测试，避免真实网络请求
type mockGuestCaller struct {
	calls    int
	lastBody []byte
	resp     *telegoapi.Response
	err      error
}

func (m *mockGuestCaller) Call(_ context.Context, _ string, data *telegoapi.RequestData) (*telegoapi.Response, error) {
	m.calls++
	if data != nil {
		m.lastBody = data.BodyRaw
	}
	return m.resp, m.err
}

// guestQueryPayload 解析 answerGuestQuery 请求体，用于断言应答内容
type guestQueryPayload struct {
	GuestQueryID string `json:"guest_query_id"`
	Result       struct {
		Type                string `json:"type"`
		ID                  string `json:"id"`
		Title               string `json:"title"`
		InputMessageContent struct {
			MessageText string `json:"message_text"`
		} `json:"input_message_content"`
	} `json:"result"`
}

func decodeGuestPayload(t *testing.T, body []byte) guestQueryPayload {
	t.Helper()
	if len(body) == 0 {
		t.Fatal("expected a non-empty answerGuestQuery request body")
	}
	var p guestQueryPayload
	if err := json.Unmarshal(body, &p); err != nil {
		t.Fatalf("decode answerGuestQuery body: %v", err)
	}
	return p
}

func newGuestTestBot(t *testing.T, caller telegoapi.Caller) *telego.Bot {
	t.Helper()
	bot, err := telego.NewBot(
		"123456:abcdefghijklmnopqrstuvwxyz012345678",
		telego.WithAPICaller(caller),
		telego.WithDiscardLogger(),
	)
	if err != nil {
		t.Fatalf("NewBot() error = %v", err)
	}
	return bot
}

// setupGuestTestService 构造一个使用 mock caller 的 BotService 及其测试数据库
func setupGuestTestService(t *testing.T, caller telegoapi.Caller) (*BotService, *gorm.DB) {
	t.Helper()
	messageService, db := setupMessageTestService(t)

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	bot := newGuestTestBot(t, caller)
	service := NewBotService(bot, messageService, nil, logger, "")
	return service, db
}

func okGuestCaller() *mockGuestCaller {
	return &mockGuestCaller{
		resp: &telegoapi.Response{
			Ok:     true,
			Result: json.RawMessage(`{"inline_message_id":"m1"}`),
		},
	}
}

func TestHandleGuestMessage_NilMessage(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	service := NewBotService(nil, nil, nil, logger, "")

	// nil 消息不应该 panic
	service.HandleGuestMessage(nil, nil)
}

func TestHandleGuestMessage_EmptyQueryID(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	service := NewBotService(nil, nil, nil, logger, "")

	// 缺少 GuestQueryID 时应提前返回，不调用 API、不 panic
	service.HandleGuestMessage(nil, &telego.Message{Text: "@gaokao_bot"})
}

func TestHandleGuestMessage_NoArg(t *testing.T) {
	caller := okGuestCaller()
	service, db := setupGuestTestService(t, caller)

	now := time.Now()
	futureDate := now.AddDate(1, 0, 0)
	db.Create(&model.ExamDate{
		ID:                1,
		ExamYear:          futureDate.Year(),
		ExamDesc:          "高考",
		ShortDesc:         "高考",
		ExamBeginDate:     futureDate,
		ExamEndDate:       futureDate.AddDate(0, 0, 3),
		ExamYearBeginDate: now,
		ExamYearEndDate:   futureDate.AddDate(0, 0, 3),
		IsDelete:          false,
	})

	service.HandleGuestMessage(service.bot, &telego.Message{
		Text:         "@gaokao_bot",
		GuestQueryID: "q-noarg",
	})

	if caller.calls != 1 {
		t.Errorf("expected 1 API call, got %d", caller.calls)
	}
}

func TestHandleGuestMessage_WithYear(t *testing.T) {
	caller := okGuestCaller()
	service, db := setupGuestTestService(t, caller)

	year := 2026
	db.Create(&model.ExamDate{
		ID:                1,
		ExamYear:          year,
		ExamDesc:          "2026年高考",
		ShortDesc:         "高考",
		ExamBeginDate:     time.Date(year, 6, 7, 9, 0, 0, 0, util.GetBJTLocation()),
		ExamEndDate:       time.Date(year, 6, 10, 17, 0, 0, 0, util.GetBJTLocation()),
		ExamYearBeginDate: time.Date(year-1, 6, 10, 17, 0, 0, 0, util.GetBJTLocation()),
		ExamYearEndDate:   time.Date(year, 6, 10, 17, 0, 0, 0, util.GetBJTLocation()),
		IsDelete:          false,
	})

	service.HandleGuestMessage(service.bot, &telego.Message{
		Text:         "@gaokao_bot 2026",
		GuestQueryID: "q-year",
	})

	if caller.calls != 1 {
		t.Errorf("expected 1 API call, got %d", caller.calls)
	}

	payload := decodeGuestPayload(t, caller.lastBody)
	if payload.GuestQueryID != "q-year" {
		t.Errorf("guest_query_id = %q, want %q", payload.GuestQueryID, "q-year")
	}
	if payload.Result.Type != telego.ResultTypeArticle {
		t.Errorf("result type = %q, want %q", payload.Result.Type, telego.ResultTypeArticle)
	}
	if payload.Result.InputMessageContent.MessageText == "" {
		t.Error("expected non-empty message text in guest answer")
	}
}

func TestHandleGuestMessage_InvalidArg(t *testing.T) {
	caller := okGuestCaller()
	service, _ := setupGuestTestService(t, caller)

	service.HandleGuestMessage(service.bot, &telego.Message{
		Text:         "@gaokao_bot hello",
		GuestQueryID: "q-invalid",
	})

	if caller.calls != 1 {
		t.Errorf("expected 1 API call, got %d", caller.calls)
	}

	// 非法参数应应答「参数暂时无法识别。」
	payload := decodeGuestPayload(t, caller.lastBody)
	if payload.GuestQueryID != "q-invalid" {
		t.Errorf("guest_query_id = %q, want %q", payload.GuestQueryID, "q-invalid")
	}
	if payload.Result.InputMessageContent.MessageText != "参数暂时无法识别。" {
		t.Errorf("message text = %q, want %q",
			payload.Result.InputMessageContent.MessageText, "参数暂时无法识别。")
	}
}

func TestHandleGuestMessage_NoData(t *testing.T) {
	caller := okGuestCaller()
	service, _ := setupGuestTestService(t, caller)

	// 数据库无考试数据时仍应答提示文本
	service.HandleGuestMessage(service.bot, &telego.Message{
		Text:         "@gaokao_bot",
		GuestQueryID: "q-nodata",
	})

	if caller.calls != 1 {
		t.Errorf("expected 1 API call, got %d", caller.calls)
	}

	payload := decodeGuestPayload(t, caller.lastBody)
	if payload.Result.InputMessageContent.MessageText != "数据库中没有可用的信息，请联系开发者。" {
		t.Errorf("message text = %q, want %q",
			payload.Result.InputMessageContent.MessageText, "数据库中没有可用的信息，请联系开发者。")
	}
}

func TestHandleGuestMessage_AnswerFails(t *testing.T) {
	caller := &mockGuestCaller{err: errors.New("network error")}
	service, _ := setupGuestTestService(t, caller)

	// API 调用失败时应记录日志、不重试、不 panic
	service.HandleGuestMessage(service.bot, &telego.Message{
		Text:         "@gaokao_bot",
		GuestQueryID: "q-fail",
	})

	if caller.calls != 1 {
		t.Errorf("expected 1 API call, got %d", caller.calls)
	}
}
