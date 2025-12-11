package database

import (
	"fmt"
	"time"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"gorm.io/gorm"
)

// Migration 数据库迁移记录
type Migration struct {
	ID        uint      `gorm:"primaryKey"`
	Version   string    `gorm:"uniqueIndex;not null"`
	AppliedAt time.Time `gorm:"not null"`
}

// MigrationFunc 迁移函数类型
type MigrationFunc func(*gorm.DB) error

// MigrationVersion 迁移版本
type MigrationVersion struct {
	Version string
	Up      MigrationFunc
	Down    MigrationFunc
}

// migrations 所有迁移版本（按顺序）
var migrations = []MigrationVersion{
	{
		Version: "001_initial_schema",
		Up:      migration001Up,
		Down:    migration001Down,
	},
}

// RunMigrations 运行所有待执行的迁移
func RunMigrations(db *gorm.DB) error {
	// 创建 migrations 表
	if err := db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// 获取已应用的迁移
	var appliedMigrations []Migration
	if err := db.Find(&appliedMigrations).Error; err != nil {
		return fmt.Errorf("failed to fetch applied migrations: %w", err)
	}

	appliedVersions := make(map[string]bool)
	for _, m := range appliedMigrations {
		appliedVersions[m.Version] = true
	}

	// 应用待执行的迁移
	for _, migration := range migrations {
		if appliedVersions[migration.Version] {
			continue
		}

		// 执行迁移
		if err := migration.Up(db); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}

		// 记录迁移
		record := Migration{
			Version:   migration.Version,
			AppliedAt: time.Now(),
		}
		if err := db.Create(&record).Error; err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
		}

		fmt.Printf("✓ Applied migration: %s\n", migration.Version)
	}

	return nil
}

// RollbackLastMigration 回滚最后一个迁移
func RollbackLastMigration(db *gorm.DB) error {
	// 获取最后一个迁移
	var lastMigration Migration
	if err := db.Order("id DESC").First(&lastMigration).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("no migrations to rollback")
		}
		return fmt.Errorf("failed to fetch last migration: %w", err)
	}

	// 找到对应的迁移版本
	var targetMigration *MigrationVersion
	for _, m := range migrations {
		if m.Version == lastMigration.Version {
			targetMigration = &m
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("migration version %s not found", lastMigration.Version)
	}

	// 执行回滚
	if err := targetMigration.Down(db); err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", lastMigration.Version, err)
	}

	// 删除迁移记录
	if err := db.Delete(&lastMigration).Error; err != nil {
		return fmt.Errorf("failed to delete migration record %s: %w", lastMigration.Version, err)
	}

	fmt.Printf("✓ Rolled back migration: %s\n", lastMigration.Version)
	return nil
}

// migration001Up 初始化数据库架构
func migration001Up(db *gorm.DB) error {
	// 自动迁移所有模型
	return db.AutoMigrate(
		&model.ExamDate{},
		&model.SendChat{},
		&model.UserTemplate{},
	)
}

// migration001Down 回滚初始化架构
func migration001Down(db *gorm.DB) error {
	// 删除所有表（按依赖顺序反向）
	return db.Migrator().DropTable(
		&model.UserTemplate{},
		&model.SendChat{},
		&model.ExamDate{},
	)
}
