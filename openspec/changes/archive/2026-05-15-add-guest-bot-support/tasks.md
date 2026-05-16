## 1. 升级 telego 依赖

- [x] 1.1 将 `go.mod` 中 `mymmrac/telego` 升级到 `v1.9.0`，执行 `go mod tidy` 更新 `go.sum`
- [x] 1.2 运行 `go build ./...`，逐一修复因跨版本破坏性变更导致的编译错误（重点 `internal/bot`、`internal/api`、`internal/handler`、`internal/middleware`）
- [x] 1.3 运行 `make test`，修复因升级导致失败的现有测试，确认全部通过

## 2. 抽取共用倒计时文本函数

- [x] 2.1 在 `internal/service` 抽出共用函数（如 `BuildCountdownText(arg string, now time.Time) (string, error)`），实现「年份参数解析 → 按年份或当前范围查考试 → 套默认模板 → 拼接多考试为单条文本」
- [x] 2.2 重构 `MessageService.GetCountDownMessage` 改为调用该共用函数，保证 `/d` 输出行为不变
- [x] 2.3 更新/补充 `message_service_test.go` 及共用函数的单元测试，覆盖无参数、合法年份、非法参数、无数据四种分支

## 3. Guest 消息文本提取

- [x] 3.1 在 `internal/util` 新增/扩展函数：剥离 Guest 消息开头的 `@botname` 提及前缀，返回剩余参数文本；空文本与回复式召唤按无参数处理
- [x] 3.2 在 `telegram_bot_util_test.go` 补充测试，覆盖 `@botname`、`@botname 2026`、纯文本、空文本等输入

## 4. Guest 消息处理逻辑

- [x] 4.1 在 `BotService` 新增 `HandleGuestMessage` 方法：提取参数 → 调共用倒计时函数 → 构造单条 `InlineQueryResultArticle` → 调用 `bot.AnswerGuestQuery`（独立 `context.WithTimeout`，失败仅记日志不重试）
- [x] 4.2 确保 Guest 应答仅用默认模板，不含 `/template`、`/debug` 或 Mini App 按钮
- [x] 4.3 在 `bot_service_test.go` 补充测试，覆盖无参数、合法年份、非法参数、无数据，以及应答失败的日志路径

## 5. Bot 处理器注册

- [x] 5.1 在 `internal/bot/bot.go` 注册 `handler.HandleGuestMessage`（使用 `AnyGuestMessage()` 预测器），回调转交 `BotService.HandleGuestMessage`，并加 Debug 级日志
- [x] 5.2 更新 `bot_test.go`，验证 Guest 处理器已注册

## 6. 文档与验收

- [x] 6.1 在 README 或部署文档新增说明：Guest 模式需在 BotFather 的 Mini App 中开启 "Guest Mode" 开关
- [x] 6.2 运行 `make test` 与 golangci-lint，确认全部通过
- [ ] 6.3 在 dev 环境实测：从一个 Bot 非成员的群 @提及 Bot（含/不含年份参数）与回复 Bot 消息，确认应答正确（需真实 Bot token 与 BotFather 开启 Guest Mode，由用户手动验证）
