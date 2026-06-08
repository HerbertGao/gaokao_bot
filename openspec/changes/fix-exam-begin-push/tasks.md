## 1. 修复 shouldSend 开考瞬间判定

- [x] 1.1 在 `internal/task/daily_send_task.go` 的 `shouldSend` 中，于 `hours <= 24 && hours > 0` 分支之前增加：`if util.IsExamBeginTime(&exam, now) { return true }`（注意 `exam` 取地址；`shouldSend` 的 `exam` 形参为值类型，取址合法）。开考分支必须**置于所有倒计时分支之前**，确保开考那一刻优先短路返回。
- [x] 1.2 确认 `util` 包已导入（现有文件已导入），无新增依赖。

## 2. 开考瞬间专用文案（抽取可测函数，消除单测接缝缺失）

- [x] 2.1 在 `internal/task/daily_send_task.go` 抽取一个**内部纯函数**（无副作用、不打日志）`buildMessage(exam *model.ExamDate, now, normalizedNow time.Time, templateContent string) string`：当 `util.IsExamBeginTime(exam, now)` 为真时返回 `fmt.Sprintf("%s开始了！", exam.ExamDesc)`；否则返回 `util.GetCountDownString(exam, templateContent, normalizedNow)`。开考分支必须先判、短路，**禁止**让开考那一刻落到 `GetCountDownString` 的「正在进行中！」分支。命中日志不放在本函数内（见 2.4 由 `execute()` 打），以保持其纯函数性便于单测。
- [x] 2.2 修改 `execute()`：**删除**原 `daily_send_task.go:106` 的内联 `message := util.GetCountDownString(...)`，改为 `message := t.buildMessage(&exam, now, normalizedNow, templateContent)`。`execute()` 内**有且仅有这一处** message 赋值来源，禁止保留旧内联调用或另写孤立分支——否则 3.2 单测虽绿、生产开考仍会发「正在进行中！」（vacuous green）。**实现完成自检 + apply 期 review 必须确认 line 106 的内联调用确已被替换、execute 唯一经 buildMessage 取文案。**
- [x] 2.3 在文件顶部 import 块补充 `"fmt"`（当前未导入）。
- [x] 2.4 在 `execute()` 内（非 `buildMessage` 内）于 `util.IsExamBeginTime(&exam, now)` 为真时记一条 info 日志，至少含 `exam`、`now`、`begin`、命中偏移 `now.Sub(begin)`（如 `t.logger.Infof("开考推送已触发: exam=%s now=%v begin=%v offset=%v", exam.ExamDesc, now, exam.ExamBeginDate, now.Sub(exam.ExamBeginDate))`），供生产事后审计是否真命中开考窗口。
- [x] 2.5 确认不改动 `internal/util/gaokao_bot_util.go` 的 `GetCountDownString` / `IsExamTime` / `IsExamBeginTime`。

## 3. 测试

- [x] 3.1 在 `internal/task/daily_send_task_test.go` 的 `shouldSend` 表驱动用例中新增：「开考后 10 秒」「开考后 59 秒」返回 `true`；「开考整点 0 秒」（不命中开考分支、`hours>0` 亦不满足）返回 `false`；「开考后 1 分钟」（已过窗）返回 `false`。
- [x] 3.2 对 `buildMessage` 做表驱动单测：开考后 10s → 「{ExamDesc}开始了！」；开考前 1 小时 → 倒计时文案（含 `{exam}`/`{time}` 渲染结果）。断言必须针对 **2.2 中 `execute()` 实际调用的同一 `buildMessage`**，不得断言一个 `execute` 未调用的孤立函数（避免 vacuous assert）。注：`now==begin` 整秒边界**不在本单测覆盖**——该 tick 由 `shouldSend` 返回 false 直接 skip、`buildMessage` 不被调用（即生产不可达），其漏发行为由 3.1 的「开考整点 0 秒 → shouldSend 返回 false」用例覆盖，不要在 `buildMessage` 层对该输入断言以免误导。
- [x] 3.3 承重前提测试放在 **`internal/repository/exam_date_repo_test.go`**（repo 层已有 sqlite in-memory 脚手架；`task` 包测试 `examDateService` 为 nil、无 DB，**不能**在此验证）。新增用例：seed 一行 `exam_year_begin_date=2025-06-10 17:00 / exam_year_end_date=2026-06-10 17:00 / exam_begin_date=2026-06-07 09:00`，断言 `GetExamsInRange(2026-06-07 09:00:00 BJT)` 返回该行——固定开考整点这一边界（现有相对日期用例只断言计数、未锁定该边界），防止将来收窄范围查询静默破坏修复。
- [x] 3.4 `make test` 全绿；`golangci-lint`（errcheck / govet / staticcheck / gosec / revive / gocyclo）无新增告警。若抽取 `buildMessage` 后 `execute()` 复杂度仍触 gocyclo，进一步下沉逻辑到 `buildMessage`。

## 4. 收尾

- [x] 4.1 自测：确认开考整点之外的发送规则（每日 9:00 / 每小时倒计时 / 进行中 / 已结束）行为不变；确认 6/8、6/9 等考试期内整点不会误发「开始了！」或「正在进行中！」。
- [ ] 4.2 按 Conventional Commits 提交：`fix: 修复开考那一刻定时推送缺失并改为「{exam}开始了！」`。
