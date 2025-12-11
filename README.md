# 高考倒计时Bot (Go版本)

[![Actions][ico-actions]][link-actions]
[![Releases][ico-releases]][link-releases]
[![Stars][ico-stars]][link-stars]
[![Software License][ico-license]](LICENSE)

[![Chat on Telegram][ico-telegram]][link-telegram]
[![Chat on Telegram][ico-telegram-channel]][link-telegram-channel]

## Introduction
高考倒计时 @gaokao_bot 是一个 Telegram Bot，可以实现通过发送 Message、Inline Query 方式获取当前时间的高考倒计时。

本项目使用 Go 语言实现，基于 Java 原项目的完整功能复刻。

## Features
- 倒计时查询 - 发送命令或 Inline Query 获取高考倒计时
- 定时推送 - 自动推送倒计时到指定群组
- Mini App - [可视化管理倒计时模板](https://github.com/HerbertGao/gaokao_bot_mini_app)
- 多环境支持 - 开发、测试、生产环境配置分离

## Quick Start

### Requirements
- Go 1.21+
- MySQL 8.0+

### Installation

```bash
# Clone repository
git clone https://github.com/HerbertGao/gaokao_bot.git
cd gaokao_bot

# Install dependencies
go mod download

# Configure environment
cp .env.example .env
# Edit .env with your configuration

# Initialize database
mysql -u root -p < sql/init.sql

# Run
go run cmd/gaokao_bot/main.go -env=dev
```

### Build

```bash
./scripts/build.sh
./bin/gaokao_bot -env=prod
```

## Tech Stack
- Go 1.21+
- [telego](https://github.com/mymmrac/telego) - Telegram Bot SDK
- GORM - ORM
- MySQL 8.0+

## Thanks to
- [mymmrac / telego](https://github.com/mymmrac/telego)

## License
The MIT License (MIT). Please see [License File](LICENSE) for more information.

[ico-actions]: https://github.com/HerbertGao/gaokao_bot/workflows/Go%20CI/badge.svg
[ico-telegram]: https://img.shields.io/badge/@gaokao__bot-2CA5E0.svg?style=flat-square&logo=telegram&label=Telegram
[ico-telegram-channel]: https://img.shields.io/badge/@GaokaoCountdown-2CA5E0.svg?style=flat-square&logo=telegram&label=Telegram
[ico-releases]: https://img.shields.io/github/release/HerbertGao/gaokao_bot
[ico-stars]: https://img.shields.io/github/stars/HerbertGao/gaokao_bot
[ico-license]: https://img.shields.io/github/license/HerbertGao/gaokao_bot

[link-actions]: https://github.com/HerbertGao/gaokao_bot/actions
[link-telegram]: https://t.me/gaokao_bot
[link-telegram-channel]: https://t.me/GaokaoCountdown
[link-releases]: https://github.com/HerbertGao/gaokao_bot/releases
[link-stars]: https://github.com/HerbertGao/gaokao_bot
[link-license]: https://opensource.org/licenses/MIT
