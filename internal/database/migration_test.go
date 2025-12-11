package database

import (
	"testing"

	"github.com/herbertgao/gaokao_bot/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	return db
}

func TestRunMigrations(t *testing.T) {
	db := setupTestDB(t)

	// 运行迁移
	err := RunMigrations(db)
	if err != nil {
		t.Errorf("RunMigrations() error = %v", err)
	}

	// 验证 migrations 表已创建
	var migrationCount int64
	err = db.Model(&Migration{}).Count(&migrationCount).Error
	if err != nil {
		t.Errorf("Failed to count migrations: %v", err)
	}

	if migrationCount == 0 {
		t.Error("Expected at least one migration record")
	}

	// 验证表已创建
	if !db.Migrator().HasTable(&model.ExamDate{}) {
		t.Error("ExamDate table not created")
	}

	if !db.Migrator().HasTable(&model.SendChat{}) {
		t.Error("SendChat table not created")
	}

	if !db.Migrator().HasTable(&model.UserTemplate{}) {
		t.Error("UserTemplate table not created")
	}
}

func TestRunMigrations_Idempotent(t *testing.T) {
	db := setupTestDB(t)

	// 第一次运行迁移
	err := RunMigrations(db)
	if err != nil {
		t.Errorf("First RunMigrations() error = %v", err)
	}

	// 第二次运行迁移（应该幂等，不会重复应用）
	err = RunMigrations(db)
	if err != nil {
		t.Errorf("Second RunMigrations() error = %v", err)
	}

	// 验证迁移记录数量没有重复
	var migrationCount int64
	err = db.Model(&Migration{}).Count(&migrationCount).Error
	if err != nil {
		t.Errorf("Failed to count migrations: %v", err)
	}

	// 应该只有一条迁移记录（001_initial_schema）
	if migrationCount != 1 {
		t.Errorf("Expected 1 migration record, got %d", migrationCount)
	}
}

func TestRollbackLastMigration(t *testing.T) {
	db := setupTestDB(t)

	// 先运行迁移
	err := RunMigrations(db)
	if err != nil {
		t.Fatalf("RunMigrations() error = %v", err)
	}

	// 验证表已创建
	if !db.Migrator().HasTable(&model.ExamDate{}) {
		t.Error("ExamDate table not created before rollback")
	}

	// 回滚最后一个迁移
	err = RollbackLastMigration(db)
	if err != nil {
		t.Errorf("RollbackLastMigration() error = %v", err)
	}

	// 验证表已删除
	if db.Migrator().HasTable(&model.ExamDate{}) {
		t.Error("ExamDate table still exists after rollback")
	}

	if db.Migrator().HasTable(&model.SendChat{}) {
		t.Error("SendChat table still exists after rollback")
	}

	if db.Migrator().HasTable(&model.UserTemplate{}) {
		t.Error("UserTemplate table still exists after rollback")
	}

	// 验证迁移记录已删除
	var migrationCount int64
	err = db.Model(&Migration{}).Count(&migrationCount).Error
	if err != nil {
		t.Errorf("Failed to count migrations: %v", err)
	}

	if migrationCount != 0 {
		t.Errorf("Expected 0 migration records after rollback, got %d", migrationCount)
	}
}

func TestRollbackLastMigration_NoMigrations(t *testing.T) {
	db := setupTestDB(t)

	// 创建 migrations 表但不运行任何迁移
	err := db.AutoMigrate(&Migration{})
	if err != nil {
		t.Fatalf("Failed to create migrations table: %v", err)
	}

	// 尝试回滚（应该返回错误）
	err = RollbackLastMigration(db)
	if err == nil {
		t.Error("Expected error when rolling back with no migrations, got nil")
	}
}

func TestMigration001Up(t *testing.T) {
	db := setupTestDB(t)

	// 执行 migration001Up
	err := migration001Up(db)
	if err != nil {
		t.Errorf("migration001Up() error = %v", err)
	}

	// 验证所有表已创建
	tables := []interface{}{
		&model.ExamDate{},
		&model.SendChat{},
		&model.UserTemplate{},
	}

	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			t.Errorf("Table %T not created", table)
		}
	}
}

func TestMigration001Down(t *testing.T) {
	db := setupTestDB(t)

	// 先执行 Up 创建表
	err := migration001Up(db)
	if err != nil {
		t.Fatalf("migration001Up() error = %v", err)
	}

	// 执行 Down 删除表
	err = migration001Down(db)
	if err != nil {
		t.Errorf("migration001Down() error = %v", err)
	}

	// 验证所有表已删除
	tables := []interface{}{
		&model.ExamDate{},
		&model.SendChat{},
		&model.UserTemplate{},
	}

	for _, table := range tables {
		if db.Migrator().HasTable(table) {
			t.Errorf("Table %T still exists after migration down", table)
		}
	}
}

func TestMigrationVersionStructure(t *testing.T) {
	// 验证 migrations 切片不为空
	if len(migrations) == 0 {
		t.Error("migrations slice is empty")
	}

	// 验证每个迁移都有必需的字段
	for i, m := range migrations {
		if m.Version == "" {
			t.Errorf("Migration %d has empty Version", i)
		}

		if m.Up == nil {
			t.Errorf("Migration %d (%s) has nil Up function", i, m.Version)
		}

		if m.Down == nil {
			t.Errorf("Migration %d (%s) has nil Down function", i, m.Version)
		}
	}
}
