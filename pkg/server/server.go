package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CC11001100/servergo/pkg/auth"
	"github.com/CC11001100/servergo/pkg/dirlist"
	"github.com/CC11001100/servergo/pkg/logger"
)

// Config 保存文件服务器的配置
type Config struct {
	Port int    // 服务器监听的端口
	Dir  string // 提供服务的目录路径
	// 认证相关配置
	AuthType        auth.AuthType // 认证类型
	Username        string        // 用户名
	Password        string        // 密码
	Token           string        // 令牌
	EnableLoginPage bool          // 是否启用登录页面
	// 目录浏览相关配置
	EnableDirListing bool   // 是否启用目录列表功能
	Theme            string // 目录列表主题
}

// FileServer 表示一个文件服务器实例
type FileServer struct {
	config        Config
	absDir        string
	engine        *gin.Engine
	authenticator auth.Authenticator
	dirTemplate   *dirlist.DirListTemplate
}

// New 创建一个新的文件服务器实例
func New(config Config) (*FileServer, error) {
	// 获取绝对路径
	absDir, err := filepath.Abs(config.Dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %v", err)
	}

	// 检查目录是否存在
	info, err := os.Stat(absDir)
	if err != nil {
		return nil, fmt.Errorf("failed to access directory %s: %v", absDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", absDir)
	}

	// 设置Gin为生产模式，避免debug信息
	gin.SetMode(gin.ReleaseMode)

	// 创建一个默认的Gin引擎
	engine := gin.New()

	// 使用自定义的日志中间件和恢复中间件
	engine.Use(logger.DefaultGinLogger(), gin.Recovery())

	// 创建认证器
	authenticator := auth.NewAuthenticator(auth.Config{
		Type:            config.AuthType,
		Username:        config.Username,
		Password:        config.Password,
		Token:           config.Token,
		EnableLoginPage: config.EnableLoginPage,
	})

	// 如果未设置主题，使用默认主题
	theme := config.Theme
	if theme == "" {
		theme = dirlist.DefaultTheme
	}

	// 创建目录列表模板
	dirTemplate, err := dirlist.NewDirListTemplate(theme)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize directory listing template: %v", err)
	}

	return &FileServer{
		config:        config,
		absDir:        absDir,
		engine:        engine,
		authenticator: authenticator,
		dirTemplate:   dirTemplate,
	}, nil
}

// Start 启动文件服务器
func (fs *FileServer) Start() error {
	// 如果是表单认证并且启用了登录页面，设置表单认证的路由
	if fs.authenticator.AuthType() == auth.FormAuth && fs.authenticator.LoginPageEnabled() {
		formAuth, ok := fs.authenticator.(*auth.FormAuthenticator)
		if ok {
			formAuth.SetupRoutes(fs.engine)
		}
	}

	// 添加认证中间件
	fs.engine.Use(fs.authenticator.Middleware())

	// 添加一个自定义路由，记录文件访问信息
	fs.engine.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// 继续处理请求
		c.Next()

		// 请求结束后，记录详细信息
		status := c.Writer.Status()
		size := c.Writer.Size()
		clientIP := c.ClientIP()
		latency := time.Since(start)

		logger.Info("GET | %s | %d | %d | %s | %.3fms",
			path, status, size, clientIP, float64(latency.Microseconds())/1000.0)
	})

	// 提供模板静态资源，使用特定路由前缀
	staticFS := dirlist.GetStaticAssets()
	fs.engine.StaticFS("/_servergo_assets", http.FS(staticFS))

	// 使用NoRoute处理所有未匹配的路由
	fs.engine.NoRoute(fs.handleFileRequest)

	// 启动服务器
	hostAddr := ":" + strconv.Itoa(fs.config.Port)
	logger.Info("Starting file server at http://localhost:%d", fs.config.Port)
	logger.Info("Serving directory: %s", fs.absDir)

	// 显示功能信息
	if fs.config.EnableDirListing {
		logger.Info("Directory listing enabled (theme: %s)", fs.dirTemplate.GetTheme())
	} else {
		logger.Info("Directory listing disabled")
	}

	// 显示认证信息
	switch fs.authenticator.AuthType() {
	case auth.NoAuth:
		logger.Info("Authentication disabled")
	case auth.BasicAuth:
		logger.Info("Basic authentication enabled")
	case auth.TokenAuth:
		logger.Info("Token authentication enabled")
		logger.Info("Access with URL parameter ?token=%s or Authorization header", fs.config.Token)
	case auth.FormAuth:
		logger.Info("Form authentication enabled")
		if fs.authenticator.LoginPageEnabled() {
			logger.Info("Login page enabled, visit /auth/login to login")
		}
	}

	logger.Info("Press Ctrl+C to stop the server")

	return fs.engine.Run(hostAddr)
}

