package logger

import (
	"fmt"
	"os"
	"time"
)

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

// 默认实例的访问日志方法
func AccessLog(method, path string, statusCode, bytes int, clientIP string, duration time.Duration) {
	Default.AccessLog(method, path, statusCode, bytes, clientIP, duration)
}
