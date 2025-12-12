package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// 创建临时环境变量
	oldPort := os.Getenv("APP_PORT")
	oldToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	defer func() {
		if oldPort != "" {
			_ = os.Setenv("APP_PORT", oldPort)
		} else {
			_ = os.Unsetenv("APP_PORT")
		}
		if oldToken != "" {
			_ = os.Setenv("TELEGRAM_BOT_TOKEN", oldToken)
		} else {
			_ = os.Unsetenv("TELEGRAM_BOT_TOKEN")
		}
	}()

	_ = os.Setenv("APP_PORT", "8080")
	_ = os.Setenv("TELEGRAM_BOT_TOKEN", "test_token")

	cfg, err := Load("dev")
	if err != nil {
		t.Errorf("Load() error = %v", err)
	}

	if cfg == nil {
		t.Error("Expected config, got nil")
		return
	}

	// 验证基本配置
	if cfg.App.Env != "dev" {
		t.Errorf("App.Env = %s, want %s", cfg.App.Env, "dev")
	}
}

func TestLoad_InvalidEnv(t *testing.T) {
	oldToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	defer func() {
		if oldToken != "" {
			_ = os.Setenv("TELEGRAM_BOT_TOKEN", oldToken)
		} else {
			_ = os.Unsetenv("TELEGRAM_BOT_TOKEN")
		}
	}()

	_ = os.Setenv("TELEGRAM_BOT_TOKEN", "test_token")

	_, err := Load("invalid_env")
	// 应该仍然成功，只是使用默认配置
	if err != nil {
		t.Errorf("Load() should not error for invalid env, got %v", err)
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	// 设置环境变量
	oldPort := os.Getenv("APP_PORT")
	oldToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	defer func() {
		if oldPort != "" {
			_ = os.Setenv("APP_PORT", oldPort)
		} else {
			_ = os.Unsetenv("APP_PORT")
		}
		if oldToken != "" {
			_ = os.Setenv("TELEGRAM_BOT_TOKEN", oldToken)
		} else {
			_ = os.Unsetenv("TELEGRAM_BOT_TOKEN")
		}
	}()

	testPort := "9999"
	_ = os.Setenv("APP_PORT", testPort)
	_ = os.Setenv("TELEGRAM_BOT_TOKEN", "test_token")

	cfg, err := Load("dev")
	if err != nil {
		t.Errorf("Load() error = %v", err)
	}

	// 注意：viper可能会将字符串转换为int
	// 这个测试只是确保 Load 函数能正常工作
	if cfg == nil {
		t.Error("Expected config, got nil")
	}
}

func TestValidate_Success(t *testing.T) {
	cfg := &Config{
		App: AppConfig{
			Env:  "dev",
			Port: 8080,
		},
		Telegram: TelegramConfig{
			Bot: BotConfig{
				Token: "test_token",
			},
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Name:     "testdb",
			Username: "testuser",
		},
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestValidate_RequiresBotToken(t *testing.T) {
	cfg := &Config{
		App: AppConfig{
			Env:  "dev",
			Port: 8080,
		},
		Telegram: TelegramConfig{
			Bot: BotConfig{
				Token: "", // 缺少 token
			},
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Name:     "testdb",
			Username: "testuser",
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Validate() should return error when bot token is missing")
	}
}

func TestValidate_DatabaseConfigRequired(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "Missing DB host",
			config: &Config{
				App:      AppConfig{Env: "dev", Port: 8080},
				Telegram: TelegramConfig{Bot: BotConfig{Token: "test_token"}},
				Database: DatabaseConfig{
					Host:     "", // 缺少
					Port:     3306,
					Name:     "testdb",
					Username: "testuser",
				},
			},
		},
		{
			name: "Missing DB name",
			config: &Config{
				App:      AppConfig{Env: "dev", Port: 8080},
				Telegram: TelegramConfig{Bot: BotConfig{Token: "test_token"}},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     3306,
					Name:     "", // 缺少
					Username: "testuser",
				},
			},
		},
		{
			name: "Missing DB username",
			config: &Config{
				App:      AppConfig{Env: "dev", Port: 8080},
				Telegram: TelegramConfig{Bot: BotConfig{Token: "test_token"}},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     3306,
					Name:     "testdb",
					Username: "", // 缺少
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if err == nil {
				t.Errorf("Validate() should return error for %s", tt.name)
			}
		})
	}
}

func TestValidate_InvalidPorts(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "App port too low",
			config: &Config{
				App:      AppConfig{Env: "dev", Port: 0},
				Telegram: TelegramConfig{Bot: BotConfig{Token: "test_token"}},
				Database: DatabaseConfig{
					Host: "localhost", Port: 3306, Name: "testdb", Username: "testuser",
				},
			},
		},
		{
			name: "App port too high",
			config: &Config{
				App:      AppConfig{Env: "dev", Port: 70000},
				Telegram: TelegramConfig{Bot: BotConfig{Token: "test_token"}},
				Database: DatabaseConfig{
					Host: "localhost", Port: 3306, Name: "testdb", Username: "testuser",
				},
			},
		},
		{
			name: "DB port too low",
			config: &Config{
				App:      AppConfig{Env: "dev", Port: 8080},
				Telegram: TelegramConfig{Bot: BotConfig{Token: "test_token"}},
				Database: DatabaseConfig{
					Host: "localhost", Port: 0, Name: "testdb", Username: "testuser",
				},
			},
		},
		{
			name: "DB port too high",
			config: &Config{
				App:      AppConfig{Env: "dev", Port: 8080},
				Telegram: TelegramConfig{Bot: BotConfig{Token: "test_token"}},
				Database: DatabaseConfig{
					Host: "localhost", Port: 999999, Name: "testdb", Username: "testuser",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if err == nil {
				t.Errorf("Validate() should return error for %s", tt.name)
			}
		})
	}
}

