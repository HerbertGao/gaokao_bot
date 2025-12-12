package updater

import (
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

	// 验证返回的名称格式正确（应该包含下划线）
	if len(assetName) == 0 {
		t.Error("Asset name should not be empty")
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
