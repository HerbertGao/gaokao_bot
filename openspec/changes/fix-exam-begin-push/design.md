## 上下文

`DailySendTask`（`internal/task/daily_send_task.go`）由 cron `"0 0 * * * *"` 驱动，每个整点触发 `execute()`。`execute()` 先 `now := util.NowBJT()`，取出当前时间范围内的考试，对每场考试用 `shouldSend(exam, now)` 判定是否推送，命中后用 `GetCountDownString` 生成文案并发往所有目标。

当前 `shouldSend` 的发送闸门：
- `hours <= 24 && hours > 0` → 每小时推送倒计时
- `hours > 24 && now.Hour()==9 && now.Minute()<=1` → 每日 9:00 推送

`hours = exam.ExamBeginDate.Sub(now).Hours()`。开考整点（如 09:00）触发时 `now` 略晚于开考时刻，`hours` 为微负，两个分支都不命中 → 开考那一刻不发送。`internal/util/gaokao_bot_util.go` 中已存在 `IsExamBeginTime`（`now.After(begin) && now.Before(begin+1min)`），语义恰为「开考后第一分钟内」，但从未被生产代码调用，仅被自身测试引用。

**承重前提（必须显式记录，否则后人改动会静默破坏本修复）**：`execute()` 处理的考试来自 `ExamDateService.GetExamsInRange(now)`，其底层查询（`internal/repository/exam_date_repo.go`）以 **`exam_year_begin_date <= now <= exam_year_end_date`** 为过滤键，而非 `exam_begin_date`。对 2026 行 `exam_year_end_date = 2026-06-10 17:00:00`，故开考时刻 `06-07 09:00` 该考试仍在结果集内、`execute()` 能取到它。本修复成立**依赖此过滤键**；若将来把范围查询改为按 `exam_begin_date` 收窄，开考推送会再次静默失效且无失败测试兜底。注意现有 `internal/repository/exam_date_repo_test.go` 的 `GetExamsInRange` 用例用**相对日期**、仅断言返回计数，并未锁定开考整点这一边界；故由 tasks 3.3 在 repo 层新增固定 `2026-06-07 09:00` 的边界用例来真正守住此前提。

## 目标 / 非目标

**目标：**
- 让定时推送在开考那一刻发出一条消息（修复缺陷）。
- 开考那一刻的定时推送文案为「{exam}开始了！」。
- 复用既有 `IsExamBeginTime`，消除其死代码状态。

**非目标：**
- 不修改 `GetCountDownString`，不改变 Inline Query / 命令 / `/start` 等任何非定时推送路径的文案与行为。
- 不改 cron 表达式、发送频率、时间标准化。
- 不改 `IsExamBeginTime` 的判定语义。

## 决策

**决策 1：在 `shouldSend` 增加开考瞬间分支，置于最前。**
```go
if util.IsExamBeginTime(&exam, now) {
    return true
}
```
放在 `hours <= 24 && hours > 0` 之前，保证开考整点优先按「开考那一刻」处理。`IsExamBeginTime` 限定开考后一分钟内，一年仅命中一次，不会在 6/8、6/9 的整点重复触发。
- 备选：把 `hours > 0` 放宽为 `hours >= 0` 或 `>= -某容差`。否决——会让倒计时分支在开考瞬间渲染「正在进行中」（经 `GetCountDownString`），无法表达「开始了！」，且语义不如显式的开考判定清晰。

**决策 2：开考那一刻在 `execute()` 内使用专用文案，不走 `GetCountDownString`。**
```go
var message string
if util.IsExamBeginTime(&exam, now) {
    message = fmt.Sprintf("%s开始了！", exam.ExamDesc)
} else {
    message = util.GetCountDownString(&exam, templateContent, normalizedNow)
}
```
- 理由：用户要求「其它不变」。若改 `GetCountDownString` 的进行中分支，会连带改变 Inline Query / 命令在考试进行中的应答（仍应为「正在进行中！」）。在任务层特判可把改动隔离到定时推送。
- 文案用 `ExamDesc`（考试全称），与既有「{ExamDesc}正在进行中！」保持一致风格。

