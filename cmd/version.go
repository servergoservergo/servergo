package cmd

import (
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/spf13/cobra"
)

// 版本信息
var (
	Version   = "0.1.0"
	BuildDate = "未知"
	GitCommit = "未知"
)

// versionCmd 显示当前版本信息和项目信息
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示ServerGo的版本信息",
	Long:  `显示ServerGo的详细版本信息、构建信息、作者信息以及问题反馈渠道等。`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("ServerGo 版本信息")
		logger.Info("==============")
		logger.Info("版本: %s", Version)
		logger.Info("构建日期: %s", BuildDate)
		logger.Info("Git提交哈希: %s", GitCommit)
		logger.Info("")
		logger.Info("ServerGo 是一个快速的HTTP文件服务器工具")
		logger.Info("类似于Python的http.server模块，但提供更好的性能和更多功能")
		logger.Info("")
		logger.Info("作者信息")
		logger.Info("--------")
		logger.Info("作者: CC11001100")
		logger.Info("GitHub: https://github.com/CC11001100")
		logger.Info("")
		logger.Info("反馈问题")
		logger.Info("--------")
		logger.Info("如果您发现任何问题或有功能请求，请访问:")
		logger.Info("https://github.com/CC11001100/servergo/issues")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
