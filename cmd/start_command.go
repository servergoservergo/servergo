package cmd

import (
	"fmt"
	"os"

	"github.com/CC11001100/servergo/pkg/auth"
	"github.com/CC11001100/servergo/pkg/dirlist"
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/CC11001100/servergo/pkg/server"
	"github.com/CC11001100/servergo/pkg/utils"
	"github.com/spf13/cobra"
)

// startCmd 表示启动服务器的命令
var startCmd = &cobra.Command{
	Use:     "start",
	Aliases: startCmdAliases,
	Short:   i18n.T("cmd.start.short"),
	Long:    i18n.T("cmd.start.long"),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 处理配置
		if err := processStartConfig(cmd); err != nil {
			return err
		}

		// 处理日志配置
		if err := processLogConfig(cmd); err != nil {
			return err
		}

		// 如果用户提供了主题标志但没有值，显示可用主题
		if cmd.Flags().Changed("theme") && theme == "" {
			// 使用dirlist包中定义的有效主题列表
			fmt.Println(i18n.T("cmd.theme.options"))
			for _, t := range dirlist.ValidThemes {
				fmt.Printf("  - %s\n", t)
			}
			fmt.Println("\n" + i18n.T("errors.theme_help"))
			os.Exit(0)
		}

		// 探测可用端口
		actualPort, err := utils.FindAvailablePort(getStartPort())
		if err != nil {
			return fmt.Errorf(i18n.T("error.no_port_available"))
		}

		// 如果使用的不是用户指定的端口，提示用户
		if port > 0 && port != actualPort {
			logger.Warning(i18n.Tf("error.port_unavailable", port))
			logger.Info(i18n.Tf("server.starting", actualPort))
		} else if port == 0 {
			startPort := getStartPort()
			if startPort > 0 && startPort != actualPort {
				// 使用了配置文件中的起始端口，但实际使用的是不同的端口
				logger.Info(i18n.Tf("server.using_config_start_port", startPort))
				logger.Info(i18n.Tf("server.starting", actualPort))
			} else {
				// 随机选择的端口
				logger.Info(i18n.Tf("server.starting", actualPort))
			}
		}

		// 转换认证类型
		var authTypeEnum auth.AuthType
		switch authType {
		case "basic":
			authTypeEnum = auth.BasicAuth
		case "token":
			authTypeEnum = auth.TokenAuth
		case "form":
			authTypeEnum = auth.FormAuth
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

func init() {
	RootCmd.AddCommand(startCmd)

	// 添加标志到start命令
	// 端口默认值为0，表示自动探测可用端口
	startCmd.Flags().IntVarP(&port, "port", "p", 0, i18n.T("flag.port"))
	startCmd.Flags().StringVarP(&dir, "dir", "d", ".", i18n.T("flag.dir"))

	// 所有配置项的命令行标志默认值设为空或false
	// 实际的默认值会从配置文件中读取，如果配置文件中没有才会使用 pkg/config/config.go 中定义的默认值
	startCmd.Flags().BoolVarP(&autoOpen, "open", "o", false, i18n.T("flag.auto_open"))
	startCmd.Flags().BoolVarP(&enableDirListing, "dir-list", "i", false, i18n.T("flag.dir_list"))
	startCmd.Flags().StringVarP(&theme, "theme", "m", "", i18n.T("flag.theme"))

	// 添加认证相关的标志
	startCmd.Flags().StringVarP(&authType, "auth", "a", "none", i18n.T("flag.auth_type"))
	startCmd.Flags().StringVarP(&username, "username", "u", "", i18n.T("flag.username"))
	startCmd.Flags().StringVarP(&password, "password", "w", "", i18n.T("flag.password"))
	startCmd.Flags().StringVarP(&token, "token", "t", "", i18n.T("flag.token"))
	startCmd.Flags().BoolVarP(&enableLoginPage, "login-page", "l", false, i18n.T("flag.login_page"))

	// 添加日志相关的标志
	startCmd.Flags().StringVar(&logLevel, "log-level", "info", i18n.T("flag.log_level"))
	startCmd.Flags().BoolVar(&enableLogPersistence, "enable-log-persistence", false, i18n.T("flag.enable_log_persistence"))
}