**决策 3：依赖 cron 触发使 `now` 严格晚于开考整点，不改 `IsExamBeginTime` 语义。**
- robfig/cron 计算下一激活时刻并 sleep 到点后才运行任务，`NowBJT()` 取到的时间为调度时刻**之后**（纳秒~毫秒级），不会早触发，故落入 `(begin, begin+1min)`，命中 `IsExamBeginTime`。
- 备选：把 `IsExamBeginTime` 改为含开考整秒（`!now.Before(begin)`）。否决——会破坏其既有语义测试，且 cron 不早触发使其实际不需要。
- 退化行为（诚实记录，已对照代码更正）：开考分支用严格 `now.After(begin)`、`≤24h` 倒计时分支用严格 `hours > 0`，两者都**不含** `now == begin`。故整秒边界 `now == begin` 时（`hours == 0`）两分支均不命中 → `shouldSend` 返回 false → **该 tick 不发送任何消息**（既不发「开始了！」也不发倒计时、不误发「正在进行中」、不崩溃）。这是「该 tick 漏发」、归入下方窗口边界漏发同类，而非「发了倒计时」。仅 `now < begin`（理论上的早触发，robfig/cron 不会发生）时 `hours` 微正才会落入倒计时分支——生产不可达。概率可忽略，spec 场景对此作显式说明。

**决策 4：开考分支命中时打一条 info 日志，使生产可观测、可事后审计。**
- 开考推送是**一年一次、不可补发**的事件；文档评审与单测都无法证明 cron 在真实负载下确实命中开考窗口。命中分支时记一条日志，至少含 `exam`、`now`、`begin` 与命中偏移 `now.Sub(begin)`，使来年（如 2027-06-07）可凭日志区分「未命中窗口（偏移 > 1min 或为负）」「命中但发送失败」「窗口宽度不足」，而非像今年这样无痕迹。

## 风险 / 权衡

- [cron 触发延迟 > 1 分钟（系统负载 / GC / 暂停）落在 `begin+1min` 之后] → 开考窗口仅 1 分钟、错过则**当年永久漏发**（开考是一次性不可恢复事件，与每小时倒计时「漏一次下个整点自愈」**本质不同**，不可用同类风险类比掩盖）。权衡后仍不扩窗（扩窗会增加误判与重复发送面），改以决策 4 的命中日志提供可观测兜底；窗口宽度列为 Open Question。
- [09:00 这一 tick 完全丢失（进程恰在 08:59–09:00 重启 / 调度被丢弃）] → robfig/cron 不补偿错过的 tick，开考消息漏发，下一个 10:00 tick 因 `IsExamBeginTime` 已过窗 → 走倒计时分支判定（此时考试进行中，但 `shouldSend` 各分支均不满足 → 不发）。此为 cron 调度的固有局限、Java 原版同样存在、非本次引入；记录在案，本次不引入持久化补发机制。
- [开考整秒 `now == begin`] → 开考分支（严格 `After`）与倒计时分支（严格 `hours>0`）均不含等号 → 该 tick **漏发**（不误发其它文案、不崩溃），归入「窗口边界漏发」同类，以决策 4 的命中日志可观测。不为此放宽等号，避免与 `IsExamBeginTime` 既有语义测试冲突；概率可忽略（cron 不早触发、取到恰好整秒概率极低）。
- [`execute()` 内两处都调用 `IsExamBeginTime`（`shouldSend` 与文案分支）] → 同一 `now` 下结果一致，重复计算开销可忽略，换取判定与文案的清晰分离。

## 待解决问题（Open Questions）

- 开考判定窗口宽度（当前 1 分钟）对「一次性不可恢复事件」是否足够？本次按现状不扩，以命中日志兜底；若来年日志显示仍漏，再评估扩窗或持久化补发。