// GetAbsDir 获取文件服务器的绝对路径
func (fs *FileServer) GetAbsDir() string {
	return fs.absDir
}

// handleFileRequest 处理文件请求
func (fs *FileServer) handleFileRequest(c *gin.Context) {
	// 获取请求路径
	reqPath := c.Request.URL.Path

	// 如果请求的是内部静态资源，跳过处理
	if strings.HasPrefix(reqPath, "/_servergo_assets/") {
		c.Next()
		return
	}

	// 处理路径中可能的URL编码
	reqPath = strings.Replace(reqPath, "%20", " ", -1)

	// 确保路径不会超出根目录
	cleanPath := filepath.Clean(reqPath)
	if !strings.HasPrefix(cleanPath, "/") {
		cleanPath = "/" + cleanPath
	}

	// 将请求路径转换为服务器文件系统上的实际路径
	fullPath := filepath.Join(fs.absDir, cleanPath)

	// 获取文件状态
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		// 文件不存在，返回404
		c.String(http.StatusNotFound, "404 Not Found: %s", reqPath)
		return
	}

	// 如果是目录，检查是否启用了目录列表功能
	if fileInfo.IsDir() {
		// 检查该目录下是否有index.html文件
		indexPath := filepath.Join(fullPath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			// 如果存在index.html，则提供该文件
			c.File(indexPath)
			return
		}

		// 如果启用了目录列表功能，则显示目录内容
		if fs.config.EnableDirListing {
			fs.renderDirectoryListing(c, fullPath, reqPath)
			return
		}

		// 未启用目录列表功能，返回403禁止访问
		c.String(http.StatusForbidden, "403 Forbidden: Directory listing disabled")
		return
	}

	// 如果是文件，则提供该文件
	c.File(fullPath)
}

// renderDirectoryListing 渲染目录列表页面
func (fs *FileServer) renderDirectoryListing(c *gin.Context, fullPath, reqPath string) {
	// 读取目录内容
	files, err := os.ReadDir(fullPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read directory content: %v", err)
		return
	}

	// 创建文件项列表
	items := make([]dirlist.FileItem, 0, len(files))
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		// 构建文件路径
		itemPath := filepath.Join(reqPath, file.Name())
		if !strings.HasPrefix(itemPath, "/") {
			itemPath = "/" + itemPath
		}

		// 格式化大小
		var sizeStr string
		var sizeBytes int64
		if file.IsDir() {
			sizeStr = "-"
			sizeBytes = 0
		} else {
			sizeBytes = info.Size()
			sizeStr = formatSize(sizeBytes)
		}

		// 添加到列表
		items = append(items, dirlist.FileItem{
			Name:         file.Name(),
			IsDir:        file.IsDir(),
			Size:         sizeStr,
			SizeBytes:    sizeBytes,
			LastModified: info.ModTime().Format("2006-01-02 15:04:05"),
			Path:         itemPath,
		})
	}

	// 按照目录在前，文件在后的方式排序
	sort.Slice(items, func(i, j int) bool {
		// 如果一个是目录一个不是，目录在前
		if items[i].IsDir != items[j].IsDir {
			return items[i].IsDir
		}
		// 否则按名称字母顺序排序
		return items[i].Name < items[j].Name
	})

	// 计算上级目录路径
	var parentDir string
	if reqPath != "/" {
		parentDir = filepath.Dir(reqPath)
		if parentDir == "." {
			parentDir = "/"
		}
		if !strings.HasPrefix(parentDir, "/") {
			parentDir = "/" + parentDir
		}
		if parentDir != "/" && strings.HasSuffix(parentDir, "/") {
			parentDir = parentDir[:len(parentDir)-1]
		}
	}

	// 准备模板数据
	data := dirlist.TemplateData{
		DirPath:     reqPath,
		Items:       items,
		ParentDir:   parentDir,
		CurrentTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	// 渲染模板
	html, err := fs.dirTemplate.Render(data)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template rendering error: %v", err)
		return
	}

	// 使用模板提供的内容类型，支持HTML和JSON格式
	c.Header("Content-Type", fs.dirTemplate.GetContentType())
	c.String(http.StatusOK, html)
}

// formatSize 格式化文件大小
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
