# ServerGo 发布指南

本指南说明如何使用 GoReleaser 发布 ServerGo 的新版本。

## 准备工作

1. 安装 GoReleaser（如果尚未安装）

```bash
# 使用 Homebrew 安装
brew install goreleaser

# 或者使用 Go 安装
go install github.com/goreleaser/goreleaser@latest
```

2. 设置 GitHub 令牌

在 GitHub 上创建个人访问令牌（Personal Access Token），并将其设置为环境变量：

```bash
export GITHUB_TOKEN=你的GitHub令牌
```

## 发布流程

### 1. 创建版本标签

可以使用 Makefile 中的 `tag` 命令来创建并推送版本标签：

```bash
make tag VERSION=1.0.0
```

这将创建一个名为 `v1.0.0` 的标签并推送到远程仓库。

### 2. 手动触发发布

标签推送后，GitHub Actions 会自动触发发布流程。但如果需要在本地手动触发，可以使用：

```bash
make goreleaser
```

### 3. 创建快照（不发布）

如果您想在不实际发布的情况下测试发布过程，可以创建一个快照：

```bash
make snapshot
```

## 发布产物

GoReleaser 会自动为以下平台构建二进制文件：

- Windows (amd64)
- macOS (amd64, arm64)
- Linux (amd64, arm64)

构建产物将包括：

- 二进制可执行文件
- 支持的发布档案（Linux/macOS 为 `.tar.gz`，Windows 为 `.zip`）
- 校验和文件（`checksums.txt`）
- Homebrew 配方

## Homebrew 安装

发布后，用户可以通过 Homebrew 安装 ServerGo：

```bash
brew install CC11001100/tap/servergo
```

## GitHub Actions 自动化

每当推送新标签（格式为 `v*`）时，GitHub Actions 工作流将自动运行，执行以下步骤：

1. 检出代码
2. 设置 Go 环境
3. 运行测试
4. 使用 GoReleaser 构建和发布
5. 上传构建产物

## 故障排除

如果发布过程中遇到问题：

1. 检查 GitHub Actions 工作流输出
2. 确保 `GITHUB_TOKEN` 环境变量已正确设置
3. 验证 `.goreleaser.yml` 配置是否正确
4. 尝试使用 `goreleaser --debug` 获取详细日志 