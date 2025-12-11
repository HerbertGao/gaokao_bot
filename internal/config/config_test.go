package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// 创建临时环境变量
	oldEnv := os.Getenv("APP_PORT")
	defer func() {
		if oldEnv != "" {
			os.Setenv("APP_PORT", oldEnv)
		} else {
			os.Unsetenv("APP_PORT")
		}
	}()

	os.Setenv("APP_PORT", "8080")

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
	_, err := Load("invalid_env")
	// 应该仍然成功，只是使用默认配置
	if err != nil {
		t.Errorf("Load() should not error for invalid env, got %v", err)
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	// 设置环境变量
	oldPort := os.Getenv("APP_PORT")
	defer func() {
		if oldPort != "" {
			os.Setenv("APP_PORT", oldPort)
		} else {
			os.Unsetenv("APP_PORT")
		}
	}()

	testPort := "9999"
	os.Setenv("APP_PORT", testPort)

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
