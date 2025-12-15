package version

import (
	"fmt"
	"runtime"
)

// Version 应用版本号，通过 ldflags 注入
var Version = "10.0.4"

// BuildTime 构建时间，通过 ldflags 注入
var BuildTime = "unknown"

// GitCommit Git 提交哈希，通过 ldflags 注入
var GitCommit = "unknown"

// GetVersion 获取完整版本信息
func GetVersion() string {
	return Version
}

// GetFullVersionInfo 获取详细版本信息
func GetFullVersionInfo() string {
	return fmt.Sprintf(
		"Version: %s\nBuild Time: %s\nGit Commit: %s\nGo Version: %s\nOS/Arch: %s/%s",
		Version,
		BuildTime,
		GitCommit,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}
