# ServerGo - 让文件分享变得简单

ServerGo 是一个轻量级的HTTP文件服务器，解决了日常开发中快速分享文件的痛点。相比Python的`http.server`，它更快、更智能，而且不会遇到端口冲突的烦恼。

## 为什么选择ServerGo？

在日常开发中，我们经常需要：
- 快速分享本地文件给同事
- 临时搭建一个文件下载服务器
- 在局域网内传输文件
- 测试前端项目的本地部署

ServerGo就是为这些场景而生的。它自动处理端口冲突，支持持久化配置，让文件分享变得前所未有的简单。

## 一键安装

### 使用Go安装（推荐）
```bash
go install github.com/CC11001100/servergo@latest
```

### 使用Homebrew（macOS）
```bash
brew tap CC11001100/servergo https://github.com/CC11001100/servergo
brew install servergo
```

## 快速开始

### 最常用的场景
```bash
# 在当前目录启动文件服务器（自动选择可用端口）
servergo

# 指定端口启动
servergo -p 8080

# 分享指定目录的文件
servergo -d ~/Documents

# 启动但不自动打开浏览器
servergo -o=false
```

启动后，你会看到类似这样的输出：
```
🚀 ServerGo 文件服务器已启动
📁 服务目录: /Users/yourname/Documents
🌐 访问地址: http://localhost:8080
🔍 自动打开浏览器: 已启用
```

## 实际使用案例

### 案例1：前端开发调试
```bash
# 在Vue/React项目目录启动服务器
servergo -p 3000
# 访问 http://localhost:3000 查看构建结果
```

### 案例2：团队文件分享
```bash
# 在项目根目录分享文档
servergo -d ./docs -p 8080
# 告诉同事访问你的IP地址即可下载文件
```

### 案例3：临时文件传输
```bash
# 快速分享桌面上的文件
servergo -d ~/Desktop -p 9999 -o=false
# 通过局域网IP让手机/其他设备访问
```

## 智能配置管理

ServerGo会记住你的使用习惯，配置文件保存在`~/.servergo/config.json`。

### 查看当前配置
```bash
servergo config list
```

### 设置默认端口
```bash
servergo config set default_port 8080
```

### 设置默认目录
```bash
servergo config set default_dir ~/Projects
```

### 关闭自动打开浏览器
```bash
servergo config set auto_open false
```

## 高级功能

### 端口自动探测
如果不指定端口，ServerGo会自动找到第一个可用端口（从配置端口开始探测）：
```bash
servergo
# 输出：端口8080被占用，自动使用8081
```

### 目录浏览美化
ServerGo提供现代化的目录浏览界面，支持：
- 文件图标显示
- 文件大小和修改时间
- 多主题切换（深色/浅色等）
- 响应式设计，手机友好

### 安全特性
- 路径遍历攻击防护
- 符号链接安全检测
- 访问日志记录

## 命令别名

ServerGo支持多种启动命令，适应不同习惯：
```bash
servergo
servergo start
servergo run
servergo serve
servergo launch
```

## 故障排除

### 端口被占用
```bash
servergo -p 8080
# 自动提示：端口8080被占用，正在寻找可用端口...
# 找到可用端口：8081
```

### 权限问题
```bash
# 确保对服务目录有读取权限
chmod +r ~/Documents
servergo -d ~/Documents
```

### 局域网访问
```bash
# 获取本机IP
ifconfig | grep inet
# 使用IP:端口让其他设备访问
```

## 从源码构建

```bash
git clone https://github.com/CC11001100/servergo.git
cd servergo
go build -o servergo
./servergo
```

## 项目状态

ServerGo正在积极开发中，欢迎提交Issue和Pull Request。当前版本特性：

✅ 基础文件服务  
✅ 端口自动探测  
✅ 配置持久化  
✅ 目录浏览美化  
✅ 安全路径验证  
✅ 多平台支持  

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

---

**ServerGo让文件分享变得简单，希望它能成为你开发工具箱中的得力助手！** 🚀