.PHONY: build clean test release goreleaser snapshot tag

# 默认版本号
VERSION ?= 0.1.0

# 主要目标
build:
	@echo "开始构建..."
	@chmod +x ./build.sh
	@VERSION=$(VERSION) ./build.sh

# 生成发布版本（带有版本标签）
release:
	@echo "构建发布版本 v$(VERSION)..."
	@chmod +x ./build.sh
	@VERSION=$(VERSION) ./build.sh

# 清理
clean:
	@echo "清理项目..."
	@rm -f servergo
	@rm -rf dist/

# 运行测试
test:
	@echo "运行测试..."
	@go test ./...

# 安装到系统
install: build
	@echo "安装到系统..."
	@cp servergo /usr/local/bin/

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