package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志级别常量
const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// 日志配置常量
const (
	// MaxLogSize 单个日志文件的最大大小（MB）
	MaxLogSize = 128
	// MaxLogBackups 保留的旧日志文件的最大数量
	MaxLogBackups = 3
	// MaxLogAge 保留的旧日志文件的最大天数
	MaxLogAge = 28
)

// 日志级别名称
var levelNames = map[int]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
}

// 日志级别对应的颜色
var levelColors = map[int]func(format string, a ...interface{}) string{
	DEBUG:   color.WhiteString,
	INFO:    color.CyanString,
	WARNING: color.YellowString,
	ERROR:   color.RedString,
	FATAL:   color.New(color.FgHiRed, color.Bold).SprintfFunc(),
}

// 访问日志的颜色
var (
	MethodColor = color.New(color.FgGreen, color.Bold).SprintfFunc()
	PathColor   = color.CyanString
	StatusColor = map[int]func(format string, a ...interface{}) string{
		// 2xx
		200: color.New(color.FgGreen).SprintfFunc(),
		201: color.New(color.FgGreen).SprintfFunc(),
		// 3xx
		301: color.New(color.FgYellow).SprintfFunc(),
		302: color.New(color.FgYellow).SprintfFunc(),
		// 4xx
		400: color.New(color.FgMagenta).SprintfFunc(),
		401: color.New(color.FgMagenta).SprintfFunc(),
		403: color.New(color.FgMagenta).SprintfFunc(),
		404: color.New(color.FgMagenta).SprintfFunc(),
		// 5xx
		500: color.New(color.FgRed, color.Bold).SprintfFunc(),
		501: color.New(color.FgRed, color.Bold).SprintfFunc(),
		502: color.New(color.FgRed, color.Bold).SprintfFunc(),
		503: color.New(color.FgRed, color.Bold).SprintfFunc(),
	}
	TimeColor     = color.New(color.FgWhite).SprintfFunc()
	ClientIPColor = color.New(color.FgWhite).SprintfFunc()
	BytesColor    = color.New(color.FgMagenta).SprintfFunc()
)

// Logger 结构表示一个日志记录器
type Logger struct {
	// 控制台输出
	console io.Writer
	// 文件输出
	file io.Writer
	// 日志级别
	level int
	// 是否启用文件日志
	fileEnabled bool
}

// LogConfig 日志配置
type LogConfig struct {
	// 日志级别
	Level int
	// 是否启用文件日志
	EnableFileLog bool
	// 日志文件名（不包含路径）
	Filename string
}

// 获取日志目录
func getLogDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("无法获取用户主目录: %v", err)
	}

	// 创建 .servergo/logs 目录
	logDir := filepath.Join(home, ".servergo", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return "", fmt.Errorf("无法创建日志目录: %v", err)
	}

	return logDir, nil
}

// 创建一个新的Logger实例
func New(config LogConfig) (*Logger, error) {
	logger := &Logger{
		console:     os.Stdout,
		level:       config.Level,
		fileEnabled: config.EnableFileLog,
	}

	// 如果启用了文件日志
	if config.EnableFileLog {
		logDir, err := getLogDir()
		if err != nil {
			return nil, err
		}

		// 使用lumberjack进行日志轮转
		logger.file = &lumberjack.Logger{
			Filename:   filepath.Join(logDir, config.Filename),
			MaxSize:    MaxLogSize,    // 每个日志文件最大128MB
			MaxBackups: MaxLogBackups, // 保留3个旧文件
			MaxAge:     MaxLogAge,     // 保留28天
			Compress:   true,          // 压缩旧的日志文件
		}
	}

	return logger, nil
}

// 默认日志实例
var Default *Logger

// 初始化默认日志实例
func init() {
	var err error
	Default, err = New(LogConfig{
		Level:         INFO,
		EnableFileLog: true,
		Filename:      "servergo.log",
	})
	if err != nil {
		// 如果无法创建文件日志，则回退到只使用控制台
		Default = &Logger{
			console: os.Stdout,
			level:   INFO,
		}
		fmt.Fprintf(os.Stderr, "警告: 无法初始化文件日志: %v\n", err)
	}
}

// SetOutput 设置控制台日志输出目标
func (l *Logger) SetOutput(w io.Writer) {
	l.console = w
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level int) {
	l.level = level
}

