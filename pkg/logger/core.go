package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

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

// SetOutput 设置控制台日志输出目标
func (l *Logger) SetOutput(w io.Writer) {
	l.console = w
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level int) {
	l.level = level
}

// GetLevel 获取当前日志级别
func (l *Logger) GetLevel() int {
	return l.level
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
