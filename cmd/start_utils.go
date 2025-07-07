package cmd

import (
	"os/exec"
	"runtime"

	"github.com/CC11001100/servergo/pkg/config"
	"github.com/CC11001100/servergo/pkg/i18n"
	"github.com/CC11001100/servergo/pkg/logger"
)

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

// getStartPort 获取起始端口，优先使用命令行参数值，否则使用配置文件中的值
func getStartPort() int {
	// 如果命令行指定了端口，使用命令行指定的端口
	if port > 0 {
		return port
	}

	// 否则从配置文件获取起始端口
	cfg := config.GetConfig()
	// 如果配置了起始端口，则使用它
	if cfg.StartPort > 0 {
		return cfg.StartPort
	}

	// 如果都没有指定，则返回0表示随机选择
	return 0
}
