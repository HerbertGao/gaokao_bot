# 高考倒计时 Bot (Go 版本)

基于 Java 原项目的完整功能复刻，使用 Go 语言实现的高考倒计时 Telegram Bot。

## 项目特点

- **现代化技术栈**: 使用 Go 1.21+ 和最新的 Telegram Bot SDK (telego v1.3.2)
- **清晰的分层架构**: Model -> Repository -> Service -> Bot
- **完整功能**: 倒计时查询、内联查询、定时推送、Mini App 支持
- **多环境支持**: 开发、测试、生产环境配置分离，自动切换 Telegram 测试服务器
- **高性能**: Go 原生并发，低内存占用

## 快速开始

### 环境要求

- Go 1.21 或更高版本
- MySQL 8.0+
- Telegram Bot Token

### 安装依赖

```bash
go mod download
```

### 配置

1. 复制环境变量模板：
```bash
cp configs/.env.example .env
```

2. 编辑 `.env` 文件，填入你的配置：
```env
BOT_USERNAME=your_bot_username
BOT_TOKEN=your_bot_token
DB_PASSWORD=your_db_password
```

3. 初始化数据库：
```bash
mysql -u root -p < sql/init.sql
```

### 运行

开发环境：
```bash
go run cmd/gaokao_bot/main.go -env=dev
```

测试环境（自动使用 Telegram 测试服务器）：
```bash
go run cmd/gaokao_bot/main.go -env=test
```

生产环境：
```bash
go run cmd/gaokao_bot/main.go -env=prod
```

> 💡 **自动测试服务器切换**: 当使用非 `prod` 环境时（如 `dev` 或 `test`），程序会自动使用 Telegram 测试服务器 API（`/bot<token>/test/`），无需手动配置。

### 构建

```bash
./scripts/build.sh
```

编译后的二进制文件位于 `bin/gaokao_bot`

## 项目结构

```
gaokao_bot/
├── cmd/gaokao_bot/          # 主程序入口
├── internal/                # 私有代码
│   ├── bot/                # Bot 核心
│   ├── service/            # 业务服务
│   ├── model/              # 数据模型
│   ├── repository/         # 数据访问
│   ├── task/               # 定时任务
│   ├── util/               # 工具函数
│   ├── config/             # 配置管理
│   └── database/           # 数据库连接
├── pkg/constant/           # 公共常量
├── configs/                # 配置文件
├── sql/                    # 数据库脚本
└── scripts/                # 构建脚本
```

## 支持的命令

- `/d [year]` - 获取高考倒计时（默认当年）
- `/ls` - 列出自定义模板（私聊）
- 内联查询：`@bot_username [year]` - 在任何聊天中快速查询

## 技术栈

- **Go**: 1.21+
- **Telegram Bot SDK**: github.com/mymmrac/telego v1.3.2
- **ORM**: GORM
- **数据库**: MySQL 8.0+
- **配置管理**: Viper
- **日志**: Logrus
- **定时任务**: robfig/cron
- **ID 生成**: Snowflake

## Mini App

本项目包含 Telegram Mini App 前端，用于可视化管理倒计时模板。

- **前端项目**: `/Users/herbertgao/WebstormProjects/gaokao-mini-app`
- **快速开始**: [5分钟快速测试指南](docs/QUICK_START_TESTING.md)
- **API 文档**: [Mini App API 参考](docs/MINI_APP_API.md)

## 开发与测试

### Telegram 测试环境

开发和测试 Mini App 时，建议使用 Telegram 官方测试环境：

- ✅ **自动切换**: 使用 `env=test` 或 `env=dev` 时自动使用测试服务器
- ✅ **数据隔离**: 测试数据不会污染生产环境
- ✅ **本地调试**: 支持 HTTP 和 localhost URL
- ✅ **真实验证**: 获取真实的 Init Data 进行身份验证

**快速设置**（5分钟）:
1. 登录 Telegram 测试环境（Desktop: Shift+Alt+右键 Add Account > Test Server）
2. 在测试环境创建 Bot（@BotFather）
3. 配置 `configs/config.test.yaml`
4. 启动：`./bin/gaokao_bot -env=test`

详细指南：
- 📖 [快速开始指南](docs/QUICK_START_TESTING.md) - 5分钟完成设置
- 📚 [Telegram 测试服务器集成](docs/TELEGRAM_TEST_SERVER.md) - 技术实现细节
- 🎨 [Mini App 设计文档](docs/MINI_APP_DESIGN.md) - 架构和设计

## 开发文档

- [Java 项目分析](docs/JAVA_PROJECT_ANALYSIS.md)
- [Go 项目设计](docs/GO_PROJECT_DESIGN.md)
- [Telegram 测试服务器集成](docs/TELEGRAM_TEST_SERVER.md)
- [快速开始测试指南](docs/QUICK_START_TESTING.md)
- [Mini App API 文档](docs/MINI_APP_API.md)
- [Mini App 设计文档](docs/MINI_APP_DESIGN.md)

## License

MIT License
