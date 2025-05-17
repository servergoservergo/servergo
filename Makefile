.PHONY: build clean test release goreleaser snapshot tag

# 版本信息
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_REF ?= $(shell git symbolic-ref -q --short HEAD 2>/dev/null || git describe --tags --exact-match 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d %H:%M:%S')

# Go命令
GO = go
GOBUILD = $(GO) build
GOCLEAN = $(GO) clean
GOTEST = $(GO) test
GOGET = $(GO) get

# 项目信息
BINARY_NAME = servergo
MAIN_PACKAGE = .

# Go构建参数
LDFLAGS = -ldflags "\
	-X 'github.com/CC11001100/servergo/pkg/version.Version=$(VERSION)' \
	-X 'github.com/CC11001100/servergo/pkg/version.BuildTime=$(BUILD_TIME)' \
	-X 'github.com/CC11001100/servergo/pkg/version.GitCommit=$(GIT_COMMIT)' \
	-X 'github.com/CC11001100/servergo/pkg/version.GitRef=$(GIT_REF)' \
	-w -s"

# 主要目标
build:
	@echo "开始构建..."
	@$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)

# 生成发布版本（带有版本标签）
release:
	@echo "构建发布版本 $(BINARY_NAME) $(VERSION)"
	@$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "创建发布压缩文件..."
	@tar -czvf $(BINARY_NAME)-$(VERSION)-$(shell go env GOOS)-$(shell go env GOARCH).tar.gz $(BINARY_NAME)

# 清理
clean:
	@echo "清理项目..."
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME)-*.tar.gz

# 运行测试
test:
	@echo "运行测试..."
	@$(GOTEST) -v ./...

# 安装到系统
install: build
	@echo "安装到系统..."
	@$(BINARY_NAME) install

# 使用 GoReleaser 进行发布
goreleaser:
	@echo "使用 GoReleaser 进行发布..."
	@goreleaser release --clean

# 使用 GoReleaser 创建快照（不发布）
snapshot:
	@echo "使用 GoReleaser 创建快照（不发布）..."
	@goreleaser release --snapshot --clean

# 创建并推送新版本标签
tag:
	@if [ -z "$(VERSION)" ]; then echo "ERROR: 请指定版本号，例如: make tag VERSION=1.0.0"; exit 1; fi
	@echo "创建版本标签 v$(VERSION)..."
	@git tag -a v$(VERSION) -m "Release v$(VERSION)"
	@echo "推送标签到远程仓库..."
	@git push origin v$(VERSION)
	@echo "标签已创建并推送。现在可以运行 'make goreleaser' 来发布此版本。"

# 帮助信息
help:
	@echo "可用命令："
	@echo "  make build         - 构建基本版本"
	@echo "  make release       - 构建发布版本"
	@echo "  make test          - 运行测试"
	@echo "  make clean         - 清理构建产物"
	@echo "  make install       - 安装到系统"
	@echo "  make goreleaser    - 使用 GoReleaser 发布当前标签"
	@echo "  make snapshot      - 使用 GoReleaser 创建快照（不发布）"
	@echo "  make tag           - 创建并推送新版本标签"
	@echo ""
	@echo "可设置的环境变量："
	@echo "  VERSION            - 指定版本号，例如: make release VERSION=1.0.0"
	@echo ""
	@echo "发布流程："
	@echo "  1. make tag VERSION=x.y.z   - 创建并推送新标签"
	@echo "  2. make goreleaser          - 发布到 GitHub"

# 默认命令
default: build 