package version

import (
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	version := GetVersion()
	if version == "" {
		t.Error("Version should not be empty")
	}

	// 默认版本应该是 "1.0.0"
	if version != "1.0.0" {
		t.Logf("Version = %s (may be overridden by ldflags)", version)
	}
}

func TestGetFullVersionInfo(t *testing.T) {
	info := GetFullVersionInfo()
	if info == "" {
		t.Error("Version info should not be empty")
	}

	// 验证包含必要的字段
	requiredFields := []string{
		"Version:",
		"Build Time:",
		"Git Commit:",
		"Go Version:",
		"OS/Arch:",
	}

	for _, field := range requiredFields {
		if !strings.Contains(info, field) {
			t.Errorf("Version info should contain %q", field)
		}
	}
}
