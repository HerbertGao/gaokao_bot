package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/version"
)

const (
	githubAPIBase = "https://api.github.com"
	repoOwner     = "HerbertGao"
	repoName      = "gaokao_bot"
	userAgent     = "gaokao_bot_updater/1.0"
)

// Updater 版本更新器
type Updater struct {
	client *http.Client
}

// Release GitHub Release 结构
type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

// Asset GitHub Release Asset 结构
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// NewUpdater 创建新的更新器
func NewUpdater() *Updater {
	return &Updater{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetCurrentVersion 获取当前版本
func (u *Updater) GetCurrentVersion() string {
	return version.GetVersion()
}

// GetLatestVersion 从 GitHub 获取最新版本
func (u *Updater) GetLatestVersion() (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", githubAPIBase, repoOwner, repoName)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := u.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("获取最新版本失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("获取最新版本失败: HTTP %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("解析版本信息失败: %w", err)
	}

	return release.TagName, nil
}

// CompareVersion 比较版本号，如果 latest > current 返回 true
// 支持 v1.0.0 和 1.0.0 格式
func (u *Updater) CompareVersion(current, latest string) bool {
	currentClean := strings.TrimPrefix(current, "v")
	latestClean := strings.TrimPrefix(latest, "v")

	currentParts := parseVersionParts(currentClean)
	latestParts := parseVersionParts(latestClean)

	// 补齐版本号长度
	maxLen := len(currentParts)
	if len(latestParts) > maxLen {
		maxLen = len(latestParts)
	}

	for len(currentParts) < maxLen {
		currentParts = append(currentParts, 0)
	}
	for len(latestParts) < maxLen {
		latestParts = append(latestParts, 0)
	}

	// 比较版本号
	for i := 0; i < maxLen; i++ {
		if latestParts[i] > currentParts[i] {
			return true
		} else if latestParts[i] < currentParts[i] {
			return false
		}
	}

	return false
}

// parseVersionParts 解析版本号为数字数组
func parseVersionParts(version string) []int {
	parts := strings.Split(version, ".")
	result := make([]int, 0, len(parts))

	for _, part := range parts {
		if num, err := strconv.Atoi(part); err == nil {
			result = append(result, num)
		}
	}

	return result
}

// GetDownloadURL 获取适合当前平台的下载 URL
func (u *Updater) GetDownloadURL(tagVersion string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/tags/%s", githubAPIBase, repoOwner, repoName, tagVersion)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := u.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("获取发布信息失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("获取发布信息失败: HTTP %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("解析发布信息失败: %w", err)
	}

	// 获取目标平台对应的文件名关键字
	targetAsset, err := u.getTargetAssetName()
	if err != nil {
		return "", err
	}

	// 查找匹配的 asset
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, targetAsset) {
			return asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("未找到适合当前平台的下载文件 (需要包含: %s)", targetAsset)
}

// getTargetAssetName 获取目标平台对应的文件名关键字
// 返回值需要匹配 release.yml 中的 artifact_name
func (u *Updater) getTargetAssetName() (string, error) {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	switch {
	case osName == "darwin" && arch == "arm64":
		return "macos_arm64", nil
	case osName == "darwin" && arch == "amd64":
		return "macos_x86_64", nil
	case osName == "linux" && arch == "amd64":
		return "linux_x86_64", nil
	case osName == "linux" && arch == "arm64":
		return "linux_arm64", nil
	case osName == "windows" && arch == "amd64":
		return "windows_x86_64", nil
	default:
		return "", fmt.Errorf("不支持的平台: %s/%s", osName, arch)
	}
}

// DownloadAndUpdate 下载并更新程序
func (u *Updater) DownloadAndUpdate(downloadURL string) error {
	fmt.Println("正在下载最新版本...")

	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		return fmt.Errorf("创建下载请求失败: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := u.client.Do(req)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败: HTTP %d", resp.StatusCode)
	}

	// 读取文件内容
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取下载内容失败: %w", err)
	}

	// 获取当前可执行文件路径
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %w", err)
	}
	currentExe, err = filepath.EvalSymlinks(currentExe)
	if err != nil {
		return fmt.Errorf("解析可执行文件路径失败: %w", err)
	}

	// 临时文件路径
	tempPath := currentExe + ".tmp"
	backupPath := currentExe + ".bak"

	// 写入临时文件（可执行文件需要 0755 权限）
	//nolint:gosec // G306: 可执行文件需要 0755 权限
	if err := os.WriteFile(tempPath, data, 0755); err != nil {
		return fmt.Errorf("写入临时文件失败: %w", err)
	}

	// 备份原文件（如果备份文件已存在则先删除）
	if _, err := os.Stat(backupPath); err == nil {
		if err := os.Remove(backupPath); err != nil {
			return fmt.Errorf("删除旧备份文件失败: %w", err)
		}
	}

	if err := os.Rename(currentExe, backupPath); err != nil {
		// 清理临时文件
		_ = os.Remove(tempPath)
		return fmt.Errorf("备份原文件失败: %w", err)
	}

	// 替换为新文件
	if err := os.Rename(tempPath, currentExe); err != nil {
		// 尝试恢复原文件
		_ = os.Rename(backupPath, currentExe)
		return fmt.Errorf("替换可执行文件失败: %w", err)
	}

	fmt.Println("更新完成！")
	fmt.Printf("原文件已备份为: %s\n", backupPath)
	fmt.Printf("新文件路径: %s\n", currentExe)

	return nil
}

// CheckUpdate 检查是否有新版本，返回最新版本号和是否需要更新
func (u *Updater) CheckUpdate() (string, bool, error) {
	currentVersion := u.GetCurrentVersion()
	latestVersion, err := u.GetLatestVersion()
	if err != nil {
		return "", false, err
	}

	hasUpdate := u.CompareVersion(currentVersion, latestVersion)
	return latestVersion, hasUpdate, nil
}

// Update 执行完整的更新流程
func (u *Updater) Update() error {
	fmt.Println("检查更新...")

	currentVersion := u.GetCurrentVersion()
	fmt.Printf("当前版本: %s\n", currentVersion)

	latestVersion, err := u.GetLatestVersion()
	if err != nil {
		return err
	}

	displayVersion := strings.TrimPrefix(latestVersion, "v")
	fmt.Printf("最新版本: %s\n", displayVersion)

	if !u.CompareVersion(currentVersion, latestVersion) {
		fmt.Println("当前已是最新版本！")
		return nil
	}

	fmt.Printf("发现新版本: %s\n", displayVersion)

	downloadURL, err := u.GetDownloadURL(latestVersion)
	if err != nil {
		return err
	}

	return u.DownloadAndUpdate(downloadURL)
}
