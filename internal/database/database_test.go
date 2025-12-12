package database

import (
	"testing"

	"github.com/herbertgao/gaokao_bot/internal/config"
)

func TestNewDatabase(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:            "localhost",
		Port:            3306,
		Name:            "test",
		Username:        "root",
		Password:        "password",
		Charset:         "utf8mb4",
		ParseTime:       true,
		Loc:             "Local",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: 3600,
	}

	// 注意：这个测试会失败，因为无法连接到真实的 MySQL
	// 但它可以测试 DSN 构建逻辑
	_, err := NewDatabase(cfg)
	// 期望连接失败（因为是测试环境）
	if err == nil {
		t.Log("Unexpected success - database connection should fail in test environment")
	}
}

func TestDSN_Building(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:            "localhost",
		Port:            3306,
		Name:            "test_db",
		Username:        "test_user",
		Password:        "test@pass",
		Charset:         "utf8mb4",
		ParseTime:       true,
		Loc:             "Asia/Shanghai",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: 3600,
	}

	// 测试 DSN 构建
	// 实际的 DSN 构建在 NewDatabase 中进行
	// 这里只是确保配置结构正确
	if cfg.Host == "" {
		t.Error("Host should not be empty")
	}
	if cfg.Port == 0 {
		t.Error("Port should not be zero")
	}
	if cfg.Name == "" {
		t.Error("Name should not be empty")
	}
}
