package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/CC11001100/servergo/pkg/auth"
	"github.com/CC11001100/servergo/pkg/dirlist"
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/logger"
)

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
