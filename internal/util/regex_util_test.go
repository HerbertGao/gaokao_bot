package util

import "testing"

func TestIsMatchCommand(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{
			name: "Simple command",
			text: "/start",
			want: true,
		},
		{
			name: "Command with arguments",
			text: "/help me",
			want: true,
		},
		{
			name: "Command lowercase",
			text: "/countdown",
			want: true,
		},
		{
			name: "Command uppercase",
			text: "/START",
			want: true,
		},
		{
			name: "Command mixed case",
			text: "/myCommand",
			want: true,
		},
		{
			name: "Not a command - no slash",
			text: "hello",
			want: false,
		},
		{
			name: "Not a command - slash in middle",
			text: "hello /world",
			want: false,
		},
		{
			name: "Not a command - slash only",
			text: "/",
			want: false,
		},
		{
			name: "Not a command - slash with number",
			text: "/123",
			want: false,
		},
		{
			name: "Not a command - slash with space",
			text: "/ start",
			want: false,
		},
		{
			name: "Empty string",
			text: "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsMatchCommand(tt.text)
			if got != tt.want {
				t.Errorf("IsMatchCommand(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestExtractTemplateName(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "Simple template name",
			text: "【默认模板】",
			want: "默认模板",
		},
		{
			name: "Template name with spaces",
			text: "【我的 模板】",
			want: "我的 模板",
		},
		{
			name: "Template name in sentence",
			text: "使用【自定义模板】发送消息",
			want: "自定义模板",
		},
		{
			name: "Multiple template names - returns first",
			text: "【模板1】和【模板2】",
			want: "模板1",
		},
		{
			name: "Template name with numbers",
			text: "【模板123】",
			want: "模板123",
		},
		{
			name: "Template name with special chars",
			text: "【模板-名称_v2】",
			want: "模板-名称_v2",
		},
		{
			name: "No template brackets",
			text: "普通文本",
			want: "",
		},
		{
			name: "Empty brackets",
			text: "【】",
			want: "",
		},
		{
			name: "Left bracket only",
			text: "【模板",
			want: "",
		},
		{
			name: "Right bracket only",
			text: "模板】",
			want: "",
		},
		{
			name: "English brackets",
			text: "[template]",
			want: "",
		},
		{
			name: "Empty string",
			text: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractTemplateName(tt.text)
			if got != tt.want {
				t.Errorf("ExtractTemplateName(%q) = %q, want %q", tt.text, got, tt.want)
			}
		})
	}
}

func TestExtractTemplateID(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "Simple template ID",
			text: "/rm_123",
			want: "123",
		},
		{
			name: "Large template ID",
			text: "/rm_999999999",
			want: "999999999",
		},
		{
			name: "Template ID in command",
			text: "/rm_456 删除模板",
			want: "456",
		},
		{
			name: "Single digit ID",
			text: "/rm_1",
			want: "1",
		},
		{
			name: "Zero ID",
			text: "/rm_0",
			want: "0",
		},
		{
			name: "Multiple IDs - returns first",
			text: "/rm_111 /rm_222",
			want: "111",
		},
		{
			name: "No ID pattern",
			text: "/rm",
			want: "",
		},
		{
			name: "Wrong pattern",
			text: "/remove_123",
			want: "",
		},
		{
			name: "ID with letters",
			text: "/rm_abc",
			want: "",
		},
		{
			name: "Empty string",
			text: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractTemplateID(tt.text)
			if got != tt.want {
				t.Errorf("ExtractTemplateID(%q) = %q, want %q", tt.text, got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		text    string
		group   int
		want    string
	}{
		{
			name:    "Simple capture group",
			pattern: `(\d+)`,
			text:    "123",
			group:   1,
			want:    "123",
		},
		{
			name:    "Multiple capture groups - get first",
			pattern: `(\d+)-(\w+)`,
			text:    "123-abc",
			group:   1,
			want:    "123",
		},
		{
			name:    "Multiple capture groups - get second",
			pattern: `(\d+)-(\w+)`,
			text:    "123-abc",
			group:   2,
			want:    "abc",
		},
		{
			name:    "Get full match (group 0)",
			pattern: `\d+`,
			text:    "123",
			group:   0,
			want:    "123",
		},
		{
			name:    "Pattern in longer text",
			pattern: `email:\s*(\S+)`,
			text:    "My email: test@example.com is here",
			group:   1,
			want:    "test@example.com",
		},
		{
			name:    "No match",
			pattern: `\d+`,
			text:    "abc",
			group:   1,
			want:    "",
		},
		{
			name:    "Group index out of range",
			pattern: `(\d+)`,
			text:    "123",
			group:   2,
			want:    "",
		},
		{
			name:    "Empty text",
			pattern: `\d+`,
			text:    "",
			group:   1,
			want:    "",
		},
		{
			name:    "Chinese characters",
			pattern: `【(.+?)】`,
			text:    "这是【中文模板】测试",
			group:   1,
			want:    "中文模板",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Get(tt.pattern, tt.text, tt.group)
			if got != tt.want {
				t.Errorf("Get(%q, %q, %d) = %q, want %q", tt.pattern, tt.text, tt.group, got, tt.want)
			}
		})
	}
}
