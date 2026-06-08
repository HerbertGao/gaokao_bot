## 为什么

定时推送从未在「高考开考那一刻」发出消息。`DailySendTask.shouldSend` 的发送闸门第一分支要求 `hours > 0`（严格大于），cron 在开考整点（如 `06-07 09:00`）触发时 `hours` 已为 0 或微负，闸门返回 false，于是开考瞬间什么都不发。重写时专门写了 `IsExamBeginTime`（开考后第一分钟内为真）来识别这一刻，却没有接入任何生产逻辑，成了只被自身测试引用的死代码。今年（2026）高考因此缺失了开考那一刻的推送。

此外，开考那一刻的推送文案当前会落到 `GetCountDownString` 的「考试进行中」分支，输出「…正在进行中！」，希望改为更明确的「…开始了！」。

## 变更内容

- 修复 `DailySendTask` 在开考那一刻不推送的缺陷：在 `shouldSend` 中接入 `IsExamBeginTime`，使开考后第一分钟内的 cron tick 命中发送（一年仅触发一次）。
- 开考那一刻的定时推送文案改为「{exam}开始了！」，与现有「正在进行中」一致使用考试全称（`ExamDesc`）。该专用文案仅作用于定时推送的开考瞬间，不修改 `GetCountDownString`，因此 Inline Query、命令应答等其它路径的「正在进行中」文案保持不变。
- 其余发送规则（>24 小时每日 9:00、≤24 小时每小时、进行中/已结束判定、时间标准化）全部不变。

## 功能 (Capabilities)

### 新增功能
- `scheduled-countdown-push`: 定时倒计时推送任务的发送判定与文案规则（含开考那一刻推送）。此前未被规范覆盖，本次将其行为固化为规范，并补充开考那一刻的需求。

### 修改功能
<!-- 无现有规范覆盖定时推送行为，故以新增功能形式固化。 -->

## 影响

- 代码：`internal/task/daily_send_task.go`（`shouldSend` 增加开考瞬间分支；`execute` 在开考瞬间改用专用文案）。
- 复用：接入既有的 `internal/util/gaokao_bot_util.go` 中的 `IsExamBeginTime`（消除死代码），不改动其语义。
- 测试：在 `internal/task/daily_send_task_test.go` 补开考瞬间用例与 `buildMessage` 单测；在 `internal/repository/exam_date_repo_test.go`（repo 层有 DB 脚手架）补「开考整点该考试仍在 `GetExamsInRange` 结果集内」的边界用例——现有 repo 用例用相对日期、仅断言计数，未锁定该边界。
- 承重前提：本修复依赖 `GetExamsInRange` 以 `exam_year_begin_date/exam_year_end_date` 为过滤键（而非 `exam_begin_date`），使开考时刻该考试仍被取到。若将来收窄该查询会静默破坏本修复，故以测试固定此前提（详见 design.md）。
- 不影响 service / repository / handler 分层、数据库、配置与 cron 表达式。

## 非目标 (Non-Goals)

- 不修改 `GetCountDownString` 的「正在进行中！」/「已经结束了。」文案，也不改变 Inline Query、`/start`、命令应答等任何非定时推送路径的行为。
- 不调整 cron 表达式、发送频率（每日 9:00 / 每小时）或时间标准化逻辑。
- 不修改 `IsExamBeginTime` 的判定语义（仍为开考后第一分钟内、不含开考整秒；依赖 cron 触发时 `now` 必然略晚于开考整点）。
- 不涉及考试日期数据、模板、发送目标等持久化数据。
