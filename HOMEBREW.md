# 通过Homebrew安装ServerGo

ServerGo可以通过Homebrew包管理器在macOS上安装。以下是安装步骤：

## 添加自定义Tap

使用以下命令添加ServerGo的Homebrew Tap：

```bash
brew tap CC11001100/servergo https://github.com/CC11001100/servergo
```

## 安装ServerGo

添加Tap后，使用以下命令安装ServerGo：

```bash
brew install servergo
```

或者可以在一步中完成：

```bash
brew install CC11001100/servergo/servergo
```

## 升级ServerGo

安装后，您可以使用以下命令升级到最新版本：

```bash
brew upgrade servergo
```

## 卸载ServerGo

如果需要卸载ServerGo，可以使用：

```bash
brew uninstall servergo
```

## 其他命令

- 查看ServerGo信息：`brew info servergo`
- 查看ServerGo是否有更新：`brew outdated` 