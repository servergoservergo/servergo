package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CC11001100/servergo/pkg/auth"
	"github.com/CC11001100/servergo/pkg/dirlist"
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/logger"
)

// Config 保存文件服务器的配置
// 该结构体包含了服务器启动和运行所需的所有配置参数
type Config struct {
	Port int    // 服务器监听的端口，例如: 8080
	Dir  string // 提供服务的目录路径，例如: "/home/user/files"

	// 认证相关配置
	AuthType        auth.AuthType // 认证类型，可选值: auth.NoAuth, auth.BasicAuth, auth.TokenAuth, auth.FormAuth
	Username        string        // 用户名，用于BasicAuth和FormAuth，例如: "admin"
	Password        string        // 密码，用于BasicAuth和FormAuth，例如: "password123"
	Token           string        // 令牌，用于TokenAuth，例如: "abcdef123456"
	EnableLoginPage bool          // 是否启用登录页面，用于FormAuth，例如: true表示启用

	// 目录浏览相关配置
	EnableDirListing bool   // 是否启用目录列表功能，例如: true表示启用
	Theme            string // 目录列表主题，可选值: "default", "bootstrap", "material" 等
}

// FileServer 表示一个文件服务器实例
// 该结构体封装了服务器的所有状态和功能
type FileServer struct {
	config        Config                   // 服务器配置信息
	absDir        string                   // 服务目录的绝对路径，例如: "/home/user/files"
	engine        *gin.Engine              // Gin引擎实例，用于处理HTTP请求
	authenticator auth.Authenticator       // 认证器实例，用于处理用户认证
	dirTemplate   *dirlist.DirListTemplate // 目录列表模板，用于渲染目录页面
}

// New 创建一个新的文件服务器实例
//
// 参数:
//   - config: 服务器配置，包含端口、目录、认证信息等
//
// 返回值:
//   - *FileServer: 初始化好的文件服务器实例
//   - error: 如有错误，返回对应错误信息，例如目录不存在或无法访问
//
// 使用示例:
// ```
//
//	config := server.Config{
//	    Port: 8080,
//	    Dir: "./files",
//	    AuthType: auth.NoAuth,
//	    EnableDirListing: true,
//	    Theme: "default",
//	}
//
// srv, err := server.New(config)
//
//	if err != nil {
//	    log.Fatalf("创建服务器失败: %v", err)
//	}
//
// ```
func New(config Config) (*FileServer, error) {
	// 获取绝对路径
	absDir, err := filepath.Abs(config.Dir)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.dir_not_exist"), config.Dir)
	}

	// 检查目录是否存在
	info, err := os.Stat(absDir)
	if err != nil {
		return nil, fmt.Errorf(i18n.Tf("error.dir_not_exist", absDir))
	}
	if !info.IsDir() {
		return nil, fmt.Errorf(i18n.Tf("error.not_a_directory", absDir))
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
	} else if !dirlist.IsValidTheme(theme) {
		// 如果提供了无效的主题，记录警告并回退到默认主题
		logger.Warning(fmt.Sprintf("无效的主题名称 '%s'，使用默认主题 '%s'", theme, dirlist.DefaultTheme))
		theme = dirlist.DefaultTheme
	}

	// 创建目录列表模板
	dirTemplate, err := dirlist.NewDirListTemplate(theme)
	if err != nil {
		return nil, fmt.Errorf(i18n.Tf("http.500_template", err))
	}

	// 在这里检查dirTemplate是否初始化成功
	if dirTemplate == nil {
		return nil, fmt.Errorf("无法初始化目录列表模板")
	}

	// 额外记录成功加载的主题信息
	logger.Info(fmt.Sprintf("成功加载目录列表主题: %s", dirTemplate.GetTheme()))

	return &FileServer{
		config:        config,
		absDir:        absDir,
		engine:        engine,
		authenticator: authenticator,
		dirTemplate:   dirTemplate,
	}, nil
}

// Start 启动文件服务器
// 此方法会阻塞执行，直到服务器停止或出错
//
// 返回值:
//   - error: 如果服务器启动失败，返回错误信息
//
// 使用示例:
// ```
// err := srv.Start()
//
//	if err != nil {
//	    log.Fatalf("服务器启动失败: %v", err)
//	}
//
// ```
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
	// /_servergo_assets 路径下的资源会被提供给客户端，如CSS、JS文件
	staticFS := dirlist.GetStaticAssets()
	fs.engine.StaticFS("/_servergo_assets", http.FS(staticFS))

	// 使用NoRoute处理所有未匹配的路由
	fs.engine.NoRoute(fs.handleFileRequest)

	// 打印服务器信息
	logger.Info(i18n.Tf("server.starting", fs.config.Port))
	logger.Info(i18n.Tf("server.serving_dir", fs.absDir))

	// 打印目录列表状态
	if fs.config.EnableDirListing {
		logger.Info(i18n.Tf("server.dir_listing_enabled", fs.dirTemplate.GetTheme()))
	} else {
		logger.Info(i18n.T("server.dir_listing_disabled"))
	}

	// 打印认证信息
	switch fs.authenticator.AuthType() {
	case auth.NoAuth:
		logger.Info(i18n.T("auth.disabled"))
	case auth.BasicAuth:
		logger.Info(i18n.T("auth.basic_enabled"))
		username, password := fs.authenticator.GetCredentials()
		logger.Info("\033[1;32m认证信息:\033[0m")
		logger.Info("\033[1;34m用户名:\033[0m \033[1;33m%s\033[0m", username)
		logger.Info("\033[1;34m密  码:\033[0m \033[1;33m%s\033[0m", password)
	case auth.TokenAuth:
		logger.Info(i18n.T("auth.token_enabled"))
		_, token := fs.authenticator.GetCredentials()
		logger.Info(i18n.Tf("auth.token_access", token))
	case auth.FormAuth:
		logger.Info(i18n.T("auth.form_enabled"))
		if fs.authenticator.LoginPageEnabled() {
			logger.Info(i18n.T("auth.login_page_enabled"))
			username, password := fs.authenticator.GetCredentials()
			logger.Info("\033[1;32m认证信息:\033[0m")
			logger.Info("\033[1;34m登录地址:\033[0m \033[1;36mhttp://localhost:%d/auth/login\033[0m", fs.config.Port)
			logger.Info("\033[1;34m用户名:\033[0m \033[1;33m%s\033[0m", username)
			logger.Info("\033[1;34m密  码:\033[0m \033[1;33m%s\033[0m", password)
		}
	}

	// 提示用户如何停止服务器
	logger.Info(i18n.T("server.press_ctrl_c"))

	// 启动服务器
	hostAddr := ":" + strconv.Itoa(fs.config.Port)
	return fs.engine.Run(hostAddr)
}

// GetAbsDir 获取文件服务器的绝对路径
//
// 返回值:
//   - string: 服务器提供服务的目录绝对路径
//
// 使用示例:
// ```
// absPath := srv.GetAbsDir()
// fmt.Printf("服务器的绝对路径: %s\n", absPath)
// ```
func (fs *FileServer) GetAbsDir() string {
	return fs.absDir
}
