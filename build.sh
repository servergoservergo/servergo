#!/bin/bash

# 获取版本信息
VERSION=${VERSION:-"0.1.0"}
BUILD_DATE=$(date "+%Y-%m-%d %H:%M:%S")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "未知")

# 显示将要构建的版本信息
echo "构建版本："
echo "- 版本号: ${VERSION}"
echo "- 构建日期: ${BUILD_DATE}"
echo "- Git提交: ${GIT_COMMIT}"
echo ""

# 构建二进制文件
go build -o servergo \
  -ldflags "-X 'github.com/CC11001100/servergo/cmd.Version=${VERSION}' \
            -X 'github.com/CC11001100/servergo/cmd.BuildDate=${BUILD_DATE}' \
            -X 'github.com/CC11001100/servergo/cmd.GitCommit=${GIT_COMMIT}'"

echo "构建完成: servergo" 