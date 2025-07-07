package logger

import (
	"io"

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