// formatMessage 格式化日志消息
func (l *Logger) formatMessage(level int, format string, args ...interface{}) (string, string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	levelStr := levelNames[level]

	// 格式化消息内容
	var message string
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format
	}

	// 控制台带颜色的消息
	colorFunc := levelColors[level]
	coloredLevel := colorFunc("[%s]", levelStr)
	timeStr := TimeColor("%s", timestamp)
	coloredMessage := fmt.Sprintf("%s %s %s", timeStr, coloredLevel, message)

	// 文件日志不带颜色的消息
	plainMessage := fmt.Sprintf("%s [%s] %s", timestamp, levelStr, message)

	return coloredMessage, plainMessage
}

// log 内部日志方法
func (l *Logger) log(level int, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	coloredMsg, plainMsg := l.formatMessage(level, format, args...)

	// 输出到控制台（带颜色）
	fmt.Fprintln(l.console, coloredMsg)

	// 如果启用了文件日志，输出到文件（不带颜色）
	if l.fileEnabled && l.file != nil {
		fmt.Fprintln(l.file, plainMsg)
	}

	// 对于FATAL级别，直接退出程序
	if level == FATAL {
		os.Exit(1)
	}
}

// Debug 记录DEBUG级别日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info 记录INFO级别日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warning 记录WARNING级别日志
func (l *Logger) Warning(format string, args ...interface{}) {
	l.log(WARNING, format, args...)
}

// Error 记录ERROR级别日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal 记录FATAL级别日志并退出程序
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

// 提供默认实例的方便方法
func Debug(format string, args ...interface{}) {
	Default.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	Default.Info(format, args...)
}

func Warning(format string, args ...interface{}) {
	Default.Warning(format, args...)
}

func Error(format string, args ...interface{}) {
	Default.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	Default.Fatal(format, args...)
}

// FormatAccessLog 格式化HTTP访问日志
func FormatAccessLog(method, path string, statusCode, bytes int, clientIP string, duration time.Duration) (string, string) {
	// 为状态码选择颜色
	var statusColorFunc func(format string, a ...interface{}) string
	var ok bool
	if statusColorFunc, ok = StatusColor[statusCode]; !ok {
		// 如果没有为特定状态码定义颜色，根据状态码范围选择颜色
		switch {
		case statusCode >= 200 && statusCode < 300:
			statusColorFunc = StatusColor[200]
		case statusCode >= 300 && statusCode < 400:
			statusColorFunc = StatusColor[301]
		case statusCode >= 400 && statusCode < 500:
			statusColorFunc = StatusColor[400]
		default:
			statusColorFunc = StatusColor[500]
		}
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	durationStr := fmt.Sprintf("%.3fms", float64(duration.Microseconds())/1000.0)

	// 带颜色的访问日志（用于控制台）
	methodStr := MethodColor("%s", method)
	pathStr := PathColor("%s", path)
	statusStr := statusColorFunc("%d", statusCode)
	bytesStr := BytesColor("%d", bytes)
	ipStr := ClientIPColor("%s", clientIP)

	coloredLog := fmt.Sprintf("%s | %s | %s | %s | %s | %s",
		methodStr, pathStr, statusStr, bytesStr, ipStr, durationStr)

	// 不带颜色的访问日志（用于文件）
	plainLog := fmt.Sprintf("%s | %s | %s | %d | %s | %s",
		timestamp, method, path, statusCode, clientIP, durationStr)

	return coloredLog, plainLog
}

// AccessLog 记录HTTP访问日志
func (l *Logger) AccessLog(method, path string, statusCode, bytes int, clientIP string, duration time.Duration) {
	coloredLog, plainLog := FormatAccessLog(method, path, statusCode, bytes, clientIP, duration)

	// 输出到控制台（带颜色）
	fmt.Fprintln(l.console, coloredLog)

	// 如果启用了文件日志，输出到文件（不带颜色）
	if l.fileEnabled && l.file != nil {
		fmt.Fprintln(l.file, plainLog)
	}
}

// 默认实例的访问日志方法
func AccessLog(method, path string, statusCode, bytes int, clientIP string, duration time.Duration) {
	Default.AccessLog(method, path, statusCode, bytes, clientIP, duration)
}
