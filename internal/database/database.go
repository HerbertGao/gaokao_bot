package database

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/config"
	"github.com/herbertgao/gaokao_bot/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDatabase 创建数据库连接
func NewDatabase(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	// 对时区和密码进行URL编码
	loc := url.QueryEscape(cfg.Loc)
	password := url.QueryEscape(cfg.Password)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.Username,
		password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Charset,
		cfg.ParseTime,
		loc,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	return db, nil
}

// NewDatabaseWithRetry 创建数据库连接（带重试逻辑）
// 适用于容器化环境中数据库可能暂时不可用的情况
func NewDatabaseWithRetry(cfg *config.DatabaseConfig, maxRetries int) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// 指数退避配置
	initialDelay := 1 * time.Second
	maxDelay := 30 * time.Second
	currentDelay := initialDelay

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err = NewDatabase(cfg)
		if err == nil {
			// 连接成功，验证数据库是否可用
			sqlDB, pingErr := db.DB()
			if pingErr == nil {
				if pingErr = sqlDB.Ping(); pingErr == nil {
					log.Printf("数据库连接成功（第 %d 次尝试）", attempt)
					return db, nil
				}
				err = pingErr
			} else {
				err = pingErr
			}
		}

		if attempt < maxRetries {
			log.Printf("数据库连接失败（第 %d/%d 次尝试）: %v，%v 后重试...",
				attempt, maxRetries, err, currentDelay)
			time.Sleep(currentDelay)

			// 指数退避：每次失败后延迟时间翻倍
			currentDelay *= 2
			if currentDelay > maxDelay {
				currentDelay = maxDelay
			}
		}
	}

	return nil, fmt.Errorf("数据库连接失败（已重试 %d 次）: %w", maxRetries, err)
}

// AutoMigrateSchema 自动迁移数据库表结构
// GORM 的 AutoMigrate 是幂等的，可以安全地多次执行
func AutoMigrateSchema(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.ExamDate{},
		&model.SendChat{},
		&model.UserTemplate{},
	)
}
