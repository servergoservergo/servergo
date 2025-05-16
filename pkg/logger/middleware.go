package logger

import (
	"time"

	"github.com/gin-gonic/gin"
)

// GinLogger 返回一个Gin中间件，用于记录HTTP请求的访问日志
func GinLogger(logger *Logger) gin.HandlerFunc {
	if logger == nil {
		logger = Default
	}

	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		// 继续处理请求
		c.Next()

		// 计算请求处理时间
		duration := time.Since(start)

		// 获取响应状态码
		statusCode := c.Writer.Status()

		// 获取客户端IP
		clientIP := c.ClientIP()

		// 获取请求方法
		method := c.Request.Method

		// 获取响应体大小
		bodySize := c.Writer.Size()
		if bodySize < 0 {
			bodySize = 0
		}

		// 记录访问日志
		logger.AccessLog(method, path, statusCode, bodySize, clientIP, duration)
	}
}

// DefaultGinLogger 使用默认日志记录器的Gin中间件
func DefaultGinLogger() gin.HandlerFunc {
	return GinLogger(Default)
}
