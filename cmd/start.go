package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/CC11001100/servergo/pkg/auth"
	"github.com/CC11001100/servergo/pkg/config"
	"github.com/CC11001100/servergo/pkg/dirlist"
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

		// 如果用户提供了主题标志但没有值，显示可用主题
		if cmd.Flags().Changed("theme") && theme == "" {
			// 使用dirlist包中定义的有效主题列表
			fmt.Println("可用主题:")
			for _, t := range dirlist.ValidThemes {
				fmt.Printf("  - %s\n", t)
			}
			fmt.Println("\n示例: --theme dark")
			os.Exit(0)
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
				return fmt.Errorf(i18n.T("auth.basic_credentials_required"))
			}
		case "token":
			authTypeEnum = auth.TokenAuth
			if token == "" {
				return fmt.Errorf(i18n.T("auth.token_required"))
			}
		case "form":
			authTypeEnum = auth.FormAuth
			if username == "" || password == "" {
				return fmt.Errorf(i18n.T("auth.form_credentials_required"))
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
		logger.Warning(i18n.Tf("server.os_not_supported", url))
		return
	}

	if err != nil {
		logger.Error(i18n.Tf("server.browser_error", err, url))
	} else {
		logger.Info(i18n.Tf("server.browser_opened", url))
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
	startCmd.Flags().StringVarP(&authType, "auth", "a", "none", i18n.T("flag.auth_type"))
	startCmd.Flags().StringVarP(&username, "username", "u", "", i18n.T("flag.username"))
	startCmd.Flags().StringVarP(&password, "password", "w", "", i18n.T("flag.password"))
	startCmd.Flags().StringVarP(&token, "token", "t", "", i18n.T("flag.token"))
	startCmd.Flags().BoolVarP(&enableLoginPage, "login-page", "l", false, i18n.T("flag.login_page"))

	// 添加目录浏览相关的标志
	startCmd.Flags().BoolVarP(&enableDirListing, "dir-list", "i", true, i18n.T("config.enable_dir_listing"))
	startCmd.Flags().StringVarP(&theme, "theme", "m", "", i18n.T("config.theme"))
}
