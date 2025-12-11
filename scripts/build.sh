#!/bin/bash

set -e

ENV=${1:-dev}

echo "Building Gaokao Bot for environment: $ENV"

# 创建输出目录
mkdir -p bin

# 获取版本信息
VERSION=$(grep '^var Version = ' internal/version/version.go | sed 's/var Version = "\(.*\)"/\1/')
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo "Version: $VERSION"
echo "Build Time: $BUILD_TIME"
echo "Git Commit: $GIT_COMMIT"

# 编译（根据系统平台编译）
CGO_ENABLED=0 go build -a -installsuffix cgo \
  -ldflags "-s -w -X 'github.com/herbertgao/gaokao_bot/internal/version.Version=${VERSION}' -X 'github.com/herbertgao/gaokao_bot/internal/version.BuildTime=${BUILD_TIME}' -X 'github.com/herbertgao/gaokao_bot/internal/version.GitCommit=${GIT_COMMIT}'" \
  -o bin/gaokao_bot \
  ./cmd/gaokao_bot

echo "Build completed: bin/gaokao_bot"