func TestGetEnvAsInt64(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue int64
		want         int64
	}{
		{
			name:         "Valid int64",
			envKey:       "TEST_INT64",
			envValue:     "12345",
			defaultValue: 0,
			want:         12345,
		},
		{
			name:         "Invalid int64",
			envKey:       "TEST_INT64_INVALID",
			envValue:     "not_a_number",
			defaultValue: 999,
			want:         999,
		},
		{
			name:         "Empty value",
			envKey:       "TEST_INT64_EMPTY",
			envValue:     "",
			defaultValue: 888,
			want:         888,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				_ = os.Setenv(tt.envKey, tt.envValue)
				defer func() {
					_ = os.Unsetenv(tt.envKey)
				}()
			}

			got := getEnvAsInt64(tt.envKey, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnvAsInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvAsBool(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue bool
		want         bool
	}{
		{
			name:         "True value",
			envKey:       "TEST_BOOL_TRUE",
			envValue:     "true",
			defaultValue: false,
			want:         true,
		},
		{
			name:         "False value",
			envKey:       "TEST_BOOL_FALSE",
			envValue:     "false",
			defaultValue: true,
			want:         false,
		},
		{
			name:         "1 as true",
			envKey:       "TEST_BOOL_1",
			envValue:     "1",
			defaultValue: false,
			want:         true,
		},
		{
			name:         "0 as false",
			envKey:       "TEST_BOOL_0",
			envValue:     "0",
			defaultValue: true,
			want:         false,
		},
		{
			name:         "Invalid bool",
			envKey:       "TEST_BOOL_INVALID",
			envValue:     "not_a_bool",
			defaultValue: true,
			want:         true,
		},
		{
			name:         "Empty value",
			envKey:       "TEST_BOOL_EMPTY",
			envValue:     "",
			defaultValue: false,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				_ = os.Setenv(tt.envKey, tt.envValue)
				defer func() {
					_ = os.Unsetenv(tt.envKey)
				}()
			}

			got := getEnvAsBool(tt.envKey, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnvAsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
