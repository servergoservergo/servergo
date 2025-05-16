package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
)

// 日志级别常量
const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
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
	output io.Writer
	level  int
}

// 创建一个新的Logger实例
func New(level int) *Logger {
	return &Logger{
		output: os.Stdout,
		level:  level,
	}
}

// 默认日志实例
var Default = New(INFO)

// SetOutput 设置日志输出目标
func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level int) {
	l.level = level
}

// formatMessage 格式化日志消息
func (l *Logger) formatMessage(level int, format string, args ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	levelStr := levelNames[level]
	colorFunc := levelColors[level]

	// 应用颜色到日志级别
	coloredLevel := colorFunc("[%s]", levelStr)
	timeStr := TimeColor("%s", timestamp)

	// 格式化消息内容
	var message string
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format
	}

	return fmt.Sprintf("%s %s %s", timeStr, coloredLevel, message)
}

// log 内部日志方法
func (l *Logger) log(level int, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	message := l.formatMessage(level, format, args...)
	fmt.Fprintln(l.output, message)

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
func FormatAccessLog(method, path string, statusCode, bytes int, clientIP string, duration time.Duration) string {
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

	// 格式化各部分
	methodStr := MethodColor("%s", method)
	pathStr := PathColor("%s", path)
	statusStr := statusColorFunc("%d", statusCode)
	bytesStr := BytesColor("%d", bytes)
	ipStr := ClientIPColor("%s", clientIP)
	durationStr := fmt.Sprintf("%.3fms", float64(duration.Microseconds())/1000.0)

	// 返回格式化后的访问日志
	return fmt.Sprintf("%s | %s | %s | %s | %s | %s",
		methodStr, pathStr, statusStr, bytesStr, ipStr, durationStr)
}

// AccessLog 记录HTTP访问日志
func (l *Logger) AccessLog(method, path string, statusCode, bytes int, clientIP string, duration time.Duration) {
	logStr := FormatAccessLog(method, path, statusCode, bytes, clientIP, duration)
	fmt.Fprintln(l.output, logStr)
}

// 默认实例的访问日志方法
func AccessLog(method, path string, statusCode, bytes int, clientIP string, duration time.Duration) {
	Default.AccessLog(method, path, statusCode, bytes, clientIP, duration)
}
