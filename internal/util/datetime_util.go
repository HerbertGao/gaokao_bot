package util

import (
	"fmt"
	"strings"
	"time"
)

const (
	// RoundingThresholdSeconds 时间标准化的四舍五入阈值（秒）
	// 当秒数 >= 此阈值时，进位到下一分钟；否则保持当前分钟
	//
	// 设计原因：
	// 1. 符合数学上的四舍五入规则（30 秒是 1 分钟的一半）
	// 2. 防止倒计时显示出现 "3天23小时59分钟59秒" 等不美观的情况
	// 3. 与 cron 定时任务的 1 分钟粒度相匹配：
	//    - cron 任务每分钟触发一次（如 "0 * * * * *" 表示每分钟的第 0 秒触发）
	//    - 实际触发时间可能在 0-59 秒之间（受系统负载影响）
	//    - 30 秒的阈值确保大部分触发都能标准化到期望的整分钟
	// 4. 对于每小时或每天触发的任务，30 秒的误差几乎可以忽略不计
	//
	// 如果未来需要支持更高精度的定时任务（如每秒触发），可能需要调整此阈值或使用不同的标准化策略
	RoundingThresholdSeconds = 30
)

// bjtLocation 北京时间时区（UTC+8）
// 使用包级变量缓存，避免重复加载
var bjtLocation *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 如果无法加载时区（罕见情况，如系统缺少 tzdata），使用固定偏移量
		loc = time.FixedZone("BJT", 8*3600)
	}
	bjtLocation = loc
}

// GetBJTLocation 返回北京时区（UTC+8）
// 用于需要明确使用北京时区的场景（如 cron 定时任务）
func GetBJTLocation() *time.Location {
	return bjtLocation
}

// NowBJT 返回当前北京时间（UTC+8）
// 用于高考倒计时等需要明确使用北京时间的场景
func NowBJT() time.Time {
	return time.Now().In(bjtLocation)
}

// FormatNormal 格式化为标准格式
func FormatNormal(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// NormalizeToMinute 标准化时间到最近的整分钟
// 四舍五入：秒 >= RoundingThresholdSeconds 进位到下一分钟，< RoundingThresholdSeconds 保持当前分钟
// 用于定时任务的倒计时显示，避免出现"3天23小时59分钟59秒"等情况
//
// 注意：此函数保留输入时间的时区（通过 t.Location()）
// 如果需要标准化到北京时间，请先使用 NowBJT() 或 time.In(bjtLocation) 转换时区
func NormalizeToMinute(t time.Time) time.Time {
	if t.Second() >= RoundingThresholdSeconds {
		t = t.Add(time.Minute)
	}
	// 截断到整分钟（保留原时区）
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}

// FormatDuration 格式化时间间隔为中文描述
// 格式：天时分钟秒，如果某项为0则不显示
// 例如：350天23小时59分钟59秒、18天3分钟
//
// 注意：如果总时长小于1秒但大于0，会显示为"1秒"
func FormatDuration(d time.Duration) string {
	// 如果为负数，返回 0秒
	if d < 0 {
		return "0秒"
	}

	totalSeconds := int64(d.Seconds())
	nanoSeconds := d.Nanoseconds() % 1_000_000_000

	// 如果有纳秒部分且总时长小于1秒，则显示为1秒
	if totalSeconds == 0 && nanoSeconds > 0 {
		return "1秒"
	}

	// 如果总时长为0，返回0秒
	if totalSeconds == 0 && nanoSeconds == 0 {
		return "0秒"
	}

	days := totalSeconds / 86400      // 24 * 60 * 60
	hours := (totalSeconds % 86400) / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	var result strings.Builder

	if days > 0 {
		result.WriteString(fmt.Sprintf("%d天", days))
	}
	if hours > 0 {
		result.WriteString(fmt.Sprintf("%d小时", hours))
	}
	if minutes > 0 {
		result.WriteString(fmt.Sprintf("%d分钟", minutes))
	}
	if seconds > 0 || result.Len() == 0 {
		result.WriteString(fmt.Sprintf("%d秒", seconds))
	}

	return result.String()
}

// FormatDurationWithMillis 格式化时间间隔为中文描述（精确到毫秒）
// 格式：天时分钟秒毫秒，如果某项为0则不显示
// 例如：350天23小时59分钟59秒500毫秒、18天3分钟200毫秒
func FormatDurationWithMillis(d time.Duration) string {
	// 如果为负数，返回 0秒
	if d < 0 {
		return "0秒"
	}

	totalSeconds := int64(d.Seconds())
	nanoSeconds := d.Nanoseconds() % 1_000_000_000
	millis := nanoSeconds / 1_000_000 // 转换为毫秒

	// 如果总时长为0，返回0秒
	if totalSeconds == 0 && millis == 0 {
		return "0秒"
	}

	days := totalSeconds / 86400      // 24 * 60 * 60
	hours := (totalSeconds % 86400) / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	var result strings.Builder

	if days > 0 {
		result.WriteString(fmt.Sprintf("%d天", days))
	}
	if hours > 0 {
		result.WriteString(fmt.Sprintf("%d小时", hours))
	}
	if minutes > 0 {
		result.WriteString(fmt.Sprintf("%d分钟", minutes))
	}
	if seconds > 0 {
		result.WriteString(fmt.Sprintf("%d秒", seconds))
	}
	if millis > 0 {
		result.WriteString(fmt.Sprintf("%d毫秒", millis))
	}

	// 如果没有任何时间单位，说明只有毫秒
	if result.Len() == 0 {
		result.WriteString(fmt.Sprintf("%d毫秒", millis))
	}

	return result.String()
}

// DaysBetween 计算两个时间之间的天数
func DaysBetween(start, end time.Time) int {
	return int(end.Sub(start).Hours() / 24)
}

// HoursBetween 计算两个时间之间的小时数
func HoursBetween(start, end time.Time) int {
	return int(end.Sub(start).Hours())
}

// StartOfDay 获取一天的开始时间
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay 获取一天的结束时间
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}