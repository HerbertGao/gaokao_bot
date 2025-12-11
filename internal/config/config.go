package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	App       AppConfig
	Telegram  TelegramConfig
	Database  DatabaseConfig
	Snowflake SnowflakeConfig
	Log       LogConfig
	Task      TaskConfig
}

// AppConfig 应用配置
type AppConfig struct {
	Name string
	Env  string
	Port int
}

// TelegramConfig Telegram 配置
type TelegramConfig struct {
	Bot     BotConfig
	MiniApp MiniAppConfig
}

// BotConfig Bot 配置
type BotConfig struct {
	Username string
	Token    string
}

// MiniAppConfig Mini App 配置
type MiniAppConfig struct {
	URL string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string
	Port            int
	Name            string
	Username        string
	Password        string
	Charset         string
	ParseTime       bool
	Loc             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

// SnowflakeConfig Snowflake ID 配置
type SnowflakeConfig struct {
	DatacenterID int64
	MachineID    int64
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string
	Format string
	Output string
}

// TaskConfig 任务配置
type TaskConfig struct {
	DailySend DailySendConfig
}

// DailySendConfig 每日发送任务配置
type DailySendConfig struct {
	Enabled bool
	Cron    string
}

// Load 加载配置
func Load(env string) (*Config, error) {
	// 尝试加载环境特定的 .env 文件
	if env != "" {
		envFile := fmt.Sprintf(".env.%s", env)
		_ = godotenv.Load(envFile)
	}

	// 加载默认 .env 文件
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "gaokao_bot"),
			Env:  getEnv("APP_ENV", env),
			Port: getEnvAsInt("APP_PORT", 8080),
		},
		Telegram: TelegramConfig{
			Bot: BotConfig{
				Username: getEnv("TELEGRAM_BOT_USERNAME", ""),
				Token:    getEnv("TELEGRAM_BOT_TOKEN", ""),
			},
			MiniApp: MiniAppConfig{
				URL: getEnv("TELEGRAM_MINIAPP_URL", ""),
			},
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "127.0.0.1"),
			Port:            getEnvAsInt("DB_PORT", 3306),
			Name:            getEnv("DB_NAME", "gaokao"),
			Username:        getEnv("DB_USERNAME", "root"),
			Password:        getEnv("DB_PASSWORD", ""),
			Charset:         getEnv("DB_CHARSET", "utf8mb4"),
			ParseTime:       getEnvAsBool("DB_PARSE_TIME", true),
			Loc:             getEnv("DB_LOC", "Asia/Shanghai"),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 20),
			ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME", 900),
		},
		Snowflake: SnowflakeConfig{
			DatacenterID: getEnvAsInt64("SNOWFLAKE_DATACENTER_ID", 1),
			MachineID:    getEnvAsInt64("SNOWFLAKE_MACHINE_ID", 1),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "text"),
			Output: getEnv("LOG_OUTPUT", "stdout"),
		},
		Task: TaskConfig{
			DailySend: DailySendConfig{
				Enabled: getEnvAsBool("TASK_DAILY_SEND_ENABLED", true),
				Cron:    getEnv("TASK_DAILY_SEND_CRON", "0 0 * * * *"),
			},
		},
	}

	return cfg, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为int，如果不存在或转换失败则返回默认值
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsInt64 获取环境变量并转换为int64，如果不存在或转换失败则返回默认值
func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool 获取环境变量并转换为bool，如果不存在或转换失败则返回默认值
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}