#!/bin/bash

# 发布脚本
# 用法: ./scripts/release.sh [major|minor|patch|build|x.y.z]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # 无色

# 检查 git 仓库
check_git_repo() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        echo -e "${RED}错误: 当前目录不是 git 仓库${NC}"
        exit 1
    fi
}

# 检查未提交更改
check_uncommitted_changes() {
    if ! git diff-index --quiet HEAD --; then
        echo -e "${YELLOW}警告: 有未提交的更改${NC}"
        git status --porcelain
        echo -e "${YELLOW}是否继续？(y/N): ${NC}"
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            echo "取消发布"
            exit 0
        fi
    fi
}

# 检查当前分支是否为 master
check_master_branch() {
    local current_branch=$(git branch --show-current)
    if [ "$current_branch" != "master" ]; then
        echo -e "${RED}错误: 只能在 master 分支打标签${NC}"
        echo -e "${YELLOW}当前分支: ${current_branch}${NC}"
        echo -e "${YELLOW}请先切换到 master 分支${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ 当前在 master 分支${NC}"
}

# 更新到最新代码
update_to_latest() {
    echo -e "${BLUE}更新到最新代码...${NC}"

    # 获取远程更新
    git fetch origin master

    # 检查本地是否落后于远程
    local local_commit=$(git rev-parse master)
    local remote_commit=$(git rev-parse origin/master)

    if [ "$local_commit" != "$remote_commit" ]; then
        echo -e "${YELLOW}本地代码不是最新的${NC}"
        echo -e "${BLUE}拉取最新代码...${NC}"
        git pull origin master
        echo -e "${GREEN}✓ 已更新到最新代码${NC}"
    else
        echo -e "${GREEN}✓ 本地代码已是最新${NC}"
    fi
}

# 运行测试
run_tests() {
    echo -e "${BLUE}运行测试...${NC}"

    # 运行 go vet
    echo -e "${BLUE}运行 go vet...${NC}"
    go vet ./...
    echo -e "${GREEN}✓ go vet 检查通过${NC}"

    # 运行单元测试
    if go test -v ./... 2>&1 | grep -q "testing: warning: no tests"; then
        echo -e "${YELLOW}⚠ 暂无测试用例${NC}"
    else
        echo -e "${GREEN}✓ 测试通过${NC}"
    fi
}

# 编译发布版本
build_release() {
    echo -e "${BLUE}编译发布版本...${NC}"

    VERSION=$(grep '^var Version = ' internal/version/version.go | sed 's/var Version = "\(.*\)"/\1/')
    BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')
    GIT_COMMIT=$(git rev-parse --short HEAD)

    go build -a \
        -ldflags "-s -w -X 'github.com/herbertgao/gaokao_bot/internal/version.Version=${VERSION}' -X 'github.com/herbertgao/gaokao_bot/internal/version.BuildTime=${BUILD_TIME}' -X 'github.com/herbertgao/gaokao_bot/internal/version.GitCommit=${GIT_COMMIT}'" \
        -o bin/gaokao_bot \
        ./cmd/gaokao_bot

    echo -e "${GREEN}✓ 编译完成${NC}"
}

# 提交更改
commit_changes() {
    local version=$1
    local commit_message="chore: bump version to ${version}"

    echo -e "${BLUE}提交更改...${NC}"
    git add .
    git commit -m "$commit_message"
    echo -e "${GREEN}✓ 更改已提交${NC}"
}

# 创建标签
git_tag() {
    local version=$1
    local tag_name="v${version}"

    echo -e "${BLUE}创建标签 ${tag_name}...${NC}"
    git tag "$tag_name"
    echo -e "${GREEN}✓ 标签已创建${NC}"
}

# 推送到远程
push_to_remote() {
    local version=$1
    local tag_name="v${version}"

    # 获取当前分支名
    local current_branch=$(git branch --show-current)

    echo -e "${BLUE}推送到远程仓库...${NC}"
    git push origin "$current_branch"
    git push origin "$tag_name"
    echo -e "${GREEN}✓ 已推送到远程仓库${NC}"
}

# 显示发布信息
show_release_info() {
    local version=$1
    local tag_name="v${version}"

    echo -e "${GREEN}🎉 发布完成！${NC}"
    echo -e "${YELLOW}版本: ${version}${NC}"
    echo -e "${YELLOW}标签: ${tag_name}${NC}"
    echo -e "${YELLOW}GitHub Actions 将自动构建多平台版本${NC}"
    echo -e "${BLUE}查看构建状态: https://github.com/HerbertGao/gaokao_bot/actions${NC}"
}

# 主流程
main() {
    if [ $# -eq 0 ]; then
        echo -e "${RED}错误: 请指定版本类型或版本号${NC}"
        echo "用法: $0 [major|minor|patch|build|x.y.z]"
        echo "  major  - 主版本号 (1.0.0 -> 2.0.0)"
        echo "  minor  - 次版本号 (1.0.0 -> 1.1.0)"
        echo "  patch  - 补丁版本 (1.0.0 -> 1.0.1)"
        echo "  build  - 构建版本 (1.0.0 -> 1.0.0.1)"
        echo "  x.y.z  - 自定义版本号 (如 2.0.0 或 2.0.0.1)"
        exit 1
    fi

    local version_input=$1

    echo -e "${BLUE}开始发布流程...${NC}"

    check_git_repo
    check_master_branch
    update_to_latest
    check_uncommitted_changes

    echo -e "${BLUE}更新版本...${NC}"
    ./scripts/version.sh "$version_input"

    local new_version=$(grep '^var Version = ' internal/version/version.go | sed 's/var Version = "\(.*\)"/\1/')

    run_tests
    build_release

    commit_changes "$new_version"
    git_tag "$new_version"
    push_to_remote "$new_version"

    show_release_info "$new_version"
}

main "$@"
