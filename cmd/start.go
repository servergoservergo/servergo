package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/CC11001100/servergo/pkg/config"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/CC11001100/servergo/pkg/server"
	"github.com/CC11001100/servergo/pkg/utils"
	"github.com/spf13/cobra"
)

// 命令行标志
var (
	// 是否自动打开浏览器（命令行标志）
	autoOpen bool
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
	Short:   "启动HTTP文件服务器",
	Long:    `启动一个HTTP文件服务器，为指定目录提供文件访问服务。目录路径可以是绝对路径或相对路径，如果是相对路径，将自动转换为绝对路径。如果不指定端口或指定的端口已被占用，将自动探测一个可用的端口。`,
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
			return fmt.Errorf("无法找到可用端口: %v", err)
		}

		// 如果使用的不是用户指定的端口，提示用户
		if port > 0 && port != actualPort {
			logger.Warning("指定的端口 %d 已被占用，使用可用端口 %d 代替", port, actualPort)
		} else if port == 0 {
			logger.Info("未指定端口，自动使用可用端口 %d", actualPort)
		}

		// 创建服务器配置
		serverConfig := server.Config{
			Port: actualPort,
			Dir:  dir,
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
	startCmd.Flags().IntVarP(&port, "port", "p", 0, "服务器要监听的端口（默认自动探测，指定端口被占用时也会自动探测）")
	startCmd.Flags().StringVarP(&dir, "dir", "d", ".", "要提供服务的目录路径（可以是绝对路径或相对路径，默认当前目录）")
	startCmd.Flags().BoolVarP(&autoOpen, "open", "o", true, "启动服务器后是否自动打开浏览器（默认使用配置中的设置）")
}
