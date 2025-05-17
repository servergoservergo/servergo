package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/CC11001100/servergo/pkg/auth"
	"github.com/CC11001100/servergo/pkg/config"
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/CC11001100/servergo/pkg/server"
	"github.com/CC11001100/servergo/pkg/utils"
	"github.com/spf13/cobra"
)

// 命令行标志
var (
	// 是否自动打开浏览器（命令行标志）
	autoOpen bool

	// 认证相关标志
	authType        string // 认证类型：none, basic, token, form
	username        string // 用户名
	password        string // 密码
	token           string // 令牌
	enableLoginPage bool   // 是否启用登录页面

	// 目录浏览相关标志
	enableDirListing bool   // 是否启用目录列表功能
	theme            string // 目录列表主题
)

// 别名列表 - 预留位置供后续扩展
var startCmdAliases = []string{
	"run",
	"serve",
	"launch",
	// 这里可以继续添加更多别名
}

// startCmd 表示启动服务器的命令
var startCmd = &cobra.Command{
	Use:     "start",
	Aliases: startCmdAliases,
	Short:   i18n.T("cmd.start.short"),
	Long:    i18n.T("cmd.start.long"),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 读取配置中的默认值
		cfg := config.GetConfig()

		// 如果命令行未指定是否自动打开浏览器，使用配置中的设置
		if !cmd.Flags().Changed("open") {
			autoOpen = cfg.AutoOpen
		}

		// 探测可用端口
		// 如果port=0，表示自动探测
		// 如果port>0，先检查指定端口是否可用，不可用则自动探测
		actualPort, err := utils.FindAvailablePort(port)
		if err != nil {
			return fmt.Errorf(i18n.T("error.no_port_available"))
		}

		// 如果使用的不是用户指定的端口，提示用户
		if port > 0 && port != actualPort {
			logger.Warning(i18n.Tf("error.port_unavailable", port))
			logger.Info(i18n.Tf("server.starting", actualPort))
		} else if port == 0 {
			logger.Info(i18n.Tf("server.starting", actualPort))
		}

		// 转换认证类型
		var authTypeEnum auth.AuthType
		switch authType {
		case "basic":
			authTypeEnum = auth.BasicAuth
			if username == "" || password == "" {
				return fmt.Errorf("使用Basic认证时必须同时提供用户名和密码")
			}
		case "token":
			authTypeEnum = auth.TokenAuth
			if token == "" {
				return fmt.Errorf("使用Token认证时必须提供令牌")
			}
		case "form":
			authTypeEnum = auth.FormAuth
			if username == "" || password == "" {
				return fmt.Errorf("使用Form认证时必须同时提供用户名和密码")
			}
		default:
			authTypeEnum = auth.NoAuth
		}

		// 创建服务器配置
		serverConfig := server.Config{
			Port:             actualPort,
			Dir:              dir,
			AuthType:         authTypeEnum,
			Username:         username,
			Password:         password,
			Token:            token,
			EnableLoginPage:  enableLoginPage,
			EnableDirListing: enableDirListing,
			Theme:            theme,
		}

		// 创建并启动文件服务器
		srv, err := server.New(serverConfig)
		if err != nil {
			return err
		}

		// 如果配置为自动打开浏览器，则在启动服务器后打开
		if autoOpen {
			// 在新的goroutine中启动浏览器，避免阻塞服务器启动
			serverURL := fmt.Sprintf("http://localhost:%d", actualPort)
			go openBrowser(serverURL)
		}

		return srv.Start()
	},
}

// 打开系统默认浏览器访问URL
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		logger.Warning("不支持的操作系统，无法自动打开浏览器，请手动访问: %s", url)
		return
	}

	if err != nil {
		logger.Error("无法打开浏览器: %v\n请手动访问: %s", err, url)
	} else {
		logger.Info("已在浏览器中打开 %s", url)
	}
}

func init() {
	RootCmd.AddCommand(startCmd)

	// 添加标志到start命令
	// 端口默认值为0，表示自动探测可用端口
	startCmd.Flags().IntVarP(&port, "port", "p", 0, i18n.T("flag.port"))
	startCmd.Flags().StringVarP(&dir, "dir", "d", ".", i18n.T("flag.dir"))
	startCmd.Flags().BoolVarP(&autoOpen, "open", "o", true, i18n.T("flag.auto_open"))

	// 添加认证相关的标志
	startCmd.Flags().StringVarP(&authType, "auth", "a", "none", "认证类型：none(不认证), basic(HTTP基本认证), token(令牌认证), form(表单认证)")
	startCmd.Flags().StringVarP(&username, "username", "u", "", "用于basic或form认证的用户名")
	startCmd.Flags().StringVarP(&password, "password", "w", "", "用于basic或form认证的密码")
	startCmd.Flags().StringVarP(&token, "token", "t", "", "用于token认证的令牌")
	startCmd.Flags().BoolVarP(&enableLoginPage, "login-page", "l", false, "是否启用登录页面（仅适用于form认证）")

	// 添加目录浏览相关的标志
	startCmd.Flags().BoolVarP(&enableDirListing, "dir-list", "i", true, i18n.T("config.enable_dir_listing"))
	startCmd.Flags().StringVarP(&theme, "theme", "m", "", i18n.T("config.theme"))
}
