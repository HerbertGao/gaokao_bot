package database

import (
	"fmt"
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

// AutoMigrateSchema 自动迁移数据库表结构
// GORM 的 AutoMigrate 是幂等的，可以安全地多次执行
func AutoMigrateSchema(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.ExamDate{},
		&model.SendChat{},
		&model.UserTemplate{},
	)
}
