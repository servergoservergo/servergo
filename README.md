# ServerGo

ServerGo 是一个简单的命令行工具，用于快速启动HTTP文件服务器，类似于Python的`http.server`模块，但使用Go实现，提供更好的性能。

## 项目结构

```
servergo/
├── main.go        # 主程序入口点
├── cmd/           # 命令行命令
│   ├── root.go    # 根命令
│   ├── start.go   # 启动服务器命令(默认命令)
│   ├── config.go  # 配置管理命令
│   └── version.go # 版本命令(包含关于信息)
├── pkg/           # 包含可重用的代码包
│   ├── server/    # 文件服务器的核心实现
│   ├── config/    # 配置管理功能
│   └── utils/     # 工具函数(端口探测等)
└── README.md      # 项目文档
```

## 安装

```bash
go install github.com/CC11001100/servergo@latest
```

## 使用方法

### 基本用法

在当前目录启动一个文件服务器，默认端口自动探测（或使用配置文件中设置的默认端口）：

```bash
# 以下命令效果相同
servergo
servergo start
servergo run      # 别名
servergo serve    # 别名
servergo launch   # 别名
```

### 指定端口

```bash
servergo -p 9000
# 或者
servergo --port 9000
# 也可以明确使用子命令
servergo start -p 9000

# 如果指定的端口已被占用，将自动探测一个可用的端口
```

### 指定目录

```bash
servergo -d /path/to/your/files
# 或者
servergo --dir /path/to/your/files
# 也可以明确使用子命令
servergo start -d /path/to/your/files

# 目录路径可以是绝对路径或相对路径
```

### 控制自动打开浏览器

默认情况下，服务器启动后会自动打开浏览器访问页面，可以使用以下参数禁用此行为：

```bash
servergo -o=false
# 或者
servergo --open=false
# 也可以明确使用子命令
servergo start -o=false
```

### 同时指定多个参数

```bash
servergo -p 9000 -d /path/to/your/files -o=false
# 或者明确使用子命令
servergo start -p 9000 -d /path/to/your/files -o=false
```

### 管理配置

ServerGo支持持久化配置，配置文件保存在用户主目录的`.servergo`目录下。

#### 列出所有配置

```bash
servergo config list
```

#### 获取指定配置

```bash
servergo config get default_port
servergo config get default_dir
servergo config get auto_open
```

#### 设置指定配置

```bash
# 设置默认端口，0表示自动探测
servergo config set default_port 8080

# 设置默认目录
servergo config set default_dir /path/to/your/files

# 设置是否自动打开浏览器
servergo config set auto_open true
# 支持多种布尔值表示
servergo config set auto_open yes
servergo config set auto_open 1
# 关闭自动打开
servergo config set auto_open false
servergo config set auto_open no
servergo config set auto_open 0
```

### 查看版本和项目信息

```bash
servergo version
```

### 显示帮助信息

```bash
servergo --help
servergo start --help
servergo config --help
```

## 参数说明

- `-p, --port`: 指定服务器监听的端口（默认使用配置中设置的端口或自动探测）
- `-d, --dir`: 指定要提供服务的目录路径（默认使用配置中设置的目录或当前目录）
- `-o, --open`: 指定是否在启动服务器后自动打开浏览器（默认使用配置中的设置）

## 特性

- **自动探测端口**: 如果不指定端口或指定端口已被占用，自动找到一个可用端口启动服务器
- **路径灵活性**: 支持绝对路径和相对路径，便于使用
- **多种启动方式**: 支持多种命令别名，适应不同习惯
- **持久化配置**: 支持保存默认配置，无需每次都指定相同的参数
- **自动打开浏览器**: 服务器启动后可以自动打开浏览器，也可以通过配置或参数禁用

## 开发

如果您想从源码构建，请执行以下命令：

```bash
# 克隆仓库
git clone https://github.com/CC11001100/servergo.git
cd servergo

# 安装依赖
go mod download

# 构建
go build -o servergo

# 运行
./servergo
```

## 许可证

参见 [LICENSE](LICENSE) 文件。 