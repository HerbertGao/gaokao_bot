## 为什么

目前用户要在某个群里查询高考倒计时，必须先把 @gaokao_bot 拉进群，或切换到与 Bot 的私聊。Telegram 在 2025 年推出的 **Guest Bots** 特性允许 Bot 在它并非成员的任意群聊/私聊中被临时召唤（通过 @提及 或回复其消息）并一次性应答。为本 Bot 接入该特性，可以让用户在任何聊天里直接召唤倒计时，无需改变群成员构成，显著降低使用门槛。

## 变更内容

- **新增** Guest 模式消息处理：接收 Telegram 推送的 `guest_message` 更新，在 Bot 非成员的聊天中应答倒计时查询。
- **新增** Guest 专用文本提取：剥离开头的 `@botname` 提及，得到可选的年份参数；回复式召唤（无提及文本）按无参数处理。
- **复用** `/d` 命令的倒计时拼接逻辑：把当前时间范围内（或指定年份）的所有考试拼成单条文本，套用默认模板。
- **新增** 通过 `answerGuestQuery` 应答：响应类型为单个 `InlineQueryResult`，构造一条 `InlineQueryResultArticle`。
- **BREAKING** 依赖升级：`mymmrac/telego` 从 `v1.3.3` 升级到 `v1.9.0`（Guest Bot 相关类型与处理器自 v1.9.0 起提供），需审计跨版本的破坏性 API 变更。
- **新增** 部署文档说明：Guest 模式需在 BotFather 的 Mini App 中开启 "Guest Mode" 开关，纯代码无法启用。

### 非目标（Non-Goals）

- 不在 Guest 模式下暴露 `/template`、`/debug` 命令，也不打开 Mini App 按钮。
- 不在 Guest 模式下提供用户自定义模板（Guest 场景无持久用户身份），仅用默认模板。
- 不沿用 Inline Query 的多候选项形式（`answerGuestQuery` 只接受单条结果）。
- 不为 Guest 交互做任何持久化（不写 send_chat、不记录 Guest 聊天）。
- 不改动现有命令、Inline Query、定时推送、Mini App 后端的行为。

## 功能 (Capabilities)

### 新增功能
- `guest-countdown`: Bot 在非成员聊天中被 @提及或被回复时，应答一条默认模板的高考倒计时，支持可选的指定年份参数。

### 修改功能
<!-- 无现有规范的需求变更 -->

## 影响

- **依赖**：`go.mod` / `go.sum` 中 `mymmrac/telego` v1.3.3 → v1.9.0；需回归 `internal/bot`、`internal/api`、`internal/handler`、`internal/middleware` 等所有引用 telego 的包。
- **service 层**：新增 Guest 消息处理逻辑（`BotService` 增加 guest 处理入口，或新增对应 service）；抽出「年份参数 → 倒计时文本」共用函数，供 `MessageService`（`/d`）与 Guest 处理共享，避免重复。
- **bot 层**：`internal/bot/bot.go` 注册 `HandleGuestMessage` + `AnyGuestMessage()` 处理器。
- **util 层**：`GetTextByMessage` 或新增辅助函数支持剥离 `@botname` 提及前缀。
- **handler / repository / middleware 层**：无行为变更，仅可能因 telego 升级需要适配编译。
- **文档**：README 或部署文档新增 BotFather 开启 Guest Mode 的说明。
- 无数据库 schema 变更，无配置项新增。
