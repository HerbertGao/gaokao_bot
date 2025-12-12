package util

import (
	"strings"

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

// RemoveFirst 删除第一个指定字符
func RemoveFirst(s string, char string) string {
	return strings.Replace(s, char, "", 1)
}