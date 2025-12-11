package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/herbertgao/gaokao_bot/internal/config"
	"github.com/herbertgao/gaokao_bot/internal/database"
)

func main() {
	// 解析命令行参数
	env := flag.String("env", "dev", "Environment (dev, test, prod)")
	action := flag.String("action", "up", "Migration action (up, down)")
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(*env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// 连接数据库
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting database connection: %v\n", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	// 执行迁移操作
	switch *action {
	case "up":
		fmt.Println("Running migrations...")
		if err := database.RunMigrations(db); err != nil {
			fmt.Fprintf(os.Stderr, "Migration failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ All migrations completed successfully")

	case "down":
		fmt.Println("Rolling back last migration...")
		if err := database.RollbackLastMigration(db); err != nil {
			fmt.Fprintf(os.Stderr, "Rollback failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Rollback completed successfully")

	default:
		fmt.Fprintf(os.Stderr, "Unknown action: %s (use 'up' or 'down')\n", *action)
		os.Exit(1)
	}
}
