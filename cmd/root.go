package cmd

import (
	"os"

	"github.com/CC11001100/servergo/pkg/config"
	"github.com/CC11001100/servergo/pkg/logger"
	"github.com/spf13/cobra"
)

// 全局配置变量
var (
	port int
	dir  string
)

// RootCmd 表示没有调用子命令时的基础命令
var RootCmd = &cobra.Command{
	Use:   "servergo",
	Short: "快速启动HTTP文件服务器",
	Long: `ServerGo 是一个简单的命令行工具，用于快速启动HTTP文件服务器，
类似于Python的http.server模块，但使用Go实现，提供更好的性能。`,
}

func init() {
	// 在这里我们可以初始化一些全局标志或配置
	// 初始化配置（忽略错误，因为首次运行可能没有配置文件）
	_ = config.InitConfig()

	// 设置默认值
	config.SetDefaults()

	// 初始化日志系统
	// 可以根据需要设置日志级别
	logger.Default.SetLevel(logger.INFO)
}

// Execute 添加所有子命令到根命令并执行
func Execute() {
	// 将在init函数中执行
	cobra.OnInitialize(func() {
		// 所有命令已经注册后执行
		if len(os.Args) == 1 {
			// 如果没有提供子命令，则使用start运行
			args := append([]string{"start"}, os.Args[1:]...)
			RootCmd.SetArgs(args)
		}
	})

	if err := RootCmd.Execute(); err != nil {
		logger.Error("执行命令时发生错误: %v", err)
		os.Exit(1)
	}
}
