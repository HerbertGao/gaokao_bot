package util

import "regexp"

var (
	commandRegex      = regexp.MustCompile(`^/[a-zA-Z]+`)
	templateNameRegex = regexp.MustCompile(`【(.+?)】`)
	templateIDRegex   = regexp.MustCompile(`/rm_(\d+)`)
)

// IsMatchCommand 检查是否匹配命令格式
func IsMatchCommand(text string) bool {
	return commandRegex.MatchString(text)
}

// ExtractTemplateName 提取【】括号内容
func ExtractTemplateName(text string) string {
	matches := templateNameRegex.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// ExtractTemplateID 提取模板 ID
func ExtractTemplateID(text string) string {
	matches := templateIDRegex.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// Get 通用正则匹配和分组提取
func Get(pattern, text string, group int) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(text)
	if len(matches) > group {
		return matches[group]
	}
	return ""
}