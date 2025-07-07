package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/CC11001100/servergo/pkg/auth"
	"github.com/CC11001100/servergo/pkg/dirlist"
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/gin-gonic/gin"
)

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
