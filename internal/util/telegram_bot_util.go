package util

import (
	"strings"
	"unicode"

	"github.com/mymmrac/telego"
)

// IsUserChat 判断是否为私聊
func IsUserChat(chat *telego.Chat) bool {
	return chat.Type == telego.ChatTypePrivate
}

// GetTextByMessage 从消息中提取文本（移除命令和 @）
func GetTextByMessage(msg *telego.Message) string {
	if msg == nil || msg.Text == "" {
		return ""
	}

	text := msg.Text

	// 如果是命令（以 / 开头），移除命令部分
	if strings.HasPrefix(text, "/") {
		parts := strings.SplitN(text, " ", 2)
		if len(parts) > 1 {
			text = parts[1]
		} else {
			text = ""
		}
	}

	return strings.TrimSpace(text)
}

// GetGuestMessageArg 从 Guest 消息中提取参数文本。
// 仅当消息以 @提及 开头时才剥离首个 @token 并返回其后的参数；
// 回复式召唤（文本不以 @ 开头）按无参数处理，返回空字符串。
// 提及与参数之间支持任意空白分隔（半角/全角空格、换行等）。
func GetGuestMessageArg(msg *telego.Message) string {
	if msg == nil {
		return ""
	}

	text := strings.TrimSpace(msg.Text)
	if !strings.HasPrefix(text, "@") {
		return ""
	}

	idx := strings.IndexFunc(text, unicode.IsSpace)
	if idx < 0 {
		return ""
	}

	return strings.TrimSpace(text[idx:])
}

// RemoveFirst 删除第一个指定字符
func RemoveFirst(s string, char string) string {
	return strings.Replace(s, char, "", 1)
}