package util

import (
	"testing"

	"github.com/mymmrac/telego"
)

func TestIsUserChat(t *testing.T) {
	tests := []struct {
		name string
		chat *telego.Chat
		want bool
	}{
		{
			name: "Private chat",
			chat: &telego.Chat{
				Type: telego.ChatTypePrivate,
			},
			want: true,
		},
		{
			name: "Group chat",
			chat: &telego.Chat{
				Type: telego.ChatTypeGroup,
			},
			want: false,
		},
		{
			name: "Supergroup chat",
			chat: &telego.Chat{
				Type: telego.ChatTypeSupergroup,
			},
			want: false,
		},
		{
			name: "Channel",
			chat: &telego.Chat{
				Type: telego.ChatTypeChannel,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsUserChat(tt.chat)
			if got != tt.want {
				t.Errorf("IsUserChat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTextByMessage(t *testing.T) {
	tests := []struct {
		name string
		msg  *telego.Message
		want string
	}{
		{
			name: "Nil message",
			msg:  nil,
			want: "",
		},
		{
			name: "Empty text",
			msg: &telego.Message{
				Text: "",
			},
			want: "",
		},
		{
			name: "Simple text",
			msg: &telego.Message{
				Text: "Hello World",
			},
			want: "Hello World",
		},
		{
			name: "Text with leading/trailing spaces",
			msg: &telego.Message{
				Text: "  Hello World  ",
			},
			want: "Hello World",
		},
		{
			name: "Command with argument",
			msg: &telego.Message{
				Text: "/start Hello World",
			},
			want: "Hello World",
		},
		{
			name: "Command with multiple words",
			msg: &telego.Message{
				Text: "/countdown 2025年高考",
			},
			want: "2025年高考",
		},
		{
			name: "Command only (no argument)",
			msg: &telego.Message{
				Text: "/help",
			},
			want: "",
		},
		{
			name: "Command with extra spaces",
			msg: &telego.Message{
				Text: "/start   Hello World  ",
			},
			want: "Hello World",
		},
		{
			name: "Command with bot username",
			msg: &telego.Message{
				Text: "/start@botname Hello",
			},
			want: "Hello",
		},
		{
			name: "Not a command (slash in middle)",
			msg: &telego.Message{
				Text: "Hello /world",
			},
			want: "Hello /world",
		},
		{
			name: "Chinese text",
			msg: &telego.Message{
				Text: "你好世界",
			},
			want: "你好世界",
		},
		{
			name: "Command with Chinese argument",
			msg: &telego.Message{
				Text: "/template 默认模板",
			},
			want: "默认模板",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTextByMessage(tt.msg)
			if got != tt.want {
				t.Errorf("GetTextByMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRemoveFirst(t *testing.T) {
	tests := []struct {
		name string
		s    string
		char string
		want string
	}{
		{
			name: "Remove first occurrence",
			s:    "hello world",
			char: "l",
			want: "helo world",
		},
		{
			name: "Remove first word",
			s:    "hello world hello",
			char: "hello",
			want: " world hello",
		},
		{
			name: "Character not found",
			s:    "hello",
			char: "x",
			want: "hello",
		},
		{
			name: "Remove from start",
			s:    "/start command",
			char: "/",
			want: "start command",
		},
		{
			name: "Multiple same characters - only first removed",
			s:    "aaa",
			char: "a",
			want: "aa",
		},
		{
			name: "Empty string",
			s:    "",
			char: "a",
			want: "",
		},
		{
			name: "Empty char",
			s:    "hello",
			char: "",
			want: "hello",
		},
		{
			name: "Chinese character",
			s:    "你好你好",
			char: "你",
			want: "好你好",
		},
		{
			name: "Remove space",
			s:    "hello world",
			char: " ",
			want: "helloworld",
		},
		{
			name: "Remove special character",
			s:    "@user@name",
			char: "@",
			want: "user@name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveFirst(tt.s, tt.char)
			if got != tt.want {
				t.Errorf("RemoveFirst(%q, %q) = %q, want %q", tt.s, tt.char, got, tt.want)
			}
		})
	}
}
