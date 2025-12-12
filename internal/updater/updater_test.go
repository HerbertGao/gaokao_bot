package updater

import (
	"runtime"
	"testing"
)

func TestCompareVersion(t *testing.T) {
	u := NewUpdater()

	tests := []struct {
		name    string
		current string
		latest  string
		want    bool
	}{
		{
			name:    "Latest is newer - major version",
			current: "1.0.0",
			latest:  "2.0.0",
			want:    true,
		},
		{
			name:    "Latest is newer - minor version",
			current: "1.0.0",
			latest:  "1.1.0",
			want:    true,
		},
		{
			name:    "Latest is newer - patch version",
			current: "1.0.0",
			latest:  "1.0.1",
			want:    true,
		},
		{
			name:    "Same version",
			current: "1.0.0",
			latest:  "1.0.0",
			want:    false,
		},
		{
			name:    "Current is newer",
			current: "2.0.0",
			latest:  "1.0.0",
			want:    false,
		},
		{
			name:    "With v prefix - latest newer",
			current: "v1.0.0",
			latest:  "v1.1.0",
			want:    true,
		},
		{
			name:    "Mixed v prefix - latest newer",
			current: "1.0.0",
			latest:  "v2.0.0",
			want:    true,
		},
		{
			name:    "Different length versions - latest newer",
			current: "1.0",
			latest:  "1.0.1",
			want:    true,
		},
		{
			name:    "Different length versions - same",
			current: "1.0.0",
			latest:  "1.0",
			want:    false,
		},
		{
			name:    "Four part version - latest newer",
			current: "1.0.0.0",
			latest:  "1.0.0.1",
			want:    true,
		},
		{
			name:    "Complex version comparison",
			current: "v1.2.3",
			latest:  "v1.2.4",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := u.CompareVersion(tt.current, tt.latest)
			if got != tt.want {
				t.Errorf("CompareVersion(%q, %q) = %v, want %v", tt.current, tt.latest, got, tt.want)
			}
		})
	}
}

func TestParseVersionParts(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    []int
	}{
		{
			name:    "Standard version",
			version: "1.0.0",
			want:    []int{1, 0, 0},
		},
		{
			name:    "Two part version",
			version: "1.0",
			want:    []int{1, 0},
		},
		{
			name:    "Four part version",
			version: "1.2.3.4",
			want:    []int{1, 2, 3, 4},
		},
		{
			name:    "Single number",
			version: "1",
			want:    []int{1},
		},
		{
			name:    "With invalid parts",
			version: "1.0.beta",
			want:    []int{1, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseVersionParts(tt.version)
			if len(got) != len(tt.want) {
				t.Errorf("parseVersionParts(%q) length = %d, want %d", tt.version, len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseVersionParts(%q)[%d] = %d, want %d", tt.version, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestGetTargetAssetName(t *testing.T) {
	u := NewUpdater()

	// 测试当前平台能否获取到 asset 名称
	assetName, err := u.getTargetAssetName()
	if err != nil {
		// 如果当前平台不支持，跳过测试
		t.Skipf("Current platform not supported: %v", err)
	}

	if assetName == "" {
		t.Error("getTargetAssetName() returned empty string")
	}

	// 验证返回的名称格式正确
	// 应该匹配 release.yml 中的 artifact_name 格式
	expectedSuffixes := []string{"_arm64", "_x86_64"}
	expectedPrefixes := []string{"macos", "linux", "windows"}

	hasValidPrefix := false
	for _, prefix := range expectedPrefixes {
		if len(assetName) >= len(prefix) && assetName[:len(prefix)] == prefix {
			hasValidPrefix = true
			break
		}
	}

	hasValidSuffix := false
	for _, suffix := range expectedSuffixes {
		if len(assetName) >= len(suffix) && assetName[len(assetName)-len(suffix):] == suffix {
			hasValidSuffix = true
			break
		}
	}

	if !hasValidPrefix {
		t.Errorf("Asset name %q should start with one of %v", assetName, expectedPrefixes)
	}
	if !hasValidSuffix {
		t.Errorf("Asset name %q should end with one of %v", assetName, expectedSuffixes)
	}

	// 验证名称格式: {os}_{arch}，例如 macos_arm64
	if !hasValidPrefix || !hasValidSuffix {
		t.Errorf("Asset name format incorrect: %q", assetName)
	}
}

// TestGetTargetAssetNameMapping 验证所有平台的资产名称映射正确
// 确保与 .github/workflows/release.yml 中的 artifact_name 保持一致
func TestGetTargetAssetNameMapping(t *testing.T) {
	tests := []struct {
		goos     string
		goarch   string
		expected string
	}{
		{goos: "darwin", goarch: "arm64", expected: "macos_arm64"},
		{goos: "darwin", goarch: "amd64", expected: "macos_x86_64"},
		{goos: "linux", goarch: "amd64", expected: "linux_x86_64"},
		{goos: "linux", goarch: "arm64", expected: "linux_arm64"},
		{goos: "windows", goarch: "amd64", expected: "windows_x86_64"},
	}

	for _, tt := range tests {
		t.Run(tt.goos+"_"+tt.goarch, func(t *testing.T) {
			// 暂时保存原始值
			oldGOOS := runtime.GOOS
			oldGOARCH := runtime.GOARCH

			// 注意：runtime.GOOS 和 runtime.GOARCH 是只读常量，无法修改
			// 这个测试主要用于文档化预期的映射关系
			// 实际测试在真实平台上运行时会验证当前平台

			// 如果当前平台匹配，验证结果
			if runtime.GOOS == tt.goos && runtime.GOARCH == tt.goarch {
				u := NewUpdater()
				got, err := u.getTargetAssetName()
				if err != nil {
					t.Errorf("getTargetAssetName() error = %v", err)
					return
				}
				if got != tt.expected {
					t.Errorf("getTargetAssetName() = %q, want %q", got, tt.expected)
				}
			}

			// 防止编译器警告
			_ = oldGOOS
			_ = oldGOARCH
		})
	}

	// 额外验证：当前平台的资产名称应该匹配 release.yml 中的某个 artifact_name
	u := NewUpdater()
	assetName, err := u.getTargetAssetName()
	if err != nil {
		t.Skipf("Current platform not supported: %v", err)
	}

	// release.yml 中定义的所有 artifact_name（去掉 gaokao_bot_ 前缀）
	validAssetNames := map[string]bool{
		"macos_arm64":   true,
		"macos_x86_64":  true,
		"linux_x86_64":  true,
		"linux_arm64":   true,
		"windows_x86_64": true,
	}

	if !validAssetNames[assetName] {
		t.Errorf("Asset name %q does not match any artifact in release.yml", assetName)
		t.Logf("Valid asset names: %v", validAssetNames)
	}
}

func TestNewUpdater(t *testing.T) {
	u := NewUpdater()

	if u == nil {
		t.Fatal("NewUpdater() returned nil")
	}

	if u.client == nil {
		t.Error("Updater client is nil")
	}

	if u.client.Timeout == 0 {
		t.Error("Updater client timeout not set")
	}
}

func TestGetCurrentVersion(t *testing.T) {
	u := NewUpdater()

	version := u.GetCurrentVersion()
	if version == "" {
		t.Error("GetCurrentVersion() returned empty string")
	}
}
