package logger

import (
	"fmt"
	"time"
)

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
