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
	cfg, err := config.Load(*env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置错误: %v\n", err)
		os.Exit(1)
	}

	// 连接数据库
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "连接数据库错误: %v\n", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取数据库连接错误: %v\n", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	// 执行迁移操作
	switch *action {
	case "up":
		fmt.Println("正在执行数据库迁移...")
		if err := database.RunMigrations(db); err != nil {
			fmt.Fprintf(os.Stderr, "数据库迁移失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ 所有迁移已成功完成")

	case "down":
		fmt.Println("正在回滚最后一次迁移...")
		if err := database.RollbackLastMigration(db); err != nil {
			fmt.Fprintf(os.Stderr, "回滚失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ 回滚已成功完成")

	default:
		fmt.Fprintf(os.Stderr, "未知操作: %s (请使用 'up' 或 'down')\n", *action)
		os.Exit(1)
	}
}
