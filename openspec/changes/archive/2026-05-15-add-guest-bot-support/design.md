## 上下文

@gaokao_bot 现有三类入口：Bot 命令（`/d`、`/template`、`/debug`）、Inline Query、定时推送。命令处理依赖 Bot 是聊天成员才能收到 `Update.Message`。

Telegram 的 Guest Bots 特性让 Bot 在它**非成员**的聊天中被 @提及或被回复时收到一条 `Update.GuestMessage`，并可一次性应答。该特性自 `mymmrac/telego v1.9.0` 起被 SDK 支持，相关类型/方法：

- `Update.GuestMessage *Message` —— Guest 召唤更新。
- `Message.GuestQueryID string` —— 应答凭据。
- `Bot.AnswerGuestQuery(ctx, *AnswerGuestQueryParams)` —— 应答方法；`Params.Result` 是单个 `InlineQueryResult`。
- `telegohandler` 提供 `HandleGuestMessage` 与 `AnyGuestMessage()` 预测器。

当前 `go.mod` 锁定 `telego v1.3.3`，缺少上述全部 API。约束：Guest Mode 的启用是 BotFather 侧的开关，代码不可控。

## 目标 / 非目标

**目标：**
- 在 Bot 非成员的聊天中，被 @提及/被回复时应答一条默认模板高考倒计时。
- 支持可选的指定年份参数（如 `@gaokao_bot 2026`）。
- 最大化复用 `/d` 命令既有的「年份解析 + 多考试拼接」逻辑。

**非目标：**
- 不支持用户自定义模板、`/template`、`/debug`、Mini App 按钮。
- 不为 Guest 交互做任何持久化。
- 不改动现有命令、Inline Query、定时推送的行为。

## 决策

### 决策 1：升级 telego v1.3.3 → v1.9.0

Guest Bot 的全部类型与处理器自 v1.9.0 起提供，无法绕过。跨 6 个 minor 版本，需在升级后全量编译并跑 `make test`，逐一适配破坏性变更（重点排查 `internal/bot`、`internal/api`、`internal/handler`、`internal/middleware` 中对 telego API 的调用）。

**替代方案**：手写裸 HTTP 调用 `answerGuestQuery` —— 否决，与项目「统一用 telego SDK」的约定冲突，且失去类型安全。

### 决策 2：抽取共用的「年份 → 倒计时文本」函数

`MessageService.GetCountDownMessage` 当前内联了「解析年份参数 → 按年份或当前范围查考试 → 套默认模板 → 拼接多考试为单条文本」的逻辑。Guest 处理需要完全相同的输出。

抽出一个共用函数（签名形如 `BuildCountdownText(arg string, now time.Time) (string, error)`），输入是已提取的纯参数文本，输出是拼接好的倒计时文本。`GetCountDownMessage` 与 Guest 处理都调用它。

**替代方案**：Guest 处理直接复制粘贴逻辑 —— 否决，违反 config rules（代码改动须配套测试且避免重复），两处年份校验/模板逻辑会漂移。

### 决策 3：Guest 文本提取 —— 剥离 `@提及` 前缀

Guest 召唤消息文本通常是 `@gaokao_bot 2026` 或仅 `@gaokao_bot`。现有 `util.GetTextByMessage` 只剥离 `/命令` 前缀，不剥离 `@提及`。

新增/扩展一个提取函数：若文本以 `@` 开头，剥离首个 `@token` 及其后空白，得到剩余参数；剩余为空则按无参数处理。回复式召唤（文本与 Bot 无关或为空）同样按无参数处理 —— 即输出当前时间范围的默认倒计时。

**为何不强解析回复文本**：回复式召唤时用户消息内容不可控，强行当作年份参数会导致大量「参数无法识别」误报。按无参数走默认输出体验最稳。

### 决策 4：应答构造 —— 单条 InlineQueryResultArticle

`AnswerGuestQuery` 只接受单个 `InlineQueryResult`。复用 `InlineQueryService` 已有的 `InlineQueryResultArticle` 构造模式：`InputTextMessageContent.MessageText` 填决策 2 产出的拼接文本，`Title` 用一个固定文案（如「高考倒计时」）。无考试数据/参数非法时，按 `/d` 的兜底文案（「参数暂时无法识别。」「数据库中没有可用的信息……」）作为正文应答。

### 决策 5：处理入口落在 service 层

在 `internal/bot/bot.go` 注册 `handler.HandleGuestMessage(...)`（用 `AnyGuestMessage()` 预测器），回调转交 `BotService` 新增的 `HandleGuestMessage` 方法。与既有 `HandleMessage` / `HandleInlineQuery` 对称，保持分层一致。

## 风险 / 权衡

- **telego 大跨版本升级引入破坏性变更** → 升级后立即 `go build ./...` + `make test` 全量回归；破坏点逐个修复并在对应 `_test.go` 补验证；先在 dev 环境验证 Bot 正常长轮询再合并。
- **Guest Mode 开关在 BotFather，代码无法验证是否开启** → 在部署文档明确写明开启步骤；代码侧仅保证收到 `GuestMessage` 时能正确应答，开关缺失时不会收到更新（无副作用）。
- **`answerGuestQuery` 有时效窗口，超时应答失败** → Guest 处理走独立 `context.WithTimeout`，失败仅记录日志（与现有 `SendMessage` 失败处理一致），不重试。
- **抽取共用函数改动 `GetCountDownMessage`** → 属内部重构，行为不变；依赖既有 `message_service_test.go` 回归保证 `/d` 输出不回归。

## 待解决问题

- telego v1.3.3 → v1.9.0 之间是否存在影响本项目的破坏性变更，需在实现首步实际升级编译后才能确认清单。
- Guest 应答 `InlineQueryResultArticle` 的 `Title` 固定文案最终用词（不影响架构，实现时定）。
