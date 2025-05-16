package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CC11001100/servergo/pkg/logger"
)

// Config 保存文件服务器的配置
type Config struct {
	Port int    // 服务器监听的端口
	Dir  string // 提供服务的目录路径
}

// FileServer 表示一个文件服务器实例
type FileServer struct {
	config Config
	absDir string
	engine *gin.Engine
}

// New 创建一个新的文件服务器实例
func New(config Config) (*FileServer, error) {
	// 获取绝对路径
	absDir, err := filepath.Abs(config.Dir)
	if err != nil {
		return nil, fmt.Errorf("无法获取绝对路径: %v", err)
	}

	// 检查目录是否存在
	info, err := os.Stat(absDir)
	if err != nil {
		return nil, fmt.Errorf("无法访问目录 %s: %v", absDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s 不是一个目录", absDir)
	}

	// 设置Gin为生产模式，避免debug信息
	gin.SetMode(gin.ReleaseMode)

	// 创建一个默认的Gin引擎
	engine := gin.New()

	// 使用自定义的日志中间件和恢复中间件
	engine.Use(logger.DefaultGinLogger(), gin.Recovery())

	return &FileServer{
		config: config,
		absDir: absDir,
		engine: engine,
	}, nil
}

// Start 启动文件服务器
func (fs *FileServer) Start() error {
	// 设置静态文件服务
	fs.engine.Static("/", fs.absDir)

	// 添加一个自定义路由，记录文件访问信息
	fs.engine.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// 继续处理请求
		c.Next()

		// 请求结束后，记录详细信息
		logger.Info("文件访问: %s (耗时: %.3fms)",
			path, float64(time.Since(start).Microseconds())/1000.0)
	})

	// 启动服务器
	hostAddr := ":" + strconv.Itoa(fs.config.Port)
	logger.Info("启动文件服务器在 http://localhost:%d", fs.config.Port)
	logger.Info("提供目录: %s", fs.absDir)
	logger.Info("按 Ctrl+C 停止服务器")

	return fs.engine.Run(hostAddr)
}

// GetAbsDir 获取文件服务器的绝对路径
func (fs *FileServer) GetAbsDir() string {
	return fs.absDir
}
